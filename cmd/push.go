package cmd

import (
	"fmt"
	"math/big"

	"github.com/hsh-719/bit/internal/chain"
	"github.com/hsh-719/bit/internal/config"
	"github.com/hsh-719/bit/internal/git"
	"github.com/hsh-719/bit/internal/ipfs"
	"github.com/hsh-719/bit/internal/manifest"
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
		headInfo, err := git.ReadHead(".")
		if err != nil {
			return fmt.Errorf("git 정보 읽기 실패: %w", err)
		}

		// 4. 체인에서 현재 브랜치 헤드 조회 (compare-and-swap용)
		expectedOldHead, err := chainClient.GetBranchHead(repoID, headInfo.Branch)
		if err != nil {
			return fmt.Errorf("브랜치 헤드 조회 실패: %w", err)
		}

		fmt.Printf("브랜치: %s, 커밋: %s\n", headInfo.Branch, headInfo.CommitHash[:8])
		fmt.Println("사전 검증 완료, 업로드 시작...")

		// ── 업로드 ─────────────────────────────────────────────

		// 5. git bundle 추출 (전체 코드를 바이트로)
		bundle, err := git.ExtractBundle(".")
		if err != nil {
			return fmt.Errorf("git bundle 추출 실패: %w", err)
		}

		// 6. bundle을 IPFS에 업로드 → bundleCID
		ipfsClient := ipfs.NewClient(cfg.IPFSURL)
		bundleCID, err := ipfsClient.Upload(bundle)
		if err != nil {
			return fmt.Errorf("IPFS 업로드 실패: %w", err)
		}
		fmt.Printf("bundle 업로드 완료 (CID: %s)\n", bundleCID)

		// 7. manifest JSON 생성 후 IPFS에 업로드 → manifestCID
		m := &manifest.Manifest{
			BundleCID:      bundleCID,
			GitCommit:      headInfo.CommitHash,
			PreviousCommit: headInfo.ParentHash,
		}
		manifestData, err := manifest.Encode(m)
		if err != nil {
			return fmt.Errorf("manifest 생성 실패: %w", err)
		}
		manifestCID, err := ipfsClient.Upload(manifestData)
		if err != nil {
			return fmt.Errorf("manifest IPFS 업로드 실패: %w", err)
		}
		fmt.Printf("manifest 업로드 완료 (CID: %s)\n", manifestCID)

		// 8. 체인에 새 manifest CID로 브랜치 헤드 갱신
		if err := chainClient.UpdateBranch(
			repoID,
			headInfo.Branch,
			expectedOldHead,
			manifestCID,
			headInfo.CommitHash,
			headInfo.ParentHash,
		); err != nil {
			return fmt.Errorf("체인 업데이트 실패 (다른 사람이 먼저 push했을 수 있습니다): %w", err)
		}

		fmt.Printf("push 완료: %s → %s\n", headInfo.Branch, manifestCID)
		return nil
	},
}
