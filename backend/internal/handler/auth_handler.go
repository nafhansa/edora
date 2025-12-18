package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"edora/backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	sessions map[string]models.User // token -> user
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{sessions: make(map[string]models.User)}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body struct{ Email, Password string }
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	// Simple hardcoded credential for now
	if body.Email != "doctor@edora.com" || body.Password != "password" {
		log.Printf("login failed for %s", body.Email)
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	u := models.User{
		ID:        "doctor-0001",
		Email:     body.Email,
		Password:  body.Password,
		Role:      "medic",
		CreatedAt: time.Now().UTC(),
	}

	// generate token
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	token := hex.EncodeToString(b)
	h.sessions[token] = u
	return c.JSON(fiber.Map{"token": token, "role": u.Role})
}

func (h *AuthHandler) AuthenticatedUser(c *fiber.Ctx) *models.User {
	auth := c.Get("Authorization")
	if auth == "" {
		return nil
	}
	// expect 'Bearer <token>'
	var token string
	if _, err := fmt.Sscanf(auth, "Bearer %s", &token); err != nil {
		return nil
	}
	u, ok := h.sessions[token]
	if !ok {
		return nil
	}
	return &u
}
