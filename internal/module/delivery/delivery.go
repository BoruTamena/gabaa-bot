package delivery

import (
	"context"
	"fmt"
	"strings"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"github.com/BoruTamena/gabaa-bot/platform"
	"go.uber.org/zap"
)

type deliveryModule struct {
	deliveryStorage storage.DeliveryStorage
	orderStorage    storage.OrderStorage
	storeStorage    storage.StoreStorage
	escrowStorage   storage.EscrowStorage
	walletStorage   storage.WalletStorage
	tele            platform.Telegram
	orderModule     module.OrderModule
}

func NewDeliveryModule(
	dStorage storage.DeliveryStorage,
	oStorage storage.OrderStorage,
	sStorage storage.StoreStorage,
	eStorage storage.EscrowStorage,
	wStorage storage.WalletStorage,
	tele platform.Telegram,
) module.DeliveryModule {
	return &deliveryModule{
		deliveryStorage: dStorage,
		orderStorage:    oStorage,
		storeStorage:    sStorage,
		escrowStorage:   eStorage,
		walletStorage:   wStorage,
		tele:            tele,
	}
}

func (m *deliveryModule) SetOrderModule(om module.OrderModule) {
	m.orderModule = om
}

func normalizeUsername(u string) string {
	u = strings.TrimSpace(strings.TrimPrefix(u, "@"))
	return strings.ToLower(u)
}

func (m *deliveryModule) ListAreaPresets(ctx context.Context) ([]dto.DeliveryAreaPreset, error) {
	presets, err := m.deliveryStorage.ListAreaPresets(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.DeliveryAreaPreset, len(presets))
	for i, p := range presets {
		out[i] = dto.DeliveryAreaPreset{ID: p.ID, Region: p.Region, City: p.City, Street: p.Street}
	}
	return out, nil
}

func (m *deliveryModule) ConnectAgent(ctx context.Context, storeID, userID int64, req dto.ConnectDeliveryAgentRequest) (*dto.DeliveryAgentResponse, error) {
	if req.Username == "" || req.FullName == "" || req.Phone == "" {
		return nil, fmt.Errorf("username, full_name, and phone are required")
	}

	agent, err := m.deliveryStorage.GetAgentByUsername(ctx, req.Username)
	if err != nil {
		agent = &db.DeliveryAgent{
			Username:     normalizeUsername(req.Username),
			FullName:     req.FullName,
			Phone:        req.Phone,
			Status:       constant.DeliveryAgentStatusPendingInvite,
			LoyaltyScore: 0,
		}
		if err := m.deliveryStorage.CreateAgent(ctx, agent); err != nil {
			return nil, err
		}
	} else {
		agent.FullName = req.FullName
		agent.Phone = req.Phone
		if err := m.deliveryStorage.UpdateAgent(ctx, agent); err != nil {
			return nil, err
		}
	}

	existing, err := m.deliveryStorage.GetStoreLink(ctx, storeID, agent.ID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("delivery agent already connected to this store")
	}

	link := &db.StoreDeliveryLink{
		StoreID:           storeID,
		DeliveryAgentID:   agent.ID,
		ShareEnabled:      req.ShareEnabled,
		ConnectedByUserID: &userID,
	}
	if err := m.deliveryStorage.CreateStoreLink(ctx, link); err != nil {
		return nil, err
	}

	for _, routeInput := range req.Routes {
		if err := m.createRouteWithLocations(ctx, link.ID, routeInput); err != nil {
			return nil, err
		}
	}

	if len(req.Routes) == 0 {
		if err := m.createRouteWithLocations(ctx, link.ID, dto.DeliveryRouteInput{
			Label: "Default Route",
			PickupLocations: []dto.DeliveryRouteLocationInput{
				{Label: "My store", UseStoreLocation: true},
			},
		}); err != nil {
			return nil, err
		}
	}

	return m.buildAgentResponse(ctx, storeID, agent.ID)
}

