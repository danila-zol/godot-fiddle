package models

import "time"

type Topic struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Thread struct { // TODO: Message relation as an array of Messages, so you can count Messages in a Thread
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	UserID         string    `json:"userID"`
	TopicID        string    `json:"topicID"`
	Tags           []string  `json:"tags"`
	CreatedAt      time.Time `json:"createdAt"`
	LastUpdate     time.Time `json:"lastUpdate"`
	TotalUpvotes   uint      `json:"totalUpvotes"`
	TotalDownvotes uint      `json:"totalDownvotes"`
}

type Message struct {
	ID        string    `json:"id"`
	ThreadID  string    `json:"threadID"`
	UserID    string    `json:"userID"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Upvotes   uint      `json:"upvotes"`
	Downvotes uint      `json:"downvotes"`
}
