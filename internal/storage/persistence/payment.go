package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/platform"
	"gorm.io/gorm"
)

type paymentPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewPaymentPersistence(db *gorm.DB, logger platform.Logger) storage.PaymentStorage {
	return &paymentPersistence{db: db, logger: logger}
}

func (p *paymentPersistence) CreatePayment(ctx context.Context, payment *db.Payment) error {
	err := p.db.WithContext(ctx).Create(payment).Error
	if err != nil {
		p.logger.Error("Failed to create payment", "error", err)
	}
	return err
}

func (p *paymentPersistence) UpdatePayment(ctx context.Context, payment *db.Payment) error {
	err := p.db.WithContext(ctx).Save(payment).Error
	if err != nil {
		p.logger.Error("Failed to update payment", "error", err, "paymentID", payment.ID)
	}
	return err
}

func (p *paymentPersistence) GetPaymentByID(ctx context.Context, id int64) (*db.Payment, error) {
	var payment db.Payment
	err := p.db.WithContext(ctx).Preload("Order").First(&payment, id).Error
	if err != nil {
		p.logger.Error("Failed to get payment by ID", "error", err, "paymentID", id)
	}
	return &payment, err
}

func (p *paymentPersistence) GetPaymentByReference(ctx context.Context, reference string) (*db.Payment, error) {
	var payment db.Payment
	err := p.db.WithContext(ctx).Preload("Order").Where("reference = ?", reference).First(&payment).Error
	if err != nil {
		p.logger.Error("Failed to get payment by reference", "error", err, "reference", reference)
	}
	return &payment, err
}

func (p *paymentPersistence) GetPaymentByTransactionID(ctx context.Context, transactionID string) (*db.Payment, error) {
	var payment db.Payment
	err := p.db.WithContext(ctx).Preload("Order").Where("transaction_id = ?", transactionID).First(&payment).Error
	if err != nil {
		p.logger.Error("Failed to get payment by transaction ID", "error", err, "transactionID", transactionID)
	}
	return &payment, err
}

type paymentWebhookPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewPaymentWebhookPersistence(db *gorm.DB, logger platform.Logger) storage.PaymentWebhookStorage {
	return &paymentWebhookPersistence{db: db, logger: logger}
}

func (p *paymentWebhookPersistence) CreateWebhookEvent(ctx context.Context, event *db.PaymentWebhook) error {
	if event.ReceivedAt.IsZero() {
		event.ReceivedAt = time.Now()
	}
	err := p.db.WithContext(ctx).Create(event).Error
	if err != nil {
		p.logger.Error("Failed to create payment webhook event", "error", err)
	}
	return err
}

func (p *paymentWebhookPersistence) MarkWebhookProcessed(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Model(&db.PaymentWebhook{}).Where("id = ?", id).Update("processed", true).Error
	if err != nil {
		p.logger.Error("Failed to mark webhook processed", "error", err, "webhookID", id)
	}
	return err
}

type escrowPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewEscrowPersistence(db *gorm.DB, logger platform.Logger) storage.EscrowStorage {
	return &escrowPersistence{db: db, logger: logger}
}

func (p *escrowPersistence) CreateEscrow(ctx context.Context, escrow *db.Escrow) error {
	if escrow.CreatedAt.IsZero() {
		escrow.CreatedAt = time.Now()
	}
	err := p.db.WithContext(ctx).Create(escrow).Error
	if err != nil {
		p.logger.Error("Failed to create escrow", "error", err)
	}
	return err
}

func (p *escrowPersistence) GetEscrowByOrderID(ctx context.Context, orderID int64) (*db.Escrow, error) {
	var escrow db.Escrow
	err := p.db.WithContext(ctx).Where("order_id = ?", orderID).First(&escrow).Error
	if err != nil {
		p.logger.Error("Failed to get escrow by order ID", "error", err, "orderID", orderID)
	}
	return &escrow, err
}

func (p *escrowPersistence) ReleaseEscrow(ctx context.Context, orderID int64) error {
	now := time.Now()
	result := p.db.WithContext(ctx).Model(&db.Escrow{}).
		Where("order_id = ? AND status = ?", orderID, constant.EscrowStatusHeld).
		Updates(map[string]interface{}{
			"status":      constant.EscrowStatusReleased,
			"released_at": now,
		})
	if result.Error != nil {
		p.logger.Error("Failed to release escrow", "error", result.Error, "orderID", orderID)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("escrow not found or already released")
	}
	return nil
}
