// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BitRegistry {
    enum Role {
        None,
        Contributor,
        Maintainer,
        Owner
    }

    struct Repo {
        address owner;
        bytes metadataCID;
        mapping(address => Role) roles;
        mapping(bytes32 => bytes) branchHeads;
        mapping(bytes32 => BranchUpdate[]) branchHistory;
        mapping(bytes32 => bytes) tags;
        mapping(bytes32 => bool) tagExists;
    }

    struct BranchUpdate {
        bytes oldHead;
        bytes newHead;
        bytes gitCommit;
        bytes previousCommit;
        address updater;
        uint256 timestamp;
    }

    uint256 public nextRepoId = 1;
    mapping(uint256 => Repo) private repos;

    event RepoCreated(uint256 indexed repoId, address indexed owner, bytes metadataCID);
    event RoleChanged(uint256 indexed repoId, address indexed user, Role role);
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
        require(repos[repoId].owner != address(0), "repo not found");
        _;
    }

    modifier onlyOwner(uint256 repoId) {
        require(repos[repoId].roles[msg.sender] == Role.Owner, "owner required");
        _;
    }

    modifier onlyMaintainer(uint256 repoId) {
        Role role = repos[repoId].roles[msg.sender];
        require(role == Role.Maintainer || role == Role.Owner, "maintainer required");
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
        require(user != address(0), "zero user");
        repos[repoId].roles[user] = role;
        emit RoleChanged(repoId, user, role);
    }

    function getBranchHead(uint256 repoId, bytes32 branch) external view repoExists(repoId) returns (bytes memory) {
        return repos[repoId].branchHeads[branch];
    }

    function getBranchHistoryLength(uint256 repoId, bytes32 branch) external view repoExists(repoId) returns (uint256) {
        return repos[repoId].branchHistory[branch].length;
    }

    function getBranchHistoryAt(
        uint256 repoId,
        bytes32 branch,
        uint256 index
    )
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
        BranchUpdate storage item = repos[repoId].branchHistory[branch][index];
        return (item.oldHead, item.newHead, item.gitCommit, item.previousCommit, item.updater, item.timestamp);
    }

    function updateBranch(
        uint256 repoId,
        bytes32 branch,
        bytes calldata expectedOldHead,
        bytes calldata newHead,
        bytes calldata gitCommit,
        bytes calldata previousCommit
    ) external repoExists(repoId) onlyMaintainer(repoId) {
        Repo storage repo = repos[repoId];
        bytes memory current = repo.branchHeads[branch];
        require(keccak256(current) == keccak256(expectedOldHead), "stale branch head");
        repo.branchHeads[branch] = newHead;
        repo.branchHistory[branch].push(BranchUpdate({
            oldHead: current,
            newHead: newHead,
            gitCommit: gitCommit,
            previousCommit: previousCommit,
            updater: msg.sender,
            timestamp: block.timestamp
        }));
        emit BranchUpdated(repoId, branch, current, newHead, gitCommit, previousCommit, msg.sender);
    }

    function createTag(uint256 repoId, bytes32 tag, bytes calldata target) external repoExists(repoId) onlyMaintainer(repoId) {
        require(!repos[repoId].tagExists[tag], "tag exists");
        repos[repoId].tagExists[tag] = true;
        repos[repoId].tags[tag] = target;
        emit TagCreated(repoId, tag, target, msg.sender);
    }

    function getTag(uint256 repoId, bytes32 tag) external view repoExists(repoId) returns (bytes memory) {
        require(repos[repoId].tagExists[tag], "tag not found");
        return repos[repoId].tags[tag];
    }
}
