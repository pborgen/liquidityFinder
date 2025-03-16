package transferEventService

import (
	"github.com/pborgen/liquidityFinder/internal/database/model/transfer_event_model"
	"github.com/pborgen/liquidityFinder/internal/types"
)



func Insert(transferEvent types.ModelTransferEvent) (int, error) {
	return transfer_event_model.Insert(transferEvent)
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

func GetLargestBlockNumber() (uint64, error) {
	return transfer_event_model.GetLargestBlockNumber()
}

func GetEventsForBlockRangeOrdered(fromBlock uint64, toBlock uint64) ([]types.ModelTransferEvent, error) {
	return transfer_event_model.GetEventsForBlockRangeOrdered(fromBlock, toBlock)
}



