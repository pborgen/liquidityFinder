package dexV3

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	abi_pancakeswapv3_factory "github.com/pborgen/liquidityFinder/abi/pancakeswapV3"
	blockchainclient "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainClient"
)

type V3PairCreated struct {
	PairAddress common.Address
	Token0 common.Address
	Token1 common.Address
	Fee big.Int
	TickSpacing big.Int
}

func GetPairsCreated(
	factoryContractAddress common.Address, 
	startBlock uint64, 
	endBlock uint64) ([]V3PairCreated, error) {

	
	client := blockchainclient.GetHttpClient()

	factory, err := abi_pancakeswapv3_factory.NewAbiPancakeswapv3FactoryFilterer(factoryContractAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create filterer: %v", err)
	}

	filter := &bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	}
	
	iter, err := factory.FilterPoolCreated(filter, nil, nil, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create iterator: %v", err)
	}

	modelPairs := []V3PairCreated{}

	for iter.Next() {
		event := iter.Event

		token0 := event.Token0
		token1 := event.Token1
		fee := event.Fee
		tickSpacing := event.TickSpacing
		pairAddress := event.Pool

		modelPairs = append(modelPairs, V3PairCreated{
			PairAddress: pairAddress,
			Token0: token0,
			Token1: token1,
			Fee: *fee,
			TickSpacing: *tickSpacing,
		})
	}

	return modelPairs, nil
}