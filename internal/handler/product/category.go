package product

import (
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
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
// @Tags product
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Produce json
// @Router /categories [get]
func (h *CategoryHandler) ListAllCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := dto.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	response, err := h.categoryModule.ListAllCategories(c.Request.Context(), params)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListStoreCategories returns categories for a specific store
// @Summary List store categories
// @Tags product
// @Param store_id path int true "Store ID"
// @Produce json
// @Router /store/:store_id/categories [get]
func (h *CategoryHandler) ListStoreCategories(c *gin.Context) {
	storeID, _ := strconv.ParseInt(c.Param("store_id"), 10, 64)

	categories, err := h.categoryModule.ListStoreCategories(c.Request.Context(), storeID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, categories)
}

// CreateCategory adds a new category for a store
// @Summary Add store category
// @Tags product
// @Param store_id path int true "Store ID"
// @Accept json
// @Produce json
// @Router /store/:store_id/category [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	storeID, _ := strconv.ParseInt(c.Param("store_id"), 10, 64)

	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.categoryModule.CreateCategory(c.Request.Context(), storeID, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, category)
}
