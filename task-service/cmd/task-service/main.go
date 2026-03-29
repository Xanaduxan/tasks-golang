package main

import (
	"log"
	"net/http"

	"github.com/Xanaduxan/tasks-golang/task-service/config"
	queue2 "github.com/Xanaduxan/tasks-golang/task-service/internal/queue"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/auth"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/deliveries"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/group_members"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/groups"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/products"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/stocks"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/tasks"
	storage2 "github.com/Xanaduxan/tasks-golang/task-service/internal/storage"
	grpcserver "github.com/Xanaduxan/tasks-golang/task-service/internal/transport/grpc"
	grpcclient "github.com/Xanaduxan/tasks-golang/task-service/internal/transport/grpc/client"
	handlers "github.com/Xanaduxan/tasks-golang/task-service/internal/transport/http-handlers"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/transport/router"
	worker2 "github.com/Xanaduxan/tasks-golang/task-service/internal/worker"
	redispkg "github.com/Xanaduxan/tasks-golang/task-service/pkg/redis"

	_ "github.com/Xanaduxan/tasks-golang/task-service/metrics"
)

func main() {
	cfg := config.MustLoad()

	db := storage2.NewPostgres(cfg.DatabaseURL)

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

	notificationClient, err := grpcclient.NewNotificationClient(cfg.NotificationAddr)
	if err != nil {
		log.Printf("notification service unavailable, starting without notifications: %v", err)
	} else {
		defer func() {
			if err := notificationClient.Close(); err != nil {
				log.Printf("failed to close notification client: %v", err)
			}
		}()
		log.Println("connected to notification service")
	}

	userStorage := storage2.NewUserStorage(db)
	taskStorage := storage2.NewTaskStorage(db)
	taskCached := storage2.NewTaskCached(taskStorage, redisClient)
	productStorage := storage2.NewProductStorage(db)
	stockStorage := storage2.NewStockStorage(db)
	deliveryStorage := storage2.NewDeliveryStorage(db)
	deliveryItemsStorage := storage2.NewDeliveryItemStorage(db)
	groupsStorage := storage2.NewGroupStorage(db)
	groupMemberStorage := storage2.NewGroupMemberStorage(db)

	authService := auth.NewService(userStorage, []byte(cfg.JWTSecret))
	groupsService := groups.NewGroupService(groupsStorage)

	tasksService := tasks.NewService(
		taskCached,
		userStorage,
		groupsService,
		groupMemberStorage,
		notificationClient,
	)
	tasksService.InitMetrics()

	productService := products.NewService(productStorage)
	stocksService := stocks.NewService(stockStorage)
	groupMemberService := group_members.NewGroupMemberService(groupMemberStorage, groupsStorage, userStorage)

	deliveryService := deliveries.NewService(
		productStorage,
		userStorage,
		deliveryStorage,
		deliveryItemsStorage,
		stockStorage,
		notificationClient,
	)

	deliveryQueue := queue2.NewDeliveryQueue()
	deliveryWorker := worker2.NewDeliveryWorker(deliveryQueue, deliveryService)
	deliveryWorker.Start()

	taskQueue := queue2.NewTaskQueue()
	taskWorker := worker2.NewTaskWorker(taskQueue, tasksService)
	taskWorker.Start()

	handlers.SetAuthService(authService)
	handlers.SetTaskService(tasksService)
	handlers.SetProductService(productService)
	handlers.SetStockService(stocksService)
	handlers.SetDeliveryService(deliveryService)
	handlers.SetDeliveryQueue(deliveryQueue)
	handlers.SetGroupService(groupsService)
	handlers.SetGroupMemberService(groupMemberService)

	grpcSrv := grpcserver.NewServer(":50051", tasksService)
	go func() {
		if err := grpcSrv.Run(); err != nil {
			log.Fatalf("grpc server error: %v", err)
		}
	}()
	log.Printf("gRPC server started")

	r := router.New([]byte(cfg.JWTSecret))

	addr := ":" + cfg.HTTPPort
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
