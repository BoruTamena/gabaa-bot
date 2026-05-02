package cart

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"go.uber.org/zap"
)

type cartModule struct {
	cartStorage    storage.CartStorage
	productStorage storage.ProductStorage
}

func NewCartModule(cStorage storage.CartStorage, pStorage storage.ProductStorage) module.CartModule {
	return &cartModule{
		cartStorage:    cStorage,
		productStorage: pStorage,
	}
}

func (m *cartModule) AddToCart(ctx context.Context, userID int64, productID int64, quantity int) error {
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

func (m *cartModule) GetCart(ctx context.Context, userID int64) (map[int64]int, error) {
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

func (m *cartModule) GetUserCart(ctx context.Context, userID int64) (*dto.CartResponse, error) {
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

func (m *cartModule) UpdateCartItem(ctx context.Context, userID int64, productID int64, action string) error {
	cart, err := m.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	currentQty, exists := cart[productID]
	if !exists {
		return fmt.Errorf("item not found in cart")
	}

	newQty := currentQty
	if action == "increment" {
		newQty++
	} else if action == "decrement" {
		newQty--
	} else {
		return fmt.Errorf("invalid action: must be increment or decrement")
	}

	if newQty <= 0 {
		return m.RemoveFromCart(ctx, userID, productID)
	}

	product, err := m.productStorage.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}
	if product.Stock < newQty {
		return fmt.Errorf("insufficient stock for product %d", productID)
	}

	err = m.cartStorage.UpdateCartItem(ctx, userID, productID, newQty)
	if err != nil {
		logger.Error("failed to update cart item", zap.Error(err), zap.Int64("user_id", userID), zap.Int64("product_id", productID))
	}
	return err
}

func (m *cartModule) RemoveFromCart(ctx context.Context, userID int64, productID int64) error {
	err := m.cartStorage.RemoveFromCart(ctx, userID, productID)
	if err != nil {
		logger.Error("failed to remove item from cart", zap.Error(err), zap.Int64("user_id", userID), zap.Int64("product_id", productID))
	}
	return err
}

func (m *cartModule) ClearCart(ctx context.Context, userID int64) error {
	err := m.cartStorage.ClearCart(ctx, userID)
	if err != nil {
		logger.Error("failed to clear cart", zap.Error(err), zap.Int64("user_id", userID))
	}
	return err
}
