package handlers

import (
	"common/network"
	"net/http"

	"github.com/gorilla/mux"
)

type AuthHandlers struct {
}

func (handlers *AuthHandlers) ping(w http.ResponseWriter, r *http.Request) {
	network.RespondJSON(w, http.StatusOK, map[string]string{"hello": "world"})
}

func (handlers *AuthHandlers) registerRoutes(router *mux.Router) {
	router.HandleFunc("/ping", handlers.ping).Methods("GET")
}

func RegisterAuthHandlers(router *mux.Router) *AuthHandlers {
	handlers := AuthHandlers{}

	handlers.registerRoutes(router)

	return &handlers
}
