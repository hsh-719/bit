package chain

import "math/big"

// Client wraps BitRegistry contract calls via go-ethereum.
type Client struct {
	rpcURL          string
	contractAddress string
	privateKey      string
}

func NewClient(rpcURL, contractAddress, privateKey string) *Client {
	return &Client{
		rpcURL:          rpcURL,
		contractAddress: contractAddress,
		privateKey:      privateKey,
	}
}

// GetBranchHead calls getBranchHead(repoId, branch) and returns the manifest CID.
func (c *Client) GetBranchHead(repoId *big.Int, branch string) (string, error) {
	// TODO: go-ethereum abigen 바인딩으로 구현
	return "", nil
}

// UpdateBranch calls updateBranch(...) on the contract.
func (c *Client) UpdateBranch(repoId *big.Int, branch, expectedOldHead, newHead, gitCommit, previousCommit string) error {
	// TODO: go-ethereum으로 트랜잭션 제출
	return nil
}

// CreateRepo calls createRepo(metadataCID) and returns the new repoId.
func (c *Client) CreateRepo(metadataCID string) (*big.Int, error) {
	// TODO: go-ethereum으로 트랜잭션 제출
	return nil, nil
}