func (m *deliveryModule) createRouteWithLocations(ctx context.Context, linkID int64, input dto.DeliveryRouteInput) error {
	label := input.Label
	if label == "" {
		label = "Route"
	}
	route := &db.DeliveryRoute{
		StoreDeliveryLinkID: linkID,
		Label:               label,
		IsActive:            true,
	}
	if err := m.deliveryStorage.CreateRoute(ctx, route); err != nil {
		return err
	}

	pickups := input.PickupLocations
	if len(pickups) == 0 {
		pickups = []dto.DeliveryRouteLocationInput{{Label: "My store", UseStoreLocation: true}}
	}
	for _, p := range pickups {
		if err := m.createLocation(ctx, route.ID, constant.DeliveryLocationTypePickup, p); err != nil {
			return err
		}
	}
	for _, d := range input.DeliveryLocations {
		if err := m.createLocation(ctx, route.ID, constant.DeliveryLocationTypeDelivery, d); err != nil {
			return err
		}
	}
	return nil
}

func (m *deliveryModule) createLocation(ctx context.Context, routeID int64, locType constant.DeliveryLocationType, input dto.DeliveryRouteLocationInput) error {
	country := input.Country
	if country == "" {
		country = "Ethiopia"
	}
	loc := &db.DeliveryRouteLocation{
		DeliveryRouteID:  routeID,
		LocationType:     locType,
		Label:            input.Label,
		Country:          country,
		Region:           input.Region,
		City:             input.City,
		Street:           input.Street,
		Landmark:         input.Landmark,
		Notes:            input.Notes,
		UseStoreLocation: input.UseStoreLocation,
	}
	return m.deliveryStorage.CreateRouteLocation(ctx, loc)
}

func (m *deliveryModule) ListAgents(ctx context.Context, storeID int64) ([]dto.DeliveryAgentResponse, error) {
	links, err := m.deliveryStorage.ListLinksWithRoutesForStore(ctx, storeID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.DeliveryAgentResponse, 0, len(links))
	for _, link := range links {
		resp, err := m.agentResponseFromLink(link)
		if err != nil {
			return nil, err
		}
		out = append(out, *resp)
	}
	return out, nil
}

func (m *deliveryModule) UpdateAgent(ctx context.Context, storeID, agentID int64, req dto.UpdateDeliveryAgentRequest) (*dto.DeliveryAgentResponse, error) {
	link, err := m.deliveryStorage.GetStoreLink(ctx, storeID, agentID)
	if err != nil {
		return nil, fmt.Errorf("delivery agent not connected to this store")
	}

	agent, err := m.deliveryStorage.GetAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}
	if req.FullName != nil {
		agent.FullName = *req.FullName
	}
	if req.Phone != nil {
		agent.Phone = *req.Phone
	}
	if err := m.deliveryStorage.UpdateAgent(ctx, agent); err != nil {
		return nil, err
	}
	if req.ShareEnabled != nil {
		link.ShareEnabled = *req.ShareEnabled
		if err := m.deliveryStorage.UpdateStoreLink(ctx, link); err != nil {
			return nil, err
		}
	}

	if len(req.Routes) > 0 {
		routes, err := m.deliveryStorage.ListRoutesByLinkID(ctx, link.ID)
		if err != nil {
			return nil, err
		}
		for _, r := range routes {
			if err := m.deliveryStorage.DeleteRoute(ctx, r.ID); err != nil {
				return nil, err
			}
		}
		for _, routeInput := range req.Routes {
			if err := m.createRouteWithLocations(ctx, link.ID, routeInput); err != nil {
				return nil, err
			}
		}
	}

	return m.buildAgentResponse(ctx, storeID, agentID)
}

func (m *deliveryModule) AddRoute(ctx context.Context, storeID, agentID int64, req dto.AddDeliveryRouteRequest) (*dto.DeliveryRouteResponse, error) {
	link, err := m.deliveryStorage.GetStoreLink(ctx, storeID, agentID)
	if err != nil {
		return nil, fmt.Errorf("delivery agent not connected to this store")
	}
	input := dto.DeliveryRouteInput{
		Label:             req.Label,
		PickupLocations:   req.PickupLocations,
		DeliveryLocations: req.DeliveryLocations,
	}
	if err := m.createRouteWithLocations(ctx, link.ID, input); err != nil {
		return nil, err
	}
	routes, err := m.deliveryStorage.ListRoutesByLinkID(ctx, link.ID)
	if err != nil || len(routes) == 0 {
		return nil, fmt.Errorf("failed to load created route")
	}
	last := routes[len(routes)-1]
	resp := mapRouteToResponse(last)
	return &resp, nil
}

