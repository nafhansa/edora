package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"edora/backend/internal/models"
	"edora/backend/internal/service"
)

type PatientHandler struct {
	svc *service.PatientService
}

func NewPatientHandler(s *service.PatientService) *PatientHandler {
	return &PatientHandler{svc: s}
}

// Struct khusus untuk menangkap Request (DTO)
// BirthDate kita set string agar "1990-01-01" bisa masuk tanpa error
type createPatientRequest struct {
	NIK       string `json:"nik"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"` // string, bukan time.Time
	Address   string `json:"address"`
}

func (h *PatientHandler) Create(c *fiber.Ctx) error {
	// 1. Gunakan struct request sementara
	var req createPatientRequest

	// Parse Body
	if err := c.BodyParser(&req); err != nil {
		// Tampilkan error detail supaya kamu tahu kalau ada salah ketik JSON
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid body: " + err.Error(),
		})
	}

	// 2. Konversi String Tanggal ("1990-01-01") ke time.Time
	parsedDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid date format, use YYYY-MM-DD",
		})
	}

	// 3. Pindahkan data ke Model asli (models.Patient)
	pt := models.Patient{
		NIK:       req.NIK,
		Name:      req.Name,
		Gender:    req.Gender,
		BirthDate: parsedDate, // Masukkan hasil parsing
		Address:   req.Address,
		// ID & CreatedAt dihandle service/DB
	}

	// 4. Panggil Service
	id, err := h.svc.CreatePatient(context.Background(), &pt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

func (h *PatientHandler) List(c *fiber.Ctx) error {
	pts, err := h.svc.ListPatients(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(pts)
}

// Update handles updating an existing patient
func (h *PatientHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "patient id required"})
	}

	var req createPatientRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body: " + err.Error()})
	}

	parsedDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid date format, use YYYY-MM-DD"})
	}

	pt := models.Patient{
		ID:        id,
		NIK:       req.NIK,
		Name:      req.Name,
		Gender:    req.Gender,
		BirthDate: parsedDate,
		Address:   req.Address,
	}

	if err := h.svc.UpdatePatient(context.Background(), &pt); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Delete handles removing a patient and cascade removing medical records
func (h *PatientHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "patient id required"})
	}
	if err := h.svc.DeletePatient(context.Background(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
