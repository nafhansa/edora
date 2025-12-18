package main

import (
	"context"
	"database/sql" // <--- TAMBAHAN PENTING
	"log"
	"os"

	"edora/backend/internal/handler"
	"edora/backend/internal/repository"
	"edora/backend/internal/service"
	"edora/backend/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1. Load Environment Variables
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://user:pass@db:5432/appdb?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	// 2. Connect to Postgres (DATABASE ASLI) 🔌
	log.Println("🔌 Menghubungkan ke Database...", dbURL)
	pgconn, err := database.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("❌ FATAL: Gagal connect ke Database: %v\nCek apakah container 'db' sudah nyala?", err)
	}
	log.Println("✅ Berhasil terhubung ke Database Postgres!")

	// ---------------------------------------------------------
	// 🔥 PERBAIKAN UTAMA DI SINI (TYPE ASSERTION) 🔥
	// Kita pastikan pgconn benar-benar *sql.DB sebelum dipakai
	// ---------------------------------------------------------
	sqlDB, ok := pgconn.(*sql.DB)
	if !ok {
		log.Fatal("❌ Error: pgconn bukan tipe *sql.DB yang valid")
	}

	// 3. Setup Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Edora Health Backend",
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// 4. Initialize Dependency Injection (Wiring)

	// Repo lain mungkin masih pakai interface pgconn (biarkan dulu jika belum direfactor)
	readingRepo := repository.NewReadingRepository(pgconn)
	deviceRepo := repository.NewDeviceRepository(pgconn)

	// Repo Patient SUDAH pakai sqlDB yang baru kita convert di atas
	patientRepo := repository.NewPatientRepository(sqlDB)

	// Service Layer
	readingSvc := service.NewReadingService(readingRepo, deviceRepo)
	dashboardSvc := service.NewDashboardService(readingRepo, deviceRepo)
	patientSvc := service.NewPatientService(patientRepo)
	deviceSvc := service.NewDeviceService(deviceRepo)

	// Handler Layer
	auth := handler.NewAuthHandler()
	readingHandler := handler.NewReadingHandler(readingSvc)
	dashHTTP := handler.NewDashboardHTTPHandler(dashboardSvc)
	patientHandler := handler.NewPatientHandler(patientSvc)
	deviceHandler := handler.NewDeviceHandler(deviceSvc)

	// 5. Define Routes
	api := app.Group("/api/v1")

	// Auth
	api.Post("/login", auth.Login)

	// Dashboard & IoT Sync
	api.Post("/sync/reading", readingHandler.SyncReading)
	api.Get("/dashboard/stats", dashHTTP.Stats)

	// Patient Management (CRUD)
	api.Get("/patients", patientHandler.List)
	api.Post("/patients", patientHandler.Create)

	// Device Management
	api.Get("/devices", deviceHandler.List)

	// 6. Start Server
	log.Printf("🚀 Server Edora berjalan di port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
