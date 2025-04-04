package user

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"gopkg.in/telebot.v4"
)

type userHandler struct {
	userModule module.UserModule
}

func InitUserHandler(uModule module.UserModule) *userHandler {
	return &userHandler{
		userModule: uModule,
	}
}

func (user *userHandler) CreateUser(c telebot.Context) error {

	// getting user info from the context
	userID := c.Sender().ID
	userName := c.Sender().Username
	userFirstName := c.Sender().FirstName
	userLastName := c.Sender().LastName

	// creating user object
	userObj := dto.User{
		TelId:     userID,
		Username:  userName,
		FirstName: userFirstName,
		LastName:  userLastName,
	}

	if err := user.userModule.CreateUser(context.Background(), userObj); err != nil {

		return c.Respond(&telebot.CallbackResponse{
			Text:      "can't register you this time",
			ShowAlert: true,
		})
	}

	return c.Respond(&telebot.CallbackResponse{
		Text:      "Congrant !!!  \n you are now offically become a gabaa seller :) ",
		ShowAlert: true,
	})

}
