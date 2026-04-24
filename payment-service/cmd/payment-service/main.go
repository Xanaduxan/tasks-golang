package main

import (
	"log"
	"net/http"

	"github.com/Xanaduxan/tasks-golang/payment-service/config"
	"github.com/Xanaduxan/tasks-golang/payment-service/internal/service/payments"
	"github.com/Xanaduxan/tasks-golang/payment-service/internal/service/shops"
	"github.com/Xanaduxan/tasks-golang/payment-service/internal/storage"
	http_handlers "github.com/Xanaduxan/tasks-golang/payment-service/internal/transport/http-handlers"
	"github.com/Xanaduxan/tasks-golang/payment-service/internal/transport/router"
	"github.com/Xanaduxan/tasks-golang/payment-service/worker"
)

func main() {
	cfg := config.MustLoad()

	db := storage.NewPostgres(cfg.DatabaseURL)

	shopStorage := storage.NewShopStorage(db)
	paymentStorage := storage.NewPaymentStorage(db)
	shopService := shops.NewService(shopStorage)
	paymentService := payments.NewPaymentsService(paymentStorage, shopStorage)

	paymentWorker := worker.NewPaymentWorker(paymentService)
	paymentWorker.Start()
	http_handlers.SetShopService(shopService)
	http_handlers.SetPaymentsService(paymentService)
	r := router.New([]byte(cfg.JWTSecret))

	addr := ":" + cfg.HTTPPort
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
