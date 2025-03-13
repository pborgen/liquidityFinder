package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pborgen/liquidityFinder/internal/api/models"
	"github.com/pborgen/liquidityFinder/internal/database/model/pair"
	"github.com/pborgen/liquidityFinder/internal/types"
	"github.com/rs/zerolog/log"
)

// GetPairs handles the GET /pairs endpoint
func GetPairs(c *gin.Context) {
	// Get query parameters
	dexId, _ := strconv.Atoi(c.Query("dexId"))
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	// Convert string parameters to appropriate types
	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)

	if limitNum > 500 {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error: struct {
				Code    string      `json:"code"`
				Message string      `json:"message"`
				Details interface{} `json:"details,omitempty"`
			}{
				Code:    "INTERNAL_ERROR",
				Message: "Limit is too high",
				Details: errors.New("limit is too high"),
			},
		})
		return
	}

	// Get pairs from database
	var pairs []types.ModelPair
	var err error

	if dexId == 0 {
		pairs, err = pair.GetAllPageAndLimit(pageNum, limitNum, true)
	} else {
		pairs, err = pair.GetAllWithDexIdPageAndLimit(dexId, pageNum, limitNum, true)
	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to get pairs")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error: struct {
				Code    string      `json:"code"`
				Message string      `json:"message"`
				Details interface{} `json:"details,omitempty"`
			}{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to get pairs",
				Details: err.Error(),
			},
		})
		return
	}


	// Calculate pagination
	total := len(pairs)
	totalPages := (total + limitNum - 1) / limitNum
	start := (pageNum - 1) * limitNum
	end := start + limitNum
	if end > total {
		end = total
	}

	// Return paginated response
	c.JSON(http.StatusOK, models.PairsResponse{
		Success: true,
		Data:    pairs[start:end],
		Pagination: models.Pagination{
			CurrentPage:  pageNum,
			TotalPages:   totalPages,
			TotalResults: total,
		},
	})
}

