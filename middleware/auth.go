package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/Kartikeya2710/to-do/utils"
)

func AuthMiddleware(jwtManager *utils.JWTManager, logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Printf("No authorization header provided in request")
				http.Error(w, "Missing authorization heaader", http.StatusBadRequest)

				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				logger.Println("Invalid Authorization header format")
				http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)

				return
			}

			jwtString := parts[1]
			claims, err := jwtManager.ValidateJWT(jwtString)
			if err != nil {
				logger.Println("Invalid JWT provided")
				http.Error(w, "invalid jwt provided", http.StatusUnauthorized)

				return
			}

			logger.Println("JWT validated successfully")

			ctx := WithUserID(r.Context(), claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
