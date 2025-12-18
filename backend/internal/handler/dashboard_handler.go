package handler

import (
	"context"

	"edora/backend/internal/models"
	"edora/backend/internal/service"
	st "edora/backend/internal/store"

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

// File-store backed dashboard handler (existing)
type DashboardHandler struct {
	store *st.Store
	auth  *AuthHandler
}

func NewDashboardHandler(s *st.Store, a *AuthHandler) *DashboardHandler {
	return &DashboardHandler{store: s, auth: a}
}

func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
	user := h.auth.AuthenticatedUser(c)
	if user == nil || user.Role != "medic" {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
	readings := h.store.Readings()
	return c.JSON(fiber.Map{"readings": readings})
}

func (h *DashboardHandler) GetUsers(c *fiber.Ctx) error {
	user := h.auth.AuthenticatedUser(c)
	if user == nil || user.Role != "medic" {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
	users := h.store.Users()
	out := make([]models.User, 0, len(users))
	for _, u := range users {
		out = append(out, u)
	}
	return c.JSON(out)
}
