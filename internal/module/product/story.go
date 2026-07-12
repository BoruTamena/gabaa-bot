package product

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/errorx"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"go.uber.org/zap"
)

type storyModule struct {
	storyStorage   storage.StoryStorage
	productStorage storage.ProductStorage
}

// NewStoryModule creates a new StoryModule.
func NewStoryModule(ss storage.StoryStorage, ps storage.ProductStorage) module.StoryModule {
	return &storyModule{
		storyStorage:   ss,
		productStorage: ps,
	}
}

// CreateStory validates and creates a new story ad for the given store.
func (m *storyModule) CreateStory(ctx context.Context, storeID int64, req dto.CreateProductStoryRequest) (*dto.ProductStory, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Parse date range
	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		return nil, fmt.Errorf("invalid starts_at format, expected RFC3339: %w", err)
	}
	endsAt, err := time.Parse(time.RFC3339, req.EndsAt)
	if err != nil {
		return nil, fmt.Errorf("invalid ends_at format, expected RFC3339: %w", err)
	}
	if !endsAt.After(startsAt) {
		return nil, fmt.Errorf("ends_at must be after starts_at")
	}
	if endsAt.Before(time.Now()) {
		return nil, fmt.Errorf("ends_at must be in the future")
	}

	// Verify product belongs to the store
	product, err := m.productStorage.GetProductByID(ctx, req.ProductID)
	if err != nil || product.StoreID == nil || *product.StoreID != storeID {
		return nil, errorx.New(errorx.ErrForbidden, "product does not belong to your store", http.StatusForbidden)
	}

	mediaBytes, _ := json.Marshal(req.MediaURLs)
	story := &db.ProductStory{
		StoreID:   storeID,
		ProductID: req.ProductID,
		Caption:   req.Caption,
		MediaURLs: string(mediaBytes),
		MediaType: req.MediaType,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
		IsActive:  true,
	}

	if err := m.storyStorage.CreateStory(ctx, story); err != nil {
		logger.Error("failed to create story", zap.Error(err))
		return nil, err
	}

	logger.Info("story created successfully", zap.Int64("story_id", story.ID), zap.Int64("store_id", storeID))
	return m.mapToDTO(story, nil), nil
}

// GetStory fetches a single story with its associated product, and increments views asynchronously.
func (m *storyModule) GetStory(ctx context.Context, id int64) (*dto.ProductStory, error) {
	story, err := m.storyStorage.GetStoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Async view increment — fire and forget
	go func() {
		if viewErr := m.storyStorage.IncrementStoryViews(context.Background(), id); viewErr != nil {
			logger.Error("failed to increment story views", zap.Error(viewErr), zap.Int64("story_id", id))
		}
	}()

	var productDTO *dto.Product
	if story.Product != nil {
		p := story.Product
		var images []string
		if p.Images != "" {
			_ = json.Unmarshal([]byte(p.Images), &images)
		}
		if images == nil {
			images = []string{}
		}
		productDTO = &dto.Product{
			ID:          p.ID,
			StoreID:     p.StoreID,
			SellerID:    p.SellerID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			Category:    p.Category,
			Images:      images,
			Status:      p.Status,
			IsPosted:    p.IsPosted,
			IsBoosted:   p.IsBoosted,
		}
	}

	return m.mapToDTO(story, productDTO), nil
}

// ListMyStories returns paginated stories scoped to the store owner.
func (m *storyModule) ListMyStories(ctx context.Context, filter dto.ProductStoryFilterParams) (*dto.PaginatedResponse, error) {
	stories, total, err := m.storyStorage.ListStoriesByStore(ctx, filter)
	if err != nil {
		return nil, err
	}

	dtoStories := make([]dto.ProductStory, len(stories))
	for i, s := range stories {
		dtoStories[i] = *m.mapToDTO(&s, nil)
	}

	return &dto.PaginatedResponse{
		Total: total,
		Data:  dtoStories,
	}, nil
}

