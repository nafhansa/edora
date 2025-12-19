-- init.sql
-- Initialize database schema for Edora (IoT Optimized)

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: patients
CREATE TABLE IF NOT EXISTS patients (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  nik VARCHAR(16) UNIQUE NOT NULL, -- Sesuai instruksi LATEST_PROGRESS
  name TEXT NOT NULL,
  gender VARCHAR(16),
  birth_date DATE, -- Handled with time.Parse in Go
  address TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_patients_nik ON patients (nik);

-- Table: medical_records (Sinkron dengan POST /sync/reading)
CREATE TABLE IF NOT EXISTS medical_records (
  id SERIAL PRIMARY KEY,
  patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
  
  -- IoT & Device Data
  device_serial TEXT,
  doctor_id UUID, -- Untuk tracking siapa dokter yang memeriksa
  
  -- Clinical Results
  bmd_result REAL, -- Bone Mineral Density
  t_score REAL NOT NULL,
  diagnosis TEXT NOT NULL, -- Server computes this from t_score
  
  -- Sensor Data
  raw_signal_data INTEGER[], -- Array sinyal dari alat
  
  -- Geolocation
  lat DECIMAL(9,6), -- Lokasi alat saat scan
  long DECIMAL(9,6),
  
  -- Metadata
  notes TEXT,
  scan_date TIMESTAMP WITH TIME ZONE DEFAULT now(),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_medical_records_patient_id ON medical_records (patient_id);
CREATE INDEX IF NOT EXISTS idx_medical_records_device_serial ON medical_records (device_serial);