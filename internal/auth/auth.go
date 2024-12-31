package auth

import (
    "encoding/json"
    "net/http"
    "strconv"
    "log"
    "github.com/blobfish465/common-circle-web-forum/internal/database"
    "github.com/blobfish465/common-circle-web-forum/internal/dataaccess/users"
    "github.com/blobfish465/common-circle-web-forum/internal/models"
    "github.com/blobfish465/common-circle-web-forum/internal/utils"
    "golang.org/x/crypto/bcrypt"
)

// Handles the neccessary authentication for user login
func Login(w http.ResponseWriter, r *http.Request) {
    var credentials models.User
    err := json.NewDecoder(r.Body).Decode(&credentials)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    db, err := database.GetDB()
    if err != nil {
        http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
        return
    }

    user, err := users.GetUserByUsername(db, credentials.Username)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Check password here
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    log.Printf("User ID to encode in JWT: %d\n", user.ID)
    token, err := utils.GenerateJWT(strconv.Itoa(user.ID))  // Convert int ID to string
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}
