package models

import "time"

type Demo struct {
	ID          *int       `json:"id,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Tags        *[]string  `json:"tags,omitempty"`
	Link        *string    `json:"link,omitempty"` // Links to an S3 bucket
	UserID      *string    `json:"userID,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitzero"`
	UpdatedAt   *time.Time `json:"updatedAt,omitzero"`
	Upvotes     *uint      `json:"upvotes,omitzero"`
	Downvotes   *uint      `json:"downvotes,omitzero"`
	ThreadID    *int       `json:"threadID,omitempty"`
}
