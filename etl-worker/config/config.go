package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	KafkaBrokers []string
	KafkaGroupID string
}

func MustLoad() Config {
	_ = godotenv.Load(".env")

	cfg := Config{
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		KafkaBrokers: getEnvSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "etl-worker"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	if len(cfg.KafkaBrokers) == 0 {
		log.Fatal("KAFKA_BROKERS not set")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvSlice(key string, fallback []string) []string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return strings.Split(val, ",")
}
