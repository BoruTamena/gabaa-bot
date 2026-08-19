package admin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	storeModule    module.StoreModule
	orderModule    module.OrderModule
	userModule     module.UserModule
	walletModule   module.WalletModule
	paymentModule  module.PaymentModule
	storyModule    module.StoryModule
	deliveryModule module.DeliveryModule
}

func NewHandler(
	storeModule module.StoreModule,
	orderModule module.OrderModule,
	userModule module.UserModule,
	walletModule module.WalletModule,
	paymentModule module.PaymentModule,
	storyModule module.StoryModule,
	deliveryModule module.DeliveryModule,
) *Handler {
	return &Handler{
		storeModule:    storeModule,
		orderModule:    orderModule,
		userModule:     userModule,
		walletModule:   walletModule,
		paymentModule:  paymentModule,
		storyModule:    storyModule,
		deliveryModule: deliveryModule,
	}
}

func parseStoreID(c *gin.Context) (int64, error) {
	storeID, err := strconv.ParseInt(c.Param("store_id"), 10, 64)
	if err != nil || storeID == 0 {
		return 0, errorx.New(errorx.ErrBadRequest, "invalid store id", http.StatusBadRequest)
	}
	return storeID, nil
}

// ListStores returns a paginated list of all stores for platform moderation
// @Summary List stores (platform admin)
// @Description Paginated store list with optional status, verification_status, category, and text search
// @Tags Platform Admin
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by store status (pending, launched, banned)"
// @Param verification_status query string false "Filter by verification status"
// @Param category query string false "Filter by category"
// @Param query query string false "Search name, description, location, or phone"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Router /admin/stores [get]
func (h *Handler) ListStores(c *gin.Context) {
	var filter dto.AdminStoreFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.CustomError(c, errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}
	resp, err := h.storeModule.AdminListStores(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// GetStore returns store detail for platform admin drill-down
// @Summary Get store detail (platform admin)
// @Description Returns full store profile by ID
// @Tags Platform Admin
// @Produce json
// @Param store_id path int true "Store ID"
// @Success 200 {object} response.BaseResponse{data=dto.Store}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 404 {object} response.BaseResponse{error=errorx.AppError}
// @Router /admin/stores/:store_id [get]
func (h *Handler) GetStore(c *gin.Context) {
	storeID, err := parseStoreID(c)
	if err != nil {
		response.CustomError(c, err.(*errorx.AppError))
		return
	}
	store, err := h.storeModule.GetStore(c.Request.Context(), storeID)
	if err != nil {
		response.CustomError(c, errorx.New(errorx.ErrNotFound, "store not found", http.StatusNotFound))
		return
	}
	response.Success(c, http.StatusOK, store)
}

// UpdateStoreStatus updates store status (ban, unban, launch, etc.)
// @Summary Update store status (platform admin)
// @Description Set store status to pending, launched, or banned
// @Tags Platform Admin
// @Accept json
// @Produce json
// @Param store_id path int true "Store ID"
// @Param request body dto.AdminUpdateStoreStatusRequest true "New status"
// @Success 200 {object} response.BaseResponse{data=dto.Store}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Router /admin/stores/:store_id/status [patch]
func (h *Handler) UpdateStoreStatus(c *gin.Context) {
	storeID, err := parseStoreID(c)
	if err != nil {
		response.CustomError(c, err.(*errorx.AppError))
		return
	}
	var req dto.AdminUpdateStoreStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.CustomError(c, errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}
	store, err := h.storeModule.AdminUpdateStoreStatus(c.Request.Context(), storeID, req.Status)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, store)
}

// ListStoreVerifications returns paginated KYC submissions for review
// @Summary List store KYC verifications (platform admin)
// @Description Paginated KYC queue with optional status and search by store name or TIN
// @Tags Platform Admin
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by verification status (pending_review, verified, rejected)"
// @Param query query string false "Search store name or TIN"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Router /admin/store-verifications [get]
func (h *Handler) ListStoreVerifications(c *gin.Context) {
	var filter dto.AdminStoreKYCFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.CustomError(c, errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}
	resp, err := h.storeModule.ListStoreVerificationsAdmin(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// UpsertStoreKYC creates or updates KYC documents for a store on behalf of platform admin
// @Summary Upsert store KYC (platform admin)
// @Description Submit or update KYC documents for a store and set verification to pending_review
// @Tags Platform Admin
// @Accept json
// @Produce json
// @Param store_id path int true "Store ID"
// @Param request body dto.SubmitStoreKYCRequest true "KYC documents"
// @Success 200 {object} response.BaseResponse{data=dto.StoreKYCResponse}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 422 {object} response.BaseResponse{error=errorx.AppError}
// @Router /admin/stores/:store_id/kyc [post]
func (h *Handler) UpsertStoreKYC(c *gin.Context) {
	storeID, err := parseStoreID(c)
	if err != nil {
		response.CustomError(c, err.(*errorx.AppError))
		return
	}
	var req dto.SubmitStoreKYCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.CustomError(c, errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}
	resp, err := h.storeModule.AdminUpsertStoreKYC(c.Request.Context(), storeID, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// ListOrders returns paginated orders across all stores
// @Summary List global orders (platform admin)
// @Description Cross-store order list with optional store_id, order_id, status, and customer search
// @Tags Platform Admin
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param store_id query int false "Filter by store ID"
// @Param order_id query int false "Search by exact order ID"
// @Param status query string false "Filter by order status"
// @Param query query string false "Search customer username"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Router /admin/orders [get]
func (h *Handler) ListOrders(c *gin.Context) {
	var filter dto.AdminOrderFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.CustomError(c, errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}
	resp, err := h.orderModule.AdminListOrders(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// ListStoreOrders returns paginated orders for a specific store
// @Summary List store orders (platform admin)
// @Description Paginated orders scoped to a store with optional status, order_id, and customer search
// @Tags Platform Admin
// @Produce json
// @Param store_id path int true "Store ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param order_id query int false "Search by exact order ID"
// @Param status query string false "Filter by order status"
// @Param query query string false "Search customer username"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Router /admin/stores/:store_id/orders [get]
func (h *Handler) ListStoreOrders(c *gin.Context) {
	storeID, err := parseStoreID(c)
	if err != nil {
		response.CustomError(c, err.(*errorx.AppError))
		return
	}
	var filter dto.AdminOrderFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.CustomError(c, errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}
	filter.StoreID = storeID
	resp, err := h.orderModule.AdminListOrders(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// ListUsers returns paginated list of all platform users
// @Summary List users (platform admin)
// @Description Paginated user list with optional role filter and username/telegram ID search
// @Tags Platform Admin
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param role query string false "Filter by user role"
// @Param query query string false "Search username or telegram user ID"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Router /admin/users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	var filter dto.AdminUserFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.CustomError(c, errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}
	resp, err := h.userModule.ListUsers(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// GetStoreWallet returns wallet summary for a store
// @Summary Get store wallet (platform admin)
// @Description Returns pending, available, and locked balances for a store
// @Tags Platform Admin
// @Produce json
// @Param store_id path int true "Store ID"
// @Success 200 {object} response.BaseResponse{data=dto.Wallet}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Router /admin/stores/:store_id/wallet [get]
func (h *Handler) GetStoreWallet(c *gin.Context) {
	storeID, err := parseStoreID(c)
	if err != nil {
		response.CustomError(c, err.(*errorx.AppError))
		return
	}
	wallet, err := h.walletModule.GetWalletSummary(c.Request.Context(), storeID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, wallet)
}

// ListStoreWithdrawals returns paginated withdrawal history for a store
// @Summary List store withdrawals (platform admin)
// @Description Paginated withdrawal history with optional status filter
// @Tags Platform Admin
// @Produce json
// @Param store_id path int true "Store ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by withdrawal status"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Router /admin/stores/:store_id/wallet/withdrawals [get]
func (h *Handler) ListStoreWithdrawals(c *gin.Context) {
	storeID, err := parseStoreID(c)
	if err != nil {
		response.CustomError(c, err.(*errorx.AppError))
		return
	}
	var filter dto.AdminWithdrawalFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.CustomError(c, errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}
	resp, err := h.walletModule.ListWithdrawals(c.Request.Context(), storeID, filter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// ListStoreTransactions returns paginated payment transactions for a store
// @Summary List store transactions (platform admin)
// @Description Paginated checkout payments with optional status, medium, and reference search
// @Tags Platform Admin
// @Produce json
// @Param store_id path int true "Store ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by payment status (initiated, pending, success, failed)"
// @Param medium query string false "Filter by payment medium"
// @Param query query string false "Search reference or transaction ID"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Router /admin/stores/:store_id/transactions [get]
func (h *Handler) ListStoreTransactions(c *gin.Context) {
	storeID, err := parseStoreID(c)
	if err != nil {
		response.CustomError(c, err.(*errorx.AppError))
		return
	}
	var filter dto.AdminTransactionFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.CustomError(c, errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}
	paymentFilter := dto.PaymentFilterParams{
		PaginationParams: filter.PaginationParams,
		StoreID:          storeID,
		Status:           filter.Status,
		Medium:           filter.Medium,
		Query:            filter.Query,
	}
	resp, err := h.paymentModule.ListStoreTransactions(c.Request.Context(), paymentFilter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// ListStoreStories returns paginated story ads for a store
// @Summary List store stories (platform admin)
// @Description Paginated story ads with optional is_active, type, and caption search
// @Tags Platform Admin
// @Produce json
// @Param store_id path int true "Store ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param is_active query bool false "Filter by active status"
// @Param type query string false "Filter by media type (image, video)"
// @Param query query string false "Search story caption"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Router /admin/stores/:store_id/stories [get]
func (h *Handler) ListStoreStories(c *gin.Context) {
	storeID, err := parseStoreID(c)
	if err != nil {
		response.CustomError(c, err.(*errorx.AppError))
		return
	}
	var filter dto.AdminStoryFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.CustomError(c, errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}
	storyFilter := dto.ProductStoryFilterParams{
		PaginationParams: filter.PaginationParams,
		StoreID:          storeID,
		IsActive:         filter.IsActive,
		Type:             filter.Type,
		Search:           filter.Query,
	}
	resp, err := h.storyModule.ListMyStories(c.Request.Context(), storyFilter)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, resp)
}

// ListStoreDeliveries returns paginated delivery agents linked to a store
// @Summary List store delivery agents (platform admin)
// @Description Paginated connected delivery agents with optional status and username search
// @Tags Platform Admin
// @Produce json
// @Param store_id path int true "Store ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Filter by agent link status"
// @Param query query string false "Search agent username or full name"
// @Success 200 {object} response.BaseResponse{data=dto.PaginatedResponse}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Router /admin/stores/:store_id/deliveries [get]
func (h *Handler) ListStoreDeliveries(c *gin.Context) {
	storeID, err := parseStoreID(c)
	if err != nil {
		response.CustomError(c, err.(*errorx.AppError))
		return
	}
	var filter dto.AdminDeliveryFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.CustomError(c, errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest))
		return
	}
	agents, err := h.deliveryModule.ListAgents(c.Request.Context(), storeID)
	if err != nil {
		response.Error(c, err)
		return
	}
	filtered := filterDeliveryAgents(agents, filter)
	page, total := paginateAgents(filtered, filter.PaginationParams)
	response.Success(c, http.StatusOK, dto.NewPaginatedResponse(page, total, filter.PaginationParams))
}

func filterDeliveryAgents(agents []dto.DeliveryAgentResponse, filter dto.AdminDeliveryFilterParams) []dto.DeliveryAgentResponse {
	out := agents
	if filter.Status != "" {
		next := make([]dto.DeliveryAgentResponse, 0)
		for _, a := range out {
			if a.Status == filter.Status {
				next = append(next, a)
			}
		}
		out = next
	}
	if filter.Query != "" {
		q := strings.ToLower(filter.Query)
		next := make([]dto.DeliveryAgentResponse, 0)
		for _, a := range out {
			if strings.Contains(strings.ToLower(a.Username), q) || strings.Contains(strings.ToLower(a.FullName), q) {
				next = append(next, a)
			}
		}
		out = next
	}
	return out
}

func paginateAgents(agents []dto.DeliveryAgentResponse, params dto.PaginationParams) ([]dto.DeliveryAgentResponse, int64) {
	total := int64(len(agents))
	start := params.GetOffset()
	if start >= len(agents) {
		return []dto.DeliveryAgentResponse{}, total
	}
	end := start + params.GetLimit()
	if end > len(agents) {
		end = len(agents)
	}
	return agents[start:end], total
}
