// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package chain

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ChainMetaData contains all meta data concerning the Chain contract.
var ChainMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"repoId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"branch\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"oldHead\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"newHead\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"gitCommit\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"previousCommit\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"updater\",\"type\":\"address\"}],\"name\":\"BranchUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"repoId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"metadataCID\",\"type\":\"bytes\"}],\"name\":\"RepoCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"repoId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumBitRegistry.Role\",\"name\":\"role\",\"type\":\"uint8\"}],\"name\":\"RoleChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"repoId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"tag\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"target\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"}],\"name\":\"TagCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"metadataCID\",\"type\":\"bytes\"}],\"name\":\"createRepo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"repoId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"repoId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"tag\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"target\",\"type\":\"bytes\"}],\"name\":\"createTag\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"repoId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"branch\",\"type\":\"bytes32\"}],\"name\":\"getBranchHead\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"repoId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"branch\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getBranchHistoryAt\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"oldHead\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"newHead\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"gitCommit\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"previousCommit\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"updater\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"repoId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"branch\",\"type\":\"bytes32\"}],\"name\":\"getBranchHistoryLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"repoId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getRole\",\"outputs\":[{\"internalType\":\"enumBitRegistry.Role\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"repoId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"tag\",\"type\":\"bytes32\"}],\"name\":\"getTag\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextRepoId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"repoId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"enumBitRegistry.Role\",\"name\":\"role\",\"type\":\"uint8\"}],\"name\":\"setRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"repoId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"branch\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"expectedOldHead\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"newHead\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"gitCommit\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"previousCommit\",\"type\":\"bytes\"}],\"name\":\"updateBranch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ChainABI is the input ABI used to generate the binding from.
// Deprecated: Use ChainMetaData.ABI instead.
var ChainABI = ChainMetaData.ABI

// Chain is an auto generated Go binding around an Ethereum contract.
type Chain struct {
	ChainCaller     // Read-only binding to the contract
	ChainTransactor // Write-only binding to the contract
	ChainFilterer   // Log filterer for contract events
}

// ChainCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChainSession struct {
	Contract     *Chain            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChainCallerSession struct {
	Contract *ChainCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ChainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChainTransactorSession struct {
	Contract     *ChainTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChainRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChainRaw struct {
	Contract *Chain // Generic contract binding to access the raw methods on
}

// ChainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChainCallerRaw struct {
	Contract *ChainCaller // Generic read-only contract binding to access the raw methods on
}

// ChainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChainTransactorRaw struct {
	Contract *ChainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChain creates a new instance of Chain, bound to a specific deployed contract.
func NewChain(address common.Address, backend bind.ContractBackend) (*Chain, error) {
	contract, err := bindChain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Chain{ChainCaller: ChainCaller{contract: contract}, ChainTransactor: ChainTransactor{contract: contract}, ChainFilterer: ChainFilterer{contract: contract}}, nil
}

// NewChainCaller creates a new read-only instance of Chain, bound to a specific deployed contract.
func NewChainCaller(address common.Address, caller bind.ContractCaller) (*ChainCaller, error) {
	contract, err := bindChain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChainCaller{contract: contract}, nil
}

// NewChainTransactor creates a new write-only instance of Chain, bound to a specific deployed contract.
func NewChainTransactor(address common.Address, transactor bind.ContractTransactor) (*ChainTransactor, error) {
	contract, err := bindChain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChainTransactor{contract: contract}, nil
}

// NewChainFilterer creates a new log filterer instance of Chain, bound to a specific deployed contract.
func NewChainFilterer(address common.Address, filterer bind.ContractFilterer) (*ChainFilterer, error) {
	contract, err := bindChain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChainFilterer{contract: contract}, nil
}

// bindChain binds a generic wrapper to an already deployed contract.
func bindChain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ChainMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chain *ChainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chain.Contract.ChainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chain *ChainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chain.Contract.ChainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chain *ChainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chain.Contract.ChainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chain *ChainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chain *ChainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chain *ChainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chain.Contract.contract.Transact(opts, method, params...)
}

