package uniswapV3Pool

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	abi_uniswapv3pool "github.com/pborgen/liquidityFinder/abi/uniswapV3/uniswapV3Pool"
	blockchainclient "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainClient"
)

type Slot0 struct {
	SqrtPriceX96               *big.Int
	Tick                       *big.Int
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
	Unlocked                   bool
}


func GetSlot0(poolAddress common.Address) (slot0 Slot0, err error) {
	client := blockchainclient.GetHttpClient()
	contract, err := abi_uniswapv3pool.NewAbiUniswapv3poolCaller(poolAddress, client)
	if err != nil {
		return Slot0{}, err
	}

	slot0, err = contract.Slot0(&bind.CallOpts{})
	if err != nil {
		return Slot0{}, err
	}

	return slot0, nil
}

