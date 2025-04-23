package models

import (
	"time"
	// "github.com/mikespook/gorbac"
)

type Session struct {
	ID     *string `json:"id,omitempty"`
	UserID *string `json:"userID,omitempty"`
}

type Role struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty" validate:"required_if=Method POST,max=90"`
	// Permissions []gorbac.Permission `json:"permissions"`
	//TODO: How do we manage permissions?
	Method string `json:"-"`
}

type User struct {
	ID          *string    `json:"id,omitempty"`
	Username    *string    `json:"username,omitempty" validate:"required_if=Method POST,omitnil,max=90"`
	DisplayName *string    `json:"displayName,omitempty" validate:"omitnil,max=200"`
	Email       *string    `json:"email,omitempty" validate:"required_if=Method POST,omitnil,email,max=50"`
	Password    *string    `json:"-"`
	Verified    *bool      `json:"verified,omitempty"`
	RoleID      *string    `json:"roleID,omitempty" validate:"required_if=Method POST,omitnil,uuid4"`
	CreatedAt   *time.Time `json:"createdAt,omitzero"`
	Karma       *int       `json:"karma,omitempty" validate:"omitnil,number"`
	Method      string     `json:"-"`
}
