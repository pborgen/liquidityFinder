package middleware

import (

	"github.com/gin-gonic/gin"
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement proper API key validation if we end up wanting to use it
		// apiKey := c.GetHeader("X-API-Key")
		// if apiKey == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"success": false,
		// 		"error": gin.H{
		// 			"code":    "UNAUTHORIZED",
		// 			"message": "Missing API key",
		// 		},
		// 	})
		// 	c.Abort()
		// 	return
		// }

		// TODO: Implement proper API key validation
		// For now, accept any non-empty key
		c.Next()
	}
} 