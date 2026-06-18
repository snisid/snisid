BEGIN;
CREATE TYPE trait_type AS ENUM ('LABOR_EXPLOITATION','SEXUAL_EXPLOITATION','FORCED_MARRIAGE','CHILD_DOMESTIC_SERVITUDE','GANG_RECRUITMENT_FORCED','IRREGULAR_MIGRATION_FACILITATION','ORGAN_TRAFFICKING','OTHER');
CREATE TYPE trait_victim_status AS ENUM ('IDENTIFIED_VICTIM','POTENTIAL_VICTIM','WITNESS','RESCUED','REPATRIATED','DECEASED','MISSING');
CREATE TABLE trait_cases (
    case_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), national_trait_id VARCHAR(25) UNIQUE NOT NULL,
    trait_type trait_type NOT NULL, status VARCHAR(20) DEFAULT 'OPEN',
    victim_count SMALLINT DEFAULT 1, minor_count SMALLINT DEFAULT 0,
    origin_country CHAR(3) DEFAULT 'HTI', transit_countries CHAR(3)[] DEFAULT '{}',
    destination_country CHAR(3), route_description TEXT,
    transport_mode TEXT[] DEFAULT '{}', mar_incident_id UUID,
    sifr_crossing_ids UUID[] DEFAULT '{}', gang_id UUID,
    recruiter_ids UUID[] DEFAULT '{}', total_amount_paid DECIMAL(12,2),
    amount_per_person DECIMAL(10,2), currency CHAR(3) DEFAULT 'USD',
    investigating_unit VARCHAR(50), case_reference VARCHAR(100),
    iom_case_ref VARCHAR(50), created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE trait_victims (
    victim_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id UUID NOT NULL REFERENCES trait_cases(case_id),
    snisid_person_id UUID, victim_status trait_victim_status NOT NULL,
    full_name VARCHAR(200), nationality CHAR(3) DEFAULT 'HTI',
    dob DATE, gender VARCHAR(10), is_minor BOOLEAN DEFAULT FALSE,
    exploitation_type TEXT, rescue_date TIMESTAMPTZ,
    rescue_location VARCHAR(300), current_location TEXT,
    assistance_provided TEXT[] DEFAULT '{}', dipe_case_id UUID,
    afis_subject_id UUID, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE trait_networks (
    network_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), network_name VARCHAR(150),
    primary_route TEXT, origin_dept CHAR(2), known_members UUID[] DEFAULT '{}',
    gang_affiliations UUID[] DEFAULT '{}', monthly_volume_est INTEGER,
    fee_per_person_usd DECIMAL(10,2), is_active BOOLEAN DEFAULT TRUE,
    intel_confidence SMALLINT, linked_cases UUID[] DEFAULT '{}',
    created_by UUID NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_trait_cases_type ON trait_cases(trait_type, status);
CREATE INDEX idx_trait_cases_gang ON trait_cases(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_trait_victims_case ON trait_victims(case_id);
CREATE INDEX idx_trait_victims_minor ON trait_victims(is_minor) WHERE is_minor = TRUE;
CREATE INDEX idx_trait_victims_snisid ON trait_victims(snisid_person_id) WHERE snisid_person_id IS NOT NULL;
CREATE INDEX idx_trait_networks_active ON trait_networks(is_active) WHERE is_active = TRUE;
COMMIT;
