// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi_pancakeswapv3_factory

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

// AbiPancakeswapv3FactoryMetaData contains all meta data concerning the AbiPancakeswapv3Factory contract.
var AbiPancakeswapv3FactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_poolDeployer\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"indexed\":true,\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"}],\"name\":\"FeeAmountEnabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"whitelistRequested\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"FeeAmountExtraInfoUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnerChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"name\":\"PoolCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"lmPoolDeployer\",\"type\":\"address\"}],\"name\":\"SetLmPoolDeployer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"verified\",\"type\":\"bool\"}],\"name\":\"WhiteListAdded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"amount0Requested\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"amount1Requested\",\"type\":\"uint128\"}],\"name\":\"collectProtocol\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"amount0\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"amount1\",\"type\":\"uint128\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"}],\"name\":\"createPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"}],\"name\":\"enableFeeAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint24\",\"name\":\"\",\"type\":\"uint24\"}],\"name\":\"feeAmountTickSpacing\",\"outputs\":[{\"internalType\":\"int24\",\"name\":\"\",\"type\":\"int24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint24\",\"name\":\"\",\"type\":\"uint24\"}],\"name\":\"feeAmountTickSpacingExtraInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"whitelistRequested\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"\",\"type\":\"uint24\"}],\"name\":\"getPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lmPoolDeployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"poolDeployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"bool\",\"name\":\"whitelistRequested\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"setFeeAmountExtraInfo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"feeProtocol0\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"feeProtocol1\",\"type\":\"uint32\"}],\"name\":\"setFeeProtocol\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"lmPool\",\"type\":\"address\"}],\"name\":\"setLmPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_lmPoolDeployer\",\"type\":\"address\"}],\"name\":\"setLmPoolDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"verified\",\"type\":\"bool\"}],\"name\":\"setWhiteListAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AbiPancakeswapv3FactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use AbiPancakeswapv3FactoryMetaData.ABI instead.
var AbiPancakeswapv3FactoryABI = AbiPancakeswapv3FactoryMetaData.ABI

// AbiPancakeswapv3Factory is an auto generated Go binding around an Ethereum contract.
type AbiPancakeswapv3Factory struct {
	AbiPancakeswapv3FactoryCaller     // Read-only binding to the contract
	AbiPancakeswapv3FactoryTransactor // Write-only binding to the contract
	AbiPancakeswapv3FactoryFilterer   // Log filterer for contract events
}

// AbiPancakeswapv3FactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type AbiPancakeswapv3FactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbiPancakeswapv3FactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AbiPancakeswapv3FactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbiPancakeswapv3FactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AbiPancakeswapv3FactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AbiPancakeswapv3FactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AbiPancakeswapv3FactorySession struct {
	Contract     *AbiPancakeswapv3Factory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// AbiPancakeswapv3FactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AbiPancakeswapv3FactoryCallerSession struct {
	Contract *AbiPancakeswapv3FactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// AbiPancakeswapv3FactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AbiPancakeswapv3FactoryTransactorSession struct {
	Contract     *AbiPancakeswapv3FactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// AbiPancakeswapv3FactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type AbiPancakeswapv3FactoryRaw struct {
	Contract *AbiPancakeswapv3Factory // Generic contract binding to access the raw methods on
}

// AbiPancakeswapv3FactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AbiPancakeswapv3FactoryCallerRaw struct {
	Contract *AbiPancakeswapv3FactoryCaller // Generic read-only contract binding to access the raw methods on
}

// AbiPancakeswapv3FactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AbiPancakeswapv3FactoryTransactorRaw struct {
	Contract *AbiPancakeswapv3FactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAbiPancakeswapv3Factory creates a new instance of AbiPancakeswapv3Factory, bound to a specific deployed contract.
func NewAbiPancakeswapv3Factory(address common.Address, backend bind.ContractBackend) (*AbiPancakeswapv3Factory, error) {
	contract, err := bindAbiPancakeswapv3Factory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AbiPancakeswapv3Factory{AbiPancakeswapv3FactoryCaller: AbiPancakeswapv3FactoryCaller{contract: contract}, AbiPancakeswapv3FactoryTransactor: AbiPancakeswapv3FactoryTransactor{contract: contract}, AbiPancakeswapv3FactoryFilterer: AbiPancakeswapv3FactoryFilterer{contract: contract}}, nil
}

// NewAbiPancakeswapv3FactoryCaller creates a new read-only instance of AbiPancakeswapv3Factory, bound to a specific deployed contract.
func NewAbiPancakeswapv3FactoryCaller(address common.Address, caller bind.ContractCaller) (*AbiPancakeswapv3FactoryCaller, error) {
	contract, err := bindAbiPancakeswapv3Factory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AbiPancakeswapv3FactoryCaller{contract: contract}, nil
}

// NewAbiPancakeswapv3FactoryTransactor creates a new write-only instance of AbiPancakeswapv3Factory, bound to a specific deployed contract.
func NewAbiPancakeswapv3FactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*AbiPancakeswapv3FactoryTransactor, error) {
	contract, err := bindAbiPancakeswapv3Factory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AbiPancakeswapv3FactoryTransactor{contract: contract}, nil
}

