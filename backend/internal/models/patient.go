package models

import "time"

type Patient struct {
	ID        string    `json:"id" db:"id"`
	NIK       string    `json:"nik" db:"nik"`
	Name      string    `json:"name" db:"name"`
	Gender    string    `json:"gender" db:"gender"`
	BirthDate time.Time `json:"birth_date" db:"birth_date"`
	Address   string    `json:"address" db:"address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Tambahkan baris ini agar error hilang:
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type MedicalRecord struct {
	ID        int       `json:"id"`
	PatientID string    `json:"patient_id"` 
	TScore    float64   `json:"t_score"`
	Diagnosis string    `json:"diagnosis"`
	ScanDate  time.Time `json:"scan_date"`
	Notes     string    `json:"notes"`
}
