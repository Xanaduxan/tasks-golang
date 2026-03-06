package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Xanaduxan/tasks-golang/internal/config"
	"github.com/Xanaduxan/tasks-golang/internal/http/handlers"
	"github.com/Xanaduxan/tasks-golang/internal/http/router"
	"github.com/Xanaduxan/tasks-golang/internal/queue"
	"github.com/Xanaduxan/tasks-golang/internal/service/auth"
	"github.com/Xanaduxan/tasks-golang/internal/service/deliveries"
	"github.com/Xanaduxan/tasks-golang/internal/service/products"
	"github.com/Xanaduxan/tasks-golang/internal/service/stocks"
	"github.com/Xanaduxan/tasks-golang/internal/service/tasks"
	"github.com/Xanaduxan/tasks-golang/internal/storage"
	"github.com/Xanaduxan/tasks-golang/internal/worker"
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
	taskStorage := storage.NewTaskStorage(db)
	productStorage := storage.NewProductStorage(db)
	stockStorage := storage.NewStockStorage(db)
	deliveryStorage := storage.NewDeliveryStorage(db)
	deliveryItemsStorage := storage.NewDeliveryItemStorage(db)

	authService := auth.NewService(userStorage, []byte(jwtSecret))
	tasksService := tasks.NewService(taskStorage, userStorage)
	productService := products.NewService(productStorage)
	stocksService := stocks.NewService(stockStorage)
	deliveryService := deliveries.NewService(
		productStorage,
		userStorage,
		deliveryStorage,
		deliveryItemsStorage,
		stockStorage,
	)

	deliveryQueue := queue.NewDeliveryQueue()
	deliveryWorker := worker.NewDeliveryWorker(deliveryQueue, deliveryService)
	deliveryWorker.Start()

	handlers.SetAuthService(authService)
	handlers.SetTaskService(tasksService)
	handlers.SetProductService(productService)
	handlers.SetStockService(stocksService)
	handlers.SetDeliveryService(deliveryService)
	handlers.SetDeliveryQueue(deliveryQueue)

	r := router.New([]byte(jwtSecret))
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
