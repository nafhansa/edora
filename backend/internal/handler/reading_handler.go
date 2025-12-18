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
