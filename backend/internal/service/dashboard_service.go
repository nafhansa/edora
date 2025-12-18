package service

import (
	"context"
	"time"

	"edora/backend/internal/repository"
)

type DashboardService struct {
	readingRepo repository.ReadingRepo
	deviceRepo  repository.DeviceRepo
}

type DashboardStats struct {
	TotalToday       int            `json:"totalPatientsToday"`
	ByClassification map[string]int `json:"byClassification"`
	ActiveDevices    int            `json:"activeDevices"`
	Osteoporosis     int            `json:"osteoporosisCases"`
	RecentScans      any            `json:"recentScans"`
}

func NewDashboardService(rr repository.ReadingRepo, dr repository.DeviceRepo) *DashboardService {
	return &DashboardService{readingRepo: rr, deviceRepo: dr}
}

func (s *DashboardService) GetStats(ctx context.Context) (*DashboardStats, error) {
	total, byClass, err := s.readingRepo.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	ost := 0
	if v, ok := byClass["Osteoporosis"]; ok {
		ost = v
	}
	// count active devices in last 5 minutes
	active := 0
	if s.deviceRepo != nil {
		cnt, err := s.deviceRepo.CountActive(ctx, 5*time.Minute)
		if err == nil {
			active = cnt
		}
	}

	return &DashboardStats{
		TotalToday:       total,
		ByClassification: byClass,
		ActiveDevices:    active,
		Osteoporosis:     ost,
		RecentScans:      []any{},
	}, nil
}
