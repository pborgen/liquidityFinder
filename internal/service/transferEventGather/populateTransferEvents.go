package transferEventGather

import (
	"context"
	"log"

	blockchainclient "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainClient"
	blockchainutil "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainutil"
	"github.com/pborgen/liquidityFinder/internal/service/transferEventService"
	"github.com/pborgen/liquidityFinder/internal/types"
)

// @param map[common.Address]bool - a map of pairs that we are monitoring
// @param bool - if true, we will only monitor the pairs in the map

func Start() {

	client := blockchainclient.GetHttpClient()

	transferEventGather, err := NewTokenEventTracker(client)
	if err != nil {
		panic(err)
	}

	maxAmountOfBlocksToProcess := uint64(100)
	largestBlockNumber, err := transferEventService.GetLargestBlockNumber()
	if err != nil {
		panic(err)
	}

	var fromBlock uint64 = 0

	if largestBlockNumber > maxAmountOfBlocksToProcess {
		fromBlock = largestBlockNumber - (maxAmountOfBlocksToProcess + 1)
	}


	if fromBlock == 0 {
		fromBlock = 460_000 // transfer events start somewhere after this block
	}


	var toBlock uint64 
	
	if fromBlock + maxAmountOfBlocksToProcess > blockchainutil.GetCurrentBlockNumber() {
		toBlock = blockchainutil.GetCurrentBlockNumber()
	} else {
		toBlock = fromBlock + maxAmountOfBlocksToProcess
	}


	for {
		transfers, err := 
			transferEventGather.GetAllTransferEventsForBlockRange(context.Background(), fromBlock, toBlock)


		if err != nil {
			log.Println("Error getting transfers for block range. Retrying... ", fromBlock, toBlock, err)
			continue
		}

		hasError := false

		transferEvents := []types.ModelTransferEvent{}
		for index, transfer := range transfers {
			
			transferEvents = append(transferEvents, types.ModelTransferEvent{
				BlockNumber: transfer.Block,
				TransactionHash: transfer.TxHash.Hex(),
				LogIndex: index,
				ContractAddress: transfer.ContractAddress,
				FromAddress: transfer.From,
				ToAddress: transfer.To,
				EventValue: transfer.Value,
			})
		}

		_, err = transferEventService.BatchInsertOrUpdate(transferEvents)
		if err != nil {
			log.Println("Error inserting transfer events. Retrying... ", fromBlock, toBlock, err)
			panic(err)
		}

		if hasError {
			log.Println("Error inserting transfer event. Retrying... ", fromBlock, toBlock)
			continue
		} else {
	
			fromBlock = toBlock + 1

			if fromBlock + maxAmountOfBlocksToProcess > blockchainutil.GetCurrentBlockNumber() {
				toBlock = blockchainutil.GetCurrentBlockNumber()
			} else {
				toBlock = fromBlock + maxAmountOfBlocksToProcess
			}
		}
	}
}

