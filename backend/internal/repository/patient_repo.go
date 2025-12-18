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
		return []models.Patient{}, nil
	}
	db, ok := p.db.(*sql.DB)
	if !ok {
		return nil, errors.New("unsupported db type")
	}
	rows, err := db.QueryContext(ctx, `SELECT id, nik, name, gender, birth_date, address, created_at, updated_at FROM patients ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Patient
	for rows.Next() {
		var pt models.Patient
		var birth sql.NullTime
		if err := rows.Scan(&pt.ID, &pt.NIK, &pt.Name, &pt.Gender, &birth, &pt.Address, &pt.CreatedAt, &pt.UpdatedAt); err != nil {
			return nil, err
		}
		if birth.Valid {
			pt.BirthDate = birth.Time
		}
		out = append(out, pt)
	}
	return out, nil
}
