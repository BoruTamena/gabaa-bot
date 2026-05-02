package telegram

import (
	"net/http"

	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/gin-gonic/gin"
	"gopkg.in/telebot.v4"
)

type WebhookHandler struct {
	tg platform.Telegram
}

func NewWebhookHandler(tg platform.Telegram) *WebhookHandler {
	return &WebhookHandler{tg: tg}
}

// HandleUpdate godoc
// @Summary Telegram Webhook
// @Description Handle incoming updates from Telegram via webhook
// @Tags Telegram
// @Accept json
// @Produce json
// @Param update body telebot.Update true "Telegram Update"
// @Success 200 {string} string "OK"
// @Router /api/v1/webhook/telegram [post]
func (h *WebhookHandler) HandleUpdate(c *gin.Context) {
	var update telebot.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid update format"})
		return
	}

	h.tg.ProcessUpdate(update)
	c.Status(http.StatusOK)
}
