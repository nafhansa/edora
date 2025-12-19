LATEST PROGRESS (concise)
=========================

Summary
-------
- Status: development server rebuilt and running; example POSTs returned an `id`.
- Recent fix: `internal/handler/patient_handler.go` now accepts `birth_date` as `"YYYY-MM-DD"` and parses it to `time.Time`.

What changed
------------
- Patient creation: DTO + `time.Parse("2006-01-02", ...)` to fix JSON date parsing.
- New endpoints added for storing and retrieving medical scan records.

Key routes (prefix `/api/v1`)
-----------------------------
LATEST PROGRESS — DETAILED FILE MAP
=================================

Overview
--------
- Status: development server rebuilt and running. Key fixes and new endpoints added (patient birth date parsing, medical record handlers).
- This document lists important files, their locations, and a short description of purpose to help navigation and maintenance.

Project root (important files)
------------------------------
- `docker-compose.yml` — Docker Compose config for services (backend, db, frontend etc.).
- `README.md` — project overview and setup notes.
- `LATEST_PROGRESS.md` — this file (progress + file map and quick commands).

Backend: `backend/` (Go)
------------------------
- `backend/cmd/api/main.go` — Application entry point. Creates Fiber app, wires repositories, services and handlers, and registers routes under `/api/v1`.
- `backend/go.mod` — Go module manifest.

Backend: internal packages
--------------------------
- `backend/internal/handler/` — HTTP handlers (Fiber). Key files:
  - `auth_handler.go`       — Authentication endpoints (login).
  - `dashboard_handler.go`  — Dashboard UI serving + stats endpoints.
  - `device_handler.go`     — Device listing and device-related endpoints.
  - `patient_handler.go`    — Patient CRUD endpoints. Contains `createPatientRequest` DTO and `time.Parse` for `birth_date`.
  - `reading_handler.go`    — Reading and scan handlers. Contains `SyncReading` for device sync; added `CreateMedicalRecord`, `GetPatientRecords`, and `determineDiagnosis` logic.
  - `product_handler.go`    — Product endpoints (if applicable).

- `backend/internal/models/` — Domain models used across repo. Key files:
  - `patient.go`            — `Patient` struct and `MedicalRecord` struct (fields: `ID`, `PatientID`, `TScore`, `Diagnosis`, `ScanDate`, `Notes`).
  - `reading.go`            — Reading model for IoT sync (t_score, classification, raw_signal_data, etc.).
  - `device.go`, `user.go`  — Device and User models.

- `backend/internal/repository/` — Repository layer (DB access). Key files:
  - `reading_repo.go`       — Reading repository. Supports mock mode (in-memory) and real DB mode. Implements `CreateReading`, `GetStats`, and newly added `CreateMedicalRecord` and `GetPatientRecords`.
  - `patient_repo.go`       — Patient repository (CRUD using SQL, returns UUID IDs).
  - `device_repo.go`        — Device repository (GetBySerial, UpdateLastSeen, etc.).

- `backend/internal/service/` — Business logic / use-cases. Key files:
  - `reading_service.go`    — Coordinates device validation and reading persistence. Added service methods `CreateMedicalRecord` and `GetPatientRecords` to wrap repository.
  - `patient_service.go`    — Patient-related business logic.
  - `dashboard_service.go`  — Aggregation logic for stats.

Backend: pkg and utilities
--------------------------
- `backend/pkg/database/db.go` — Database connection helper (Connect function returning DB connection).
- `backend/pkg/redis/redis.go` — Redis helper (if used).

Backend: SQL migrations and init
-------------------------------
- `backend/sql/init.sql` — (added) Schema initializer for Postgres. Creates `uuid-ossp` extension, `patients` table (UUID primary key), `medical_records` table (serial id, patient_id FK, t_score, diagnosis, notes, device_serial, raw_signal_data INTEGER[]), and indexes on `nik` and `patient_id`.
- `backend/sql/migrations/` — existing migration files (e.g., `0001_create_edora_tables.up.sql.sql`, `0002_create_patients_devices_readings.sql`). Review and reconcile with `init.sql` before applying to production DB.

Frontend
--------
- `frontend/` — React + Vite app.
  - `frontend/src/main.jsx`, `App.jsx` — entry and main components.
  - `frontend/src/App.css`, `index.css` — styles.
  - `frontend/package.json` — frontend dependencies and scripts.

Flutter (mobile)
----------------
- `flutter/` — Flutter app (if used).
  - `lib/main.dart` — Flutter entry.
  - `lib/models/reading_model.dart` — mobile model mirroring backend reading.
  - `lib/services/api_service.dart` — client to backend API.

Load testing
------------
- `k6/load-test.js` — k6 script for load testing endpoints.

Quick commands (developer)
-------------------------
- Rebuild and run with Docker Compose (project root):

```bash
docker compose down && docker compose build --no-cache && docker compose up -d
```

- Build backend locally (no Docker):

```bash
cd backend
go build ./...
```

- Run backend tests:

```bash
cd backend
go test ./...
```

- Apply `init.sql` to Postgres (example using psql):

```bash
# from host where psql can reach the DB
psql "postgres://user:pass@localhost:5432/appdb?sslmode=disable" -f backend/sql/init.sql
```

Notes & recommendations
-----------------------
- Consistency: handlers now reference `models.MedicalRecord` (package-qualified). If you add more models move them into `internal/models` and import consistently.
- Database: `init.sql` creates both `patients` (UUID primary key) and `medical_records`. If you use existing migrations, merge carefully to avoid duplicate table definitions.
- Indexes: `idx_patients_nik` and `idx_medical_records_patient_id` added for faster lookups; add more indexes later based on query patterns (e.g., `scan_date`).

Next actions I can take
----------------------
1. Create a proper SQL migration file compatible with your migration tooling (flyway/liquibase/migrate).
2. Add `scripts/integration_test.sh` to automatically run curl checks and assert HTTP statuses.
3. Run the DB init against a running Postgres container and verify tables exist.

If you want me to proceed with any of the next actions, tell me which one.

Generated on: 2025-12-18
