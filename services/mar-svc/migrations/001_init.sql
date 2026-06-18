BEGIN;

CREATE TYPE mar_vessel_type AS ENUM ('CARGO_SHIP','TANKER','FISHING_BOAT','GO_FAST','SAILBOAT','YACHT','FERRY','PATROL_BOAT','WOODEN_BOAT','CANOE','UNKNOWN');
CREATE TYPE mar_vessel_status AS ENUM ('REGISTERED','STOLEN','SUSPECTED','DETAINED','SUNK','DESTROYED','MISSING','INTERPOL_ALERT');
CREATE TYPE mar_incident_type AS ENUM ('DRUG_SEIZURE','ARMS_SEIZURE','MIGRANT_INTERDICTION','SMUGGLING','SUSPICIOUS_ACTIVITY','DISTRESS','PIRACY','ILLEGAL_FISHING','HUMAN_TRAFFICKING');

CREATE TABLE mar_vessels (
    vessel_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_mar_id VARCHAR(25) UNIQUE NOT NULL,
    vessel_name VARCHAR(150),
    imo_number VARCHAR(20),
    mmsi VARCHAR(15),
    call_sign VARCHAR(15),
    vessel_type mar_vessel_type NOT NULL,
    flag_country CHAR(3),
    hull_color VARCHAR(50),
    length_m DECIMAL(8,2),
    tonnage_gt INTEGER,
    engine_count SMALLINT,
    horsepower INTEGER,
    owner_name VARCHAR(200),
    owner_snisid_id UUID,
    registration_number VARCHAR(50),
    registration_port VARCHAR(100),
    status mar_vessel_status NOT NULL DEFAULT 'REGISTERED',
    gang_id UUID,
    interpol_svd_ref VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mar_ais_sightings (
    sighting_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vessel_id UUID REFERENCES mar_vessels(vessel_id),
    mmsi VARCHAR(15),
    vessel_name VARCHAR(150),
    sighting_timestamp TIMESTAMPTZ NOT NULL,
    lat DECIMAL(10,7) NOT NULL,
    lng DECIMAL(10,7) NOT NULL,
    speed_knots DECIMAL(5,2),
    heading_degrees SMALLINT,
    destination VARCHAR(100),
    source_type VARCHAR(30),
    zone_code VARCHAR(20),
    alert_triggered BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mar_incidents (
    incident_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vessel_id UUID REFERENCES mar_vessels(vessel_id),
    incident_type mar_incident_type NOT NULL,
    incident_date TIMESTAMPTZ NOT NULL,
    lat DECIMAL(10,7),
    lng DECIMAL(10,7),
    zone_desc VARCHAR(100),
    responding_unit VARCHAR(50),
    outcome TEXT,
    persons_involved INTEGER DEFAULT 0,
    snisid_person_ids UUID[] DEFAULT '{}',
    drug_types TEXT[] DEFAULT '{}',
    drug_weight_kg DECIMAL(12,3),
    weapons_found BOOLEAN DEFAULT FALSE,
    weapons_count INTEGER DEFAULT 0,
    migrants_count INTEGER DEFAULT 0,
    biar_refs UUID[] DEFAULT '{}',
    case_reference VARCHAR(100),
    photo_refs TEXT[] DEFAULT '{}',
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mar_watch_vessels (
    watch_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vessel_id UUID REFERENCES mar_vessels(vessel_id),
    mmsi VARCHAR(15),
    vessel_name VARCHAR(150),
    watch_reason TEXT NOT NULL,
    alert_level VARCHAR(20) DEFAULT 'CAUTION',
    requesting_unit VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    expiry_date TIMESTAMPTZ,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_mar_ais_timestamp ON mar_ais_sightings(sighting_timestamp DESC);
CREATE INDEX idx_mar_ais_mmsi ON mar_ais_sightings(mmsi);
CREATE INDEX idx_mar_ais_coords ON mar_ais_sightings(lat, lng);
CREATE INDEX idx_mar_incidents_date ON mar_incidents(incident_date DESC);
CREATE INDEX idx_mar_incidents_type ON mar_incidents(incident_type);
CREATE INDEX idx_mar_watch_active ON mar_watch_vessels(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_mar_vessels_status ON mar_vessels(status);

COMMIT;
