package tokenAmountService

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/database/model/token_amount_model"
	"github.com/pborgen/liquidityFinder/internal/myConfig"
	"github.com/pborgen/liquidityFinder/internal/service/transferEventService"
	"github.com/pborgen/liquidityFinder/internal/types"
	"github.com/rs/zerolog/log"
)

var tokenAmountServiceBatchSize = myConfig.GetInstance().TokenAmountServiceBatchSize


func Start() {


	
	largestTokenAmountBlockNumberUpdated, err := GetLargestLastBlockNumberUpdated()
	if err != nil {
		panic(err)
	}

	var fromBlockNumber uint64 = 0
	var toBlockNumber uint64 = 0

	// Special case for the first run
	if largestTokenAmountBlockNumberUpdated == 0 {
		fromBlockNumber = 460000
		toBlockNumber = fromBlockNumber + tokenAmountServiceBatchSize
	} else {
		fromBlockNumber = largestTokenAmountBlockNumberUpdated - tokenAmountServiceBatchSize
		toBlockNumber = fromBlockNumber + tokenAmountServiceBatchSize

		fromBlockNumber, toBlockNumber = 
			calculateNextFromAndToBlockNumbers(
				fromBlockNumber, 
				toBlockNumber,
		)
	}

	for {
		startTime := time.Now()
		
		log.Info().Msgf("Processing block range: %d - %d", fromBlockNumber, toBlockNumber)
		modelTransferEventList, err := transferEventService.GetEventsForBlockRangeOrdered(fromBlockNumber, toBlockNumber)

		if err != nil {
			log.Error().Msgf("Error getting transfers for block range. Retrying... %d, %d, %v", fromBlockNumber, toBlockNumber, err)
			continue
		}

		if len(modelTransferEventList) > 0 {

			startTimeGetModelTokenAmounts := time.Now()
			modelTokenAmounts, err := getModelTokenAmountsFromTransferEvents(modelTransferEventList)
			durationGetModelTokenAmounts := time.Since(startTimeGetModelTokenAmounts).Seconds()
			log.Debug().Msgf("Gathered transfer events: %d", len(modelTransferEventList))
			log.Debug().Msgf("Time taken for getModelTokenAmountsFromTransferEvents: %f seconds", durationGetModelTokenAmounts)
			if err != nil {
				panic(err)
			}

			for _, modelTransferEvent := range modelTransferEventList {
				fromAddress := modelTransferEvent.FromAddress
				toAddress := modelTransferEvent.ToAddress
				contractAddress := modelTransferEvent.ContractAddress
				amount := modelTransferEvent.EventValue

				// Subtract the amount from the from address
				modelTokenAmountFrom := modelTokenAmounts[contractAddress][fromAddress]
				modelTokenAmountFrom.Amount = modelTokenAmountFrom.Amount.Sub(modelTokenAmountFrom.Amount, amount)
				// If amount becomes negative, set it to 0
				if modelTokenAmountFrom.Amount.Sign() < 0 {
					modelTokenAmountFrom.Amount = new(big.Int).SetInt64(0)
				}

				modelTokenAmountFrom.LastBlockNumberUpdated = modelTransferEvent.BlockNumber
				modelTokenAmountFrom.LastLogIndexUpdated = modelTransferEvent.LogIndex

				// Add the amount to the to address
				modelTokenAmountTo := modelTokenAmounts[contractAddress][toAddress]
				modelTokenAmountTo.Amount = modelTokenAmountTo.Amount.Add(modelTokenAmountTo.Amount, amount)

				modelTokenAmountTo.LastBlockNumberUpdated = modelTransferEvent.BlockNumber
				modelTokenAmountTo.LastLogIndexUpdated = modelTransferEvent.LogIndex


				modelTokenAmounts[contractAddress][fromAddress] = modelTokenAmountFrom
				modelTokenAmounts[contractAddress][toAddress] = modelTokenAmountTo
				
			}

			tokenAmounts := make([]types.ModelTokenAmount, 0)
			for _, ownerAddressMap := range modelTokenAmounts {
				for _, modelTokenAmount := range ownerAddressMap {
					tokenAmounts = append(tokenAmounts, *modelTokenAmount)
				}
			}
			startTimeBatchInsert := time.Now()
			err = BatchInsertOrUpdate(tokenAmounts)
			log.Info().Msgf("Inserted %d token amounts", len(tokenAmounts))
			
			if err != nil {
				panic(err)
			}
			durationBatchInsert := time.Since(startTimeBatchInsert).Seconds()
			log.Info().Msgf("Time taken to BatchInsertOrUpdate: %f seconds", durationBatchInsert)
		}

		fromBlockNumberTemp, toBlockNumberTemp := 
			calculateNextFromAndToBlockNumbers(
				fromBlockNumber, 
				toBlockNumber,
			)
	
		if fromBlockNumberTemp == fromBlockNumber && toBlockNumberTemp == toBlockNumber {
			log.Info().Msgf("Reached the largest transfer event block number.")
			break
		} else {
			fromBlockNumber = fromBlockNumberTemp
			toBlockNumber = toBlockNumberTemp
		}

		duration := time.Since(startTime).Seconds()
		
		numberOfBlocksProcessed := toBlockNumber - fromBlockNumber
		blocksPerSecond := float64(numberOfBlocksProcessed) / duration
		log.Info().Msgf(
			"Time taken to process block range: %d - %d (%d blocks): at %f blocks/sec with a total time of %f seconds", 
			fromBlockNumber, toBlockNumber, numberOfBlocksProcessed, blocksPerSecond, duration,
		)
	}
}

