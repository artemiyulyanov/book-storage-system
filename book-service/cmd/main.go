package main

import (
	"book-service/internal/handlers"
	"common/database"
	"common/events"
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	connectionString := os.Getenv("DATABASE_URL")

	pool, err := database.RegisterPool(context.Background(), connectionString)

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	if err := database.RunMigrations(connectionString, "migrations"); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	kafkaProducer := events.NewProducer(brokers)
	defer kafkaProducer.Close()

	handlers.RegisterBookHandlers(r, pool, kafkaProducer)

	log.Println("Listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
