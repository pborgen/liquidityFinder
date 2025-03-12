package pairService

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	dexUniswapV2Factory "github.com/pborgen/liquidityFinder/internal/blockchain/dex/v2"
	"github.com/pborgen/liquidityFinder/internal/blockchain/dex/v2/dexUniswapV2Pair"
	"github.com/pborgen/liquidityFinder/internal/database/model/dex"
	"github.com/pborgen/liquidityFinder/internal/database/model/erc20"
	"github.com/pborgen/liquidityFinder/internal/database/model/pair"
	"github.com/pborgen/liquidityFinder/internal/myUtil"
	cacheService "github.com/pborgen/liquidityFinder/internal/service/CacheService"
	"github.com/pborgen/liquidityFinder/internal/service/erc20Service"
	"github.com/pborgen/liquidityFinder/internal/service/pairServiceV2"
	"github.com/pborgen/liquidityFinder/internal/service/pairServiceV3"

	"github.com/pborgen/liquidityFinder/internal/types"
	"github.com/pborgen/liquidityFinder/myConst"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

var highLiquidityThreshold *big.Int

func init() {
	highLiquidityThreshold = new(big.Int).Mul(myConst.GetOneWplsBigint(), big.NewInt(25_000))
}

func PlsPairWithHighAmountOfPls(
	dexIds []int, 
	minAmountOfPls *big.Int, 
	shouldHydrate bool) []types.ModelPair {

	return pair.PlsPairWithHighAmountOfPls(dexIds, minAmountOfPls, shouldHydrate)
}

func NonPlsPairWithHighLiquidity(
	dexIds []int, 
	shouldHydrate bool) ([]types.ModelPair, error) {

	return pair.GetAllNonPlsPairsWithHighLiquidity(dexIds, false)
}

func GetPairsByPairContractAddresses(
	pairContractAddresses []common.Address, 
	shouldHydrate bool) ([]types.ModelPair, error) {

	return pair.GetPairsByPairContractAddresses(pairContractAddresses, shouldHydrate)
}

func GetPairsByPairContractAddressesAsMap(
	pairContractAddresses []common.Address, 
	shouldHydrate bool) (map[common.Address]types.ModelPair, error) {
		
	modelPairs, err := pair.GetPairsByPairContractAddresses(pairContractAddresses, shouldHydrate)

	if err != nil {
		return nil, err
	}

	modelPairMap := make(map[common.Address]types.ModelPair)

	for _, modelPair := range modelPairs {
		modelPairMap[modelPair.PairContractAddress] = modelPair
	}

	return modelPairMap, nil
}

func FixPairOrdering() {
	allPairs, err := GetAllPairs(false, false)

	if err != nil {
		log.Error().Msgf("Error getting all pairs: %v", err)
		panic("Error getting all pairs")
	}

	log.Info().Msgf("Found %d pairs", len(allPairs))

	counter := 0

	startTime := time.Now()

	wg := sync.WaitGroup{}
	sem := semaphore.NewWeighted(250)

	for _, modelPair := range allPairs {
		token0Address := modelPair.Token0Address
		token1Address := modelPair.Token1Address

		sortedToken0Address, sortedToken1Address := myUtil.SortTokens(token0Address, token1Address)
	

		if sortedToken0Address != token0Address || sortedToken1Address != token1Address {

			wg.Add(1) 
			sem.Acquire(context.Background(), 1)

			go func() {
				defer sem.Release(1)
				defer wg.Done()

				modelPair.Token0Address = sortedToken0Address
				modelPair.Token1Address = sortedToken1Address

				token0Erc20, err := erc20.GetByContractAddress(sortedToken0Address, 3)
				if err != nil {
					panic(err)
				}
				modelPair.Token0Erc20Id = token0Erc20.Erc20Id
				
				token1Erc20, err := erc20.GetByContractAddress(sortedToken1Address, 3)
				if err != nil {
					panic(err)
				}
				modelPair.Token1Erc20Id = token1Erc20.Erc20Id

				pair.UpdateForFix(modelPair)
			}()

			// Sleep to not blast the db too much
			time.Sleep(10 * time.Millisecond)

			log.Info().Msgf("Pair %s has tokens %s and %s, sorted to %s and %s", modelPair.PairContractAddress.String(), token0Address.String(), token1Address.String(), sortedToken0Address.String(), sortedToken1Address.String())
			
		} else {
			log.Info().Msgf("Pair %s is already sorted", modelPair.PairContractAddress.String())
		}

		counter++
		if counter % 100 == 0 {
			amountOfTimeTook := time.Since(startTime)
			amountOfTimeTookInSeconds := amountOfTimeTook.Seconds()
			amountOfPairsPerSecond := float64(counter) / amountOfTimeTookInSeconds

			log.Info().Msgf(
				"Processed %d pairs at a rate of %f pairs per second", 
				counter, 
				amountOfPairsPerSecond)
		}

		//pair.UpdateReserves(pair.PairId, pair.Token0Reserves, pair.Token1Reserves)
	}
}

