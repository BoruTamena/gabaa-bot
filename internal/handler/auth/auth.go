package auth

import (
	"net/http"

	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
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
// @Accept json
// @Produce json
// @Param body body map[string]string true "initData"
// @Router /auth/telegram [post]
func (h *AuthHandler) TelegramAuth(c *gin.Context) {
	var req struct {
		InitData string `json:"initData" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, "Missing or invalid request body", http.StatusBadRequest)
		c.JSON(appErr.Status, appErr)
		return
	}

	resp, err := h.authModule.TelegramAuth(c.Request.Context(), req.InitData)
	if err != nil {
		status, appErr := errorx.ErrorResponse(err)
		if appErr.Code == errorx.ErrInternal {
			appErr = errorx.New(errorx.ErrUnauthorized, err.Error(), http.StatusUnauthorized)
			status = http.StatusUnauthorized
		}
		c.JSON(status, appErr)
		return
	}

	c.JSON(http.StatusOK, resp)
}
