package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"

	"edora/backend/internal/models"
	"edora/backend/internal/service"
)

type ReadingHandler struct {
	rs *service.ReadingService
}

func NewReadingHandler(rs *service.ReadingService) *ReadingHandler {
	return &ReadingHandler{rs: rs}
}

type syncPayload struct {
	DeviceSerial   string          `json:"device_serial"`
	PatientID      string          `json:"patient_id"`
	DoctorID       string          `json:"doctor_id"`
	BMDResult      float64         `json:"bmd_result"`
	TScore         float64         `json:"t_score"`
	Classification string          `json:"classification"`
	RawSignalData  json.RawMessage `json:"raw_signal_data"`
	Lat            float64         `json:"lat"`
	Long           float64         `json:"long"`
	Timestamp      string          `json:"timestamp"`
}

func (h *ReadingHandler) SyncReading(c *fiber.Ctx) error {
	var p syncPayload
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	rd := &models.Reading{
		PatientID:      p.PatientID,
		DoctorID:       p.DoctorID,
		BMDResult:      p.BMDResult,
		TScore:         p.TScore,
		Classification: p.Classification,
		RawSignalData:  p.RawSignalData,
		Latitude:       p.Lat,
		Longitude:      p.Long,
	}
	if p.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, p.Timestamp); err == nil {
			rd.CreatedAt = t
		} else {
			rd.CreatedAt = time.Now().UTC()
		}
	} else {
		rd.CreatedAt = time.Now().UTC()
	}

	id, err := h.rs.SyncReading(context.Background(), rd, p.DeviceSerial)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

// --- FUNGSI LOGIKA DIAGNOSIS (Sesuai Standar WHO) ---
func determineDiagnosis(score float64) string {
	if score >= -1.0 {
		return "Normal"
	} else if score > -2.5 {
		return "Osteopenia"
	}
	return "Osteoporosis"
}

// CreateMedicalRecord handler untuk menyimpan hasil scan/medical record
func (h *ReadingHandler) CreateMedicalRecord(c *fiber.Ctx) error {
	var input models.MedicalRecord
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	// If client posted to /patients/:id/medical_records without patient_id in body,
	// allow using the path parameter as the patient id.
	if input.PatientID == "" {
		if pid := c.Params("id"); pid != "" {
			input.PatientID = pid
		}
	}

	// 1. Hitung diagnosis otomatis
	input.Diagnosis = determineDiagnosis(input.TScore)
	if input.ScanDate.IsZero() {
		input.ScanDate = time.Now().UTC()
	}

	// 2. Simpan via service
	mr, err := h.rs.CreateMedicalRecord(context.Background(), &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(mr)
}

// GetPatientRecords handler untuk melihat riwayat scan pasien
func (h *ReadingHandler) GetPatientRecords(c *fiber.Ctx) error {
	patientID := c.Params("id")
	if patientID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "patient id required"})
	}

	records, err := h.rs.GetPatientRecords(context.Background(), patientID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if records == nil {
		records = []models.MedicalRecord{}
	}
	return c.Status(fiber.StatusOK).JSON(records)
}
