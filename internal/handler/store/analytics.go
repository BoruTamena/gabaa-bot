package store

import (
	"net/http"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsModule module.AnalyticsModule
}

func NewAnalyticsHandler(aModule module.AnalyticsModule) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsModule: aModule,
	}
}

func getMerchantStoreContext(c *gin.Context) (int64, error) {
	storeID := c.GetInt64("store_id")
	if storeID == 0 {
		return 0, errorx.New(errorx.ErrForbidden, "Store context missing", http.StatusForbidden)
	}

	role := c.GetString("role")
	if role != "admin" {
		return 0, errorx.New(errorx.ErrForbidden, "Merchant access required", http.StatusForbidden)
	}

	return storeID, nil
}

// GetSalesAnalytics retrieves sales analytics for the store
// @Summary Get sales analytics
// @Description Returns sales metrics, trend data by period, and top selling products for the authenticated store.
// @Tags Store Analytics
// @Produce json
// @Param from query string false "Start date (RFC3339 format, e.g. 2026-07-01T00:00:00Z)"
// @Param to query string false "End date (RFC3339 format, e.g. 2026-07-08T23:59:59Z)"
// @Success 200 {object} response.BaseResponse{data=dto.SalesAnalytics}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/analytics/sales [get]
func (h *AnalyticsHandler) GetSalesAnalytics(c *gin.Context) {
	storeID, err := getMerchantStoreContext(c)
	if err != nil {
		response.CustomError(c, err.(*errorx.AppError))
		return
	}

	var filter dto.AnalyticsFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		response.CustomError(c, appErr)
		return
	}

	data, err := h.analyticsModule.GetSalesAnalytics(c.Request.Context(), storeID, filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, data)
}

// GetOrderAnalytics retrieves order analytics for the store
// @Summary Get order analytics
// @Description Returns order status breakdown, cancellation rates, and recent order volume for the authenticated store.
// @Tags Store Analytics
// @Produce json
// @Param from query string false "Start date (RFC3339 format, e.g. 2026-07-01T00:00:00Z)"
// @Param to query string false "End date (RFC3339 format, e.g. 2026-07-08T23:59:59Z)"
// @Success 200 {object} response.BaseResponse{data=dto.OrderAnalytics}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/analytics/orders [get]
func (h *AnalyticsHandler) GetOrderAnalytics(c *gin.Context) {
	storeID, err := getMerchantStoreContext(c)
	if err != nil {
		response.CustomError(c, err.(*errorx.AppError))
		return
	}

	var filter dto.AnalyticsFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		response.CustomError(c, appErr)
		return
	}

	data, err := h.analyticsModule.GetOrderAnalytics(c.Request.Context(), storeID, filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, data)
}

// GetProductAnalytics retrieves product catalog analytics for the store
// @Summary Get product analytics
// @Description Returns product catalog status counts, stock alerts, and top viewed products for the authenticated store.
// @Tags Store Analytics
// @Produce json
// @Param from query string false "Start date (RFC3339 format, e.g. 2026-07-01T00:00:00Z)"
// @Param to query string false "End date (RFC3339 format, e.g. 2026-07-08T23:59:59Z)"
// @Success 200 {object} response.BaseResponse{data=dto.ProductAnalytics}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/analytics/products [get]
func (h *AnalyticsHandler) GetProductAnalytics(c *gin.Context) {
	storeID, err := getMerchantStoreContext(c)
	if err != nil {
		response.CustomError(c, err.(*errorx.AppError))
		return
	}

	var filter dto.AnalyticsFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		response.CustomError(c, appErr)
		return
	}

	data, err := h.analyticsModule.GetProductAnalytics(c.Request.Context(), storeID, filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, data)
}

// GetStoryAnalytics retrieves stories analytics for the store
// @Summary Get stories analytics
// @Description Returns stats on active/expired stories and top viewed product stories for the authenticated store.
// @Tags Store Analytics
// @Produce json
// @Param from query string false "Start date (RFC3339 format, e.g. 2026-07-01T00:00:00Z)"
// @Param to query string false "End date (RFC3339 format, e.g. 2026-07-08T23:59:59Z)"
// @Success 200 {object} response.BaseResponse{data=dto.StoryAnalytics}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/analytics/stories [get]
func (h *AnalyticsHandler) GetStoryAnalytics(c *gin.Context) {
	storeID, err := getMerchantStoreContext(c)
	if err != nil {
		response.CustomError(c, err.(*errorx.AppError))
		return
	}

	var filter dto.AnalyticsFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		response.CustomError(c, appErr)
		return
	}

	data, err := h.analyticsModule.GetStoryAnalytics(c.Request.Context(), storeID, filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, data)
}
