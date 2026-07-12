package persistence

import (
	"context"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
	"gorm.io/gorm"
)

type withdrawalPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewWithdrawalPersistence(db *gorm.DB, logger platform.Logger) storage.WithdrawalStorage {
	return &withdrawalPersistence{db: db, logger: logger}
}

func (p *withdrawalPersistence) CreateWithdrawal(ctx context.Context, withdrawal *db.Withdrawal) error {
	err := p.db.WithContext(ctx).Create(withdrawal).Error
	if err != nil {
		p.logger.Error("Failed to create withdrawal", "error", err)
	}
	return err
}

func (p *withdrawalPersistence) UpdateWithdrawal(ctx context.Context, withdrawal *db.Withdrawal) error {
	err := p.db.WithContext(ctx).Save(withdrawal).Error
	if err != nil {
		p.logger.Error("Failed to update withdrawal", "error", err, "withdrawalID", withdrawal.ID)
	}
	return err
}

func (p *withdrawalPersistence) GetWithdrawalByID(ctx context.Context, id int64) (*db.Withdrawal, error) {
	var withdrawal db.Withdrawal
	err := p.db.WithContext(ctx).First(&withdrawal, id).Error
	if err != nil {
		p.logger.Error("Failed to get withdrawal by ID", "error", err, "withdrawalID", id)
	}
	return &withdrawal, err
}

func (p *withdrawalPersistence) GetWithdrawalByReference(ctx context.Context, reference string) (*db.Withdrawal, error) {
	var withdrawal db.Withdrawal
	err := p.db.WithContext(ctx).Where("reference = ?", reference).First(&withdrawal).Error
	if err != nil {
		p.logger.Error("Failed to get withdrawal by reference", "error", err, "reference", reference)
	}
	return &withdrawal, err
}

func (p *withdrawalPersistence) GetWithdrawalByTransactionID(ctx context.Context, transactionID string) (*db.Withdrawal, error) {
	var withdrawal db.Withdrawal
	err := p.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&withdrawal).Error
	if err != nil {
		p.logger.Error("Failed to get withdrawal by transaction ID", "error", err, "transactionID", transactionID)
	}
	return &withdrawal, err
}

func (p *withdrawalPersistence) ListWithdrawalsByStoreID(ctx context.Context, storeID int64, limit, offset int) ([]db.Withdrawal, int64, error) {
	var withdrawals []db.Withdrawal
	var total int64

	query := p.db.WithContext(ctx).Model(&db.Withdrawal{}).Where("store_id = ?", storeID)
	if err := query.Count(&total).Error; err != nil {
		p.logger.Error("Failed to count withdrawals", "error", err, "storeID", storeID)
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&withdrawals).Error
	if err != nil {
		p.logger.Error("Failed to list withdrawals", "error", err, "storeID", storeID)
	}
	return withdrawals, total, err
}
