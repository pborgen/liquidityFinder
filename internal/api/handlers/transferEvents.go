package handlers

import (
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

	transferEvents, err := transferEventService.GetAllForAddressGroupBy(common.HexToAddress(address))

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

