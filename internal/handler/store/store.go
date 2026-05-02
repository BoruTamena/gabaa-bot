package store

import (
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/BoruTamena/gabaa-bot/pkg/response"
	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeModule module.StoreModule
}

func NewStoreHandler(sModule module.StoreModule) *StoreHandler {
	return &StoreHandler{storeModule: sModule}
}

// CreateStore handles first-time store setup
// @Summary Create store
// @Description Setup a new store. Admin only.
// @Tags Store
// @Accept json
// @Produce json
// @Param request body dto.CreateStoreRequest true "Store details"
// @Success 201 {object} response.BaseResponse{data=dto.Store}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 422 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/from-chat [post]
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var req dto.CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		response.CustomError(c, appErr)
		return
	}

	userID := c.GetInt64("user_id")
	
	store, err := h.storeModule.CreateStore(c.Request.Context(), userID, req)
	if err != nil {
		// Ozzo validation returns errors that we should handle
		appErr := errorx.New(errorx.ErrValidation, err.Error(), http.StatusUnprocessableEntity)
		response.CustomError(c, appErr)
		return
	}

	response.Success(c, http.StatusCreated, store)
}

// GetStore retrieves store profile
// @Summary Get store by ID
// @Description Returns store details
// @Tags Store
// @Produce json
// @Param store_id path int true "Store ID"
// @Success 200 {object} response.BaseResponse{data=dto.Store}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 404 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/:store_id [get]
func (h *StoreHandler) GetStore(c *gin.Context) {
	idStr := c.Param("store_id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	store, err := h.storeModule.GetStore(c.Request.Context(), id)
	if err != nil {
		appErr := errorx.New(errorx.ErrNotFound, "Store not found", http.StatusNotFound)
		response.CustomError(c, appErr)
		return
	}

	response.Success(c, http.StatusOK, store)
}

// GetStoreStatus retrieves store status
// @Summary Get store status
// @Description Returns store status: 'pending' or 'launched'
// @Tags Store
// @Produce json
// @Param store_id path int true "Store ID"
// @Success 200 {object} response.BaseResponse{data=string}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 404 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/:store_id/status [get]
func (h *StoreHandler) GetStoreStatus(c *gin.Context) {
	idStr := c.Param("store_id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	status, err := h.storeModule.GetStoreStatus(c.Request.Context(), id)
	if err != nil {
		appErr := errorx.New(errorx.ErrNotFound, "Store not found", http.StatusNotFound)
		response.CustomError(c, appErr)
		return
	}

	response.Success(c, http.StatusOK, status)
}

// UpdateStore updates store profile
// @Summary Update store
// @Description Update store details. Admin only.
// @Tags Store
// @Accept json
// @Produce json
// @Param store_id path int true "Store ID"
// @Param request body dto.UpdateStoreRequest true "Store details"
// @Success 200 {object} response.BaseResponse{data=dto.Store}
// @Failure 400 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 403 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 422 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/:store_id [put]
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	idStr := c.Param("store_id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var req dto.UpdateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		response.CustomError(c, appErr)
		return
	}

	role := c.GetString("role")
	userStoreID := c.GetInt64("store_id")

	if role != "admin" || (userStoreID != 0 && userStoreID != id) {
		appErr := errorx.New(errorx.ErrForbidden, "Unauthorized to update this store", http.StatusForbidden)
		response.CustomError(c, appErr)
		return
	}

	store, err := h.storeModule.UpdateStore(c.Request.Context(), id, req)
	if err != nil {
		appErr := errorx.New(errorx.ErrValidation, err.Error(), http.StatusUnprocessableEntity)
		response.CustomError(c, appErr)
		return
	}

	response.Success(c, http.StatusOK, store)
}

// GetDashboard returns the appropriate dashboard type for the user
// @Summary Get admin dashboard info
// @Description Returns dashboardType: 'setup', 'manage', or 'storefront'
// @Tags Store
// @Produce json
// @Param chat_id path int true "Telegram Chat ID"
// @Success 200 {object} response.BaseResponse{data=map[string]interface{}}
// @Failure 401 {object} response.BaseResponse{error=errorx.AppError}
// @Failure 500 {object} response.BaseResponse{error=errorx.AppError}
// @Router /store/dashboard/:chat_id [get]
func (h *StoreHandler) GetDashboard(c *gin.Context) {
	chatIDStr := c.Param("chat_id")
	chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)
	userID := c.GetInt64("user_id")

	dashboardType, store, err := h.storeModule.GetAdminDashboard(c.Request.Context(), userID, chatID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"dashboard_type": dashboardType,
		"store":          store,
	})
}