// NewAbiPancakeswapv3FactoryFilterer creates a new log filterer instance of AbiPancakeswapv3Factory, bound to a specific deployed contract.
func NewAbiPancakeswapv3FactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*AbiPancakeswapv3FactoryFilterer, error) {
	contract, err := bindAbiPancakeswapv3Factory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AbiPancakeswapv3FactoryFilterer{contract: contract}, nil
}

// bindAbiPancakeswapv3Factory binds a generic wrapper to an already deployed contract.
func bindAbiPancakeswapv3Factory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AbiPancakeswapv3FactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AbiPancakeswapv3Factory.Contract.AbiPancakeswapv3FactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.AbiPancakeswapv3FactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.AbiPancakeswapv3FactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AbiPancakeswapv3Factory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.contract.Transact(opts, method, params...)
}

// FeeAmountTickSpacing is a free data retrieval call binding the contract method 0x22afcccb.
//
// Solidity: function feeAmountTickSpacing(uint24 ) view returns(int24)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryCaller) FeeAmountTickSpacing(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AbiPancakeswapv3Factory.contract.Call(opts, &out, "feeAmountTickSpacing", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeAmountTickSpacing is a free data retrieval call binding the contract method 0x22afcccb.
//
// Solidity: function feeAmountTickSpacing(uint24 ) view returns(int24)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) FeeAmountTickSpacing(arg0 *big.Int) (*big.Int, error) {
	return _AbiPancakeswapv3Factory.Contract.FeeAmountTickSpacing(&_AbiPancakeswapv3Factory.CallOpts, arg0)
}

// FeeAmountTickSpacing is a free data retrieval call binding the contract method 0x22afcccb.
//
// Solidity: function feeAmountTickSpacing(uint24 ) view returns(int24)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryCallerSession) FeeAmountTickSpacing(arg0 *big.Int) (*big.Int, error) {
	return _AbiPancakeswapv3Factory.Contract.FeeAmountTickSpacing(&_AbiPancakeswapv3Factory.CallOpts, arg0)
}

// FeeAmountTickSpacingExtraInfo is a free data retrieval call binding the contract method 0x88e8006d.
//
// Solidity: function feeAmountTickSpacingExtraInfo(uint24 ) view returns(bool whitelistRequested, bool enabled)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryCaller) FeeAmountTickSpacingExtraInfo(opts *bind.CallOpts, arg0 *big.Int) (struct {
	WhitelistRequested bool
	Enabled            bool
}, error) {
	var out []interface{}
	err := _AbiPancakeswapv3Factory.contract.Call(opts, &out, "feeAmountTickSpacingExtraInfo", arg0)

	outstruct := new(struct {
		WhitelistRequested bool
		Enabled            bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.WhitelistRequested = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Enabled = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// FeeAmountTickSpacingExtraInfo is a free data retrieval call binding the contract method 0x88e8006d.
//
// Solidity: function feeAmountTickSpacingExtraInfo(uint24 ) view returns(bool whitelistRequested, bool enabled)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) FeeAmountTickSpacingExtraInfo(arg0 *big.Int) (struct {
	WhitelistRequested bool
	Enabled            bool
}, error) {
	return _AbiPancakeswapv3Factory.Contract.FeeAmountTickSpacingExtraInfo(&_AbiPancakeswapv3Factory.CallOpts, arg0)
}

// FeeAmountTickSpacingExtraInfo is a free data retrieval call binding the contract method 0x88e8006d.
//
// Solidity: function feeAmountTickSpacingExtraInfo(uint24 ) view returns(bool whitelistRequested, bool enabled)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryCallerSession) FeeAmountTickSpacingExtraInfo(arg0 *big.Int) (struct {
	WhitelistRequested bool
	Enabled            bool
}, error) {
	return _AbiPancakeswapv3Factory.Contract.FeeAmountTickSpacingExtraInfo(&_AbiPancakeswapv3Factory.CallOpts, arg0)
}

// GetPool is a free data retrieval call binding the contract method 0x1698ee82.
//
// Solidity: function getPool(address , address , uint24 ) view returns(address)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryCaller) GetPool(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AbiPancakeswapv3Factory.contract.Call(opts, &out, "getPool", arg0, arg1, arg2)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPool is a free data retrieval call binding the contract method 0x1698ee82.
//
// Solidity: function getPool(address , address , uint24 ) view returns(address)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) GetPool(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (common.Address, error) {
	return _AbiPancakeswapv3Factory.Contract.GetPool(&_AbiPancakeswapv3Factory.CallOpts, arg0, arg1, arg2)
}

// GetPool is a free data retrieval call binding the contract method 0x1698ee82.
//
// Solidity: function getPool(address , address , uint24 ) view returns(address)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryCallerSession) GetPool(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (common.Address, error) {
	return _AbiPancakeswapv3Factory.Contract.GetPool(&_AbiPancakeswapv3Factory.CallOpts, arg0, arg1, arg2)
}

// LmPoolDeployer is a free data retrieval call binding the contract method 0x5e492ac8.
//
// Solidity: function lmPoolDeployer() view returns(address)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryCaller) LmPoolDeployer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AbiPancakeswapv3Factory.contract.Call(opts, &out, "lmPoolDeployer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LmPoolDeployer is a free data retrieval call binding the contract method 0x5e492ac8.
//
// Solidity: function lmPoolDeployer() view returns(address)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) LmPoolDeployer() (common.Address, error) {
	return _AbiPancakeswapv3Factory.Contract.LmPoolDeployer(&_AbiPancakeswapv3Factory.CallOpts)
}

