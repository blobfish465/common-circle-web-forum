package models

type Comment struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
	UserID int `json:"user_id"`
	ThreadID  int    `json:"thread_id"`
}