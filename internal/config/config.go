package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configDir = ".bit"
const configFile = "config.json"

// Remote는 bit remote add로 등록한 원격 저장소 정보를 담는다.
type Remote struct {
	URL    string `json:"url"`    // bit://<network>/<registry>/<repoId>
	RepoID uint64 `json:"repoId"` // 체인에서 발급된 저장소 ID
}

// Config는 .bit/config.json에 저장되는 프로젝트 설정이다.
type Config struct {
	RPCURL          string            `json:"rpcURL"`          // 이더리움 노드 주소
	ContractAddress string            `json:"contractAddress"` // BitRegistry 컨트랙트 주소
	PrivateKey      string            `json:"privateKey"`      // 트랜잭션 서명용 개인키
	IPFSURL         string            `json:"ipfsURL"`         // IPFS 노드 주소
	RepoID          uint64            `json:"repoId"`          // 체인에서 발급된 저장소 ID
	Remotes         map[string]Remote `json:"remotes"`         // remote 이름 → Remote 정보
}

// configPath는 .bit/config.json 경로를 반환한다.
func configPath(repoPath string) string {
	return filepath.Join(repoPath, configDir, configFile)
}

// Load는 .bit/config.json을 읽어 Config를 반환한다.
func Load(repoPath string) (*Config, error) {
	data, err := os.ReadFile(configPath(repoPath))
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Remotes == nil {
		cfg.Remotes = make(map[string]Remote)
	}
	return &cfg, nil
}

// Save는 Config를 .bit/config.json에 저장한다.
func Save(repoPath string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Join(repoPath, configDir), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(repoPath), data, 0600)
}

// AddRemote는 config에 remote를 추가하고 저장한다.
func AddRemote(repoPath, name string, remote Remote) error {
	cfg, err := Load(repoPath)
	if err != nil {
		return err
	}
	cfg.Remotes[name] = remote
	return Save(repoPath, cfg)
}
