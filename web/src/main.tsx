import React, { useEffect, useMemo, useState } from "react";
import { createRoot } from "react-dom/client";
import {
  createPublicClient,
  createWalletClient,
  custom,
  http,
  keccak256,
  stringToBytes,
  type Address,
  type Hex,
} from "viem";
import bitRegistryArtifact from "../../internal/chain/artifacts/BitRegistry.json";
import "./styles.css";

type RepoMetadata = {
  version?: number;
  name?: string;
  description?: string;
  defaultBranch?: string;
};

type RepoSummary = {
  id: bigint;
  owner: Address;
  metadataCID: string;
  metadata: RepoMetadata | null;
};

type BranchSummary = {
  name: string;
  branchHash: Hex;
  commitCount: number;
  headCommit: string;
};

type CommitSummary = {
  hash: string;
  treeHash: string;
  updater: Address;
  chainTimestamp: bigint;
  message: string;
  authorName: string;
  authorEmail: string;
  authorDate: string;
  committerName: string;
  committerEmail: string;
  committerDate: string;
  parents: string[];
};

type PullRequestSummary = {
  id: bigint;
  targetRepoId: bigint;
  targetBranch: Hex;
  sourceRepoId: bigint;
  sourceBranch: Hex;
  baseCommit: Hex;
  sourceHeadCommit: Hex;
  author: Address;
  status: bigint;
  createdAt: bigint;
  updatedAt: bigint;
};

type Manifest = {
  gitCommit: string;
  treeHash: string;
  branch?: string;
  parentCommits?: string[];
  author?: Identity;
  committer?: Identity;
  message?: string;
};

type Identity = {
  name?: string;
  email?: string;
  date?: string;
};

type EthereumProvider = {
  request: (args: { method: string; params?: unknown[] }) => Promise<unknown>;
};

type PageState = "home" | "project";

declare global {
  interface Window {
    ethereum?: EthereumProvider;
  }
}

const abi = bitRegistryArtifact.abi;
const defaultRpcURL = localStorage.getItem("bit.rpcURL") ?? "http://127.0.0.1:8545";
const defaultContract = localStorage.getItem("bit.contract") ?? "";
const defaultGateway = localStorage.getItem("bit.ipfsGateway") ?? "http://127.0.0.1:8080/ipfs";
const defaultIpfsAPI = localStorage.getItem("bit.ipfsAPI") ?? "http://127.0.0.1:5001";
const ROLE_LABELS = ["None", "Contributor", "Maintainer", "Owner"] as const;
type RoleLabel = (typeof ROLE_LABELS)[number];

