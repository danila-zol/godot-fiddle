package models

import (
	"time"

	"github.com/google/uuid"
)

type Demo struct {
	ID          *int       `json:"id,omitempty"`
	Title       *string    `json:"title,omitempty" validate:"required_if=Method POST,omitnil,lt=90"`
	Description *string    `json:"description,omitempty" validate:"omitnil,max=5000"`
	Link        *string    `json:"link,omitempty" validate:"required_if=Method POST,omitnil,url"` // Links to an S3 bucket
	Tags        *[]string  `json:"tags,omitempty" validate:"omitnil,unique,max=40"`
	UserID      *uuid.UUID `json:"userID,omitempty" validate:"required_if=Method POST,omitnil,uuid4"`
	ThreadID    *int       `json:"threadID,omitempty" validate:"required_if=Method PATCH,omitnil,number"`
	CreatedAt   *time.Time `json:"createdAt,omitzero"`
	UpdatedAt   *time.Time `json:"updatedAt,omitzero"`
	Upvotes     *uint      `json:"upvotes,omitzero" validate:"omitnil,number,min=0"`
	Downvotes   *uint      `json:"downvotes,omitzero" validate:"omitnil,number,min=0"`
	Views       *uint      `json:"views,omitzero" validate:"omitnil,number,min=0"`
	Method      string     `json:"-"`
}
