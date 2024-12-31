package models

import "fmt"

type User struct {
	ID   int    `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"-"` // Not serialized when sent over JSON
	Password      string `json:"password,omitempty"` // Used only for binding incoming JSON data
}

func (user *User) Greet() string {
	return fmt.Sprintf("Hello, I am %s", user.Username)
}