// UpdateStory partially updates a story, enforcing store ownership.
func (m *storyModule) UpdateStory(ctx context.Context, storeID int64, storyID int64, req dto.UpdateProductStoryRequest) (*dto.ProductStory, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	story, err := m.storyStorage.GetStoryByID(ctx, storyID)
	if err != nil {
		return nil, err
	}

	if story.StoreID != storeID {
		return nil, errorx.New(errorx.ErrForbidden, "story does not belong to your store", http.StatusForbidden)
	}

	// Apply partial updates
	if req.Caption != "" {
		story.Caption = req.Caption
	}
	if len(req.MediaURLs) > 0 {
		mediaBytes, _ := json.Marshal(req.MediaURLs)
		story.MediaURLs = string(mediaBytes)
	}
	if req.MediaType != "" {
		story.MediaType = req.MediaType
	}
	if req.StartsAt != "" {
		t, err := time.Parse(time.RFC3339, req.StartsAt)
		if err != nil {
			return nil, fmt.Errorf("invalid starts_at format: %w", err)
		}
		story.StartsAt = t
	}
	if req.EndsAt != "" {
		t, err := time.Parse(time.RFC3339, req.EndsAt)
		if err != nil {
			return nil, fmt.Errorf("invalid ends_at format: %w", err)
		}
		story.EndsAt = t
	}
	if req.IsActive != nil {
		story.IsActive = *req.IsActive
	}

	// Validate date range after update
	if !story.EndsAt.After(story.StartsAt) {
		return nil, fmt.Errorf("ends_at must be after starts_at")
	}

	if err := m.storyStorage.UpdateStory(ctx, story); err != nil {
		logger.Error("failed to update story", zap.Error(err), zap.Int64("story_id", storyID))
		return nil, err
	}

	logger.Info("story updated successfully", zap.Int64("story_id", storyID))
	return m.mapToDTO(story, nil), nil
}

// DeleteStory soft-deletes a story, enforcing store ownership.
func (m *storyModule) DeleteStory(ctx context.Context, storeID int64, storyID int64) error {
	story, err := m.storyStorage.GetStoryByID(ctx, storyID)
	if err != nil {
		return err
	}
	if story.StoreID != storeID {
		return errorx.New(errorx.ErrForbidden, "story does not belong to your store", http.StatusForbidden)
	}

	if err := m.storyStorage.DeleteStory(ctx, storyID); err != nil {
		logger.Error("failed to delete story", zap.Error(err), zap.Int64("story_id", storyID))
		return err
	}

	logger.Info("story deleted successfully", zap.Int64("story_id", storyID))
	return nil
}

// StartExpiryJob runs a background scheduler that deactivates stories whose ends_at has passed.
// It runs once on startup, then twice daily at midnight and noon (server local time).
func (m *storyModule) StartExpiryJob(ctx context.Context) {
	go func() {
		m.runStoryExpiry(ctx)

		for {
			wait := durationUntilNextStoryExpiryRun(time.Now())
			timer := time.NewTimer(wait)
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
				m.runStoryExpiry(ctx)
			}
		}
	}()
}

func durationUntilNextStoryExpiryRun(now time.Time) time.Duration {
	loc := now.Location()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	noon := midnight.Add(12 * time.Hour)
	midnightNext := midnight.Add(24 * time.Hour)

	var next time.Time
	if now.Before(noon) {
		next = noon
	} else {
		next = midnightNext
	}
	return next.Sub(now)
}

func (m *storyModule) runStoryExpiry(ctx context.Context) {
	count, err := m.storyStorage.ExpireEndedStories(ctx)
	if err != nil {
		logger.Error("story expiry job failed", zap.Error(err))
		return
	}
	if count > 0 {
		logger.Info("expired stories deactivated", zap.Int64("count", count))
	}
}

// ListActiveStories returns publicly visible, currently-active stories with pagination.
func (m *storyModule) ListActiveStories(ctx context.Context, params dto.PaginationParams) (*dto.PaginatedResponse, error) {
	stories, total, err := m.storyStorage.ListActiveStories(ctx, params)
	if err != nil {
		return nil, err
	}

	dtoStories := make([]dto.ProductStory, len(stories))
	for i, s := range stories {
		dtoStories[i] = *m.mapToDTO(&s, nil)
	}

	return &dto.PaginatedResponse{
		Total: total,
		Data:  dtoStories,
	}, nil
}

// mapToDTO converts a db.ProductStory to dto.ProductStory.
func (m *storyModule) mapToDTO(s *db.ProductStory, product *dto.Product) *dto.ProductStory {
	var mediaURLs []string
	if s.MediaURLs != "" {
		_ = json.Unmarshal([]byte(s.MediaURLs), &mediaURLs)
	}
	if mediaURLs == nil {
		mediaURLs = []string{}
	}

	return &dto.ProductStory{
		ID:        s.ID,
		StoreID:   s.StoreID,
		ProductID: s.ProductID,
		Caption:   s.Caption,
		MediaURLs: mediaURLs,
		MediaType: s.MediaType,
		StartsAt:  s.StartsAt.Format(time.RFC3339),
		EndsAt:    s.EndsAt.Format(time.RFC3339),
		IsActive:  s.IsActive,
		Views:     s.Views,
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
		Product:   product,
	}
}
