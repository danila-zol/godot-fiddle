package database

import "time"

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Display_name string    `json:"display_name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Role_id      string    `json:"role_id"`
	Created_at   time.Time `json:"created_at"`
	Karma        uint      `json:"karma"`
}
