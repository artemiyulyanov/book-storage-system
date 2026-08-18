package main

import (
	"book-service/internal/database"
	"book-service/internal/handlers"
	"context"
	"log"
	"net/http"
	"os"

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

	handlers.RegisterBookHandlers(r, pool)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
