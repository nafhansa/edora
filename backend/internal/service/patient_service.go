package service

import (
	"context"

	"edora/backend/internal/models"
	"edora/backend/internal/repository"
)

type PatientService struct {
	repo *repository.PatientRepository
}

func NewPatientService(pr *repository.PatientRepository) *PatientService {
	return &PatientService{repo: pr}
}

func (s *PatientService) CreatePatient(ctx context.Context, pt *models.Patient) (string, error) {
	return s.repo.CreatePatient(ctx, pt)
}

func (s *PatientService) ListPatients(ctx context.Context) ([]models.Patient, error) {
	return s.repo.ListPatients(ctx)
}

func (s *PatientService) UpdatePatient(ctx context.Context, pt *models.Patient) error {
	return s.repo.UpdatePatient(ctx, pt)
}

func (s *PatientService) DeletePatient(ctx context.Context, id string) error {
	return s.repo.DeletePatient(ctx, id)
}
