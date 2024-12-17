package database

import "time"

type Demo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	User_id     string    `json:"user_id"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
	Upvotes     uint      `json:"upvotes"`
	Downvotes   uint      `json:"downvotes"`
	Topic_id    string    `json:"topic_id"`
}

type Asset struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Created_at  time.Time `json:"created_at"`
}
