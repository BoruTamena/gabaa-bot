package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/module/auth"
	"github.com/BoruTamena/gabaa-bot/internal/module/order"
	"github.com/BoruTamena/gabaa-bot/internal/module/product"
	"github.com/BoruTamena/gabaa-bot/internal/module/store"
	"github.com/BoruTamena/gabaa-bot/internal/module/user"
	"github.com/BoruTamena/gabaa-bot/internal/module/wallet"
)

type Module struct {
	AuthModule    module.AuthModule
	StoreModule   module.StoreModule
	ProductModule module.ProductModule
	OrderModule   module.OrderModule
	WalletModule  module.WalletModule
	UserModule    module.UserModule
}

func InitModule(persistence Persistence, platform PlatFormLayer) Module {
	return Module{
		AuthModule:    auth.NewAuthModule(persistence.UserStorage, persistence.StoreStorage, platform.tg),
		StoreModule:   store.NewStoreModule(persistence.StoreStorage, platform.tg),
		ProductModule: product.NewProductModule(persistence.ProductStorage, platform.logger),
		OrderModule:   order.NewOrderModule(persistence.OrderStorage, persistence.ProductStorage, persistence.CartStorage, persistence.WalletStorage),
		WalletModule:  wallet.NewWalletModule(persistence.WalletStorage),
		UserModule:    user.NewUserModule(persistence.UserStorage),
	}
}
