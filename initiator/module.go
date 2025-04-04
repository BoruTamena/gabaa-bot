package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/module/order"
	"github.com/BoruTamena/gabaa-bot/internal/module/product"
	"github.com/BoruTamena/gabaa-bot/internal/module/user"
)

type Module struct {
	productModule module.ProductModule
	orderModule   module.OrderModule
	userModule    module.UserModule
}

func InitModule(persistence Persistance, platform PlatFormLayer) Module {

	return Module{
		userModule: user.InitUserModule(persistence.userStorage),
		productModule: product.InitProductModule(persistence.productStorage,
			platform.tg),

		orderModule: order.InitOrderModule(persistence.productStorage,
			persistence.orderStorage,
			platform.cach),
	}

}
