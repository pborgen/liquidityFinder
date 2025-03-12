package dexUniswapV2Router

import (
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	uniswapv2router "github.com/pborgen/liquidityFinder/abi/V2Router"
	blockchainclient "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainClient"
	blockchainutil "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainutil"
	dexUniswapV2Factory "github.com/pborgen/liquidityFinder/internal/blockchain/dex/v2"
	erc20Helper "github.com/pborgen/liquidityFinder/internal/blockchain/erc20"
	"github.com/pborgen/liquidityFinder/internal/database/model/dex"
	"github.com/rs/zerolog/log"
)


func AddLiquidity(
	auth *bind.TransactOpts,
	routerAddress common.Address,
	tokenA common.Address,
	tokenB common.Address,
	amountADesired *big.Int,
	amountBDesired *big.Int,
	amountAMin *big.Int,
	amountBMin *big.Int,
	to common.Address,
) (bool, error) {
	client := blockchainclient.GetHttpClient()	
	router, err := uniswapv2router.NewUniswapv2routerTransactor(routerAddress, client)

	if err != nil {
		log.Fatal().Msgf("Error creating router instance: %v", err)
	}

	blockchainutil.SetupAuthForTransaction(auth, false)


	log.Info().Msgf("Adding liquidity to router: %s", routerAddress.String())
	log.Info().Msgf("Token A: %s", tokenA.String())
	log.Info().Msgf("Token B: %s", tokenB.String())
	log.Info().Msgf("Amount A Desired: %s", amountADesired.String())
	log.Info().Msgf("Amount B Desired: %s", amountBDesired.String())
	log.Info().Msgf("Amount A Min: %s", amountAMin.String())
	log.Info().Msgf("Amount B Min: %s", amountBMin.String())
	log.Info().Msgf("To: %s", to.String())

	var tx *types.Transaction = nil
	

	attempts := 0
	for attempts < 2 {
		currentBlockTimestamp, err := blockchainutil.GetCurrentBlockTimestamp()
		if err != nil {
			return false, err
		}
	
		deadline := big.NewInt(int64(currentBlockTimestamp + 1020))

		
		tx, err = 
			router.AddLiquidity(
				auth, 
				tokenA, 
				tokenB, 
				amountADesired, 
				amountBDesired, 
				amountAMin, 
				amountBMin, 
				to, 
				deadline,
			)

		if err == nil {
			break
		}

		log.Info().Msgf("Sleeping. Error adding liquidity: %v", err)
		time.Sleep(2 * time.Second)
		attempts++
	}
	
	if err != nil {
		log.Fatal().Msgf("Error adding liquidity: %v", err)
		return false, err
	}

	log.Info().Msgf("Set add liquidity transaction: %s", tx.Hash().String())

	success, err := blockchainutil.GetTransactionResults(auth, tx)
	if err != nil {
		log.Fatal().Msgf("Error getting transaction results: %v", err)
		return false, err
	}

	log.Info().Msgf("Added liquidity: %s", tx.Hash().String())

	return success, nil
}

func RemoveAllLiquidity(
	auth *bind.TransactOpts,
	dexModel dex.ModelDex,
	tokenA common.Address,
	tokenB common.Address,
) (success bool, err error) {

	pairAddress := dexUniswapV2Factory.GetPairAddressByTokens(dexModel, tokenA, tokenB)

	log.Info().Msgf("Remove all liquidity from Pair address: %s", pairAddress.String())
	lpTokenBalance, err := erc20Helper.BalanceOf(pairAddress, auth.From)
	if err != nil {
		return false, err
	}

	erc20Helper.Approve(auth, pairAddress, dexModel.RouterContractAddress, lpTokenBalance)
	
	client := blockchainclient.GetHttpClient()	
	router, err := uniswapv2router.NewUniswapv2routerTransactor(dexModel.RouterContractAddress, client)

	if err != nil {
		return false, err
	}

	blockchainutil.SetupAuthForTransaction(auth, false)

	currentBlockTimestamp, err := blockchainutil.GetCurrentBlockTimestamp()
	if err != nil {
		return false, err
	}

	deadline := big.NewInt(int64(currentBlockTimestamp + 30))


	tx, err := 
		router.RemoveLiquidity(
			auth, 
			tokenA, 
			tokenB, 
			lpTokenBalance, 
			big.NewInt(0), 
			big.NewInt(0), 
			auth.From, 
			deadline,
		)

	if err != nil {
		log.Fatal().Msgf("Error adding liquidity: %v", err)
		return false, err
	}

	if tx == nil {
		return false, errors.New("transaction is nil")
	}

	log.Info().Msgf("Set add liquidity transaction: %s", tx.Hash().String())

	success, err = blockchainutil.GetTransactionResults(auth, tx)
	if err != nil {
		log.Fatal().Msgf("Error getting transaction results: %v", err)
		return false, err
	}

	log.Info().Msgf("Added liquidity: %s", tx.Hash().String())

	return success, nil
}

func SwapExactTokensForTokens(
	auth *bind.TransactOpts,
	dexModel dex.ModelDex,
	pairAddress common.Address,
	amountIn *big.Int,
	amountOutMin *big.Int,
	tokenIn common.Address,
	tokenOut common.Address,
	
	) (success bool, err error) {

	if amountOutMin == nil || amountOutMin.Cmp(big.NewInt(0)) == 0 {
		panic("amountOutMin is required")
	}

	approved, err := erc20Helper.Approve(auth, tokenIn, dexModel.RouterContractAddress, amountIn)
	if err != nil {
		return false, err
	}

	if !approved {
		return false, errors.New("token not approved")
	}
	

	// Check that we have the amountIn tokens
	balance, err := erc20Helper.BalanceOf(tokenIn, auth.From)
	if err != nil {
		return false, err
	}

	if balance.Cmp(amountIn) < 0 {
		return false, errors.New("not enough " + tokenIn.Hex() + " tokens")
	}

	client := blockchainclient.GetHttpClient()	
	router, err := 
		uniswapv2router.NewUniswapv2routerTransactor(dexModel.RouterContractAddress, client)

	if err != nil {
		log.Fatal().Msgf("Error creating router instance: %v", err)
	}
	
	err = blockchainutil.SetupAuthForTransaction(auth, false)
	if err != nil {
		return false, err
	}

	currentBlockTimestamp, err := blockchainutil.GetCurrentBlockTimestamp()
	if err != nil {
		return false, err
	}

	deadline := big.NewInt(int64(currentBlockTimestamp + 1000))

	path := []common.Address{tokenIn, tokenOut}


	log.Info().Msgf("AmountIn: %s", amountIn.String())
	log.Info().Msgf("AmountOutMin: %s", amountOutMin.String())
	tx, err := router.SwapExactTokensForTokens(
		auth,
		amountIn,
		amountOutMin,
		path,
		auth.From,
		deadline,
	)

	if err != nil {
		log.Fatal().Msgf("Error swapping tokens: %v", err)
		return false, err
	}

	log.Info().Msgf("Swapped tokens transaction: %s", tx.Hash().String())

	// success, err = blockchainutil.GetTransactionResults(auth, tx)
	// if err != nil {
	// 	log.Fatal().Msgf("Error getting transaction results: %v", err)
	// 	return false, err
	// }

	// log.Info().Msgf("Swapped tokens transaction: %s", tx.Hash().String())

	return true, nil
}
