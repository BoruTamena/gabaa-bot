package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/module/auth"
	"github.com/BoruTamena/gabaa-bot/internal/module/cart"
	"github.com/BoruTamena/gabaa-bot/internal/module/order"
	"github.com/BoruTamena/gabaa-bot/internal/module/product"
	"github.com/BoruTamena/gabaa-bot/internal/module/store"
	"github.com/BoruTamena/gabaa-bot/internal/module/telegram"
	"github.com/BoruTamena/gabaa-bot/internal/module/user"
	"github.com/BoruTamena/gabaa-bot/internal/module/wallet"
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
}

func InitModule(persistence Persistence, platform PlatFormLayer) Module {
	return Module{
		AuthModule:     auth.NewAuthModule(persistence.UserStorage, persistence.StoreStorage, platform.tg),
		StoreModule:    store.NewStoreModule(persistence.StoreStorage, persistence.UserStorage, platform.tg),
		ProductModule:  product.NewProductModule(persistence.ProductStorage),
		OrderModule:    order.NewOrderModule(persistence.OrderStorage, persistence.ProductStorage, persistence.CartStorage, persistence.WalletStorage),
		CartModule:     cart.NewCartModule(persistence.CartStorage, persistence.ProductStorage),
		WalletModule:   wallet.NewWalletModule(persistence.WalletStorage),
		UserModule:     user.NewUserModule(persistence.UserStorage),
		CategoryModule: product.NewCategoryModule(persistence.CategoryStorage),
		BotModule:      telegram.NewBotModule(persistence.UserStorage, persistence.StoreStorage, platform.tg),
	}
}
