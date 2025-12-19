package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"time"

	"edora/backend/internal/models"
)

type PatientRepository struct {
	db *sql.DB
}

func NewPatientRepository(db *sql.DB) *PatientRepository {
	return &PatientRepository{
		db: db,
	}
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// PERBAIKAN DI SINI: models.Patient (bukan models.models.Patient)
func (r *PatientRepository) CreatePatient(ctx context.Context, p *models.Patient) (string, error) {
	if p.ID == "" {
		p.ID = generateID()
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	query := `
		INSERT INTO patients (id, nik, name, gender, birth_date, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		p.ID,
		p.NIK,
		p.Name,
		p.Gender,
		p.BirthDate,
		p.Address,
		p.CreatedAt,
		p.UpdatedAt,
	)

	if err != nil {
		log.Printf("❌ Gagal Insert ke DB: %v", err)
		return "", err
	}

	return p.ID, nil
}

// PERBAIKAN DI SINI JUGA: models.Patient
func (r *PatientRepository) ListPatients(ctx context.Context) ([]models.Patient, error) {
	query := `
		SELECT id, nik, name, gender, birth_date, address, created_at, updated_at 
		FROM patients 
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []models.Patient
	for rows.Next() {
		var p models.Patient
		if err := rows.Scan(
			&p.ID,
			&p.NIK,
			&p.Name,
			&p.Gender,
			&p.BirthDate,
			&p.Address,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *PatientRepository) UpdatePatient(ctx context.Context, p *models.Patient) error {
	p.UpdatedAt = time.Now()
	query := `
		UPDATE patients
		SET nik = $1, name = $2, gender = $3, birth_date = $4, address = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.db.ExecContext(ctx, query,
		p.NIK,
		p.Name,
		p.Gender,
		p.BirthDate,
		p.Address,
		p.UpdatedAt,
		p.ID,
	)
	return err
}

func (r *PatientRepository) DeletePatient(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM patients WHERE id = $1`, id)
	return err
}
