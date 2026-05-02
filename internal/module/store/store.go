package store

import (
	"context"
	"fmt"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/BoruTamena/gabaa-bot/platform"
	"go.uber.org/zap"
)

type storeModule struct {
	storeStorage storage.StoreStorage
	userStorage  storage.UserStorage
	tele         platform.Telegram
}

func NewStoreModule(sStorage storage.StoreStorage, uStorage storage.UserStorage, tele platform.Telegram) module.StoreModule {
	return &storeModule{
		storeStorage: sStorage,
		userStorage:  uStorage,
		tele:         tele,
	}
}

func (m *storeModule) CreateStore(ctx context.Context, userID int64, req dto.CreateStoreRequest) (*dto.Store, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Security: Verify user is admin of the target chat if provided
	if req.TelegramChatID != 0 {
		user, err := m.userStorage.GetUserByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("user not found")
		}

		if user.TelegramUserID == nil {
			return nil, fmt.Errorf("user has no telegram id linked")
		}

		isAdmin, err := m.tele.IsChatAdmin(req.TelegramChatID, *user.TelegramUserID)
		if err != nil || !isAdmin {
			logger.Warn("user attempted to create store for chat without admin rights", 
				zap.Int64("user_id", userID), zap.Int64("chat_id", req.TelegramChatID))
			return nil, fmt.Errorf("you must be an admin of the chat to create a store")
		}
	}

	store := &db.Store{
		SellerID:       userID,
		TelegramChatID: req.TelegramChatID,
		Name:           req.Name,
		Category:       req.Category,
		Description:    req.Description,
		LogoImage:      req.LogoImage,
		CoverImage:     req.CoverImage,
		Phone:          req.Phone,
		Email:          req.Email,
		Location:       req.Location,
	}

	if err := m.storeStorage.CreateStore(ctx, store); err != nil {
		logger.Error("failed to create store", zap.Error(err), zap.Int64("seller_id", userID))
		return nil, err
	}

	logger.Info("store created successfully", zap.Int64("store_id", store.ID), zap.Int64("seller_id", userID))

	return &dto.Store{
		ID:                store.ID,
		SellerID:          store.SellerID,
		TelegramChatID:    store.TelegramChatID,
		TelegramChatTitle: store.TelegramChatTitle,
		Status:            store.Status,
		Name:              store.Name,
		Category:          store.Category,
		Description:       store.Description,
		LogoImage:         store.LogoImage,
		CoverImage:        store.CoverImage,
		Phone:             store.Phone,
		Email:             store.Email,
		Location:          store.Location,
	}, nil
}

func (m *storeModule) GetAdminDashboard(ctx context.Context, userID int64, chatID int64) (string, *dto.Store, error) {
	// 1. Try to get store by specific chat ID (e.g. Group/Channel)
	store, err := m.storeStorage.GetStoreByChatID(ctx, chatID)
	
	user, userErr := m.userStorage.GetUserByID(ctx, userID)
	tgUserID := int64(0)
	if userErr == nil && user.TelegramUserID != nil {
		tgUserID = *user.TelegramUserID
	}

	// 2. If store exists for this chat, check if user is admin
	if err == nil {
		if tgUserID != 0 {
			isAdmin, _ := m.tele.IsChatAdmin(chatID, tgUserID)
			if isAdmin {
				logger.Info("dashboard: manage (group admin)", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
				return "manage", m.mapToDTO(store), nil
			}
		}
		logger.Info("dashboard: storefront", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
		return "storefront", m.mapToDTO(store), nil
	}

	// 3. No store for this chat. Check if the user already HAS any store (Private Chat fallback)
	stores, _ := m.storeStorage.GetStoresBySellerID(ctx, userID)
	if len(stores) > 0 {
		// Merchant already has at least one store. 
		// If they are in a private chat (chatID > 0), show them their management dashboard.
		if chatID > 0 {
			logger.Info("dashboard: manage (private chat fallback)", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
			return "manage", m.mapToDTO(&stores[0]), nil
		}
	}

	// 4. Truly no store found. Determine if they should see setup or storefront
	if tgUserID != 0 {
		isAdmin, _ := m.tele.IsChatAdmin(chatID, tgUserID)
		if isAdmin {
			logger.Info("dashboard: setup", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
			return "setup", nil, nil
		}
	}

	logger.Info("dashboard: storefront (no store)", zap.Int64("chat_id", chatID), zap.Int64("user_id", userID))
	return "storefront", nil, nil
}

func (m *storeModule) GetStore(ctx context.Context, id int64) (*dto.Store, error) {
	store, err := m.storeStorage.GetStoreByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return m.mapToDTO(store), nil
}

func (m *storeModule) GetStoreStatus(ctx context.Context, id int64) (string, error) {
	store, err := m.storeStorage.GetStoreByID(ctx, id)
	if err != nil {
		return "", err
	}
	return store.Status, nil
}

func (m *storeModule) UpdateStore(ctx context.Context, id int64, req dto.UpdateStoreRequest) (*dto.Store, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	store, err := m.storeStorage.GetStoreByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		store.Name = req.Name
	}
	if req.Category != "" {
		store.Category = req.Category
	}
	if req.Description != "" {
		store.Description = req.Description
	}
	if req.LogoImage != "" {
		store.LogoImage = req.LogoImage
	}
	if req.CoverImage != "" {
		store.CoverImage = req.CoverImage
	}
	if req.Phone != "" {
		store.Phone = req.Phone
	}
	if req.Email != "" {
		store.Email = req.Email
	}
	if req.Location != "" {
		store.Location = req.Location
	}

	if err := m.storeStorage.UpdateStore(ctx, store); err != nil {
		logger.Error("failed to update store", zap.Error(err), zap.Int64("store_id", id))
		return nil, err
	}

	logger.Info("store updated successfully", zap.Int64("store_id", id))

	return m.mapToDTO(store), nil
}

func (m *storeModule) mapToDTO(store *db.Store) *dto.Store {
	return &dto.Store{
		ID:                store.ID,
		SellerID:          store.SellerID,
		TelegramChatID:    store.TelegramChatID,
		TelegramChatTitle: store.TelegramChatTitle,
		Status:            store.Status,
		Name:              store.Name,
		Category:          store.Category,
		Description:       store.Description,
		LogoImage:         store.LogoImage,
		CoverImage:        store.CoverImage,
		Phone:             store.Phone,
		Email:             store.Email,
		Location:          store.Location,
	}
}
