BEGIN;
CREATE TYPE expl_type AS ENUM ('IED','GRENADE','RPG','MORTAR','LANDMINE','DYNAMITE','BLASTING_CAP','AMMUNITION_BULK','MILITARY_ORDNANCE','UNKNOWN');
CREATE TYPE expl_status AS ENUM ('RECOVERED','DESTROYED','DETONATED','STORED_EVIDENCE','TRANSFERRED');
CREATE TABLE expl_incidents (
    incident_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), national_expl_id VARCHAR(25) UNIQUE NOT NULL,
    incident_type VARCHAR(30) NOT NULL, explosive_type expl_type NOT NULL, status expl_status NOT NULL DEFAULT 'RECOVERED',
    quantity INTEGER DEFAULT 1, weight_kg DECIMAL(10,3), manufacturer VARCHAR(100), lot_number VARCHAR(50),
    manufacture_country CHAR(3), estimated_date DATE, incident_date TIMESTAMPTZ NOT NULL,
    location_desc VARCHAR(300), dept_code CHAR(2), commune VARCHAR(100), lat DECIMAL(10,7), lng DECIMAL(10,7),
    responding_unit VARCHAR(50), eod_officer UUID, casualties SMALLINT DEFAULT 0, gang_id UUID,
    from_person_id UUID, case_reference VARCHAR(100), dna_sample_taken BOOLEAN DEFAULT FALSE,
    bio_sample_ref VARCHAR(100), photo_refs TEXT[] DEFAULT '{}', interpol_exploint_ref VARCHAR(50),
    notes TEXT, created_by UUID NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE expl_legal_stocks (
    stock_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), holder_entity VARCHAR(200) NOT NULL,
    holder_type VARCHAR(30), explosive_type expl_type NOT NULL, quantity_kg DECIMAL(12,3) NOT NULL,
    storage_location TEXT NOT NULL, dept_code CHAR(2), license_ref VARCHAR(50),
    last_audit_date DATE, next_audit_date DATE, is_secured BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_expl_type ON expl_incidents(explosive_type, incident_date DESC);
CREATE INDEX idx_expl_dept ON expl_incidents(dept_code);
CREATE INDEX idx_expl_gang ON expl_incidents(gang_id) WHERE gang_id IS NOT NULL;
COMMIT;
