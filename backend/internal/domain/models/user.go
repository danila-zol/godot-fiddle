package models

import (
	"time"
	// "github.com/mikespook/gorbac"
)

type Session struct {
	ID      *string `json:"id,omitempty"`
	UserID  *string `json:"userID,omitempty"`
}

type LoginForm struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password"`
}

type Role struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	// Permissions []gorbac.Permission `json:"permissions"`
	//TODO: How do we manage permissions?
}

type User struct {
	ID          *string    `json:"id,omitempty"`
	Username    *string    `json:"username,omitempty"`
	DisplayName *string    `json:"displayName,omitempty"`
	Email       *string    `json:"email,omitempty"`
	Password    *string    `json:"password"`
	Verified    *bool      `json:"verified,omitempty"` // TODO: Verification endpoint
	RoleID      *string    `json:"roleID,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitzero"`
	Karma       *int       `json:"karma,omitempty"`
}
