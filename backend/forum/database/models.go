package database

import "time"

type Topic struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

type Message struct {
	ID         string    `json:"id"`
	Thread_id  string    `json:"thread_id"`
	User_id    string    `json:"user_id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	Tag        string    `json:"tag"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Upvotes    uint      `json:"upvotes"`
	Downvotes  uint      `json:"downvotes"`
}
