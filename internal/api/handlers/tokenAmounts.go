package handlers

import (
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pborgen/liquidityFinder/internal/api/models"
	"github.com/pborgen/liquidityFinder/internal/service/tokenAmountService"
	"github.com/pborgen/liquidityFinder/internal/types"
	"github.com/rs/zerolog/log"
)

// GetPairs handles the GET /pairs endpoint
func GetTokenAmountsForTokenAddress(c *gin.Context) {
	// Get query parameters
	tokenAddress := c.Param("tokenAddress")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	// Convert string parameters to appropriate types
	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)

	if limitNum > 500 {
		HandleErrorLimitToHighError(c)
		return
	}

	// Get pairs from database
	var tokenAmounts []types.ModelTokenAmount
	var err error

	tokenAmounts, err = tokenAmountService.GetByTokenAddress(common.HexToAddress(tokenAddress), limitNum, pageNum)


	if err != nil {
		log.Error().Err(err).Msg("Failed to get pairs")
		HandleGeneralError(c, "Failed to get pairs", err)
		return
	}


	// Calculate pagination
	total := len(tokenAmounts)
	totalPages := (total + limitNum - 1) / limitNum
	start := (pageNum - 1) * limitNum
	end := start + limitNum
	if end > total {
		end = total
	}

	// Return paginated response
	c.JSON(http.StatusOK, models.TokenAmountsResponse{
		Success: true,
		Data:    tokenAmounts,
		Pagination: models.Pagination{
			CurrentPage:  pageNum,
			TotalPages:   totalPages,
			TotalResults: total,
		},
	})
}

