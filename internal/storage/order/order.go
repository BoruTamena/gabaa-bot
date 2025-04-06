package order

import (
	"context"

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
		BuyerID:    order.BuyerID,
		Status:     "pending",
		OrderItems: []db.OrderItem{},
	}

	for _, items := range order.OrderItems {

		orderModel.OrderItems = append(orderModel.OrderItems, db.OrderItem{
			ProductId: items.ProductId,
			Quantity:  items.Quantity,
			Price:     items.Price,
		})

	}

	res := s.db.WithContext(ctx).Create(&orderModel)

	if err := res.Error; err != nil {

		return err, uuid.Nil
	}

	return nil, orderModel.ID
}
