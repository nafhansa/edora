package handler

import (
    "edora/backend/internal/models"
    st "edora/backend/internal/store"
    "github.com/gofiber/fiber/v2"
)

type DashboardHandler struct{
    store *st.Store
    auth *AuthHandler
}

func NewDashboardHandler(s *st.Store, a *AuthHandler) *DashboardHandler {
    return &DashboardHandler{store: s, auth: a}
}

func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
    user := h.auth.AuthenticatedUser(c)
    if user == nil || user.Role != "medic" {
        return c.Status(403).JSON(fiber.Map{"error":"forbidden"})
    }
    readings := h.store.Readings()
    return c.JSON(fiber.Map{"readings": readings})
}

func (h *DashboardHandler) GetUsers(c *fiber.Ctx) error {
    user := h.auth.AuthenticatedUser(c)
    if user == nil || user.Role != "medic" {
        return c.Status(403).JSON(fiber.Map{"error":"forbidden"})
    }
    users := h.store.Users()
    // convert []models.User to interface for JSON
    out := make([]models.User, 0, len(users))
    for _, u := range users { out = append(out, u) }
    return c.JSON(out)
}
