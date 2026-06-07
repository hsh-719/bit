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

var pullCmd = &cobra.Command{
	Use:   "pull <remote> <branch>",
	Short: "Pull branch from remote",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		remoteName, branch := args[0], args[1]

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

		// 2. 체인에서 브랜치 헤드 manifest CID 조회
		chainClient, err := chain.NewClient(cfg.RPCURL, cfg.ContractAddress, cfg.PrivateKey)
		if err != nil {
			return fmt.Errorf("체인 연결 실패: %w", err)
		}
		manifestCID, err := chainClient.GetBranchHead(repoID, branch)
		if err != nil {
			return fmt.Errorf("브랜치 헤드 조회 실패: %w", err)
		}
		if manifestCID == "" {
			return fmt.Errorf("브랜치 '%s' 가 아직 push된 적 없습니다", branch)
		}
		fmt.Printf("manifest CID: %s\n", manifestCID)

		// 3. IPFS에서 manifest JSON 다운로드
		ipfsClient := ipfs.NewClient(cfg.IPFSURL)
		manifestData, err := ipfsClient.Download(manifestCID)
		if err != nil {
			return fmt.Errorf("manifest 다운로드 실패: %w", err)
		}

		// 4. manifest 파싱 → bundleCID 추출
		m, err := manifest.Decode(manifestData)
		if err != nil {
			return fmt.Errorf("manifest 파싱 실패: %w", err)
		}
		fmt.Printf("bundle CID: %s, 커밋: %s\n", m.BundleCID, m.GitCommit[:8])

		// 5. IPFS에서 bundle 다운로드
		bundle, err := ipfsClient.Download(m.BundleCID)
		if err != nil {
			return fmt.Errorf("bundle 다운로드 실패: %w", err)
		}

		// 6. 로컬 .git에 bundle 반영
		if err := git.ApplyBundle(".", bundle); err != nil {
			return fmt.Errorf("bundle 적용 실패: %w", err)
		}

		fmt.Printf("pull 완료: %s (%s)\n", branch, m.GitCommit[:8])
		return nil
	},
}
