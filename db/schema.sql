
-- =====================================================
-- Database: PostgreSQL 16+
-- File: db/schema.sql
-- =====================================================

-- =====================================================
-- RESET SCHEMA (DEVELOPMENT ONLY)
-- PERINGATAN: Seluruh data tabel di bawah akan dihapus.
-- =====================================================

BEGIN;

DROP TABLE IF EXISTS attendance_records CASCADE;
DROP TABLE IF EXISTS auth_sessions CASCADE;
DROP TABLE IF EXISTS employees CASCADE;
DROP TABLE IF EXISTS positions CASCADE;
DROP TABLE IF EXISTS departments CASCADE;

-- Hapus function dan trigger versi sebelumnya.
DROP FUNCTION IF EXISTS set_updated_at() CASCADE;

-- =====================================================
-- LANJUTKAN DENGAN CREATE TABLE
-- =====================================================

-- =====================================================
-- 1. DEPARTMENTS
-- Master data departemen
-- =====================================================
CREATE TABLE departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(30) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 2. POSITIONS
-- Master data jabatan
-- =====================================================
CREATE TABLE positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(30) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 3. EMPLOYEES
-- Data registrasi, profil, dan autentikasi karyawan
-- =====================================================
CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_no VARCHAR(30) NOT NULL UNIQUE,

    -- Personal information
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE,
    phone VARCHAR(30),
    date_of_birth DATE,
    gender VARCHAR(20)
        CHECK (gender IN (
            'male', 'female'
        )),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(15),

    -- Authentication and access control
    password_hash VARCHAR(255) NOT NULL,
    account_role VARCHAR(20) NOT NULL DEFAULT 'employee'
        CHECK (account_role IN ('employee', 'admin')),
    account_status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (account_status IN ('active', 'locked', 'disabled')),
    failed_login_attempts INTEGER NOT NULL DEFAULT 0
        CHECK (failed_login_attempts >= 0),
    locked_until TIMESTAMPTZ,

    -- Employment information
department_code VARCHAR(20)
    REFERENCES departments(code) ON DELETE RESTRICT,
position_code VARCHAR(20)
    REFERENCES positions(code) ON DELETE RESTRICT,
    employment_status VARCHAR(20) NOT NULL DEFAULT 'active'
        CHECK (employment_status IN (
            'active', 'inactive', 'suspended', 'resigned'
        )),
    hire_date DATE,

    -- Profile photo metadata
    photo_url TEXT,
    photo_mime_type VARCHAR(100),
    photo_size_bytes BIGINT
        CHECK (
            photo_size_bytes IS NULL
            OR photo_size_bytes BETWEEN 0 AND 5242880
        ),

    -- Audit and soft delete
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- Index untuk pencarian karyawan berdasarkan departemen.
CREATE INDEX idx_employees_department
    ON employees(department_code);

-- Index untuk pencarian karyawan berdasarkan jabatan.
CREATE INDEX idx_employees_position
    ON employees(position_code);

-- Index untuk filter status kerja dan status penghapusan.
CREATE INDEX idx_employees_status_deleted
    ON employees(employment_status, deleted_at);

-- Partial index untuk karyawan yang belum dihapus.
CREATE INDEX idx_employees_active
    ON employees(id)
    WHERE deleted_at IS NULL;


-- =====================================================
-- 4. AUTH_SESSIONS
-- Sesi login dan refresh token
-- =====================================================
CREATE TABLE auth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL
        REFERENCES employees(id) ON DELETE CASCADE,

    -- Hash refresh token, bukan token mentah.
    refresh_token_hash VARCHAR(64) NOT NULL UNIQUE,

    device_id VARCHAR(150),
    ip_address INET,
    user_agent TEXT,

    expires_at TIMESTAMPTZ NOT NULL,
    last_used_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk mencari sesi berdasarkan karyawan.
CREATE INDEX idx_auth_sessions_employee
    ON auth_sessions(employee_id);

-- Partial index untuk sesi yang belum dicabut.
CREATE INDEX idx_auth_sessions_active
    ON auth_sessions(employee_id, expires_at)
    WHERE revoked_at IS NULL;


-- =====================================================
-- 5. ATTENDANCE_RECORDS
-- Transaksi absensi dengan geotagging
-- =====================================================
CREATE TABLE attendance_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL
        REFERENCES employees(id) ON DELETE RESTRICT,

    attendance_type VARCHAR(10) NOT NULL
        CHECK (attendance_type IN ('check_in', 'check_out')),

    -- Waktu transaksi dicatat oleh server/database.
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Koordinat GPS dari perangkat karyawan.
    latitude NUMERIC(10, 7) NOT NULL
        CHECK (latitude BETWEEN -90 AND 90),
    longitude NUMERIC(10, 7) NOT NULL
        CHECK (longitude BETWEEN -180 AND 180),
    accuracy_m NUMERIC(8, 2)
        CHECK (accuracy_m IS NULL OR accuracy_m >= 0),
    location_captured_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk riwayat absensi karyawan, terbaru dahulu.
CREATE INDEX idx_attendance_employee_time
    ON attendance_records(employee_id, recorded_at DESC);

-- Index untuk laporan absensi berdasarkan waktu.
CREATE INDEX idx_attendance_time
    ON attendance_records(recorded_at DESC);

-- Index untuk pencarian berdasarkan karyawan, jenis, dan waktu.
CREATE INDEX idx_attendance_employee_type_time
    ON attendance_records(
        employee_id,
        attendance_type,
        recorded_at DESC
    );


-- =====================================================
-- 6. TRIGGER: AUTO UPDATE updated_at
-- Memperbarui updated_at setiap kali karyawan diubah.
-- =====================================================
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_employees_updated_at
BEFORE UPDATE ON employees
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

COMMIT;
