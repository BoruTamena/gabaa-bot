package product

import (
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *ProductHandler) CreateInquiry(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, "invalid product ID", http.StatusBadRequest))
		return
	}

	var req dto.CreateProductInquiryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}

	inquiry, err := h.productModule.CreateProductInquiry(c.Request.Context(), c.GetInt64("user_id"), productID, req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusCreated, inquiry)
}

func (h *ProductHandler) ListMyInquiries(c *gin.Context) {
	var filter dto.ProductInquiryFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}

	resp, err := h.productModule.ListMyInquiries(c.Request.Context(), c.GetInt64("user_id"), filter)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, resp)
}

func (h *ProductHandler) ListStoreInquiries(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		c.Error(errorx.New(errorx.ErrUnauthorized, "Store context missing", http.StatusUnauthorized))
		return
	}

	var filter dto.ProductInquiryFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}

	resp, err := h.productModule.ListStoreInquiries(c.Request.Context(), storeID, filter)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, resp)
}

func (h *ProductHandler) GetStoreInquiry(c *gin.Context) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		c.Error(errorx.New(errorx.ErrUnauthorized, "Store context missing", http.StatusUnauthorized))
		return
	}

	inquiryID, err := strconv.ParseInt(c.Param("inquiry_id"), 10, 64)
	if err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, "invalid inquiry ID", http.StatusBadRequest))
		return
	}

	resp, err := h.productModule.GetStoreInquiry(c.Request.Context(), storeID, inquiryID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, resp)
}

func (h *ProductHandler) ApproveInquiry(c *gin.Context) {
	h.reviewInquiry(c, true)
}

func (h *ProductHandler) RejectInquiry(c *gin.Context) {
	h.reviewInquiry(c, false)
}

func (h *ProductHandler) reviewInquiry(c *gin.Context, approve bool) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		c.Error(errorx.New(errorx.ErrUnauthorized, "Store context missing", http.StatusUnauthorized))
		return
	}

	role := c.GetString("role")
	if role != "admin" {
		c.Error(errorx.New(errorx.ErrForbidden, "Unauthorized to review inquiries", http.StatusForbidden))
		return
	}

	inquiryID, err := strconv.ParseInt(c.Param("inquiry_id"), 10, 64)
	if err != nil {
		c.Error(errorx.New(errorx.ErrBadRequest, "invalid inquiry ID", http.StatusBadRequest))
		return
	}

	var req dto.ReviewProductInquiryRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		c.Error(errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}

	var resp *dto.ProductInquiry
	if approve {
		resp, err = h.productModule.ApproveProductInquiry(c.Request.Context(), storeID, c.GetInt64("user_id"), inquiryID, req)
	} else {
		resp, err = h.productModule.RejectProductInquiry(c.Request.Context(), storeID, c.GetInt64("user_id"), inquiryID, req)
	}
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, resp)
}
