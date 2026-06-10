// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BitRegistry {
    enum Role {
        None,
        Contributor,
        Maintainer,
        Owner
    }

    enum PullRequestStatus {
        None,
        Open,
        Approved,
        Rejected,
        Closed
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

    struct PullRequest {
        uint256 id;
        uint256 targetRepoId;
        bytes32 targetBranch;
        uint256 sourceRepoId;
        bytes32 sourceBranch;
        bytes20 baseCommit;
        bytes20 sourceHeadCommit;
        address author;
        PullRequestStatus status;
        uint256 createdAt;
        uint256 updatedAt;
    }

    uint256 public nextRepoId = 1;
    uint256 public nextPullRequestId = 1;
    mapping(uint256 => Repo) private repos;
    mapping(uint256 => PullRequest) private pullRequests;
    mapping(uint256 => uint256[]) private repoPullRequests;

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
    error PullRequestNotFound();
    error PullRequestNotOpen();
    error EmptySourceBranch();
    error NoCommitsToMerge();
    error SourceHistoryDoesNotContainBase();
    error SourceHeadNotFound();
    error UnauthorizedPullRequestAction();
    error MergeCommitNotSupported();

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
    event PullRequestCreated(
        uint256 indexed prId,
        uint256 indexed targetRepoId,
        uint256 indexed sourceRepoId,
        bytes32 targetBranch,
        bytes32 sourceBranch,
        bytes20 baseCommit,
        bytes20 sourceHeadCommit,
        address author
    );
    event PullRequestApproved(
        uint256 indexed prId,
        uint256 indexed targetRepoId,
        bytes32 indexed targetBranch,
        bytes20 oldHead,
        bytes20 newHead,
        address approver
    );
    event PullRequestRejected(uint256 indexed prId, address indexed rejectedBy);
    event PullRequestClosed(uint256 indexed prId, address indexed closedBy);

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

    function _isMaintainer(uint256 repoId, address user) private view returns (bool) {
        Role role = repos[repoId].roles[user];
        return role == Role.Maintainer || role == Role.Owner;
    }

    function _requireOpenPullRequest(uint256 prId) private view returns (PullRequest storage pr) {
        pr = pullRequests[prId];
        if (pr.status == PullRequestStatus.None) revert PullRequestNotFound();
        if (pr.status != PullRequestStatus.Open) revert PullRequestNotOpen();
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

    function createPullRequest(
        uint256 targetRepoId,
        bytes32 targetBranch,
        uint256 sourceRepoId,
        bytes32 sourceBranch
    ) external repoExists(targetRepoId) repoExists(sourceRepoId) returns (uint256 prId) {
        Repo storage targetRepo = repos[targetRepoId];
        Repo storage sourceRepo = repos[sourceRepoId];
        bytes20 baseCommit = targetRepo.branchCommits[targetBranch];
        bytes20 sourceHeadCommit = sourceRepo.branchCommits[sourceBranch];
        if (sourceHeadCommit == bytes20(0)) revert EmptySourceBranch();

        _getPullRequestCommitRange(sourceRepo, sourceBranch, baseCommit, sourceHeadCommit);

        prId = nextPullRequestId++;
        pullRequests[prId] = PullRequest({
            id: prId,
            targetRepoId: targetRepoId,
            targetBranch: targetBranch,
            sourceRepoId: sourceRepoId,
            sourceBranch: sourceBranch,
            baseCommit: baseCommit,
            sourceHeadCommit: sourceHeadCommit,
            author: msg.sender,
            status: PullRequestStatus.Open,
            createdAt: block.timestamp,
            updatedAt: block.timestamp
        });
        repoPullRequests[targetRepoId].push(prId);

        emit PullRequestCreated(
            prId,
            targetRepoId,
            sourceRepoId,
            targetBranch,
            sourceBranch,
            baseCommit,
            sourceHeadCommit,
            msg.sender
        );
    }

    function approvePullRequest(uint256 prId) external {
        PullRequest storage pr = _requireOpenPullRequest(prId);
        if (!_isMaintainer(pr.targetRepoId, msg.sender)) revert MaintainerRequired();

        Repo storage targetRepo = repos[pr.targetRepoId];
        Repo storage sourceRepo = repos[pr.sourceRepoId];
        bytes20 oldHead = targetRepo.branchCommits[pr.targetBranch];
        if (oldHead != pr.baseCommit) revert StaleBranchHead();

        (uint256 start, uint256 end) =
            _getPullRequestCommitRange(sourceRepo, pr.sourceBranch, pr.baseCommit, pr.sourceHeadCommit);

        bytes20 currentHead = oldHead;
        bytes20 newHead = oldHead;
        bytes20[] storage sourceHistory = sourceRepo.branchHistory[pr.sourceBranch];
        for (uint256 i = start; i < end; i++) {
            bytes20 commitHash = sourceHistory[i];
            _copyCommitRecord(sourceRepo, targetRepo, commitHash);
            CommitRecord storage item = targetRepo.commits[commitHash];
            if (item.parents.length > 1) revert MergeCommitNotSupported();
            if (currentHead != bytes20(0)) {
                if (item.parents.length == 0) revert MissingParent();
                if (item.parents[0] != currentHead) revert FirstParentMismatch();
            }

            targetRepo.branchCommits[pr.targetBranch] = commitHash;
            targetRepo.branchHistory[pr.targetBranch].push(commitHash);
            newHead = commitHash;

            emit CommitRecorded(
                pr.targetRepoId,
                pr.targetBranch,
                commitHash,
                item.treeHash,
                _parentsToMemory(item),
                item.manifestDigest,
                item.diffDigest,
                msg.sender
            );
            emit BranchUpdated(
                pr.targetRepoId,
                pr.targetBranch,
                currentHead == bytes20(0) ? bytes("") : abi.encodePacked(targetRepo.commits[currentHead].manifestDigest),
                abi.encodePacked(item.manifestDigest),
                abi.encodePacked(commitHash),
                item.parents.length == 0 ? bytes("") : abi.encodePacked(item.parents[0]),
                msg.sender
            );

            currentHead = commitHash;
        }

        pr.status = PullRequestStatus.Approved;
        pr.updatedAt = block.timestamp;
        emit PullRequestApproved(prId, pr.targetRepoId, pr.targetBranch, oldHead, newHead, msg.sender);
    }

    function rejectPullRequest(uint256 prId) external {
        PullRequest storage pr = _requireOpenPullRequest(prId);
        if (!_isMaintainer(pr.targetRepoId, msg.sender)) revert MaintainerRequired();

        pr.status = PullRequestStatus.Rejected;
        pr.updatedAt = block.timestamp;
        emit PullRequestRejected(prId, msg.sender);
    }

    function closePullRequest(uint256 prId) external {
        PullRequest storage pr = _requireOpenPullRequest(prId);
        if (msg.sender != pr.author && !_isMaintainer(pr.targetRepoId, msg.sender)) {
            revert UnauthorizedPullRequestAction();
        }

        pr.status = PullRequestStatus.Closed;
        pr.updatedAt = block.timestamp;
        emit PullRequestClosed(prId, msg.sender);
    }

    function getPullRequest(uint256 prId) external view returns (PullRequest memory) {
        PullRequest storage pr = pullRequests[prId];
        if (pr.status == PullRequestStatus.None) revert PullRequestNotFound();
        return pr;
    }

    function getRepoPullRequestCount(uint256 repoId) external view repoExists(repoId) returns (uint256) {
        return repoPullRequests[repoId].length;
    }

    function getRepoPullRequestAt(uint256 repoId, uint256 index) external view repoExists(repoId) returns (uint256) {
        return repoPullRequests[repoId][index];
    }

    function _getPullRequestCommitRange(
        Repo storage sourceRepo,
        bytes32 sourceBranch,
        bytes20 baseCommit,
        bytes20 sourceHeadCommit
    ) private view returns (uint256 start, uint256 end) {
        bytes20[] storage history = sourceRepo.branchHistory[sourceBranch];
        if (sourceHeadCommit == bytes20(0)) revert EmptySourceBranch();

        bool foundBase = baseCommit == bytes20(0);
        bool foundHead;
        start = 0;
        for (uint256 i = 0; i < history.length; i++) {
            bytes20 commitHash = history[i];
            if (!foundBase && commitHash == baseCommit) {
                foundBase = true;
                start = i + 1;
                continue;
            }
            if (foundBase && commitHash == sourceHeadCommit) {
                foundHead = true;
                end = i + 1;
                break;
            }
        }

        if (!foundBase) revert SourceHistoryDoesNotContainBase();
        if (!foundHead) revert SourceHeadNotFound();
        if (start >= end) revert NoCommitsToMerge();
    }

    function _copyCommitRecord(Repo storage sourceRepo, Repo storage targetRepo, bytes20 commitHash) private {
        CommitRecord storage source = sourceRepo.commits[commitHash];
        if (!source.exists) revert CommitNotFound();

        CommitRecord storage target = targetRepo.commits[commitHash];
        if (target.exists) {
            if (
                target.treeHash != source.treeHash || target.manifestDigest != source.manifestDigest
                    || target.diffDigest != source.diffDigest || target.parents.length != source.parents.length
            ) {
                revert CommitMetadataMismatch();
            }
            for (uint256 i = 0; i < source.parents.length; i++) {
                if (target.parents[i] != source.parents[i]) revert CommitMetadataMismatch();
            }
            return;
        }

        target.treeHash = source.treeHash;
        target.manifestDigest = source.manifestDigest;
        target.diffDigest = source.diffDigest;
        target.updater = source.updater;
        target.timestamp = source.timestamp;
        target.exists = true;
        for (uint256 i = 0; i < source.parents.length; i++) {
            target.parents.push(source.parents[i]);
        }
    }

    function _parentsToMemory(CommitRecord storage item) private view returns (bytes20[] memory parents) {
        parents = new bytes20[](item.parents.length);
        for (uint256 i = 0; i < item.parents.length; i++) {
            parents[i] = item.parents[i];
        }
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
