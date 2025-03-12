package dexpairgather

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/pborgen/liquidityFinder/internal/database/model/dex"
	"github.com/pborgen/liquidityFinder/internal/database/model/pair"
	"github.com/pborgen/liquidityFinder/internal/mylogger"
	"github.com/pborgen/liquidityFinder/internal/service/pairService"
	"golang.org/x/sync/semaphore"

	"github.com/pborgen/liquidityFinder/internal/types"
	"github.com/rs/zerolog/log"
)

type PossibleSandwich struct {
    DexName  string
	DexRouterAddress string
	DexFactoryAddress string

	Token0Address string
	Token1Address string
	PairContractAddress string
}

func Start() {
	mylogger.Init()

	allDexes := dex.GetAllByNetworkId(1)

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(5) // Limit to 5 concurrent goroutines

	

	for _, modelDex := range allDexes {
		wg.Add(1)

		// Acquire a semaphore slot
		if err := sem.Acquire(context.Background(), 1); err != nil {
			log.Error().Msgf("Failed to acquire semaphore: %v", err)
			wg.Done()
			continue
		}
		
		go func(modelDex dex.ModelDex) {
			defer wg.Done()
			defer sem.Release(1) // Release the semaphore slot when done

			Gather(modelDex.DexId)
		}(modelDex)
	
	}

	wg.Wait()
	log.Info().Msgf("Completed gathering new pairs for all dexes")
}



func WriteToFilePlsPairsByDexId(dexIds []int) {

	millionPls := big.NewInt(0)
    millionPls.SetString("1000000000000000000000000", 10)
	tenMillionPls := new(big.Int).Mul(millionPls, big.NewInt(10))


	var allPairs []types.ModelPair = pair.PlsPairWithHighAmountOfPls(dexIds, tenMillionPls, false)

	log.Info().Msgf("Found " + strconv.Itoa(len(allPairs)) + " pairs with high amount of pls")

	var myMap = make(map[string]PossibleSandwich)
	for _, pair := range allPairs {
		dex := dex.GetById(pair.DexId)
		myMap[pair.PairContractAddress.String()] = PossibleSandwich{
			DexName: dex.Name,
			DexRouterAddress: dex.RouterContractAddress.String(),
			DexFactoryAddress: dex.FactoryContractAddress.String(),

			Token0Address: pair.Token0Erc20.ContractAddress.String(),
			Token1Address: pair.Token1Erc20.ContractAddress.String(),
			PairContractAddress: pair.PairContractAddress.String(),
		}
	}

	// Convert map to json
	jsonData, err := json.MarshalIndent(myMap, "", "  ")
	if err != nil {
		log.Error().Msgf(err.Error())
	}


	// Write to file
	err = os.WriteFile("possibleSandwiches.json", jsonData, 0644)
	if err != nil {
		log.Error().Msgf(err.Error())
	}

	log.Info().Msgf(string(jsonData))
}

func Gather(dexId int) {
	log.Info().Msgf("Starting dexpairgather.Start with Id:" + strconv.Itoa(dexId))
	dex := dex.GetById(dexId)

	log.Info().Msgf("Gathering new pairs for dex with name:" + dex.Name)

	
	maxRetries := 5 // Define the maximum number of retries
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Warn().Msgf("Attempt %d: Recovered from panic: %v", attempt, r)
				}
			}()

			log.Info().Msgf("Gathering new pairs for dex with name:" + dex.Name)

			pairService.PopulatePairsInDb(dex)
		}()

		// If the operation panicked, wait before retrying
		if attempt < maxRetries {
			log.Info().Msgf("Retrying operation after panic, attempt %d", attempt+1)
			time.Sleep(10 * time.Second) // Adjust the sleep duration as needed
		}
	}

	
	
}
