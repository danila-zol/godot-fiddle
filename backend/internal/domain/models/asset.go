package models

import "time"

type Asset struct {
	ID          *int       `json:"id,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Link        *string    `json:"link,omitempty"` // Links to an S3 bucket
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
}

// Pointers return nil if a field is omitted (e.g. in PATCH request)
