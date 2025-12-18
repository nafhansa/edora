package handler
//go:build ignore
// This handler was removed during refactor. Kept only for history.

package handler

// product handler removed

import (
    "context"
    "time"

    "github.com/gofiber/fiber/v2"
    "edora/backend/internal/service"
)

type Handler struct{
    prodSvc *service.ProductService
}

func NewHandler(prod *service.ProductService) *Handler {
    return &Handler{prodSvc: prod}
}

func (h *Handler) GetProducts(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    products, err := h.prodSvc.List(ctx, 100)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(products)
}

import (
    "context"
    "time"

    "edora/backend/internal/service"

    "github.com/gofiber/fiber/v2"
)

type Handler struct{
    prodSvc *service.ProductService
}

func NewHandler(prod *service.ProductService) *Handler {
    return &Handler{prodSvc: prod}
}

func (h *Handler) GetProducts(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    products, err := h.prodSvc.List(ctx, 100)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(products)
}