func GetAllPairsWithLimit(limit int, shouldHydrate bool) ([]types.ModelPair, error) {
	cacheServiceInstance := cacheService.GetInstance()
	cacheStruct := cacheService.CacheType_PairService_AllPairs_WithLimit
	cacheKey := cacheStruct.Name

	returnValue, err := cacheService.GetObject[[]types.ModelPair](context.Background(), cacheKey, cacheStruct)

	if err == nil && returnValue != nil && len(returnValue) > 0 {
		return returnValue, nil
	}

	returnValue, err = pair.GetAllWithLimit(limit, shouldHydrate)

	if err != nil {
		return nil, err
	}

	cacheServiceInstance.SetObject(context.Background(), cacheKey, returnValue, cacheStruct)

	return returnValue, nil
}

func GetAllPairs(shouldHydrate bool, shouldUseCache bool) ([]types.ModelPair, error) {

	var returnValue []types.ModelPair
	var err error
	var cacheServiceInstance *cacheService.CacheService
	var cacheStruct cacheService.CacheType
	var cacheKey string

	if shouldUseCache {
		cacheServiceInstance = cacheService.GetInstance()
		cacheStruct = cacheService.CacheType_PairService_AllPairs
		cacheKey = cacheStruct.Name

		returnValue, err = cacheService.GetObject[[]types.ModelPair](context.Background(), cacheKey, cacheStruct)

		if err == nil && returnValue != nil && len(returnValue) > 0 {
			return returnValue, nil
		}
	}

	returnValue, err = pair.GetAll(shouldHydrate)

	if err != nil {
		return nil, err
	}

	if shouldUseCache {
		cacheServiceInstance.SetObject(context.Background(), cacheKey, returnValue, cacheStruct)
	}

	return returnValue, nil
}

func GetAllPlsHighLiquidityPairsThatAreNotTaxToken(shouldHydrate bool) ([]types.ModelPair, error) {
	pairs, err := GetAllHighLiquidityPairsThatAreNotTaxToken(shouldHydrate)

	if err != nil {
		return nil, err
	}

	plsPairs := []types.ModelPair{}

	for _, pair := range pairs {
		if isPlsPair(pair) {
			plsPairs = append(plsPairs, pair)
		}
	}

	return plsPairs, nil
}

func GetAllNonPlsHighLiquidityPairsThatAreNotTaxToken(shouldHydrate bool) ([]types.ModelPair, error) {
	pairs, err := GetAllHighLiquidityPairsThatAreNotTaxToken(shouldHydrate)

	if err != nil {
		return nil, err
	}

	plsPairs := []types.ModelPair{}

	for _, pair := range pairs {
		if !isPlsPair(pair) {
			plsPairs = append(plsPairs, pair)
		}
	}

	return plsPairs, nil
}

func GetAllHighLiquidityPairsThatAreNotTaxToken(shouldHydrate bool) ([]types.ModelPair, error) {
	cacheServiceInstance := cacheService.GetInstance()
	cacheStruct := cacheService.CacheType_PairService_AllHighLiquidityPairs
	cacheKey := cacheStruct.Name

	returnValue, err := cacheService.GetObject[[]types.ModelPair](context.Background(), cacheKey, cacheStruct)

	if err == nil && returnValue != nil && len(returnValue) > 0 {
		return returnValue, nil
	}

	returnValue, err = pair.GetAllHighLiquidityPairsThatAreNotTaxToken(shouldHydrate)

	if err != nil {
		return nil, err
	}

	cacheServiceInstance.SetObject(context.Background(), cacheKey, returnValue, cacheStruct)

	return returnValue, nil
}

