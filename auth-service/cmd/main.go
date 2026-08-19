package main

import (
	"auth-service/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	handlers.RegisterAuthHandlers(r)

	log.Println("Listening on :8083")
	log.Fatal(http.ListenAndServe(":8083", r))
}
