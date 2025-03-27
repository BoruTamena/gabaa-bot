package initiator

import (
	"github.com/BoruTamena/gabaa-bot/platform"
	"github.com/BoruTamena/gabaa-bot/platform/telegram"
)

type PlatFormLayer struct {
	tg platform.Telegram
}

func InitPlatFormLayer() PlatFormLayer {

	return PlatFormLayer{

		tg: telegram.InitTelBot(),
	}
}