func (m *deliveryModule) UpdateRoute(ctx context.Context, storeID, routeID int64, req dto.UpdateDeliveryRouteRequest) (*dto.DeliveryRouteResponse, error) {
	route, err := m.deliveryStorage.GetRouteByID(ctx, routeID)
	if err != nil {
		return nil, err
	}
	links, err := m.deliveryStorage.ListStoreLinksByStoreID(ctx, storeID)
	if err != nil {
		return nil, err
	}
	var ownerLink *db.StoreDeliveryLink
	for i := range links {
		if links[i].ID == route.StoreDeliveryLinkID {
			ownerLink = &links[i]
			break
		}
	}
	if ownerLink == nil {
		return nil, fmt.Errorf("route not found for this store")
	}

	if req.Label != nil {
		route.Label = *req.Label
	}
	if req.IsActive != nil {
		route.IsActive = *req.IsActive
	}
	if err := m.deliveryStorage.UpdateRoute(ctx, route); err != nil {
		return nil, err
	}

	if len(req.PickupLocations) > 0 || len(req.DeliveryLocations) > 0 {
		if err := m.deliveryStorage.DeleteRouteLocationsByRouteID(ctx, routeID); err != nil {
			return nil, err
		}
		pickups := req.PickupLocations
		if len(pickups) == 0 {
			pickups = []dto.DeliveryRouteLocationInput{{Label: "My store", UseStoreLocation: true}}
		}
		for _, p := range pickups {
			if err := m.createLocation(ctx, routeID, constant.DeliveryLocationTypePickup, p); err != nil {
				return nil, err
			}
		}
		for _, d := range req.DeliveryLocations {
			if err := m.createLocation(ctx, routeID, constant.DeliveryLocationTypeDelivery, d); err != nil {
				return nil, err
			}
		}
	}

	updated, err := m.deliveryStorage.GetRouteByID(ctx, routeID)
	if err != nil {
		return nil, err
	}
	resp := mapRouteToResponse(*updated)
	return &resp, nil
}

func (m *deliveryModule) DeleteRoute(ctx context.Context, storeID, routeID int64) error {
	route, err := m.deliveryStorage.GetRouteByID(ctx, routeID)
	if err != nil {
		return err
	}
	links, err := m.deliveryStorage.ListStoreLinksByStoreID(ctx, storeID)
	if err != nil {
		return err
	}
	found := false
	for _, l := range links {
		if l.ID == route.StoreDeliveryLinkID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("route not found for this store")
	}
	return m.deliveryStorage.DeleteRoute(ctx, routeID)
}

func (m *deliveryModule) DisconnectAgent(ctx context.Context, storeID, agentID int64) error {
	_, err := m.deliveryStorage.GetStoreLink(ctx, storeID, agentID)
	if err != nil {
		return fmt.Errorf("delivery agent not connected to this store")
	}
	return m.deliveryStorage.DeleteStoreLink(ctx, storeID, agentID)
}

func (m *deliveryModule) ListSharedAgents(ctx context.Context, storeID int64) ([]dto.DeliveryAgentResponse, error) {
	agents, err := m.deliveryStorage.ListSharedAgents(ctx, storeID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.DeliveryAgentResponse, 0, len(agents))
	for _, agent := range agents {
		resp := dto.DeliveryAgentResponse{
			ID:           agent.ID,
			Username:     agent.Username,
			FullName:     agent.FullName,
			Phone:        agent.Phone,
			Status:       string(agent.Status),
			LoyaltyScore: agent.LoyaltyScore,
			ShareEnabled: true,
		}
		ownerLinks, err := m.deliveryStorage.ListShareEnabledLinksForAgent(ctx, agent.ID)
		if err == nil {
			resp.Routes = mapRoutesFromLinks(ownerLinks)
		}
		out = append(out, resp)
	}
	return out, nil
}

