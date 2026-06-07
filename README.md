# bit

IPFS와 블록체인 위에서 동작하는 탈중앙화 버전 관리 시스템입니다.

코드는 IPFS에, 브랜치 상태는 스마트 컨트랙트(BitRegistry)에 저장됩니다.

---

## 설치

**사전 조건**

- Go 1.21+
- IPFS 데몬 (`ipfs daemon`)
- 이더리움 노드 접근 (로컬: Anvil, 테스트넷: Sepolia 등)

```bash
git clone https://github.com/hsh-719/bit.git
cd bit
go mod tidy
go build -o bit .

# 전역 명령어로 등록 (선택)
sudo cp ./bit /usr/local/bin/bit
```

---

## 빠른 시작

### 1. 저장소 초기화

기존 git 저장소 디렉토리에서 실행합니다.

```bash
cd my-project      # .git이 있는 디렉토리
git init           # 아직 git init 안 했다면

bit init \
  --rpc http://127.0.0.1:8545 \
  --contract 0xYourContractAddress \
  --key YourPrivateKeyHex
```

성공하면 `.bit/config.json`이 생성되고 체인에 저장소가 등록됩니다.

> `--ipfs` 플래그 생략 시 기본값 `http://localhost:5001` 사용

### 2. remote 추가

push/pull 대상을 등록합니다. URL 형식: `bit://<network>/<contractAddress>/<repoId>`

```bash
bit remote add origin bit://local/0xYourContractAddress/1
```

`repoId`는 `bit init` 실행 후 출력되는 숫자입니다.

### 3. push

현재 브랜치의 전체 코드를 IPFS에 올리고, 브랜치 상태를 체인에 기록합니다.

```bash
bit push origin
```

### 4. pull

다른 디렉토리(또는 다른 팀원)에서 코드를 받아옵니다.

```bash
mkdir other-project && cd other-project
git init

bit init \
  --rpc http://127.0.0.1:8545 \
  --contract 0xYourContractAddress \
  --key YourPrivateKeyHex

bit remote add origin bit://local/0xYourContractAddress/1

bit pull origin main
```

---

## 명령어 요약

| 명령어 | 설명 |
|---|---|
| `bit init --rpc <url> --contract <addr> --key <key>` | 저장소 초기화 및 체인 등록 |
| `bit remote add <name> <url>` | remote 추가 |
| `bit push <remote>` | 현재 브랜치를 push |
| `bit pull <remote> <branch>` | 지정 브랜치를 pull |

---

## 동작 원리

```
push:
  로컬 .git → git bundle → IPFS (bundleCID)
                         → manifest JSON → IPFS (manifestCID)
                                         → 스마트 컨트랙트 (브랜치 헤드 갱신)

pull:
  스마트 컨트랙트 (manifestCID 조회)
    → IPFS manifest 다운로드 (bundleCID 추출)
    → IPFS bundle 다운로드
    → 로컬 .git 적용
```

---

## 로컬 테스트 환경 (Anvil)

```bash
# 1. Anvil 실행
anvil

# 2. 컨트랙트 배포
cd contracts
forge create --broadcast --rpc-url http://127.0.0.1:8545 \
  --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
  src/BitRegistry.sol:BitRegistry

# 3. IPFS 데몬 실행
ipfs daemon
```

배포 후 출력되는 `Deployed to:` 주소를 `--contract` 플래그에 사용합니다.


---

## 프로젝트 구조

```
bit/
├── main.go
├── cmd/                  # CLI 명령어 (init, push, pull, remote)
├── internal/
│   ├── chain/            # 스마트 컨트랙트 연동 (go-ethereum)
│   ├── git/              # .git 읽기/쓰기 (go-git)
│   ├── ipfs/             # IPFS 업로드/다운로드
│   ├── manifest/         # manifest JSON 인코딩/디코딩
│   └── config/           # .bit/config.json 관리
└── contracts/
    └── src/BitRegistry.sol
```
