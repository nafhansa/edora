-- Create patients, devices, readings
BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Patients
CREATE TABLE IF NOT EXISTS patients (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  email VARCHAR(255) UNIQUE,
  phone VARCHAR(50),
  dob DATE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_patients_email ON patients(email);

-- Devices
CREATE TABLE IF NOT EXISTS devices (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  patient_id UUID REFERENCES patients(id) ON DELETE CASCADE,
  device_type VARCHAR(100) NOT NULL,
  serial VARCHAR(255) UNIQUE,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_devices_patient ON devices(patient_id);

-- Readings
CREATE TABLE IF NOT EXISTS readings (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  patient_id UUID REFERENCES patients(id) ON DELETE CASCADE,
  device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
  metric VARCHAR(100) NOT NULL,
  value DOUBLE PRECISION,
  unit VARCHAR(50),
  recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_readings_patient ON readings(patient_id);
CREATE INDEX IF NOT EXISTS idx_readings_device ON readings(device_id);
CREATE INDEX IF NOT EXISTS idx_readings_metric_recorded_at ON readings(metric, recorded_at);

COMMIT;
