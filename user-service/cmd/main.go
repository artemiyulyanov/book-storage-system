package main

import (
	"common/database"
	"context"
	"log"
	"net/http"
	"os"

	"user-service/internal/handlers"

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

	handlers.RegisterUserHandlers(r, pool)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
