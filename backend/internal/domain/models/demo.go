package models

import (
	"time"

	"github.com/google/uuid"
)

type Demo struct {
	ID           *int       `form:"id" json:"id,omitempty"`
	Title        *string    `form:"title" json:"title,omitempty" validate:"required_if=Method POST,omitnil,lt=90"`
	Description  *string    `form:"description" json:"description,omitempty" validate:"omitnil,max=5000"`
	Tags         *[]string  `form:"tags" json:"tags,omitempty" validate:"omitnil,unique,max=40"`
	UserID       *uuid.UUID `form:"userID" json:"userID,omitempty" validate:"required_if=Method POST,omitnil,uuid4"`
	ThreadID     *int       `form:"threadID" json:"threadID,omitempty" validate:"required_if=Method PATCH,omitnil,number"`
	CreatedAt    *time.Time `json:"createdAt,omitzero"`
	UpdatedAt    *time.Time `json:"updatedAt,omitzero"`
	Upvotes      *uint      `form:"upvotes" json:"upvotes,omitzero" validate:"omitnil,number,min=0"`
	Downvotes    *uint      `form:"downvotes" json:"downvotes,omitzero" validate:"omitnil,number,min=0"`
	Rating       *float64   `json:"rating,omitzero" validate:"omitnil,excluded_if=Method GET"`
	Views        *uint      `json:"views,omitzero" validate:"omitnil,number,min=0"`
	Key          *string    `json:"key"` // Links to an S3 bucket
	ThumbnailKey *string    `json:"thumbnailKey"`
	Method       string     `json:"-"`
}
