package initiator

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
)

func seedDeliveryData(ctx context.Context, p Persistence, store *db.Store, seller *db.User) {
	if store == nil || store.ID == 0 || seller == nil {
		return
	}

	existing, _ := p.DeliveryStorage.ListAgentsForStore(ctx, store.ID)
	if len(existing) > 0 {
		return
	}

	if store.Location == "" {
		store.Location = "Bole Road, Addis Ababa, Bole"
		_ = p.StoreStorage.UpdateStore(ctx, store)
	}

	courierTelegramID := int64(987654321)
	courierUser, err := p.UserStorage.GetUserByTelegramID(ctx, courierTelegramID)
	if err != nil {
		courierUser = &db.User{
			TelegramUserID: &courierTelegramID,
			Username:       "courier_john",
			Role:           "customer",
			BotStarted:     true,
		}
		_ = p.UserStorage.CreateUser(ctx, courierUser)
		courierUser, _ = p.UserStorage.GetUserByTelegramID(ctx, courierTelegramID)
	}

	agentJohn := &db.DeliveryAgent{
		Username:       "courier_john",
		FullName:       "John Doe",
		Phone:          "251911234567",
		UserID:         &courierUser.ID,
		TelegramUserID: &courierTelegramID,
		LoyaltyScore:   3,
		Status:         constant.DeliveryAgentStatusActive,
	}
	if err := p.DeliveryStorage.CreateAgent(ctx, agentJohn); err != nil {
		return
	}

	agentSara := &db.DeliveryAgent{
		Username:     "courier_sara",
		FullName:     "Sara Bekele",
		Phone:        "251922334455",
		LoyaltyScore: 1,
		Status:       constant.DeliveryAgentStatusPendingInvite,
	}
	if err := p.DeliveryStorage.CreateAgent(ctx, agentSara); err != nil {
		return
	}

	linkJohn := &db.StoreDeliveryLink{
		StoreID:           store.ID,
		DeliveryAgentID:   agentJohn.ID,
		ShareEnabled:      true,
		ConnectedByUserID: &seller.ID,
	}
	if err := p.DeliveryStorage.CreateStoreLink(ctx, linkJohn); err != nil {
		return
	}

	linkSara := &db.StoreDeliveryLink{
		StoreID:           store.ID,
		DeliveryAgentID:   agentSara.ID,
		ShareEnabled:      false,
		ConnectedByUserID: &seller.ID,
	}
	if err := p.DeliveryStorage.CreateStoreLink(ctx, linkSara); err != nil {
		return
	}

	seedJohnRoutes(ctx, p, linkJohn.ID)
	seedSaraRoutes(ctx, p, linkSara.ID)
}

func seedJohnRoutes(ctx context.Context, p Persistence, linkID int64) {
	route := &db.DeliveryRoute{
		StoreDeliveryLinkID: linkID,
		Label:               "Bole + Atlas Route",
		IsActive:            true,
	}
	if err := p.DeliveryStorage.CreateRoute(ctx, route); err != nil {
		return
	}

	pickups := []db.DeliveryRouteLocation{
		{
			DeliveryRouteID:  route.ID,
			LocationType:     constant.DeliveryLocationTypePickup,
			Label:            "My store",
			UseStoreLocation: true,
		},
		{
			DeliveryRouteID: route.ID,
			LocationType:    constant.DeliveryLocationTypePickup,
			Label:           "Warehouse",
			Country:         "Ethiopia",
			Region:          "Bole",
			City:            "Addis Ababa",
			Street:          "Bole Road",
			Landmark:        "Near Edna Mall",
		},
	}
	deliveries := []db.DeliveryRouteLocation{
		{
			DeliveryRouteID: route.ID,
			LocationType:    constant.DeliveryLocationTypeDelivery,
			Label:           "Bole area",
			Country:         "Ethiopia",
			Region:          "Bole",
			City:            "Addis Ababa",
		},
		{
			DeliveryRouteID: route.ID,
			LocationType:    constant.DeliveryLocationTypeDelivery,
			Label:           "Atlas street",
			Country:         "Ethiopia",
			Region:          "Bole",
			City:            "Addis Ababa",
			Street:          "Atlas Avenue",
		},
		{
			DeliveryRouteID: route.ID,
			LocationType:    constant.DeliveryLocationTypeDelivery,
			Label:           "CMC",
			Country:         "Ethiopia",
			Region:          "CMC",
			City:            "Addis Ababa",
			Street:          "CMC Road",
		},
	}
	for i := range pickups {
		_ = p.DeliveryStorage.CreateRouteLocation(ctx, &pickups[i])
	}
	for i := range deliveries {
		_ = p.DeliveryStorage.CreateRouteLocation(ctx, &deliveries[i])
	}
}

func seedSaraRoutes(ctx context.Context, p Persistence, linkID int64) {
	route := &db.DeliveryRoute{
		StoreDeliveryLinkID: linkID,
		Label:               "CMC + Megenagna Route",
		IsActive:            true,
	}
	if err := p.DeliveryStorage.CreateRoute(ctx, route); err != nil {
		return
	}

	pickups := []db.DeliveryRouteLocation{
		{
			DeliveryRouteID:  route.ID,
			LocationType:     constant.DeliveryLocationTypePickup,
			Label:            "Main shop",
			UseStoreLocation: true,
		},
	}
	deliveries := []db.DeliveryRouteLocation{
		{
			DeliveryRouteID: route.ID,
			LocationType:    constant.DeliveryLocationTypeDelivery,
			Label:           "CMC area",
			Country:         "Ethiopia",
			Region:          "CMC",
			City:            "Addis Ababa",
		},
		{
			DeliveryRouteID: route.ID,
			LocationType:    constant.DeliveryLocationTypeDelivery,
			Label:           "Megenagna",
			Country:         "Ethiopia",
			Region:          "Megenagna",
			City:            "Addis Ababa",
			Street:          "Megenagna Road",
			Landmark:        "Near Bole Medhanialem",
		},
		{
			DeliveryRouteID: route.ID,
			LocationType:    constant.DeliveryLocationTypeDelivery,
			Label:           "Kazanchis",
			Country:         "Ethiopia",
			Region:          "Kazanchis",
			City:            "Addis Ababa",
		},
	}
	for i := range pickups {
		_ = p.DeliveryStorage.CreateRouteLocation(ctx, &pickups[i])
	}
	for i := range deliveries {
		_ = p.DeliveryStorage.CreateRouteLocation(ctx, &deliveries[i])
	}
}
