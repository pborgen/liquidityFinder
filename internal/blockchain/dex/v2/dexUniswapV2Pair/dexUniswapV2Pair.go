package dexUniswapV2Pair

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	abi_uniswapv2pair "github.com/pborgen/liquidityFinder/abi/dexV2"
	"github.com/pborgen/liquidityFinder/myConst"

	blockchainclient "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainClient"
	"github.com/pborgen/liquidityFinder/internal/database/model/dex"
	"github.com/pborgen/liquidityFinder/internal/myUtil"
	"github.com/rs/zerolog/log"
)

type Reserve struct {
	Reserve0 big.Int
	Reserve1 big.Int
}


var gasLimitMaxForSkim = 100000
var maxCostInPlsForSkim = new(big.Int).Mul(myConst.GetOneWplsBigint(), big.NewInt(250))

func GetTokenAddressesForPair(dex dex.ModelDex, pairAddress common.Address) struct {
	Token0Address common.Address
	Token1Address common.Address
} {

	client := blockchainclient.GetHttpClient()
	contract, err := abi_uniswapv2pair.NewAbiUniswapv2pairCaller(pairAddress, client)

	if err != nil {
		panic(err)
	}

	maxRetries := 10

	var token0Address, token1Address common.Address

	for attempt := 1; attempt <= maxRetries; attempt++ {

		// Call Contract
		token0Address, err = contract.Token0(nil)

		if err == nil {
			token1Address, err = contract.Token1(nil)
		}

		if err == nil {
			outStruct := new(struct {
				Token0Address common.Address
				Token1Address common.Address
			})

			sortedToken0, sortedToken1 := myUtil.SortTokens(token0Address, token1Address)
			outStruct.Token0Address = sortedToken0
			outStruct.Token1Address = sortedToken1

			return *outStruct
		} else {
			log.Warn().Msgf("Failed to get token addresses for pair: %v", err)
			time.Sleep(10 * time.Second)
		}

	}

	panic("Failed to get token addresses for pair")
}

func estimateGasLimitForSkim(
	pairAddress common.Address,
	to common.Address) (uint64, error) {

	abi, err := abi_uniswapv2pair.AbiUniswapv2pairMetaData.GetAbi()

	if err != nil {
		return 0, err
	}

	data, err := abi.Pack("skim", to)

	if err != nil {
		return 0, err
	}
	
	msg := ethereum.CallMsg{
		From:  common.HexToAddress("0x948bBD062ad9D33E0b06eEeCA5b575e13a0b8f7C"),
		To:    &pairAddress,
		Value: big.NewInt(0), // No ETH is being transferred, just tokens
		Data:  data,
	}
	
	client := blockchainclient.GetHttpClient()
	gasLimit, err := client.EstimateGas(context.Background(), msg)

	if err != nil {
		return 0, err
	}

	return gasLimit, nil
	
}

func GetReservesForPair(
	pairAddress common.Address, 
	tokenA common.Address, 
	tokenB common.Address, 
	retries int) (reserveTokenA *big.Int, reserveTokenB *big.Int, err error) {

	token0, _ := myUtil.SortTokens(tokenA, tokenB)

	reserves, err := populateReservesForPair(pairAddress, retries)
	if err != nil {
		panic(err)
	}

	if tokenA == token0 {
		reserveTokenA = &reserves.Reserve0
		reserveTokenB = &reserves.Reserve1
	} else {
		reserveTokenA = &reserves.Reserve1
		reserveTokenB = &reserves.Reserve0
	}


	return reserveTokenA, reserveTokenB, nil
}

func populateReservesForPair(pairAddress common.Address, maxTries int) (*struct {
	Reserve0 big.Int
	Reserve1 big.Int
}, error) {

	var lastErr error

	for i := 0; i < maxTries; i++ {
		client := blockchainclient.GetHttpClient()
		contract, err := abi_uniswapv2pair.NewAbiUniswapv2pairCaller(pairAddress, client)

		if err != nil {
			lastErr = err
			continue
		}
		
		reserves, err := contract.GetReserves(nil)
		if err != nil {
			lastErr = err
			continue
		}

		outStruct := new(struct {
			Reserve0 big.Int
			Reserve1 big.Int
		})

		outStruct.Reserve0 = *reserves.Reserve0
		outStruct.Reserve1 = *reserves.Reserve1

		return outStruct, nil
	}

	return nil, fmt.Errorf("failed after %d retries, last error: %v", maxTries, lastErr)
}

