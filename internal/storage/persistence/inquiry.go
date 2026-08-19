package persistence

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
	"gorm.io/gorm"
)

type inquiryPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewInquiryPersistence(db *gorm.DB, logger platform.Logger) storage.InquiryStorage {
	return &inquiryPersistence{db: db, logger: logger}
}

func (p *inquiryPersistence) CreateInquiry(ctx context.Context, inquiry *db.ProductInquiry) error {
	err := p.db.WithContext(ctx).Create(inquiry).Error
	if err != nil {
		p.logger.Error("Failed to create product inquiry", "error", err)
	}
	return err
}

func (p *inquiryPersistence) GetInquiryByID(ctx context.Context, id int64) (*db.ProductInquiry, error) {
	var inquiry db.ProductInquiry
	err := p.db.WithContext(ctx).
		Preload("Product").
		First(&inquiry, id).Error
	if err != nil {
		p.logger.Error("Failed to get product inquiry by ID", "error", err, "inquiryID", id)
		return nil, err
	}
	return &inquiry, nil
}

func (p *inquiryPersistence) ListInquiriesByStore(ctx context.Context, storeID int64, filter dto.ProductInquiryFilterParams) ([]db.ProductInquiry, int64, error) {
	var inquiries []db.ProductInquiry
	var count int64

	query := p.db.WithContext(ctx).Model(&db.ProductInquiry{}).Where("store_id = ?", storeID)
	if filter.ProductID != 0 {
		query = query.Where("product_id = ?", filter.ProductID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if err := query.Count(&count).Error; err != nil {
		p.logger.Error("Failed to count store product inquiries", "error", err, "storeID", storeID)
		return nil, 0, err
	}

	err := query.Preload("Product").
		Order("created_at DESC").
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Find(&inquiries).Error
	if err != nil {
		p.logger.Error("Failed to list store product inquiries", "error", err, "storeID", storeID)
		return nil, 0, err
	}
	return inquiries, count, nil
}

func (p *inquiryPersistence) ListInquiriesByCustomer(ctx context.Context, customerID int64, filter dto.ProductInquiryFilterParams) ([]db.ProductInquiry, int64, error) {
	var inquiries []db.ProductInquiry
	var count int64

	query := p.db.WithContext(ctx).Model(&db.ProductInquiry{}).Where("customer_id = ?", customerID)
	if filter.ProductID != 0 {
		query = query.Where("product_id = ?", filter.ProductID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if err := query.Count(&count).Error; err != nil {
		p.logger.Error("Failed to count customer product inquiries", "error", err, "customerID", customerID)
		return nil, 0, err
	}

	err := query.Preload("Product").
		Order("created_at DESC").
		Limit(filter.GetLimit()).
		Offset(filter.GetOffset()).
		Find(&inquiries).Error
	if err != nil {
		p.logger.Error("Failed to list customer product inquiries", "error", err, "customerID", customerID)
		return nil, 0, err
	}
	return inquiries, count, nil
}

func (p *inquiryPersistence) UpdateInquiry(ctx context.Context, inquiry *db.ProductInquiry) error {
	err := p.db.WithContext(ctx).Save(inquiry).Error
	if err != nil {
		p.logger.Error("Failed to update product inquiry", "error", err, "inquiryID", inquiry.ID)
	}
	return err
}
