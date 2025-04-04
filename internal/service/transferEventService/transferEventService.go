package transferEventService

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/database/model/transfer_event_model"
	cacheService "github.com/pborgen/liquidityFinder/internal/service/CacheService"
	"github.com/pborgen/liquidityFinder/internal/types"
	"github.com/rs/zerolog/log"
)

func GetAllForAddressGroupBy(address common.Address, viewMode string) ([]types.TransferEventGroupBy, error) {
	cacheServiceInstance := cacheService.GetInstance()
	cacheStruct := cacheService.CacheType_TransferEventService_GetAllForAddressGroupBy
	cacheKey := cacheStruct.Name + "_" + address.Hex() + "_" + viewMode

	// Try to get from cache first
	results, err := cacheService.GetObject[[]types.TransferEventGroupBy](context.Background(), cacheKey, cacheStruct)
	if err != nil {
		return nil, err
	}

	if results != nil {
		return results, nil
	}

	// If not in cache, get from database
	results, err = transfer_event_model.GetAllForAddressGroupBy(address, viewMode)
	if err != nil {
		return nil, err
	}

	// Store in cache
	err = cacheServiceInstance.SetObject(context.Background(), cacheKey, results, cacheStruct)
	if err != nil {
		// Log error but don't fail the request
		log.Error().Err(err).Msg("Failed to cache transfer events")
	}

	return results, nil
}

func GetAllForAddress(address common.Address, limit int, offset int) ([]types.ModelTransferEvent, error) {
	return transfer_event_model.GetAllForAddress(address, limit, offset)
}

func BatchInsertOrUpdate(transferEvents []types.ModelTransferEvent) ([]int, error) {
	return transfer_event_model.BatchInsertOrUpdate(transferEvents)
}


func GetById(id int) (*types.ModelTransferEvent, error) {
	return transfer_event_model.GetById(id)
}

func GetSmallestBlockNumber() (uint64, error) {
	return transfer_event_model.GetSmallestBlockNumber()
}

func DoAnyRowsExists() (bool, error) {
	return transfer_event_model.DoAnyRowsExists()
}

func GetLargestBlockNumber() (uint64, error) {
	return transfer_event_model.GetLargestBlockNumber()
}

func GetEventsForBlockRangeOrdered(fromBlock uint64, toBlock uint64) ([]types.ModelTransferEvent, error) {
	return transfer_event_model.GetEventsForBlockRangeOrdered(fromBlock, toBlock)
}