func GetPairsByContractAddresses(pairContractAddresses []common.Address) ([]types.ModelPair, error) {
	return pair.GetPairsByContractAddresses(pairContractAddresses)
}

func GetByContractAddress(pairContractAddress common.Address) (*types.ModelPair, error) {
	return pair.GetByContractAddress(pairContractAddress)
}

func GetSpecificPair(dex dex.ModelDex, token0ContractAddress common.Address, token1ContractAddress common.Address) (types.ModelPair, error) {
	
	return pair.GetSpecificPair(dex, token0ContractAddress, token1ContractAddress)
}

func UpdateIsHighLiquidityForAllPairs(shouldUseCache bool) {
	allPairs, err := GetAllPairs(false, shouldUseCache)

	if err != nil {
		log.Error().Msgf("Error getting all pairs: %v", err)
		panic("Error getting all pairs")
	}

	var wg sync.WaitGroup
	counter := 0
	mu := sync.Mutex{} // Mutex to protect shared resources

	concurrencyLimit := 10 // Set your desired concurrency limit
	sem := make(chan struct{}, concurrencyLimit)


	for _, element := range allPairs {
		wg.Add(1)

		sem <- struct{}{}

		go func(element types.ModelPair) {
			defer wg.Done()
			defer func() { <-sem }()

			UpdateReservesForPair(&element, 10)
			UpdateIsHighLiquidityForPair(element, allPairs)

			mu.Lock()
			counter++
			if counter % 100 == 0 {
					log.Info().Msgf("Processed %d pairs", counter)
			}
			mu.Unlock()
		
		}(element)
	}

	wg.Wait()
	log.Info().Msgf("Completed updating is high liquidity for all pairs")
}

func UpdateIsHighLiquidityForPair(element types.ModelPair, allPairs []types.ModelPair) {
	if isPlsPair(element) {
		wplsReserves := big.Int{}

		if element.Token0Address == common.HexToAddress(myConst.WPLS_ADDRESS_STRING) {
			wplsReserves = element.Token0Reserves
		} else if element.Token1Address == common.HexToAddress(myConst.WPLS_ADDRESS_STRING) {
			wplsReserves = element.Token1Reserves
		} else {
			panic("Error one of the tokens should have been WPLS")
		}

		isHighLiquidity := wplsReserves.Cmp(highLiquidityThreshold) > 0

		pair.UpdateIsHighLiquidity(element.PairId, isHighLiquidity)
		
	} else {
		// find all neighbors of this element that have pls. 
		allNeighborsThatHavePls := []types.ModelPair{}

		for _, neighbor := range allPairs {
			nonCurrentPair := element.PairContractAddress != neighbor.PairContractAddress
			if isPlsPair(neighbor) && nonCurrentPair && isNeighbor(element, neighbor) {
				allNeighborsThatHavePls = append(allNeighborsThatHavePls, neighbor)
			}
		}

		for _, neighborWithPls := range allNeighborsThatHavePls {
			wplsReserves := big.Int{}
			if neighborWithPls.Token0Address == common.HexToAddress(myConst.WPLS_ADDRESS_STRING) {
				wplsReserves = neighborWithPls.Token0Reserves
			} else if neighborWithPls.Token1Address == common.HexToAddress(myConst.WPLS_ADDRESS_STRING) {
				wplsReserves = neighborWithPls.Token1Reserves
			} else {
				panic("Error one of the tokens should have been WPLS")
			}

			if wplsReserves.Cmp(highLiquidityThreshold) > 0 {
				pair.UpdateIsHighLiquidity(element.PairId, true)
				break
			}
		}
	}
}

func isNeighbor(element types.ModelPair, neighbor types.ModelPair) bool {
	return element.Token0Address == neighbor.Token0Address || 
		   element.Token0Address == neighbor.Token1Address ||
		   element.Token1Address == neighbor.Token0Address ||
		   element.Token1Address == neighbor.Token1Address
}

func isPlsPair(pair types.ModelPair) bool {
	return pair.Token0Address == common.HexToAddress(myConst.WPLS_ADDRESS_STRING) || 
		   pair.Token1Address == common.HexToAddress(myConst.WPLS_ADDRESS_STRING)
}



