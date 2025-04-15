package models

import "time"

type Demo struct {
	ID          *int       `json:"id,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Link        *string    `json:"link,omitempty"` // Links to an S3 bucket
	Tags        *[]string  `json:"tags,omitempty"`
	UserID      *string    `json:"userID,omitempty"`
	ThreadID    *int       `json:"threadID,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitzero"`
	UpdatedAt   *time.Time `json:"updatedAt,omitzero"`
	Upvotes     *uint      `json:"upvotes,omitzero"`
	Downvotes   *uint      `json:"downvotes,omitzero"`
}
