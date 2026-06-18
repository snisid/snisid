-- ============================================================
-- SNI-SIDE: Firearms Intelligence Database
-- PostgreSQL 16
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_firearms;
SET search_path TO snisid_firearms;

-- ============ FIREARMS ============
CREATE TABLE firearms (
    firearm_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    serial_number VARCHAR(100) UNIQUE NOT NULL,
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    caliber VARCHAR(50) NOT NULL,
    type VARCHAR(50) CHECK (type IN (
        'PISTOL','REVOLVER','RIFLE','SHOTGUN','SUBMACHINE_GUN','ASSAULT_RIFLE',
        'MACHINE_GUN','SNIPER_RIFLE','GRENADE_LAUNCHER','OTHER'
    )),
    action_type VARCHAR(30) CHECK (action_type IN ('SEMI_AUTO','FULL_AUTO','BOLT_ACTION','PUMP','BREAK','REVOLVER','LEVER','OTHER')),
    barrel_length_cm DECIMAL(5,2),
    total_length_cm DECIMAL(5,2),
    weight_kg DECIMAL(5,2),
    magazine_capacity INT,
    country_of_origin VARCHAR(100),
    year_of_manufacture INT,
    manufacturer VARCHAR(255),
    import_mark VARCHAR(100),
    owner_niu VARCHAR(10),
    owner_name VARCHAR(255),
    registration_number VARCHAR(100) UNIQUE,
    registration_date DATE,
    registration_expiry DATE,
    license_number VARCHAR(100),
    license_holder VARCHAR(255),
    status VARCHAR(30) CHECK (status IN (
        'REGISTERED','STOLEN','RECOVERED','EVIDENCE','DESTROYED','EXPORTED','MISSING','DECOMMISSIONED'
    )),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_owner FOREIGN KEY (owner_niu) REFERENCES snisid_identity.citizens(niu)
);

CREATE INDEX idx_firearm_serial ON firearms(serial_number);
CREATE INDEX idx_firearm_make ON firearms(make);
CREATE INDEX idx_firearm_caliber ON firearms(caliber);
CREATE INDEX idx_firearm_type ON firearms(type);
CREATE INDEX idx_firearm_owner ON firearms(owner_niu);
CREATE INDEX idx_firearm_status ON firearms(status);

-- ============ BALLISTIC EVIDENCE ============
CREATE TABLE ballistic_evidence (
    ballistic_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firearm_id UUID REFERENCES firearms(firearm_id),
    case_id UUID,
    evidence_number VARCHAR(100) UNIQUE NOT NULL,
    evidence_type VARCHAR(50) CHECK (evidence_type IN (
        'FIRED_CASING','BULLET','FRAGMENT','FIREARM','GUNSHOT_RESIDUE','IMPRESSION'
    )),
    caliber VARCHAR(50),
    brand_headstamp VARCHAR(50),
    lands_grooves INT,
    twist_direction VARCHAR(10) CHECK (twist_direction IN ('LEFT','RIGHT')),
    rifling_impressions VARCHAR(500),
    firing_pin_impression VARCHAR(500),
    ejector_mark VARCHAR(500),
    extractor_mark VARCHAR(500),
    breach_face_mark VARCHAR(500),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_ballistic_case ON ballistic_evidence(case_id);
CREATE INDEX idx_ballistic_firearm ON ballistic_evidence(firearm_id);
CREATE INDEX idx_ballistic_caliber ON ballistic_evidence(caliber);

-- ============ BALLISTIC MATCHES ============
CREATE TABLE ballistic_matches (
    match_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    evidence_id_1 UUID NOT NULL REFERENCES ballistic_evidence(ballistic_id),
    evidence_id_2 UUID NOT NULL REFERENCES ballistic_evidence(ballistic_id),
    match_score DECIMAL(8,4) NOT NULL,
    correlation_method VARCHAR(50),
    examiner_id VARCHAR(100),
    verification_status VARCHAR(20) CHECK (verification_status IN ('PENDING','CONFIRMED','REJECTED')),
    case_links TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_ballistic_match UNIQUE (evidence_id_1, evidence_id_2)
);

-- ============ STOLEN FIREARMS ============
CREATE TABLE stolen_firearms (
    theft_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firearm_id UUID NOT NULL REFERENCES firearms(firearm_id),
    theft_date TIMESTAMPTZ NOT NULL,
    theft_location VARCHAR(500),
    theft_method VARCHAR(100),
    reported_date TIMESTAMPTZ DEFAULT NOW(),
    reported_by VARCHAR(255),
    reporting_agency VARCHAR(100),
    recovery_date TIMESTAMPTZ,
    recovery_location VARCHAR(500),
    recovery_circumstances TEXT,
    status VARCHAR(20) CHECK (status IN ('STOLEN','RECOVERED','TRACED'))
);

ALTER TABLE firearms ENABLE ROW LEVEL SECURITY;
CREATE POLICY firearms_pnh_select ON firearms FOR SELECT USING (
    current_setting('snisid.agency') IN ('PNH','DCPJ','SNISID_ADMIN')
);
