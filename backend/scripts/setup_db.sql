-- ============================================
-- BAYARIN Database Setup Script
-- ============================================

-- Drop existing database if exists (HATI-HATI: akan hapus semua data)
DROP DATABASE IF EXISTS bayarin_db;

-- Drop existing user if exists
DROP USER IF EXISTS bayarin_user;

-- Create user with password
CREATE USER bayarin_user WITH PASSWORD 'bayarin_pass_2024';

-- Create database
CREATE DATABASE bayarin_db OWNER bayarin_user;

-- Grant all privileges
GRANT ALL PRIVILEGES ON DATABASE bayarin_db TO bayarin_user;

-- Connect to database
\c bayarin_db

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO bayarin_user;

-- Grant default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT ALL ON TABLES TO bayarin_user;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT ALL ON SEQUENCES TO bayarin_user;