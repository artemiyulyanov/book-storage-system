package handlers

import (
	"book-service/internal/database/repository"
	"common/events"
	"common/network"
	"common/network/requests"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	kafkaPublisher "common/publisher"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookHandlers struct {
	kafkaPublisher *kafkaPublisher.KafkaAsyncPublisher
	repo           *repository.BookRepository
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
	userId, err := network.ParseUserID(r)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req *requests.BookCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		network.RespondError(w, http.StatusBadRequest, "Incorrect request body!")
		return
	}

	if err := network.Validate.Struct(req); err != nil {
		network.RespondValidationError(w, err)
		return
	}

	req.AuthorID = userId
	id, err := handlers.repo.CreateBook(ctx, req)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	go handlers.kafkaPublisher.PublishAsync(events.BookCreated, id, events.BookCreatedPayload{
		Title:       req.Title,
		Description: req.Description,
		AuthorID:    userId,
	})

	network.RespondJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (handlers *BookHandlers) updateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := network.ParseUserID(r)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := network.ParseID(r)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, "Incorrect id!")
		return
	}

	var req *requests.BookUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		network.RespondError(w, http.StatusBadRequest, "Incorrect request body!")
		return
	}

	if err := network.Validate.Struct(req); err != nil {
		network.RespondValidationError(w, err)
		return
	}

	req.AuthorID = userId
	rowsAffected, err := handlers.repo.UpdateBook(ctx, id, req)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if rowsAffected == 0 {
		network.RespondError(w, http.StatusNotFound, "Book not found!")
		return
	}

	go handlers.kafkaPublisher.PublishAsync(events.BookUpdated, id, events.BookUpdatedPayload{
		Title:       req.Title,
		Description: req.Description,
	})

	req.ID = id
	network.RespondJSON(w, http.StatusOK, req)
}

func (handlers *BookHandlers) deleteBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := network.ParseUserID(r)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := network.ParseID(r)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, "Incorrect id!")
		return
	}

	rowsAffected, err := handlers.repo.DeleteBook(ctx, id, userId)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if rowsAffected == 0 {
		network.RespondError(w, http.StatusNotFound, "Book not found!")
		return
	}

	go handlers.kafkaPublisher.PublishAsync(events.BookDeleted, id, events.BookDeletedPayload{})

	w.WriteHeader(http.StatusNoContent)
}

func (handlers *BookHandlers) registerRoutes(router *mux.Router) {
	router.HandleFunc("/", handlers.getBooks).Methods("GET")
	router.HandleFunc("/{id}", handlers.getBook).Methods("GET")

	router.HandleFunc("/", handlers.createBook).Methods("POST")

	router.HandleFunc("/{id}", handlers.updateBook).Methods("PUT")

	router.HandleFunc("/{id}", handlers.deleteBook).Methods("DELETE")
}

func RegisterBookHandlers(router *mux.Router, pool *pgxpool.Pool, kafkaProducer *events.Producer) *BookHandlers {
	repo := repository.NewBookRepository(pool)

	handlers := BookHandlers{
		kafkaPublisher: kafkaPublisher.NewKafkaAsyncPublisher(kafkaProducer),
		repo:           repo,
	}

	handlers.registerRoutes(router)

	return &handlers
}