// LmPoolDeployer is a free data retrieval call binding the contract method 0x5e492ac8.
//
// Solidity: function lmPoolDeployer() view returns(address)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryCallerSession) LmPoolDeployer() (common.Address, error) {
	return _AbiPancakeswapv3Factory.Contract.LmPoolDeployer(&_AbiPancakeswapv3Factory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AbiPancakeswapv3Factory.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) Owner() (common.Address, error) {
	return _AbiPancakeswapv3Factory.Contract.Owner(&_AbiPancakeswapv3Factory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryCallerSession) Owner() (common.Address, error) {
	return _AbiPancakeswapv3Factory.Contract.Owner(&_AbiPancakeswapv3Factory.CallOpts)
}

// PoolDeployer is a free data retrieval call binding the contract method 0x3119049a.
//
// Solidity: function poolDeployer() view returns(address)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryCaller) PoolDeployer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AbiPancakeswapv3Factory.contract.Call(opts, &out, "poolDeployer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolDeployer is a free data retrieval call binding the contract method 0x3119049a.
//
// Solidity: function poolDeployer() view returns(address)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) PoolDeployer() (common.Address, error) {
	return _AbiPancakeswapv3Factory.Contract.PoolDeployer(&_AbiPancakeswapv3Factory.CallOpts)
}

// PoolDeployer is a free data retrieval call binding the contract method 0x3119049a.
//
// Solidity: function poolDeployer() view returns(address)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryCallerSession) PoolDeployer() (common.Address, error) {
	return _AbiPancakeswapv3Factory.Contract.PoolDeployer(&_AbiPancakeswapv3Factory.CallOpts)
}

// CollectProtocol is a paid mutator transaction binding the contract method 0x43db87da.
//
// Solidity: function collectProtocol(address pool, address recipient, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactor) CollectProtocol(opts *bind.TransactOpts, pool common.Address, recipient common.Address, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.contract.Transact(opts, "collectProtocol", pool, recipient, amount0Requested, amount1Requested)
}

// CollectProtocol is a paid mutator transaction binding the contract method 0x43db87da.
//
// Solidity: function collectProtocol(address pool, address recipient, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) CollectProtocol(pool common.Address, recipient common.Address, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.CollectProtocol(&_AbiPancakeswapv3Factory.TransactOpts, pool, recipient, amount0Requested, amount1Requested)
}

// CollectProtocol is a paid mutator transaction binding the contract method 0x43db87da.
//
// Solidity: function collectProtocol(address pool, address recipient, uint128 amount0Requested, uint128 amount1Requested) returns(uint128 amount0, uint128 amount1)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactorSession) CollectProtocol(pool common.Address, recipient common.Address, amount0Requested *big.Int, amount1Requested *big.Int) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.CollectProtocol(&_AbiPancakeswapv3Factory.TransactOpts, pool, recipient, amount0Requested, amount1Requested)
}

// CreatePool is a paid mutator transaction binding the contract method 0xa1671295.
//
// Solidity: function createPool(address tokenA, address tokenB, uint24 fee) returns(address pool)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactor) CreatePool(opts *bind.TransactOpts, tokenA common.Address, tokenB common.Address, fee *big.Int) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.contract.Transact(opts, "createPool", tokenA, tokenB, fee)
}

// CreatePool is a paid mutator transaction binding the contract method 0xa1671295.
//
// Solidity: function createPool(address tokenA, address tokenB, uint24 fee) returns(address pool)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) CreatePool(tokenA common.Address, tokenB common.Address, fee *big.Int) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.CreatePool(&_AbiPancakeswapv3Factory.TransactOpts, tokenA, tokenB, fee)
}

// CreatePool is a paid mutator transaction binding the contract method 0xa1671295.
//
// Solidity: function createPool(address tokenA, address tokenB, uint24 fee) returns(address pool)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactorSession) CreatePool(tokenA common.Address, tokenB common.Address, fee *big.Int) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.CreatePool(&_AbiPancakeswapv3Factory.TransactOpts, tokenA, tokenB, fee)
}

// EnableFeeAmount is a paid mutator transaction binding the contract method 0x8a7c195f.
//
// Solidity: function enableFeeAmount(uint24 fee, int24 tickSpacing) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactor) EnableFeeAmount(opts *bind.TransactOpts, fee *big.Int, tickSpacing *big.Int) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.contract.Transact(opts, "enableFeeAmount", fee, tickSpacing)
}

// EnableFeeAmount is a paid mutator transaction binding the contract method 0x8a7c195f.
//
// Solidity: function enableFeeAmount(uint24 fee, int24 tickSpacing) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) EnableFeeAmount(fee *big.Int, tickSpacing *big.Int) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.EnableFeeAmount(&_AbiPancakeswapv3Factory.TransactOpts, fee, tickSpacing)
}

// EnableFeeAmount is a paid mutator transaction binding the contract method 0x8a7c195f.
//
// Solidity: function enableFeeAmount(uint24 fee, int24 tickSpacing) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactorSession) EnableFeeAmount(fee *big.Int, tickSpacing *big.Int) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.EnableFeeAmount(&_AbiPancakeswapv3Factory.TransactOpts, fee, tickSpacing)
}

