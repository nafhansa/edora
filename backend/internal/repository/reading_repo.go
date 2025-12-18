package repository

import (
	"context"
	"crypto/rand" // Tambahan untuk random ID
	"database/sql"
	"encoding/hex" // Tambahan untuk random ID
	"encoding/json"
	"errors"

	"edora/backend/internal/models"
)

type ReadingRepository struct {
	db interface{}
}

func NewReadingRepository(db interface{}) *ReadingRepository {
	return &ReadingRepository{db: db}
}

// ReadingRepo defines methods used by services and tests.
type ReadingRepo interface {
	CreateReading(ctx context.Context, rd *models.Reading) (string, error)
	GetStats(ctx context.Context) (int, map[string]int, error)
}

// CreateReading inserts a reading; if db is nil this returns a generated Mock ID.
func (r *ReadingRepository) CreateReading(ctx context.Context, rd *models.Reading) (string, error) {
	// JIKA DB TIDAK KONEK (NIL): Generate Dummy ID agar Mobile App happy
	if r.db == nil {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		mockID := hex.EncodeToString(b)
		return mockID, nil
	}

	// JIKA DB KONEK: Lakukan Insert sesungguhnya
	db, ok := r.db.(*sql.DB)
	if !ok {
		return "", errors.New("unsupported db type")
	}

	q := `INSERT INTO readings (device_id, patient_id, doctor_id, bmd_result, t_score, classification, raw_signal_data, latitude, longitude, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`
	raw := rd.RawSignalData
	if len(raw) == 0 {
		raw = json.RawMessage("[]")
	}
	var id string
	if err := db.QueryRowContext(ctx, q, rd.DeviceID, rd.PatientID, rd.DoctorID, rd.BMDResult, rd.TScore, rd.Classification, raw, rd.Latitude, rd.Longitude, rd.CreatedAt).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

// GetStats returns total scans today and counts grouped by classification
func (r *ReadingRepository) GetStats(ctx context.Context) (int, map[string]int, error) {
	if r.db == nil {
		// Mock Data Stats untuk Dashboard
		return 15, map[string]int{"Normal": 10, "Osteopenia": 3, "Osteoporosis": 2}, nil
	}
	db, ok := r.db.(*sql.DB)
	if !ok {
		return 0, nil, errors.New("unsupported db type")
	}

	var total int
	// Hitung total hari ini
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM readings WHERE created_at >= date_trunc('day', now())`).Scan(&total); err != nil {
		return 0, nil, err
	}

	// Hitung breakdown per klasifikasi
	rows, err := db.QueryContext(ctx, `SELECT classification, COUNT(*) FROM readings WHERE created_at >= date_trunc('day', now()) GROUP BY classification`)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var cls string
		var cnt int
		if err := rows.Scan(&cls, &cnt); err != nil {
			return 0, nil, err
		}
		stats[cls] = cnt
	}
	return total, stats, nil
}
