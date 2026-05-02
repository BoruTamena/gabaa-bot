package main

import (
	"github.com/BoruTamena/gabaa-bot/initiator"
	"github.com/BoruTamena/gabaa-bot/pkg/logger"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	// defer logger.Sync()

	// call your initator init
	initiator.Init()

}
