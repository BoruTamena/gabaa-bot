package recommendation

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/BoruTamena/gabaa-bot/platform"
	"go.uber.org/zap"
)

const recommendationSendDelay = 30 * time.Millisecond

type recommendationModule struct {
	preferenceStorage     storage.PreferenceStorage
	recommendationStorage storage.RecommendationStorage
	userStorage           storage.UserStorage
	storeStorage          storage.StoreStorage
	tele                  platform.Telegram
}

func NewRecommendationModule(
	prefStorage storage.PreferenceStorage,
	recStorage storage.RecommendationStorage,
	userStorage storage.UserStorage,
	storeStorage storage.StoreStorage,
	tele platform.Telegram,
) module.RecommendationModule {
	return &recommendationModule{
		preferenceStorage:     prefStorage,
		recommendationStorage: recStorage,
		userStorage:           userStorage,
		storeStorage:          storeStorage,
		tele:                  tele,
	}
}

func (m *recommendationModule) GetPreferences(ctx context.Context, userID int64) (*dto.UserPreferences, error) {
	user, err := m.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	categories, err := m.preferenceStorage.GetUserPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}
	if categories == nil {
		categories = []string{}
	}

	return &dto.UserPreferences{
		Enabled:    user.RecommendationsEnabled,
		Categories: categories,
	}, nil
}

func (m *recommendationModule) SetPreferences(ctx context.Context, userID int64, req dto.UpdateUserPreferencesRequest) (*dto.UserPreferences, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if err := m.preferenceStorage.SetUserPreferences(ctx, userID, req.Categories); err != nil {
		return nil, err
	}

	user, err := m.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.RecommendationsEnabled = req.Enabled
	if err := m.userStorage.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return m.GetPreferences(ctx, userID)
}

func (m *recommendationModule) SetBotStarted(ctx context.Context, telegramUserID int64, username string) error {
	user, err := m.userStorage.GetUserByTelegramID(ctx, telegramUserID)
	if err != nil {
		user = &db.User{
			TelegramUserID: &telegramUserID,
			Username:       username,
			Role:           "customer",
			BotStarted:     true,
		}
		return m.userStorage.CreateUser(ctx, user)
	}

	if user.BotStarted {
		return nil
	}

	user.BotStarted = true
	if username != "" && user.Username == "" {
		user.Username = username
	}
	return m.userStorage.UpdateUser(ctx, user)
}

func (m *recommendationModule) SetRecommendationsEnabled(ctx context.Context, telegramUserID int64, enabled bool) error {
	user, err := m.userStorage.GetUserByTelegramID(ctx, telegramUserID)
	if err != nil {
		return err
	}

	user.RecommendationsEnabled = enabled
	return m.userStorage.UpdateUser(ctx, user)
}

func (m *recommendationModule) ToggleCategory(ctx context.Context, telegramUserID int64, category string) (bool, error) {
	user, err := m.userStorage.GetUserByTelegramID(ctx, telegramUserID)
	if err != nil {
		return false, err
	}

	return m.preferenceStorage.ToggleUserCategory(ctx, user.ID, category)
}

func (m *recommendationModule) NotifyMatchingUsers(ctx context.Context, product *db.Product, sellerUserID int64) {
	if product == nil || strings.TrimSpace(product.Category) == "" {
		return
	}

	users, err := m.preferenceStorage.GetUsersByCategories(ctx, []string{product.Category})
	if err != nil {
		logger.Error("failed to fetch users for product recommendation", zap.Error(err), zap.Int64("product_id", product.ID))
		return
	}

	storeName := "Gabaa Place"
	if product.StoreID != nil {
		store, storeErr := m.storeStorage.GetStoreByID(ctx, *product.StoreID)
		if storeErr == nil && store.Name != "" {
			storeName = store.Name
		}
	}

	productDTO := m.mapProductToDTO(product)

	for _, user := range users {
		if user.ID == sellerUserID {
			continue
		}
		if user.TelegramUserID == nil {
			continue
		}

		notified, checkErr := m.recommendationStorage.WasNotified(ctx, user.ID, product.ID)
		if checkErr != nil {
			logger.Error("failed to check recommendation dedup", zap.Error(checkErr), zap.Int64("user_id", user.ID))
			continue
		}
		if notified {
			continue
		}

		sendErr := m.tele.SendProductRecommendation(*user.TelegramUserID, *productDTO, storeName)
		if sendErr != nil {
			if isTelegramBlockedError(sendErr) {
				logger.Warn("skipping recommendation: user cannot receive bot messages",
					zap.Int64("user_id", user.ID),
					zap.Int64("telegram_user_id", *user.TelegramUserID),
				)
			} else {
				logger.Error("failed to send product recommendation",
					zap.Error(sendErr),
					zap.Int64("user_id", user.ID),
					zap.Int64("product_id", product.ID),
				)
			}
			continue
		}

		if recordErr := m.recommendationStorage.RecordNotification(ctx, user.ID, product.ID); recordErr != nil {
			logger.Error("failed to record product recommendation", zap.Error(recordErr), zap.Int64("user_id", user.ID))
		}

		time.Sleep(recommendationSendDelay)
	}
}

func (m *recommendationModule) mapProductToDTO(p *db.Product) *dto.Product {
	var images []string
	if p.Images != "" {
		_ = json.Unmarshal([]byte(p.Images), &images)
	}
	if images == nil {
		images = []string{}
	}

	return &dto.Product{
		ID:          p.ID,
		SellerID:    p.SellerID,
		StoreID:     p.StoreID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		Images:      images,
		Status:      p.Status,
		IsPosted:    p.IsPosted,
		IsBoosted:   p.IsBoosted,
	}
}

func isTelegramBlockedError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "403") ||
		strings.Contains(msg, "blocked") ||
		strings.Contains(msg, "bot was blocked") ||
		strings.Contains(msg, "user is deactivated") ||
		strings.Contains(msg, "chat not found")
}
