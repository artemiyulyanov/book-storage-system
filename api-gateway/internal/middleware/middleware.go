package middleware

import (
	"common/auth"
	"common/network"
	"net/http"
	"strconv"
	"strings"
)

func JWTAuthForMethods(secret string, protectedMethods map[string]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !protectedMethods[r.Method] {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authHeader, "Bearer ") {
				network.RespondError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := auth.ParseToken(tokenString, secret)
			if err != nil {
				network.RespondError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			r.Header.Set("X-User-Id", strconv.FormatInt(claims.UserID, 10))

			next.ServeHTTP(w, r)
		})
	}
}
