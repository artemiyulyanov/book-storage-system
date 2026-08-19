package main

import (
	"auth-service/internal/grpcclients"
	"auth-service/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	userClient, err := grpcclients.NewUserClient(os.Getenv("USER_SERVICE_GRPC_URL"))

	if err != nil {
		log.Fatalf("failed to connect to user-service: %v", err)
	}

	defer userClient.Close()

	r := mux.NewRouter()
	handlers.RegisterAuthHandlers(r, userClient)

	log.Println("Listening on :8083")
	log.Fatal(http.ListenAndServe(":8083", r))
}
