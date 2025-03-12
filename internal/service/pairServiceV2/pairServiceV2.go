package pairServiceV2

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/blockchain/dex/v2/dexUniswapV2Pair"
)


func GetReservesForPair(
	pairAddress common.Address, 
	tokenA common.Address, 
	tokenB common.Address, 
	retries int) (reserveTokenA *big.Int, reserveTokenB *big.Int, err error) {

	reserveTokenA, reserveTokenB, err = dexUniswapV2Pair.GetReservesForPair(
		pairAddress, 
		tokenA, 
		tokenB, 
		retries,
	)

	return reserveTokenA, reserveTokenB, err
}