package database

import "time"

type Demo struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
	Upvotes     uint      `json:"upvotes"`
	Downvotes   uint      `json:"downvotes"`
}
