package initiator

import (
	"fmt"
	"log"

	"github.com/BoruTamena/gabaa-bot/internal/constant/persistencedb"
	"github.com/BoruTamena/gabaa-bot/internal/glue/routing"
	"github.com/spf13/viper"
)

func Init() {
	err := InitViper("./")
	if err != nil {
		panic(err)
	}

	dbPersistence := persistencedb.NewPersistenceDb()

	// Migration setup
	mg := InitMigiration(viper.GetString("migiration.path"), viper.GetString("db.url"))
	UpMigiration(mg)
	fmt.Println("Migrations applied")

	platformLayer := InitPlatFormLayer()

	// Persistence layer needs DB and Redis
	persistenceLayer := InitPersistence(dbPersistence, platformLayer.cach, platformLayer.logger)

	moduleLayer := InitModule(persistenceLayer, platformLayer)

	handlerLayer := InitHandler(moduleLayer, platformLayer)

	// Start bot in goroutine
	go platformLayer.tg.Start()
	fmt.Println("Telegram bot started")

	// Initialize Gin Router
	r := routing.NewGinRouter(
		handlerLayer.AuthHandler,
		handlerLayer.StoreHandler,
		handlerLayer.ProductHandler,
		handlerLayer.OrderHandler,
		handlerLayer.PaymentHandler,
		handlerLayer.AuthMiddleware,
	)

	port := viper.GetString("server.port")
	if port == "" {
		port = "8085"
	}

	fmt.Printf("Starting HTTP server on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
