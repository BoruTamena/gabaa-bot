package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
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
		storeIDStr := strings.TrimPrefix(payload, "link_store_")
		storeID, err := strconv.ParseInt(storeIDStr, 10, 64)
		if err != nil {
			return c.Send("❌ Invalid store ID.")
		}

		ctx := context.Background()
		user, err := m.userStorage.GetUserByTelegramID(ctx, c.Sender().ID)
		if err != nil {
			return c.Send("❌ User record not found. Please ensure you are logged in to the web app.")
		}

		// Save the store we are currently trying to link
		user.PendingStoreID = &storeID
		if err := m.userStorage.UpdateUser(ctx, user); err != nil {
			logger.Error("failed to update user pending store", zap.Error(err))
			return c.Send("❌ Failed to initiate linking. Please try again from the dashboard.")
		}

		return c.Send("🚀 *Store linking initiated!*\n\nNow, add me to your Group or Channel as an *Administrator* to complete the setup.")
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

		// Check if the chat is already linked to ANOTHER store
		existing, err := m.storeStorage.GetStoreByChatID(ctx, chatID)
		if err == nil && existing.ID != 0 {
			m.tele.GetBot().Send(&telebot.User{ID: merchantTGID},
				fmt.Sprintf("❌ Linking failed! This group is already linked to your other store: '%s'.", existing.Name))
			return nil
		}

		// Determine which store to link
		var targetStoreID int64
		if user.PendingStoreID != nil {
			targetStoreID = *user.PendingStoreID
		} else {
			// Fallback: If merchant has only one store, link that one
			stores, _ := m.storeStorage.GetStoresBySellerID(ctx, user.ID)
			if len(stores) == 1 {
				targetStoreID = stores[0].ID
			} else {
				m.tele.GetBot().Send(&telebot.User{ID: merchantTGID},
					"⚠️ Could not determine which store to link. Please go to your store dashboard and click 'Connect Bot' again.")
				return nil
			}
		}

		store, err := m.storeStorage.GetStoreByID(ctx, targetStoreID)
		if err != nil {
			logger.Error("failed to get store for linking", zap.Error(err), zap.Int64("store_id", targetStoreID))
			return nil
		}

		logger.Info("Linking store to chat",
			zap.Int64("store_id", store.ID),
			zap.Int64("chat_id", chatID),
			zap.String("chat_title", chatTitle))

		store.TelegramChatID = chatID
		store.TelegramChatTitle = chatTitle
		store.Status = constant.StoreStatusLaunched

		if err := m.storeStorage.UpdateStore(ctx, store); err != nil {
			logger.Error("failed to update store after linking", zap.Error(err), zap.Int64("store_id", store.ID))
			m.tele.GetBot().Send(&telebot.User{ID: merchantTGID}, "❌ Failed to link store due to a technical error.")
			return nil
		}

		// Reset pending store
		user.PendingStoreID = nil
		m.userStorage.UpdateUser(ctx, user)

		// Send PRIVATE confirmation to merchant
		m.tele.GetBot().Send(&telebot.User{ID: merchantTGID},
			fmt.Sprintf("✅ Success! Store '%s' is now linked to '%s' and is live! 🚀", store.Name, chatTitle))
	}

	return nil
}
