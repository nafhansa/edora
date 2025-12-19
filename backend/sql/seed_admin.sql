-- seed_admin.sql
-- Create users table and seed admin user (dev only)

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  username TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  role TEXT DEFAULT 'admin',
  created_at TIMESTAMPTZ DEFAULT now()
);

INSERT INTO users (username, password, role)
VALUES ('adminedora', crypt('adminedora', gen_salt('bf')), 'admin')
ON CONFLICT (username) DO NOTHING;
