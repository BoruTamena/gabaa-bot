package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/module/address"
	"github.com/BoruTamena/gabaa-bot/internal/module/analytics"
	"github.com/BoruTamena/gabaa-bot/internal/module/auth"
	"github.com/BoruTamena/gabaa-bot/internal/module/cart"
	"github.com/BoruTamena/gabaa-bot/internal/module/delivery"
	"github.com/BoruTamena/gabaa-bot/internal/module/order"
	"github.com/BoruTamena/gabaa-bot/internal/module/payment"
	"github.com/BoruTamena/gabaa-bot/internal/module/product"
	"github.com/BoruTamena/gabaa-bot/internal/module/recommendation"
	"github.com/BoruTamena/gabaa-bot/internal/module/store"
	"github.com/BoruTamena/gabaa-bot/internal/module/telegram"
	"github.com/BoruTamena/gabaa-bot/internal/module/upload"
	"github.com/BoruTamena/gabaa-bot/internal/module/user"
	"github.com/BoruTamena/gabaa-bot/internal/module/wallet"
	"github.com/spf13/viper"
)

type Module struct {
	AuthModule           module.AuthModule
	StoreModule          module.StoreModule
	ProductModule        module.ProductModule
	OrderModule          module.OrderModule
	PaymentModule        module.PaymentModule
	CartModule           module.CartModule
	WalletModule         module.WalletModule
	UserModule           module.UserModule
	CategoryModule       module.CategoryModule
	BotModule            module.BotModule
	UploadModule         module.UploadModule
	AddressModule        module.AddressModule
	StoryModule          module.StoryModule
	FavoriteModule       module.FavoriteModule
	RecommendationModule module.RecommendationModule
	AnalyticsModule      module.AnalyticsModule
	DeliveryModule       module.DeliveryModule
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
		persistence.DeliveryStorage,
		platform.tg,
		persistence.AuthSessionStorage,
	)

	orderMod := order.NewOrderModule(
		persistence.OrderStorage,
		persistence.ProductStorage,
		persistence.CartStorage,
		persistence.WalletStorage,
		persistence.EscrowStorage,
		persistence.AddressStorage,
		persistence.StoreStorage,
		persistence.UserStorage,
		platform.tg,
	)

	deliveryMod := delivery.NewDeliveryModule(
		persistence.DeliveryStorage,
		persistence.OrderStorage,
		persistence.StoreStorage,
		persistence.EscrowStorage,
		persistence.WalletStorage,
		platform.tg,
	)
	deliveryMod.SetOrderModule(orderMod)
	orderMod.SetDeliveryModule(deliveryMod)

	paymentMod := payment.NewPaymentModule(
		persistence.PaymentStorage,
		persistence.PaymentWebhookStorage,
		persistence.EscrowStorage,
		persistence.WalletStorage,
		persistence.WithdrawalStorage,
		persistence.OrderStorage,
		platform.lakipay,
		orderMod,
	)
	orderMod.SetPaymentModule(paymentMod)

	return Module{
		AuthModule:    authModule,
		StoreModule:   store.NewStoreModule(persistence.StoreStorage, persistence.StoreKYCStorage, persistence.UserStorage, platform.tg),
		ProductModule: product.NewProductModule(persistence.ProductStorage, persistence.StoreStorage, platform.tg, viper.GetString("app.url"), recommendationModule),
		OrderModule:   orderMod,
		PaymentModule: paymentMod,
		CartModule:    cart.NewCartModule(persistence.CartStorage, persistence.ProductStorage),
		WalletModule: wallet.NewWalletModule(
			persistence.WalletStorage,
			persistence.WithdrawalStorage,
			persistence.StoreStorage,
			platform.lakipay,
		),
		UserModule:     user.NewUserModule(persistence.UserStorage),
		CategoryModule: product.NewCategoryModule(persistence.CategoryStorage),
		BotModule: telegram.NewBotModule(
			persistence.UserStorage,
			persistence.StoreStorage,
			persistence.CategoryStorage,
			recommendationModule,
			authModule,
			deliveryMod,
			platform.tg,
		),
		UploadModule:         upload.NewUploadModule(platform.uploader),
		AddressModule:        address.NewAddressModule(persistence.AddressStorage),
		StoryModule:          product.NewStoryModule(persistence.StoryStorage, persistence.ProductStorage),
		FavoriteModule:       product.NewFavoriteModule(persistence.FavoriteStorage),
		RecommendationModule: recommendationModule,
		AnalyticsModule:      analytics.NewAnalyticsModule(persistence.AnalyticsStorage),
		DeliveryModule:       deliveryMod,
	}
}
