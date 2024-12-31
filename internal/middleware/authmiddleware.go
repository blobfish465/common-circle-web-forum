package middleware

import (
	"context"
	"net/http"
	"strings"
	"log" 

	"github.com/golang-jwt/jwt/v4"
	"github.com/blobfish465/common-circle-web-forum/internal/utils"
)

var jwtSecret = utils.GetJWTSecret()

// Claims defines the structure of the JWT payload
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates the JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the "Authorization" header
		authHeader := r.Header.Get("Authorization")
		log.Println("Authorization Header:", authHeader) 
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Ensure the token follows the "Bearer <token>" format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenStr := tokenParts[1]

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		log.Printf("Authenticated User ID: %s\n", claims.UserID)

		// Attach the user ID to the request context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
