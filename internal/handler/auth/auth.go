package auth

import (
	"net/http"

	authmod "github.com/BoruTamena/gabaa-bot/internal/module/auth"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
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

// AdminLogin authenticates a platform admin with static credentials
// @Summary Platform admin login
// @Description Validates static admin credentials and returns a JWT with role platform_admin
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.AdminLoginRequest true "Admin credentials"
// @Success 200 {object} response.BaseResponse{data=dto.AuthResponse}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Router /auth/admin/login [post]
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req dto.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		response.CustomError(c, appErr)
		return
	}

	resp, err := h.authModule.AdminLogin(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		appErr := errorx.New(errorx.ErrUnauthorized, err.Error(), http.StatusUnauthorized)
		response.CustomError(c, appErr)
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// StartTelegramLoginSession starts a bot-mediated login session for web/mobile clients.
func (h *AuthHandler) StartTelegramLoginSession(c *gin.Context) {
	resp, err := h.authModule.StartBotLoginSession(c.Request.Context())
	if err != nil {
		appErr := errorx.New(errorx.ErrInternal, err.Error(), http.StatusInternalServerError)
		response.CustomError(c, appErr)
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// PollTelegramLoginSession polls a bot login session and returns a JWT when completed.
func (h *AuthHandler) PollTelegramLoginSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		appErr := errorx.New(errorx.ErrBadRequest, "Missing session id", http.StatusBadRequest)
		response.CustomError(c, appErr)
		return
	}

	resp, err := h.authModule.PollBotLoginSession(c.Request.Context(), sessionID)
	if err != nil {
		if authmod.IsSessionNotFound(err) {
			appErr := errorx.New(errorx.ErrNotFound, "Session not found or expired", http.StatusNotFound)
			response.CustomError(c, appErr)
			return
		}
		appErr := errorx.New(errorx.ErrInternal, err.Error(), http.StatusInternalServerError)
		response.CustomError(c, appErr)
		return
	}

	response.Success(c, http.StatusOK, resp)
}
