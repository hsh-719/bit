package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hsh-719/bit/internal/chain"
	"github.com/hsh-719/bit/internal/config"
	"github.com/hsh-719/bit/internal/ipfs"
	"github.com/hsh-719/bit/internal/repo"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new bit repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		rpcURL, _ := cmd.Flags().GetString("rpc")
		contractAddress, _ := cmd.Flags().GetString("contract")
		privateKey, _ := cmd.Flags().GetString("key")
		ipfsURL, _ := cmd.Flags().GetString("ipfs")
		repoName, _ := cmd.Flags().GetString("name")
		repoDescription, _ := cmd.Flags().GetString("description")
		defaultBranch, _ := cmd.Flags().GetString("branch")

		// 1. .git 존재 여부 확인 (git init이 되어있어야 함)
		if _, err := os.Stat(".git"); os.IsNotExist(err) {
			return fmt.Errorf(".git 디렉토리가 없습니다. 먼저 git init을 실행하세요")
		}
		if repoName == "" {
			repoName = filepath.Base(mustGetwd())
		}
		if defaultBranch == "" {
			defaultBranch = "main"
		}

		metadata := &repo.Metadata{
			Version:       repo.MetadataVersion,
			Name:          repoName,
			Description:   repoDescription,
			DefaultBranch: defaultBranch,
		}
		metadataData, err := repo.EncodeMetadata(metadata)
		if err != nil {
			return fmt.Errorf("repo metadata 생성 실패: %w", err)
		}
		metadataCID, err := ipfs.NewClient(ipfsURL).Upload(metadataData)
		if err != nil {
			return fmt.Errorf("repo metadata IPFS 업로드 실패: %w", err)
		}

		// 2. 체인 연결 확인 및 저장소 등록 → repoId 발급
		chainClient, err := chain.NewClient(rpcURL, contractAddress, privateKey)
		if err != nil {
			return fmt.Errorf("체인 연결 실패: %w", err)
		}

		repoID, err := chainClient.CreateRepo(metadataCID)
		if err != nil {
			return fmt.Errorf("저장소 생성 실패: %w", err)
		}
		if repoID == nil {
			return fmt.Errorf("저장소 생성 실패: repoId를 확인할 수 없습니다")
		}
		fmt.Printf("체인 저장소 생성 완료 (repoId: %s)\n", repoID.String())
		fmt.Printf("repo metadata 등록 완료: %s\n", metadataCID)

		// 3. .bit/config.json 생성
		cfg := &config.Config{
			RPCURL:          rpcURL,
			ContractAddress: contractAddress,
			PrivateKey:      privateKey,
			IPFSURL:         ipfsURL,
			Remotes:         make(map[string]config.Remote),
			RepoID:          repoID.Uint64(),
		}
		if err := config.Save(".", cfg); err != nil {
			return fmt.Errorf("config 저장 실패: %w", err)
		}
		fmt.Println("config 저장 완료 (.bit/config.json)")

		return nil
	},
}

func init() {
	initCmd.Flags().String("rpc", "", "이더리움 RPC URL (예: https://mainnet.infura.io/v3/...)")
	initCmd.Flags().String("contract", "", "BitRegistry 컨트랙트 주소")
	initCmd.Flags().String("key", "", "지갑 개인키 (0x 제외)")
	initCmd.Flags().String("ipfs", "http://localhost:5001", "IPFS 노드 주소")
	initCmd.Flags().String("name", "", "웹에 표시할 저장소 이름 (기본값: 현재 디렉토리명)")
	initCmd.Flags().String("description", "", "웹에 표시할 저장소 설명")
	initCmd.Flags().String("branch", "main", "웹에서 기본으로 조회할 브랜치")

	initCmd.MarkFlagRequired("rpc")
	initCmd.MarkFlagRequired("contract")
	initCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(initCmd)
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}
