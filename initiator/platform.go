package initiator

import (
	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/BoruTamena/gabaa-bot/platform/arifpay"
	"github.com/BoruTamena/gabaa-bot/platform/logger"
	"github.com/BoruTamena/gabaa-bot/platform/rediscache"
	"github.com/BoruTamena/gabaa-bot/platform/telegram"
)

type PlatFormLayer struct {
	tg      platform.Telegram
	cach    platform.Redis
	payment platform.Payment
	logger  platform.Logger
}

func InitPlatFormLayer() PlatFormLayer {

	return PlatFormLayer{

		tg: telegram.InitTelBot(),
		cach: rediscache.NewRedis(
			&rediscache.RedisClient{
				Addr:     "localhost:7979",
				Password: "",
				DB:       0,
			},
		),

		payment: arifpay.NewPayment(),
		logger:  logger.NewZapLogger(),
	}
}
