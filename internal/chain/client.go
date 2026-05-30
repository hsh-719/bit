package chain

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

// Client는 BitRegistry 스마트 컨트랙트와 통신하는 클라이언트다.
// contract: abigen으로 생성된 Go 바인딩 (bitregistry.go)
// auth: 트랜잭션 서명에 사용하는 지갑 정보
// conn: 이더리움 노드 RPC 연결
type Client struct {
	contract *Chain
	auth     *bind.TransactOpts
	conn     *ethclient.Client
}

// NewClient는 RPC 노드에 연결하고 개인키로 지갑을 초기화해 Client를 반환한다.
// rpcURL: 이더리움 노드 주소 (예: "http://localhost:8545")
// contractAddress: 배포된 BitRegistry 컨트랙트 주소
// privateKeyHex: 트랜잭션 서명에 사용할 지갑 개인키 (0x 제외 hex 문자열)
func NewClient(rpcURL, contractAddress, privateKeyHex string) (*Client, error) {
	// 이더리움 노드에 RPC 연결
	conn, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	// hex 문자열 개인키를 ECDSA 키 객체로 변환
	privKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	// 체인 ID 조회 (트랜잭션 서명 시 리플레이 공격 방지용)
	chainID, err := conn.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	// 개인키와 체인 ID로 트랜잭션 서명 객체 생성
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		return nil, err
	}

	// 컨트랙트 주소로 Go 바인딩 인스턴스 생성
	contract, err := NewChain(common.HexToAddress(contractAddress), conn)
	if err != nil {
		return nil, err
	}

	return &Client{contract: contract, auth: auth, conn: conn}, nil
}

// branchNameToBytes32는 브랜치명 문자열을 keccak256 해시한 [32]byte로 변환한다.
// 컨트랙트는 브랜치명을 keccak256(branchName) 형태의 bytes32로 저장하기 때문에
// Go에서 컨트랙트를 호출할 때도 동일하게 변환해서 넘겨야 한다.
func branchNameToBytes32(name string) [32]byte {
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(name))
	var result [32]byte
	copy(result[:], h.Sum(nil))
	return result
}

// GetBranchHead는 특정 브랜치의 현재 헤드 manifest CID를 체인에서 조회한다.
// pull 시 "어떤 CID를 IPFS에서 받아올지" 결정하는 첫 번째 단계에서 사용한다.
// 브랜치가 아직 한 번도 push된 적 없으면 빈 문자열을 반환한다.
func (c *Client) GetBranchHead(repoId *big.Int, branch string) (string, error) {
	cid, err := c.contract.GetBranchHead(&bind.CallOpts{}, repoId, branchNameToBytes32(branch))
	if err != nil {
		return "", err
	}
	return string(cid), nil
}

// UpdateBranch는 브랜치 헤드를 새 manifest CID로 갱신하는 트랜잭션을 제출한다.
// push 시 IPFS 업로드가 완료된 후 마지막 단계에서 호출한다.
//
// compare-and-swap 방식으로 동작하기 때문에,
// expectedOldHead가 현재 체인 상태와 다르면 컨트랙트가 revert한다.
// (다른 사람이 먼저 push한 경우 발생)
func (c *Client) UpdateBranch(repoId *big.Int, branch, expectedOldHead, newHead, gitCommit, previousCommit string) error {
	_, err := c.contract.UpdateBranch(
		c.auth,
		repoId,
		branchNameToBytes32(branch),
		[]byte(expectedOldHead),
		[]byte(newHead),
		[]byte(gitCommit),
		[]byte(previousCommit),
	)
	return err
}

// CreateRepo는 체인에 새 저장소를 생성하고 repoId를 반환한다.
// 저장소를 처음 만들 때 한 번 호출한다. (git init에 해당)
// 트랜잭션이 채굴될 때까지 대기한 후 RepoCreated 이벤트에서 repoId를 추출한다.
func (c *Client) CreateRepo(metadataCID string) (*big.Int, error) {
	tx, err := c.contract.CreateRepo(c.auth, []byte(metadataCID))
	if err != nil {
		return nil, err
	}

	// 트랜잭션이 블록에 포함될 때까지 대기
	receipt, err := bind.WaitMined(context.Background(), c.conn, tx)
	if err != nil {
		return nil, err
	}

	// 트랜잭션 영수증의 로그에서 RepoCreated 이벤트를 찾아 repoId 추출
	for _, log := range receipt.Logs {
		event, err := c.contract.ParseRepoCreated(*log)
		if err != nil {
			continue
		}
		return event.RepoId, nil
	}

	return nil, nil
}

// GetPublicAddress는 개인키로부터 지갑 주소를 계산해 반환한다.
// 현재 어떤 지갑으로 트랜잭션을 보내는지 확인할 때 사용한다.
func GetPublicAddress(privateKeyHex string) (string, error) {
	privKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	pubKey := privKey.Public().(*ecdsa.PublicKey)
	return crypto.PubkeyToAddress(*pubKey).Hex(), nil
}
