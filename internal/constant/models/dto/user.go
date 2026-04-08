package dto


type User struct {
	ID             int64  `json:"id"`
	TelegramUserID int64  `json:"telegram_user_id"`
	Username       string `json:"username"`
	Role           string `json:"role"`
}

type UserList struct {
	Users []User `json:"users"`
}

