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

// StoryHandler handles HTTP requests for product story ads.
type StoryHandler struct {
	storyModule module.StoryModule
}

// NewStoryHandler creates a new StoryHandler.
func NewStoryHandler(sm module.StoryModule) *StoryHandler {
	return &StoryHandler{storyModule: sm}
}

// CreateStory godoc
// @Summary Create a product story ad
// @Description Store owner creates a new story ad (image/video) for a specific product within a date range
// @Tags Story
// @Accept json
// @Produce json
// @Param request body dto.CreateProductStoryRequest true "Story ad data"
// @Success 201 {object} response.BaseResponse{data=dto.ProductStory}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 422 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/stories [post]
func (h *StoryHandler) CreateStory(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		c.Error(errorx.New(errorx.ErrUnauthorized, "store context missing", http.StatusUnauthorized))
		return
	}

	role := c.GetString("role")
	if role != "admin" {
		c.Error(errorx.New(errorx.ErrForbidden, "only store admins can create story ads", http.StatusForbidden))
		return
	}

	var req dto.CreateProductStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}

	story, err := h.storyModule.CreateStory(c.Request.Context(), storeID, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, story)
}

// ListMyStories godoc
// @Summary List my store's story ads
// @Description Returns paginated story ads for the authenticated merchant's store
// @Tags Story
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param product_id query int false "Filter by product ID"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/stories [get]
func (h *StoryHandler) ListMyStories(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		c.Error(errorx.New(errorx.ErrUnauthorized, "store context missing", http.StatusUnauthorized))
		return
	}

	var filter dto.ProductStoryFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}

	// Inject store ID server-side — never trust client
	filter.StoreID = storeID

	resp, err := h.storyModule.ListMyStories(c.Request.Context(), filter)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// GetMyStory godoc
// @Summary Get a single story ad (merchant)
// @Description Returns one story ad belonging to the authenticated merchant's store
// @Tags Story
// @Produce json
// @Param id path int true "Story ID"
// @Success 200 {object} response.BaseResponse{data=dto.ProductStory}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 404 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/stories/:id [get]
func (h *StoryHandler) GetMyStory(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		c.Error(errorx.New(errorx.ErrUnauthorized, "store context missing", http.StatusUnauthorized))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, "invalid story ID", http.StatusBadRequest))
		return
	}

	story, err := h.storyModule.GetStory(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	// Ownership check
	if story.StoreID != storeID {
		c.Error(errorx.New(errorx.ErrForbidden, "story does not belong to your store", http.StatusForbidden))
		return
	}

	response.Success(c, http.StatusOK, story)
}

// UpdateStory godoc
// @Summary Update a story ad
// @Description Partially update an existing story ad (admin only, store-scoped)
// @Tags Story
// @Accept json
// @Produce json
// @Param id path int true "Story ID"
// @Param request body dto.UpdateProductStoryRequest true "Update story data"
// @Success 200 {object} response.BaseResponse{data=dto.ProductStory}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 422 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/stories/:id [put]
func (h *StoryHandler) UpdateStory(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		c.Error(errorx.New(errorx.ErrUnauthorized, "store context missing", http.StatusUnauthorized))
		return
	}

	role := c.GetString("role")
	if role != "admin" {
		c.Error(errorx.New(errorx.ErrForbidden, "only store admins can update story ads", http.StatusForbidden))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, "invalid story ID", http.StatusBadRequest))
		return
	}

	var req dto.UpdateProductStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}

	story, err := h.storyModule.UpdateStory(c.Request.Context(), storeID, id, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, story)
}

// DeleteStory godoc
// @Summary Delete a story ad
// @Description Soft-delete a story ad (admin only, store-scoped)
// @Tags Story
// @Produce json
// @Param id path int true "Story ID"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/stories/:id [delete]
func (h *StoryHandler) DeleteStory(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		c.Error(errorx.New(errorx.ErrUnauthorized, "store context missing", http.StatusUnauthorized))
		return
	}

	role := c.GetString("role")
	if role != "admin" {
		c.Error(errorx.New(errorx.ErrForbidden, "only store admins can delete story ads", http.StatusForbidden))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, "invalid story ID", http.StatusBadRequest))
		return
	}

	if err := h.storyModule.DeleteStory(c.Request.Context(), storeID, id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "story deleted"})
}

// PublicListActiveStories godoc
// @Summary List active story ads (public)
// @Description Returns all currently active story ads ordered by recency (date-range enforced)
// @Tags Story
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /stories [get]
func (h *StoryHandler) PublicListActiveStories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := dto.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := h.storyModule.ListActiveStories(c.Request.Context(), params)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// PublicGetStory godoc
// @Summary Get a single story ad (public)
// @Description Returns one story with full product details; increments view count
// @Tags Story
// @Produce json
// @Param id path int true "Story ID"
// @Success 200 {object} response.BaseResponse{data=dto.ProductStory}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 404 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /stories/:id [get]
func (h *StoryHandler) PublicGetStory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, "invalid story ID", http.StatusBadRequest))
		return
	}

	story, err := h.storyModule.GetStory(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, story)
}
