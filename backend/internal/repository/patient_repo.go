package repository

import (
	"context"
	"crypto/rand" // Tambahan: untuk generate ID acak
	"database/sql"
	"encoding/hex" // Tambahan: untuk format ID
	"errors"
	"time"

	"edora/backend/internal/models"
)

type PatientRepository struct {
	db interface{}
}

func NewPatientRepository(db interface{}) *PatientRepository {
	return &PatientRepository{db: db}
}

// PatientRepo defines methods used by services and tests.
type PatientRepo interface {
	CreatePatient(ctx context.Context, pt *models.Patient) (string, error)
	ListPatients(ctx context.Context) ([]models.Patient, error)
}

func (p *PatientRepository) CreatePatient(ctx context.Context, pt *models.Patient) (string, error) {
	// Set CreatedAt/UpdatedAt jika belum ada
	if pt.CreatedAt.IsZero() {
		pt.CreatedAt = time.Now().UTC()
		pt.UpdatedAt = pt.CreatedAt
	}

	// JIKA DB TIDAK KONEK (NIL): Generate Dummy ID agar Frontend tidak error
	if p.db == nil {
		// Generate random ID (simulasi UUID)
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		mockID := hex.EncodeToString(b)

		// Kita kembalikan Mock ID ini
		return mockID, nil
	}

	// JIKA DB KONEK: Lakukan Insert sesungguhnya
	db, ok := p.db.(*sql.DB)
	if !ok {
		return "", errors.New("unsupported db type")
	}

	q := `INSERT INTO patients (nik, name, gender, birth_date, address, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`
	var id string
	if err := db.QueryRowContext(ctx, q, pt.NIK, pt.Name, pt.Gender, pt.BirthDate, pt.Address, pt.CreatedAt, pt.UpdatedAt).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

// ListPatients returns all patients. In smoke-test mode returns empty slice.
func (p *PatientRepository) ListPatients(ctx context.Context) ([]models.Patient, error) {
	if p.db == nil {
		// --- MOCK DATA START ---
		// Kita buat 3 pasien dummy untuk testing tampilan di HP
		mockPatients := []models.Patient{
			{
				ID:        "mock-pat-001",
				NIK:       "3204010101900001",
				Name:      "Budi Santoso",
				Gender:    "male",
				BirthDate: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
				Address:   "Jl. Merdeka No. 10, Bandung",
				CreatedAt: time.Now().Add(-24 * time.Hour), // Dibuat kemarin
			},
			{
				ID:        "mock-pat-002",
				NIK:       "3204010202850002",
				Name:      "Siti Aminah",
				Gender:    "female",
				BirthDate: time.Date(1985, 5, 20, 0, 0, 0, 0, time.UTC),
				Address:   "Komp. Setiabudi Regency Blok A",
				CreatedAt: time.Now().Add(-48 * time.Hour),
			},
			{
				ID:        "mock-pat-003",
				NIK:       "3204011212700003",
				Name:      "Hartono",
				Gender:    "male",
				BirthDate: time.Date(1970, 12, 12, 0, 0, 0, 0, time.UTC),
				Address:   "Jl. Dago Atas No. 99",
				CreatedAt: time.Now().Add(-72 * time.Hour),
			},
		}
		return mockPatients, nil
		// --- MOCK DATA END ---
	}

	// Logic DB Asli (PostgreSQL)
	db, ok := p.db.(*sql.DB)
	if !ok {
		return nil, errors.New("unsupported db type")
	}
	rows, err := db.QueryContext(ctx, `SELECT id, nik, name, gender, birth_date, address, created_at, updated_at FROM patients ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []models.Patient
	for rows.Next() {
		var pt models.Patient
		// Scan updated_at juga
		if err := rows.Scan(&pt.ID, &pt.NIK, &pt.Name, &pt.Gender, &pt.BirthDate, &pt.Address, &pt.CreatedAt, &pt.UpdatedAt); err != nil {
			return nil, err
		}
		patients = append(patients, pt)
	}
	return patients, nil
}
