package cmd

import (
	"fmt"
	"math/big"

	"github.com/opendasom/bit/internal/chain"
	compactcid "github.com/opendasom/bit/internal/cid"
	"github.com/opendasom/bit/internal/config"
	"github.com/opendasom/bit/internal/git"
	"github.com/opendasom/bit/internal/ipfs"
	"github.com/opendasom/bit/internal/manifest"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push <remote>",
	Short: "Push current branch to remote",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		remoteName := args[0]

		// ── 사전 검증 ──────────────────────────────────────────

		// 1. config 읽기
		cfg, err := config.Load(".")
		if err != nil {
			return fmt.Errorf("config 읽기 실패 (.bit/config.json): %w", err)
		}

		remote, ok := cfg.Remotes[remoteName]
		if !ok {
			return fmt.Errorf("remote '%s' 를 찾을 수 없습니다", remoteName)
		}
		repoID := new(big.Int).SetUint64(remote.RepoID)

		// 2. 체인 연결 확인
		chainClient, err := chain.NewClient(cfg.RPCURL, cfg.ContractAddress, cfg.PrivateKey)
		if err != nil {
			return fmt.Errorf("체인 연결 실패: %w", err)
		}

		// 3. git 정보 읽기
		branch, err := git.CurrentBranch(".")
		if err != nil {
			return fmt.Errorf("git 정보 읽기 실패: %w", err)
		}

		// 4. 체인에서 현재 브랜치 커밋 조회 (compare-and-swap용)
		expectedOldCommit, err := chainClient.GetBranchCommit(repoID, branch)
		if err != nil {
			return fmt.Errorf("브랜치 헤드 조회 실패: %w", err)
		}
		baseCommit := ""
		if !chain.IsZeroBytes20(expectedOldCommit) {
			baseCommit = chain.Bytes20ToGitHash(expectedOldCommit)
		}

		commits, err := git.CommitsAfter(".", baseCommit)
		if err != nil {
			return fmt.Errorf("push 대상 커밋 계산 실패: %w", err)
		}
		if len(commits) == 0 {
			fmt.Printf("%s is already up to date\n", branch)
			return nil
		}

		fmt.Printf("브랜치: %s, 업로드 커밋: %d개\n", branch, len(commits))
		fmt.Println("사전 검증 완료, diff 업로드 시작...")

		// ── 업로드 ─────────────────────────────────────────────

		ipfsClient := ipfs.NewClient(cfg.IPFSURL)
		for _, commit := range commits {
			info, err := git.ReadCommit(".", commit)
			if err != nil {
				return fmt.Errorf("커밋 메타데이터 읽기 실패 (%s): %w", commit, err)
			}
			if len(info.ParentHashes) > 1 {
				return fmt.Errorf("merge commit diff push is not supported yet: %s", info.Hash)
			}

			diff, err := git.ExtractCommitDiff(".", info)
			if err != nil {
				return fmt.Errorf("커밋 diff 추출 실패 (%s): %w", info.Hash, err)
			}
			diffCID, err := ipfsClient.Upload(diff)
			if err != nil {
				return fmt.Errorf("diff IPFS 업로드 실패 (%s): %w", info.Hash, err)
			}
			diffDigest, err := compactcid.DigestFromCIDV0(diffCID)
			if err != nil {
				return fmt.Errorf("diff CID 압축 실패 (%s): %w", diffCID, err)
			}

			m := git.ManifestForCommit(branch, info, diffCID)
			manifestData, err := manifest.Encode(m)
			if err != nil {
				return fmt.Errorf("manifest 생성 실패 (%s): %w", info.Hash, err)
			}
			manifestCID, err := ipfsClient.Upload(manifestData)
			if err != nil {
				return fmt.Errorf("manifest IPFS 업로드 실패 (%s): %w", info.Hash, err)
			}
			manifestDigest, err := compactcid.DigestFromCIDV0(manifestCID)
			if err != nil {
				return fmt.Errorf("manifest CID 압축 실패 (%s): %w", manifestCID, err)
			}

			commitHash, err := chain.GitHashToBytes20(info.Hash)
			if err != nil {
				return fmt.Errorf("커밋 해시 변환 실패 (%s): %w", info.Hash, err)
			}
			treeHash, err := chain.GitHashToBytes20(info.TreeHash)
			if err != nil {
				return fmt.Errorf("트리 해시 변환 실패 (%s): %w", info.TreeHash, err)
			}
			parentHashes := make([][20]byte, 0, len(info.ParentHashes))
			for _, parent := range info.ParentHashes {
				parentHash, err := chain.GitHashToBytes20(parent)
				if err != nil {
					return fmt.Errorf("부모 커밋 해시 변환 실패 (%s): %w", parent, err)
				}
				parentHashes = append(parentHashes, parentHash)
			}

			if err := chainClient.RecordCommit(
				repoID,
				branch,
				expectedOldCommit,
				commitHash,
				treeHash,
				parentHashes,
				manifestDigest,
				diffDigest,
			); err != nil {
				return fmt.Errorf("체인 커밋 기록 실패 (다른 사람이 먼저 push했을 수 있습니다, commit %s): %w", info.Hash, err)
			}

			expectedOldCommit = commitHash
			fmt.Printf("commit push 완료: %s diff=%s manifest=%s\n", info.Hash[:8], diffCID, manifestCID)
		}

		fmt.Printf("push 완료: %s → %s\n", branch, commits[len(commits)-1][:8])
		return nil
	},
}
