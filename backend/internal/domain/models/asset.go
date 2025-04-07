package models

import "time"

type Asset struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"` // Links to an S3 bucket
	CreatedAt   time.Time `json:"createdAt"`
}
