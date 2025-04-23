package models

import "time"

type Asset struct {
	ID          *int       `json:"id,omitempty"`
	Name        *string    `json:"name,omitempty" validate:"required_if=Method POST,omitnil,max=90"`
	Description *string    `json:"description,omitempty"`
	Link        *string    `json:"link,omitempty" validate:"required_if=Method POST,omitnil,url"` // Links to an S3 bucket
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	Method      string     `json:"-"`
}

// Pointers return nil if a field is omitted (e.g. in PATCH request)
