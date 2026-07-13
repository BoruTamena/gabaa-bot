package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	authmod "github.com/BoruTamena/gabaa-bot/internal/module/auth"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/BoruTamena/gabaa-bot/platform"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
)

type botModule struct {
	userStorage          storage.UserStorage
	storeStorage         storage.StoreStorage
	categoryStorage      storage.CategoryStorage
	recommendationModule module.RecommendationModule
	authModule           module.AuthModule
	deliveryModule       module.DeliveryModule
	tele                 platform.Telegram
}

func NewBotModule(
	uStorage storage.UserStorage,
	sStorage storage.StoreStorage,
	cStorage storage.CategoryStorage,
	rModule module.RecommendationModule,
	aModule module.AuthModule,
	dModule module.DeliveryModule,
	tele platform.Telegram,
) module.BotModule {
	m := &botModule{
		userStorage:          uStorage,
		storeStorage:         sStorage,
		categoryStorage:      cStorage,
		recommendationModule: rModule,
		authModule:           aModule,
		deliveryModule:       dModule,
		tele:                 tele,
	}
	m.registerHandlers()
	return m
}

func (m *botModule) registerHandlers() {
	bot := m.tele.GetBot()

	bot.Handle("/start", m.handleStart)
	bot.Handle("/preferences", m.handlePreferences)
	bot.Handle("/recommendations", m.handleRecommendations)
	bot.Handle(&telebot.Btn{Unique: "pref_toggle"}, m.handlePreferenceToggle)
	bot.Handle(telebot.OnMyChatMember, m.handleMyChatMember)
}

func (m *botModule) handleStart(c telebot.Context) error {
	ctx := context.Background()
	username := ""
	if c.Sender() != nil {
		username = c.Sender().Username
	}

	if err := m.recommendationModule.SetBotStarted(ctx, c.Sender().ID, username); err != nil {
		logger.Error("failed to mark bot started for user", zap.Error(err), zap.Int64("telegram_id", c.Sender().ID))
	}

	payload := c.Message().Payload
	if strings.HasPrefix(payload, "login_") {
		sessionID := strings.TrimPrefix(payload, "login_")
		tgUser := &dto.TelegramUser{
			ID:       c.Sender().ID,
			Username: username,
		}
		if err := m.authModule.CompleteBotLoginSession(ctx, sessionID, tgUser); err != nil {
			if authmod.IsSessionNotFound(err) {
				return c.Send("❌ This login link is invalid or has expired. Please request a new one from the app.")
			}
			logger.Error("failed to complete bot login session", zap.Error(err), zap.String("session_id", sessionID))
			return c.Send("❌ Failed to complete login. Please try again from the app.")
		}
		return c.Send("✅ You're logged in to Gabaa! Return to the app to continue.")
	}

	if strings.HasPrefix(payload, "link_store_") {
		storeIDStr := strings.TrimPrefix(payload, "link_store_")
		storeID, err := strconv.ParseInt(storeIDStr, 10, 64)
		if err != nil {
			return c.Send("❌ Invalid store ID.")
		}

		user, err := m.userStorage.GetUserByTelegramID(ctx, c.Sender().ID)
		if err != nil {
			return c.Send("❌ User record not found. Please ensure you are logged in to the web app.")
		}

		user.PendingStoreID = &storeID
		if err := m.userStorage.UpdateUser(ctx, user); err != nil {
			logger.Error("failed to update user pending store", zap.Error(err))
			return c.Send("❌ Failed to initiate linking. Please try again from the dashboard.")
		}

		return c.Send("🚀 *Store linking initiated!*\n\nNow, add me to your Group or Channel as an *Administrator* to complete the setup.")
	}

	user, err := m.userStorage.GetUserByTelegramID(ctx, c.Sender().ID)
	if err == nil && m.deliveryModule != nil {
		activated, actErr := m.deliveryModule.ActivatePendingInvite(ctx, c.Sender().ID, username, user.ID)
		if actErr != nil {
			logger.Error("failed to activate delivery invite", zap.Error(actErr))
		}
		if activated {
			deliveryURL := m.tele.DeliveryAppURL()
			msg := "🚚 *Welcome to Gabaa Delivery!*\n\nYou've been connected as a delivery partner. Open the delivery app to view assigned orders."
			if deliveryURL != "" {
				msg += "\n\n" + deliveryURL
			}
			return c.Send(msg, telebot.ModeMarkdown)
		}
	}

	return c.Send("Welcome to Gabaa Bot! 🛍️\n\nUse /preferences to choose product categories and /recommendations on to get new product alerts in this chat.")
}

func (m *botModule) handlePreferences(c telebot.Context) error {
	ctx := context.Background()
	if err := m.recommendationModule.SetBotStarted(ctx, c.Sender().ID, c.Sender().Username); err != nil {
		logger.Error("failed to mark bot started for user", zap.Error(err))
	}

	user, err := m.userStorage.GetUserByTelegramID(ctx, c.Sender().ID)
	if err != nil {
		return c.Send("❌ Please open the Gabaa Mini App first so we can create your account, then try /preferences again.")
	}

	prefs, err := m.recommendationModule.GetPreferences(ctx, user.ID)
	if err != nil {
		return c.Send("❌ Failed to load your preferences. Please try again.")
	}

	markup, err := m.buildPreferencesKeyboard(ctx, prefs.Categories)
	if err != nil {
		return c.Send("❌ Failed to load categories. Please try again.")
	}

	message := m.preferencesMessage(prefs.Enabled, prefs.Categories)

	return c.Send(message, markup, telebot.ModeMarkdown)
}

