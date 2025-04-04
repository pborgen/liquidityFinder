package api

import (
	"github.com/pborgen/liquidityFinder/internal/api/router"
	"github.com/pborgen/liquidityFinder/internal/myConfig"
	"github.com/pborgen/liquidityFinder/internal/mylogger"
	"github.com/rs/zerolog/log"
)

func Start() {
	// Initialize logger
	mylogger.Init()
	log.Info().Msg("Starting API server...")

	// Load configuration
	config := myConfig.GetInstance()
	if config == nil {
		log.Fatal().Msg("Failed to load configuration")
	}

	// Initialize and start the router
	r := router.SetupRouter()
	
	// Start server
	port := ":3001"
	log.Info().Msgf("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
} 