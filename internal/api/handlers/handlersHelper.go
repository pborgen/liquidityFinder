package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pborgen/liquidityFinder/internal/api/models"
)


func HandleErrorLimitToHighError(c *gin.Context) {
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
}

func HandleGeneralError(c *gin.Context, message string, err error) {
	c.JSON(http.StatusInternalServerError, models.ErrorResponse{
		Success: false,
		Error: struct {
			Code    string      `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details,omitempty"`
		}{
			Code:    "INTERNAL_ERROR",
			Message: message,
			Details: err.Error(),
		},
	})
}