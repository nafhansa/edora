package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"edora/backend/internal/handler"
	"edora/backend/internal/store"
	"edora/backend/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load envs
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://user:pass@db:5432/appdb"
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	ctx := context.Background()

	// Connect to Postgres (stubbed in this build)
	pgconn, err := database.Connect(ctx, dbURL)
	if err != nil {
		log.Printf("postgres connect error: %v", err)
	}

	// Redis/client is optional; we run without Redis in smoke tests.

	// Initialize store (file-based) for users/readings
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	dataDir := filepath.Join(exeDir, "data")
	st, serr := store.New(dataDir)
	if serr != nil {
		log.Printf("store init error: %v", serr)
	}

	app := fiber.New()

	// auth & dashboard handlers (file-store backed)
	auth := handler.NewAuthHandler(st)
	dash := handler.NewDashboardHandler(st, auth)

	api := app.Group("/api/v1")
	api.Post("/login", auth.Login)
	api.Get("/dashboard", dash.GetDashboard)
	api.Get("/users", dash.GetUsers)
	api.Get("/debug/users", func(c *fiber.Ctx) error {
		us := st.Users()
		return c.JSON(us)
	})
	api.Post("/debug/users", func(c *fiber.Ctx) error {
		us := st.Users()
		return c.JSON(us)
	})
	// product APIs removed per refactor

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Print simple health info
	info := map[string]string{"status": "ok", "port": port}
	b, _ := json.Marshal(info)
	fmt.Printf("Starting backend: %s\n", string(b))

	// Start server
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Block until terminated or error
	for {
		time.Sleep(10 * time.Second)
	}
}
