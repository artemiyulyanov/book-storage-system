package handlers

import (
	"auth-service/internal/grpcclients"
	"common/auth"
	"common/events"
	"common/network"
	"common/network/requests"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gorilla/mux"
)

type AuthHandlers struct {
	userClient    *grpcclients.UserClient
	kafkaProducer *events.Producer
}

func (handlers *AuthHandlers) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req *requests.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		network.RespondError(w, http.StatusBadRequest, "Incorrect request body!")
		return
	}

	if err := network.Validate.Struct(req); err != nil {
		network.RespondValidationError(w, err)
		return
	}

	user, err := handlers.userClient.GetUserByEmail(ctx, req.Email)

	if err != nil {
		st, ok := status.FromError(err)

		if ok && st.Code() == codes.NotFound {
			network.RespondError(w, http.StatusUnauthorized, "Incorrect email or password")
			return
		}

		network.RespondError(w, http.StatusBadGateway, "Failed to reach user-service")
		return
	}

	if !auth.PasswordsHashEqual(user.PasswordHash, req.Password) {
		network.RespondError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	token, err := auth.GenerateToken(user.Id, os.Getenv("JWT_SECRET"), 24*time.Hour)
	if err != nil {
		network.RespondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	clientIP := network.ParseClientIP(r)

	go handlers.kafkaProducer.Publish(ctx, events.UserLoggedIn, user.Id, events.UserLoggedInPayload{
		IP: clientIP,
	})

	network.RespondJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (handlers *AuthHandlers) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req *requests.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		network.RespondError(w, http.StatusBadRequest, "Incorrect request body!")
		return
	}

	if err := network.Validate.Struct(req); err != nil {
		network.RespondValidationError(w, err)
		return
	}

	res, err := handlers.userClient.CreateUser(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)

		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				network.RespondError(w, http.StatusBadRequest, "Passwords do not match!")
				return
			case codes.AlreadyExists:
				network.RespondError(w, http.StatusConflict, "Email already registered!")
				return
			}
		}

		network.RespondError(w, http.StatusBadGateway, "Failed to reach user-service")
		return
	}

	token, err := auth.GenerateToken(res.Id, os.Getenv("JWT_SECRET"), 24*time.Hour)
	if err != nil {
		network.RespondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	clientIP := network.ParseClientIP(r)

	go handlers.kafkaProducer.Publish(ctx, events.UserRegistered, res.Id, events.UserRegisteredPayload{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		IP:        clientIP,
	})

	network.RespondJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func (handlers *AuthHandlers) registerRoutes(router *mux.Router) {
	router.HandleFunc("/login", handlers.login).Methods("POST")
	router.HandleFunc("/register", handlers.register).Methods("POST")
}

func RegisterAuthHandlers(router *mux.Router, userClient *grpcclients.UserClient, kafkaProducer *events.Producer) *AuthHandlers {
	handlers := AuthHandlers{
		userClient:    userClient,
		kafkaProducer: kafkaProducer,
	}

	handlers.registerRoutes(router)

	return &handlers
}
