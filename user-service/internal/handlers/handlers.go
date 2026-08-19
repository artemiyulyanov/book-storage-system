package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"user-service/internal/database/repository"

	"common/network"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandlers struct {
	repo *repository.UserRepository
}

func (handlers *UserHandlers) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := handlers.repo.GetUsers(ctx)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	network.RespondJSON(w, http.StatusOK, users)
}

func (handlers *UserHandlers) getUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := network.ParseID(r)

	if err != nil {
		network.RespondError(w, http.StatusBadRequest, "Incorrect id!")
		return
	}

	user, err := handlers.repo.GetUser(ctx, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			network.RespondError(w, http.StatusNotFound, "User not found!")
		} else {
			network.RespondError(w, http.StatusBadRequest, err.Error())
		}

		return
	}

	network.RespondJSON(w, http.StatusOK, user)
}

func (handlers *UserHandlers) registerRoutes(router *mux.Router) {
	router.HandleFunc("/", handlers.getUsers).Methods("GET")
	router.HandleFunc("/{id}", handlers.getUser).Methods("GET")
}

func RegisterUserHandlers(router *mux.Router, pool *pgxpool.Pool) *UserHandlers {
	repo := repository.NewUserRepository(pool)

	handlers := UserHandlers{
		repo,
	}

	handlers.registerRoutes(router)

	return &handlers
}
