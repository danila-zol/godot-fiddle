package models

import "time"

type Asset struct {
	ID          *int       `json:"id,omitempty"`
	Name        *string    `json:"name,omitempty" validate:"required,lt=90"`
	Description *string    `json:"description,omitempty"`
	Link        *string    `json:"link,omitempty" validate:"required,url"` // Links to an S3 bucket
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
}

// Pointers return nil if a field is omitted (e.g. in PATCH request)
