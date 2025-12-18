package handler

import (
	"context"

	"edora/backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

// HTTP-facing dashboard handler (aggregates from DB)
type DashboardHTTPHandler struct {
	ds *service.DashboardService
}

func NewDashboardHTTPHandler(ds *service.DashboardService) *DashboardHTTPHandler {
	return &DashboardHTTPHandler{ds: ds}
}

func (h *DashboardHTTPHandler) Stats(c *fiber.Ctx) error {
	stt, err := h.ds.GetStats(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stt)
}

// Note: legacy file-store backed DashboardHandler removed as part of store refactor.
