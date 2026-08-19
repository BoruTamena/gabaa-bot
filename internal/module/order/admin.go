package order

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
)

func (m *orderModule) AdminListOrders(ctx context.Context, filter dto.AdminOrderFilterParams) (*dto.PaginatedResponse, error) {
	orders, total, err := m.orderStorage.GetOrdersByAdminFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]dto.Order, len(orders))
	for i := range orders {
		items[i] = *m.mapOrderToDTO(&orders[i])
	}
	return dto.NewPaginatedResponse(items, total, filter.PaginationParams), nil
}
