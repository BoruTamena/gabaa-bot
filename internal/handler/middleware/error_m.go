package middleware

import (
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/gin-gonic/gin"
)

// ErrorMiddleware catches errors added to the Gin context, 
// and formats them into a standard HTTP JSON response.
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Execute all handlers in the chain

		// Check if there are any errors attached to the request
		if len(c.Errors) > 0 {
			// Get the most recent error
			err := c.Errors.Last().Err
			
			// Format it using errorx, which correctly maps standard HTTP response
			status, appErr := errorx.ErrorResponse(err)
			
			// Send the JSON response
			c.JSON(status, gin.H{
				"success": false,
				"data":    nil,
				"error":   appErr,
			})
			c.Abort()
		}
	}
}
