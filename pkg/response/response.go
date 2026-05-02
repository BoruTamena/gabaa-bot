package response

import (
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/gin-gonic/gin"
)

// BaseResponse represents the standard response structure
type BaseResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

// Success returns a standardized success response
func Success(c *gin.Context, status int, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(status, BaseResponse{
		Success: true,
		Data:    data,
		Error:   nil,
	})
}

// Error returns a standardized error response based on an error interface
func Error(c *gin.Context, err error) {
	status, appErr := errorx.ErrorResponse(err)
	c.JSON(status, BaseResponse{
		Success: false,
		Data:    nil,
		Error:   appErr,
	})
}

// CustomError returns a standardized error response based on an AppError
func CustomError(c *gin.Context, appErr *errorx.AppError) {
	c.JSON(appErr.Status, BaseResponse{
		Success: false,
		Data:    nil,
		Error:   appErr,
	})
}