// SetFeeAmountExtraInfo is a paid mutator transaction binding the contract method 0x8ff38e80.
//
// Solidity: function setFeeAmountExtraInfo(uint24 fee, bool whitelistRequested, bool enabled) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactor) SetFeeAmountExtraInfo(opts *bind.TransactOpts, fee *big.Int, whitelistRequested bool, enabled bool) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.contract.Transact(opts, "setFeeAmountExtraInfo", fee, whitelistRequested, enabled)
}

// SetFeeAmountExtraInfo is a paid mutator transaction binding the contract method 0x8ff38e80.
//
// Solidity: function setFeeAmountExtraInfo(uint24 fee, bool whitelistRequested, bool enabled) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) SetFeeAmountExtraInfo(fee *big.Int, whitelistRequested bool, enabled bool) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.SetFeeAmountExtraInfo(&_AbiPancakeswapv3Factory.TransactOpts, fee, whitelistRequested, enabled)
}

// SetFeeAmountExtraInfo is a paid mutator transaction binding the contract method 0x8ff38e80.
//
// Solidity: function setFeeAmountExtraInfo(uint24 fee, bool whitelistRequested, bool enabled) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactorSession) SetFeeAmountExtraInfo(fee *big.Int, whitelistRequested bool, enabled bool) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.SetFeeAmountExtraInfo(&_AbiPancakeswapv3Factory.TransactOpts, fee, whitelistRequested, enabled)
}

// SetFeeProtocol is a paid mutator transaction binding the contract method 0x7e8435e6.
//
// Solidity: function setFeeProtocol(address pool, uint32 feeProtocol0, uint32 feeProtocol1) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactor) SetFeeProtocol(opts *bind.TransactOpts, pool common.Address, feeProtocol0 uint32, feeProtocol1 uint32) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.contract.Transact(opts, "setFeeProtocol", pool, feeProtocol0, feeProtocol1)
}

// SetFeeProtocol is a paid mutator transaction binding the contract method 0x7e8435e6.
//
// Solidity: function setFeeProtocol(address pool, uint32 feeProtocol0, uint32 feeProtocol1) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) SetFeeProtocol(pool common.Address, feeProtocol0 uint32, feeProtocol1 uint32) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.SetFeeProtocol(&_AbiPancakeswapv3Factory.TransactOpts, pool, feeProtocol0, feeProtocol1)
}

// SetFeeProtocol is a paid mutator transaction binding the contract method 0x7e8435e6.
//
// Solidity: function setFeeProtocol(address pool, uint32 feeProtocol0, uint32 feeProtocol1) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactorSession) SetFeeProtocol(pool common.Address, feeProtocol0 uint32, feeProtocol1 uint32) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.SetFeeProtocol(&_AbiPancakeswapv3Factory.TransactOpts, pool, feeProtocol0, feeProtocol1)
}

// SetLmPool is a paid mutator transaction binding the contract method 0x11ff5e8d.
//
// Solidity: function setLmPool(address pool, address lmPool) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactor) SetLmPool(opts *bind.TransactOpts, pool common.Address, lmPool common.Address) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.contract.Transact(opts, "setLmPool", pool, lmPool)
}

// SetLmPool is a paid mutator transaction binding the contract method 0x11ff5e8d.
//
// Solidity: function setLmPool(address pool, address lmPool) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) SetLmPool(pool common.Address, lmPool common.Address) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.SetLmPool(&_AbiPancakeswapv3Factory.TransactOpts, pool, lmPool)
}

// SetLmPool is a paid mutator transaction binding the contract method 0x11ff5e8d.
//
// Solidity: function setLmPool(address pool, address lmPool) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactorSession) SetLmPool(pool common.Address, lmPool common.Address) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.SetLmPool(&_AbiPancakeswapv3Factory.TransactOpts, pool, lmPool)
}

// SetLmPoolDeployer is a paid mutator transaction binding the contract method 0x80d6a792.
//
// Solidity: function setLmPoolDeployer(address _lmPoolDeployer) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactor) SetLmPoolDeployer(opts *bind.TransactOpts, _lmPoolDeployer common.Address) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.contract.Transact(opts, "setLmPoolDeployer", _lmPoolDeployer)
}

// SetLmPoolDeployer is a paid mutator transaction binding the contract method 0x80d6a792.
//
// Solidity: function setLmPoolDeployer(address _lmPoolDeployer) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) SetLmPoolDeployer(_lmPoolDeployer common.Address) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.SetLmPoolDeployer(&_AbiPancakeswapv3Factory.TransactOpts, _lmPoolDeployer)
}

// SetLmPoolDeployer is a paid mutator transaction binding the contract method 0x80d6a792.
//
// Solidity: function setLmPoolDeployer(address _lmPoolDeployer) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactorSession) SetLmPoolDeployer(_lmPoolDeployer common.Address) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.SetLmPoolDeployer(&_AbiPancakeswapv3Factory.TransactOpts, _lmPoolDeployer)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(address _owner) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactor) SetOwner(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.contract.Transact(opts, "setOwner", _owner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(address _owner) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) SetOwner(_owner common.Address) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.SetOwner(&_AbiPancakeswapv3Factory.TransactOpts, _owner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(address _owner) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactorSession) SetOwner(_owner common.Address) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.SetOwner(&_AbiPancakeswapv3Factory.TransactOpts, _owner)
}

