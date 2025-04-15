package models

import (
	"time"
	// "github.com/mikespook/gorbac"
)

type Session struct {
	ID     *string `json:"id,omitempy"`
	UserID *string `json:"userID,omitempy"`
}

type LoginForm struct {
	Username *string `json:"username,omitempy"`
	Email    *string `json:"email,omitempy"`
	Password *string `json:"password"`
}

type Role struct {
	ID   *string `json:"id,omitempy"`
	Name *string `json:"name,omitempy"`
	// Permissions []gorbac.Permission `json:"permissions"`
	//TODO: How do we manage permissions?
}

type User struct {
	ID          *string    `json:"id,omitempy"`
	Username    *string    `json:"username,omitempy"`
	DisplayName *string    `json:"displayName,omitempy"`
	Email       *string    `json:"email,omitempy"`
	Password    *string    `json:"password"`
	Verified    *bool      `json:"verified,omitempty"` // TODO: Verification endpoint
	RoleID      *string    `json:"roleID,omitempy"`
	CreatedAt   *time.Time `json:"createdAt,omitzero"`
	Karma       *int       `json:"karma,omitempy"`
}
