package auth

import (
	"net/http"

	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authModule module.AuthModule
}

func NewAuthHandler(aModule module.AuthModule) *AuthHandler {
	return &AuthHandler{authModule: aModule}
}

// TelegramAuth handles Telegram MiniApp authentication
// @Summary Authenticate via Telegram
// @Description Validates initData and returns JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "initData"
// @Success 200 {object} response.BaseResponse{data=dto.AuthResponse}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /auth/telegram [post]
func (h *AuthHandler) TelegramAuth(c *gin.Context) {
	var req struct {
		InitData string `json:"initData" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, "Missing or invalid request body", http.StatusBadRequest)
		response.CustomError(c, appErr)
		return
	}

	resp, err := h.authModule.TelegramAuth(c.Request.Context(), req.InitData)
	if err != nil {
		appErr, ok := err.(*errorx.AppError)
		if !ok || appErr.Code == errorx.ErrInternal {
			appErr = errorx.New(errorx.ErrUnauthorized, err.Error(), http.StatusUnauthorized)
		}
		response.CustomError(c, appErr)
		return
	}

	response.Success(c, http.StatusOK, resp)
}
