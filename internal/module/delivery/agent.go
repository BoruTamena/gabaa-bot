package delivery

import (
	"context"
	"fmt"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
)

func (m *deliveryModule) GetProfile(ctx context.Context, agentID int64) (*dto.DeliveryProfileResponse, error) {
	agent, err := m.deliveryStorage.GetAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}

	var allRoutes []dto.DeliveryRouteResponse
	storeLinks, err := m.deliveryStorage.ListLinksByAgentID(ctx, agentID)
	if err != nil {
		return nil, err
	}
	for _, link := range storeLinks {
		for _, r := range link.Routes {
			allRoutes = append(allRoutes, mapRouteToResponse(r))
		}
	}

	return &dto.DeliveryProfileResponse{
		ID:           agent.ID,
		Username:     agent.Username,
		FullName:     agent.FullName,
		Phone:        agent.Phone,
		LoyaltyScore: agent.LoyaltyScore,
		Status:       string(agent.Status),
		Routes:       allRoutes,
	}, nil
}

func (m *deliveryModule) ListAssignedOrders(ctx context.Context, agentID int64, status string, params dto.PaginationParams) (*dto.PaginatedResponse, error) {
	orders, err := m.orderStorage.GetOrdersByDeliveryAgentID(ctx, agentID, status, params.GetLimit(), params.GetOffset())
	if err != nil {
		return nil, err
	}
	total, err := m.orderStorage.GetOrdersTotalByDeliveryAgentID(ctx, agentID, status)
	if err != nil {
		return nil, err
	}
	items := make([]dto.DeliveryOrderResponse, len(orders))
	for i, o := range orders {
		items[i] = m.mapDeliveryOrder(ctx, &o)
	}
	return &dto.PaginatedResponse{
		Total: total,
		Data:  items,
	}, nil
}

func (m *deliveryModule) GetAssignedOrder(ctx context.Context, agentID, orderID int64) (*dto.DeliveryOrderResponse, error) {
	order, err := m.orderStorage.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.DeliveryAgentID == nil || *order.DeliveryAgentID != agentID {
		return nil, fmt.Errorf("order not assigned to you")
	}
	resp := m.mapDeliveryOrder(ctx, order)
	return &resp, nil
}

func (m *deliveryModule) UpdateDeliveryOrderStatus(ctx context.Context, agentID, orderID int64, status string) error {
	if status != "delivered" && status != "failed" {
		return fmt.Errorf("status must be 'delivered' or 'failed'")
	}

	order, err := m.orderStorage.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order.DeliveryAgentID == nil || *order.DeliveryAgentID != agentID {
		return fmt.Errorf("order not assigned to you")
	}
	if order.Status != "shipped" {
		return fmt.Errorf("order must be in shipped status")
	}

	if status == "delivered" {
		if err := releaseOrderDelivery(ctx, m.escrowStorage, m.walletStorage, order.StoreID, orderID); err != nil {
			return err
		}
		if err := m.orderStorage.UpdateOrderStatus(ctx, orderID, "delivered"); err != nil {
			return err
		}
		return m.deliveryStorage.AdjustLoyaltyScore(ctx, agentID, 1)
	}

	// failed
	if err := m.orderStorage.UpdateOrderStatus(ctx, orderID, "cancelled"); err != nil {
		return err
	}
	return m.deliveryStorage.AdjustLoyaltyScore(ctx, agentID, -1)
}

func releaseOrderDelivery(ctx context.Context, escrowStorage storage.EscrowStorage, walletStorage storage.WalletStorage, storeID, orderID int64) error {
	escrow, err := escrowStorage.GetEscrowByOrderID(ctx, orderID)
	if err == nil && escrow.Status == constant.EscrowStatusHeld {
		if err := escrowStorage.ReleaseEscrow(ctx, orderID); err != nil {
			return err
		}
		if err := walletStorage.ReleaseEscrowFunds(ctx, storeID, escrow.Amount); err != nil {
			return err
		}
	}
	return nil
}

func (m *deliveryModule) mapDeliveryOrder(ctx context.Context, o *db.Order) dto.DeliveryOrderResponse {
	pickup := ""
	storeName := ""
	if o.Store.ID != 0 {
		storeName = o.Store.Name
		pickup = o.Store.Location
	} else if o.StoreID != 0 {
		if store, err := m.storeStorage.GetStoreByID(ctx, o.StoreID); err == nil {
			storeName = store.Name
			pickup = store.Location
		}
	}

	if o.DeliveryRouteID != nil {
		if route, err := m.deliveryStorage.GetRouteByID(ctx, *o.DeliveryRouteID); err == nil {
			for _, loc := range route.Locations {
				if loc.LocationType == constant.DeliveryLocationTypePickup {
					if loc.UseStoreLocation {
						break
					}
					pickup = formatLocation(loc.Street, loc.City, loc.Region, loc.Landmark)
					break
				}
			}
		}
	}

	resp := dto.DeliveryOrderResponse{
		PickupLocation: pickup,
		StoreName:      storeName,
	}
	resp.ID = o.ID
	resp.StoreID = o.StoreID
	resp.UserID = o.UserID
	resp.Status = o.Status
	resp.TotalPrice = o.TotalPrice
	resp.CreatedAt = o.CreatedAt

	if o.ShippingAddress != nil {
		resp.ShippingAddress = &dto.Address{
			ID:            o.ShippingAddress.ID,
			RecipientName: o.ShippingAddress.RecipientName,
			Phone:         o.ShippingAddress.Phone,
			Street:        o.ShippingAddress.Street,
			City:          o.ShippingAddress.City,
			Region:        o.ShippingAddress.Region,
			Country:       o.ShippingAddress.Country,
		}
	}

	items := make([]dto.OrderItem, len(o.Items))
	for i, item := range o.Items {
		items[i] = dto.OrderItem{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
		if item.Product.ID != 0 {
			items[i].Product = &dto.OrderProductDetail{
				ID:   item.Product.ID,
				Name: item.Product.Name,
			}
		}
	}
	resp.OrderItems = items
	return resp
}
