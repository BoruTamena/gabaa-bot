package initiator

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BoruTamena/gabaa-bot/internal/constant/models/db"
	"github.com/BoruTamena/gabaa-bot/internal/constant/persistencedb"
	"github.com/BoruTamena/gabaa-bot/internal/glue/routing"
	"github.com/spf13/viper"
)

func SeedData(p Persistence) {
	ctx := context.Background()

	// 1. Seed Categories
	defaultCategories := []string{"Electronics", "Fashion", "Home & Garden", "Beauty", "Toys"}
	for _, name := range defaultCategories {
		cat, _ := p.CategoryStorage.GetCategoryByName(ctx, name, 0)
		if cat == nil {
			p.CategoryStorage.CreateCategory(ctx, &db.Category{
				StoreID: 0,
				Name:    name,
			})
		}
	}

	// 2. Seed Default User (Seller)
	sellerID := int64(123456789)
	user, err := p.UserStorage.GetUserByTelegramID(ctx, sellerID)
	if err != nil {
		user = &db.User{
			TelegramUserID: sellerID,
			Username:       "sample_seller",
			Role:           "admin",
		}
		p.UserStorage.CreateUser(ctx, user)
		// Re-fetch to get ID
		user, _ = p.UserStorage.GetUserByTelegramID(ctx, sellerID)
	}

	// 3. Seed Default Store
	store, _ := p.StoreStorage.GetStoreByChatID(ctx, sellerID)
	if store.ID == 0 {
		store = &db.Store{
			SellerID:       user.ID,
			TelegramChatID: sellerID,
			Name:           "Gabaa Sample Store",
			Category:       "Electronics",
			Description:    "A sample store for testing product listings.",
		}
		p.StoreStorage.CreateStore(ctx, store)
	}

	// 4. Seed Products (if store has no products)
	count, _ := p.ProductStorage.GetProductsTotal(ctx, store.ID)
	if count == 0 {
		products := []db.Product{
			{StoreID: store.ID, Name: "Smartphone X", Category: "Electronics", Price: 799.99, Stock: 50, Description: "Latest flagship smartphone", Images: `["https://picsum.photos/id/160/600/400"]`},
			{StoreID: store.ID, Name: "Laptop Pro", Category: "Electronics", Price: 1299.99, Stock: 20, Description: "Powerful laptop for professionals", Images: `["https://picsum.photos/id/119/600/400"]`},
			{StoreID: store.ID, Name: "Wireless Earbuds", Category: "Electronics", Price: 149.99, Stock: 100, Description: "Noise cancelling earbuds", Images: `["https://picsum.photos/id/211/600/400"]`},
			{StoreID: store.ID, Name: "Summer T-Shirt", Category: "Fashion", Price: 19.99, Stock: 200, Description: "Breathable cotton t-shirt", Images: `["https://picsum.photos/id/22/600/400"]`},
			{StoreID: store.ID, Name: "Jeans Slim Fit", Category: "Fashion", Price: 49.99, Stock: 80, Description: "Denim jeans for everyday wear", Images: `["https://picsum.photos/id/23/600/400"]`},
			{StoreID: store.ID, Name: "Running Shoes", Category: "Fashion", Price: 89.99, Stock: 40, Description: "Lightweight running sneakers", Images: `["https://picsum.photos/id/24/600/400"]`},
			{StoreID: store.ID, Name: "Coffee Maker", Category: "Home & Garden", Price: 59.99, Stock: 30, Description: "Brews perfect coffee every time", Images: `["https://picsum.photos/id/25/600/400"]`},
			{StoreID: store.ID, Name: "Garden Tool Set", Category: "Home & Garden", Price: 34.99, Stock: 60, Description: "Essential tools for gardening", Images: `["https://picsum.photos/id/26/600/400"]`},
			{StoreID: store.ID, Name: "Lipstick Matte", Category: "Beauty", Price: 14.99, Stock: 150, Description: "Long lasting matte finish", Images: `["https://picsum.photos/id/27/600/400"]`},
			{StoreID: store.ID, Name: "Action Figure", Category: "Toys", Price: 24.99, Stock: 70, Description: "Collectible action hero toy", Images: `["https://picsum.photos/id/28/600/400"]`},
		}

		for _, prod := range products {
			p.ProductStorage.CreateProduct(ctx, &prod)
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

	// Seed default data
	SeedData(persistenceLayer)
	fmt.Println("Sample data seeded")

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

	port := os.Getenv("PORT")
	if port == "" {
		port = viper.GetString("server.port")
	}
	if port == "" {
		port = "8085"
	}

	fmt.Printf("Starting HTTP server on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
