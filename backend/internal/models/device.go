package models

import "time"

type Device struct {
	ID           string    `json:"id" db:"id"`
	SerialNumber string    `json:"serial_number" db:"serial_number"`
	Name         string    `json:"name" db:"name"`
	Status       string    `json:"status" db:"status"`
	LastSeen     time.Time `json:"last_seen" db:"last_seen"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
