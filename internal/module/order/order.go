package order

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"go.uber.org/zap"
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

	// check if product is in stock
	product, err := m.productStorage.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}
	if product.Stock < quantity {
		return fmt.Errorf("insufficient stock for product %d", productID)
	}

	err = m.cartStorage.AddToCart(ctx, userID, productID, quantity)
	if err != nil {
		logger.Error("failed to add item to cart", zap.Error(err), zap.Int64("user_id", userID), zap.Int64("product_id", productID))
	} else {
		logger.Info("item added to cart successfully", zap.Int64("user_id", userID), zap.Int64("product_id", productID), zap.Int("quantity", quantity))
	}
	return err
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

func (m *orderModule) GetUserCart(ctx context.Context, userID int64) (*dto.CartResponse, error) {
	cart, err := m.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := &dto.CartResponse{
		Items:      make([]dto.CartItem, 0, len(cart)),
		TotalPrice: 0,
	}

	for pID, qty := range cart {
		product, err := m.productStorage.GetProductByID(ctx, pID)
		if err != nil {
			// If product not found, we skip it or handle error
			continue
		}

		var images []string
		if product.Images != "" {
			_ = json.Unmarshal([]byte(product.Images), &images)
		}
		if images == nil {
			images = []string{}
		}

		res.Items = append(res.Items, dto.CartItem{
			Product: dto.Product{
				ID:          product.ID,
				StoreID:     product.StoreID,
				SellerID:    product.SellerID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Stock:       product.Stock,
				Category:    product.Category,
				Images:      images,
				IsPosted:    product.IsPosted,
				IsBoosted:   product.IsBoosted,
			},
			Quantity: qty,
		})
		res.TotalPrice += product.Price * float64(qty)
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
		logger.Error("failed to create order", zap.Error(err), zap.Int64("user_id", userID))
		return nil, err
	}

	if err := m.cartStorage.ClearCart(ctx, userID); err != nil {
		logger.Error("failed to clear cart after checkout", zap.Error(err), zap.Int64("user_id", userID))
		return nil, err
	}

	logger.Info("checkout completed successfully", zap.Int64("order_id", order.ID), zap.Int64("user_id", userID), zap.Float64("total_price", order.TotalPrice))

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
			logger.Error("failed to get order for status update", zap.Error(err), zap.Int64("order_id", orderID))
			return err
		}
		if err := m.walletStorage.UpdateWalletBalance(ctx, order.StoreID, order.TotalPrice); err != nil {
			logger.Error("failed to update wallet balance on order completion", zap.Error(err), zap.Int64("store_id", order.StoreID))
			return err
		}
	}
	err := m.orderStorage.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		logger.Error("failed to update order status", zap.Error(err), zap.Int64("order_id", orderID), zap.String("status", status))
	} else {
		logger.Info("order status updated successfully", zap.Int64("order_id", orderID), zap.String("status", status))
	}
	return err
}

func (m *orderModule) GetUserOrders(ctx context.Context, userID int64, params dto.PaginationParams) (*dto.PaginatedResponse, error) {
	limit := params.GetLimit()
	offset := params.GetOffset()

	orders, err := m.orderStorage.GetOrdersByCustomerID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := m.orderStorage.GetOrdersTotalByUserID(ctx, userID)
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
