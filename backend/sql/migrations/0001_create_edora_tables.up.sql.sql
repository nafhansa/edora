-- Hapus tabel lama jika ada biar bersih
DROP TABLE IF EXISTS readings;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS patients;

-- Enable UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Tabel Pasien
CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nik VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    gender CHAR(1) NOT NULL CHECK (gender IN ('M', 'F')),
    birth_date DATE NOT NULL,
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Tabel Alat (Device)
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    serial_number VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100),
    status VARCHAR(20) DEFAULT 'offline', 
    last_seen TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 3. Tabel Pembacaan (Readings - INTI DATA)
CREATE TABLE readings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
    patient_id UUID REFERENCES patients(id) ON DELETE CASCADE,
    doctor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Hasil Medis
    bmd_result FLOAT NOT NULL,
    t_score FLOAT NOT NULL,
    classification VARCHAR(50) NOT NULL, 
    
    -- Data Sinyal (Array disimpan sebagai JSONB)
    raw_signal_data JSONB NOT NULL, 
    
    -- Lokasi GPS
    latitude FLOAT,
    longitude FLOAT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);