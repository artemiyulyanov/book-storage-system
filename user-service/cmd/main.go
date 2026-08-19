package main

import (
	"common/database"
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"user-service/internal/database/repository"
	"user-service/internal/grpcserver"
	"user-service/internal/handlers"

	pb "common/proto/user"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
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

	repo := repository.NewUserRepository(pool)

	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatalf("failed to listen on :9090: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterUserServiceServer(grpcServer, grpcserver.NewUserServer(repo))

		log.Println("gRPC listening on :9090")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("grpc server failed: %v", err)
		}
	}()

	r := mux.NewRouter()
	handlers.RegisterUserHandlers(r, pool, repo)

	log.Println("Listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}
