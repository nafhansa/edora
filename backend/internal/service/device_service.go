package service

import (
	"context"
	"time"

	"edora/backend/internal/repository"
)

type DeviceService struct {
	repo *repository.DeviceRepository
}

func NewDeviceService(dr *repository.DeviceRepository) *DeviceService {
	return &DeviceService{repo: dr}
}

// CountActive returns number of devices active within the given duration.
func (s *DeviceService) CountActive(ctx context.Context, since time.Duration) (int, error) {
	return s.repo.CountActive(ctx, since)
}
