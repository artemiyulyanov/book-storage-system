package main

import (
	"book-service/internal/handlers"
	"common/database"
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

	log.Println("Listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
