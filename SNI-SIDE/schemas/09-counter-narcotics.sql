-- ============================================================
-- SNI-SIDE: Counter Narcotics Intelligence Database
-- PostgreSQL 16
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_narcotics;
SET search_path TO snisid_narcotics;

-- ============ CARTELS ============
CREATE TABLE cartels (
    cartel_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    alias TEXT[],
    country_of_origin VARCHAR(100),
    regions_operating TEXT[],
    primary_drugs TEXT[],
    estimated_size VARCHAR(100),
    estimated_revenue_usd DECIMAL(20,2),
    leader_name VARCHAR(255),
    hierarchy JSONB DEFAULT '{}',
    known_factions TEXT[],
    rival_cartels TEXT[],
    allied_groups TEXT[],
    corruption_network JSONB DEFAULT '{}',
    violence_level VARCHAR(20) CHECK (violence_level IN ('EXTREME','HIGH','MEDIUM','LOW')),
    status VARCHAR(20) CHECK (status IN ('ACTIVE','DISRUPTED','DISMANTLED','INVESTIGATING')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_cartel_name ON cartels(name);
CREATE INDEX idx_cartel_country ON cartels(country_of_origin);
CREATE INDEX idx_cartel_status ON cartels(status);
CREATE INDEX idx_cartel_drugs ON cartels USING gin(primary_drugs);

-- ============ NARCOTICS ROUTES ============
CREATE TABLE narcotics_routes (
    route_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    route_type VARCHAR(30) CHECK (route_type IN ('AIR','MARITIME','LAND','MIXED','POSTAL')),
    origin_country VARCHAR(100) NOT NULL,
    origin_region VARCHAR(255),
    transit_countries TEXT[],
    destination_country VARCHAR(100) NOT NULL,
    destination_region VARCHAR(255),
    primary_drugs TEXT[],
    estimated_annual_volume_kg DECIMAL(15,2),
    estimated_value_usd DECIMAL(20,2),
    method_of_transport TEXT[],
    known_cartels TEXT[],
    risk_level VARCHAR(20),
    status VARCHAR(20) CHECK (status IN ('ACTIVE','MONITORED','DISRUPTED')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_route_origin ON narcotics_routes(origin_country);
CREATE INDEX idx_route_destination ON narcotics_routes(destination_country);
CREATE INDEX idx_route_risk ON narcotics_routes(risk_level);

-- ============ VESSELS ============
CREATE TABLE narcotics_vessels (
    vessel_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    imo_number VARCHAR(50),
    mmsi_number VARCHAR(50),
    flag_country VARCHAR(100),
    vessel_type VARCHAR(50) CHECK (vessel_type IN (
        'FISHING','CARGO','CONTAINER','PLEASURE','TANKER','GO_FAST','SUBmersible','OTHER'
    )),
    gross_tonnage DECIMAL(10,2),
    length_m DECIMAL(8,2),
    owner_company VARCHAR(255),
    operator_company VARCHAR(255),
    known_drug_links BOOLEAN DEFAULT FALSE,
    last_known_location GEOMETRY(Point, 4326),
    last_known_port VARCHAR(255),
    risk_score DECIMAL(5,2),
    status VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_vessel_name ON narcotics_vessels(name);
CREATE INDEX idx_vessel_imo ON narcotics_vessels(imo_number);
CREATE INDEX idx_vessel_flag ON narcotics_vessels(flag_country);
CREATE INDEX idx_vessel_risk ON narcotics_vessels(risk_score DESC);

-- ============ AIRCRAFT ============
CREATE TABLE narcotics_aircraft (
    aircraft_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tail_number VARCHAR(50) UNIQUE,
    icao_code VARCHAR(4),
    make VARCHAR(100),
    model VARCHAR(100),
    aircraft_type VARCHAR(50) CHECK (aircraft_type IN (
        'PRIVATE','CARGO','COMMERCIAL','CHARTER','MILITARY','DRONE','ULTRA_LIGHT'
    )),
    owner_company VARCHAR(255),
    operator_company VARCHAR(255),
    known_drug_links BOOLEAN DEFAULT FALSE,
    last_known_airport VARCHAR(255),
    risk_score DECIMAL(5,2),
    status VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_aircraft_tail ON narcotics_aircraft(tail_number);
CREATE INDEX idx_aircraft_risk ON narcotics_aircraft(risk_score DESC);

-- ============ SEIZURES ============
CREATE TABLE narcotics_seizures (
    seizure_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seizure_date TIMESTAMPTZ NOT NULL,
    seizure_location VARCHAR(500),
    latitude DECIMAL(10,7),
    longitude DECIMAL(10,7),
    seizure_type VARCHAR(30) CHECK (seizure_type IN ('AIRPORT','PORT','BORDER','ROAD','RESIDENTIAL','HIDDEN_COMPARTMENT','MAIL')),
    drug_types TEXT[] NOT NULL,
    drug_quantities JSONB NOT NULL,
    estimated_street_value_usd DECIMAL(20,2),
    concealment_method VARCHAR(255),
    seizure_agency VARCHAR(100) NOT NULL,
    operation_name VARCHAR(255),
    arrests_made INT DEFAULT 0,
    arrested_persons TEXT[],
    related_vehicle_plate VARCHAR(50),
    related_vessel_id UUID REFERENCES narcotics_vessels(vessel_id),
    related_aircraft_id UUID REFERENCES narcotics_aircraft(aircraft_id),
    related_route_id UUID REFERENCES narcotics_routes(route_id),
    related_cartel_id UUID REFERENCES cartels(cartel_id),
    case_id UUID,
    status VARCHAR(20) CHECK (status IN ('PENDING','EVIDENCE_PROCESSING','INTELLIGENCE_EXTRACTED','CLOSED')),
    created_at TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY RANGE (seizure_date);

CREATE INDEX idx_seizure_date ON narcotics_seizures(seizure_date DESC);
CREATE INDEX idx_seizure_drug ON narcotics_seizures USING gin(drug_types);
CREATE INDEX idx_seizure_agency ON narcotics_seizures(seizure_agency);
CREATE INDEX idx_seizure_cartel ON narcotics_seizures(related_cartel_id);
CREATE INDEX idx_seizure_location ON narcotics_seizures(seizure_location);

ALTER TABLE narcotics_seizures ENABLE ROW LEVEL SECURITY;
CREATE POLICY narcotics_pnh_select ON narcotics_seizures FOR SELECT USING (
    current_setting('snisid.agency') IN ('PNH','DCPJ','CUSTOMS','ANTI_NARCOTICS','SNISID_ADMIN')
);
