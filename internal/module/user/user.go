package user

import (
	"context"
	"log"

	"github.com/BoruTamena/gabaa-bot/internal/constant/errors"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
)

type userModule struct {
	userStorage storage.UserStorage
}

func InitUserModule(uStorage storage.UserStorage) userModule {
	return userModule{
		userStorage: uStorage,
	}
}

func (u userModule) CreateUser(ctx context.Context, userDto dto.User) error {

	if err := userDto.Validate(); err != nil {

		err = errors.BadInput.Wrap(err, "bad user input")

		log.Println("can't register user :: ", err)
		return err
	}

	err := u.userStorage.CreateUser(ctx, userDto)

	if err != nil {

		return err

	}

	return nil
}
