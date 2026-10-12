BEGIN;

-- =====================================================
-- SEED: DEPARTMENTS
-- =====================================================
INSERT INTO departments (code, name)
VALUES
    ('IT', 'Information Technology'),
    ('HR', 'Human Resources'),
    ('FIN', 'Finance'),
    ('OPS', 'Operations'),
    ('MKT', 'Marketing'),
    ('SCM', 'Supply Chain Management'),
    ('STORE', 'Store Operations')
ON CONFLICT (code) DO NOTHING;


-- =====================================================
-- SEED: POSITIONS
-- =====================================================
INSERT INTO positions (code, name)
VALUES
    ('IT-DEV', 'IT Developer'),
    ('IT-LEAD', 'IT Development Lead'),
    ('IT-ADMIN', 'IT Administrator'),
    ('HR-STAFF', 'HR Staff'),
    ('FIN-STAFF', 'Finance Staff'),
    ('OPS-STAFF', 'Operations Staff'),
    ('STORE-MGR', 'Store Manager'),
    ('STORE-STAFF', 'Store Staff'),
    ('MGR', 'Manager'),
    ('SUPERVISOR', 'Supervisor')
ON CONFLICT (code) DO NOTHING;

COMMIT;