package dto


type User struct {
	ID             int64  `json:"id"`
	TelegramUserID int64  `json:"telegram_user_id"`
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
	Token    string `json:"token"`
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	HasStore bool   `json:"hasStore"`
	StoreID  int64  `json:"storeId,omitempty"`
}

type UserList struct {
	Users []User `json:"users"`
}

