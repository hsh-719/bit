package cmd

import (
	"fmt"
	"math/big"

	"github.com/opendasom/bit/internal/chain"
	"github.com/opendasom/bit/internal/config"
	"github.com/opendasom/bit/internal/git"
	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests",
}

var prCreateCmd = &cobra.Command{
	Use:   "create <remote> <branch>",
	Short: "Create a pull request from the current branch to a remote branch",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		remoteName, targetBranch := args[0], args[1]

		cfg, err := config.Load(".")
		if err != nil {
			return fmt.Errorf("config 읽기 실패 (.bit/config.json): %w", err)
		}
		if cfg.RepoID == 0 {
			return fmt.Errorf("현재 저장소 repoId가 설정되지 않았습니다. 먼저 bit init 또는 bit fork를 실행하세요")
		}

		remote, ok := cfg.Remotes[remoteName]
		if !ok {
			return fmt.Errorf("remote '%s' 를 찾을 수 없습니다", remoteName)
		}

		sourceBranch, err := git.CurrentBranch(".")
		if err != nil {
			return fmt.Errorf("현재 브랜치 조회 실패: %w", err)
		}

		chainClient, err := chain.NewClient(cfg.RPCURL, cfg.ContractAddress, cfg.PrivateKey)
		if err != nil {
			return fmt.Errorf("체인 연결 실패: %w", err)
		}

		sourceRepoID := new(big.Int).SetUint64(cfg.RepoID)
		targetRepoID := new(big.Int).SetUint64(remote.RepoID)

		sourceHistoryLen, err := chainClient.GetBranchHistoryLength(sourceRepoID, sourceBranch)
		if err != nil {
			return fmt.Errorf("source 브랜치 히스토리 조회 실패: %w", err)
		}
		if sourceHistoryLen.Sign() == 0 {
			return fmt.Errorf("현재 브랜치 '%s' 에 커밋이 없습니다", sourceBranch)
		}

		targetHead, err := chainClient.GetBranchCommit(targetRepoID, targetBranch)
		if err != nil {
			return fmt.Errorf("target 브랜치 헤드 조회 실패: %w", err)
		}

		sourceRecords, err := loadBranchRecords(chainClient, sourceRepoID, sourceBranch, sourceHistoryLen.Int64())
		if err != nil {
			return err
		}

		if !chain.IsZeroBytes20(targetHead) {
			targetHeadHex := chain.Bytes20ToGitHash(targetHead)
			targetIndex := -1
			for i, record := range sourceRecords {
				if chain.Bytes20ToGitHash(record.CommitHash) == targetHeadHex {
					targetIndex = i
					break
				}
			}
			if targetIndex == -1 {
				return fmt.Errorf("target 브랜치 '%s' 가 source 히스토리에 없습니다. upstream이 먼저 앞서 나갔습니다", targetBranch)
			}
			if targetIndex >= len(sourceRecords)-1 {
				return fmt.Errorf("새로 추가된 커밋이 없습니다. '%s' 는 이미 최신입니다", sourceBranch)
			}
		}

		prID, err := chainClient.CreatePullRequest(targetRepoID, targetBranch, sourceRepoID, sourceBranch)
		if err != nil {
			return fmt.Errorf("PR 생성 실패: %w", err)
		}

		fmt.Printf("PR 생성 완료: #%s\n", prID.String())
		fmt.Printf("  source: current repo #%d (%s)\n", cfg.RepoID, sourceBranch)
		fmt.Printf("  target: %s #%d (%s)\n", remoteName, remote.RepoID, targetBranch)
		return nil
	},
}

func init() {
	prCmd.AddCommand(prCreateCmd)
}
