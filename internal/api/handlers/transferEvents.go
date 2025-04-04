package handlers

import (
	"errors"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pborgen/liquidityFinder/internal/api/models"
	"github.com/pborgen/liquidityFinder/internal/service/transferEventService"
	"github.com/rs/zerolog/log"
)

// GetPairs handles the GET /pairs endpoint
func GetTransferEventsForAddress(c *gin.Context) {
	// Get query parameters
	address := c.Param("address")
	viewMode := c.Query("viewMode")

	if viewMode == "" {
		viewMode = "all"
	}

	if viewMode != "all" && viewMode != "in" && viewMode != "out" {
		HandleGeneralError(c, "Invalid view mode", errors.New("invalid view mode"))
		return
	}

	transferEvents, err := 
		transferEventService.GetAllForAddressGroupBy(common.HexToAddress(address), viewMode)

	if err != nil {
		log.Error().Err(err).Msg("Failed to get pairs")
		HandleGeneralError(c, "Failed to get pairs", err)
		return
	}

	// Return paginated response
	c.JSON(http.StatusOK, models.TransferEventsGroupByResponse{
		Success: true,
		Data:    transferEvents,
	})
}

