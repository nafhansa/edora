package service

import (
	"context"
	"testing"
	"time"

	"edora/backend/internal/models"
)

type mockReadingRepoForStats struct{}

func (m *mockReadingRepoForStats) CreateReading(ctx context.Context, rd *models.Reading) (string, error) {
	return "", nil
}
func (m *mockReadingRepoForStats) GetStats(ctx context.Context) (int, map[string]int, error) {
	return 10, map[string]int{"Normal": 7, "Osteoporosis": 3}, nil
}

type mockDeviceRepoForStats struct{}

func (m *mockDeviceRepoForStats) GetBySerial(ctx context.Context, serial string) (*models.Device, error) {
	return nil, nil
}
func (m *mockDeviceRepoForStats) UpdateLastSeen(ctx context.Context, id string, t time.Time) error {
	return nil
}
func (m *mockDeviceRepoForStats) CountActive(ctx context.Context, since time.Duration) (int, error) {
	return 5, nil
}

func TestDashboardService_GetStats(t *testing.T) {
	rr := &mockReadingRepoForStats{}
	dr := &mockDeviceRepoForStats{}
	svc := NewDashboardService(rr, dr)
	st, err := svc.GetStats(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st.TotalToday != 10 {
		t.Fatalf("expected total 10 got %d", st.TotalToday)
	}
	if st.Osteoporosis != 3 {
		t.Fatalf("expected osteoporosis 3 got %d", st.Osteoporosis)
	}
	if st.ActiveDevices != 5 {
		t.Fatalf("expected active devices 5 got %d", st.ActiveDevices)
	}
}