// SetWhiteListAddress is a paid mutator transaction binding the contract method 0xe4a86a99.
//
// Solidity: function setWhiteListAddress(address user, bool verified) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactor) SetWhiteListAddress(opts *bind.TransactOpts, user common.Address, verified bool) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.contract.Transact(opts, "setWhiteListAddress", user, verified)
}

// SetWhiteListAddress is a paid mutator transaction binding the contract method 0xe4a86a99.
//
// Solidity: function setWhiteListAddress(address user, bool verified) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactorySession) SetWhiteListAddress(user common.Address, verified bool) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.SetWhiteListAddress(&_AbiPancakeswapv3Factory.TransactOpts, user, verified)
}

// SetWhiteListAddress is a paid mutator transaction binding the contract method 0xe4a86a99.
//
// Solidity: function setWhiteListAddress(address user, bool verified) returns()
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryTransactorSession) SetWhiteListAddress(user common.Address, verified bool) (*types.Transaction, error) {
	return _AbiPancakeswapv3Factory.Contract.SetWhiteListAddress(&_AbiPancakeswapv3Factory.TransactOpts, user, verified)
}

// AbiPancakeswapv3FactoryFeeAmountEnabledIterator is returned from FilterFeeAmountEnabled and is used to iterate over the raw logs and unpacked data for FeeAmountEnabled events raised by the AbiPancakeswapv3Factory contract.
type AbiPancakeswapv3FactoryFeeAmountEnabledIterator struct {
	Event *AbiPancakeswapv3FactoryFeeAmountEnabled // Event containing the contract specifics and raw log

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
func (it *AbiPancakeswapv3FactoryFeeAmountEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiPancakeswapv3FactoryFeeAmountEnabled)
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
		it.Event = new(AbiPancakeswapv3FactoryFeeAmountEnabled)
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
func (it *AbiPancakeswapv3FactoryFeeAmountEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiPancakeswapv3FactoryFeeAmountEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiPancakeswapv3FactoryFeeAmountEnabled represents a FeeAmountEnabled event raised by the AbiPancakeswapv3Factory contract.
type AbiPancakeswapv3FactoryFeeAmountEnabled struct {
	Fee         *big.Int
	TickSpacing *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterFeeAmountEnabled is a free log retrieval operation binding the contract event 0xc66a3fdf07232cdd185febcc6579d408c241b47ae2f9907d84be655141eeaecc.
//
// Solidity: event FeeAmountEnabled(uint24 indexed fee, int24 indexed tickSpacing)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) FilterFeeAmountEnabled(opts *bind.FilterOpts, fee []*big.Int, tickSpacing []*big.Int) (*AbiPancakeswapv3FactoryFeeAmountEnabledIterator, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}
	var tickSpacingRule []interface{}
	for _, tickSpacingItem := range tickSpacing {
		tickSpacingRule = append(tickSpacingRule, tickSpacingItem)
	}

	logs, sub, err := _AbiPancakeswapv3Factory.contract.FilterLogs(opts, "FeeAmountEnabled", feeRule, tickSpacingRule)
	if err != nil {
		return nil, err
	}
	return &AbiPancakeswapv3FactoryFeeAmountEnabledIterator{contract: _AbiPancakeswapv3Factory.contract, event: "FeeAmountEnabled", logs: logs, sub: sub}, nil
}

// WatchFeeAmountEnabled is a free log subscription operation binding the contract event 0xc66a3fdf07232cdd185febcc6579d408c241b47ae2f9907d84be655141eeaecc.
//
// Solidity: event FeeAmountEnabled(uint24 indexed fee, int24 indexed tickSpacing)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) WatchFeeAmountEnabled(opts *bind.WatchOpts, sink chan<- *AbiPancakeswapv3FactoryFeeAmountEnabled, fee []*big.Int, tickSpacing []*big.Int) (event.Subscription, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}
	var tickSpacingRule []interface{}
	for _, tickSpacingItem := range tickSpacing {
		tickSpacingRule = append(tickSpacingRule, tickSpacingItem)
	}

	logs, sub, err := _AbiPancakeswapv3Factory.contract.WatchLogs(opts, "FeeAmountEnabled", feeRule, tickSpacingRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiPancakeswapv3FactoryFeeAmountEnabled)
				if err := _AbiPancakeswapv3Factory.contract.UnpackLog(event, "FeeAmountEnabled", log); err != nil {
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

// ParseFeeAmountEnabled is a log parse operation binding the contract event 0xc66a3fdf07232cdd185febcc6579d408c241b47ae2f9907d84be655141eeaecc.
//
// Solidity: event FeeAmountEnabled(uint24 indexed fee, int24 indexed tickSpacing)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) ParseFeeAmountEnabled(log types.Log) (*AbiPancakeswapv3FactoryFeeAmountEnabled, error) {
	event := new(AbiPancakeswapv3FactoryFeeAmountEnabled)
	if err := _AbiPancakeswapv3Factory.contract.UnpackLog(event, "FeeAmountEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdatedIterator is returned from FilterFeeAmountExtraInfoUpdated and is used to iterate over the raw logs and unpacked data for FeeAmountExtraInfoUpdated events raised by the AbiPancakeswapv3Factory contract.
type AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdatedIterator struct {
	Event *AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdated // Event containing the contract specifics and raw log

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
func (it *AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdated)
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
		it.Event = new(AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdated)
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
func (it *AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdated represents a FeeAmountExtraInfoUpdated event raised by the AbiPancakeswapv3Factory contract.
type AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdated struct {
	Fee                *big.Int
	WhitelistRequested bool
	Enabled            bool
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterFeeAmountExtraInfoUpdated is a free log retrieval operation binding the contract event 0xed85b616dbfbc54d0f1180a7bd0f6e3bb645b269b234e7a9edcc269ef1443d88.
//
// Solidity: event FeeAmountExtraInfoUpdated(uint24 indexed fee, bool whitelistRequested, bool enabled)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) FilterFeeAmountExtraInfoUpdated(opts *bind.FilterOpts, fee []*big.Int) (*AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdatedIterator, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _AbiPancakeswapv3Factory.contract.FilterLogs(opts, "FeeAmountExtraInfoUpdated", feeRule)
	if err != nil {
		return nil, err
	}
	return &AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdatedIterator{contract: _AbiPancakeswapv3Factory.contract, event: "FeeAmountExtraInfoUpdated", logs: logs, sub: sub}, nil
}

// WatchFeeAmountExtraInfoUpdated is a free log subscription operation binding the contract event 0xed85b616dbfbc54d0f1180a7bd0f6e3bb645b269b234e7a9edcc269ef1443d88.
//
// Solidity: event FeeAmountExtraInfoUpdated(uint24 indexed fee, bool whitelistRequested, bool enabled)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) WatchFeeAmountExtraInfoUpdated(opts *bind.WatchOpts, sink chan<- *AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdated, fee []*big.Int) (event.Subscription, error) {

	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _AbiPancakeswapv3Factory.contract.WatchLogs(opts, "FeeAmountExtraInfoUpdated", feeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdated)
				if err := _AbiPancakeswapv3Factory.contract.UnpackLog(event, "FeeAmountExtraInfoUpdated", log); err != nil {
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

// ParseFeeAmountExtraInfoUpdated is a log parse operation binding the contract event 0xed85b616dbfbc54d0f1180a7bd0f6e3bb645b269b234e7a9edcc269ef1443d88.
//
// Solidity: event FeeAmountExtraInfoUpdated(uint24 indexed fee, bool whitelistRequested, bool enabled)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) ParseFeeAmountExtraInfoUpdated(log types.Log) (*AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdated, error) {
	event := new(AbiPancakeswapv3FactoryFeeAmountExtraInfoUpdated)
	if err := _AbiPancakeswapv3Factory.contract.UnpackLog(event, "FeeAmountExtraInfoUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbiPancakeswapv3FactoryOwnerChangedIterator is returned from FilterOwnerChanged and is used to iterate over the raw logs and unpacked data for OwnerChanged events raised by the AbiPancakeswapv3Factory contract.
type AbiPancakeswapv3FactoryOwnerChangedIterator struct {
	Event *AbiPancakeswapv3FactoryOwnerChanged // Event containing the contract specifics and raw log

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
func (it *AbiPancakeswapv3FactoryOwnerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiPancakeswapv3FactoryOwnerChanged)
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
		it.Event = new(AbiPancakeswapv3FactoryOwnerChanged)
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
func (it *AbiPancakeswapv3FactoryOwnerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiPancakeswapv3FactoryOwnerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiPancakeswapv3FactoryOwnerChanged represents a OwnerChanged event raised by the AbiPancakeswapv3Factory contract.
type AbiPancakeswapv3FactoryOwnerChanged struct {
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnerChanged is a free log retrieval operation binding the contract event 0xb532073b38c83145e3e5135377a08bf9aab55bc0fd7c1179cd4fb995d2a5159c.
//
// Solidity: event OwnerChanged(address indexed oldOwner, address indexed newOwner)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) FilterOwnerChanged(opts *bind.FilterOpts, oldOwner []common.Address, newOwner []common.Address) (*AbiPancakeswapv3FactoryOwnerChangedIterator, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AbiPancakeswapv3Factory.contract.FilterLogs(opts, "OwnerChanged", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AbiPancakeswapv3FactoryOwnerChangedIterator{contract: _AbiPancakeswapv3Factory.contract, event: "OwnerChanged", logs: logs, sub: sub}, nil
}

// WatchOwnerChanged is a free log subscription operation binding the contract event 0xb532073b38c83145e3e5135377a08bf9aab55bc0fd7c1179cd4fb995d2a5159c.
//
// Solidity: event OwnerChanged(address indexed oldOwner, address indexed newOwner)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) WatchOwnerChanged(opts *bind.WatchOpts, sink chan<- *AbiPancakeswapv3FactoryOwnerChanged, oldOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AbiPancakeswapv3Factory.contract.WatchLogs(opts, "OwnerChanged", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiPancakeswapv3FactoryOwnerChanged)
				if err := _AbiPancakeswapv3Factory.contract.UnpackLog(event, "OwnerChanged", log); err != nil {
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

// ParseOwnerChanged is a log parse operation binding the contract event 0xb532073b38c83145e3e5135377a08bf9aab55bc0fd7c1179cd4fb995d2a5159c.
//
// Solidity: event OwnerChanged(address indexed oldOwner, address indexed newOwner)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) ParseOwnerChanged(log types.Log) (*AbiPancakeswapv3FactoryOwnerChanged, error) {
	event := new(AbiPancakeswapv3FactoryOwnerChanged)
	if err := _AbiPancakeswapv3Factory.contract.UnpackLog(event, "OwnerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbiPancakeswapv3FactoryPoolCreatedIterator is returned from FilterPoolCreated and is used to iterate over the raw logs and unpacked data for PoolCreated events raised by the AbiPancakeswapv3Factory contract.
type AbiPancakeswapv3FactoryPoolCreatedIterator struct {
	Event *AbiPancakeswapv3FactoryPoolCreated // Event containing the contract specifics and raw log

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
func (it *AbiPancakeswapv3FactoryPoolCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiPancakeswapv3FactoryPoolCreated)
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
		it.Event = new(AbiPancakeswapv3FactoryPoolCreated)
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
func (it *AbiPancakeswapv3FactoryPoolCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiPancakeswapv3FactoryPoolCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiPancakeswapv3FactoryPoolCreated represents a PoolCreated event raised by the AbiPancakeswapv3Factory contract.
type AbiPancakeswapv3FactoryPoolCreated struct {
	Token0      common.Address
	Token1      common.Address
	Fee         *big.Int
	TickSpacing *big.Int
	Pool        common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPoolCreated is a free log retrieval operation binding the contract event 0x783cca1c0412dd0d695e784568c96da2e9c22ff989357a2e8b1d9b2b4e6b7118.
//
// Solidity: event PoolCreated(address indexed token0, address indexed token1, uint24 indexed fee, int24 tickSpacing, address pool)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) FilterPoolCreated(opts *bind.FilterOpts, token0 []common.Address, token1 []common.Address, fee []*big.Int) (*AbiPancakeswapv3FactoryPoolCreatedIterator, error) {

	var token0Rule []interface{}
	for _, token0Item := range token0 {
		token0Rule = append(token0Rule, token0Item)
	}
	var token1Rule []interface{}
	for _, token1Item := range token1 {
		token1Rule = append(token1Rule, token1Item)
	}
	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _AbiPancakeswapv3Factory.contract.FilterLogs(opts, "PoolCreated", token0Rule, token1Rule, feeRule)
	if err != nil {
		return nil, err
	}
	return &AbiPancakeswapv3FactoryPoolCreatedIterator{contract: _AbiPancakeswapv3Factory.contract, event: "PoolCreated", logs: logs, sub: sub}, nil
}

// WatchPoolCreated is a free log subscription operation binding the contract event 0x783cca1c0412dd0d695e784568c96da2e9c22ff989357a2e8b1d9b2b4e6b7118.
//
// Solidity: event PoolCreated(address indexed token0, address indexed token1, uint24 indexed fee, int24 tickSpacing, address pool)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) WatchPoolCreated(opts *bind.WatchOpts, sink chan<- *AbiPancakeswapv3FactoryPoolCreated, token0 []common.Address, token1 []common.Address, fee []*big.Int) (event.Subscription, error) {

	var token0Rule []interface{}
	for _, token0Item := range token0 {
		token0Rule = append(token0Rule, token0Item)
	}
	var token1Rule []interface{}
	for _, token1Item := range token1 {
		token1Rule = append(token1Rule, token1Item)
	}
	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _AbiPancakeswapv3Factory.contract.WatchLogs(opts, "PoolCreated", token0Rule, token1Rule, feeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiPancakeswapv3FactoryPoolCreated)
				if err := _AbiPancakeswapv3Factory.contract.UnpackLog(event, "PoolCreated", log); err != nil {
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

// ParsePoolCreated is a log parse operation binding the contract event 0x783cca1c0412dd0d695e784568c96da2e9c22ff989357a2e8b1d9b2b4e6b7118.
//
// Solidity: event PoolCreated(address indexed token0, address indexed token1, uint24 indexed fee, int24 tickSpacing, address pool)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) ParsePoolCreated(log types.Log) (*AbiPancakeswapv3FactoryPoolCreated, error) {
	event := new(AbiPancakeswapv3FactoryPoolCreated)
	if err := _AbiPancakeswapv3Factory.contract.UnpackLog(event, "PoolCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbiPancakeswapv3FactorySetLmPoolDeployerIterator is returned from FilterSetLmPoolDeployer and is used to iterate over the raw logs and unpacked data for SetLmPoolDeployer events raised by the AbiPancakeswapv3Factory contract.
type AbiPancakeswapv3FactorySetLmPoolDeployerIterator struct {
	Event *AbiPancakeswapv3FactorySetLmPoolDeployer // Event containing the contract specifics and raw log

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
func (it *AbiPancakeswapv3FactorySetLmPoolDeployerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiPancakeswapv3FactorySetLmPoolDeployer)
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
		it.Event = new(AbiPancakeswapv3FactorySetLmPoolDeployer)
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
func (it *AbiPancakeswapv3FactorySetLmPoolDeployerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiPancakeswapv3FactorySetLmPoolDeployerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiPancakeswapv3FactorySetLmPoolDeployer represents a SetLmPoolDeployer event raised by the AbiPancakeswapv3Factory contract.
type AbiPancakeswapv3FactorySetLmPoolDeployer struct {
	LmPoolDeployer common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSetLmPoolDeployer is a free log retrieval operation binding the contract event 0x4c912280cda47bed324de14f601d3f125a98254671772f3f1f491e50fa0ca407.
//
// Solidity: event SetLmPoolDeployer(address indexed lmPoolDeployer)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) FilterSetLmPoolDeployer(opts *bind.FilterOpts, lmPoolDeployer []common.Address) (*AbiPancakeswapv3FactorySetLmPoolDeployerIterator, error) {

	var lmPoolDeployerRule []interface{}
	for _, lmPoolDeployerItem := range lmPoolDeployer {
		lmPoolDeployerRule = append(lmPoolDeployerRule, lmPoolDeployerItem)
	}

	logs, sub, err := _AbiPancakeswapv3Factory.contract.FilterLogs(opts, "SetLmPoolDeployer", lmPoolDeployerRule)
	if err != nil {
		return nil, err
	}
	return &AbiPancakeswapv3FactorySetLmPoolDeployerIterator{contract: _AbiPancakeswapv3Factory.contract, event: "SetLmPoolDeployer", logs: logs, sub: sub}, nil
}

// WatchSetLmPoolDeployer is a free log subscription operation binding the contract event 0x4c912280cda47bed324de14f601d3f125a98254671772f3f1f491e50fa0ca407.
//
// Solidity: event SetLmPoolDeployer(address indexed lmPoolDeployer)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) WatchSetLmPoolDeployer(opts *bind.WatchOpts, sink chan<- *AbiPancakeswapv3FactorySetLmPoolDeployer, lmPoolDeployer []common.Address) (event.Subscription, error) {

	var lmPoolDeployerRule []interface{}
	for _, lmPoolDeployerItem := range lmPoolDeployer {
		lmPoolDeployerRule = append(lmPoolDeployerRule, lmPoolDeployerItem)
	}

	logs, sub, err := _AbiPancakeswapv3Factory.contract.WatchLogs(opts, "SetLmPoolDeployer", lmPoolDeployerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiPancakeswapv3FactorySetLmPoolDeployer)
				if err := _AbiPancakeswapv3Factory.contract.UnpackLog(event, "SetLmPoolDeployer", log); err != nil {
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

// ParseSetLmPoolDeployer is a log parse operation binding the contract event 0x4c912280cda47bed324de14f601d3f125a98254671772f3f1f491e50fa0ca407.
//
// Solidity: event SetLmPoolDeployer(address indexed lmPoolDeployer)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) ParseSetLmPoolDeployer(log types.Log) (*AbiPancakeswapv3FactorySetLmPoolDeployer, error) {
	event := new(AbiPancakeswapv3FactorySetLmPoolDeployer)
	if err := _AbiPancakeswapv3Factory.contract.UnpackLog(event, "SetLmPoolDeployer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AbiPancakeswapv3FactoryWhiteListAddedIterator is returned from FilterWhiteListAdded and is used to iterate over the raw logs and unpacked data for WhiteListAdded events raised by the AbiPancakeswapv3Factory contract.
type AbiPancakeswapv3FactoryWhiteListAddedIterator struct {
	Event *AbiPancakeswapv3FactoryWhiteListAdded // Event containing the contract specifics and raw log

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
func (it *AbiPancakeswapv3FactoryWhiteListAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AbiPancakeswapv3FactoryWhiteListAdded)
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
		it.Event = new(AbiPancakeswapv3FactoryWhiteListAdded)
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
func (it *AbiPancakeswapv3FactoryWhiteListAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AbiPancakeswapv3FactoryWhiteListAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AbiPancakeswapv3FactoryWhiteListAdded represents a WhiteListAdded event raised by the AbiPancakeswapv3Factory contract.
type AbiPancakeswapv3FactoryWhiteListAdded struct {
	User     common.Address
	Verified bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWhiteListAdded is a free log retrieval operation binding the contract event 0xaec42ac7f1bb8651906ae6522f50a19429e124e8ea678ef59fd27750759288a2.
//
// Solidity: event WhiteListAdded(address indexed user, bool verified)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) FilterWhiteListAdded(opts *bind.FilterOpts, user []common.Address) (*AbiPancakeswapv3FactoryWhiteListAddedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _AbiPancakeswapv3Factory.contract.FilterLogs(opts, "WhiteListAdded", userRule)
	if err != nil {
		return nil, err
	}
	return &AbiPancakeswapv3FactoryWhiteListAddedIterator{contract: _AbiPancakeswapv3Factory.contract, event: "WhiteListAdded", logs: logs, sub: sub}, nil
}

// WatchWhiteListAdded is a free log subscription operation binding the contract event 0xaec42ac7f1bb8651906ae6522f50a19429e124e8ea678ef59fd27750759288a2.
//
// Solidity: event WhiteListAdded(address indexed user, bool verified)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) WatchWhiteListAdded(opts *bind.WatchOpts, sink chan<- *AbiPancakeswapv3FactoryWhiteListAdded, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _AbiPancakeswapv3Factory.contract.WatchLogs(opts, "WhiteListAdded", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AbiPancakeswapv3FactoryWhiteListAdded)
				if err := _AbiPancakeswapv3Factory.contract.UnpackLog(event, "WhiteListAdded", log); err != nil {
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

// ParseWhiteListAdded is a log parse operation binding the contract event 0xaec42ac7f1bb8651906ae6522f50a19429e124e8ea678ef59fd27750759288a2.
//
// Solidity: event WhiteListAdded(address indexed user, bool verified)
func (_AbiPancakeswapv3Factory *AbiPancakeswapv3FactoryFilterer) ParseWhiteListAdded(log types.Log) (*AbiPancakeswapv3FactoryWhiteListAdded, error) {
	event := new(AbiPancakeswapv3FactoryWhiteListAdded)
	if err := _AbiPancakeswapv3Factory.contract.UnpackLog(event, "WhiteListAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