func UpdateReservesForAllPairs(shouldUseCache bool) {
	allPairs, err := GetAllPairs(false, shouldUseCache)

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(3) // Limit to 5 concurrent goroutines


	if err != nil {
		log.Error().Msgf("Error getting all pairs: %v", err)
		panic("Error getting all pairs")
	}

	counter := 0
	for _, pair := range allPairs {

		wg.Add(1)

		// Acquire a semaphore slot
		if err := sem.Acquire(context.Background(), 1); err != nil {
			log.Error().Msgf("Failed to acquire semaphore: %v", err)
			wg.Done()
			continue
		}
		
		go func(pair types.ModelPair) { // Replace `pairType` with the actual type of `pair`
			defer wg.Done()
			defer sem.Release(1) // Release the semaphore slot when done

			UpdateReservesForPair(&pair, 10)
		}(pair)

		counter++
		if counter % 100 == 0 {
			log.Info().Msgf("Processed %d pairs", counter)
		}
	}
}

func PopulatePairsInDb(dex dex.ModelDex) {

	currentLength := dexUniswapV2Factory.GetAllPairsLength(dex)
	largestPairIndex, err := pair.GetLargestPairIndex(dex)

	if err != nil {
		panic(err)
	}

	counter := largestPairIndex

	// Add the new pairs
	for pairIndex := largestPairIndex + 1; pairIndex < currentLength; pairIndex++ {

		iAsBigInt := *big.NewInt(int64(pairIndex))

		pairAddress := dexUniswapV2Factory.GetPairAddress(dex, iAsBigInt)

		// Get the 2 taken Addresses for the pair
		result := dexUniswapV2Pair.GetTokenAddressesForPair(dex, pairAddress)
		token0Address := result.Token0Address
		token1Address := result.Token1Address

		erc20Token0, errToken0 := erc20Service.GetByContractAddress(token0Address, dex.NetworkId)
		erc20Token1, errToken1 := erc20Service.GetByContractAddress(token1Address, dex.NetworkId)

		if errToken0 == nil && errToken1 == nil {
			isPlsPair := 
				myUtil.IsWPLS(token0Address)|| 
				myUtil.IsWPLS(token1Address)


			temp := types.ModelPair{
				DexId:               dex.DexId,
				PairIndex:           pairIndex,
				PairContractAddress: pairAddress,
				Token0Erc20Id:       erc20Token0.Erc20Id,
				Token1Erc20Id:       erc20Token1.Erc20Id,
				Token0Address:       token0Address,
				Token1Address:       token1Address,
				Token0Reserves:      big.Int{},
				Token1Reserves:      big.Int{},
				ShouldFindArb:       false, // set the correct value
				HasTaxToken:    	 false,
				IsPlsPair:           isPlsPair,
				IsHighLiquidity:     false,
				UniswapV3Fee:        big.Int{},
				UniswapV3TickSpacings: big.Int{},
				InsertedAt:          time.Now(),
				LastUpdated:         time.Now(),
				LastTimeReservesUpdated: time.Unix(0, 0),
			}

			_, err := pair.Insert(temp)

			if err != nil {
				panic(err)
			}

			// var wg sync.WaitGroup
		
			// detectAndUpdate := func(address common.Address) {
			// 	defer wg.Done()
			// 	taxDetector, err := taxTokenDetector.NewTaxDetector()
			// 	if err != nil {
			// 		panic(err)
			// 	}

			// 	isTaxToken, err := taxDetector.DetectTaxToken(context.Background(), address)
			// 	if err != nil {
			// 		panic(err)
			// 	}
			
			// 	erc20.UpdateTax(address, isTaxToken, -1)

				
			// 	// This query is slow and could cause a deadlock if not done in a mutex
			// 	pair.UpdateAllPairsWithTokenToAPairWithATaxToken(address, isTaxToken)
		
			// }

			// wg.Add(3)
			// go detectAndUpdate(token0Address)
			// go detectAndUpdate(token1Address)
			// go func() {
			// 	defer wg.Done()
			// 	UpdateReservesForPair(modelPair, 10)
			// }()

			// wg.Wait()

		} else {
			if errToken0 != nil {
				log.Warn().Msgf("Could not get Erc20 with Address:" + token0Address.String())
			}
			if errToken1 != nil {
				log.Warn().Msgf("Could not get Erc20 with Address:" + token1Address.String())
			}
		}

		if counter % 10 == 0 {
			log.Info().Msgf("Processed %d pairs of %d for dex_id %d", counter, currentLength, dex.DexId)
		}

		counter++
	}

	log.Info().Msgf("Completed gathering pairs for " + dex.Name)
}



