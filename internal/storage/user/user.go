package user

import (
	"context"
	"log"

	"github.com/BoruTamena/gabaa-bot/internal/constant/errors"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/constant/persistencedb"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
)

type userStorage struct {
	db persistencedb.PersistenceDb
}

func InitUserStorage(db persistencedb.PersistenceDb) storage.UserStorage {
	return userStorage{
		db: db,
	}
}
func (u userStorage) CreateUser(ctx context.Context, userDto dto.User) error {
	// implement the logic to create a user

	userModel := db.User{
		TelID:     userDto.TelId,
		Username:  userDto.Username,
		FirstName: userDto.FirstName,
		LastName:  userDto.LastName,
	}

	res := u.db.WithContext(ctx).Create(&userModel)

	if err := res.Error; err != nil {

		err := errors.WriteErr.Wrap(err, "can't register you a seller")

		log.Println(" unable to register user on the system ::", err)
		return err

	}

	return nil
}
