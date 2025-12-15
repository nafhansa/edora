package models

import "time"

type User struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    Password  string    `json:"password"` // plain text for prototype — replace with hash for production
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}
