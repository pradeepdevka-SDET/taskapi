package model

import "time"

type Task struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}
