BEGIN;
CREATE TYPE port_risk_level AS ENUM ('LOW','MEDIUM','HIGH','CRITICAL');
CREATE TYPE port_container_status AS ENUM ('PENDING_INSPECTION','CLEARED','HELD_FOR_INSPECTION','SEIZED','RELEASED_AFTER_INSPECTION');
CREATE TABLE port_vessels_arrivals (
    arrival_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), port_code VARCHAR(10) NOT NULL,
    vessel_imo VARCHAR(20), vessel_name VARCHAR(150) NOT NULL, flag_country CHAR(3),
    shipping_company VARCHAR(200), arrival_date TIMESTAMPTZ NOT NULL,
    origin_port VARCHAR(100), origin_country CHAR(3), container_count INTEGER DEFAULT 0,
    manifest_ref VARCHAR(100), mar_vessel_id UUID, risk_score SMALLINT DEFAULT 0,
    risk_level port_risk_level DEFAULT 'LOW', cbp_targeting_ref VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE port_containers (
    container_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), arrival_id UUID NOT NULL REFERENCES port_vessels_arrivals(arrival_id),
    container_number VARCHAR(20) NOT NULL, container_type VARCHAR(10),
    declared_content TEXT NOT NULL, declared_weight_kg DECIMAL(12,3),
    declared_value_usd DECIMAL(15,2), shipper_name VARCHAR(200), shipper_country CHAR(3),
    consignee_name VARCHAR(200), consignee_snisid_id UUID,
    status port_container_status NOT NULL DEFAULT 'PENDING_INSPECTION',
    risk_score SMALLINT DEFAULT 0, risk_level port_risk_level DEFAULT 'LOW',
    risk_flags TEXT[] DEFAULT '{}', selected_for_scan BOOLEAN DEFAULT FALSE,
    scan_date TIMESTAMPTZ, scan_result TEXT, seized BOOLEAN DEFAULT FALSE,
    seizure_description TEXT, case_reference VARCHAR(100),
    cbp_targeting_match BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE port_risk_factors (
    factor_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), container_id UUID NOT NULL REFERENCES port_containers(container_id),
    factor_type VARCHAR(50) NOT NULL, description TEXT NOT NULL,
    weight_score SMALLINT NOT NULL, source VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_port_containers_risk ON port_containers(risk_level, status);
CREATE INDEX idx_port_containers_arrival ON port_containers(arrival_id);
CREATE INDEX idx_port_arrivals_date ON port_vessels_arrivals(arrival_date DESC);
CREATE INDEX idx_port_arrivals_port ON port_vessels_arrivals(port_code, arrival_date DESC);
COMMIT;
