package persistence

import (
	"context"
	"strconv"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"gorm.io/gorm"
)

func (p *storePersistence) ListStoresAdmin(ctx context.Context, filter dto.AdminStoreFilterParams) ([]db.Store, int64, error) {
	var stores []db.Store
	var count int64

	query := p.db.WithContext(ctx).Model(&db.Store{})
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.VerificationStatus != "" {
		query = query.Where("verification_status = ?", filter.VerificationStatus)
	}
	if filter.Category != "" {
		query = query.Where("category ILIKE ?", filter.Category)
	}
	if filter.Query != "" {
		term := "%" + filter.Query + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ? OR location ILIKE ? OR phone ILIKE ?", term, term, term, term)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("created_at DESC").
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Find(&stores).Error
	return stores, count, err
}

func (p *storePersistence) UpdateStoreStatus(ctx context.Context, storeID int64, status string) error {
	result := p.db.WithContext(ctx).Model(&db.Store{}).Where("id = ?", storeID).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (p *persistence) ListUsers(ctx context.Context, filter dto.AdminUserFilterParams) ([]db.User, int64, error) {
	var users []db.User
	var count int64

	query := p.db.WithContext(ctx).Model(&db.User{})
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}
	if filter.Query != "" {
		term := "%" + filter.Query + "%"
		if id, err := strconv.ParseInt(filter.Query, 10, 64); err == nil {
			query = query.Where("username ILIKE ? OR telegram_user_id = ?", term, id)
		} else {
			query = query.Where("username ILIKE ?", term)
		}
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	err := query.Order("created_at DESC").
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Find(&users).Error
	return users, count, err
}

func (p *storeKYCPersistence) ListStoreKYCAdmin(ctx context.Context, filter dto.AdminStoreKYCFilterParams) ([]db.StoreKYC, int64, error) {
	var records []db.StoreKYC
	var count int64

	status := filter.Status
	if status == "" {
		status = "pending_review"
	}

	query := p.db.WithContext(ctx).Model(&db.StoreKYC{}).
		Joins("JOIN stores ON stores.id = store_kyc.store_id").
		Where("stores.verification_status = ?", status)

	if filter.Query != "" {
		term := "%" + filter.Query + "%"
		query = query.Where("stores.name ILIKE ? OR store_kyc.tin_number ILIKE ?", term, term)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Store").
		Order("store_kyc.submitted_at DESC").
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Find(&records).Error
	return records, count, err
}

func (p *orderPersistence) GetOrdersByAdminFilter(ctx context.Context, filter dto.AdminOrderFilterParams) ([]db.Order, int64, error) {
	var orders []db.Order
	var count int64

	base := p.db.WithContext(ctx).Model(&db.Order{})
	if filter.StoreID != 0 {
		base = base.Where("orders.store_id = ?", filter.StoreID)
	}
	if filter.OrderID != nil {
		base = base.Where("orders.id = ?", *filter.OrderID)
	}
	if filter.Status != "" {
		base = base.Where("orders.status = ?", filter.Status)
	}
	if filter.Query != "" {
		term := "%" + filter.Query + "%"
		base = base.Joins("LEFT JOIN users ON users.id = orders.user_id").
			Where("users.username ILIKE ?", term)
	}

	if err := base.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	err := base.
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Preload("ShippingAddress").
		Order("orders.created_at DESC").
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Find(&orders).Error
	return orders, count, err
}
