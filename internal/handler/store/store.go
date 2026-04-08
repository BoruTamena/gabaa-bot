package store

import (
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
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
// @Accept json
// @Produce json
// @Router /store/from-chat [post]
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var req dto.CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		c.JSON(appErr.Status, appErr)
		return
	}

	userID := c.GetInt64("user_id")
	role := c.GetString("role")
	if role != "admin" {
		appErr := errorx.New(errorx.ErrForbidden, "Only admins can create stores", http.StatusForbidden)
		c.JSON(appErr.Status, appErr)
		return
	}

	store, err := h.storeModule.CreateStore(c.Request.Context(), userID, req)
	if err != nil {
		// Ozzo validation returns errors that we should handle
		appErr := errorx.New(errorx.ErrValidation, err.Error(), http.StatusUnprocessableEntity)
		c.JSON(appErr.Status, appErr)
		return
	}

	c.JSON(http.StatusCreated, store)
}

// GetStore retrieves store profile
// @Summary Get store by ID
// @Description Returns store details
// @Produce json
// @Router /store/:id [get]
func (h *StoreHandler) GetStore(c *gin.Context) {
	idStr := c.Param("store_id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	store, err := h.storeModule.GetStore(c.Request.Context(), id)
	if err != nil {
		appErr := errorx.New(errorx.ErrNotFound, "Store not found", http.StatusNotFound)
		c.JSON(appErr.Status, appErr)
		return
	}

	c.JSON(http.StatusOK, store)
}

// UpdateStore updates store profile
// @Summary Update store
// @Description Update store details. Admin only.
// @Accept json
// @Produce json
// @Router /store/:id [put]
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	idStr := c.Param("store_id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var req dto.UpdateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errorx.New(errorx.ErrBadRequest, err.Error(), http.StatusBadRequest)
		c.JSON(appErr.Status, appErr)
		return
	}

	role := c.GetString("role")
	userStoreID := c.GetInt64("store_id")

	if role != "admin" || (userStoreID != 0 && userStoreID != id) {
		appErr := errorx.New(errorx.ErrForbidden, "Unauthorized to update this store", http.StatusForbidden)
		c.JSON(appErr.Status, appErr)
		return
	}

	store, err := h.storeModule.UpdateStore(c.Request.Context(), id, req)
	if err != nil {
		appErr := errorx.New(errorx.ErrValidation, err.Error(), http.StatusUnprocessableEntity)
		c.JSON(appErr.Status, appErr)
		return
	}

	c.JSON(http.StatusOK, store)
}

// GetDashboard returns the appropriate dashboard type for the user
// @Summary Get admin dashboard info
// @Description Returns dashboardType: 'setup', 'manage', or 'storefront'
// @Produce json
// @Router /store/dashboard/:chat_id [get]
func (h *StoreHandler) GetDashboard(c *gin.Context) {
	chatIDStr := c.Param("chat_id")
	chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)
	userID := c.GetInt64("user_id")

	dashboardType, store, err := h.storeModule.GetAdminDashboard(c.Request.Context(), userID, chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dashboard_type": dashboardType,
		"store":          store,
	})
}
