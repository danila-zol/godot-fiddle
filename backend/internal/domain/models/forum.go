package models

import "time"

type Topic struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Thread struct { // TODO: Message relation as an array of Messages, so you can count Messages in a Thread
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	UserID          string    `json:"userID"`
	TopicID         string    `json:"topicID"`
	Tag             []string  `json:"tag"`
	Created_at      time.Time `json:"created_at"`
	Last_update     time.Time `json:"last_update"`
	Total_upvotes   uint      `json:"total_upvotes"`
	Total_downvotes uint      `json:"total_downvotes"`
}

type Message struct {
	ID         string    `json:"id"`
	ThreadID   string    `json:"threadID"`
	UserID     string    `json:"userID"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	Tag        []string  `json:"tag"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Upvotes    uint      `json:"upvotes"`
	Downvotes  uint      `json:"downvotes"`
}
