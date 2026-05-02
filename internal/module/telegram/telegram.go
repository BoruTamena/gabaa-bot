package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/BoruTamena/gabaa-bot/platform"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
)

type botModule struct {
	userStorage  storage.UserStorage
	storeStorage storage.StoreStorage
	tele         platform.Telegram
}

func NewBotModule(uStorage storage.UserStorage, sStorage storage.StoreStorage, tele platform.Telegram) module.BotModule {
	m := &botModule{
		userStorage:  uStorage,
		storeStorage: sStorage,
		tele:         tele,
	}
	m.registerHandlers()
	return m
}

func (m *botModule) registerHandlers() {
	bot := m.tele.GetBot()

	// Handle /start command
	bot.Handle("/start", m.handleStart)

	// Handle my_chat_member updates
	bot.Handle(telebot.OnMyChatMember, m.handleMyChatMember)
}

func (m *botModule) handleStart(c telebot.Context) error {
	payload := c.Message().Payload
	if strings.HasPrefix(payload, "link_store_") {
		// This is where we could handle deep-linked store association
		// For now just acknowledge
		return c.Send("Starting store linking process... Please add me to your group as an admin.")
	}
	return c.Send("Welcome to Gabaa Bot! 🛍️\n\nUse our web dashboard to manage your store.")
}

func (m *botModule) handleMyChatMember(c telebot.Context) error {
	update := c.ChatMember()
	
	// Check if the bot was added as an administrator
	if update.NewChatMember.Role == telebot.Administrator {
		merchantTGID := update.Sender.ID
		chatID := c.Chat().ID
		chatTitle := c.Chat().Title
		
		ctx := context.Background()
		user, err := m.userStorage.GetUserByTelegramID(ctx, merchantTGID)
		if err != nil {
			logger.Warn("Bot added to group by unknown user", zap.Int64("telegram_id", merchantTGID), zap.Int64("chat_id", chatID))
			return nil
		}

		stores, err := m.storeStorage.GetStoresBySellerID(ctx, user.ID)
		if err != nil || len(stores) == 0 {
			logger.Warn("Merchant has no stores to link", zap.Int64("user_id", user.ID))
			return nil
		}

		// Linking logic: For now, link to the first store found for the merchant
		// In a production scenario, we'd use the state saved during /start deep-linking
		store := stores[0]
		
		logger.Info("Silently linking store to chat", 
			zap.Int64("store_id", store.ID), 
			zap.Int64("chat_id", chatID), 
			zap.String("chat_title", chatTitle))
		
		// TODO: Update storage to save the linked chat
		// For now, we update the store's TelegramChatID (legacy mapping)
		store.TelegramChatID = chatID
		m.storeStorage.UpdateStore(ctx, &store)

		// Send PRIVATE confirmation to merchant
		m.tele.GetBot().Send(&telebot.User{ID: merchantTGID}, 
			fmt.Sprintf("✅ Store '%s' successfully linked to group '%s'!", store.Name, chatTitle))
	}
	
	return nil
}
