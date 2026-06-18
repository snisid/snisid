BEGIN;

CREATE TABLE sipep_facilities (
    facility_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code              VARCHAR(20) UNIQUE NOT NULL,
    name              VARCHAR(200) NOT NULL,
    department        VARCHAR(100) NOT NULL,
    dept_code         CHAR(2) NOT NULL,
    facility_type     facility_type NOT NULL,
    capacity          INTEGER NOT NULL,
    current_occupancy INTEGER DEFAULT 0,
    address           TEXT,
    phone             VARCHAR(30),
    is_active         BOOLEAN DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO sipep_facilities (code, name, department, dept_code, facility_type, capacity) VALUES
    ('PNPP',  'Pénitencier National P-au-P', 'Ouest',      'OU', 'NATIONAL',     3500),
    ('PCCH',  'Prison Civile Cap-Haïtien',   'Nord',       'ND', 'DEPARTMENTAL', 800),
    ('PCGO',  'Prison Civile Gonaïves',      'Artibonite', 'AR', 'DEPARTMENTAL', 400),
    ('PCLC',  'Prison Civile Les Cayes',     'Sud',        'SD', 'DEPARTMENTAL', 300),
    ('CML',   'CERMICOL (Mineurs)',          'Ouest',      'OU', 'SPECIALIZED',  100),
    ('RSEK',  'Établissement femmes (RESEK)','Ouest',      'OU', 'SPECIALIZED',  150);

COMMIT;