func UpdateReservesForPairWithAddress(pairAddress common.Address, retries int) (reserveToken0 *big.Int, reserveToken1 *big.Int, err error) {
	modelPair, err := pair.GetByContractAddress(pairAddress)

	if err != nil {
		return nil, nil, err
	}

	reserveToken0, reserveToken1, err = UpdateReservesForPair(modelPair, retries)

	if err != nil {
		return nil, nil, err
	}

	return reserveToken0, reserveToken1, nil
}

func UpdateReservesForPair(modelPair *types.ModelPair, retries int) (reserveToken0 *big.Int, reserveToken1 *big.Int, err error) {
	var token0Reserve *big.Int
	var token1Reserve *big.Int

	if types.IsV2Pair(modelPair) {
		token0Reserve, token1Reserve, err = pairServiceV2.GetReservesForPair(
			modelPair.PairContractAddress, 
			modelPair.Token0Erc20.ContractAddress, 
			modelPair.Token1Erc20.ContractAddress, 
			retries,
		)
	} else if types.IsV3Pair(modelPair) {
		token0Reserve, token1Reserve, err = pairServiceV3.GetReservesForPair(
			modelPair.PairContractAddress, 
			modelPair.Token0Erc20.ContractAddress, 
			modelPair.Token1Erc20.ContractAddress, 
			retries,
		)
	} else {
		panic("Error updating reserves for pair " + modelPair.PairContractAddress.String() + " because it is not a v2 or v3 pair")
	}

	if err == nil {
		modelPair.Token0Reserves = *token0Reserve
		modelPair.Token1Reserves = *token1Reserve

		err = pair.UpdateReserves(modelPair.PairId, modelPair.Token0Reserves, modelPair.Token1Reserves)

		if err != nil {
			log.Error().Msgf("Error updating reserves for pair %s: %v", modelPair.PairContractAddress.String(), err)
			return nil, nil, err
		}
	} else {
		log.Error().Msgf("Error updating reserves for pair %s: %v", modelPair.PairContractAddress.String(), err)
		return nil, nil, err
	}

	return &modelPair.Token0Reserves, &modelPair.Token1Reserves, nil
}

func FindPairsInV1ThatAreNotInV2() {
	dexV1 := dex.GetById(3)
	dexV2 := dex.GetById(4)
	amountOfPls := myConst.GetOneHundredMillionWplsBigint()

	dexV1Pairs := PlsPairWithHighAmountOfPls([]int{dexV1.DexId}, amountOfPls, false)
	dexV2Pairs, _ := pair.GetAllPairsThatHavePlsByDexId(dexV2.DexId, false)

	for _, dexV1Pair := range dexV1Pairs {

		if dexV1Pair.HasTaxToken {
			continue
		}

		found := false
		for _, dexV2Pair := range dexV2Pairs {

			if dexV2Pair.HasTaxToken {
				continue
			}

			if (dexV1Pair.Token0Address == dexV2Pair.Token0Address && 
			   dexV1Pair.Token1Address == dexV2Pair.Token1Address) ||
			   (dexV1Pair.Token0Address == dexV2Pair.Token1Address && 
			   dexV1Pair.Token1Address == dexV2Pair.Token0Address) {

				found = true
				break
			}

			
		}
		if !found {
			UpdateReservesForPair(&dexV1Pair, 5)
			updatedPair, _ :=pair.GetById(dexV1Pair.PairId)

			hasHighAmountOfPls := pair.PlsPairHasHighAmountOfPls(*updatedPair, amountOfPls)

			if hasHighAmountOfPls {
				log.Info().Msgf("Pair %s has high amount of pls", dexV1Pair.PairContractAddress.String())
			}

			
		}
	}

	log.Info().Msgf("Completed finding pairs in v1 that are not in v2")
	
}

