package ipfs

// Client wraps IPFS HTTP API calls.
type Client struct {
	apiURL string // e.g. "http://localhost:5001"
}

func NewClient(apiURL string) *Client {
	return &Client{apiURL: apiURL}
}

// Upload uploads data to IPFS and returns the CID string.
func (c *Client) Upload(data []byte) (string, error) {
	// TODO: multipart POST to /api/v0/add
	return "", nil
}

// Download fetches data from IPFS by CID.
func (c *Client) Download(cid string) ([]byte, error) {
	// TODO: POST to /api/v0/cat?arg=<cid>
	return nil, nil
}
