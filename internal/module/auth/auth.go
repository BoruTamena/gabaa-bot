package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/internal/storage/persistence"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const loginSessionTTL = 5 * time.Minute

type authModule struct {
	userStorage        storage.UserStorage
	storeStorage       storage.StoreStorage
	deliveryStorage    storage.DeliveryStorage
	authSessionStorage storage.AuthSessionStorage
	tele               platform.Telegram
}

func NewAuthModule(
	uStorage storage.UserStorage,
	sStorage storage.StoreStorage,
	dStorage storage.DeliveryStorage,
	tele platform.Telegram,
	authSessionStorage storage.AuthSessionStorage,
) module.AuthModule {
	return &authModule{
		userStorage:        uStorage,
		storeStorage:       sStorage,
		deliveryStorage:    dStorage,
		authSessionStorage: authSessionStorage,
		tele:               tele,
	}
}

func (m *authModule) TelegramAuth(ctx context.Context, initData string) (*dto.AuthResponse, error) {
	logger.Info("telegram auth", zap.String("init_data", initData))

	valid, err := m.tele.ValidateInitData(initData)
	if err != nil || !valid {
		logger.Error("invalid telegram init data", zap.Error(err))
		return nil, fmt.Errorf("invalid telegram init data")
	}

	tgUser, chatID, err := m.tele.ParseInitData(initData)
	if err != nil {
		logger.Error("failed to parse init data", zap.Error(err))
		return nil, err
	}

	return m.authenticateTelegramUser(ctx, tgUser, chatID)
}

func (m *authModule) StartBotLoginSession(ctx context.Context) (*dto.TelegramLoginSessionResponse, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		logger.Error("failed to generate login session id", zap.Error(err))
		return nil, err
	}

	expiresAt := time.Now().Add(loginSessionTTL)
	if err := m.authSessionStorage.CreateSession(ctx, sessionID, expiresAt); err != nil {
		return nil, err
	}

	return &dto.TelegramLoginSessionResponse{
		SessionID: sessionID,
		BotURL:    m.botLoginURL(sessionID),
		ExpiresAt: expiresAt,
	}, nil
}

func (m *authModule) CompleteBotLoginSession(ctx context.Context, sessionID string, tgUser *dto.TelegramUser) error {
	if tgUser == nil {
		return persistence.ErrAuthSessionNotFound
	}

	username := tgUser.Username
	if username == "" {
		username = tgUser.FirstName
	}

	if err := m.authSessionStorage.CompleteSession(ctx, sessionID, tgUser.ID, username); err != nil {
		return err
	}

	logger.Info("telegram bot login session completed",
		zap.String("session_id", sessionID),
		zap.Int64("telegram_user_id", tgUser.ID),
	)
	return nil
}

func (m *authModule) PollBotLoginSession(ctx context.Context, sessionID string) (*dto.TelegramLoginPollResponse, error) {
	session, err := m.authSessionStorage.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session.Status != db.TelegramLoginSessionStatusCompleted {
		return &dto.TelegramLoginPollResponse{Status: db.TelegramLoginSessionStatusPending}, nil
	}

	tgUser := &dto.TelegramUser{
		ID:       session.TelegramUserID,
		Username: session.Username,
	}

	authResp, err := m.authenticateTelegramUser(ctx, tgUser, 0)
	if err != nil {
		return nil, err
	}

	if err := m.authSessionStorage.DeleteSession(ctx, sessionID); err != nil {
		logger.Error("failed to delete login session after poll", zap.Error(err), zap.String("session_id", sessionID))
	}

	return dto.NewTelegramLoginPollResponseFromAuth(db.TelegramLoginSessionStatusCompleted, authResp), nil
}

