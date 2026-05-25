package git

// HeadInfo holds the current branch and commit hash read from .git
type HeadInfo struct {
	Branch     string
	CommitHash string
}

// ReadHead reads the current HEAD info from a local git repository at repoPath.
func ReadHead(repoPath string) (*HeadInfo, error) {
	// TODO: go-git으로 구현
	// repo, err := gogit.PlainOpen(repoPath)
	// head, err := repo.Head()
	// return &HeadInfo{Branch: head.Name().Short(), CommitHash: head.Hash().String()}, nil
	return nil, nil
}