function App() {
  const initialRoute = routeFromLocation(window.location.pathname, window.location.search);
  const [page, setPage] = useState<PageState>(initialRoute.page);
  const [selectedRepoId, setSelectedRepoId] = useState<bigint | null>(initialRoute.repoId);
  const [selectedBranch, setSelectedBranch] = useState<string>(initialRoute.branch ?? "");
  const [activeTab, setActiveTab] = useState<"commits" | "prs">("commits");
  const [rpcURL, setRpcURL] = useState(defaultRpcURL);
  const [contractAddress, setContractAddress] = useState(defaultContract);
  const [ipfsGateway, setIpfsGateway] = useState(defaultGateway);
  const [repos, setRepos] = useState<RepoSummary[]>([]);
  const [branches, setBranches] = useState<BranchSummary[]>([]);
  const [commits, setCommits] = useState<CommitSummary[]>([]);
  const [pullRequests, setPullRequests] = useState<PullRequestSummary[]>([]);
  const [repoRole, setRepoRole] = useState<RoleLabel>("None");
  const [walletAddress, setWalletAddress] = useState<Address | null>(null);
  const [walletChainId, setWalletChainId] = useState<string>("");
  const [loadingRepos, setLoadingRepos] = useState(false);
  const [loadingDetail, setLoadingDetail] = useState(false);
  const [loadingAction, setLoadingAction] = useState<string | null>(null);
  const [copyState, setCopyState] = useState<"idle" | "copied">("idle");
  const [error, setError] = useState("");

  const publicClient = useMemo(() => createPublicClient({ transport: http(rpcURL) }), [rpcURL]);
  const selectedRepo = repos.find((repo) => repo.id === selectedRepoId) ?? null;
  const branchNameByHash = useMemo(() => {
    const map = new Map<string, string>();
    for (const branch of branches) {
      map.set(branch.branchHash.toLowerCase(), branch.name);
    }
    return map;
  }, [branches]);

  useEffect(() => {
    const onPopState = () => {
      const nextRoute = routeFromLocation(window.location.pathname, window.location.search);
      setPage(nextRoute.page);
      setSelectedRepoId(nextRoute.repoId);
      setSelectedBranch(nextRoute.branch ?? "");
      setActiveTab("commits");
      setCopyState("idle");

      if (nextRoute.page === "home") {
        setCommits([]);
        setPullRequests([]);
        setBranches([]);
        setRepoRole("None");
        return;
      }

      if (nextRoute.repoId && repos.length > 0) {
        void loadRepoDetail(nextRoute.repoId, repos, nextRoute.branch ?? undefined);
      }
    };

    window.addEventListener("popstate", onPopState);
    return () => window.removeEventListener("popstate", onPopState);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [repos, rpcURL, contractAddress, ipfsGateway, walletAddress]);

  async function loadRepos() {
    setLoadingRepos(true);
    setError("");
    try {
      const address = parseAddress(contractAddress);
      localStorage.setItem("bit.rpcURL", rpcURL);
      localStorage.setItem("bit.contract", contractAddress);
      localStorage.setItem("bit.ipfsGateway", ipfsGateway);
      localStorage.setItem("bit.ipfsAPI", defaultIpfsAPI);

      const count = (await publicClient.readContract({
        address,
        abi,
        functionName: "getRepoCount",
      })) as bigint;

      const nextRepos: RepoSummary[] = [];
      for (let repoId = 1n; repoId <= count; repoId++) {
        const [owner, metadataBytes] = (await publicClient.readContract({
          address,
          abi,
          functionName: "getRepo",
          args: [repoId],
        })) as [Address, Hex];
        const metadataCID = bytesHexToString(metadataBytes);
        const metadata = metadataCID ? await fetchJson<RepoMetadata>(ipfsURL(ipfsGateway, metadataCID)) : null;
        nextRepos.push({ id: repoId, owner, metadataCID, metadata });
      }

      setRepos(nextRepos);
      setPage("home");
      setSelectedRepoId(null);
      setSelectedBranch("");
      setActiveTab("commits");
      setCommits([]);
      setPullRequests([]);
      setBranches([]);
      setRepoRole("None");
      window.history.replaceState({}, "", "/");
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setLoadingRepos(false);
    }
  }

  async function openRepo(repoId: bigint) {
    const repo = repos.find((item) => item.id === repoId);
    const branch = repo?.metadata?.defaultBranch || "main";
    const nextPath = `/projects/${repoId.toString()}?branch=${encodeURIComponent(branch)}`;
    window.history.pushState({}, "", nextPath);
    setPage("project");
    setSelectedRepoId(repoId);
    setSelectedBranch(branch);
    setActiveTab("commits");
    setCopyState("idle");
    await loadRepoDetail(repoId, repos, branch);
  }

  async function goHome() {
    window.history.pushState({}, "", "/");
    setPage("home");
    setSelectedRepoId(null);
    setSelectedBranch("");
    setActiveTab("commits");
    setCopyState("idle");
    setCommits([]);
    setPullRequests([]);
    setBranches([]);
    setRepoRole("None");
    setError("");
  }

  async function loadRepoDetail(repoId: bigint, repoLookup: RepoSummary[] = repos, branchName?: string) {
    setLoadingDetail(true);
    setError("");
    try {
      const address = parseAddress(contractAddress);
      const repo = repoLookup.find((item) => item.id === repoId) ?? selectedRepo;
      const branch = branchName || selectedBranch || repo?.metadata?.defaultBranch || "main";
      setSelectedBranch(branch);
      const branchKey = keccak256(stringToBytes(branch));

      const historyLength = (await publicClient.readContract({
        address,
        abi,
        functionName: "getBranchHistoryLength",
        args: [repoId, branchKey],
      })) as bigint;

      const pageSize = historyLength > 50n ? 50n : historyLength;
      const start = historyLength > pageSize ? historyLength - pageSize : 0n;
      const [hashes, treeHashes, manifestDigests] = (await publicClient.readContract({
        address,
        abi,
        functionName: "getBranchCommitsWithMetadata",
        args: [repoId, branchKey, start, pageSize],
      })) as [Hex[], Hex[], Hex[], Hex[]];

      const nextCommits = await Promise.all(
        hashes.map(async (hash, index) => {
          const [,,, updater, chainTimestamp] = (await publicClient.readContract({
            address,
            abi,
            functionName: "getCommit",
            args: [repoId, hash],
          })) as [Hex, Hex, Hex, Address, bigint];
          const manifestCID = cidV0FromDigest(manifestDigests[index]);
          const manifest = await fetchJson<Manifest>(ipfsURL(ipfsGateway, manifestCID));

          return {
            hash: bytes20HexToGitHash(hash),
            treeHash: bytes20HexToGitHash(treeHashes[index]),
            updater,
            chainTimestamp,
            message: manifest.message ?? "",
            authorName: manifest.author?.name ?? "",
            authorEmail: manifest.author?.email ?? "",
            authorDate: manifest.author?.date ?? "",
            committerName: manifest.committer?.name ?? "",
            committerEmail: manifest.committer?.email ?? "",
            committerDate: manifest.committer?.date ?? "",
            parents: manifest.parentCommits ?? [],
          };
        }),
      );

      const prCount = (await publicClient.readContract({
        address,
        abi,
        functionName: "getRepoPullRequestCount",
        args: [repoId],
      })) as bigint;

      const nextPullRequests: PullRequestSummary[] = [];
      for (let index = 0n; index < prCount; index++) {
        const prId = (await publicClient.readContract({
          address,
          abi,
          functionName: "getRepoPullRequestAt",
          args: [repoId, index],
        })) as bigint;
        const pr = (await publicClient.readContract({
          address,
          abi,
          functionName: "getPullRequest",
          args: [prId],
        })) as [
          bigint,
          bigint,
          Hex,
          bigint,
          Hex,
          Hex,
          Hex,
          Address,
          bigint,
          bigint,
          bigint,
        ];
        if (pr[8] !== 1n) continue;
        nextPullRequests.push({
          id: pr[0],
          targetRepoId: pr[1],
          targetBranch: pr[2],
          sourceRepoId: pr[3],
          sourceBranch: pr[4],
          baseCommit: pr[5],
          sourceHeadCommit: pr[6],
          author: pr[7],
          status: pr[8],
          createdAt: pr[9],
          updatedAt: pr[10],
        });
      }

      setCommits(nextCommits.reverse());
      setPullRequests(nextPullRequests);
      await loadRepoBranches(repoId, repoLookup, branch);

      if (walletAddress) {
        const role = (await publicClient.readContract({
          address,
          abi,
          functionName: "getRole",
          args: [repoId, walletAddress],
        })) as bigint;
        setRepoRole(roleToLabel(role));
      } else {
        setRepoRole("None");
      }
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setLoadingDetail(false);
    }
  }

  async function loadRepoBranches(repoId: bigint, repoLookup: RepoSummary[] = repos, defaultBranch = "main") {
    try {
      const repo = repoLookup.find((item) => item.id === repoId) ?? selectedRepo;
      const logs = await publicClient.getContractEvents({
        address: parseAddress(contractAddress),
        abi,
        eventName: "CommitRecorded",
        args: { repoId },
        fromBlock: 0n,
        toBlock: "latest",
      });

      const nextBranches = new Map<string, BranchSummary>();
      for (const log of logs) {
        const manifestCID = cidV0FromDigest(log.args.manifestDigest);
        const manifest = await getManifest(ipfsGateway, manifestCID);
        const branchName = manifest.branch || defaultBranch || repo?.metadata?.defaultBranch || "main";
        nextBranches.set(branchName, {
          name: branchName,
          branchHash: keccak256(stringToBytes(branchName)),
          commitCount: (nextBranches.get(branchName)?.commitCount ?? 0) + 1,
          headCommit: bytes20HexToGitHash(log.args.commitHash),
        });
      }

      const sorted = Array.from(nextBranches.values()).sort((left, right) => {
        if (left.name === defaultBranch) return -1;
        if (right.name === defaultBranch) return 1;
        return left.name.localeCompare(right.name);
      });
      setBranches(sorted);
    } catch (err) {
      setError(errorMessage(err));
    }
  }

  async function connectWallet() {
    try {
      if (!window.ethereum) {
        setError("MetaMask가 설치되어 있지 않습니다.");
        return;
      }
      const accounts = (await window.ethereum.request({ method: "eth_requestAccounts" })) as string[];
      const chainId = (await window.ethereum.request({ method: "eth_chainId" })) as string;
      if (!accounts?.[0]) {
        setError("지갑 계정을 가져오지 못했습니다.");
        return;
      }
      setWalletAddress(accounts[0] as Address);
      setWalletChainId(chainId);
      setError("");
    } catch (err) {
      setError(errorMessage(err));
    }
  }

  async function approvePullRequest(prId: bigint) {
    if (!window.ethereum) {
      setError("MetaMask가 필요합니다.");
      return;
    }
    if (!walletAddress) {
      setError("먼저 MetaMask를 연결하세요.");
      return;
    }
    if (!selectedRepoId) {
      setError("선택된 프로젝트가 없습니다.");
      return;
    }

    setLoadingAction(`approve-${prId.toString()}`);
    setError("");
    try {
      const address = parseAddress(contractAddress);
      const walletClient = createWalletClient({
        account: walletAddress,
        transport: custom(window.ethereum),
      });
      const txHash = await walletClient.writeContract({
        address,
        abi,
        functionName: "approvePullRequest",
        args: [prId],
      });
      await publicClient.waitForTransactionReceipt({ hash: txHash });
      await loadRepoDetail(selectedRepoId);
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setLoadingAction(null);
    }
  }

  async function copyForkCommand() {
    if (!selectedRepo) return;

    const command = buildForkCommand({
      contractAddress,
      repoId: selectedRepo.id,
      branch: selectedBranch || selectedRepo.metadata?.defaultBranch || "main",
      rpcURL,
      ipfsAPI: defaultIpfsAPI,
    });

    try {
      await navigator.clipboard.writeText(command);
      setCopyState("copied");
      window.setTimeout(() => setCopyState("idle"), 1500);
    } catch (err) {
      setError(errorMessage(err));
    }
  }

  const walletSummary = walletAddress ? `${shortAddress(walletAddress)} · ${formatChainId(walletChainId)}` : "";
  const contractSummary = contractAddress ? shortAddress(contractAddress) : "n/a";

  return (
    <main className="page">
      <header className="siteHeader">
        <a className="siteBrand" href="/" aria-label="Go to home">
          <div className="brandMark" aria-hidden="true">
            <span />
            <span />
            <span />
          </div>
          <div>
            <div className="brandName">BIT</div>
            <div className="brandTag">Blockchain-based version control</div>
          </div>
        </a>

        <div className="headerActions">
          <div className="headerChip">
            <span>Contract</span>
            <strong className="mono">{contractSummary}</strong>
          </div>
          {walletAddress ? (
            <div className="headerChip headerWalletChip">
              <span>MetaMask</span>
              <strong className="mono">{walletSummary}</strong>
            </div>
          ) : (
            <button type="button" className="ghostButton headerButton" onClick={connectWallet}>
              Connect MetaMask
            </button>
          )}
        </div>
      </header>

      {error && <div className="errorBanner">{error}</div>}

      {page === "home" && (
        <>
          <section className="heroBand">
            <div className="heroCopy">
              <div className="heroKicker">01. Landing Page</div>
              <h1>BIT</h1>
              <p>Repository history and pull request metadata, recorded on-chain and surfaced read-only in the browser.</p>
              <div className="heroActions">
                <button type="button" className="primaryButton heroButton" onClick={loadRepos} disabled={loadingRepos}>
                  {loadingRepos ? "Loading..." : "Load repositories"}
                </button>
                <button
                  type="button"
                  className="ghostButton heroButton"
                  onClick={() => document.getElementById("projects-band")?.scrollIntoView({ behavior: "smooth", block: "start" })}
                >
                  View projects
                </button>
              </div>
            </div>

            <div className="heroVisual" aria-hidden="true">
              <div className="cube">
                <span />
                <span />
                <span />
              </div>
            </div>
          </section>

          <section className="projectsBand" id="projects-band">
            <div className="bandHeading">
              <div>
                <div className="eyebrow">Projects</div>
                <h2>Repository list</h2>
                <p>Each entry is a repository registered in the current BitRegistry contract.</p>
              </div>
              <div className="panelBadge">{repos.length} repositories</div>
            </div>

            <div className="projectList">
              {repos.length === 0 && <div className="emptyStage">Load repositories to see the list.</div>}
              {repos.map((repo) => (
                <button
                  key={repo.id.toString()}
                  type="button"
                  className="projectRow"
                  onClick={() => {
                    void openRepo(repo.id);
                  }}
                >
                  <div className="projectMain">
                    <div className="projectTitle">{repo.metadata?.name || `Repo #${repo.id}`}</div>
                    <div className="projectDescription">{repo.metadata?.description || "No description provided."}</div>
                  </div>
                  <div className="projectMeta">
                    <span className="mono">{shortAddress(repo.owner)}</span>
                    <span>{repo.metadata?.defaultBranch || "main"}</span>
                  </div>
                </button>
              ))}
            </div>
          </section>
        </>
      )}

      {page === "project" && selectedRepo && (
        <section className="detailBand" id="detail-band">
          <header className="detailHeader">
            <div>
              <div className="eyebrow">Project</div>
              <h2>{selectedRepo.metadata?.name || "Project detail"}</h2>
              <p>{selectedRepo.metadata?.description || "Metadata and pull request state for the selected repository."}</p>
            </div>
            <div className="detailHeaderActions">
              <button type="button" className="ghostButton copyButton" onClick={() => void copyForkCommand()} disabled={!selectedRepo}>
                {copyState === "copied" ? "Copied fork command" : "Copy fork command"}
              </button>
            </div>
          </header>

          <div className="detailShell">
            <aside className="detailNav">
              <button type="button" className={activeTab === "commits" ? "tab active" : "tab"} onClick={() => setActiveTab("commits")}>
                Commits
              </button>
              <button type="button" className={activeTab === "prs" ? "tab active" : "tab"} onClick={() => setActiveTab("prs")}>
                Pull Requests
              </button>
            </aside>

            <div className="detailContent">
              {activeTab === "commits" && (
                <section className="detailGrid">
                  <div className="panel panelTall">
                    <div className="panelHeading">
                      <div>
                        <span className="eyebrow">Commit History</span>
                        <h3>Metadata only</h3>
                      </div>
                      <div className="panelBadge">{commits.length} commits</div>
                    </div>
                    <div className="timeline">
                      {loadingDetail && commits.length === 0 && <div className="emptyState">Loading commits...</div>}
                      {commits.map((commit) => (
                        <article className="timelineItem" key={commit.hash}>
                          <div className="timelineMark" />
                          <div className="timelineBody">
                            <div className="timelineTop">
                              <h4>{commit.message || "(no message)"}</h4>
                              <span className="mono commitHash">{shortHex(commit.hash)}</span>
                            </div>
                            <div className="timelineMeta">
                              <span>{commit.authorName || "Unknown author"}</span>
                              <span>{formatDate(commit.authorDate)}</span>
                              <span>{shortAddress(commit.updater)}</span>
                            </div>
                            <div className="timelineFoot">
                              <span className="mono">tree {shortHex(commit.treeHash)}</span>
                              <span className="mono">parents {commit.parents.length > 0 ? commit.parents.join(", ") : "none"}</span>
                            </div>
                          </div>
                        </article>
                      ))}
                    </div>
                  </div>

                  <aside className="panel stackPanel">
                    <div className="panelHeading">
                      <div>
                        <span className="eyebrow">Branches</span>
                        <h3>Branch list</h3>
                      </div>
                    </div>
                    <div className="branchList">
                      {branches.length === 0 && <div className="emptyState">No branches discovered yet.</div>}
                      {branches.map((branch) => (
                        <button
                          key={branch.branchHash}
                          type="button"
                          className={branch.name === selectedBranch ? "branchItem active" : "branchItem"}
                          onClick={() => {
                            const nextPath = `/projects/${selectedRepo.id.toString()}?branch=${encodeURIComponent(branch.name)}`;
                            window.history.pushState({}, "", nextPath);
                            setSelectedBranch(branch.name);
                            setCopyState("idle");
                            void loadRepoDetail(selectedRepo.id, repos, branch.name);
                          }}
                        >
                          <div className="branchItemTop">
                            <strong>{branch.name}</strong>
                            {branch.name === selectedBranch && <span className="statusChip ok">CURRENT</span>}
                          </div>
                          <div className="branchItemBottom">
                            <span className="mono">{branch.commitCount} commits</span>
                            <span className="mono">{shortHex(branch.headCommit)}</span>
                          </div>
                        </button>
                      ))}
                    </div>

                    <div className="panelHeading">
                      <div>
                        <span className="eyebrow">Repository</span>
                        <h3>{selectedRepo.metadata?.name || `Repo #${selectedRepo.id}`}</h3>
                      </div>
                    </div>
                    <div className="kvList">
                      <div>
                        <span>Owner</span>
                        <strong className="mono">{shortAddress(selectedRepo.owner)}</strong>
                      </div>
                      <div>
                        <span>Default branch</span>
                        <strong>{selectedBranch || selectedRepo.metadata?.defaultBranch || "main"}</strong>
                      </div>
                      <div>
                        <span>Metadata CID</span>
                        <strong className="mono">{selectedRepo.metadataCID || "n/a"}</strong>
                      </div>
                    </div>
                  </aside>
                </section>
              )}

              {activeTab === "prs" && (
                <section className="detailGrid">
                  <div className="panel panelTall">
                    <div className="panelHeading">
                      <div>
                        <span className="eyebrow">Pull Requests</span>
                        <h3>Open PRs</h3>
                      </div>
                      <div className="panelBadge">{pullRequests.length} active</div>
                    </div>
                    <div className="prList">
                      {loadingDetail && pullRequests.length === 0 && <div className="emptyState">Loading pull requests...</div>}
                      {pullRequests.length === 0 && !loadingDetail && <div className="emptyState">No open pull requests.</div>}
                      {pullRequests.map((pr) => (
                        <article className="prCard" key={pr.id.toString()}>
                          <div className="prHeader">
                            <div>
                              <div className="prTitle">#{pr.id.toString()}</div>
                              <div className="prSubtitle">
                                <span className="mono">
                                  from {branchNameByHash.get(pr.sourceBranch.toLowerCase()) ?? `#${pr.sourceRepoId.toString()}`}
                                </span>
                                <span className="mono">
                                  to {branchNameByHash.get(pr.targetBranch.toLowerCase()) ?? `#${pr.targetRepoId.toString()}`}
                                </span>
                              </div>
                            </div>
                            <span className="statusChip open">OPEN</span>
                          </div>
                          <div className="prMetaGrid">
                            <div>
                              <span>Author</span>
                              <strong className="mono">{shortAddress(pr.author)}</strong>
                            </div>
                            <div>
                              <span>Created</span>
                              <strong>{formatUnix(pr.createdAt)}</strong>
                            </div>
                            <div>
                              <span>Base</span>
                              <strong className="mono">{shortHex(pr.baseCommit)}</strong>
                            </div>
                            <div>
                              <span>Head</span>
                              <strong className="mono">{shortHex(pr.sourceHeadCommit)}</strong>
                            </div>
                          </div>
                          <div className="prActions">
                            <button
                              type="button"
                              className="primaryButton"
                              disabled={!walletAddress || loadingAction === `approve-${pr.id.toString()}`}
                              onClick={() => void approvePullRequest(pr.id)}
                            >
                              {loadingAction === `approve-${pr.id.toString()}` ? "Approving..." : "Approve"}
                            </button>
                            {!walletAddress && <span className="helperText">Connect MetaMask to sign.</span>}
                          </div>
                        </article>
                      ))}
                    </div>
                  </div>

                  <aside className="panel stackPanel">
                    <div className="panelHeading">
                      <div>
                        <span className="eyebrow">Authorization</span>
                        <h3>Signer</h3>
                      </div>
                    </div>
                    <div className="kvList">
                      <div>
                        <span>Wallet</span>
                        <strong>{walletAddress ? shortAddress(walletAddress) : "Not connected"}</strong>
                      </div>
                      <div>
                        <span>Role on repo</span>
                        <strong>{repoRole}</strong>
                      </div>
                      <div>
                        <span>Approval rule</span>
                        <strong>Owner or Maintainer</strong>
                      </div>
                      <div>
                        <span>Chain</span>
                        <strong>{formatChainId(walletChainId) || "unknown"}</strong>
                      </div>
                    </div>
                    <p className="helperText">Approve opens MetaMask and submits the transaction for the connected account.</p>
                  </aside>
                </section>
              )}
            </div>
          </div>
        </section>
      )}
    </main>
  );
}

function routeFromLocation(
  pathname: string,
  search: string,
): { page: PageState; repoId: bigint | null; branch: string | null } {
  const match = pathname.match(/^\/projects\/(\d+)\/?$/);
  if (!match) {
    return { page: "home", repoId: null, branch: null };
  }
  const params = new URLSearchParams(search);
  return { page: "project", repoId: BigInt(match[1]), branch: params.get("branch") };
}

function buildForkCommand(options: {
  contractAddress: string;
  repoId: bigint;
  branch: string;
  rpcURL: string;
  ipfsAPI: string;
}): string {
  const contract = options.contractAddress || "0xYourContractAddress";
  return [
    "bit fork",
    `bit://local/${contract}/${options.repoId.toString()}`,
    `--rpc ${options.rpcURL}`,
    `--contract ${contract}`,
    `--key <YOUR_PRIVATE_KEY>`,
    `--ipfs ${options.ipfsAPI}`,
    `--branch ${options.branch}`,
  ].join(" ");
}

function parseAddress(value: string): Address {
  if (!/^0x[a-fA-F0-9]{40}$/.test(value)) {
    throw new Error("Contract address must be a 20-byte hex address.");
  }
  return value as Address;
}

function bytesHexToString(value: Hex): string {
  const hex = value.startsWith("0x") ? value.slice(2) : value;
  let out = "";
  for (let index = 0; index < hex.length; index += 2) {
    const code = Number.parseInt(hex.slice(index, index + 2), 16);
    if (code !== 0) out += String.fromCharCode(code);
  }
  return out;
}

function bytes20HexToGitHash(value: Hex): string {
  return value.startsWith("0x") ? value.slice(2) : value;
}

function ipfsURL(gateway: string, cid: string): string {
  return `${gateway.replace(/\/$/, "")}/${cid}`;
}

async function fetchJson<T>(url: string): Promise<T> {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`IPFS fetch failed (${response.status})`);
  }
  return response.json() as Promise<T>;
}

