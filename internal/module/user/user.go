package user

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
)

type userModule struct {
	userStorage storage.UserStorage
}

func NewUserModule(uStorage storage.UserStorage) module.UserModule {
	return &userModule{
		userStorage: uStorage,
	}
}

func (m *userModule) GetOrCreateUser(ctx context.Context, telegramID int64, username string) (*dto.User, error) {
	user, err := m.userStorage.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		// If not found, create
		user = &db.User{
			TelegramUserID: &telegramID,
			Username:       username,
			Role:           "customer",
		}
		if err := m.userStorage.CreateUser(ctx, user); err != nil {
			return nil, err
		}
	}

	return &dto.User{
		ID:             user.ID,
		TelegramUserID: user.TelegramUserID,
		Username:       user.Username,
		Role:           user.Role,
	}, nil
}
