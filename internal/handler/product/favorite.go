package product

import (
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	favoriteModule module.FavoriteModule
}

func NewFavoriteHandler(fm module.FavoriteModule) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteModule: fm,
	}
}

// AddFavorite godoc
// @Summary Add a product to favorites
// @Description Adds a specific product to the authenticated user's favorites list
// @Tags Favorite
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 201 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /favorites/{product_id} [post]
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.Error(errorx.New(errorx.ErrUnauthorized, "unauthorized access", http.StatusUnauthorized))
		return
	}

	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, "invalid product ID", http.StatusBadRequest))
		return
	}

	if err := h.favoriteModule.AddFavorite(c.Request.Context(), userID, productID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, gin.H{"message": "product added to favorites"})
}

// RemoveFavorite godoc
// @Summary Remove a product from favorites
// @Description Removes a specific product from the authenticated user's favorites list
// @Tags Favorite
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /favorites/{product_id} [delete]
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.Error(errorx.New(errorx.ErrUnauthorized, "unauthorized access", http.StatusUnauthorized))
		return
	}

	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, "invalid product ID", http.StatusBadRequest))
		return
	}

	if err := h.favoriteModule.RemoveFavorite(c.Request.Context(), userID, productID); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "product removed from favorites"})
}

// ListUserFavorites godoc
// @Summary List user's favorite products
// @Description Returns a paginated list of the authenticated user's favorite products
// @Tags Favorite
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /favorites [get]
func (h *FavoriteHandler) ListUserFavorites(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.Error(errorx.New(errorx.ErrUnauthorized, "unauthorized access", http.StatusUnauthorized))
		return
	}

	var params dto.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}

	resp, err := h.favoriteModule.ListUserFavorites(c.Request.Context(), userID, params)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, resp)
}
