package models

import (
	"time"

	"github.com/google/uuid"
)

type Topic struct {
	ID      *int    `json:"id,omitempty"`
	Name    *string `json:"name,omitempty" validate:"required,lt=90"`
	Version *int    `json:"version,omitempty" validate:"required_if=Method PATCH,omitnil,number,gt=0"`
	Method  string  `json:"-"`
}

type Thread struct {
	ID        *int       `json:"id,omitempty"`
	Title     *string    `json:"title,omitempty" validate:"required_if=Method POST,omitnil,lt=90"`
	UserID    *uuid.UUID `json:"userID,omitempty" validate:"required_if=Method POST,omitnil,uuid4"`
	TopicID   *int       `json:"topicID,omitempty" validate:"required_if=Method POST,omitnil,number"`
	Tags      *[]string  `json:"tags,omitempty" validate:"omitnil,unique,max=40"`
	CreatedAt *time.Time `json:"createdAt,omitzero"`
	UpdatedAt *time.Time `json:"updatedAt,omitzero"`
	Upvotes   *uint      `json:"upvotes,omitempty" validate:"omitnil,number,min=0"`
	Downvotes *uint      `json:"downvotes,omitempty" validate:"omitnil,number,min=0"`
	Method    string     `json:"-"`
}

type Message struct {
	ID        *int       `json:"id,omitempty"`
	ThreadID  *int       `json:"threadID,omitempty" validate:"required_if=Method POST,omitnil,number"`
	UserID    *uuid.UUID `json:"userID,omitempty" validate:"required_if=Method POST,omitnil,uuid4"`
	Title     *string    `json:"title,omitempty" validate:"required_if=Method POST,omitnil,lt=90"`
	Body      *string    `json:"body,omitempty" validate:"omitnil,max=10000"`
	Tags      *[]string  `json:"tags,omitempty" validate:"omitnil,unique,max=40"`
	CreatedAt *time.Time `json:"createdAt,omitzero"`
	UpdatedAt *time.Time `json:"updatedAt,omitzero"`
	Upvotes   *uint      `json:"upvotes,omitempty" validate:"omitnil,number,min=0"`
	Downvotes *uint      `json:"downvotes,omitempty" validate:"omitnil,number,min=0"`
	Method    string     `json:"-"`
}