const manifestCache = new Map<string, Manifest>();

async function getManifest(gateway: string, cid: string): Promise<Manifest> {
  const cacheKey = `${gateway}|${cid}`;
  const cached = manifestCache.get(cacheKey);
  if (cached) return cached;
  const manifest = await fetchJson<Manifest>(ipfsURL(gateway, cid));
  manifestCache.set(cacheKey, manifest);
  return manifest;
}

function cidV0FromDigest(digestHex: Hex): string {
  const digest = hexToBytes(digestHex);
  return base58btcEncode(new Uint8Array([0x12, 0x20, ...digest]));
}

function hexToBytes(value: Hex): Uint8Array {
  const hex = value.startsWith("0x") ? value.slice(2) : value;
  const out = new Uint8Array(hex.length / 2);
  for (let index = 0; index < out.length; index += 1) {
    out[index] = Number.parseInt(hex.slice(index * 2, index * 2 + 2), 16);
  }
  return out;
}

function base58btcEncode(bytes: Uint8Array): string {
  const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz";
  let value = 0n;
  for (const byte of bytes) {
    value = (value << 8n) + BigInt(byte);
  }
  let encoded = "";
  while (value > 0n) {
    const mod = value % 58n;
    encoded = alphabet[Number(mod)] + encoded;
    value /= 58n;
  }
  for (const byte of bytes) {
    if (byte !== 0) break;
    encoded = alphabet[0] + encoded;
  }
  return encoded;
}

function formatDate(value: string): string {
  if (!value) return "unknown";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString();
}

function formatUnix(value: bigint): string {
  return new Date(Number(value) * 1000).toLocaleString();
}

function formatChainId(chainId: string): string {
  if (!chainId) return "";
  try {
    return String(Number.parseInt(chainId, 16));
  } catch {
    return chainId;
  }
}

function shortAddress(value: string): string {
  if (!value) return "n/a";
  return `${value.slice(0, 6)}...${value.slice(-4)}`;
}

function shortHex(value: string): string {
  return `${value.slice(0, 8)}...${value.slice(-6)}`;
}

function roleToLabel(value: bigint): RoleLabel {
  const index = Number(value);
  return ROLE_LABELS[index] ?? "None";
}

function errorMessage(err: unknown): string {
  return err instanceof Error ? err.message : String(err);
}

createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
