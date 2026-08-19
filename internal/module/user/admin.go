package user

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
)

func (m *userModule) ListUsers(ctx context.Context, filter dto.AdminUserFilterParams) (*dto.PaginatedResponse, error) {
	users, total, err := m.userStorage.ListUsers(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]dto.User, len(users))
	for i := range users {
		email := ""
		if users[i].Email != nil {
			email = *users[i].Email
		}
		items[i] = dto.User{
			ID:             users[i].ID,
			TelegramUserID: users[i].TelegramUserID,
			Email:          email,
			Username:       users[i].Username,
			Role:           users[i].Role,
		}
	}
	return dto.NewPaginatedResponse(items, total, filter.PaginationParams), nil
}
