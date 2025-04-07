package models

import (
	"time"

	"github.com/mikespook/gorbac"
)

type Session struct {
	ID     string `json:"id"`
	UserID string `json:"userID"`
}

type Role struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Permissions []gorbac.Permission `json:"permissions"`
	//TODO: How do we manage permissions?
}

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Display_name string    `json:"displayName"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	Salt         string    `json:"-"`
	Role_id      string    `json:"roleID"`
	Created_at   time.Time `json:"createdAt"`
	Karma        int       `json:"karma"`
}
