package main

import (
	"log"
	"net/http"

	"github.com/Xanaduxan/tasks-golang/config"

	"github.com/Xanaduxan/tasks-golang/internal/queue"
	"github.com/Xanaduxan/tasks-golang/internal/service/auth"
	"github.com/Xanaduxan/tasks-golang/internal/service/deliveries"
	"github.com/Xanaduxan/tasks-golang/internal/service/group_members"
	"github.com/Xanaduxan/tasks-golang/internal/service/groups"
	"github.com/Xanaduxan/tasks-golang/internal/service/products"
	"github.com/Xanaduxan/tasks-golang/internal/service/stocks"
	"github.com/Xanaduxan/tasks-golang/internal/service/tasks"
	"github.com/Xanaduxan/tasks-golang/internal/storage"
	handlers "github.com/Xanaduxan/tasks-golang/internal/transport/http-handlers"
	"github.com/Xanaduxan/tasks-golang/internal/transport/router"
	"github.com/Xanaduxan/tasks-golang/internal/transport/websocket"
	"github.com/Xanaduxan/tasks-golang/internal/worker"
	_ "github.com/Xanaduxan/tasks-golang/metrics"
	redispkg "github.com/Xanaduxan/tasks-golang/pkg/redis"
)

func main() {
	cfg := config.MustLoad()

	db := storage.NewPostgres(cfg.DatabaseURL)

	redisClient, err := redispkg.NewRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Printf("redis is unavailable, starting without cache: %v", err)
	} else {
		defer func() {
			if err := redispkg.Close(redisClient); err != nil {
				log.Printf("failed to close redis: %v", err)
			}
		}()
		log.Println("connected to redis")
	}

	userStorage := storage.NewUserStorage(db)
	taskStorage := storage.NewTaskStorage(db)
	taskCached := storage.NewTaskCached(taskStorage, redisClient)
	productStorage := storage.NewProductStorage(db)
	stockStorage := storage.NewStockStorage(db)
	deliveryStorage := storage.NewDeliveryStorage(db)
	deliveryItemsStorage := storage.NewDeliveryItemStorage(db)
	groupsStorage := storage.NewGroupStorage(db)
	groupMemberStorage := storage.NewGroupMemberStorage(db)

	wsManager := websocket.NewManager()
	wsNotifier := websocket.NewNotifier(wsManager, groupMemberStorage)

	authService := auth.NewService(userStorage, []byte(cfg.JWTSecret))

	tasksService := tasks.NewService(
		taskCached,
		userStorage,
		groupsStorage,
		groupMemberStorage,
		wsNotifier,
	)
	tasksService.InitMetrics()
	productService := products.NewService(productStorage)
	stocksService := stocks.NewService(stockStorage)
	groupsService := groups.NewGroupService(groupsStorage)
	groupMemberService := group_members.NewGroupMemberService(groupMemberStorage, groupsStorage, userStorage)

	deliveryService := deliveries.NewService(
		productStorage,
		userStorage,
		deliveryStorage,
		deliveryItemsStorage,
		stockStorage,
		wsNotifier,
	)

	deliveryQueue := queue.NewDeliveryQueue()
	deliveryWorker := worker.NewDeliveryWorker(deliveryQueue, deliveryService)
	deliveryWorker.Start()

	taskQueue := queue.NewTaskQueue()
	taskWorker := worker.NewTaskWorker(taskQueue, tasksService)
	taskWorker.Start()

	handlers.SetAuthService(authService)
	handlers.SetTaskService(tasksService)
	handlers.SetProductService(productService)
	handlers.SetStockService(stocksService)
	handlers.SetDeliveryService(deliveryService)
	handlers.SetDeliveryQueue(deliveryQueue)
	handlers.SetGroupService(groupsService)
	handlers.SetGroupMemberService(groupMemberService)

	wsHandler := websocket.NewHandler(wsManager)
	r := router.New([]byte(cfg.JWTSecret), wsHandler)

	addr := ":" + cfg.HTTPPort
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
