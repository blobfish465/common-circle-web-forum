package utils

import (
	"log"
	"os"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

// jwtSecret holds the secret used to sign JWTs. It is initialized in the init function.
var jwtSecret []byte

// Claims defines the structure of JWT claims.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// init loads the .env file and initializes jwtSecret with the JWT_SECRET environment variable.
func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file, falling back to system environment variables")
	}

	// Get JWT_SECRET from environment variables
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set in the environment variables")
	}
	jwtSecret = []byte(secret)
}

// GetJWTSecret returns the JWT secret as a byte slice
func GetJWTSecret() []byte {
	return jwtSecret
}

// GenerateJWT generates a JWT token for a given user ID.
func GenerateJWT(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Token expiration (e.g., 24 hours)
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
