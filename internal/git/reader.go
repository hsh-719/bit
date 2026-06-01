package git

import (
	"bytes"
	"fmt"
	"os/exec"

	gogit "github.com/go-git/go-git/v5"
)

// HeadInfo는 로컬 .git에서 읽은 현재 브랜치와 커밋 정보를 담는다.
type HeadInfo struct {
	Branch     string // 현재 브랜치명 (예: "main")
	CommitHash string // 현재 커밋 해시
	ParentHash string // 부모 커밋 해시 (초기 커밋이면 빈 문자열)
}

// ReadHead는 로컬 git 저장소에서 현재 HEAD 정보를 읽어 반환한다.
// repoPath: git 저장소 루트 경로 (예: ".")
// push 시 커밋 해시와 브랜치명을 읽을 때 사용한다.
func ReadHead(repoPath string) (*HeadInfo, error) {
	repo, err := gogit.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commitHash := head.Hash().String()
	branch := head.Name().Short()

	// 부모 커밋 해시 읽기 (초기 커밋이면 ParentHashes가 비어있음)
	parentHash := ""
	commit, err := repo.CommitObject(head.Hash())
	if err == nil && len(commit.ParentHashes) > 0 {
		parentHash = commit.ParentHashes[0].String()
	}

	return &HeadInfo{
		Branch:     branch,
		CommitHash: commitHash,
		ParentHash: parentHash,
	}, nil
}

// ExtractBundle은 현재 저장소의 모든 Git 객체를 bundle 형태로 추출해 반환한다.
// push 시 코드 전체를 IPFS에 업로드하기 위해 사용한다.
// 내부적으로 "git bundle create - --all"을 실행한다.
func ExtractBundle(repoPath string) ([]byte, error) {
	cmd := exec.Command("git", "bundle", "create", "-", "--all")
	cmd.Dir = repoPath
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// ApplyBundle은 IPFS에서 받은 bundle 데이터를 로컬 git 저장소에 반영한다.
// pull 시 다운로드한 코드를 로컬에 적용할 때 사용한다.
// 내부적으로 "git bundle unbundle -" 후 "git checkout -b <branch> <commit>"을 실행한다.
func ApplyBundle(repoPath string, data []byte) error {
	// bundle을 로컬 .git에 풀기
	unbundle := exec.Command("git", "bundle", "unbundle", "-")
	unbundle.Dir = repoPath
	unbundle.Stdin = bytes.NewReader(data)
	var out bytes.Buffer
	unbundle.Stdout = &out
	if err := unbundle.Run(); err != nil {
		return err
	}

	// unbundle 출력에서 커밋 해시와 브랜치명 추출
	// 출력 형식: "<commitHash> refs/heads/<branch>"
	var commitHash, refName string
	fmt.Sscanf(out.String(), "%s %s", &commitHash, &refName)
	branch := refName[len("refs/heads/"):]

	// 해당 브랜치로 checkout
	checkout := exec.Command("git", "checkout", "-b", branch, commitHash)
	checkout.Dir = repoPath
	return checkout.Run()
}
