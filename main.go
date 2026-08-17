package main

import (
	"book-storage-system/internal/database"
	"book-storage-system/internal/handlers"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	connectionString := os.Getenv("DATABASE_URL")

	pool, err := database.NewDatabaseInstance(context.Background(), connectionString)

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	r := mux.NewRouter()

	handlers.RegisterBookHandlers(r, pool)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
