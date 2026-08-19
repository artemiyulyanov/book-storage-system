package main

import (
	"api-gateway/internal/proxy"
	"log"
	"net/http"
	"os"
)

func main() {
	router := proxy.NewRouter()

	bookServiceURL := os.Getenv("BOOK_SERVICE_URL")
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")

	if err := router.RegisterService("/api/books", bookServiceURL); err != nil {
		log.Fatalf("Failed to register service book-service: %v", err)
	}

	if err := router.RegisterService("/api/users/", userServiceURL); err != nil {
		log.Fatalf("Failed to register service user-service: %v", err)
	}

	if err := router.RegisterService("/api/auth/", authServiceURL); err != nil {
		log.Fatalf("Failed to register service auth-service: %v", err)
	}

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router.Handler()))
}
