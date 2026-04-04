package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Xanaduxan/tasks-golang/etl-worker/config"
	"github.com/Xanaduxan/tasks-golang/etl-worker/internal/kafka"
	"github.com/Xanaduxan/tasks-golang/etl-worker/internal/service"
	"github.com/Xanaduxan/tasks-golang/etl-worker/internal/storage"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustLoad()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	eventStorage := storage.NewTaskEventLogStorage(db)
	analyticsStorage := storage.NewTaskUserAnalyticsStorage(db)

	processor := service.NewProcessor(eventStorage, analyticsStorage)

	consumer, err := kafka.NewConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaGroupID,
		processor,
	)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("failed to close consumer: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("shutting down...")
		cancel()
	}()

	log.Println("etl-worker started")

	if err := consumer.Run(ctx); err != nil {
		log.Fatalf("consumer error: %v", err)
	}
}
