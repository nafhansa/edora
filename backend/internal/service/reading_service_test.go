package service

import (
	"context"
	"testing"
	"time"

	"edora/backend/internal/models"
)

// Note: use simple mocks implementing repository interfaces
type mockReadingRepo struct{}

func (m *mockReadingRepo) CreateReading(ctx context.Context, rd *models.Reading) (string, error) {
	return "rid-123", nil
}
func (m *mockReadingRepo) GetStats(ctx context.Context) (int, map[string]int, error) {
	return 0, map[string]int{}, nil
}

type mockDeviceRepo struct{ dev *models.Device }

func (m *mockDeviceRepo) GetBySerial(ctx context.Context, serial string) (*models.Device, error) {
	return m.dev, nil
}
func (m *mockDeviceRepo) UpdateLastSeen(ctx context.Context, id string, t time.Time) error {
	return nil
}
func (m *mockDeviceRepo) CountActive(ctx context.Context, since time.Duration) (int, error) {
	return 0, nil
}

func TestReadingService_DeviceNotFound(t *testing.T) {
	rr := &mockReadingRepo{}
	dr := &mockDeviceRepo{dev: nil}
	svc := NewReadingService(rr, dr)
	rd := &models.Reading{}
	_, err := svc.SyncReading(context.Background(), rd, "unknown-serial")
	if err == nil {
		t.Fatalf("expected error when device not found")
	}
}

func TestReadingService_Success(t *testing.T) {
	rr := &mockReadingRepo{}
	dr := &mockDeviceRepo{dev: &models.Device{ID: "dev-1"}}
	svc := NewReadingService(rr, dr)
	rd := &models.Reading{PatientID: "p1", DoctorID: "d1"}
	id, err := svc.SyncReading(context.Background(), rd, "ESP32-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "rid-123" {
		t.Fatalf("unexpected id: %s", id)
	}
}