// GetBranchHead is a free data retrieval call binding the contract method 0x3e3bfbe0.
//
// Solidity: function getBranchHead(uint256 repoId, bytes32 branch) view returns(bytes)
func (_Chain *ChainCaller) GetBranchHead(opts *bind.CallOpts, repoId *big.Int, branch [32]byte) ([]byte, error) {
	var out []interface{}
	err := _Chain.contract.Call(opts, &out, "getBranchHead", repoId, branch)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetBranchHead is a free data retrieval call binding the contract method 0x3e3bfbe0.
//
// Solidity: function getBranchHead(uint256 repoId, bytes32 branch) view returns(bytes)
func (_Chain *ChainSession) GetBranchHead(repoId *big.Int, branch [32]byte) ([]byte, error) {
	return _Chain.Contract.GetBranchHead(&_Chain.CallOpts, repoId, branch)
}

// GetBranchHead is a free data retrieval call binding the contract method 0x3e3bfbe0.
//
// Solidity: function getBranchHead(uint256 repoId, bytes32 branch) view returns(bytes)
func (_Chain *ChainCallerSession) GetBranchHead(repoId *big.Int, branch [32]byte) ([]byte, error) {
	return _Chain.Contract.GetBranchHead(&_Chain.CallOpts, repoId, branch)
}

// GetBranchHistoryAt is a free data retrieval call binding the contract method 0x1700cceb.
//
// Solidity: function getBranchHistoryAt(uint256 repoId, bytes32 branch, uint256 index) view returns(bytes oldHead, bytes newHead, bytes gitCommit, bytes previousCommit, address updater, uint256 timestamp)
func (_Chain *ChainCaller) GetBranchHistoryAt(opts *bind.CallOpts, repoId *big.Int, branch [32]byte, index *big.Int) (struct {
	OldHead        []byte
	NewHead        []byte
	GitCommit      []byte
	PreviousCommit []byte
	Updater        common.Address
	Timestamp      *big.Int
}, error) {
	var out []interface{}
	err := _Chain.contract.Call(opts, &out, "getBranchHistoryAt", repoId, branch, index)

	outstruct := new(struct {
		OldHead        []byte
		NewHead        []byte
		GitCommit      []byte
		PreviousCommit []byte
		Updater        common.Address
		Timestamp      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.OldHead = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.NewHead = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.GitCommit = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.PreviousCommit = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.Updater = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Timestamp = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetBranchHistoryAt is a free data retrieval call binding the contract method 0x1700cceb.
//
// Solidity: function getBranchHistoryAt(uint256 repoId, bytes32 branch, uint256 index) view returns(bytes oldHead, bytes newHead, bytes gitCommit, bytes previousCommit, address updater, uint256 timestamp)
func (_Chain *ChainSession) GetBranchHistoryAt(repoId *big.Int, branch [32]byte, index *big.Int) (struct {
	OldHead        []byte
	NewHead        []byte
	GitCommit      []byte
	PreviousCommit []byte
	Updater        common.Address
	Timestamp      *big.Int
}, error) {
	return _Chain.Contract.GetBranchHistoryAt(&_Chain.CallOpts, repoId, branch, index)
}

// GetBranchHistoryAt is a free data retrieval call binding the contract method 0x1700cceb.
//
// Solidity: function getBranchHistoryAt(uint256 repoId, bytes32 branch, uint256 index) view returns(bytes oldHead, bytes newHead, bytes gitCommit, bytes previousCommit, address updater, uint256 timestamp)
func (_Chain *ChainCallerSession) GetBranchHistoryAt(repoId *big.Int, branch [32]byte, index *big.Int) (struct {
	OldHead        []byte
	NewHead        []byte
	GitCommit      []byte
	PreviousCommit []byte
	Updater        common.Address
	Timestamp      *big.Int
}, error) {
	return _Chain.Contract.GetBranchHistoryAt(&_Chain.CallOpts, repoId, branch, index)
}

// GetBranchHistoryLength is a free data retrieval call binding the contract method 0x8e00061b.
//
// Solidity: function getBranchHistoryLength(uint256 repoId, bytes32 branch) view returns(uint256)
func (_Chain *ChainCaller) GetBranchHistoryLength(opts *bind.CallOpts, repoId *big.Int, branch [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Chain.contract.Call(opts, &out, "getBranchHistoryLength", repoId, branch)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBranchHistoryLength is a free data retrieval call binding the contract method 0x8e00061b.
//
// Solidity: function getBranchHistoryLength(uint256 repoId, bytes32 branch) view returns(uint256)
func (_Chain *ChainSession) GetBranchHistoryLength(repoId *big.Int, branch [32]byte) (*big.Int, error) {
	return _Chain.Contract.GetBranchHistoryLength(&_Chain.CallOpts, repoId, branch)
}

// GetBranchHistoryLength is a free data retrieval call binding the contract method 0x8e00061b.
//
// Solidity: function getBranchHistoryLength(uint256 repoId, bytes32 branch) view returns(uint256)
func (_Chain *ChainCallerSession) GetBranchHistoryLength(repoId *big.Int, branch [32]byte) (*big.Int, error) {
	return _Chain.Contract.GetBranchHistoryLength(&_Chain.CallOpts, repoId, branch)
}

// GetRole is a free data retrieval call binding the contract method 0x93552a3d.
//
// Solidity: function getRole(uint256 repoId, address user) view returns(uint8)
func (_Chain *ChainCaller) GetRole(opts *bind.CallOpts, repoId *big.Int, user common.Address) (uint8, error) {
	var out []interface{}
	err := _Chain.contract.Call(opts, &out, "getRole", repoId, user)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetRole is a free data retrieval call binding the contract method 0x93552a3d.
//
// Solidity: function getRole(uint256 repoId, address user) view returns(uint8)
func (_Chain *ChainSession) GetRole(repoId *big.Int, user common.Address) (uint8, error) {
	return _Chain.Contract.GetRole(&_Chain.CallOpts, repoId, user)
}

// GetRole is a free data retrieval call binding the contract method 0x93552a3d.
//
// Solidity: function getRole(uint256 repoId, address user) view returns(uint8)
func (_Chain *ChainCallerSession) GetRole(repoId *big.Int, user common.Address) (uint8, error) {
	return _Chain.Contract.GetRole(&_Chain.CallOpts, repoId, user)
}

// GetTag is a free data retrieval call binding the contract method 0x3104b124.
//
// Solidity: function getTag(uint256 repoId, bytes32 tag) view returns(bytes)
func (_Chain *ChainCaller) GetTag(opts *bind.CallOpts, repoId *big.Int, tag [32]byte) ([]byte, error) {
	var out []interface{}
	err := _Chain.contract.Call(opts, &out, "getTag", repoId, tag)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetTag is a free data retrieval call binding the contract method 0x3104b124.
//
// Solidity: function getTag(uint256 repoId, bytes32 tag) view returns(bytes)
func (_Chain *ChainSession) GetTag(repoId *big.Int, tag [32]byte) ([]byte, error) {
	return _Chain.Contract.GetTag(&_Chain.CallOpts, repoId, tag)
}

// GetTag is a free data retrieval call binding the contract method 0x3104b124.
//
// Solidity: function getTag(uint256 repoId, bytes32 tag) view returns(bytes)
func (_Chain *ChainCallerSession) GetTag(repoId *big.Int, tag [32]byte) ([]byte, error) {
	return _Chain.Contract.GetTag(&_Chain.CallOpts, repoId, tag)
}

// NextRepoId is a free data retrieval call binding the contract method 0x51734fc5.
//
// Solidity: function nextRepoId() view returns(uint256)
func (_Chain *ChainCaller) NextRepoId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Chain.contract.Call(opts, &out, "nextRepoId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextRepoId is a free data retrieval call binding the contract method 0x51734fc5.
//
// Solidity: function nextRepoId() view returns(uint256)
func (_Chain *ChainSession) NextRepoId() (*big.Int, error) {
	return _Chain.Contract.NextRepoId(&_Chain.CallOpts)
}

// NextRepoId is a free data retrieval call binding the contract method 0x51734fc5.
//
// Solidity: function nextRepoId() view returns(uint256)
func (_Chain *ChainCallerSession) NextRepoId() (*big.Int, error) {
	return _Chain.Contract.NextRepoId(&_Chain.CallOpts)
}

// CreateRepo is a paid mutator transaction binding the contract method 0xb1096ace.
//
// Solidity: function createRepo(bytes metadataCID) returns(uint256 repoId)
func (_Chain *ChainTransactor) CreateRepo(opts *bind.TransactOpts, metadataCID []byte) (*types.Transaction, error) {
	return _Chain.contract.Transact(opts, "createRepo", metadataCID)
}

// CreateRepo is a paid mutator transaction binding the contract method 0xb1096ace.
//
// Solidity: function createRepo(bytes metadataCID) returns(uint256 repoId)
func (_Chain *ChainSession) CreateRepo(metadataCID []byte) (*types.Transaction, error) {
	return _Chain.Contract.CreateRepo(&_Chain.TransactOpts, metadataCID)
}

// CreateRepo is a paid mutator transaction binding the contract method 0xb1096ace.
//
// Solidity: function createRepo(bytes metadataCID) returns(uint256 repoId)
func (_Chain *ChainTransactorSession) CreateRepo(metadataCID []byte) (*types.Transaction, error) {
	return _Chain.Contract.CreateRepo(&_Chain.TransactOpts, metadataCID)
}

// CreateTag is a paid mutator transaction binding the contract method 0xaa4731a4.
//
// Solidity: function createTag(uint256 repoId, bytes32 tag, bytes target) returns()
func (_Chain *ChainTransactor) CreateTag(opts *bind.TransactOpts, repoId *big.Int, tag [32]byte, target []byte) (*types.Transaction, error) {
	return _Chain.contract.Transact(opts, "createTag", repoId, tag, target)
}

// CreateTag is a paid mutator transaction binding the contract method 0xaa4731a4.
//
// Solidity: function createTag(uint256 repoId, bytes32 tag, bytes target) returns()
func (_Chain *ChainSession) CreateTag(repoId *big.Int, tag [32]byte, target []byte) (*types.Transaction, error) {
	return _Chain.Contract.CreateTag(&_Chain.TransactOpts, repoId, tag, target)
}

// CreateTag is a paid mutator transaction binding the contract method 0xaa4731a4.
//
// Solidity: function createTag(uint256 repoId, bytes32 tag, bytes target) returns()
func (_Chain *ChainTransactorSession) CreateTag(repoId *big.Int, tag [32]byte, target []byte) (*types.Transaction, error) {
	return _Chain.Contract.CreateTag(&_Chain.TransactOpts, repoId, tag, target)
}

// SetRole is a paid mutator transaction binding the contract method 0x5aaa3abc.
//
// Solidity: function setRole(uint256 repoId, address user, uint8 role) returns()
func (_Chain *ChainTransactor) SetRole(opts *bind.TransactOpts, repoId *big.Int, user common.Address, role uint8) (*types.Transaction, error) {
	return _Chain.contract.Transact(opts, "setRole", repoId, user, role)
}

// SetRole is a paid mutator transaction binding the contract method 0x5aaa3abc.
//
// Solidity: function setRole(uint256 repoId, address user, uint8 role) returns()
func (_Chain *ChainSession) SetRole(repoId *big.Int, user common.Address, role uint8) (*types.Transaction, error) {
	return _Chain.Contract.SetRole(&_Chain.TransactOpts, repoId, user, role)
}

// SetRole is a paid mutator transaction binding the contract method 0x5aaa3abc.
//
// Solidity: function setRole(uint256 repoId, address user, uint8 role) returns()
func (_Chain *ChainTransactorSession) SetRole(repoId *big.Int, user common.Address, role uint8) (*types.Transaction, error) {
	return _Chain.Contract.SetRole(&_Chain.TransactOpts, repoId, user, role)
}

// UpdateBranch is a paid mutator transaction binding the contract method 0xa24b3f41.
//
// Solidity: function updateBranch(uint256 repoId, bytes32 branch, bytes expectedOldHead, bytes newHead, bytes gitCommit, bytes previousCommit) returns()
func (_Chain *ChainTransactor) UpdateBranch(opts *bind.TransactOpts, repoId *big.Int, branch [32]byte, expectedOldHead []byte, newHead []byte, gitCommit []byte, previousCommit []byte) (*types.Transaction, error) {
	return _Chain.contract.Transact(opts, "updateBranch", repoId, branch, expectedOldHead, newHead, gitCommit, previousCommit)
}

// UpdateBranch is a paid mutator transaction binding the contract method 0xa24b3f41.
//
// Solidity: function updateBranch(uint256 repoId, bytes32 branch, bytes expectedOldHead, bytes newHead, bytes gitCommit, bytes previousCommit) returns()
func (_Chain *ChainSession) UpdateBranch(repoId *big.Int, branch [32]byte, expectedOldHead []byte, newHead []byte, gitCommit []byte, previousCommit []byte) (*types.Transaction, error) {
	return _Chain.Contract.UpdateBranch(&_Chain.TransactOpts, repoId, branch, expectedOldHead, newHead, gitCommit, previousCommit)
}

// UpdateBranch is a paid mutator transaction binding the contract method 0xa24b3f41.
//
// Solidity: function updateBranch(uint256 repoId, bytes32 branch, bytes expectedOldHead, bytes newHead, bytes gitCommit, bytes previousCommit) returns()
func (_Chain *ChainTransactorSession) UpdateBranch(repoId *big.Int, branch [32]byte, expectedOldHead []byte, newHead []byte, gitCommit []byte, previousCommit []byte) (*types.Transaction, error) {
	return _Chain.Contract.UpdateBranch(&_Chain.TransactOpts, repoId, branch, expectedOldHead, newHead, gitCommit, previousCommit)
}

// ChainBranchUpdatedIterator is returned from FilterBranchUpdated and is used to iterate over the raw logs and unpacked data for BranchUpdated events raised by the Chain contract.
type ChainBranchUpdatedIterator struct {
	Event *ChainBranchUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChainBranchUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainBranchUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChainBranchUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChainBranchUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChainBranchUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChainBranchUpdated represents a BranchUpdated event raised by the Chain contract.
type ChainBranchUpdated struct {
	RepoId         *big.Int
	Branch         [32]byte
	OldHead        []byte
	NewHead        []byte
	GitCommit      []byte
	PreviousCommit []byte
	Updater        common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterBranchUpdated is a free log retrieval operation binding the contract event 0xb43f28d3a34853ca9427870699d0723b06ba455fa0565b8fcf9bfca57d710081.
//
// Solidity: event BranchUpdated(uint256 indexed repoId, bytes32 indexed branch, bytes oldHead, bytes newHead, bytes gitCommit, bytes previousCommit, address indexed updater)
func (_Chain *ChainFilterer) FilterBranchUpdated(opts *bind.FilterOpts, repoId []*big.Int, branch [][32]byte, updater []common.Address) (*ChainBranchUpdatedIterator, error) {

	var repoIdRule []interface{}
	for _, repoIdItem := range repoId {
		repoIdRule = append(repoIdRule, repoIdItem)
	}
	var branchRule []interface{}
	for _, branchItem := range branch {
		branchRule = append(branchRule, branchItem)
	}

	var updaterRule []interface{}
	for _, updaterItem := range updater {
		updaterRule = append(updaterRule, updaterItem)
	}

	logs, sub, err := _Chain.contract.FilterLogs(opts, "BranchUpdated", repoIdRule, branchRule, updaterRule)
	if err != nil {
		return nil, err
	}
	return &ChainBranchUpdatedIterator{contract: _Chain.contract, event: "BranchUpdated", logs: logs, sub: sub}, nil
}

// WatchBranchUpdated is a free log subscription operation binding the contract event 0xb43f28d3a34853ca9427870699d0723b06ba455fa0565b8fcf9bfca57d710081.
//
// Solidity: event BranchUpdated(uint256 indexed repoId, bytes32 indexed branch, bytes oldHead, bytes newHead, bytes gitCommit, bytes previousCommit, address indexed updater)
func (_Chain *ChainFilterer) WatchBranchUpdated(opts *bind.WatchOpts, sink chan<- *ChainBranchUpdated, repoId []*big.Int, branch [][32]byte, updater []common.Address) (event.Subscription, error) {

	var repoIdRule []interface{}
	for _, repoIdItem := range repoId {
		repoIdRule = append(repoIdRule, repoIdItem)
	}
	var branchRule []interface{}
	for _, branchItem := range branch {
		branchRule = append(branchRule, branchItem)
	}

	var updaterRule []interface{}
	for _, updaterItem := range updater {
		updaterRule = append(updaterRule, updaterItem)
	}

	logs, sub, err := _Chain.contract.WatchLogs(opts, "BranchUpdated", repoIdRule, branchRule, updaterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChainBranchUpdated)
				if err := _Chain.contract.UnpackLog(event, "BranchUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBranchUpdated is a log parse operation binding the contract event 0xb43f28d3a34853ca9427870699d0723b06ba455fa0565b8fcf9bfca57d710081.
//
// Solidity: event BranchUpdated(uint256 indexed repoId, bytes32 indexed branch, bytes oldHead, bytes newHead, bytes gitCommit, bytes previousCommit, address indexed updater)
func (_Chain *ChainFilterer) ParseBranchUpdated(log types.Log) (*ChainBranchUpdated, error) {
	event := new(ChainBranchUpdated)
	if err := _Chain.contract.UnpackLog(event, "BranchUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChainRepoCreatedIterator is returned from FilterRepoCreated and is used to iterate over the raw logs and unpacked data for RepoCreated events raised by the Chain contract.
type ChainRepoCreatedIterator struct {
	Event *ChainRepoCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChainRepoCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainRepoCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChainRepoCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChainRepoCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChainRepoCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChainRepoCreated represents a RepoCreated event raised by the Chain contract.
type ChainRepoCreated struct {
	RepoId      *big.Int
	Owner       common.Address
	MetadataCID []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRepoCreated is a free log retrieval operation binding the contract event 0x9084d44a05e074c4b623c67dcdd18849edb840144de35f5df63873048fc31d1e.
//
// Solidity: event RepoCreated(uint256 indexed repoId, address indexed owner, bytes metadataCID)
func (_Chain *ChainFilterer) FilterRepoCreated(opts *bind.FilterOpts, repoId []*big.Int, owner []common.Address) (*ChainRepoCreatedIterator, error) {

	var repoIdRule []interface{}
	for _, repoIdItem := range repoId {
		repoIdRule = append(repoIdRule, repoIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Chain.contract.FilterLogs(opts, "RepoCreated", repoIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &ChainRepoCreatedIterator{contract: _Chain.contract, event: "RepoCreated", logs: logs, sub: sub}, nil
}

// WatchRepoCreated is a free log subscription operation binding the contract event 0x9084d44a05e074c4b623c67dcdd18849edb840144de35f5df63873048fc31d1e.
//
// Solidity: event RepoCreated(uint256 indexed repoId, address indexed owner, bytes metadataCID)
func (_Chain *ChainFilterer) WatchRepoCreated(opts *bind.WatchOpts, sink chan<- *ChainRepoCreated, repoId []*big.Int, owner []common.Address) (event.Subscription, error) {

	var repoIdRule []interface{}
	for _, repoIdItem := range repoId {
		repoIdRule = append(repoIdRule, repoIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _Chain.contract.WatchLogs(opts, "RepoCreated", repoIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChainRepoCreated)
				if err := _Chain.contract.UnpackLog(event, "RepoCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRepoCreated is a log parse operation binding the contract event 0x9084d44a05e074c4b623c67dcdd18849edb840144de35f5df63873048fc31d1e.
//
// Solidity: event RepoCreated(uint256 indexed repoId, address indexed owner, bytes metadataCID)
func (_Chain *ChainFilterer) ParseRepoCreated(log types.Log) (*ChainRepoCreated, error) {
	event := new(ChainRepoCreated)
	if err := _Chain.contract.UnpackLog(event, "RepoCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChainRoleChangedIterator is returned from FilterRoleChanged and is used to iterate over the raw logs and unpacked data for RoleChanged events raised by the Chain contract.
type ChainRoleChangedIterator struct {
	Event *ChainRoleChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChainRoleChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainRoleChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChainRoleChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChainRoleChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChainRoleChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChainRoleChanged represents a RoleChanged event raised by the Chain contract.
type ChainRoleChanged struct {
	RepoId *big.Int
	User   common.Address
	Role   uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRoleChanged is a free log retrieval operation binding the contract event 0xc98b3c1f2e93a069aa9220f35ddceb1fc071af0645ca6b79b9436a5b139447d9.
//
// Solidity: event RoleChanged(uint256 indexed repoId, address indexed user, uint8 role)
func (_Chain *ChainFilterer) FilterRoleChanged(opts *bind.FilterOpts, repoId []*big.Int, user []common.Address) (*ChainRoleChangedIterator, error) {

	var repoIdRule []interface{}
	for _, repoIdItem := range repoId {
		repoIdRule = append(repoIdRule, repoIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Chain.contract.FilterLogs(opts, "RoleChanged", repoIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &ChainRoleChangedIterator{contract: _Chain.contract, event: "RoleChanged", logs: logs, sub: sub}, nil
}

// WatchRoleChanged is a free log subscription operation binding the contract event 0xc98b3c1f2e93a069aa9220f35ddceb1fc071af0645ca6b79b9436a5b139447d9.
//
// Solidity: event RoleChanged(uint256 indexed repoId, address indexed user, uint8 role)
func (_Chain *ChainFilterer) WatchRoleChanged(opts *bind.WatchOpts, sink chan<- *ChainRoleChanged, repoId []*big.Int, user []common.Address) (event.Subscription, error) {

	var repoIdRule []interface{}
	for _, repoIdItem := range repoId {
		repoIdRule = append(repoIdRule, repoIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Chain.contract.WatchLogs(opts, "RoleChanged", repoIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChainRoleChanged)
				if err := _Chain.contract.UnpackLog(event, "RoleChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleChanged is a log parse operation binding the contract event 0xc98b3c1f2e93a069aa9220f35ddceb1fc071af0645ca6b79b9436a5b139447d9.
//
// Solidity: event RoleChanged(uint256 indexed repoId, address indexed user, uint8 role)
func (_Chain *ChainFilterer) ParseRoleChanged(log types.Log) (*ChainRoleChanged, error) {
	event := new(ChainRoleChanged)
	if err := _Chain.contract.UnpackLog(event, "RoleChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ChainTagCreatedIterator is returned from FilterTagCreated and is used to iterate over the raw logs and unpacked data for TagCreated events raised by the Chain contract.
type ChainTagCreatedIterator struct {
	Event *ChainTagCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ChainTagCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ChainTagCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ChainTagCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ChainTagCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ChainTagCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ChainTagCreated represents a TagCreated event raised by the Chain contract.
type ChainTagCreated struct {
	RepoId  *big.Int
	Tag     [32]byte
	Target  []byte
	Creator common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTagCreated is a free log retrieval operation binding the contract event 0xfaab59977056b9bafe405414d233180de7c9cc3a3ba57beb2963ba80cde9c94c.
//
// Solidity: event TagCreated(uint256 indexed repoId, bytes32 indexed tag, bytes target, address indexed creator)
func (_Chain *ChainFilterer) FilterTagCreated(opts *bind.FilterOpts, repoId []*big.Int, tag [][32]byte, creator []common.Address) (*ChainTagCreatedIterator, error) {

	var repoIdRule []interface{}
	for _, repoIdItem := range repoId {
		repoIdRule = append(repoIdRule, repoIdItem)
	}
	var tagRule []interface{}
	for _, tagItem := range tag {
		tagRule = append(tagRule, tagItem)
	}

	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _Chain.contract.FilterLogs(opts, "TagCreated", repoIdRule, tagRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &ChainTagCreatedIterator{contract: _Chain.contract, event: "TagCreated", logs: logs, sub: sub}, nil
}

// WatchTagCreated is a free log subscription operation binding the contract event 0xfaab59977056b9bafe405414d233180de7c9cc3a3ba57beb2963ba80cde9c94c.
//
// Solidity: event TagCreated(uint256 indexed repoId, bytes32 indexed tag, bytes target, address indexed creator)
func (_Chain *ChainFilterer) WatchTagCreated(opts *bind.WatchOpts, sink chan<- *ChainTagCreated, repoId []*big.Int, tag [][32]byte, creator []common.Address) (event.Subscription, error) {

	var repoIdRule []interface{}
	for _, repoIdItem := range repoId {
		repoIdRule = append(repoIdRule, repoIdItem)
	}
	var tagRule []interface{}
	for _, tagItem := range tag {
		tagRule = append(tagRule, tagItem)
	}

	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _Chain.contract.WatchLogs(opts, "TagCreated", repoIdRule, tagRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ChainTagCreated)
				if err := _Chain.contract.UnpackLog(event, "TagCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTagCreated is a log parse operation binding the contract event 0xfaab59977056b9bafe405414d233180de7c9cc3a3ba57beb2963ba80cde9c94c.
//
// Solidity: event TagCreated(uint256 indexed repoId, bytes32 indexed tag, bytes target, address indexed creator)
func (_Chain *ChainFilterer) ParseTagCreated(log types.Log) (*ChainTagCreated, error) {
	event := new(ChainTagCreated)
	if err := _Chain.contract.UnpackLog(event, "TagCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
