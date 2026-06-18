CREATE TABLE bie_stolen_vehicles (
    record_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_number   VARCHAR(50) UNIQUE NOT NULL,
    vin             VARCHAR(17) UNIQUE,
    plate_number    VARCHAR(20),
    plate_dept      VARCHAR(50),
    vehicle_make    VARCHAR(100),
    vehicle_model   VARCHAR(100),
    vehicle_year    SMALLINT,
    vehicle_color   VARCHAR(50),
    vehicle_type    VARCHAR(50),
    theft_date      DATE NOT NULL,
    theft_location  VARCHAR(200) NOT NULL,
    theft_department VARCHAR(50),
    owner_niu       VARCHAR(20),
    owner_name      VARCHAR(200),
    owner_phone     VARCHAR(50),
    foves_record_id UUID,
    status          VARCHAR(20) DEFAULT 'STOLEN' CHECK (status IN ('STOLEN','RECOVERED','CANCELLED')),
    recovered_date  DATE,
    recovered_location VARCHAR(200),
    entering_agency VARCHAR(100) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE bie_stolen_firearms (
    record_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_number   VARCHAR(50) UNIQUE NOT NULL,
    serial_number   VARCHAR(100) UNIQUE NOT NULL,
    make            VARCHAR(100),
    model           VARCHAR(100),
    caliber         VARCHAR(50),
    firearm_type    VARCHAR(50),
    barrel_length   DECIMAL(5,2),
    theft_date      DATE NOT NULL,
    theft_location  VARCHAR(200),
    owner_niu       VARCHAR(20),
    status          VARCHAR(20) DEFAULT 'STOLEN' CHECK (status IN ('STOLEN','RECOVERED','CANCELLED')),
    recovered_date  DATE,
    entering_agency VARCHAR(100) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE bie_stolen_documents (
    record_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_number   VARCHAR(50) UNIQUE NOT NULL,
    document_type   VARCHAR(50) NOT NULL CHECK (document_type IN ('PASSPORT','CIN','ACTE_NAISSANCE','PERMIS_CONDUIRE','TITRE_FONCIER','AUTRE')),
    document_number VARCHAR(100),
    issuing_agency  VARCHAR(100),
    issue_date      DATE,
    expiry_date     DATE,
    owner_niu       VARCHAR(20),
    owner_name      VARCHAR(200),
    report_date     DATE NOT NULL,
    report_location VARCHAR(200),
    theft_type      VARCHAR(20) CHECK (theft_type IN ('STOLEN','LOST','FORGED')),
    status          VARCHAR(20) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','RECOVERED','CANCELLED')),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE bie_stolen_vessels (
    record_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_number   VARCHAR(50) UNIQUE NOT NULL,
    vessel_name     VARCHAR(200),
    registration_number VARCHAR(100),
    hull_id_number  VARCHAR(50),
    vessel_type     VARCHAR(50),
    vessel_make     VARCHAR(100),
    vessel_length_m DECIMAL(6,2),
    hull_color      VARCHAR(50),
    home_port       VARCHAR(200),
    theft_location  VARCHAR(200) NOT NULL,
    theft_date      DATE NOT NULL,
    owner_niu       VARCHAR(20),
    owner_name      VARCHAR(200),
    status          VARCHAR(20) DEFAULT 'STOLEN',
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE bio_audit_log (
    log_id          UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    event_type      VARCHAR(100) NOT NULL,
    table_name      VARCHAR(100),
    record_id       UUID,
    officer_niu     VARCHAR(20) NOT NULL,
    agency_code     VARCHAR(50) NOT NULL,
    purpose         VARCHAR(200) NOT NULL,
    case_number     VARCHAR(100),
    ip_hash         VARCHAR(64),
    action          VARCHAR(20) CHECK (action IN ('CREATE','READ','UPDATE','DELETE','SEARCH','HIT')),
    details         JSONB,
    signature       TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE TABLE bio_audit_log_2026_06 PARTITION OF bio_audit_log
    FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');

ALTER TABLE bio_str_profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE per_wanted_persons ENABLE ROW LEVEL SECURITY;
ALTER TABLE per_gang_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE bio_identity_links ENABLE ROW LEVEL SECURITY;

CREATE POLICY bio_identity_links_policy ON bio_identity_links
    USING (current_user = 'snisid_dcpj_director' OR current_user = 'snisid_admin');
