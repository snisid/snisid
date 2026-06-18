-- ============================================================
-- SNI-SIDE: GEOINT Database
-- PostgreSQL 16 + PostGIS
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_geoint;
SET search_path TO snisid_geoint;

-- ============ GEOSPATIAL INTELLIGENCE ============
CREATE TABLE geoint_layers (
    layer_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    layer_name VARCHAR(255) NOT NULL,
    layer_type VARCHAR(30) CHECK (layer_type IN (
        'SATELLITE','DRONE','MAP','TERRAIN','HEATMAP','RISK_ZONE','BOUNDARY','INFRASTRUCTURE'
    )),
    source VARCHAR(255),
    acquisition_date TIMESTAMPTZ,
    geometry GEOMETRY(GEOMETRY, 4326),
    properties JSONB DEFAULT '{}',
    resolution_m DECIMAL(10,2),
    classification VARCHAR(20) CHECK (classification IN ('UNCLASSIFIED','RESTRICTED','CONFIDENTIAL','SECRET','TOP_SECRET')),
    status VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_geoint_layer ON geoint_layers USING gist(geometry);
CREATE INDEX idx_geoint_type ON geoint_layers(layer_type);
CREATE INDEX idx_geoint_date ON geoint_layers(acquisition_date DESC);

-- ============ HOTSPOTS ============
CREATE TABLE hotspots (
    hotspot_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hotspot_type VARCHAR(50) CHECK (hotspot_type IN (
        'CRIME_HOTSPOT','NARCOTICS_ROUTE','BORDER_CROSSING','HUMAN_TRAFFICKING',
        'WEAPONS_SMUGGLING','TERRORIST_ACTIVITY','DRONE_SIGHTING','SUSPICIOUS_VESSEL'
    )),
    name VARCHAR(255),
    description TEXT,
    location GEOMETRY(POINT, 4326),
    region GEOMETRY(POLYGON, 4326),
    risk_level VARCHAR(20),
    frequency INT,
    last_incident_date TIMESTAMPTZ,
    intelligence_feed TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_hotspot_point ON hotspots USING gist(location);
CREATE INDEX idx_hotspot_region ON hotspots USING gist(region);
CREATE INDEX idx_hotspot_type ON hotspots(hotspot_type);
CREATE INDEX idx_hotspot_risk ON hotspots(risk_level);
CREATE INDEX idx_hotspot_last ON hotspots(last_incident_date DESC);

-- ============ DRONE INTELLIGENCE ============
CREATE TABLE drone_intelligence (
    drone_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    drone_model VARCHAR(255),
    serial_number VARCHAR(100),
    operator_name VARCHAR(255),
    flight_date TIMESTAMPTZ,
    flight_path GEOMETRY(LINESTRING, 4326),
    altitude_m DECIMAL(8,2),
    speed_ms DECIMAL(6,2),
    camera_footprint GEOMETRY(POLYGON, 4326),
    captured_imagery TEXT[],
    analysis_notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_drone_path ON drone_intelligence USING gist(flight_path);
CREATE INDEX idx_drone_date ON drone_intelligence(flight_date DESC);

-- ============ RISK ZONES ============
CREATE TABLE risk_zones (
    zone_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_name VARCHAR(255) NOT NULL,
    zone_type VARCHAR(50) CHECK (zone_type IN (
        'HIGH_CRIME','TERRORIST','BORDER','NARCOTICS','HUMAN_TRAFFICKING',
        'KIDNAPPING','GANG_TERRITORY','NATURAL_DISASTER','CRITICAL_INFRASTRUCTURE'
    )),
    boundary GEOMETRY(POLYGON, 4326) NOT NULL,
    risk_level VARCHAR(20),
    risk_score DECIMAL(5,2),
    population_affected INT,
    advisories TEXT[],
    last_assessed TIMESTAMPTZ,
    status VARCHAR(20) CHECK (status IN ('ACTIVE','MONITORED','CLEARED')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_risk_zone ON risk_zones USING gist(boundary);
CREATE INDEX idx_risk_type ON risk_zones(zone_type);
CREATE INDEX idx_risk_level ON risk_zones(risk_level);
CREATE INDEX idx_risk_score ON risk_zones(risk_score DESC);

ALTER TABLE geoint_layers ENABLE ROW LEVEL SECURITY;
CREATE POLICY geoint_national_select ON geoint_layers FOR SELECT USING (
    current_setting('snisid.agency') IN ('DEFENSE','PNH','INTELLIGENCE','SNISID_ADMIN')
);
