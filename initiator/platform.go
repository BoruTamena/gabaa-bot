package initiator

import (
	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/BoruTamena/gabaa-bot/platform/arifpay"
	"github.com/BoruTamena/gabaa-bot/platform/logger"
	"github.com/BoruTamena/gabaa-bot/platform/rediscache"
	"github.com/BoruTamena/gabaa-bot/platform/cloudinary"
	"github.com/BoruTamena/gabaa-bot/platform/telegram"
	"github.com/spf13/viper"
)

type PlatFormLayer struct {
	tg      platform.Telegram
	cach    platform.Redis
	payment  platform.Payment
	logger   platform.Logger
	uploader platform.FileUploader
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
		uploader: func() platform.FileUploader {
			u, _ := cloudinary.NewCloudinaryService(
				viper.GetString("cloudinary.cloud_name"),
				viper.GetString("cloudinary.api_key"),
				viper.GetString("cloudinary.api_secret"),
			)
			return u
		}(),
	}
}
