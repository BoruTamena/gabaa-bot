package preference

import (
	"net/http"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

type PreferenceHandler struct {
	recommendationModule module.RecommendationModule
}

func NewPreferenceHandler(rm module.RecommendationModule) *PreferenceHandler {
	return &PreferenceHandler{recommendationModule: rm}
}

// GetPreferences godoc
// @Summary Get user category preferences
// @Description Returns the authenticated user's recommendation opt-in status and preferred product categories
// @Tags Preference
// @Produce json
// @Success 200 {object} response.BaseResponse{data=dto.UserPreferences}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/preferences [get]
func (h *PreferenceHandler) GetPreferences(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.Error(errorx.New(errorx.ErrUnauthorized, "unauthorized access", http.StatusUnauthorized))
		return
	}

	prefs, err := h.recommendationModule.GetPreferences(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, prefs)
}

// UpdatePreferences godoc
// @Summary Update user category preferences
// @Description Replaces the authenticated user's preferred categories and recommendation opt-in status
// @Tags Preference
// @Accept json
// @Produce json
// @Param request body dto.UpdateUserPreferencesRequest true "Preference settings"
// @Success 200 {object} response.BaseResponse{data=dto.UserPreferences}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /user/preferences [put]
func (h *PreferenceHandler) UpdatePreferences(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.Error(errorx.New(errorx.ErrUnauthorized, "unauthorized access", http.StatusUnauthorized))
		return
	}

	var req dto.UpdateUserPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}

	prefs, err := h.recommendationModule.SetPreferences(c.Request.Context(), userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, prefs)
}