func (m *authModule) authenticateTelegramUser(ctx context.Context, tgUser *dto.TelegramUser, chatID int64) (*dto.AuthResponse, error) {
	user, err := m.userStorage.GetUserByTelegramID(ctx, tgUser.ID)
	if err != nil {
		logger.Info("user not found, creating new user", zap.Int64("telegram_id", tgUser.ID))
		user = &db.User{
			TelegramUserID: &tgUser.ID,
			Username:       tgUser.Username,
			Role:           "customer",
			BotStarted:     true,
		}
		if err := m.userStorage.CreateUser(ctx, user); err != nil {
			logger.Error("failed to create user", zap.Error(err), zap.Int64("telegram_id", tgUser.ID))
			return nil, err
		}
		logger.Info("user created successfully", zap.Int64("user_id", user.ID))
	} else if !user.BotStarted {
		user.BotStarted = true
		if tgUser.Username != "" && user.Username == "" {
			user.Username = tgUser.Username
		}
		if err := m.userStorage.UpdateUser(ctx, user); err != nil {
			logger.Error("failed to mark bot started for user", zap.Error(err), zap.Int64("user_id", user.ID))
		}
	}

	role := "customer"
	var deliveryAgentID int64
	if chatID != 0 {
		isAdmin, _ := m.tele.IsChatAdmin(chatID, tgUser.ID)
		if isAdmin {
			role = "admin"
		}
	} else {
		stores, err := m.storeStorage.GetStoresBySellerID(ctx, user.ID)
		if err == nil && len(stores) > 0 {
			role = "admin"
		}
	}

	if role != "admin" {
		agent, err := m.deliveryStorage.GetAgentByUserID(ctx, user.ID)
		if err != nil {
			agent, err = m.deliveryStorage.GetAgentByTelegramUserID(ctx, tgUser.ID)
		}
		if err == nil && agent != nil && agent.Status == constant.DeliveryAgentStatusActive {
			role = constant.RoleDelivery
			deliveryAgentID = agent.ID
		}
	}

	var storeID int64
	hasStore := false
	if role == "admin" {
		if chatID != 0 {
			store, err := m.storeStorage.GetStoreByChatID(ctx, chatID)
			if err == nil {
				hasStore = true
				storeID = store.ID
			}
		} else {
			stores, err := m.storeStorage.GetStoresBySellerID(ctx, user.ID)
			if err == nil && len(stores) > 0 {
				hasStore = true
				storeID = stores[0].ID
			}
		}
	}

	token, err := m.generateJWT(user.ID, role, storeID, deliveryAgentID)
	if err != nil {
		logger.Error("failed to generate JWT", zap.Error(err), zap.Int64("user_id", user.ID))
		return nil, err
	}

	logger.Info("user authenticated successfully", zap.Int64("user_id", user.ID), zap.String("role", role), zap.Int64("store_id", storeID))

	isDelivery := role == constant.RoleDelivery
	return &dto.AuthResponse{
		Token:           token,
		UserID:          user.ID,
		TelegramUserID:  tgUser.ID,
		Username:        user.Username,
		Role:            role,
		HasStore:        hasStore,
		StoreID:         storeID,
		DeliveryAgentID: deliveryAgentID,
		IsDelivery:      isDelivery,
	}, nil
}

func (m *authModule) botLoginURL(sessionID string) string {
	username := viper.GetString("tg.bot_username")
	if username == "" {
		username = "gabaaBot"
	}

	bot := m.tele.GetBot()
	if bot != nil && bot.Me != nil && bot.Me.Username != "" {
		username = bot.Me.Username
	}

	return fmt.Sprintf("https://t.me/%s?start=login_%s", username, sessionID)
}

func generateSessionID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (m *authModule) generateJWT(userID int64, role string, storeID, deliveryAgentID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"role":     role,
		"store_id": storeID,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	if deliveryAgentID > 0 {
		claims["delivery_agent_id"] = deliveryAgentID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(viper.GetString("jwt.secret")))
}

func IsSessionNotFound(err error) bool {
	return errors.Is(err, persistence.ErrAuthSessionNotFound)
}
