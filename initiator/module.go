package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/module/address"
	"github.com/BoruTamena/gabaa-bot/internal/module/auth"
	"github.com/BoruTamena/gabaa-bot/internal/module/cart"
	"github.com/BoruTamena/gabaa-bot/internal/module/order"
	"github.com/BoruTamena/gabaa-bot/internal/module/product"
	"github.com/BoruTamena/gabaa-bot/internal/module/recommendation"
	"github.com/BoruTamena/gabaa-bot/internal/module/store"
	"github.com/BoruTamena/gabaa-bot/internal/module/telegram"
	"github.com/BoruTamena/gabaa-bot/internal/module/user"
	"github.com/BoruTamena/gabaa-bot/internal/module/upload"
	"github.com/BoruTamena/gabaa-bot/internal/module/wallet"
	"github.com/spf13/viper"
)

type Module struct {
	AuthModule     module.AuthModule
	StoreModule    module.StoreModule
	ProductModule  module.ProductModule
	OrderModule    module.OrderModule
	CartModule     module.CartModule
	WalletModule   module.WalletModule
	UserModule     module.UserModule
	CategoryModule module.CategoryModule
	BotModule      module.BotModule
	UploadModule   module.UploadModule
	AddressModule  module.AddressModule
	StoryModule          module.StoryModule
	FavoriteModule       module.FavoriteModule
	RecommendationModule module.RecommendationModule
}

func InitModule(persistence Persistence, platform PlatFormLayer) Module {
	recommendationModule := recommendation.NewRecommendationModule(
		persistence.PreferenceStorage,
		persistence.RecommendationStorage,
		persistence.UserStorage,
		persistence.StoreStorage,
		platform.tg,
	)

	authModule := auth.NewAuthModule(
		persistence.UserStorage,
		persistence.StoreStorage,
		platform.tg,
		persistence.AuthSessionStorage,
	)

	return Module{
		AuthModule: authModule,
		StoreModule:    store.NewStoreModule(persistence.StoreStorage, persistence.StoreKYCStorage, persistence.UserStorage, platform.tg),
		ProductModule:  product.NewProductModule(persistence.ProductStorage, persistence.StoreStorage, platform.tg, viper.GetString("app.url"), recommendationModule),
		OrderModule:    order.NewOrderModule(persistence.OrderStorage, persistence.ProductStorage, persistence.CartStorage, persistence.WalletStorage, persistence.AddressStorage),
		CartModule:     cart.NewCartModule(persistence.CartStorage, persistence.ProductStorage),
		WalletModule:   wallet.NewWalletModule(persistence.WalletStorage),
		UserModule:     user.NewUserModule(persistence.UserStorage),
		CategoryModule: product.NewCategoryModule(persistence.CategoryStorage),
		BotModule: telegram.NewBotModule(
			persistence.UserStorage,
			persistence.StoreStorage,
			persistence.CategoryStorage,
			recommendationModule,
			authModule,
			platform.tg,
		),
		UploadModule:         upload.NewUploadModule(platform.uploader),
		AddressModule:        address.NewAddressModule(persistence.AddressStorage),
		StoryModule:          product.NewStoryModule(persistence.StoryStorage, persistence.ProductStorage),
		FavoriteModule:       product.NewFavoriteModule(persistence.FavoriteStorage),
		RecommendationModule: recommendationModule,
	}
}

