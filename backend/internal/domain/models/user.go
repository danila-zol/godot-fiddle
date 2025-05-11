package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID     *uuid.UUID `json:"id,omitempty"`
	UserID *uuid.UUID `json:"userID,omitempty"`
}

type User struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	Username    *string    `json:"username,omitempty" validate:"required_if=Method POST,omitnil,max=90"`
	DisplayName *string    `json:"displayName,omitempty" validate:"omitnil,max=200"`
	Email       *string    `json:"email,omitempty" validate:"required_if=Method POST,omitnil,email,max=50"`
	Password    *string    `json:"-"`
	Verified    *bool      `json:"verified,omitempty"`
	Role        *string    `json:"role,omitempty" validate:"required_if=Method POST,omitnil,max=255"`
	CreatedAt   *time.Time `json:"createdAt,omitzero"`
	Karma       *int       `json:"karma,omitempty" validate:"omitnil,number"`
	Method      string     `json:"-"`
}
