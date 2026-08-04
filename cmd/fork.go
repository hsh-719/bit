package cmd

import (
	"fmt"
	"math/big"
	"os"
	"os/exec"

	"github.com/opendasom/bit/internal/chain"
	compactcid "github.com/opendasom/bit/internal/cid"
	"github.com/opendasom/bit/internal/config"
	"github.com/opendasom/bit/internal/git"
	"github.com/opendasom/bit/internal/ipfs"
	"github.com/opendasom/bit/internal/manifest"
	"github.com/spf13/cobra"
)

var forkCmd = &cobra.Command{
	Use:   "fork <bitURL>",
	Short: "Fork a remote repository into a new local repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceURL := args[0]

		rpcURL, _ := cmd.Flags().GetString("rpc")
		contractAddress, _ := cmd.Flags().GetString("contract")
		privateKey, _ := cmd.Flags().GetString("key")
		ipfsURL, _ := cmd.Flags().GetString("ipfs")
		branch, _ := cmd.Flags().GetString("branch")

		// 1. .git 없으면 자동으로 git init
		if _, err := os.Stat(".git"); os.IsNotExist(err) {
			fmt.Println("git 저장소가 없습니다. git init을 실행합니다...")
			initCmd := exec.Command("git", "init")
			initCmd.Stdout = os.Stdout
			initCmd.Stderr = os.Stderr
			if err := initCmd.Run(); err != nil {
				return fmt.Errorf("git init 실패: %w", err)
			}
		}

		// 2. 플래그 미입력 시 기존 .bit/config.json에서 읽기
		if rpcURL == "" || contractAddress == "" || privateKey == "" {
			cfg, err := config.Load(".")
			if err != nil {
				return fmt.Errorf("--rpc, --contract, --key 플래그가 필요합니다 (또는 .bit/config.json이 있어야 합니다): %w", err)
			}
			if rpcURL == "" {
				rpcURL = cfg.RPCURL
			}
			if contractAddress == "" {
				contractAddress = cfg.ContractAddress
			}
			if privateKey == "" {
				privateKey = cfg.PrivateKey
			}
			if ipfsURL == "" {
				ipfsURL = cfg.IPFSURL
			}
		}
		if ipfsURL == "" {
			ipfsURL = "http://localhost:5001"
		}

		// 3. source repoId 파싱
		sourceRepoID, err := parseRepoID(sourceURL)
		if err != nil {
			return fmt.Errorf("bitURL 파싱 실패: %w", err)
		}

		// 4. 체인 연결
		chainClient, err := chain.NewClient(rpcURL, contractAddress, privateKey)
		if err != nil {
			return fmt.Errorf("체인 연결 실패: %w", err)
		}

		// 5. source 브랜치 히스토리 조회
		srcRepoID := new(big.Int).SetUint64(sourceRepoID)
		historyLen, err := chainClient.GetBranchHistoryLength(srcRepoID, branch)
		if err != nil {
			return fmt.Errorf("source 브랜치 히스토리 조회 실패: %w", err)
		}
		if historyLen.Sign() == 0 {
			return fmt.Errorf("source 브랜치 '%s'에 커밋이 없습니다", branch)
		}
		records, err := loadBranchRecords(chainClient, srcRepoID, branch, historyLen.Int64())
		if err != nil {
			return err
		}
		fmt.Printf("source 커밋 %d개 발견\n", len(records))

		// 6. B 저장소 생성 (체인에 새 repoId 발급)
		fmt.Println("fork 저장소를 체인에 생성 중...")
		forkRepoID, err := chainClient.CreateRepo("")
		if err != nil {
			return fmt.Errorf("fork 저장소 생성 실패: %w", err)
		}
		fmt.Printf("fork 저장소 생성 완료 (repoId: %s)\n", forkRepoID.String())

		// 7. 각 커밋: IPFS 다운로드 → 로컬 git 복원 → B 체인에 기록
		ipfsClient := ipfs.NewClient(ipfsURL)
		expectedOldCommit := [20]byte{}

		for i, record := range records {
			manifestCID := compactcid.CIDV0FromDigest(record.ManifestDigest)
			diffCID := compactcid.CIDV0FromDigest(record.DiffDigest)

			// manifest 다운로드
			manifestData, err := ipfsClient.Download(manifestCID)
			if err != nil {
				return fmt.Errorf("manifest 다운로드 실패 (커밋 %d, %s): %w", i+1, manifestCID, err)
			}
			m, err := manifest.Decode(manifestData)
			if err != nil {
				return fmt.Errorf("manifest 파싱 실패 (%s): %w", manifestCID, err)
			}

			// 검증
			expectedCommit := chain.Bytes20ToGitHash(record.CommitHash)
			if m.GitCommit != expectedCommit {
				return fmt.Errorf("manifest commit mismatch: got %s, want %s", m.GitCommit, expectedCommit)
			}
			if m.DiffCID != diffCID {
				return fmt.Errorf("manifest diff CID mismatch for %s", expectedCommit)
			}

			// diff 다운로드
			diff, err := ipfsClient.Download(m.DiffCID)
			if err != nil {
				return fmt.Errorf("diff 다운로드 실패 (%s): %w", m.DiffCID, err)
			}

			// 로컬 git에 커밋 복원
			if err := git.ApplyCommitDiff(".", m, diff); err != nil {
				return fmt.Errorf("commit diff 적용 실패 (%s): %w", expectedCommit, err)
			}

			// parents 변환 (manifest.ParentCommits → [][20]byte)
			parentHashes := make([][20]byte, 0, len(m.ParentCommits))
			for _, parent := range m.ParentCommits {
				h, err := chain.GitHashToBytes20(parent)
				if err != nil {
					return fmt.Errorf("parent 해시 변환 실패 (%s): %w", parent, err)
				}
				parentHashes = append(parentHashes, h)
			}

			// B 체인에 기록 (IPFS 재업로드 없이 기존 digest 그대로 사용)
			if err := chainClient.RecordCommit(
				forkRepoID,
				branch,
				expectedOldCommit,
				record.CommitHash,
				record.TreeHash,
				parentHashes,
				record.ManifestDigest,
				record.DiffDigest,
			); err != nil {
				return fmt.Errorf("체인 커밋 기록 실패 (커밋 %d, %s): %w", i+1, expectedCommit, err)
			}

			expectedOldCommit = record.CommitHash
			fmt.Printf("커밋 복원 완료 (%d/%d): %s\n", i+1, len(records), expectedCommit[:8])
		}

		// 8. .bit/config.json 저장
		cfg := &config.Config{
			RPCURL:          rpcURL,
			ContractAddress: contractAddress,
			PrivateKey:      privateKey,
			IPFSURL:         ipfsURL,
			RepoID:          forkRepoID.Uint64(),
			Remotes:         make(map[string]config.Remote),
		}
		if err := config.Save(".", cfg); err != nil {
			return fmt.Errorf("config 저장 실패: %w", err)
		}

		// 9. remote 설정
		forkURL := fmt.Sprintf("bit://local/%s/%s", contractAddress, forkRepoID.String())
		if err := config.AddRemote(".", "origin", config.Remote{URL: forkURL, RepoID: forkRepoID.Uint64()}); err != nil {
			return fmt.Errorf("origin remote 저장 실패: %w", err)
		}
		if err := config.AddRemote(".", "upstream", config.Remote{URL: sourceURL, RepoID: sourceRepoID}); err != nil {
			return fmt.Errorf("upstream remote 저장 실패: %w", err)
		}

		fmt.Printf("\nfork 완료!\n")
		fmt.Printf("  origin   → %s (repoId: %s)\n", forkURL, forkRepoID.String())
		fmt.Printf("  upstream → %s (repoId: %d)\n", sourceURL, sourceRepoID)
		fmt.Printf("\n새 커밋 후 push: bit push origin\n")
		fmt.Printf("PR 생성:         bit pr create upstream %s\n", branch)

		return nil
	},
}

func init() {
	forkCmd.Flags().String("rpc", "", "이더리움 RPC URL")
	forkCmd.Flags().String("contract", "", "BitRegistry 컨트랙트 주소")
	forkCmd.Flags().String("key", "", "지갑 개인키 (0x 제외)")
	forkCmd.Flags().String("ipfs", "", "IPFS 노드 주소 (기본값: http://localhost:5001)")
	forkCmd.Flags().String("branch", "main", "fork할 브랜치")
}
