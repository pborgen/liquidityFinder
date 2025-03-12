package addLiquidityV2

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pborgen/liquidityFinder/internal/blockchain/dex/v2/dexUniswapV2Router"
	erc20Helper "github.com/pborgen/liquidityFinder/internal/blockchain/erc20"
	"github.com/pborgen/liquidityFinder/internal/database/model/dex"
	"github.com/pborgen/liquidityFinder/internal/myUtil"
)
func Run() {
	
}

func ApproveTokens(
	auth *bind.TransactOpts, 
	routerAddress common.Address, 
	token0Address common.Address, 
	token1Address common.Address, 
	amount0 *big.Int, 
	amount1 *big.Int,
) (success bool, err error) {
	type Result struct {
		Success bool
		Error   error
	}
		
	approve1Channel := myUtil.MyAsync(func() Result {
		success, err = erc20Helper.Approve(auth, token0Address, routerAddress, amount0)
		if err != nil {
			return Result{Success: false, Error: err}
		}
		
		if !success {
			return Result{Success: false, Error: errors.New("Failed to approve token0 at address: " + token0Address.Hex())}
		}
		return Result{Success: true, Error: nil}
	})


	approve2Channel := myUtil.MyAsync(func() Result {
		success, err = erc20Helper.Approve(auth, token1Address, routerAddress, amount1)
		if err != nil {
			return Result{Success: false, Error: err}
		}
		
		if !success {
			return Result{Success: false, Error: errors.New("Failed to approve token1 at address: " + token1Address.Hex())}
		}
		return Result{Success: true, Error: nil}
	})

	approve1Result := <-approve1Channel
	approve2Result := <-approve2Channel

	if approve1Result.Error != nil || approve2Result.Error != nil {
		return false, errors.New("failed to approve tokens")
	}	

	return true, nil
}

func AddInitialLiquidity(
	auth *bind.TransactOpts, 
	shouldCheckApprovals bool,
	routerAddress common.Address, 
	token0Address common.Address, 
	token1Address common.Address, 
	amount0 *big.Int, 
	amount1 *big.Int,
	toAddressForLpTokens common.Address,
) (success bool, err error) {
	
	if shouldCheckApprovals {
		success, err = ApproveTokens(
			auth,
			routerAddress,
			token0Address,
			token1Address,
			amount0,
			amount1,
		)
	}

	if err != nil {
		return false, err
	}

	if !success {
		return false, errors.New("failed to approve tokens")
	}

	success, err = dexUniswapV2Router.AddLiquidity(
		auth,
		routerAddress,
		token0Address, 
		token1Address, 
		amount0, 
		amount1,
		big.NewInt(0), 
		big.NewInt(0), 
		toAddressForLpTokens,
	)

	if err != nil {
		return false, err
	}

	if !success {
		return false, errors.New("Failed to add liquidity")
	}

	return true, nil
}

func RemoveAllLiquidity(auth *bind.TransactOpts, dexModel dex.ModelDex, token0Address common.Address, token1Address common.Address) (success bool, err error) {
	success, err = dexUniswapV2Router.RemoveAllLiquidity(
		auth, 
		dexModel, 
		token0Address, 
		token1Address,
	)

	if err != nil {
		return false, err
	}

	return success, nil
}