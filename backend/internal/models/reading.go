package models

import (
	"encoding/json"
	"time"
)

type Reading struct {
	ID        string `json:"id" db:"id"`
	DeviceID  string `json:"device_id" db:"device_id"`
	PatientID string `json:"patient_id" db:"patient_id"`
	DoctorID  string `json:"doctor_id" db:"doctor_id"`

	BMDResult      float64 `json:"bmd_result" db:"bmd_result"`
	TScore         float64 `json:"t_score" db:"t_score"`
	Classification string  `json:"classification" db:"classification"`

	RawSignalData json.RawMessage `json:"raw_signal_data" db:"raw_signal_data"`

	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Request Payload untuk Sync dari Mobile App
type CreateReadingRequest struct {
	DeviceSerial   string          `json:"device_serial"`
	PatientID      string          `json:"patient_id"`
	DoctorID       string          `json:"doctor_id"`
	BMDResult      float64         `json:"bmd_result"`
	TScore         float64         `json:"t_score"`
	Classification string          `json:"classification"`
	RawSignalData  json.RawMessage `json:"raw_signal_data"`
	Lat            float64         `json:"lat"`
	Long           float64         `json:"long"`
}
