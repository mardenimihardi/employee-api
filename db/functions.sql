CREATE OR REPLACE FUNCTION fn_register_employee(
    p_employee_no VARCHAR,
    p_full_name VARCHAR,
    p_email VARCHAR,
    p_phone VARCHAR,
    p_date_of_birth DATE,
    p_gender VARCHAR,
    p_address TEXT,
    p_city VARCHAR,
    p_province VARCHAR,
    p_postal_code VARCHAR,
    p_password_hash TEXT,
    p_department_code VARCHAR,
    p_position_code VARCHAR
)
RETURNS UUID
LANGUAGE plpgsql
AS $$
DECLARE
    v_employee_id UUID;
BEGIN
    INSERT INTO employees (
        employee_no,
        full_name,
        email,
        phone,
        date_of_birth,
        gender,
        address,
        city,
        province,
        postal_code,
        password_hash,
        department_code,
        position_code
    )
    VALUES (
        p_employee_no,
        p_full_name,
        p_email,
        NULLIF(p_phone, ''),
        p_date_of_birth,
        NULLIF(p_gender, ''),
        NULLIF(p_address, ''),
        NULLIF(p_city, ''),
        NULLIF(p_province, ''),
        NULLIF(p_postal_code, ''),
        p_password_hash,
        NULLIF(p_department_code, ''),
        NULLIF(p_position_code, '')
    )
    RETURNING id INTO v_employee_id;

    RETURN v_employee_id;
END;
$$;


DROP FUNCTION IF EXISTS fn_login_employee(VARCHAR);

CREATE FUNCTION fn_login_employee(
    p_employee_no VARCHAR
)
RETURNS TABLE (
    id UUID,
    employee_no VARCHAR,
    full_name VARCHAR,
    email VARCHAR,
    password_hash TEXT,
    account_role VARCHAR,
    account_status VARCHAR,
    failed_login_attempts INTEGER,
    locked_until TIMESTAMPTZ
)
LANGUAGE sql
STABLE
AS $$
    SELECT
        e.id,
        e.employee_no,
        e.full_name,
        e.email,
        e.password_hash,
        e.account_role,
        e.account_status,
        e.failed_login_attempts,
        e.locked_until
    FROM employees e
    WHERE e.employee_no = p_employee_no
      AND e.deleted_at IS NULL
    LIMIT 1;
$$;





-- =====================================================
-- ATTENDANCE TRANSACTION
-- =====================================================

CREATE OR REPLACE FUNCTION fn_record_attendance(
    p_employee_id UUID,
    p_attendance_type VARCHAR,
    p_latitude NUMERIC,
    p_longitude NUMERIC,
    p_accuracy_m NUMERIC,
    p_location_captured_at TIMESTAMPTZ
)
RETURNS SETOF attendance_records
LANGUAGE plpgsql
AS $$
BEGIN
    IF p_attendance_type NOT IN ('check_in', 'check_out') THEN
        RAISE EXCEPTION 'Invalid attendance type'
            USING ERRCODE = '22023';
    END IF;

    RETURN QUERY
    INSERT INTO attendance_records (
        employee_id,
        attendance_type,
        recorded_at,
        latitude,
        longitude,
        accuracy_m,
        location_captured_at
    )
    VALUES (
        p_employee_id,
        p_attendance_type,
        NOW(),
        p_latitude,
        p_longitude,
        p_accuracy_m,
        p_location_captured_at
    )
    RETURNING *;
END;
$$;



