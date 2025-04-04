package order

import (
	"context"
	"log"

	"github.com/BoruTamena/gabaa-bot/internal/constant/errors"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/constant/persistencedb"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/google/uuid"
)

type orderStorage struct {
	db persistencedb.PersistenceDb
}

func InitOrderStorage(db persistencedb.PersistenceDb) storage.OrderStorage {

	return &orderStorage{
		db: db,
	}
}

func (s *orderStorage) CreateOrder(ctx context.Context, order dto.Order) (error, uuid.UUID) {

	orderModel := db.Order{
		UserID:    order.UserID,
		ProductID: order.ProductID,
		Quantity:  order.Quantity,
		Status:    order.Status,
	}
	res := s.db.WithContext(ctx).Create(&orderModel)
	if err := res.Error; err != nil {
		err := errors.WriteErr.Wrap(err, "can't create order")
		log.Println("can't create order ::", err)
		return err, uuid.Nil
	}

	return nil, orderModel.ID
}
