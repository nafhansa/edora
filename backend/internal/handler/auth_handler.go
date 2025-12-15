package handler

import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "log"

    "edora/backend/internal/store"
    "edora/backend/internal/models"
    "github.com/gofiber/fiber/v2"
)

type AuthHandler struct{
    store *store.Store
    sessions map[string]string // token -> userID
}

func NewAuthHandler(s *store.Store) *AuthHandler {
    return &AuthHandler{store: s, sessions: make(map[string]string)}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var body struct{ Email, Password string }
    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error":"invalid body"})
    }
    u := h.store.FindUserByEmail(body.Email)
    if u == nil {
        log.Printf("login: user not found for %s", body.Email)
    } else {
        log.Printf("login: found user %s with password='%s'", u.Email, u.Password)
    }
    if u == nil || u.Password != body.Password {
        return c.Status(401).JSON(fiber.Map{"error":"invalid credentials"})
    }
    // generate token
    b := make([]byte, 16)
    rand.Read(b)
    token := hex.EncodeToString(b)
    h.sessions[token] = u.ID
    return c.JSON(fiber.Map{"token": token, "role": u.Role})
}

func (h *AuthHandler) AuthenticatedUser(c *fiber.Ctx) *models.User {
    auth := c.Get("Authorization")
    if auth == "" { return nil }
    // expect 'Bearer <token>'
    var token string
    if _, err := fmt.Sscanf(auth, "Bearer %s", &token); err != nil {
        return nil
    }
    uid, ok := h.sessions[token]
    if !ok { return nil }
    // find user
    return h.store.FindUserByID(uid)
}
