package product

import (
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productModule module.ProductModule
}

func NewProductHandler(pModule module.ProductModule) *ProductHandler {
	return &ProductHandler{productModule: pModule}
}

// ListProducts returns all products for a store with pagination
// @Summary List products
// @Tags product
// @Param store_id path int true "Store ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Produce json
// @Router /store/:store_id/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, _ := strconv.ParseInt(storeIDStr, 10, 64)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := dto.PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}

	response, err := h.productModule.ListProducts(c.Request.Context(), storeID, params)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateProduct adds a new product to a store
// @Summary Add product
// @Tags product
// @Accept json
// @Produce json
// @Router /store/:store_id/product [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, _ := strconv.ParseInt(storeIDStr, 10, 64)

	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		c.Error(appErr)
		return
	}

	role := c.GetString("role")
	userStoreID := c.GetInt64("store_id")

	if role != "admin" || (userStoreID != 0 && userStoreID != storeID) {
		appErr := errorx.New(errorx.ErrForbidden, "Unauthorized to add products to this store", http.StatusForbidden)
		c.Error(appErr)
		return
	}

	product, err := h.productModule.CreateProduct(c.Request.Context(), storeID, req)
	if err != nil {
		appErr := errorx.New(errorx.ErrValidation, err.Error(), http.StatusUnprocessableEntity)
		c.Error(appErr)
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct edits an existing product
// @Summary Edit product
// @Tags product
// @Router /store/:store_id/product/:id [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, _ := strconv.ParseInt(storeIDStr, 10, 64)
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		c.Error(appErr)
		return
	}

	role := c.GetString("role")
	userStoreID := c.GetInt64("store_id")

	if role != "admin" || (userStoreID != 0 && userStoreID != storeID) {
		appErr := errorx.New(errorx.ErrForbidden, "Unauthorized to update products for this store", http.StatusForbidden)
		c.Error(appErr)
		return
	}

	product, err := h.productModule.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		appErr := errorx.New(errorx.ErrValidation, err.Error(), http.StatusUnprocessableEntity)
		c.Error(appErr)
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct removes a product
// @Summary Delete product
// @Tags product
// @Router /store/:store_id/product/:id [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	storeIDStr := c.Param("store_id")
	storeID, _ := strconv.ParseInt(storeIDStr, 10, 64)
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	role := c.GetString("role")
	userStoreID := c.GetInt64("store_id")

	if role != "admin" || (userStoreID != 0 && userStoreID != storeID) {
		appErr := errorx.New(errorx.ErrForbidden, "Unauthorized to delete products from this store", http.StatusForbidden)
		c.Error(appErr)
		return
	}

	err := h.productModule.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}
