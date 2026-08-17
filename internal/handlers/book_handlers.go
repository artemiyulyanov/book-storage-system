package handlers

import (
	"book-storage-system/internal/database/repositories"
	"book-storage-system/internal/network"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookHandlers struct {
	repo *repositories.BookRepository
}

func (handlers *BookHandlers) getBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	books, err := handlers.repo.GetBooks(ctx)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	network.RespondJSON(w, http.StatusOK, books)
}

func (handlers *BookHandlers) getBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := network.ParseID(r)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, "Incorrect id!")
		return
	}

	book, err := handlers.repo.GetBook(ctx, id)

	if err != nil {
		network.RespondError(w, http.StatusNotFound, "The book is not found!")
		return
	}

	network.RespondJSON(w, http.StatusOK, book)
}

func (handlers *BookHandlers) registerRoutes(router *mux.Router) {
	router.HandleFunc("/books", handlers.getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", handlers.getBook).Methods("GET")
}

func RegisterBookHandlers(router *mux.Router, pool *pgxpool.Pool) *BookHandlers {
	repo := repositories.NewBookRepository(pool)

	handlers := BookHandlers{
		repo,
	}

	handlers.registerRoutes(router)

	return &handlers
}
