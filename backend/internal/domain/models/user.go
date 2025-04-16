package models

import (
	"time"
	// "github.com/mikespook/gorbac"
)

type Session struct {
	ID     *string `json:"id,omitempty"`
	UserID *string `json:"userID,omitempty"`
}

type LoginForm struct {
	Username *string `json:"username,omitempty" validate:"omitempty,required_without=Email,max=90"`
	Email    *string `json:"email,omitempty" validate:"omitempty,required_without=Username,email,max=50"`
	Password *string `json:"password" validate:"required,min=8,alphanum"`
}

type Role struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty" validate:"required,max=90"`
	// Permissions []gorbac.Permission `json:"permissions"`
	//TODO: How do we manage permissions?
}

type User struct {
	ID          *string    `json:"id,omitempty"`
	Username    *string    `json:"username,omitempty" validate:"required,max=90"`
	DisplayName *string    `json:"displayName,omitempty" validate:"required,max=200"`
	Email       *string    `json:"email,omitempty" validate:"required,email,max=50"`
	Password    *string    `json:"password" validate:"required,min=8,alphanum"`
	Verified    *bool      `json:"verified,omitempty"`
	RoleID      *string    `json:"roleID,omitempty" validate:"required,uuid4"`
	CreatedAt   *time.Time `json:"createdAt,omitzero"`
	Karma       *int       `json:"karma,omitempty" validate:"omitnil,number"`
}
