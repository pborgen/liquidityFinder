package tokenAmountService

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/database/model/token_amount_model"
	"github.com/pborgen/liquidityFinder/internal/types"
)



func GetLargestLastBlockNumberUpdated() (uint64, error) {
	
	largestLastBlockNumberUpdated, err := token_amount_model.GetLargestLastBlockNumberUpdated()
	if err != nil {
		return 0, err
	}

	return largestLastBlockNumberUpdated, nil
}

func BatchInsertOrUpdate(tokenAmounts []types.ModelTokenAmount) (int64, error) {
	return token_amount_model.BatchInsertOrUpdate(tokenAmounts)
}

func GetByContractAddressAndOwner(contractAddressList []common.Address, ownerList []common.Address) (map[common.Address]map[common.Address]*types.ModelTokenAmount, error) {
	return token_amount_model.GetByContractAddressAndOwner(contractAddressList, ownerList)
}

