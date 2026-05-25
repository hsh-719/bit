package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull <remote> <branch>",
	Short: "Pull branch from remote",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		remote, branch := args[0], args[1]

		// 1. 컨트랙트 getBranchHead 호출 → manifest CID
		// manifestCID, err := chain.GetBranchHead(repoId, branch)

		// 2. IPFS에서 manifest JSON 다운로드
		// manifest, err := ipfs.Download(manifestCID)

		// 3. manifest에서 코드 CID 추출 후 Git 객체 다운로드
		// objects, err := ipfs.Download(manifest.BundleCID)

		// 4. 로컬 .git에 반영
		// err = git.Import(objects)

		fmt.Printf("TODO: pull branch %s from remote %s\n", branch, remote)
		return nil
	},
}
