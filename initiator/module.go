package initiator

import (
	"github.com/BoruTamena/gabaa-bot/internal/module"
	"github.com/BoruTamena/gabaa-bot/internal/module/product"
)

type Module struct {
	productModule module.ProductModule
}

func InitModule(persistence Persistance, platform PlatFormLayer) Module {

	return Module{
		productModule: product.InitProductModule(persistence.productStorage,
			platform.tg),
	}

}