func (m *deliveryModule) AdoptAgent(ctx context.Context, storeID, userID, agentID int64) (*dto.DeliveryAgentResponse, error) {
	has, err := m.deliveryStorage.StoreHasAgent(ctx, storeID, agentID)
	if err != nil {
		return nil, err
	}
	if has {
		return nil, fmt.Errorf("agent already connected to this store")
	}

	agent, err := m.deliveryStorage.GetAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}

	ownerLink, err := m.deliveryStorage.GetSharedOwnerLink(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("agent is not available for sharing")
	}

	link := &db.StoreDeliveryLink{
		StoreID:           storeID,
		DeliveryAgentID:   agentID,
		ShareEnabled:      false,
		ConnectedByUserID: &userID,
	}
	if err := m.deliveryStorage.CreateStoreLink(ctx, link); err != nil {
		return nil, err
	}

	if err := m.deliveryStorage.CreateAgentShare(ctx, &db.DeliveryAgentShare{
		OwnerStoreID:    ownerLink.StoreID,
		DeliveryAgentID: agentID,
		AdoptedStoreID:  storeID,
	}); err != nil {
		return nil, err
	}

	if err := m.copyRoutesFromLink(ctx, ownerLink.ID, link.ID); err != nil {
		return nil, err
	}

	return m.buildAgentResponse(ctx, storeID, agent.ID)
}

