package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"edora/backend/internal/service"
)

type DeviceHandler struct {
	svc *service.DeviceService
}

func NewDeviceHandler(s *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{svc: s}
}

// List returns a simple summary: number of active devices in the last 5 minutes.
func (h *DeviceHandler) List(c *fiber.Ctx) error {
	since := 5 * time.Minute
	cnt, err := h.svc.CountActive(context.Background(), since)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"active_count": cnt})
}
