package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type authModule struct {
	userStorage  storage.UserStorage
	storeStorage storage.StoreStorage
	tele         platform.Telegram
}

func NewAuthModule(uStorage storage.UserStorage, sStorage storage.StoreStorage, tele platform.Telegram) module.AuthModule {
	return &authModule{
		userStorage:  uStorage,
		storeStorage: sStorage,
		tele:         tele,
	}
}

func (m *authModule) TelegramAuth(ctx context.Context, initData string) (*dto.AuthResponse, error) {
	
	logger.Info("telegram auth", zap.String("init_data", initData))
	
	// 1. Validate Telegram initData
	valid, err := m.tele.ValidateInitData(initData)
	if err != nil || !valid {
		logger.Error("invalid telegram init data", zap.Error(err))
		return nil, fmt.Errorf("invalid telegram init data")
	}

	// 2. Extract user data
	tgUser, chatID, err := m.tele.ParseInitData(initData)
	if err != nil {
		logger.Error("failed to parse init data", zap.Error(err))
		return nil, err
	}

	// 3. Check if user exists -> create if not
	user, err := m.userStorage.GetUserByTelegramID(ctx, tgUser.ID)
	if err != nil {
		logger.Info("user not found, creating new user", zap.Int64("telegram_id", tgUser.ID))
		// Create user
		user = &db.User{
			TelegramUserID: tgUser.ID,
			Username:       tgUser.Username,
			Role:           "customer", // Default role
		}
		if err := m.userStorage.CreateUser(ctx, user); err != nil {
			logger.Error("failed to create user", zap.Error(err), zap.Int64("telegram_id", tgUser.ID))
			return nil, err
		}
		logger.Info("user created successfully", zap.Int64("user_id", user.ID))
	}

	// 4. Determine role
	role := "customer"
	if chatID != 0 {
		isAdmin, _ := m.tele.IsChatAdmin(chatID, tgUser.ID)
		if isAdmin {
			role = "admin"
		}
	} else {
		// Personal chat - check if user has any stores
		stores, err := m.storeStorage.GetStoresBySellerID(ctx, user.ID)
		if err == nil && len(stores) > 0 {
			role = "admin"
		}
	}

	// 5. Check store existence for admins
	var storeID int64
	hasStore := false
	if role == "admin" {
		targetChatID := chatID
		if targetChatID == 0 {
			targetChatID = tgUser.ID // Check for personal store
		}
		store, err := m.storeStorage.GetStoreByChatID(ctx, targetChatID)
		if err == nil {
			hasStore = true
			storeID = store.ID
		}
	}

	// 6. Generate JWT
	token, err := m.generateJWT(user.ID, role, storeID)
	if err != nil {
		logger.Error("failed to generate JWT", zap.Error(err), zap.Int64("user_id", user.ID))
		return nil, err
	}

	logger.Info("user authenticated successfully", zap.Int64("user_id", user.ID), zap.String("role", role), zap.Int64("store_id", storeID))

	return &dto.AuthResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Role:     role,
		HasStore: hasStore,
		StoreID:  storeID,
	}, nil
}

func (m *authModule) generateJWT(userID int64, role string, storeID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"role":     role,
		"store_id": storeID,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 1 week
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString("jwt.secret")))
}
