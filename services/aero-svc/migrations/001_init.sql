BEGIN;
CREATE TYPE aero_aircraft_type AS ENUM ('COMMERCIAL_JET','TURBOPROP','PISTON_SINGLE','PISTON_TWIN','HELICOPTER','ULTRALIGHT','DRONE_LARGE','UNKNOWN');
CREATE TYPE aero_strip_status AS ENUM ('ACTIVE','INACTIVE','DESTROYED','LEGALIZED','UNDER_SURVEILLANCE');
CREATE TABLE aero_aircraft_registry (
    aircraft_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), registration_mark VARCHAR(20),
    icao_hex_code VARCHAR(10), aircraft_type aero_aircraft_type NOT NULL,
    make VARCHAR(100), model VARCHAR(100), manufacture_year SMALLINT, flag_country CHAR(3),
    owner_name VARCHAR(200), owner_snisid_id UUID, operator_name VARCHAR(200),
    is_registered BOOLEAN DEFAULT FALSE, is_suspected BOOLEAN DEFAULT FALSE,
    is_stolen BOOLEAN DEFAULT FALSE, gang_id UUID, drug_trafficking BOOLEAN DEFAULT FALSE,
    interpol_ref VARCHAR(50), faa_registry_ref VARCHAR(50), notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE aero_clandestine_strips (
    strip_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), strip_name VARCHAR(150),
    dept_code CHAR(2) NOT NULL, commune VARCHAR(100), lat DECIMAL(10,7) NOT NULL,
    lng DECIMAL(10,7) NOT NULL, length_m INTEGER, surface_type VARCHAR(30),
    status aero_strip_status NOT NULL DEFAULT 'ACTIVE', capable_aircraft TEXT[] DEFAULT '{}',
    gang_id UUID, first_detected DATE, last_activity_date DATE, source_intel TEXT,
    satellite_image_ref VARCHAR(500), created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE aero_suspicious_flights (
    flight_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), aircraft_id UUID REFERENCES aero_aircraft_registry(aircraft_id),
    registration_mark VARCHAR(20), flight_date TIMESTAMPTZ NOT NULL,
    origin_airport VARCHAR(10), destination_airport VARCHAR(10),
    origin_country CHAR(3), destination_country CHAR(3) DEFAULT 'HTI',
    landing_strip_id UUID REFERENCES aero_clandestine_strips(strip_id),
    landing_location VARCHAR(300), flight_type VARCHAR(30), cargo_suspected TEXT,
    source_radar VARCHAR(50), source_informant BOOLEAN DEFAULT FALSE,
    case_reference VARCHAR(100), created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_aero_registry_mark ON aero_aircraft_registry(registration_mark);
CREATE INDEX idx_aero_registry_gang ON aero_aircraft_registry(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_aero_strips_dept ON aero_clandestine_strips(dept_code) WHERE status = 'ACTIVE';
CREATE INDEX idx_aero_strips_coords ON aero_clandestine_strips(lat, lng);
CREATE INDEX idx_aero_flights_date ON aero_suspicious_flights(flight_date DESC);
COMMIT;
