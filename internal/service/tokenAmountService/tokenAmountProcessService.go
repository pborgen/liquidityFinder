package tokenAmountService

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/service/transferEventService"
	"github.com/pborgen/liquidityFinder/internal/types"
	"github.com/rs/zerolog/log"
)

const MAX_AMOUNT_OF_BLOCKS_TO_PROCESS = 100

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
		toBlockNumber = fromBlockNumber + MAX_AMOUNT_OF_BLOCKS_TO_PROCESS
	} else {
		fromBlockNumber = largestTokenAmountBlockNumberUpdated - MAX_AMOUNT_OF_BLOCKS_TO_PROCESS
		toBlockNumber = fromBlockNumber + MAX_AMOUNT_OF_BLOCKS_TO_PROCESS

		fromBlockNumber, toBlockNumber = 
			calculateNextFromAndToBlockNumbers(
				fromBlockNumber, 
				toBlockNumber,
		)
	}

	for {
		modelTransferEventList, err := transferEventService.GetEventsForBlockRangeOrdered(fromBlockNumber, toBlockNumber)

		if err != nil {
			log.Error().Msgf("Error getting transfers for block range. Retrying... %d, %d, %v", fromBlockNumber, toBlockNumber, err)
			continue
		}

		if len(modelTransferEventList) > 0 {

			modelTokenAmounts, err := getModelTokenAmountsFromTransferEvents(modelTransferEventList)

			if err != nil {
				panic(err)
			}

			for _, modelTransferEvent := range modelTransferEventList {
				fromAddress := modelTransferEvent.FromAddress
				toAddress := modelTransferEvent.ToAddress
				contractAddress := modelTransferEvent.ContractAddress
				amount := modelTransferEvent.EventValue

				log.Info().Msgf("EventValue: %s", amount.String())
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

				// Only add if balance is greater than 0
				if modelTokenAmountFrom.Amount.Cmp(big.NewInt(0)) > 0 {
					modelTokenAmounts[contractAddress][fromAddress] = modelTokenAmountFrom
				}
				if modelTokenAmountTo.Amount.Cmp(big.NewInt(0)) > 0 {
					modelTokenAmounts[contractAddress][toAddress] = modelTokenAmountTo
				}
			}

			tokenAmounts := make([]types.ModelTokenAmount, 0)
			for _, ownerAddressMap := range modelTokenAmounts {
				for _, modelTokenAmount := range ownerAddressMap {
					tokenAmounts = append(tokenAmounts, *modelTokenAmount)
				}
			}

			_, err = BatchInsertOrUpdate(tokenAmounts)
			if err != nil {
				panic(err)
			}
		}

		fromBlockNumber, toBlockNumber = 
			calculateNextFromAndToBlockNumbers(
				fromBlockNumber, 
				toBlockNumber,
			)
	}
}

func getModelTokenAmountsFromTransferEvents(
	modelTransferEventList []types.ModelTransferEvent) (map[common.Address]map[common.Address]*types.ModelTokenAmount, error) {

	// Get the current token balances
	contractAddressList := make([]common.Address, 0)
	ownerAddressList := make([]common.Address, 0)
	contractAddressOwnersMap := make(map[common.Address][] common.Address)

	for _, transfer := range modelTransferEventList {
		contractAddress := transfer.ContractAddress
		fromAddress := transfer.FromAddress
		toAddress := transfer.ToAddress

		contractAddressList = append(contractAddressList, contractAddress)
		ownerAddressList = append(ownerAddressList, fromAddress, toAddress)

		// Used for the initialize
		contractAddressOwnersMap[contractAddress] = append(contractAddressOwnersMap[contractAddress], fromAddress)
		contractAddressOwnersMap[contractAddress] = append(contractAddressOwnersMap[contractAddress], toAddress)
	}

	modelTokenAmounts, err := GetByContractAddressAndOwner(contractAddressList, ownerAddressList)
	
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
	nextToBlockNumber = nextFromBlockNumber + MAX_AMOUNT_OF_BLOCKS_TO_PROCESS

	if nextFromBlockNumber > maxBlockNumberInTransferEventTable {
		nextFromBlockNumber = maxBlockNumberInTransferEventTable
	}

	if nextToBlockNumber > maxBlockNumberInTransferEventTable {
		nextToBlockNumber = maxBlockNumberInTransferEventTable
	}


	return nextFromBlockNumber, nextToBlockNumber
}
