package initiator

import (
	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/BoruTamena/gabaa-bot/platform/arifpay"
	"github.com/BoruTamena/gabaa-bot/platform/logger"
	"github.com/BoruTamena/gabaa-bot/platform/rediscache"
	"github.com/BoruTamena/gabaa-bot/platform/telegram"
	"github.com/spf13/viper"
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
				Addr:     viper.GetString("redis.addr"),
				Password: viper.GetString("redis.password"),
				DB:       viper.GetInt("redis.db"),
			},
		),

		payment: arifpay.NewPayment(),
		logger:  logger.NewZapLogger(),
	}
}
