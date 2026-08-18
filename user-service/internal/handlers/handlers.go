package handlers

import (
	"user-service/internal/database/repository"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandlers struct {
	repo *repository.UserRepository
}

func (handlers *UserHandlers) registerRoutes(router *mux.Router) {
	// router.HandleFunc("/books", handlers.getBooks).Methods("GET")
	// router.HandleFunc("/books/{id}", handlers.getBook).Methods("GET")

	// router.HandleFunc("/books", handlers.createBook).Methods("POST")

	// router.HandleFunc("/books/{id}", handlers.updateBook).Methods("PUT")

	// router.HandleFunc("/books/{id}", handlers.deleteBook).Methods("DELETE")
}

func RegisterUserHandlers(router *mux.Router, pool *pgxpool.Pool) *UserHandlers {
	repo := repository.NewUserRepository(pool)

	handlers := UserHandlers{
		repo,
	}

	handlers.registerRoutes(router)

	return &handlers
}
