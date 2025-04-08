package transferEventGather

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"math/big"

	blockchainclient "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainClient"
	blockchainutil "github.com/pborgen/liquidityFinder/internal/blockchain/blockchainutil"
	"github.com/pborgen/liquidityFinder/internal/myConfig"
	"github.com/pborgen/liquidityFinder/internal/service/transferEventService"
	"github.com/pborgen/liquidityFinder/internal/types"
	"github.com/rs/zerolog/log"
)

// @param map[common.Address]bool - a map of pairs that we are monitoring
// @param bool - if true, we will only monitor the pairs in the map
var BATCH_SIZE = myConfig.GetInstance().GetTransferEventGatherBatchSize()

func Start() {

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create channel for signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start a goroutine to handle shutdown signals
	go func() {
		sig := <-sigChan
		log.Info().Msgf("Received shutdown signal: %v", sig)
		cancel()
	}()
	
	client := blockchainclient.GetHttpClient()

	transferEventGather, err := NewTokenEventTracker(client)
	if err != nil {
		panic(err)
	}

	maxAmountOfBlocksToProcess := uint64(BATCH_SIZE)

	doAnyRowsExist, err := transferEventService.DoAnyRowsExists()
	if err != nil {
		panic(err)
	}

	var fromBlock uint64 = 0
	var largestBlockNumber uint64 = 0

	if doAnyRowsExist {
		largestBlockNumber, err = transferEventService.GetLargestBlockNumber()
		if err != nil {
			panic(err)
		}
	} else {
		largestBlockNumber = 0
	}


	if largestBlockNumber > maxAmountOfBlocksToProcess {
		fromBlock = largestBlockNumber - (maxAmountOfBlocksToProcess + 1)
	}


	if fromBlock == 0 {
		fromBlock = 462_000 // transfer events start somewhere after this block
	}


	var toBlock uint64 
	
	if fromBlock + maxAmountOfBlocksToProcess > blockchainutil.GetCurrentBlockNumber() {
		toBlock = blockchainutil.GetCurrentBlockNumber()
	} else {
		toBlock = fromBlock + maxAmountOfBlocksToProcess
	}

	
	for {
		start := time.Now()

		// Check if context is done before processing
		if ctx.Err() != nil {
			log.Info().Msg("Shutting down transfer event gatherer...")
			return
		}
		
		gatherEventsTimeStart := time.Now()
		transfers, err := 
			transferEventGather.GetAllTransferEventsForBlockRange(context.Background(), fromBlock, toBlock)
		gatherEventsTimeEnd := time.Now()
		gatherEventsTime := gatherEventsTimeEnd.Sub(gatherEventsTimeStart).Seconds()
		log.Debug().Msgf("Gathering events time: %.2f seconds", gatherEventsTime)

		if err != nil {
			log.Error().Msgf("Error getting transfers for block range. Retrying... %d, %d, %v", fromBlock, toBlock, err)
			time.Sleep(20 * time.Second)
			continue
		}

		shouldInsert := false
		hasError := false

		if len(transfers) > 0 {
			shouldInsert = true
		}

		if shouldInsert {
			transferEvents := []types.ModelTransferEvent{}
			for index, transfer := range transfers {
				
				// Do not insert transfer events with value 0
				if transfer.Value.Cmp(big.NewInt(0)) == 0 {
					continue
				}

				transferEvents = append(transferEvents, types.ModelTransferEvent{
					TransactionHash: transfer.TxHash,
					BlockNumber: transfer.Block,
					LogIndex: index,
					ContractAddress: transfer.ContractAddress,
					FromAddress: transfer.From,
					ToAddress: transfer.To,
					EventValue: transfer.Value,
				})
			}

			numberOfEvents := len(transferEvents)
			log.Debug().Msgf("Number of events to insert: %d", numberOfEvents)

			insertTimeStart := time.Now()
			_, err = transferEventService.BatchInsertOrUpdate(transferEvents)
			insertTimeEnd := time.Now()
			insertTime := insertTimeEnd.Sub(insertTimeStart).Seconds()
			log.Debug().Msgf("Inserting events time: %.2f seconds", insertTime)

			if err != nil {
				log.Error().Msgf("Error inserting transfer events. Retrying... %d, %d, %v", fromBlock, toBlock, err)
				panic(err)
			}
		}

		if hasError {
			log.Error().Msgf("Error inserting transfer event. Retrying... %d, %d", fromBlock, toBlock)
			continue
		} else {
	

			totalTimeSeconds := time.Since(start).Seconds()
			blocksPerSecond := float64(toBlock - fromBlock) / totalTimeSeconds

			log.Debug().Msgf("Processed block range: %d - %d, %.2f blocks/sec, %.2f seconds", fromBlock, toBlock, blocksPerSecond, totalTimeSeconds)

			fromBlock = toBlock + 1

			if fromBlock + maxAmountOfBlocksToProcess > blockchainutil.GetCurrentBlockNumber() {
				toBlock = blockchainutil.GetCurrentBlockNumber()
			} else {
				toBlock = fromBlock + maxAmountOfBlocksToProcess
			}

			if fromBlock > toBlock {
				log.Info().Msg("No more blocks to process. Exiting...")
				return
			}

			if fromBlock == toBlock {
				log.Info().Msg("No more blocks to process. Exiting...")
				return
			}
		}
	}
}

