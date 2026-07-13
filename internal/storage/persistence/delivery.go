package persistence

import (
	"context"
	"strings"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
	"gorm.io/gorm"
)

type deliveryPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewDeliveryStorage(db *gorm.DB, logger platform.Logger) storage.DeliveryStorage {
	return &deliveryPersistence{db: db, logger: logger}
}

func normalizeUsername(u string) string {
	u = strings.TrimSpace(strings.TrimPrefix(u, "@"))
	return strings.ToLower(u)
}

func (p *deliveryPersistence) CreateAgent(ctx context.Context, agent *db.DeliveryAgent) error {
	agent.Username = normalizeUsername(agent.Username)
	return p.db.WithContext(ctx).Create(agent).Error
}

func (p *deliveryPersistence) GetAgentByID(ctx context.Context, id int64) (*db.DeliveryAgent, error) {
	var agent db.DeliveryAgent
	err := p.db.WithContext(ctx).First(&agent, id).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (p *deliveryPersistence) GetAgentByUsername(ctx context.Context, username string) (*db.DeliveryAgent, error) {
	var agent db.DeliveryAgent
	err := p.db.WithContext(ctx).
		Where("LOWER(username) = ?", normalizeUsername(username)).
		First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (p *deliveryPersistence) GetAgentByTelegramUserID(ctx context.Context, telegramUserID int64) (*db.DeliveryAgent, error) {
	var agent db.DeliveryAgent
	err := p.db.WithContext(ctx).Where("telegram_user_id = ?", telegramUserID).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (p *deliveryPersistence) GetAgentByUserID(ctx context.Context, userID int64) (*db.DeliveryAgent, error) {
	var agent db.DeliveryAgent
	err := p.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, constant.DeliveryAgentStatusActive).
		First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (p *deliveryPersistence) UpdateAgent(ctx context.Context, agent *db.DeliveryAgent) error {
	return p.db.WithContext(ctx).Save(agent).Error
}

func (p *deliveryPersistence) AdjustLoyaltyScore(ctx context.Context, agentID int64, delta int) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var agent db.DeliveryAgent
		if err := tx.First(&agent, agentID).Error; err != nil {
			return err
		}
		score := agent.LoyaltyScore + delta
		if score < 0 {
			score = 0
		}
		return tx.Model(&agent).Update("loyalty_score", score).Error
	})
}

func (p *deliveryPersistence) CreateStoreLink(ctx context.Context, link *db.StoreDeliveryLink) error {
	return p.db.WithContext(ctx).Create(link).Error
}

func (p *deliveryPersistence) GetStoreLink(ctx context.Context, storeID, agentID int64) (*db.StoreDeliveryLink, error) {
	var link db.StoreDeliveryLink
	err := p.db.WithContext(ctx).
		Where("store_id = ? AND delivery_agent_id = ?", storeID, agentID).
		First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (p *deliveryPersistence) ListStoreLinksByStoreID(ctx context.Context, storeID int64) ([]db.StoreDeliveryLink, error) {
	var links []db.StoreDeliveryLink
	err := p.db.WithContext(ctx).
		Preload("DeliveryAgent").
		Preload("Routes", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true)
		}).
		Preload("Routes.Locations").
		Where("store_id = ?", storeID).
		Find(&links).Error
	return links, err
}

func (p *deliveryPersistence) UpdateStoreLink(ctx context.Context, link *db.StoreDeliveryLink) error {
	return p.db.WithContext(ctx).Save(link).Error
}

func (p *deliveryPersistence) DeleteStoreLink(ctx context.Context, storeID, agentID int64) error {
	return p.db.WithContext(ctx).
		Where("store_id = ? AND delivery_agent_id = ?", storeID, agentID).
		Delete(&db.StoreDeliveryLink{}).Error
}

func (p *deliveryPersistence) StoreHasAgent(ctx context.Context, storeID, agentID int64) (bool, error) {
	var count int64
	err := p.db.WithContext(ctx).Model(&db.StoreDeliveryLink{}).
		Where("store_id = ? AND delivery_agent_id = ?", storeID, agentID).
		Count(&count).Error
	return count > 0, err
}

func (p *deliveryPersistence) CreateRoute(ctx context.Context, route *db.DeliveryRoute) error {
	return p.db.WithContext(ctx).Create(route).Error
}

