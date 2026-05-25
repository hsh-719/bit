package manifest

import "encoding/json"

// Manifest is uploaded to IPFS and its CID is recorded on-chain as the branch head.
type Manifest struct {
	BundleCID      string `json:"bundleCid"`      // IPFS CID of the Git object pack
	GitCommit      string `json:"gitCommit"`      // Git commit hash
	PreviousCommit string `json:"previousCommit"` // first-parent commit hash
}

func Encode(m *Manifest) ([]byte, error) {
	return json.Marshal(m)
}

func Decode(data []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
