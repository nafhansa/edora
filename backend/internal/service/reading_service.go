package service

import (
	"context"
	"errors"
	"time"

	"edora/backend/internal/models"
	"edora/backend/internal/repository"
)

type ReadingService struct {
	readingRepo repository.ReadingRepo
	deviceRepo  repository.DeviceRepo
}

func NewReadingService(rr repository.ReadingRepo, dr repository.DeviceRepo) *ReadingService {
	return &ReadingService{readingRepo: rr, deviceRepo: dr}
}

// SyncReading validates device serial, inserts reading and updates device last seen
func (s *ReadingService) SyncReading(ctx context.Context, rd *models.Reading, deviceSerial string) (string, error) {
	if deviceSerial == "" {
		return "", errors.New("device_serial required")
	}
	dev, err := s.deviceRepo.GetBySerial(ctx, deviceSerial)
	if err != nil {
		return "", err
	}
	if dev == nil {
		return "", errors.New("device not registered")
	}

	rd.DeviceID = dev.ID
	if rd.CreatedAt.IsZero() {
		rd.CreatedAt = time.Now().UTC()
	}

	id, err := s.readingRepo.CreateReading(ctx, rd)
	if err != nil {
		return "", err
	}

	// best-effort update last seen
	_ = s.deviceRepo.UpdateLastSeen(ctx, dev.ID, rd.CreatedAt)
	return id, nil
}

// CreateMedicalRecord membuat medical record baru melalui repository
func (s *ReadingService) CreateMedicalRecord(ctx context.Context, mr *models.MedicalRecord) (*models.MedicalRecord, error) {
	_, err := s.readingRepo.CreateMedicalRecord(ctx, mr)
	if err != nil {
		return nil, err
	}
	return mr, nil
}

// GetPatientRecords mengembalikan semua medical record untuk pasien
func (s *ReadingService) GetPatientRecords(ctx context.Context, patientID string) ([]models.MedicalRecord, error) {
	return s.readingRepo.GetPatientRecords(ctx, patientID)
}
