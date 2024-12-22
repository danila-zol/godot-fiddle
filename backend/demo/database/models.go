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
	Thread_id   string    `json:"thread_id"`
}

type Thread struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	User_id         string    `json:"user_id"`
	Topic_id        string    `json:"topic_id"`
	Tag             string    `json:"tag"`
	Created_at      time.Time `json:"created_at"`
	Last_update     time.Time `json:"last_update"`
	Total_upvotes   uint      `json:"total_upvotes"`
	Total_downvotes uint      `json:"total_downvotes"`
}
