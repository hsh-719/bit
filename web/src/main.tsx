import React, { useMemo, useState } from "react";
import { createRoot } from "react-dom/client";
import { createPublicClient, http, keccak256, stringToBytes, type Address, type Hex } from "viem";
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

type Manifest = {
  gitCommit: string;
  treeHash: string;
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

const abi = bitRegistryArtifact.abi;
const defaultRpcURL = localStorage.getItem("bit.rpcURL") ?? "http://127.0.0.1:8545";
const defaultContract = localStorage.getItem("bit.contract") ?? "";
const defaultGateway = localStorage.getItem("bit.ipfsGateway") ?? "http://127.0.0.1:8080/ipfs";

function App() {
  const [rpcURL, setRpcURL] = useState(defaultRpcURL);
  const [contractAddress, setContractAddress] = useState(defaultContract);
  const [ipfsGateway, setIpfsGateway] = useState(defaultGateway);
  const [branch, setBranch] = useState("main");
  const [repos, setRepos] = useState<RepoSummary[]>([]);
  const [selectedRepoId, setSelectedRepoId] = useState<bigint | null>(null);
  const [commits, setCommits] = useState<CommitSummary[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const client = useMemo(() => {
    return createPublicClient({ transport: http(rpcURL) });
  }, [rpcURL]);

  const selectedRepo = repos.find((repo) => repo.id === selectedRepoId) ?? null;

  async function loadRepos() {
    setLoading(true);
    setError("");
    setCommits([]);
    try {
      const address = parseAddress(contractAddress);
      localStorage.setItem("bit.rpcURL", rpcURL);
      localStorage.setItem("bit.contract", contractAddress);
      localStorage.setItem("bit.ipfsGateway", ipfsGateway);

      const count = await client.readContract({
        address,
        abi,
        functionName: "getRepoCount",
      }) as bigint;

      const nextRepos: RepoSummary[] = [];
      for (let repoId = 1n; repoId <= count; repoId++) {
        const [owner, metadataBytes] = await client.readContract({
          address,
          abi,
          functionName: "getRepo",
          args: [repoId],
        }) as [Address, Hex];
        const metadataCID = bytesHexToString(metadataBytes);
        const metadata = metadataCID ? await fetchJson<RepoMetadata>(ipfsURL(ipfsGateway, metadataCID)) : null;
        nextRepos.push({ id: repoId, owner, metadataCID, metadata });
      }
      setRepos(nextRepos);
      if (nextRepos.length > 0) {
        setSelectedRepoId(nextRepos[0].id);
      }
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  async function loadCommits(repoId: bigint) {
    setLoading(true);
    setError("");
    try {
      const address = parseAddress(contractAddress);
      const branchKey = keccak256(stringToBytes(branch));
      const length = await client.readContract({
        address,
        abi,
        functionName: "getBranchHistoryLength",
        args: [repoId, branchKey],
      }) as bigint;
      const pageSize = length > 50n ? 50n : length;
      const start = length > pageSize ? length - pageSize : 0n;
      const [hashes, treeHashes, manifestDigests] = await client.readContract({
        address,
        abi,
        functionName: "getBranchCommitsWithMetadata",
        args: [repoId, branchKey, start, pageSize],
      }) as [Hex[], Hex[], Hex[], Hex[]];

      const nextCommits = await Promise.all(hashes.map(async (hash, index) => {
        const [,, , updater, chainTimestamp] = await client.readContract({
          address,
          abi,
          functionName: "getCommit",
          args: [repoId, hash],
        }) as [Hex, Hex, Hex, Address, bigint];
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
      }));
      setSelectedRepoId(repoId);
      setCommits(nextCommits.reverse());
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="app">
      <section className="toolbar">
        <div>
          <h1>bit explorer</h1>
          <p>Read-only repository and commit metadata</p>
        </div>
        <button type="button" onClick={loadRepos} disabled={loading}>
          {loading ? "Loading" : "Load repositories"}
        </button>
      </section>

      <section className="settings">
        <label>
          RPC URL
          <input value={rpcURL} onChange={(event) => setRpcURL(event.target.value)} />
        </label>
        <label>
          Contract
          <input value={contractAddress} onChange={(event) => setContractAddress(event.target.value)} placeholder="0x..." />
        </label>
        <label>
          IPFS Gateway
          <input value={ipfsGateway} onChange={(event) => setIpfsGateway(event.target.value)} />
        </label>
        <label>
          Branch
          <input value={branch} onChange={(event) => setBranch(event.target.value)} />
        </label>
      </section>

      {error && <p className="error">{error}</p>}

      <section className="content">
        <aside className="repos">
          <div className="sectionTitle">Repositories</div>
          {repos.map((repo) => (
            <button
              className={repo.id === selectedRepoId ? "repo selected" : "repo"}
              key={repo.id.toString()}
              type="button"
              onClick={() => loadCommits(repo.id)}
            >
              <span className="repoName">{repo.metadata?.name || `Repo #${repo.id}`}</span>
              <span className="mono">#{repo.id.toString()}</span>
              <span className="owner">{shortAddress(repo.owner)}</span>
            </button>
          ))}
        </aside>

        <section className="history">
          <div className="historyHeader">
            <div>
              <div className="sectionTitle">Commit History</div>
              <h2>{selectedRepo?.metadata?.name || (selectedRepo ? `Repo #${selectedRepo.id}` : "No repository selected")}</h2>
              {selectedRepo?.metadata?.description && <p>{selectedRepo.metadata.description}</p>}
            </div>
            {selectedRepo && (
              <button type="button" onClick={() => loadCommits(selectedRepo.id)} disabled={loading}>
                Refresh
              </button>
            )}
          </div>

          <div className="commitList">
            {commits.map((commit) => (
              <article className="commit" key={commit.hash}>
                <div>
                  <h3>{commit.message || "(no message)"}</h3>
                  <p>
                    {commit.authorName || "Unknown author"}
                    {commit.authorEmail ? ` <${commit.authorEmail}>` : ""}
                  </p>
                </div>
                <dl>
                  <dt>Commit</dt>
                  <dd className="mono">{commit.hash}</dd>
                  <dt>Author Date</dt>
                  <dd>{formatDate(commit.authorDate)}</dd>
                  <dt>Recorded By</dt>
                  <dd className="mono">{shortAddress(commit.updater)}</dd>
                  <dt>On-chain Time</dt>
                  <dd>{formatUnix(commit.chainTimestamp)}</dd>
                  <dt>Parents</dt>
                  <dd className="mono">{commit.parents.length > 0 ? commit.parents.join(", ") : "none"}</dd>
                </dl>
              </article>
            ))}
          </div>
        </section>
      </section>
    </main>
  );
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
  for (let i = 0; i < hex.length; i += 2) {
    const code = Number.parseInt(hex.slice(i, i + 2), 16);
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

function cidV0FromDigest(digestHex: Hex): string {
  const digest = hexToBytes(digestHex);
  return base58btcEncode(new Uint8Array([0x12, 0x20, ...digest]));
}

function hexToBytes(value: Hex): Uint8Array {
  const hex = value.startsWith("0x") ? value.slice(2) : value;
  const out = new Uint8Array(hex.length / 2);
  for (let i = 0; i < out.length; i++) {
    out[i] = Number.parseInt(hex.slice(i * 2, i * 2 + 2), 16);
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
    value = value / 58n;
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

function shortAddress(value: string): string {
  return `${value.slice(0, 6)}...${value.slice(-4)}`;
}

function errorMessage(err: unknown): string {
  return err instanceof Error ? err.message : String(err);
}

createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
