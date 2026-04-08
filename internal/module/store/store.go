package store

import (
	"context"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
)

type storeModule struct {
	storeStorage storage.StoreStorage
	tele         platform.Telegram
}

func NewStoreModule(sStorage storage.StoreStorage, tele platform.Telegram) module.StoreModule {
	return &storeModule{
		storeStorage: sStorage,
		tele:         tele,
	}
}

func (m *storeModule) CreateStore(ctx context.Context, userID int64, chatID int64, chatType string, name string) (*dto.Store, error) {
	store := &db.Store{
		SellerID: userID,
		ChatID:   chatID,
		ChatType: chatType,
		Name:     name,
	}

	if err := m.storeStorage.CreateStore(ctx, store); err != nil {
		return nil, err
	}

	return &dto.Store{
		ID:       store.ID,
		SellerID: store.SellerID,
		ChatID:   store.ChatID,
		ChatType: store.ChatType,
		Name:     store.Name,
	}, nil
}

func (m *storeModule) GetAdminDashboard(ctx context.Context, userID int64, chatID int64) (string, *dto.Store, error) {
	store, err := m.storeStorage.GetStoreByChatID(ctx, chatID)
	if err != nil {
		// No store associated with this chat yet
		isAdmin, _ := m.tele.IsChatAdmin(chatID, userID)
		if isAdmin {
			return "setup", nil, nil
		}
		return "storefront", nil, nil
	}

	isAdmin, _ := m.tele.IsChatAdmin(chatID, userID)
	if isAdmin {
		return "manage", &dto.Store{
			ID:       store.ID,
			SellerID: store.SellerID,
			ChatID:   store.ChatID,
			ChatType: store.ChatType,
			Name:     store.Name,
		}, nil
	}

	return "storefront", &dto.Store{
		ID:       store.ID,
		SellerID: store.SellerID,
		ChatID:   store.ChatID,
		ChatType: store.ChatType,
		Name:     store.Name,
	}, nil
}

