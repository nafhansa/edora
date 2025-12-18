STEP 1. Project Context & Objective

Project Name: EDORA (Early Detection of Risk Assessment - Osteoporosis) Goal: Create a high-performance IoT ecosystem where an ESP32 device sends bone impedance data to a Flutter App, which syncs to a Go Backend, and is visualized on a Next.js Dashboard.

Critical Requirement:

    No Boilerplate: Remove all "Product" or generic template code. Replace with specific domain logic (Patients, Readings, Devices).

    Latency Control: The Mobile App must buffer data. The Backend must handle high-throughput writes.

    Data Integrity: Ensure raw signal data (arrays of impedance) is stored accurately for medical analysis.

STEP 2. Domain & Database Architecture (PostgreSQL)

Action: Refactor internal/models and sql/migrations to match this schema strictly.
A. Conceptual ERD

    Users: Doctors/Admins who access the system.

    Devices: The physical ESP32 hardware units (tracked for inventory/status).

    Patients: The subjects being scanned.

    Readings: The core transaction. One scan event containing signal data and AI results.

B. Schema Specification
1. Table: users (Existing - Keep but Refine)

    id (UUID, PK)

    email (VARCHAR, Unique)

    password_hash (VARCHAR)

    role (ENUM: 'doctor', 'admin')

    full_name (VARCHAR)

2. Table: patients (NEW)

    id (UUID, PK)

    nik (VARCHAR, Unique) - National ID

    name (VARCHAR)

    gender (CHAR(1)) - 'M'/'F'

    birth_date (DATE)

    address (TEXT) - For geolocation mapping context

3. Table: devices (NEW)

    id (UUID, PK)

    serial_number (VARCHAR, Unique) - Matches ESP32 Mac/ID

    name (VARCHAR) - e.g., "Unit Lab A"

    status (VARCHAR) - 'online', 'offline', 'maintenance'

    last_seen (TIMESTAMP)

4. Table: readings (NEW - The Core Data)

    id (UUID, PK)

    device_id (UUID, FK -> devices)

    patient_id (UUID, FK -> patients)

    doctor_id (UUID, FK -> users)

    raw_signal_data (JSONB) - CRITICAL: Stores array [{freq: 50, z: 120}, ...]

    bmd_result (FLOAT) - Bone Mineral Density value

    t_score (FLOAT) - Classification score

    classification (VARCHAR) - 'Normal', 'Osteopenia', 'Osteoporosis'

    latitude (FLOAT) - GPS Lat from Mobile App

    longitude (FLOAT) - GPS Long from Mobile App

    created_at (TIMESTAMP)

STEP 3. Backend Implementation (Go + Fiber)

Directory: backend/
Cleanup Tasks

    DELETE: internal/models/product.go

    DELETE: internal/handler/product_handler.go

    DELETE: internal/repository/product_repo.go

    DELETE: internal/service/product_service.go

Implementation Tasks
A. Models (internal/models/)

Create struct files for Patient, Device, and Reading.

    Note: Reading.RawSignalData must be type json.RawMessage or []byte to efficiently handle JSONB without deep parsing overhead during write.

B. API Endpoints

1. Mobile Sync Endpoint (POST /api/v1/sync/reading)

    Purpose: Receives scan results from Flutter.

    Payload:
    JSON

    {
      "device_serial": "ESP32-001",
      "patient_id": "uuid...",
      "doctor_id": "uuid...",
      "timestamp": "2023-10-10T10:00:00Z",
      "location": {"lat": -6.91, "lng": 107.61},
      "readings": [
        {"freq": 5000, "impedance": 120.5, "phase": 12.0},
        {"freq": 50000, "impedance": 90.2, "phase": 8.5}
      ],
      "analysis": {
        "bmd": 0.85,
        "t_score": -1.2,
        "class": "Osteopenia"
      }
    }

    Logic:

        Validate Device Serial.

        Insert into readings table.

        Update devices table (last_seen = now).

        Return 201 OK.

2. Web Dashboard Endpoint (GET /api/v1/dashboard/stats)

    Purpose: Aggregated data for the web chart.

    Response:

        Total Scans Today.

        Breakdown: Normal vs Osteoporosis count.

        Recent 5 Readings (for list view).

3. Patient History (GET /api/v1/patients/:id/history)

    Purpose: Show scan history curve.

4. Frontend Implementation (Next.js)

Directory: frontend/
UI/UX Design System (Medical Theme)

    Colors: White background, Slate-500 text, Blue-600 (Primary), Red-500 (Alert/Osteoporosis), Green-500 (Normal).

    Framework: Tailwind CSS + ShadcnUI (recommended) or native Tailwind.

Key Pages
A. Dashboard (/dashboard)

    Layout: Sidebar (Menu), Main Content.

    Top Cards: "Total Patients", "Devices Active", "Risk Alert (Osteoporosis Detected)".

    Main Chart: Use Recharts. A Line Chart showing the "Impedance Curve" of the most recent scan.

    Recent Scans Table: Date, Patient Name, T-Score (Color coded), Status.

B. Device Monitor (/devices)

    Grid view of devices.

    Status indicator (Green dot = Online, Gray = Offline).

Data Fetching Strategy

    Use SWR or React Query.

    Auto-Refresh: Set dashboard to poll GET /api/v1/dashboard/stats every 10 seconds.

5. Mobile App Implementation (Flutter)

Directory: flutter/
Architecture: MVVM + Repository Pattern
Core Features
A. Scanning Logic (The "Anti-Lag" Strategy)

    State 1: Connection. Connect to ESP32 via BLE or connect to local Hotspot.

    State 2: Acquisition. App receives stream of bytes. DO NOT render chart on every single byte. Buffer data into a List.

    State 3: Visualization. Update UI chart only every 100ms (throttle) or after scan finishes.

    State 4: Upload.

        Check Internet.

        If Online: POST to /api/v1/sync/reading.

        If Offline: Save to SQLite (sqflite). Background service (workmanager) retries upload later.

B. UI Screens

    Home: "Start Scan", "Patient List".

    Scan Page: Real-time wave animation (simplified), "Analyzing..." loader.

    Result Page:

        Big T-Score Gauge.

        Recommendation Text.

        "Save & Upload" button.

6. Execution Instructions for Copilot

Step 1: Backend Refactor

    Analyze backend/go.mod and ensure dependencies (gorm, fiber, pgx).

    Execute the "Cleanup Tasks" (Delete product files).

    Generate sql/migrations/0002_create_edora_tables.sql based on Schema Specification.

    Create Models in internal/models.

    Create ReadingHandler in internal/handler.

Step 2: Frontend Dashboard

    Create app/dashboard/page.tsx.

    Implement Recharts for signal visualization.

    Create a types file types/api.ts matching the Go backend JSON response.

Step 3: Mobile Service

    Create lib/services/api_service.dart.

    Implement the syncReading function with error handling.