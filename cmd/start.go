package main

import (
	"os"
	"time"

	"github.com/pborgen/liquidityFinder/cmd/dexpairgather"
	"github.com/pborgen/liquidityFinder/internal/api/router"

	"github.com/pborgen/liquidityFinder/internal/myConfig"
	"github.com/pborgen/liquidityFinder/internal/mylogger"
	"github.com/pborgen/liquidityFinder/internal/service/pairService"
	"github.com/pborgen/liquidityFinder/internal/service/taxTokenDetector"

	"github.com/rs/zerolog/log"

	"github.com/pborgen/liquidityFinder/internal/service/pairServiceV3"
)

func main() {
	pulsechainNetworkId := 1

	mylogger.Init()

	baseDir := os.Getenv("BASE_DIR")
	log.Info().Msgf("BASE_DIR: " + baseDir)
	log.Info().Msgf("Starting...")

	log.Info().Msgf("-----------------------------")
	log.Info().Msgf("POSTGRES_HOST: " + os.Getenv("POSTGRES_HOST"))
	log.Info().Msgf("-----------------------------")

	args := os.Args
	log.Info().Msgf("-----------------------------")
	log.Info().Msgf("Type of Args = %T\n", args)

	if len(args) < 2 {
		panic("Invalid params passed")
	}
	log.Info().Msgf(args[0], args[1])
	log.Info().Msgf("-----------------------------")

	processName := args[1]
	log.Info().Msgf("About to run process with Name: " + processName)

	myConfig.GetInstance()
	
	if processName == "api" {
		// Initialize and start the API server
		r := router.SetupRouter()
		port := ":8080"
		log.Info().Msgf("API server starting on port %s", port)
		if err := r.Run(port); err != nil {
			log.Fatal().Err(err).Msg("Failed to start API server")
		}
	} else if processName == "test" {
		
	} else if processName == "gatherPairs" {
		for {	
			dexpairgather.Start()
			time.Sleep(15 * time.Second)
		}
	} else if processName == "updateReservesAndHighLiquidity" {
		for {
			pairService.UpdateIsHighLiquidityForAllPairs(false)
			time.Sleep(10 * time.Hour)
		}
	} else if processName == "addNewPairsForAllV3Dexes" {
		pairServiceV3.AddNewPairForAllV3Dexes(pulsechainNetworkId)

	} else if processName == "writePlsPairsByDexId" {
		dexpairgather.WriteToFilePlsPairsByDexId([]int{3, 4})
	} else if processName == "fixPairOrdering" {
		pairService.FixPairOrdering()
	
	} else if processName == "detectTaxToken" {
		for {
			detector, err := taxTokenDetector.NewTaxDetector()
			if err != nil {
				log.Error().Msgf("Error in NewTaxDetector: " + err.Error())
			}
			detector.Run()
			time.Sleep(30 * time.Second)
		}
	} else {
		log.Error().Msgf("Invalid process Name: " + processName)
	}
	log.Info().Msgf("Completed...")

}
