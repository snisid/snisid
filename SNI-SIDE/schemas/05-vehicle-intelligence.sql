-- ============================================================
-- SNI-SIDE: Vehicle Intelligence Database
-- PostgreSQL 16
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_vehicle;
SET search_path TO snisid_vehicle;

-- ============ VEHICLES ============
CREATE TABLE vehicles (
    vehicle_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vin VARCHAR(17) UNIQUE NOT NULL,
    plate_number VARCHAR(50),
    plate_country VARCHAR(10) DEFAULT 'HT',
    plate_state VARCHAR(100),
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER CHECK (year >= 1900 AND year <= 2035),
    color VARCHAR(50),
    color_secondary VARCHAR(50),
    body_type VARCHAR(50) CHECK (body_type IN (
        'SEDAN','SUV','TRUCK','MOTORCYCLE','BUS','VAN','COUPE','CONVERTIBLE','HATCHBACK','WAGON','PICKUP','OTHER'
    )),
    engine_number VARCHAR(100),
    engine_type VARCHAR(50),
    fuel_type VARCHAR(20),
    weight_kg DECIMAL(10,2),
    seats INTEGER,
    registration_number VARCHAR(100),
    registration_date DATE,
    registration_expiry DATE,
    registration_agency VARCHAR(100),
    owner_niu VARCHAR(10),
    owner_name VARCHAR(255),
    insurance_provider VARCHAR(255),
    insurance_policy VARCHAR(100),
    insurance_expiry DATE,
    status VARCHAR(30) CHECK (status IN (
        'REGISTERED','STOLEN','WANTED','SUSPENDED','DEREGISTERED','SALVAGE','EXPORTED'
    )),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_owner FOREIGN KEY (owner_niu) REFERENCES snisid_identity.citizens(niu)
);

CREATE INDEX idx_vehicle_vin ON vehicles(vin);
CREATE INDEX idx_vehicle_plate ON vehicles(plate_number);
CREATE INDEX idx_vehicle_make ON vehicles(make);
CREATE INDEX idx_vehicle_year ON vehicles(year);
CREATE INDEX idx_vehicle_color ON vehicles(color);
CREATE INDEX idx_vehicle_owner ON vehicles(owner_niu);
CREATE INDEX idx_vehicle_status ON vehicles(status);
CREATE INDEX idx_vehicle_reg ON vehicles(registration_number);

-- ============ OWNERSHIP HISTORY ============
CREATE TABLE vehicle_ownership_history (
    ownership_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id UUID NOT NULL REFERENCES vehicles(vehicle_id),
    owner_niu VARCHAR(10) NOT NULL,
    owner_name VARCHAR(255) NOT NULL,
    ownership_start DATE NOT NULL,
    ownership_end DATE,
    transfer_reason VARCHAR(100),
    transfer_location VARCHAR(500),
    sale_price DECIMAL(20,2),
    notary_ref VARCHAR(100),
    agency VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_owner_hist_vehicle ON vehicle_ownership_history(vehicle_id);
CREATE INDEX idx_owner_hist_owner ON vehicle_ownership_history(owner_niu);

-- ============ STOLEN VEHICLES ============
CREATE TABLE stolen_vehicles (
    theft_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id UUID NOT NULL REFERENCES vehicles(vehicle_id),
    theft_date TIMESTAMPTZ NOT NULL,
    theft_location VARCHAR(500),
    theft_lat DECIMAL(10,7),
    theft_lng DECIMAL(10,7),
    theft_method VARCHAR(100),
    reported_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reported_by VARCHAR(255),
    reporting_agency VARCHAR(100) NOT NULL,
    insurance_claim_ref VARCHAR(100),
    recovery_date TIMESTAMPTZ,
    recovery_location VARCHAR(500),
    recovery_condition VARCHAR(100),
    recovery_by_agency VARCHAR(100),
    status VARCHAR(20) CHECK (status IN ('STOLEN','RECOVERED','WRITTEN_OFF','EXPORTED')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_stolen_vehicle ON stolen_vehicles(vehicle_id);
CREATE INDEX idx_stolen_status ON stolen_vehicles(status);
CREATE INDEX idx_stolen_date ON stolen_vehicles(theft_date DESC);

-- ============ VEHICLE NARCOTICS LINKS ============
CREATE TABLE vehicle_narcotics_links (
    link_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id UUID NOT NULL REFERENCES vehicles(vehicle_id),
    case_id UUID,
    detection_date TIMESTAMPTZ NOT NULL,
    detection_location VARCHAR(500),
    drug_type VARCHAR(100),
    quantity_kg DECIMAL(15,4),
    concealment_method VARCHAR(255),
    seizure_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============ VEHICLE TERRORISM LINKS ============
CREATE TABLE vehicle_terrorism_links (
    link_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id UUID NOT NULL REFERENCES vehicles(vehicle_id),
    case_id UUID,
    incident_date TIMESTAMPTZ,
    incident_location VARCHAR(500),
    link_type VARCHAR(50),
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ============ VEHICLE SEARCH LOG ============
CREATE TABLE vehicle_search_log (
    search_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    search_type VARCHAR(20) CHECK (search_type IN ('VIN','PLATE','OWNER','NETWORK')),
    search_query VARCHAR(255) NOT NULL,
    result_count INT,
    user_id VARCHAR(100),
    agency VARCHAR(100),
    search_duration_ms INT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_vehicle_search_type ON vehicle_search_log(search_type);
CREATE INDEX idx_vehicle_search_agency ON vehicle_search_log(agency);

ALTER TABLE vehicles ENABLE ROW LEVEL SECURITY;
CREATE POLICY vehicle_pnh_select ON vehicles FOR SELECT USING (
    current_setting('snisid.agency') IN ('PNH','DCPJ','TRANSPORT','SNISID_ADMIN')
);
