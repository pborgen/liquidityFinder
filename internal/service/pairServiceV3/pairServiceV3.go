package pairServiceV3

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"

	"math/big"
	"strconv"
	"time"

	"github.com/pborgen/liquidityFinder/internal/blockchain/blockchainutil"
	"github.com/pborgen/liquidityFinder/internal/blockchain/dex/dexV3"
	erc20Helper "github.com/pborgen/liquidityFinder/internal/blockchain/erc20"
	"github.com/pborgen/liquidityFinder/internal/database/model/dex"
	"github.com/pborgen/liquidityFinder/internal/database/model/pair"
	"github.com/pborgen/liquidityFinder/internal/service/erc20Service"
	"github.com/pborgen/liquidityFinder/internal/service/nameValueService"
	"github.com/pborgen/liquidityFinder/internal/types"
	"github.com/pborgen/liquidityFinder/myConst"
)


func GetReservesForPair(
	pairAddress common.Address, 
	tokenA common.Address, 
	tokenB common.Address, 
	retries int) (reserveTokenA *big.Int, reserveTokenB *big.Int, err error) {
	
	balanceTokenA, err := erc20Helper.BalanceOf(tokenA, pairAddress)

	if err != nil {
		return nil, nil, err
	}

	balanceTokenB, err := erc20Helper.BalanceOf(tokenB, pairAddress)

	if err != nil {
		return nil, nil, err
	}
	
	
	return balanceTokenA, balanceTokenB, nil
}

func AddNewPairForAllV3Dexes(networkId int) {

	allV3Dexes := dex.GetAllV3Dexes(networkId)

	for _, element := range allV3Dexes {
		err := AddNewPairs(element.DexId)

		if err != nil {
			panic(err)
		}
	}
}

func AddNewPairs(dexId int) error {


	nameValue := "LAST_BLOCK_WHERE_NEW_POOLS_PROCESSED_BY_DEX_ID=" + strconv.Itoa(dexId)

	modelDex := dex.GetById(dexId)

	if modelDex.DexType != dex.UniswapV3 {
		panic("Dex is not UniswapV3")
	}

	if !nameValueService.Exists(nameValue) {
		nameValueService.Insert(nameValue, "0", nameValueService.DataTypeInt)
	}

	maxBlockBatch := uint64(1000)
	currentLatestBlock := blockchainutil.GetCurrentBlockNumberAsUInt64()

	lastBlockProcessedTemp, err := nameValueService.GetValueInt(nameValue)
	lastProcessedBlock := uint64(lastBlockProcessedTemp)


	if err != nil {
		return err
	}

	if lastProcessedBlock >= currentLatestBlock {
		return nil
	}
	
	for lastProcessedBlock < currentLatestBlock {
		fromBlock := lastProcessedBlock
		var toBlock uint64

		if lastProcessedBlock + maxBlockBatch > currentLatestBlock {
			toBlock = currentLatestBlock
		} else {
			toBlock = lastProcessedBlock + maxBlockBatch
		}

		log.Info().Msgf("Processing pairs from block %d to block %d", fromBlock, toBlock)

		v3PairCreatedList, err := 
			dexV3.GetPairsCreated(
				modelDex.FactoryContractAddress, 
				fromBlock, 
				toBlock,
			)

		if err != nil {
			return err
		}
		
		processPairs(v3PairCreatedList, modelDex.NetworkId, dexId)
		
		nameValueService.UpdateValueInt(nameValue, int(toBlock))
		lastProcessedBlock = toBlock

		currentLatestBlock = blockchainutil.GetCurrentBlockNumberAsUInt64()
	}
	
	log.Info().Msgf("Finished processing pairs for dexId %d", dexId)
	
	return nil
}

func processPairs(
	v3PairCreatedList []dexV3.V3PairCreated, 
	networkId int, 
	dexId int) {

	for _, element := range v3PairCreatedList {

		// For some reason this address is of a nft v3 position
		if element.Token0 == common.HexToAddress("0xC36442b4a4522E871399CD717aBDD847Ab11FE88") ||
			element.Token1 == common.HexToAddress("0xC36442b4a4522E871399CD717aBDD847Ab11FE88") ||
			element.Token0 == common.HexToAddress("0xEB9951021698B42e4399f9cBb6267Aa35F82D59D") ||
			element.Token1 == common.HexToAddress("0xEB9951021698B42e4399f9cBb6267Aa35F82D59D") {
			continue
		}

		exists := pair.ExistsByContractAddress(element.PairAddress)

		if exists {
			continue
		} else {
			isPlsPair := element.Token0 == 
				myConst.GetWplsAddress() || element.Token1 == myConst.GetWplsAddress()

			erc20Token0, errToken0 := erc20Service.GetByContractAddress(element.Token0, networkId)
			
			if errToken0 != nil {
				if errToken0.Error() == "no contract code at given address" || 
				   strings.HasPrefix(errToken0.Error(), "could not get the name of the token:") ||
				   strings.HasPrefix(errToken0.Error(), "could not get the decimals of the token:") {
						
					log.Error().Msgf("Token with address %s not found", element.Token0)
					continue
				} else {
					continue
				}
			}

			erc20Token1, errToken1 := erc20Service.GetByContractAddress(element.Token1, networkId)

			if errToken1 != nil {

				if errToken1.Error() == "no contract code at given address" ||
					strings.HasPrefix(errToken1.Error(), "could not get the name of the token:") ||
					strings.HasPrefix(errToken1.Error(), "could not get the decimals of the token:") {

					log.Error().Msgf("Token with address %s not found", element.Token1)
					continue
				} else {
					continue
				}
			}

			temp := types.ModelPair{
				DexId:               dexId,
				PairIndex:           -1,
				PairContractAddress: element.PairAddress,
				Token0Erc20Id:       erc20Token0.Erc20Id,
				Token1Erc20Id:       erc20Token1.Erc20Id,
				Token0Address:       element.Token0,
				Token1Address:       element.Token1,
				Token0Reserves:      big.Int{},
				Token1Reserves:      big.Int{},
				ShouldFindArb:       true, // set the correct value
				HasTaxToken:    	 false,
				IsPlsPair:           isPlsPair,
				IsHighLiquidity:     false,
				UniswapV3Fee:        element.Fee,
				UniswapV3TickSpacings: element.TickSpacing,
				InsertedAt:          time.Now(),
				LastUpdated:         time.Now(),
				LastTimeReservesUpdated: time.Unix(0, 0),
			}


			_, err := pair.Insert(temp)

			if err != nil {
				log.Error().Msgf("Error inserting pair address %s", element.PairAddress)
				log.Error().Msgf("token: %s", err)
				panic(err)
			}


		}
	}

}

