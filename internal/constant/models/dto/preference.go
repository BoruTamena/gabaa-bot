package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const MaxUserPreferenceCategories = 10

type UserPreferences struct {
	Enabled    bool     `json:"enabled"`
	Categories []string `json:"categories"`
}

type UpdateUserPreferencesRequest struct {
	Enabled    bool     `json:"enabled"`
	Categories []string `json:"categories"`
}

func (r UpdateUserPreferencesRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Categories, validation.Length(0, MaxUserPreferenceCategories)),
	)
}
