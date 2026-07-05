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

var ErrAuthSessionNotFound = errors.New("auth session not found or expired")

type authSessionPersistence struct {
	db     *gorm.DB
	logger platform.Logger
}

func NewAuthSessionPersistence(db *gorm.DB, logger platform.Logger) storage.AuthSessionStorage {
	return &authSessionPersistence{db: db, logger: logger}
}

func (p *authSessionPersistence) CreateSession(ctx context.Context, sessionID string, expiresAt time.Time) error {
	session := &db.TelegramLoginSession{
		ID:        sessionID,
		Status:    db.TelegramLoginSessionStatusPending,
		ExpiresAt: expiresAt,
	}
	if err := p.db.WithContext(ctx).Create(session).Error; err != nil {
		p.logger.Error("Failed to create telegram login session", "error", err, "sessionID", sessionID)
		return err
	}
	return nil
}

func (p *authSessionPersistence) CompleteSession(ctx context.Context, sessionID string, telegramUserID int64, username string) error {
	now := time.Now()
	result := p.db.WithContext(ctx).
		Model(&db.TelegramLoginSession{}).
		Where("id = ? AND status = ? AND expires_at > ?", sessionID, db.TelegramLoginSessionStatusPending, now).
		Updates(map[string]interface{}{
			"status":           db.TelegramLoginSessionStatusCompleted,
			"telegram_user_id": telegramUserID,
			"username":         username,
			"completed_at":     now,
		})
	if result.Error != nil {
		p.logger.Error("Failed to complete telegram login session", "error", result.Error, "sessionID", sessionID)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrAuthSessionNotFound
	}
	return nil
}

func (p *authSessionPersistence) GetSession(ctx context.Context, sessionID string) (*storage.AuthSession, error) {
	var session db.TelegramLoginSession
	err := p.db.WithContext(ctx).Where("id = ?", sessionID).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuthSessionNotFound
		}
		p.logger.Error("Failed to get telegram login session", "error", err, "sessionID", sessionID)
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, ErrAuthSessionNotFound
	}

	authSession := &storage.AuthSession{
		ID:        session.ID,
		Status:    session.Status,
		Username:  session.Username,
		ExpiresAt: session.ExpiresAt,
	}
	if session.TelegramUserID != nil {
		authSession.TelegramUserID = *session.TelegramUserID
	}
	if session.CompletedAt != nil {
		completedAt := *session.CompletedAt
		authSession.CompletedAt = &completedAt
	}
	return authSession, nil
}

func (p *authSessionPersistence) DeleteSession(ctx context.Context, sessionID string) error {
	result := p.db.WithContext(ctx).Where("id = ?", sessionID).Delete(&db.TelegramLoginSession{})
	if result.Error != nil {
		p.logger.Error("Failed to delete telegram login session", "error", result.Error, "sessionID", sessionID)
		return result.Error
	}
	return nil
}
