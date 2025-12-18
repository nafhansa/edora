package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"edora/backend/internal/handler"
	"edora/backend/internal/repository"
	"edora/backend/internal/service"
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

	// legacy file-store removed; using repository/service pattern

	app := fiber.New()

	// Initialize repositories and services (pgconn may be nil in smoke tests)
	readingRepo := repository.NewReadingRepository(pgconn)
	deviceRepo := repository.NewDeviceRepository(pgconn)
	patientRepo := repository.NewPatientRepository(pgconn)

	readingSvc := service.NewReadingService(readingRepo, deviceRepo)
	dashboardSvc := service.NewDashboardService(readingRepo, deviceRepo)
	patientSvc := service.NewPatientService(patientRepo)
	deviceSvc := service.NewDeviceService(deviceRepo)

	// auth & handlers
	auth := handler.NewAuthHandler()
	// HTTP handlers for new APIs
	readingHandler := handler.NewReadingHandler(readingSvc)
	dashHTTP := handler.NewDashboardHTTPHandler(dashboardSvc)
	patientHandler := handler.NewPatientHandler(patientSvc)
	deviceHandler := handler.NewDeviceHandler(deviceSvc)

	api := app.Group("/api/v1")
	api.Post("/login", auth.Login)
	// product APIs removed per refactor
	api.Post("/sync/reading", readingHandler.SyncReading)
	api.Get("/dashboard/stats", dashHTTP.Stats)

	// Patient & Device management
	api.Get("/patients", patientHandler.List)
	api.Post("/patients", patientHandler.Create)
	api.Get("/devices", deviceHandler.List)

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
