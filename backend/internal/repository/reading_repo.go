package repository

import (
	"context"
	"database/sql"
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

// CreateReading inserts a reading; if db is nil this is a no-op returning a generated id placeholder.
func (r *ReadingRepository) CreateReading(ctx context.Context, rd *models.Reading) (string, error) {
	if r.db == nil {
		// no DB available in this environment (smoke tests) — return empty id
		return "", nil
	}
	// Expecting *sql.DB for concrete implementation
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
		return 0, map[string]int{}, nil
	}
	db, ok := r.db.(*sql.DB)
	if !ok {
		return 0, nil, errors.New("unsupported db type")
	}

	var total int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM readings WHERE created_at >= date_trunc('day', now())`).Scan(&total); err != nil {
		return 0, nil, err
	}

	rows, err := db.QueryContext(ctx, `SELECT classification, COUNT(*) FROM readings WHERE created_at >= date_trunc('day', now()) GROUP BY classification`)
	if err != nil {
		return total, nil, err
	}
	defer rows.Close()

	m := make(map[string]int)
	for rows.Next() {
		var cls sql.NullString
		var cnt int
		if err := rows.Scan(&cls, &cnt); err != nil {
			return total, nil, err
		}
		key := "unknown"
		if cls.Valid {
			key = cls.String
		}
		m[key] = cnt
	}
	return total, m, nil
}
