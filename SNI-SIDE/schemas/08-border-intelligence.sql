-- ============================================================
-- SNI-SIDE: Border Intelligence Database
-- PostgreSQL 16
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_border;
SET search_path TO snisid_border;

-- ============ BORDER CROSSINGS ============
CREATE TABLE border_crossings (
    crossing_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    niu VARCHAR(10),
    passport_number VARCHAR(100),
    travel_document_type VARCHAR(50) CHECK (travel_document_type IN (
        'PASSPORT','CIN','VISA','REFUGEE_TRAVEL_DOC','LAISSEZ_PASSER','DIPLOMATIC'
    )),
    travel_document_code VARCHAR(100),
    full_name VARCHAR(255) NOT NULL,
    nationality VARCHAR(100),
    date_of_birth DATE,
    gender VARCHAR(1),
    crossing_direction VARCHAR(10) CHECK (crossing_direction IN ('ENTRY','EXIT')),
    border_point VARCHAR(255) NOT NULL,
    border_point_code VARCHAR(50),
    crossing_date TIMESTAMPTZ NOT NULL,
    crossing_method VARCHAR(50) CHECK (crossing_method IN (
        'AIR','LAND','SEA','RAIL','PEDESTRIAN'
    )),
    flight_number VARCHAR(50),
    vessel_name VARCHAR(255),
    vehicle_plate VARCHAR(50),
    vehicle_vin VARCHAR(17),
    port_of_entry VARCHAR(255),
    port_of_departure VARCHAR(255),
    purpose_of_travel VARCHAR(100),
    duration_of_stay_days INT,
    visa_type VARCHAR(50),
    visa_number VARCHAR(100),
    visa_issue_date DATE,
    visa_expiry_date DATE,
    biometric_verified BOOLEAN DEFAULT FALSE,
    biometric_match_score DECIMAL(5,2),
    risk_score DECIMAL(5,2),
    alert_triggered BOOLEAN DEFAULT FALSE,
    officer_id VARCHAR(100),
    officer_agency VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_citizen FOREIGN KEY (niu) REFERENCES snisid_identity.citizens(niu)
) PARTITION BY RANGE (crossing_date);

CREATE INDEX idx_border_cross_niu ON border_crossings(niu);
CREATE INDEX idx_border_cross_passport ON border_crossings(passport_number);
CREATE INDEX idx_border_cross_name ON border_crossings USING gin(to_tsvector('french', full_name));
CREATE INDEX idx_border_cross_date ON border_crossings(crossing_date DESC);
CREATE INDEX idx_border_cross_direction ON border_crossings(crossing_direction);
CREATE INDEX idx_border_cross_point ON border_crossings(border_point);
CREATE INDEX idx_border_cross_method ON border_crossings(crossing_method);
CREATE INDEX idx_border_cross_nationality ON border_crossings(nationality);
CREATE INDEX idx_border_cross_risk ON border_crossings(risk_score DESC);
CREATE INDEX idx_border_cross_alert ON border_crossings(alert_triggered);
CREATE INDEX idx_border_cross_visa ON border_crossings(visa_number);
CREATE INDEX idx_border_cross_vehicle ON border_crossings(vehicle_plate);

-- ============ VISAS ============
CREATE TABLE visas (
    visa_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    visa_number VARCHAR(100) UNIQUE NOT NULL,
    niu VARCHAR(10),
    passport_number VARCHAR(100) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    nationality VARCHAR(100) NOT NULL,
    visa_type VARCHAR(50) CHECK (visa_type IN (
        'DIPLOMATIC','OFFICIAL','TOURIST','BUSINESS','STUDENT','WORK','TRANSIT','HUMANITARIAN','FAMILY'
    )),
    visa_category VARCHAR(50),
    issuing_post VARCHAR(255) NOT NULL,
    issuing_country VARCHAR(100),
    issue_date DATE NOT NULL,
    expiry_date DATE NOT NULL,
    max_stay_days INT,
    entries_allowed INT,
    entries_used INT DEFAULT 0,
    conditions TEXT,
    status VARCHAR(20) CHECK (status IN ('VALID','EXPIRED','REVOKED','CANCELLED','USED')),
    biometric_enrolled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_visa_number ON visas(visa_number);
CREATE INDEX idx_visa_niu ON visas(niu);
CREATE INDEX idx_visa_passport ON visas(passport_number);
CREATE INDEX idx_visa_type ON visas(visa_type);
CREATE INDEX idx_visa_status ON visas(status);
CREATE INDEX idx_visa_expiry ON visas(expiry_date);

-- ============ DEPORTATIONS ============
CREATE TABLE deportations (
    deportation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    niu VARCHAR(10),
    passport_number VARCHAR(100),
    full_name VARCHAR(255) NOT NULL,
    nationality VARCHAR(100) NOT NULL,
    deportation_date TIMESTAMPTZ NOT NULL,
    deportation_reason VARCHAR(500),
    deportation_type VARCHAR(30) CHECK (deportation_type IN ('ADMINISTRATIVE','CRIMINAL','SECURITY','OVERSTAY')),
    origin_country VARCHAR(100),
    destination_country VARCHAR(100),
    flight_number VARCHAR(50),
    escort_agency VARCHAR(100),
    escort_officers TEXT[],
    departure_point VARCHAR(255),
    arrival_point VARCHAR(255),
    biometric_confirmed BOOLEAN DEFAULT FALSE,
    case_reference VARCHAR(100),
    status VARCHAR(20) CHECK (status IN ('PENDING','IN_TRANSIT','COMPLETED','CANCELLED','ESCAPED')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_deportation_niu ON deportations(niu);
CREATE INDEX idx_deportation_date ON deportations(deportation_date DESC);
CREATE INDEX idx_deportation_nationality ON deportations(nationality);

-- ============ BORDER WATCHLIST MATCHES ============
CREATE TABLE border_watchlist_matches (
    match_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crossing_id UUID NOT NULL REFERENCES border_crossings(crossing_id),
    watchlist_type VARCHAR(30) CHECK (watchlist_type IN (
        'INTERPOL_RED','INTERPOL_BLUE','NATIONAL_WANTED','TERRORIST','SANCTIONS','AMBER_ALERT'
    )),
    matched_person VARCHAR(255),
    matched_niuk VARCHAR(10),
    match_confidence DECIMAL(5,2),
    alert_level VARCHAR(20) CHECK (alert_level IN ('CRITICAL','HIGH','MEDIUM','LOW')),
    action_taken VARCHAR(500),
    status VARCHAR(20) CHECK (status IN ('NEW','ACTIONED','ESCALATED','FALSE_POSITIVE')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE border_crossings ENABLE ROW LEVEL SECURITY;
CREATE POLICY border_immigration_select ON border_crossings FOR SELECT USING (
    current_setting('snisid.agency') IN ('IMMIGRATION','PNH','DCPJ','CUSTOMS','SNISID_ADMIN')
);
