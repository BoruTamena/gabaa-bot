package product

import (
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryModule module.CategoryModule
}

func NewCategoryHandler(cModule module.CategoryModule) *CategoryHandler {
	return &CategoryHandler{categoryModule: cModule}
}

// ListAllCategories returns all categories
// @Summary List all categories
// @Description Retrieve a paginated list of all global categories
// @Tags Category
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /categories [get]
func (h *CategoryHandler) ListAllCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := dto.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := h.categoryModule.ListAllCategories(c.Request.Context(), params)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// ListStoreCategories returns categories for a specific store
// @Summary List store categories
// @Description Retrieve all categories specific to a store
// @Tags Category
// @Produce json
// @Param store_id path int true "Store ID"
// @Success 200 {object} response.BaseResponse{data=[]dto.Category}
// @Failure 404 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/:store_id/categories [get]
func (h *CategoryHandler) ListStoreCategories(c *gin.Context) {
	storeID, _ := strconv.ParseInt(c.Param("store_id"), 10, 64)

	categories, err := h.categoryModule.ListStoreCategories(c.Request.Context(), storeID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, categories)
}

// CreateCategory adds a new category for a store
// @Summary Add store category
// @Description Creates a new custom category for a store
// @Tags Category
// @Accept json
// @Produce json
// @Param store_id path int true "Store ID"
// @Param request body dto.CreateCategoryRequest true "Category Details"
// @Success 201 {object} response.BaseResponse{data=dto.Category}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 422 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/:store_id/category [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	storeID, _ := strconv.ParseInt(c.Param("store_id"), 10, 64)

	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Replace inline errors with c.Error to be caught by ErrorMiddleware or explicit
		// Let's use standard formatting. But wait, if ErrorMiddleware catches c.Error, we can just return it.
		// Wait, c.Error only works well if we set the correct errorx.AppError.
		// For simplicity, let's use c.Error but create AppError first. Wait, product.go and category.go do c.Error.
		// I'll just keep it consistent with what I did in other places: response.Error.
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    nil,
			"error":   gin.H{"code": "BAD_REQUEST", "message": err.Error()},
		})
		return
	}

	category, err := h.categoryModule.CreateCategory(c.Request.Context(), storeID, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, category)
}