func getModelTokenAmountsFromTransferEvents(
	modelTransferEventList []types.ModelTransferEvent) (map[common.Address]map[common.Address]*types.ModelTokenAmount, error) {

	// Get the current token balances

	contractAddressOwnersMap := make(map[common.Address][] common.Address)

	var tokenAddressOwnerAddressList []token_amount_model.TokenAddressOwnerAddress
	
	for _, transfer := range modelTransferEventList {
		contractAddress := transfer.ContractAddress
		fromAddress := transfer.FromAddress
		toAddress := transfer.ToAddress

		tokenAddressOwnerAddressList = append(tokenAddressOwnerAddressList, token_amount_model.TokenAddressOwnerAddress{
			TokenAddress: contractAddress,
			OwnerAddress: fromAddress,
		})
		

		// Used for the initialize
		contractAddressOwnersMap[contractAddress] = append(contractAddressOwnersMap[contractAddress], fromAddress)
		contractAddressOwnersMap[contractAddress] = append(contractAddressOwnersMap[contractAddress], toAddress)
	}

	modelTokenAmounts, err := GetByContractAddressAndOwner(tokenAddressOwnerAddressList)
	
	if err != nil {
		return nil, err
	}

	// initialize model token amounts that were no in the db.
	for contractAddress, ownerAddressList := range contractAddressOwnersMap {
		for _, ownerAddress := range ownerAddressList {

			_, exists := modelTokenAmounts[contractAddress][ownerAddress]

			// Check if we need to initialize
			if !exists {
				toInsert := types.ModelTokenAmount{
					TokenAddress: contractAddress,
					OwnerAddress: ownerAddress,
					Amount: big.NewInt(0),
					LastBlockNumberUpdated: 0,
					LastLogIndexUpdated: -1,
				}
				if _, exists := modelTokenAmounts[contractAddress]; !exists {
					modelTokenAmounts[contractAddress] = make(map[common.Address]*types.ModelTokenAmount)
				}
				modelTokenAmounts[contractAddress][ownerAddress] = &toInsert
			}
		}
	} 

	return modelTokenAmounts, nil
}


func calculateNextFromAndToBlockNumbers(fromBlockNumber uint64, toBlockNumber uint64) (nextFromBlockNumber uint64, nextToBlockNumber uint64) {
	if fromBlockNumber <= 0 {
		panic("fromBlockNumber must be greater than 0")
	}

	if toBlockNumber <= 0 {
		panic("toBlockNumber must be greater than 0")
	}

	maxBlockNumberInTransferEventTable, err := transferEventService.GetLargestBlockNumber()

 	if err != nil {
		panic(err)
	}

	nextFromBlockNumber = toBlockNumber + 1
	nextToBlockNumber = nextFromBlockNumber + tokenAmountServiceBatchSize

	if nextFromBlockNumber > maxBlockNumberInTransferEventTable {
		nextFromBlockNumber = maxBlockNumberInTransferEventTable
	}

	if nextToBlockNumber > maxBlockNumberInTransferEventTable {
		nextToBlockNumber = maxBlockNumberInTransferEventTable
	}


	return nextFromBlockNumber, nextToBlockNumber
}
