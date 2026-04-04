package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	HTTPPort    string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	NotificationAddr string
	KafkaBrokers     string
}

func MustLoad() Config {
	_ = LoadEnv(".env")

	cfg := Config{
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		HTTPPort:         getEnv("HTTP_PORT", "8080"),
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          getEnvInt("REDIS_DB", 0),
		NotificationAddr: getEnv("NOTIFICATION_ADDR", "localhost:50052"),
		KafkaBrokers:     getEnv("KAFKA_BROKERS", "localhost:9092"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	return cfg
}

func LoadEnv(path string) error {
	return godotenv.Load(path)
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("%s must be int, got: %s", key, val)
	}

	return num
}
