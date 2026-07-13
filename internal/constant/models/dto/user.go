package dto

import "time"

type User struct {
	ID             int64  `json:"id"`
	TelegramUserID *int64 `json:"telegram_user_id"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	Role           string `json:"role"`
}

type TelegramUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type AuthResponse struct {
	Token           string `json:"token"`
	UserID          int64  `json:"userId"`
	TelegramUserID  int64  `json:"telegramUserId,omitempty"`
	Username        string `json:"username"`
	Role            string `json:"role"`
	HasStore        bool   `json:"hasStore"`
	StoreID         int64  `json:"storeId,omitempty"`
	DeliveryAgentID int64  `json:"deliveryAgentId,omitempty"`
	IsDelivery      bool   `json:"isDelivery"`
}

type UserList struct {
	Users []User `json:"users"`
}

type TelegramLoginSessionResponse struct {
	SessionID string    `json:"sessionId"`
	BotURL    string    `json:"botUrl"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type TelegramLoginPollResponse struct {
	Status         string `json:"status"`
	Token          string `json:"token,omitempty"`
	UserID         int64  `json:"userId,omitempty"`
	TelegramUserID int64  `json:"telegramUserId,omitempty"`
	Username       string `json:"username,omitempty"`
	Role           string `json:"role,omitempty"`
	HasStore       bool   `json:"hasStore,omitempty"`
	StoreID        int64  `json:"storeId,omitempty"`
}

func NewTelegramLoginPollResponseFromAuth(status string, auth *AuthResponse) *TelegramLoginPollResponse {
	if auth == nil {
		return &TelegramLoginPollResponse{Status: status}
	}
	return &TelegramLoginPollResponse{
		Status:         status,
		Token:          auth.Token,
		UserID:         auth.UserID,
		TelegramUserID: auth.TelegramUserID,
		Username:       auth.Username,
		Role:           auth.Role,
		HasStore:       auth.HasStore,
		StoreID:        auth.StoreID,
	}
}

