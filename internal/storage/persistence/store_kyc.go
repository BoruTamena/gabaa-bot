package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
	"gorm.io/gorm"
)

type storeKYCPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewStoreKYCPersistence(db *gorm.DB, logger platform.Logger) storage.StoreKYCStorage {
	return &storeKYCPersistence{db: db, logger: logger}
}

func (p *storePersistence) UpdateStoreVerificationStatus(ctx context.Context, storeID int64, status string) error {
	result := p.db.WithContext(ctx).
		Model(&db.Store{}).
		Where("id = ?", storeID).
		Update("verification_status", status)
	if result.Error != nil {
		p.logger.Error("Failed to update store verification status", "error", result.Error, "storeID", storeID)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (p *storeKYCPersistence) UpsertStoreKYC(ctx context.Context, kyc *db.StoreKYC) error {
	var existing db.StoreKYC
	err := p.db.WithContext(ctx).Where("store_id = ?", kyc.StoreID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := p.db.WithContext(ctx).Create(kyc).Error; err != nil {
			p.logger.Error("Failed to create store KYC", "error", err, "storeID", kyc.StoreID)
			return err
		}
		return nil
	}
	if err != nil {
		p.logger.Error("Failed to lookup store KYC", "error", err, "storeID", kyc.StoreID)
		return err
	}

	existing.TINNumber = kyc.TINNumber
	existing.BusinessRegistrationNumber = kyc.BusinessRegistrationNumber
	existing.TINCertificateURL = kyc.TINCertificateURL
	existing.BusinessLicenseURL = kyc.BusinessLicenseURL
	existing.ReviewNote = ""
	existing.SubmittedAt = kyc.SubmittedAt
	existing.ReviewedAt = nil

	if err := p.db.WithContext(ctx).Save(&existing).Error; err != nil {
		p.logger.Error("Failed to update store KYC", "error", err, "storeID", kyc.StoreID)
		return err
	}
	return nil
}

func (p *storeKYCPersistence) GetStoreKYCByStoreID(ctx context.Context, storeID int64) (*db.StoreKYC, error) {
	var kyc db.StoreKYC
	err := p.db.WithContext(ctx).Preload("Store").Where("store_id = ?", storeID).First(&kyc).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			p.logger.Error("Failed to get store KYC", "error", err, "storeID", storeID)
		}
		return nil, err
	}
	return &kyc, nil
}

func (p *storeKYCPersistence) ListStoreKYCByVerificationStatus(ctx context.Context, status string) ([]db.StoreKYC, error) {
	var records []db.StoreKYC
	err := p.db.WithContext(ctx).
		Preload("Store").
		Joins("JOIN stores ON stores.id = store_kyc.store_id").
		Where("stores.verification_status = ?", status).
		Find(&records).Error
	if err != nil {
		p.logger.Error("Failed to list store KYC by verification status", "error", err, "status", status)
		return nil, err
	}
	return records, nil
}

func (p *storeKYCPersistence) UpdateStoreKYCReview(ctx context.Context, storeID int64, reviewNote string, reviewedAt time.Time) error {
	result := p.db.WithContext(ctx).
		Model(&db.StoreKYC{}).
		Where("store_id = ?", storeID).
		Updates(map[string]interface{}{
			"review_note": reviewNote,
			"reviewed_at": reviewedAt,
		})
	if result.Error != nil {
		p.logger.Error("Failed to update store KYC review", "error", result.Error, "storeID", storeID)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
