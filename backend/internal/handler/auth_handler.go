package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"edora/backend/internal/models"
	"edora/backend/internal/repository"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	sessions map[string]models.User // token -> user
	users    repository.UserRepo
}

func NewAuthHandler(users repository.UserRepo) *AuthHandler {
	return &AuthHandler{sessions: make(map[string]models.User), users: users}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body struct{ Username, Password string }
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	// Lookup user from DB
	u, err := h.users.GetByUsername(context.Background(), body.Username)
	if err != nil {
		log.Printf("login lookup error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "internal"})
	}
	if u == nil {
		log.Printf("login failed for %s: user not found", body.Username)
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// Verify password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(body.Password)); err != nil {
		log.Printf("login failed for %s: bad password", body.Username)
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// generate token
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	token := hex.EncodeToString(b)

	// store user without password in session map
	sessionUser := *u
	sessionUser.Password = ""
	h.sessions[token] = sessionUser

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
