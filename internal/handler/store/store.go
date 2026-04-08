package store

import (
	"net/http"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeModule module.StoreModule
}

func NewStoreHandler(sModule module.StoreModule) *StoreHandler {
	return &StoreHandler{storeModule: sModule}
}

// CreateStoreFromChat creates a new store from a Telegram group/channel
// @Summary Create store from chat
// @Description Create a new store linked to a Telegram chat. Admin only.
// @Accept json
// @Produce json
// @Header 200 {string} X-Telegram-Init-Data "MiniApp init data"
// @Router /store/from-chat [post]
func (h *StoreHandler) CreateStoreFromChat(c *gin.Context) {
	var req struct {
		ChatID   int64  `json:"chat_id" binding:"required"`
		ChatType string `json:"chat_type" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt64("user_id")

	store, err := h.storeModule.CreateStore(c.Request.Context(), userID, req.ChatID, req.ChatType, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
