package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pborgen/liquidityFinder/internal/api/handlers"
	"github.com/pborgen/liquidityFinder/internal/api/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API key middleware
	r.Use(middleware.APIKeyAuth())

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Pairs endpoints
		pairs := v1.Group("/pairs")
		{
			pairs.GET("", handlers.GetPairs)
		}




	}


	return r
} 