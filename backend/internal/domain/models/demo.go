package models

import "time"

type Demo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"` // Links to an S3 bucket
	UserID      string    `json:"userID"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Upvotes     uint      `json:"upvotes"`
	Downvotes   uint      `json:"downvotes"`
	ThreadID    string    `json:"threadID"`
}
