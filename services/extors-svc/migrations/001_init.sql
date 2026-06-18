BEGIN;
CREATE TYPE extors_type AS ENUM ('KIDNAPPING_RANSOM','ROAD_TOLL_ILLEGAL','BUSINESS_PROTECTION_RACKET','REAL_ESTATE_EXTORTION','PUBLIC_SERVANT_EXTORTION','NGO_EXTORTION','FUEL_TRUCK_HIJACK','OTHER');
CREATE TYPE extors_payment_channel AS ENUM ('MONCASH','NATCASH','DIGICEL_MONEY','WIRE_TRANSFER','CASH_DROP','CRYPTOCURRENCY','INTERMEDIARY','UNKNOWN');
CREATE TYPE extors_status AS ENUM ('ACTIVE','PAID','REFUSED','NEGOTIATING','LAW_ENFORCEMENT_INVOLVED','RESOLVED','VICTIM_HARMED');
CREATE TABLE extors_cases (
    case_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), national_extors_id VARCHAR(25) UNIQUE NOT NULL,
    extors_type extors_type NOT NULL, status extors_status NOT NULL DEFAULT 'ACTIVE',
    gang_id UUID, gang_name VARCHAR(150), perpetrator_ids UUID[] DEFAULT '{}',
    chef_member_ids UUID[] DEFAULT '{}', victim_count SMALLINT DEFAULT 1,
    victim_snisid_ids UUID[] DEFAULT '{}', victim_types TEXT[] DEFAULT '{}',
    victim_nationality CHAR(3)[] DEFAULT '{}', is_foreigner_victim BOOLEAN DEFAULT FALSE,
    incident_location VARCHAR(300), dept_code CHAR(2), commune VARCHAR(100),
    lat DECIMAL(10,7), lng DECIMAL(10,7), route_number VARCHAR(10),
    demanded_amount DECIMAL(15,2), demanded_currency CHAR(3) DEFAULT 'USD',
    paid_amount DECIMAL(15,2), paid_currency CHAR(3),
    payment_channel extors_payment_channel, payment_ref VARCHAR(200),
    payment_date TIMESTAMPTZ, first_contact_date TIMESTAMPTZ NOT NULL,
    resolution_date TIMESTAMPTZ, case_reference VARCHAR(100),
    investigating_unit VARCHAR(50), ucref_str_id UUID, blan_case_id UUID,
    notes TEXT, created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE extors_road_toll_points (
    toll_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), gang_id UUID NOT NULL,
    location_desc VARCHAR(300) NOT NULL, route_number VARCHAR(10),
    dept_code CHAR(2) NOT NULL, commune VARCHAR(100), lat DECIMAL(10,7),
    lng DECIMAL(10,7), daily_revenue_usd DECIMAL(10,2),
    vehicle_types_taxed TEXT[] DEFAULT '{}', toll_rates JSONB,
    active_since DATE, is_active BOOLEAN DEFAULT TRUE, source_intel TEXT,
    last_confirmed_at TIMESTAMPTZ, created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE extors_negotiations (
    neg_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), case_id UUID NOT NULL REFERENCES extors_cases(case_id),
    negotiation_date TIMESTAMPTZ NOT NULL, contact_method VARCHAR(50),
    contact_number VARCHAR(30), demand_updated DECIMAL(15,2),
    demand_currency CHAR(3), position_update TEXT, recorded_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_extors_cases_gang ON extors_cases(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_extors_cases_type ON extors_cases(extors_type, status);
CREATE INDEX idx_extors_cases_dept ON extors_cases(dept_code, first_contact_date DESC);
CREATE INDEX idx_extors_cases_channel ON extors_cases(payment_channel) WHERE paid_amount IS NOT NULL;
CREATE INDEX idx_extors_tolls_route ON extors_road_toll_points(route_number) WHERE is_active = TRUE;
CREATE INDEX idx_extors_tolls_dept ON extors_road_toll_points(dept_code) WHERE is_active = TRUE;
COMMIT;
