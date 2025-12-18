package models

import (
	"encoding/json"
	"time"
)

type Patient struct {
	ID        string     `json:"id" db:"id"`
	FirstName string     `json:"first_name" db:"first_name"`
	LastName  string     `json:"last_name" db:"last_name"`
	Email     *string    `json:"email,omitempty" db:"email"`
	Phone     *string    `json:"phone,omitempty" db:"phone"`
	DOB       *time.Time `json:"dob,omitempty" db:"dob"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type Device struct {
	ID        string          `json:"id" db:"id"`
	PatientID string          `json:"patient_id" db:"patient_id"`
	Type      string          `json:"device_type" db:"device_type"`
	Serial    *string         `json:"serial,omitempty" db:"serial"`
	Metadata  json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type Reading struct {
	ID         string    `json:"id" db:"id"`
	PatientID  string    `json:"patient_id" db:"patient_id"`
	DeviceID   *string   `json:"device_id,omitempty" db:"device_id"`
	Metric     string    `json:"metric" db:"metric"`
	Value      *float64  `json:"value,omitempty" db:"value"`
	Unit       *string   `json:"unit,omitempty" db:"unit"`
	RecordedAt time.Time `json:"recorded_at" db:"recorded_at"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