func (p *deliveryPersistence) GetRouteByID(ctx context.Context, routeID int64) (*db.DeliveryRoute, error) {
	var route db.DeliveryRoute
	err := p.db.WithContext(ctx).
		Preload("Locations").
		First(&route, routeID).Error
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (p *deliveryPersistence) UpdateRoute(ctx context.Context, route *db.DeliveryRoute) error {
	return p.db.WithContext(ctx).Save(route).Error
}

func (p *deliveryPersistence) DeleteRoute(ctx context.Context, routeID int64) error {
	return p.db.WithContext(ctx).Delete(&db.DeliveryRoute{}, routeID).Error
}

func (p *deliveryPersistence) ListRoutesByLinkID(ctx context.Context, linkID int64) ([]db.DeliveryRoute, error) {
	var routes []db.DeliveryRoute
	err := p.db.WithContext(ctx).
		Preload("Locations").
		Where("store_delivery_link_id = ?", linkID).
		Find(&routes).Error
	return routes, err
}

func (p *deliveryPersistence) CreateRouteLocation(ctx context.Context, loc *db.DeliveryRouteLocation) error {
	return p.db.WithContext(ctx).Create(loc).Error
}

func (p *deliveryPersistence) DeleteRouteLocationsByRouteID(ctx context.Context, routeID int64) error {
	return p.db.WithContext(ctx).
		Where("delivery_route_id = ?", routeID).
		Delete(&db.DeliveryRouteLocation{}).Error
}

func (p *deliveryPersistence) ListRouteLocationsByRouteID(ctx context.Context, routeID int64) ([]db.DeliveryRouteLocation, error) {
	var locs []db.DeliveryRouteLocation
	err := p.db.WithContext(ctx).
		Where("delivery_route_id = ?", routeID).
		Find(&locs).Error
	return locs, err
}

func (p *deliveryPersistence) CreateAgentShare(ctx context.Context, share *db.DeliveryAgentShare) error {
	share.CreatedAt = time.Now()
	return p.db.WithContext(ctx).Create(share).Error
}

func (p *deliveryPersistence) ListSharedAgents(ctx context.Context, excludeStoreID int64) ([]db.DeliveryAgent, error) {
	var agents []db.DeliveryAgent
	err := p.db.WithContext(ctx).
		Distinct().
		Joins("JOIN store_delivery_links ON store_delivery_links.delivery_agent_id = delivery_agents.id").
		Where("store_delivery_links.share_enabled = ?", true).
		Where("delivery_agents.id NOT IN (?)",
			p.db.Model(&db.StoreDeliveryLink{}).
				Select("delivery_agent_id").
				Where("store_id = ?", excludeStoreID),
		).
		Find(&agents).Error
	return agents, err
}

func (p *deliveryPersistence) ListAgentsForStore(ctx context.Context, storeID int64) ([]db.DeliveryAgent, error) {
	var agents []db.DeliveryAgent
	err := p.db.WithContext(ctx).
		Joins("JOIN store_delivery_links ON store_delivery_links.delivery_agent_id = delivery_agents.id").
		Where("store_delivery_links.store_id = ?", storeID).
		Find(&agents).Error
	return agents, err
}

func (p *deliveryPersistence) ListLinksWithRoutesForStore(ctx context.Context, storeID int64) ([]db.StoreDeliveryLink, error) {
	return p.ListStoreLinksByStoreID(ctx, storeID)
}

func (p *deliveryPersistence) ListAreaPresets(ctx context.Context) ([]db.DeliveryAreaPreset, error) {
	var presets []db.DeliveryAreaPreset
	err := p.db.WithContext(ctx).Where("is_active = ?", true).Find(&presets).Error
	return presets, err
}

func (p *deliveryPersistence) GetSharedOwnerLink(ctx context.Context, agentID int64) (*db.StoreDeliveryLink, error) {
	var link db.StoreDeliveryLink
	err := p.db.WithContext(ctx).
		Preload("Routes").
		Preload("Routes.Locations").
		Where("delivery_agent_id = ? AND share_enabled = ?", agentID, true).
		First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (p *deliveryPersistence) ListLinksByAgentID(ctx context.Context, agentID int64) ([]db.StoreDeliveryLink, error) {
	var links []db.StoreDeliveryLink
	err := p.db.WithContext(ctx).
		Preload("Routes").
		Preload("Routes.Locations").
		Where("delivery_agent_id = ?", agentID).
		Find(&links).Error
	return links, err
}

func (p *deliveryPersistence) ListShareEnabledLinksForAgent(ctx context.Context, agentID int64) ([]db.StoreDeliveryLink, error) {
	var links []db.StoreDeliveryLink
	err := p.db.WithContext(ctx).
		Preload("Routes").
		Preload("Routes.Locations").
		Where("delivery_agent_id = ? AND share_enabled = ?", agentID, true).
		Find(&links).Error
	return links, err
}
