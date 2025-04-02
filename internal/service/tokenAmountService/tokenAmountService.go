package tokenAmountService

import (
	"context"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/database/model/token_amount_model"
	cacheService "github.com/pborgen/liquidityFinder/internal/service/CacheService"
	"github.com/pborgen/liquidityFinder/internal/types"
)

func GetByTokenAddress(tokenAddress common.Address, limit int, offset int) ([]types.ModelTokenAmount, error) {

	cacheServiceInstance := cacheService.GetInstance()
	cacheStruct := cacheService.CacheType_TokenAmountService_GetByTokenAddress

	cacheKey := cacheStruct.Name + "_" + tokenAddress.Hex() + "_" + strconv.Itoa(limit) + "_" + strconv.Itoa(offset)

	tokenAmounts, err := cacheService.GetObject[[]types.ModelTokenAmount](context.Background(), cacheKey, cacheStruct)
	if err != nil {
		return nil, err
	}

	if tokenAmounts != nil {
		return tokenAmounts, nil
	}
	
	tokenAmounts, err = token_amount_model.GetByTokenAddress(tokenAddress, limit, offset)
	if err != nil {
		return nil, err
	}

	cacheServiceInstance.SetObject(context.Background(), cacheKey, tokenAmounts, cacheStruct)

	return tokenAmounts, nil
}

func GetByOwnerAddress(ownerAddress string, limit int, offset int) ([]types.ModelTokenAmount, error) {

	cacheServiceInstance := cacheService.GetInstance()
	cacheStruct := cacheService.CacheType_TokenAmountService_GetByOwnerAddress

	cacheKey := cacheStruct.Name + "_" + ownerAddress + "_" + strconv.Itoa(limit) + "_" + strconv.Itoa(offset)

	tokenAmounts, err := cacheService.GetObject[[]types.ModelTokenAmount](context.Background(), cacheKey, cacheStruct)
	if err != nil {
		return nil, err
	}

	if tokenAmounts != nil {
		return tokenAmounts, nil
	}

	tokenAmounts, err = token_amount_model.GetByOwnerAddress(ownerAddress, limit, offset)
	if err != nil {
		return nil, err
	}

	cacheServiceInstance.SetObject(context.Background(), cacheKey, tokenAmounts, cacheStruct)
	return tokenAmounts, nil
}

func GetLargestLastBlockNumberUpdated() (uint64, error) {
	
	largestLastBlockNumberUpdated, err := token_amount_model.GetLargestLastBlockNumberUpdated()
	if err != nil {
		return 0, err
	}

	return largestLastBlockNumberUpdated, nil
}

func BatchInsertOrUpdate(tokenAmounts []types.ModelTokenAmount) (error) {
	return token_amount_model.BatchInsertOrUpdate(tokenAmounts)
}

func GetByContractAddressAndOwner(tokenAddressOwnerAddressList []token_amount_model.TokenAddressOwnerAddress) (map[common.Address]map[common.Address]*types.ModelTokenAmount, error) {
	return token_amount_model.GetByContractAddressAndOwner(tokenAddressOwnerAddressList)
}

