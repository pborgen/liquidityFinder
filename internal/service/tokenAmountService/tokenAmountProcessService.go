package tokenAmountService

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pborgen/liquidityFinder/internal/service/transferEventService"
	"github.com/pborgen/liquidityFinder/internal/types"
	"github.com/rs/zerolog/log"
)


func Start() {


	maxAmountOfBlocksToProcess := uint64(1000)
	largestLastBlockNumberUpdated, err := GetLargestLastBlockNumberUpdated()
	if err != nil {
		panic(err)
	}

	var fromBlock uint64 = 0
	var toBlock uint64 = 0
	if largestLastBlockNumberUpdated == 0 {
		fromBlock = 1_000_000
		toBlock = fromBlock + maxAmountOfBlocksToProcess
	} else {
		if largestLastBlockNumberUpdated > maxAmountOfBlocksToProcess {
			fromBlock = largestLastBlockNumberUpdated - (maxAmountOfBlocksToProcess)
			toBlock = fromBlock + maxAmountOfBlocksToProcess
		}
	}

	


	for {
		modelTransferEventList, err := transferEventService.GetEventsForBlockRangeOrdered(fromBlock, toBlock)

		if err != nil {
			log.Error().Msgf("Error getting transfers for block range. Retrying... %d, %d, %v", fromBlock, toBlock, err)
			continue
		}

		if len(modelTransferEventList) == 0 {
			fromBlock = toBlock + 1
			toBlock = fromBlock + maxAmountOfBlocksToProcess - 1
			continue
		}

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

			modelTokenAmounts[contractAddress][fromAddress] = modelTokenAmountFrom
			modelTokenAmounts[contractAddress][toAddress] = modelTokenAmountTo
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



