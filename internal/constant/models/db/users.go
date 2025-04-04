package db

type User struct {
	BaseModel
	TelID     int64    `json:" tel_id"`
	Username  string   `json:"username"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	UserRole  UserRole `json:"user_role" gorm:"references:ID"`
}

type UserRole struct {
	BaseModel
	Name string `json:"name"`
}
