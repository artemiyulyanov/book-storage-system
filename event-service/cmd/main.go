package main

import (
	"common/database"
	"context"
	"event-service/internal/clickhouse"
	"event-service/internal/consumer"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	clickhouseDSN := os.Getenv("CLICKHOUSE_DSN")

	if err := database.RunClickHouseMigrations(clickhouseDSN, "migrations/clickhouse"); err != nil {
		log.Fatalf("clickhouse migrations failed: %v", err)
	}

	store, err := clickhouse.NewStore(clickhouseDSN)
	if err != nil {
		log.Fatalf("failed to connect to clickhouse: %v", err)
	}
	defer store.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("event-service started, consuming events...")
	consumer.Run(ctx, brokers, "event-service-group", store)
}
