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
	Username    *string    `form:"username" json:"username,omitempty" validate:"required_if=Method POST,omitnil,max=90"`
	DisplayName *string    `form:"displayName" json:"displayName,omitempty" validate:"omitnil,max=200"`
	Email       *string    `form:"email" json:"email,omitempty" validate:"required_if=Method POST,omitnil,email,max=50"`
	Password    *string    `json:"-"`
	Verified    *bool      `form:"verified" json:"verified,omitempty"`
	Role        *string    `form:"role" json:"role,omitempty" validate:"required_if=Method POST,omitnil,max=255"`
	CreatedAt   *time.Time `json:"createdAt,omitzero"`
	Karma       *int       `json:"karma,omitempty" validate:"omitnil,number"`
	ProfilePic  string     `json:"profilePic"`
	Method      string     `json:"-"`
}
