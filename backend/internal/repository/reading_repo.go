package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"sync" // Tambah sync untuk keamanan data
	"time"

	"edora/backend/internal/models"
)

type ReadingRepository struct {
	db interface{}
	// --- SMART MOCK STORAGE ---
	mu                 sync.Mutex
	mockReadings       []models.Reading
	mockMedicalRecords []models.MedicalRecord
}

func NewReadingRepository(db interface{}) *ReadingRepository {
	// Kita isi data awal (Dummy) supaya dashboard tidak kosong pas pertama buka
	initialData := []models.Reading{}
	if db == nil {
		initialData = []models.Reading{
			{Classification: "Normal", CreatedAt: time.Now()},
			{Classification: "Normal", CreatedAt: time.Now()},
			{Classification: "Osteopenia", CreatedAt: time.Now()},
		}
	}

	// contoh data medical record mock
	initialMR := []models.MedicalRecord{}
	if db == nil {
		initialMR = []models.MedicalRecord{
			{ID: 1, PatientID: "patient-1", TScore: -0.5, Diagnosis: "Normal", ScanDate: time.Now(), Notes: "demo"},
		}
	}

	return &ReadingRepository{
		db:                 db,
		mockReadings:       initialData,
		mockMedicalRecords: initialMR,
	}
}

type ReadingRepo interface {
	CreateReading(ctx context.Context, rd *models.Reading) (string, error)
	GetStats(ctx context.Context) (int, map[string]int, error)
	CreateMedicalRecord(ctx context.Context, mr *models.MedicalRecord) (int, error)
	GetPatientRecords(ctx context.Context, patientID string) ([]models.MedicalRecord, error)
}

func (r *ReadingRepository) CreateReading(ctx context.Context, rd *models.Reading) (string, error) {
	// MODE MOCK (DB MATI)
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()

		// Generate ID
		b := make([]byte, 16)
		rand.Read(b)
		rd.ID = hex.EncodeToString(b)

		// Kalau timestamp kosong, isi sekarang
		if rd.CreatedAt.IsZero() {
			rd.CreatedAt = time.Now()
		}

		// SIMPAN KE MEMORI (RAM)
		r.mockReadings = append(r.mockReadings, *rd)

		return rd.ID, nil
	}

	// MODE REAL DB
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

func (r *ReadingRepository) GetStats(ctx context.Context) (int, map[string]int, error) {
	// MODE MOCK (HITUNG DARI MEMORI)
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()

		totalToday := 0
		stats := map[string]int{
			"Normal":       0,
			"Osteopenia":   0,
			"Osteoporosis": 0,
		}

		// Mulai hitung jam 00:00 hari ini
		now := time.Now()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		for _, reading := range r.mockReadings {
			// Hanya hitung data hari ini
			if reading.CreatedAt.After(startOfDay) || reading.CreatedAt.Equal(startOfDay) {
				totalToday++
				// Normalisasi string (biar "osteoporosis" dan "Osteoporosis" sama)
				classKey := strings.Title(strings.ToLower(reading.Classification))
				if _, exists := stats[classKey]; exists {
					stats[classKey]++
				} else {
					// Fallback kalau ada klasifikasi aneh
					stats[classKey] = 1
				}
			}
		}

		return totalToday, stats, nil
	}

	// MODE REAL DB
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

// CreateMedicalRecord menyimpan record medis (mock atau DB)
func (r *ReadingRepository) CreateMedicalRecord(ctx context.Context, mr *models.MedicalRecord) (int, error) {
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()

		id := len(r.mockMedicalRecords) + 1
		mr.ID = id
		if mr.ScanDate.IsZero() {
			mr.ScanDate = time.Now()
		}
		r.mockMedicalRecords = append(r.mockMedicalRecords, *mr)
		return id, nil
	}

	db, ok := r.db.(*sql.DB)
	if !ok {
		return 0, errors.New("unsupported db type")
	}

	q := `INSERT INTO medical_records (patient_id, t_score, diagnosis, notes, scan_date) VALUES ($1,$2,$3,$4,$5) RETURNING id, scan_date`
	if mr.ScanDate.IsZero() {
		mr.ScanDate = time.Now()
	}
	var id int
	var scanDate time.Time
	if err := db.QueryRowContext(ctx, q, mr.PatientID, mr.TScore, mr.Diagnosis, mr.Notes, mr.ScanDate).Scan(&id, &scanDate); err != nil {
		return 0, err
	}
	mr.ID = id
	mr.ScanDate = scanDate
	return id, nil
}

// GetPatientRecords mengambil riwayat medical records untuk pasien
func (r *ReadingRepository) GetPatientRecords(ctx context.Context, patientID string) ([]models.MedicalRecord, error) {
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()

		out := []models.MedicalRecord{}
		for _, mr := range r.mockMedicalRecords {
			if mr.PatientID == patientID {
				out = append(out, mr)
			}
		}
		return out, nil
	}

	db, ok := r.db.(*sql.DB)
	if !ok {
		return nil, errors.New("unsupported db type")
	}

	rows, err := db.QueryContext(ctx, `SELECT id, t_score, diagnosis, scan_date, notes FROM medical_records WHERE patient_id = $1 ORDER BY scan_date DESC`, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MedicalRecord
	for rows.Next() {
		var rcd models.MedicalRecord
		rcd.PatientID = patientID
		if err := rows.Scan(&rcd.ID, &rcd.TScore, &rcd.Diagnosis, &rcd.ScanDate, &rcd.Notes); err != nil {
			continue
		}
		records = append(records, rcd)
	}
	return records, nil
}
