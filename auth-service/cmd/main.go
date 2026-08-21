package main

import (
	"auth-service/internal/grpcclients"
	"auth-service/internal/handlers"
	"common/events"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	userClient, err := grpcclients.NewUserClient(os.Getenv("USER_SERVICE_GRPC_URL"))

	if err != nil {
		log.Fatalf("failed to connect to user-service: %v", err)
	}

	defer userClient.Close()

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	kafkaProducer := events.NewProducer(brokers)

	r := mux.NewRouter()
	handlers.RegisterAuthHandlers(r, userClient, kafkaProducer)

	log.Println("Listening on :8083")
	log.Fatal(http.ListenAndServe(":8083", r))
}
