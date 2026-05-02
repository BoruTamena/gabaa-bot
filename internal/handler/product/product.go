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

type ProductHandler struct {
	productModule module.ProductModule
}

func NewProductHandler(pModule module.ProductModule) *ProductHandler {
	return &ProductHandler{productModule: pModule}
}

// ListProducts returns all products for a store with pagination
// @Summary List products
// @Description Retrieve a paginated list of all products for a store
// @Tags Product
// @Produce json
// @Param store_id path int true "Store ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
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

	resp, err := h.productModule.ListProducts(c.Request.Context(), storeID, params)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// CreateProduct adds a new product to a store
// @Summary Add product
// @Description Add a new product to a store (Admin only)
// @Tags Product
// @Accept json
// @Produce json
// @Param store_id path int true "Store ID"
// @Param request body dto.CreateProductRequest true "Product details"
// @Success 201 {object} response.BaseResponse{data=dto.Product}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 422 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
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

	userID := c.GetInt64("user_id")
	product, err := h.productModule.CreateProduct(c.Request.Context(), userID, storeID, req)
	if err != nil {
		appErr := errorx.New(errorx.ErrValidation, err.Error(), http.StatusUnprocessableEntity)
		c.Error(appErr)
		return
	}

	response.Success(c, http.StatusCreated, product)
}

// UpdateProduct edits an existing product
// @Summary Edit product
// @Description Modify details of an existing product (Admin only)
// @Tags Product
// @Accept json
// @Produce json
// @Param store_id path int true "Store ID"
// @Param id path int true "Product ID"
// @Param request body dto.UpdateProductRequest true "Product update details"
// @Success 200 {object} response.BaseResponse{data=dto.Product}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 422 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
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

	response.Success(c, http.StatusOK, product)
}

// DeleteProduct removes a product
// @Summary Delete product
// @Description Completely remove a product from a store (Admin only)
// @Tags Product
// @Produce json
// @Param store_id path int true "Store ID"
// @Param id path int true "Product ID"
// @Success 200 {object} response.BaseResponse{data=map[string]string}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
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

	response.Success(c, http.StatusOK, gin.H{"message": "product deleted"})
}

// PublicGetProductByID returns a single product by its ID (Public)
// @Summary Get a product by ID (Public)
// @Description Fetch details of a single product (public endpoint)
// @Tags Product
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.BaseResponse{data=dto.Product}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 404 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /products/{id} [get]
func (h *ProductHandler) PublicGetProductByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    nil,
			"error":   gin.H{"code": "BAD_REQUEST", "message": "invalid product id"},
		})
		return
	}

	product, err := h.productModule.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, product)
}

// PublicListProducts returns all products with filtering and pagination
// @Summary List all products (public)
// @Description Fetch all available products across stores with filtering/pagination
// @Tags Product
// @Produce json
// @Param category query string false "Category"
// @Param query query string false "Search query"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /products [get]
func (h *ProductHandler) PublicListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	query := c.Query("query")

	params := dto.ProductFilterParams{
		PaginationParams: dto.PaginationParams{
			Page:     page,
			PageSize: pageSize,
		},
		Category: category,
		Query:    query,
	}

	resp, err := h.productModule.ListAllProducts(c.Request.Context(), params)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, resp)
}