func (m *botModule) preferencesMessage(enabled bool, selected []string) string {
	status := "off"
	if enabled {
		status = "on"
	}

	selectedLine := "_None selected yet_"
	if len(selected) > 0 {
		selectedLine = strings.Join(selected, ", ")
	}

	return fmt.Sprintf(
		"🎯 *Product Recommendations*\n\n"+
			"Status: *%s*\n"+
			"Selected: %s\n\n"+
			"Tap categories below to subscribe or unsubscribe.\n"+
			"✅ = selected\n\n"+
			"Use /recommendations on or /recommendations off to control alerts.",
		status,
		selectedLine,
	)
}

func (m *botModule) handleRecommendations(c telebot.Context) error {
	ctx := context.Background()
	if err := m.recommendationModule.SetBotStarted(ctx, c.Sender().ID, c.Sender().Username); err != nil {
		logger.Error("failed to mark bot started for user", zap.Error(err))
	}

	payload := strings.TrimSpace(strings.ToLower(c.Message().Payload))
	switch payload {
	case "on":
		if err := m.recommendationModule.SetRecommendationsEnabled(ctx, c.Sender().ID, true); err != nil {
			return c.Send("❌ Failed to enable recommendations. Please open the Mini App first.")
		}
		return c.Send("✅ Recommendations are *on*. You'll get new product alerts for your selected categories.\n\nUse /preferences to choose categories.", telebot.ModeMarkdown)
	case "off":
		if err := m.recommendationModule.SetRecommendationsEnabled(ctx, c.Sender().ID, false); err != nil {
			return c.Send("❌ Failed to disable recommendations.")
		}
		return c.Send("🔕 Recommendations are *off*. You won't receive product alerts.", telebot.ModeMarkdown)
	default:
		return c.Send("Usage:\n/recommendations on\n/recommendations off")
	}
}

func (m *botModule) handlePreferenceToggle(c telebot.Context) error {
	categoryID, err := strconv.ParseInt(strings.TrimSpace(c.Data()), 10, 64)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Invalid category"})
	}

	ctx := context.Background()
	category, err := m.categoryStorage.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Category not found"})
	}

	added, err := m.recommendationModule.ToggleCategory(ctx, c.Sender().ID, category.Name)
	if err != nil {
		logger.Error("failed to toggle preference category", zap.Error(err))
		return c.Respond(&telebot.CallbackResponse{Text: "Failed to update preference"})
	}

	user, err := m.userStorage.GetUserByTelegramID(ctx, c.Sender().ID)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "User not found"})
	}

	prefs, err := m.recommendationModule.GetPreferences(ctx, user.ID)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Failed to refresh preferences"})
	}

	markup, err := m.buildPreferencesKeyboard(ctx, prefs.Categories)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Failed to refresh categories"})
	}

	message := m.preferencesMessage(prefs.Enabled, prefs.Categories)

	if err := c.Edit(message, markup, telebot.ModeMarkdown); err != nil {
		if err := c.Edit(markup); err != nil {
			logger.Error("failed to edit preferences keyboard", zap.Error(err))
		}
	}

	action := "removed"
	if added {
		action = "selected"
	}
	return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s %s", category.Name, action)})
}

func (m *botModule) buildPreferencesKeyboard(ctx context.Context, selected []string) (*telebot.ReplyMarkup, error) {
	categories, _, err := m.categoryStorage.GetAllCategories(ctx, 100, 0)
	if err != nil {
		return nil, err
	}

	selectedMap := make(map[string]bool)
	for _, category := range selected {
		selectedMap[strings.ToLower(strings.TrimSpace(category))] = true
	}

	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	for _, category := range categories {
		if category.StoreID != 0 {
			continue
		}

		label := category.Name
		if selectedMap[strings.ToLower(category.Name)] {
			label = "✅ " + category.Name
		}

		btn := markup.Data(label, "pref_toggle", strconv.FormatInt(category.ID, 10))
		rows = append(rows, markup.Row(btn))
	}

	if len(rows) == 0 {
		return markup, nil
	}

	markup.Inline(rows...)
	return markup, nil
}

func (m *botModule) handleMyChatMember(c telebot.Context) error {
	update := c.ChatMember()

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

		existing, err := m.storeStorage.GetStoreByChatID(ctx, chatID)
		if err == nil && existing.ID != 0 {
			m.tele.GetBot().Send(&telebot.User{ID: merchantTGID},
				fmt.Sprintf("❌ Linking failed! This group is already linked to your other store: '%s'.", existing.Name))
			return nil
		}

		var targetStoreID int64
		if user.PendingStoreID != nil {
			targetStoreID = *user.PendingStoreID
		} else {
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

		user.PendingStoreID = nil
		m.userStorage.UpdateUser(ctx, user)

		m.tele.GetBot().Send(&telebot.User{ID: merchantTGID},
			fmt.Sprintf("✅ Success! Store '%s' is now linked to '%s' and is live! 🚀", store.Name, chatTitle))
	}

	return nil
}
