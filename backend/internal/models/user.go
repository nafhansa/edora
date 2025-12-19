package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"` // stored as bcrypt hash
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
