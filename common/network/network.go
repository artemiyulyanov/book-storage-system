package network

import (
	"common/models"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var Validate = validator.New()
var ErrMissingUserID = errors.New("missing X-User-Id header")
var ErrMissingUserRole = errors.New("missing X-User-Role header")

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{"error": message})
}

func RespondValidationError(w http.ResponseWriter, err error) {
	fieldErrors := make(map[string]string)

	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		for _, fe := range ve {
			fieldErrors[fe.Field()] = validationMessage(fe)
		}
	} else {
		fieldErrors["_"] = err.Error()
	}

	RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"errors": fieldErrors})
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "required"
	case "email":
		return "incorrect email"
	case "min":
		return "too short value (min " + fe.Param() + ")"
	case "max":
		return "too long value (max " + fe.Param() + ")"
	default:
		return "incorrect value"
	}
}

func ParseID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	return strconv.ParseInt(vars["id"], 10, 64)
}

func ParseUserID(r *http.Request) (int64, error) {
	userIDHeader := r.Header.Get("X-User-Id")

	if userIDHeader == "" {
		return 0, ErrMissingUserID
	}

	return strconv.ParseInt(userIDHeader, 10, 64)
}

func ParseUserRole(r *http.Request) (models.UserRole, error) {
	userRoleHeader := r.Header.Get("X-User-Role")

	if userRoleHeader == "" {
		return models.ROLE_USER, ErrMissingUserRole
	}

	return models.UserRole(userRoleHeader), nil
}

func ParseClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
