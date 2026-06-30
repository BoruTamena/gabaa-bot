package product

import (
	"context"
	"encoding/json"
	"time"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/models/dto"
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
	"go.uber.org/zap"
)

type favoriteModule struct {
	favoriteStorage storage.FavoriteStorage
}

// NewFavoriteModule creates a new FavoriteModule.
func NewFavoriteModule(fs storage.FavoriteStorage) module.FavoriteModule {
	return &favoriteModule{
		favoriteStorage: fs,
	}
}

func (m *favoriteModule) AddFavorite(ctx context.Context, userID, productID int64) error {
	// Check if already favorited to make it idempotent
	isFav, err := m.favoriteStorage.IsFavorite(ctx, userID, productID)
	if err != nil {
		return err
	}
	if isFav {
		return nil // Already a favorite
	}

	fav := &db.Favorite{
		UserID:    userID,
		ProductID: productID,
	}

	if err := m.favoriteStorage.AddFavorite(ctx, fav); err != nil {
		logger.Error("failed to add favorite", zap.Error(err), zap.Int64("user_id", userID), zap.Int64("product_id", productID))
		return err
	}

	return nil
}

func (m *favoriteModule) RemoveFavorite(ctx context.Context, userID, productID int64) error {
	if err := m.favoriteStorage.RemoveFavorite(ctx, userID, productID); err != nil {
		logger.Error("failed to remove favorite", zap.Error(err), zap.Int64("user_id", userID), zap.Int64("product_id", productID))
		return err
	}
	return nil
}

func (m *favoriteModule) ListUserFavorites(ctx context.Context, userID int64, params dto.PaginationParams) (*dto.PaginatedResponse, error) {
	favorites, total, err := m.favoriteStorage.ListUserFavorites(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	dtoFavorites := make([]dto.FavoriteResponse, len(favorites))
	for i, f := range favorites {
		var productDTO *dto.Product
		if f.Product != nil {
			p := f.Product
			var images []string
			if p.Images != "" {
				_ = json.Unmarshal([]byte(p.Images), &images)
			}
			if images == nil {
				images = []string{}
			}
			productDTO = &dto.Product{
				ID:          p.ID,
				SellerID:    p.SellerID,
				StoreID:     p.StoreID,
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

		dtoFavorites[i] = dto.FavoriteResponse{
			ID:        f.ID,
			UserID:    f.UserID,
			ProductID: f.ProductID,
			CreatedAt: f.CreatedAt.Format(time.RFC3339),
			Product:   productDTO,
		}
	}

	return &dto.PaginatedResponse{
		Total: total,
		Data:  dtoFavorites,
	}, nil
}
