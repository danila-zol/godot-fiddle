package models

import "time"

type Demo struct {
	ID          *int       `json:"id,omitempty"`
	Title       *string    `json:"title,omitempty" validate:"required,lt=90"`
	Description *string    `json:"description,omitempty" validate:"omitnil,max=5000"`
	Link        *string    `json:"link,omitempty" validate:"required,url"` // Links to an S3 bucket
	Tags        *[]string  `json:"tags,omitempty" validate:"omitnil,unique,max=40"`
	UserID      *string    `json:"userID,omitempty" validate:"required"`
	ThreadID    *int       `json:"threadID,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitzero"`
	UpdatedAt   *time.Time `json:"updatedAt,omitzero"`
	Upvotes     *uint      `json:"upvotes,omitzero" validate:"omitnil,number,min=0"`
	Downvotes   *uint      `json:"downvotes,omitzero" validate:"omitnil,number,min=0"`
}
