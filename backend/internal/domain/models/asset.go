package models

import "time"

type Asset struct {
	ID           *int       `form:"id" json:"id,omitempty"`
	Name         *string    `form:"name" json:"name,omitempty" validate:"required_if=Method POST,omitnil,max=90"`
	Description  *string    `form:"description" json:"description,omitempty"`
	Tags         *[]string  `form:"tags" json:"tags,omitempty" validate:"omitnil,unique,max=40"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitzero"`
	Version      *int       `form:"version" json:"version,omitempty" validate:"required_if=Method PATCH,omitnil,number,gt=0"`
	Upvotes      *uint      `form:"upvotes" json:"upvotes,omitzero" validate:"omitnil,number,min=0"`
	Downvotes    *uint      `form:"downvotes" json:"downvotes,omitzero" validate:"omitnil,number,min=0"`
	Rating       *float64   `json:"rating,omitzero" validate:"omitnil,excluded_if=Method GET"`
	Views        *uint      `json:"views,omitzero" validate:"omitnil,number,min=0"`
	Key          *string    `json:"key"` // Links to an S3 bucket
	ThumbnailKey *string    `json:"thumbnailKey"`
	Method       string     `json:"-"`
}

// Pointers return nil if a field is omitted (e.g. in PATCH request)
