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

// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		appErr := errorx.New(errorx.ErrUnauthorized, "Store context missing", http.StatusUnauthorized)
		c.Error(appErr)
		return
	}

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

// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/product [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		appErr := errorx.New(errorx.ErrUnauthorized, "Store context missing", http.StatusUnauthorized)
		c.Error(appErr)
		return
	}

	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		c.Error(appErr)
		return
	}

	role := c.GetString("role")

	if role != "admin" {
		appErr := errorx.New(errorx.ErrForbidden, "Unauthorized to add products", http.StatusForbidden)
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

// @Success 200 {object} response.BaseResponse{data=dto.Product}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 422 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/product/:id [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		appErr := errorx.New(errorx.ErrUnauthorized, "Store context missing", http.StatusUnauthorized)
		c.Error(appErr)
		return
	}
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		c.Error(appErr)
		return
	}

	role := c.GetString("role")

	if role != "admin" {
		appErr := errorx.New(errorx.ErrForbidden, "Unauthorized to update products", http.StatusForbidden)
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

// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /my-store/product/:id [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		appErr := errorx.New(errorx.ErrUnauthorized, "Store context missing", http.StatusUnauthorized)
		c.Error(appErr)
		return
	}
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	role := c.GetString("role")

	if role != "admin" {
		appErr := errorx.New(errorx.ErrForbidden, "Unauthorized to delete products", http.StatusForbidden)
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
