package order

import (
	"context"
	"fmt"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
)

type orderModule struct {
	orderStorage   storage.OrderStorage
	productStorage storage.ProductStorage
	cartStorage    storage.CartStorage
	walletStorage  storage.WalletStorage
}

func NewOrderModule(oStorage storage.OrderStorage, pStorage storage.ProductStorage, cStorage storage.CartStorage, wStorage storage.WalletStorage) module.OrderModule {
	return &orderModule{
		orderStorage:   oStorage,
		productStorage: pStorage,
		cartStorage:    cStorage,
		walletStorage:  wStorage,
	}
}

func (m *orderModule) AddToCart(ctx context.Context, userID int64, productID int64, quantity int) error {
	return m.cartStorage.AddToCart(ctx, userID, productID, quantity)
}

func (m *orderModule) GetCart(ctx context.Context, userID int64) (map[int64]int, error) {
	cart, err := m.cartStorage.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]int)
	for k, v := range cart {
		pID, _ := strconv.ParseInt(k[2:], 10, 64) // k is "p:ID"
		res[pID] = v
	}
	return res, nil
}

func (m *orderModule) Checkout(ctx context.Context, userID int64, storeID int64) (*dto.Order, error) {
	cart, err := m.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(cart) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	order := &db.Order{
		UserID:     userID,
		StoreID:    storeID,
		Status:     "pending",
		TotalPrice: 0,
	}

	for pID, qty := range cart {
		product, err := m.productStorage.GetProductByID(ctx, pID)
		if err != nil {
			return nil, err
		}
		if product.Stock < qty {
			return nil, fmt.Errorf("insufficient stock for product %d", pID)
		}
		order.TotalPrice += product.Price * float64(qty)
		order.Items = append(order.Items, db.OrderItem{
			ProductID: pID,
			Quantity:  qty,
			Price:     product.Price,
		})

		// Decrement stock
		product.Stock -= qty
		if err := m.productStorage.UpdateProduct(ctx, product); err != nil {
			return nil, err
		}
	}

	if err := m.orderStorage.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	if err := m.cartStorage.ClearCart(ctx, userID); err != nil {
		return nil, err
	}

	return m.mapOrderToDTO(order), nil
}

func (m *orderModule) ListOrders(ctx context.Context, storeID int64, params dto.PaginationParams) (*dto.PaginatedResponse, error) {
	limit := params.GetLimit()
	offset := params.GetOffset()

	orders, err := m.orderStorage.GetOrdersByStoreID(ctx, storeID, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := m.orderStorage.GetOrdersTotalByStoreID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	dtoOrders := make([]dto.Order, len(orders))
	for i, o := range orders {
		dtoOrders[i] = *m.mapOrderToDTO(&o)
	}

	return &dto.PaginatedResponse{
		Total: total,
		Data:  dtoOrders,
	}, nil
}

func (m *orderModule) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	if status == "completed" {
		order, err := m.orderStorage.GetOrderByID(ctx, orderID)
		if err != nil {
			return err
		}
		if err := m.walletStorage.UpdateWalletBalance(ctx, order.StoreID, order.TotalPrice); err != nil {
			return err
		}
	}
	return m.orderStorage.UpdateOrderStatus(ctx, orderID, status)
}

func (m *orderModule) mapOrderToDTO(o *db.Order) *dto.Order {
	items := make([]dto.OrderItem, len(o.Items))
	for i, item := range o.Items {
		items[i] = dto.OrderItem{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return &dto.Order{
		ID:         o.ID,
		StoreID:    o.StoreID,
		UserID:     o.UserID,
		Status:     o.Status,
		TotalPrice: o.TotalPrice,
		CreatedAt:  o.CreatedAt,
		OrderItems: items,
	}
}
