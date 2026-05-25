package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push <remote> <branch>",
	Short: "Push current branch to remote",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		remote, branch := args[0], args[1]

		// 1. .git에서 커밋 해시, 브랜치 정보 읽기
		// commit, err := git.ReadHead(".")

		// 2. Git 객체를 IPFS에 업로드 → CID 반환
		// cid, err := ipfs.Upload(objects)

		// 3. manifest JSON 생성 후 IPFS에 업로드 → manifest CID
		// manifestCID, err := manifest.Create(cid, commit)

		// 4. 컨트랙트 updateBranch 호출
		// err = chain.UpdateBranch(repoId, branch, oldHead, manifestCID, commitHash)

		fmt.Printf("TODO: push branch %s to remote %s\n", branch, remote)
		return nil
	},
}
