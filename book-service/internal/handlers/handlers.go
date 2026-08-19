package handlers

import (
	"book-service/internal/database/repository"
	"common/models"
	"common/network"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookHandlers struct {
	repo *repository.BookRepository
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
		if errors.Is(err, sql.ErrNoRows) {
			network.RespondError(w, http.StatusNotFound, "Book not found!")
		} else {
			network.RespondError(w, http.StatusBadRequest, err.Error())
		}

		return
	}

	network.RespondJSON(w, http.StatusOK, book)
}

func (handlers *BookHandlers) createBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var book *models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		network.RespondError(w, http.StatusBadRequest, "Incorrect request body!")
		return
	}

	if err := network.Validate.Struct(book); err != nil {
		network.RespondValidationError(w, err)
		return
	}

	id, err := handlers.repo.CreateBook(ctx, book)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	network.RespondJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (handlers *BookHandlers) updateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := network.ParseID(r)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, "Incorrect id!")
		return
	}

	var book *models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		network.RespondError(w, http.StatusBadRequest, "Incorrect request body!")
		return
	}

	if err := network.Validate.Struct(book); err != nil {
		network.RespondValidationError(w, err)
		return
	}

	rowsAffected, err := handlers.repo.UpdateBook(ctx, id, book)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if rowsAffected == 0 {
		network.RespondError(w, http.StatusNotFound, "Book not found!")
		return
	}

	book.ID = id
	network.RespondJSON(w, http.StatusOK, book)
}

func (handlers *BookHandlers) deleteBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := network.ParseID(r)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, "Incorrect id!")
		return
	}

	rowsAffected, err := handlers.repo.DeleteBook(ctx, id)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if rowsAffected == 0 {
		network.RespondError(w, http.StatusNotFound, "Book not found!")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handlers *BookHandlers) registerRoutes(router *mux.Router) {
	router.HandleFunc("/", handlers.getBooks).Methods("GET")
	router.HandleFunc("/{id}", handlers.getBook).Methods("GET")

	router.HandleFunc("/", handlers.createBook).Methods("POST")

	router.HandleFunc("/{id}", handlers.updateBook).Methods("PUT")

	router.HandleFunc("/{id}", handlers.deleteBook).Methods("DELETE")
}

func RegisterBookHandlers(router *mux.Router, pool *pgxpool.Pool) *BookHandlers {
	repo := repository.NewBookRepository(pool)

	handlers := BookHandlers{
		repo,
	}

	handlers.registerRoutes(router)

	return &handlers
}
