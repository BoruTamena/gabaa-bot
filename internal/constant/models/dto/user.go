package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type User struct {
	TelId     int64  `json:"tel_id,omitempty"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// Phone     string `json:"phone"`
	// Email     string `json:"email"`
	UserRole string `json:"user_role,omitempty"`
}
type UserList struct {
	Users []User `json:"users"`
}

func (user User) Validate() error {
	return validation.ValidateStruct(&user,
		validation.Field(&user.Username, validation.Required),
		validation.Field(&user.FirstName, validation.Required),
		validation.Field(&user.LastName, validation.Required),
		// validation.Field(&user.Phone, validation.Required),
	)
}
