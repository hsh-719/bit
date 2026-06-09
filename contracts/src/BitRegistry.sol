// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BitRegistry {
    enum Role {
        None,
        Contributor,
        Maintainer,
        Owner
    }

    struct CommitRecord {
        bytes20 treeHash;
        bytes32 manifestDigest;
        bytes32 diffDigest;
        address updater;
        uint256 timestamp;
        bytes20[] parents;
        bool exists;
    }

    struct Repo {
        address owner;
        bytes metadataCID;
        mapping(address => Role) roles;
        mapping(bytes32 => bytes20) branchCommits;
        mapping(bytes32 => bytes20[]) branchHistory;
        mapping(bytes20 => CommitRecord) commits;
        mapping(bytes32 => bytes) tags;
        mapping(bytes32 => bool) tagExists;
    }

    uint256 public nextRepoId = 1;
    mapping(uint256 => Repo) private repos;

    error RepoNotFound();
    error OwnerRequired();
    error MaintainerRequired();
    error ZeroUser();
    error ZeroCommit();
    error StaleBranchHead();
    error MissingParent();
    error FirstParentMismatch();
    error CommitMetadataMismatch();
    error CommitNotFound();
    error TagExists();
    error TagNotFound();

    event RepoCreated(uint256 indexed repoId, address indexed owner, bytes metadataCID);
    event RoleChanged(uint256 indexed repoId, address indexed user, Role role);
    event CommitRecorded(
        uint256 indexed repoId,
        bytes32 indexed branch,
        bytes20 indexed commitHash,
        bytes20 treeHash,
        bytes20[] parents,
        bytes32 manifestDigest,
        bytes32 diffDigest,
        address updater
    );
    event BranchUpdated(
        uint256 indexed repoId,
        bytes32 indexed branch,
        bytes oldHead,
        bytes newHead,
        bytes gitCommit,
        bytes previousCommit,
        address indexed updater
    );
    event TagCreated(uint256 indexed repoId, bytes32 indexed tag, bytes target, address indexed creator);

    modifier repoExists(uint256 repoId) {
        if (repos[repoId].owner == address(0)) revert RepoNotFound();
        _;
    }

    modifier onlyOwner(uint256 repoId) {
        if (repos[repoId].roles[msg.sender] != Role.Owner) revert OwnerRequired();
        _;
    }

    modifier onlyMaintainer(uint256 repoId) {
        Role role = repos[repoId].roles[msg.sender];
        if (role != Role.Maintainer && role != Role.Owner) revert MaintainerRequired();
        _;
    }

    function createRepo(bytes calldata metadataCID) external returns (uint256 repoId) {
        repoId = nextRepoId++;
        Repo storage repo = repos[repoId];
        repo.owner = msg.sender;
        repo.metadataCID = metadataCID;
        repo.roles[msg.sender] = Role.Owner;
        emit RepoCreated(repoId, msg.sender, metadataCID);
        emit RoleChanged(repoId, msg.sender, Role.Owner);
    }

    function getRole(uint256 repoId, address user) external view repoExists(repoId) returns (Role) {
        return repos[repoId].roles[user];
    }

    function setRole(uint256 repoId, address user, Role role) external repoExists(repoId) onlyOwner(repoId) {
        if (user == address(0)) revert ZeroUser();
        repos[repoId].roles[user] = role;
        emit RoleChanged(repoId, user, role);
    }

    function getBranchHead(uint256 repoId, bytes32 branch) external view repoExists(repoId) returns (bytes memory) {
        bytes20 commitHash = repos[repoId].branchCommits[branch];
        if (commitHash == bytes20(0)) {
            return bytes("");
        }
        return abi.encodePacked(repos[repoId].commits[commitHash].manifestDigest);
    }

    function getBranchCommit(uint256 repoId, bytes32 branch) external view repoExists(repoId) returns (bytes20) {
        return repos[repoId].branchCommits[branch];
    }

    function getBranchHistoryLength(uint256 repoId, bytes32 branch) external view repoExists(repoId) returns (uint256) {
        return repos[repoId].branchHistory[branch].length;
    }

    function getBranchCommitAt(uint256 repoId, bytes32 branch, uint256 index)
        external
        view
        repoExists(repoId)
        returns (bytes20)
    {
        return repos[repoId].branchHistory[branch][index];
    }

    function getBranchCommitsWithMetadata(uint256 repoId, bytes32 branch, uint256 start, uint256 limit)
        external
        view
        repoExists(repoId)
        returns (
            bytes20[] memory commitHashes,
            bytes20[] memory treeHashes,
            bytes32[] memory manifestDigests,
            bytes32[] memory diffDigests
        )
    {
        bytes20[] storage history = repos[repoId].branchHistory[branch];
        if (start >= history.length || limit == 0) {
            return (new bytes20[](0), new bytes20[](0), new bytes32[](0), new bytes32[](0));
        }
        uint256 count = history.length - start;
        if (count > limit) {
            count = limit;
        }
        commitHashes = new bytes20[](count);
        treeHashes = new bytes20[](count);
        manifestDigests = new bytes32[](count);
        diffDigests = new bytes32[](count);
        for (uint256 i = 0; i < count; i++) {
            bytes20 commitHash = history[start + i];
            CommitRecord storage item = repos[repoId].commits[commitHash];
            commitHashes[i] = commitHash;
            treeHashes[i] = item.treeHash;
            manifestDigests[i] = item.manifestDigest;
            diffDigests[i] = item.diffDigest;
        }
    }

    function getCommit(uint256 repoId, bytes20 commitHash)
        external
        view
        repoExists(repoId)
        returns (bytes20 treeHash, bytes32 manifestDigest, bytes32 diffDigest, address updater, uint256 timestamp)
    {
        CommitRecord storage item = repos[repoId].commits[commitHash];
        if (!item.exists) revert CommitNotFound();
        return (item.treeHash, item.manifestDigest, item.diffDigest, item.updater, item.timestamp);
    }

    function getCommitParentCount(uint256 repoId, bytes20 commitHash)
        external
        view
        repoExists(repoId)
        returns (uint256)
    {
        CommitRecord storage item = repos[repoId].commits[commitHash];
        if (!item.exists) revert CommitNotFound();
        return item.parents.length;
    }

    function getCommitParentAt(uint256 repoId, bytes20 commitHash, uint256 index)
        external
        view
        repoExists(repoId)
        returns (bytes20)
    {
        CommitRecord storage item = repos[repoId].commits[commitHash];
        if (!item.exists) revert CommitNotFound();
        return item.parents[index];
    }

    function getBranchHistoryAt(uint256 repoId, bytes32 branch, uint256 index)
        external
        view
        repoExists(repoId)
        returns (
            bytes memory oldHead,
            bytes memory newHead,
            bytes memory gitCommit,
            bytes memory previousCommit,
            address updater,
            uint256 timestamp
        )
    {
        Repo storage repo = repos[repoId];
        bytes20 commitHash = repo.branchHistory[branch][index];
        CommitRecord storage item = repo.commits[commitHash];
        bytes20 previous;
        if (item.parents.length > 0) {
            previous = item.parents[0];
            oldHead = abi.encodePacked(repo.commits[previous].manifestDigest);
            previousCommit = abi.encodePacked(previous);
        }
        return (oldHead, abi.encodePacked(item.manifestDigest), abi.encodePacked(commitHash), previousCommit, item.updater, item.timestamp);
    }

    function recordCommit(
        uint256 repoId,
        bytes32 branch,
        bytes20 expectedOldCommit,
        bytes20 commitHash,
        bytes20 treeHash,
        bytes20[] calldata parents,
        bytes32 manifestDigest,
        bytes32 diffDigest
    ) external repoExists(repoId) onlyMaintainer(repoId) {
        if (commitHash == bytes20(0)) revert ZeroCommit();

        Repo storage repo = repos[repoId];
        bytes20 currentCommit = repo.branchCommits[branch];
        if (currentCommit != expectedOldCommit) revert StaleBranchHead();
        if (currentCommit != bytes20(0)) {
            if (parents.length == 0) revert MissingParent();
            if (parents[0] != currentCommit) revert FirstParentMismatch();
        }

        CommitRecord storage item = repo.commits[commitHash];
        if (item.exists) {
            if (item.treeHash != treeHash) revert CommitMetadataMismatch();
        } else {
            item.treeHash = treeHash;
            item.manifestDigest = manifestDigest;
            item.diffDigest = diffDigest;
            item.updater = msg.sender;
            item.timestamp = block.timestamp;
            item.exists = true;
            for (uint256 i = 0; i < parents.length; i++) {
                item.parents.push(parents[i]);
            }
        }

        bytes20 oldCommit = repo.branchCommits[branch];
        bytes memory oldHead = oldCommit == bytes20(0) ? bytes("") : abi.encodePacked(repo.commits[oldCommit].manifestDigest);
        bytes memory previousCommit = parents.length == 0 ? bytes("") : abi.encodePacked(parents[0]);
        repo.branchCommits[branch] = commitHash;
        repo.branchHistory[branch].push(commitHash);

        emit CommitRecorded(repoId, branch, commitHash, treeHash, parents, manifestDigest, diffDigest, msg.sender);
        emit BranchUpdated(
            repoId,
            branch,
            oldHead,
            abi.encodePacked(manifestDigest),
            abi.encodePacked(commitHash),
            previousCommit,
            msg.sender
        );
    }

    function createTag(uint256 repoId, bytes32 tag, bytes calldata target) external repoExists(repoId) onlyMaintainer(repoId) {
        if (repos[repoId].tagExists[tag]) revert TagExists();
        repos[repoId].tagExists[tag] = true;
        repos[repoId].tags[tag] = target;
        emit TagCreated(repoId, tag, target, msg.sender);
    }

    function getTag(uint256 repoId, bytes32 tag) external view repoExists(repoId) returns (bytes memory) {
        if (!repos[repoId].tagExists[tag]) revert TagNotFound();
        return repos[repoId].tags[tag];
    }
}
