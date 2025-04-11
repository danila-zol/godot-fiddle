package models

import "time"

type Topic struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type Thread struct { // TODO: Message relation as an array of Messages, so you can count Messages in a Thread
	ID             *string    `json:"id,omitempty"`
	Title          *string    `json:"title,omitempty"`
	UserID         *string    `json:"userID,omitempty"`
	TopicID        *string    `json:"topicID,omitempty"`
	Tags           *[]string  `json:"tags,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitzero"`
	LastUpdate     *time.Time `json:"lastUpdate,omitzero"`
	TotalUpvotes   *uint      `json:"totalUpvotes,omitempty"`
	TotalDownvotes *uint      `json:"totalDownvotes,omitempty"`
}

type Message struct {
	ID        *string    `json:"id,omitempty"`
	ThreadID  *string    `json:"threadID,omitempty"`
	UserID    *string    `json:"userID,omitempty"`
	Title     *string    `json:"title,omitempty"`
	Body      *string    `json:"body,omitempty"`
	Tags      *[]string  `json:"tags,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitzero"`
	UpdatedAt *time.Time `json:"updatedAt,omitzero"`
	Upvotes   *uint      `json:"upvotes,omitempty"`
	Downvotes *uint      `json:"downvotes,omitempty"`
}
