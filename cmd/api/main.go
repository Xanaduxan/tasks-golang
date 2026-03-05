package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Xanaduxan/tasks-golang/internal/config"
	"github.com/Xanaduxan/tasks-golang/internal/http/handlers"
	"github.com/Xanaduxan/tasks-golang/internal/http/router"
	"github.com/Xanaduxan/tasks-golang/internal/service/auth"
	"github.com/Xanaduxan/tasks-golang/internal/service/deliveries"
	"github.com/Xanaduxan/tasks-golang/internal/service/products"
	"github.com/Xanaduxan/tasks-golang/internal/service/stocks"
	"github.com/Xanaduxan/tasks-golang/internal/service/tasks"
	"github.com/Xanaduxan/tasks-golang/internal/storage"
)

func main() {
	config.LoadEnv(".env")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	db := storage.NewPostgres(dsn)
	userStorage := storage.NewUserStorage(db)

	authService := auth.NewService(userStorage, []byte(jwtSecret))
	handlers.SetAuthService(authService)
	taskStorage := storage.NewTaskStorage(db)
	tasksService := tasks.NewService(taskStorage, userStorage)
	handlers.SetTaskService(tasksService)

	productStorage := storage.NewProductStorage(db)
	productService := products.NewService(productStorage)
	handlers.SetProductService(productService)

	stockStorage := storage.NewStockStorage(db)
	stocksService := stocks.NewService(stockStorage)
	handlers.SetStockService(stocksService)

	deliveryStorage := storage.NewDeliveryStorage(db)
	deliveryItemsStorage := storage.NewDeliveryItemStorage(db)
	deliveryService := deliveries.NewService(productStorage, userStorage, deliveryStorage, deliveryItemsStorage, stockStorage)
	handlers.SetDeliveryService(deliveryService)

	r := router.New([]byte(jwtSecret))
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
