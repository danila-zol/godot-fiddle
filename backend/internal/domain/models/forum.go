package models

import "time"

type Topic struct {
	ID   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type Thread struct { // TODO: Message relation as an array of Messages, so you can count Messages in a Thread
	ID        *int       `json:"id,omitempty"`
	Title     *string    `json:"title,omitempty"`
	UserID    *string    `json:"userID,omitempty"`
	TopicID   *int       `json:"topicID,omitempty"`
	Tags      *[]string  `json:"tags,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitzero"`
	UpdatedAt *time.Time `json:"updatedAt,omitzero"`
	Upvotes   *uint      `json:"upvotes,omitempty"`
	Downvotes *uint      `json:"downvotes,omitempty"`
}

type Message struct {
	ID        *int       `json:"id,omitempty"`
	ThreadID  *int       `json:"threadID,omitempty"`
	UserID    *string    `json:"userID,omitempty"`
	Title     *string    `json:"title,omitempty"`
	Body      *string    `json:"body,omitempty"`
	Tags      *[]string  `json:"tags,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitzero"`
	UpdatedAt *time.Time `json:"updatedAt,omitzero"`
	Upvotes   *uint      `json:"upvotes,omitempty"`
	Downvotes *uint      `json:"downvotes,omitempty"`
}
