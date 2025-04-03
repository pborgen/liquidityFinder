package transferEventService

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/database/model/transfer_event_model"
	"github.com/pborgen/liquidityFinder/internal/types"
)

func GetAllForAddressGroupBy(address common.Address) ([]types.TransferEventGroupBy, error) {
	return transfer_event_model.GetAllForAddressGroupBy(address)
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