func (m *deliveryModule) copyRoutesFromLink(ctx context.Context, fromLinkID, toLinkID int64) error {
	routes, err := m.deliveryStorage.ListRoutesByLinkID(ctx, fromLinkID)
	if err != nil {
		return err
	}
	for _, r := range routes {
		newRoute := &db.DeliveryRoute{
			StoreDeliveryLinkID: toLinkID,
			Label:               r.Label,
			IsActive:            r.IsActive,
		}
		if err := m.deliveryStorage.CreateRoute(ctx, newRoute); err != nil {
			return err
		}
		for _, loc := range r.Locations {
			newLoc := loc
			newLoc.ID = 0
			newLoc.DeliveryRouteID = newRoute.ID
			if err := m.deliveryStorage.CreateRouteLocation(ctx, &newLoc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *deliveryModule) GetDeliverySuggestions(ctx context.Context, storeID, orderID int64) ([]dto.DeliverySuggestion, error) {
	order, err := m.orderStorage.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order.StoreID != storeID {
		return nil, fmt.Errorf("order does not belong to your store")
	}

	store, err := m.storeStorage.GetStoreByID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	deliveryStreet, deliveryCity, deliveryRegion, deliveryCountry := "", "", "", ""
	if order.ShippingAddress != nil {
		deliveryStreet = order.ShippingAddress.Street
		deliveryCity = order.ShippingAddress.City
		deliveryRegion = order.ShippingAddress.Region
		deliveryCountry = order.ShippingAddress.Country
	}

	links, err := m.deliveryStorage.ListLinksWithRoutesForStore(ctx, storeID)
	if err != nil {
		return nil, err
	}

	type candidate struct {
		suggestion dto.DeliverySuggestion
		totalScore int
	}
	var candidates []candidate

	for _, link := range links {
		if link.DeliveryAgent.Status != constant.DeliveryAgentStatusActive &&
			link.DeliveryAgent.Status != constant.DeliveryAgentStatusPendingInvite {
			continue
		}
		agent := link.DeliveryAgent
		bestScore := 0
		var bestRouteID int64
		var bestPickup, bestDelivery string

		for _, route := range link.Routes {
			if !route.IsActive {
				continue
			}
			ok, rm := scoreRoute(route, store.Location, "", "",
				deliveryStreet, deliveryCity, deliveryRegion, deliveryCountry)
			if !ok {
				continue
			}
			total := rm.PickupScore + rm.DeliveryScore + agent.LoyaltyScore
			if total > bestScore {
				bestScore = total
				bestRouteID = rm.RouteID
				bestPickup = rm.PickupSummary
				bestDelivery = rm.DeliverySummary
			}
		}

		if bestScore > 0 {
			candidates = append(candidates, candidate{
				suggestion: dto.DeliverySuggestion{
					AgentID:      agent.ID,
					RouteID:      bestRouteID,
					Username:     agent.Username,
					FullName:     agent.FullName,
					Phone:        agent.Phone,
					LoyaltyScore: agent.LoyaltyScore,
					Score:        bestScore,
					MatchSummary: fmt.Sprintf("pickup: %s, delivery: %s", bestPickup, bestDelivery),
				},
				totalScore: bestScore,
			})
		}
	}

	// sort descending by score
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].totalScore > candidates[i].totalScore {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	out := make([]dto.DeliverySuggestion, len(candidates))
	for i, c := range candidates {
		c.suggestion.Suggested = i == 0
		out[i] = c.suggestion
	}
	return out, nil
}

func (m *deliveryModule) DispatchOnShip(ctx context.Context, storeID, orderID int64, agentID, routeID *int64) error {
	order, err := m.orderStorage.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order.StoreID != storeID {
		return fmt.Errorf("order does not belong to your store")
	}

	var finalAgentID, finalRouteID int64
	if agentID != nil {
		finalAgentID = *agentID
		has, err := m.deliveryStorage.StoreHasAgent(ctx, storeID, finalAgentID)
		if err != nil || !has {
			return fmt.Errorf("delivery agent not connected to this store")
		}
		if routeID != nil {
			finalRouteID = *routeID
		}
	} else {
		suggestions, err := m.GetDeliverySuggestions(ctx, storeID, orderID)
		if err != nil {
			return err
		}
		if len(suggestions) == 0 {
			return fmt.Errorf("no_matching_courier")
		}
		finalAgentID = suggestions[0].AgentID
		finalRouteID = suggestions[0].RouteID
	}

	if err := m.orderStorage.UpdateOrderDispatch(ctx, orderID, "shipped", finalAgentID, finalRouteID); err != nil {
		return err
	}

	go m.notifyDeliveryAgent(context.Background(), orderID, storeID, finalAgentID, finalRouteID)
	return nil
}

func (m *deliveryModule) notifyDeliveryAgent(ctx context.Context, orderID, storeID, agentID, routeID int64) {
	agent, err := m.deliveryStorage.GetAgentByID(ctx, agentID)
	if err != nil || agent.TelegramUserID == nil {
		logger.Warn("skipping delivery dispatch notification: agent has no telegram id",
			zap.Int64("order_id", orderID), zap.Int64("agent_id", agentID))
		return
	}

	store, err := m.storeStorage.GetStoreByID(ctx, storeID)
	if err != nil {
		return
	}
	order, err := m.orderStorage.GetOrderByID(ctx, orderID)
	if err != nil {
		return
	}

	pickupLocation := store.Location
	if routeID > 0 {
		route, err := m.deliveryStorage.GetRouteByID(ctx, routeID)
		if err == nil {
			for _, loc := range route.Locations {
				if loc.LocationType == constant.DeliveryLocationTypePickup {
					if loc.UseStoreLocation {
						pickupLocation = store.Location
					} else {
						pickupLocation = formatLocation(loc.Street, loc.City, loc.Region, loc.Landmark)
					}
					break
				}
			}
		}
	}

	var orderDTO dto.Order
	if m.orderModule != nil {
		if o, err := m.orderModule.GetOrder(ctx, orderID); err == nil {
			orderDTO = *o
		}
	}
	if orderDTO.ID == 0 {
		orderDTO = dto.Order{ID: order.ID, TotalPrice: order.TotalPrice}
	}

	var addr *dto.Address
	if order.ShippingAddress != nil {
		addr = &dto.Address{
			RecipientName: order.ShippingAddress.RecipientName,
			Phone:         order.ShippingAddress.Phone,
			Street:        order.ShippingAddress.Street,
			City:          order.ShippingAddress.City,
			Region:        order.ShippingAddress.Region,
			Country:       order.ShippingAddress.Country,
		}
	}

	if err := m.tele.SendDeliveryDispatchNotification(*agent.TelegramUserID, orderDTO, addr, store.Name, pickupLocation); err != nil {
		logger.Error("failed to send delivery dispatch notification", zap.Error(err), zap.Int64("order_id", orderID))
	}
}

func formatLocation(street, city, region, landmark string) string {
	parts := []string{}
	if street != "" {
		parts = append(parts, street)
	}
	if landmark != "" {
		parts = append(parts, landmark)
	}
	if city != "" {
		parts = append(parts, city)
	}
	if region != "" {
		parts = append(parts, region)
	}
	return strings.Join(parts, ", ")
}

func (m *deliveryModule) buildAgentResponse(ctx context.Context, storeID, agentID int64) (*dto.DeliveryAgentResponse, error) {
	link, err := m.deliveryStorage.GetStoreLink(ctx, storeID, agentID)
	if err != nil {
		return nil, err
	}
	agent, err := m.deliveryStorage.GetAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}
	link.DeliveryAgent = *agent
	routes, err := m.deliveryStorage.ListRoutesByLinkID(ctx, link.ID)
	if err != nil {
		return nil, err
	}
	link.Routes = routes
	return m.agentResponseFromLink(*link)
}

func (m *deliveryModule) agentResponseFromLink(link db.StoreDeliveryLink) (*dto.DeliveryAgentResponse, error) {
	routes := make([]dto.DeliveryRouteResponse, len(link.Routes))
	for i, r := range link.Routes {
		routes[i] = mapRouteToResponse(r)
	}
	return &dto.DeliveryAgentResponse{
		ID:           link.DeliveryAgent.ID,
		Username:     link.DeliveryAgent.Username,
		FullName:     link.DeliveryAgent.FullName,
		Phone:        link.DeliveryAgent.Phone,
		Status:       string(link.DeliveryAgent.Status),
		LoyaltyScore: link.DeliveryAgent.LoyaltyScore,
		ShareEnabled: link.ShareEnabled,
		Routes:       routes,
	}, nil
}

func mapRoutesFromLinks(links []db.StoreDeliveryLink) []dto.DeliveryRouteResponse {
	var routes []dto.DeliveryRouteResponse
	for _, link := range links {
		for _, r := range link.Routes {
			routes = append(routes, mapRouteToResponse(r))
		}
	}
	return routes
}

func mapRouteToResponse(r db.DeliveryRoute) dto.DeliveryRouteResponse {
	resp := dto.DeliveryRouteResponse{
		ID:       r.ID,
		Label:    r.Label,
		IsActive: r.IsActive,
	}
	for _, loc := range r.Locations {
		item := dto.DeliveryRouteLocationResponse{
			ID:               loc.ID,
			LocationType:     string(loc.LocationType),
			Label:            loc.Label,
			Country:          loc.Country,
			Region:           loc.Region,
			City:             loc.City,
			Street:           loc.Street,
			Landmark:         loc.Landmark,
			Notes:            loc.Notes,
			UseStoreLocation: loc.UseStoreLocation,
		}
		if loc.LocationType == constant.DeliveryLocationTypePickup {
			resp.PickupLocations = append(resp.PickupLocations, item)
		} else {
			resp.DeliveryLocations = append(resp.DeliveryLocations, item)
		}
	}
	return resp
}

// ActivatePendingInvite links telegram user on bot /start
func (m *deliveryModule) ActivatePendingInvite(ctx context.Context, telegramUserID int64, username string, userID int64) (bool, error) {
	agent, err := m.deliveryStorage.GetAgentByUsername(ctx, username)
	if err != nil {
		agent, err = m.deliveryStorage.GetAgentByTelegramUserID(ctx, telegramUserID)
		if err != nil {
			return false, nil
		}
	}
	if agent.Status != constant.DeliveryAgentStatusPendingInvite && agent.Status != constant.DeliveryAgentStatusActive {
		return false, nil
	}
	agent.TelegramUserID = &telegramUserID
	agent.UserID = &userID
	agent.Status = constant.DeliveryAgentStatusActive
	if err := m.deliveryStorage.UpdateAgent(ctx, agent); err != nil {
		return false, err
	}
	return true, nil
}
