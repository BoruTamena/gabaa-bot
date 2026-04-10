package initiator

import (
	"context"
	"fmt"
	"log"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/persistencedb"
	"github.com/BoruTamena/gabaa-bot/internal/glue/routing"
	"github.com/BoruTamena/gabaa-bot/internal/storage"
	"github.com/spf13/viper"
)

func SeedCategories(s storage.CategoryStorage) {
	ctx := context.Background()
	defaultCategories := []string{"Electronics", "Fashion", "Home & Garden", "Beauty", "Toys"}

	for _, name := range defaultCategories {
		// Check if exists (storeID 0 means global)
		cat, _ := s.GetCategoryByName(ctx, name, 0)
		if cat == nil {
			s.CreateCategory(ctx, &db.Category{
				StoreID: 0,
				Name:    name,
			})
		}
	}
}

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

	// Seed default categories
	SeedCategories(persistenceLayer.CategoryStorage)
	fmt.Println("Default categories seeded")

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
		handlerLayer.CategoryHandler,
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
