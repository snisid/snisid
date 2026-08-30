-- GANG-HT Migration: Registre National des Organisations Criminelles et Gangs

BEGIN;

CREATE TYPE gang_structure_type AS ENUM (
    'HIERARCHY','NETWORK','CELL','COALITION','FRANCHISE'
);

CREATE TYPE gang_activity_level AS ENUM (
    'DORMANT','LOW','MODERATE','HIGH','EXTREME'
);

CREATE TYPE gang_primary_activity AS ENUM (
    'KIDNAPPING','DRUG_TRAFFICKING','ARMS_TRAFFICKING',
    'EXTORTION','TERRITORY_CONTROL','CONTRACT_KILLING',
    'HUMAN_TRAFFICKING','MONEY_LAUNDERING','MIXED'
);

CREATE TABLE gang_organizations (
    gang_id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_gang_id         VARCHAR(25) UNIQUE NOT NULL,  -- GANG-HT-NNNNNN
    name                     VARCHAR(150) NOT NULL,
    aliases                  TEXT[] DEFAULT '{}',
    structure_type           gang_structure_type,
    primary_activity         gang_primary_activity NOT NULL,
    activity_level           gang_activity_level NOT NULL DEFAULT 'HIGH',
    estimated_members        INTEGER,
    armed_members_pct        SMALLINT,
    heavy_weapons            BOOLEAN DEFAULT FALSE,
    primary_dept_code        CHAR(2) NOT NULL,
    territory_communes       TEXT[] DEFAULT '{}',
    territory_geojson        JSONB,
    estimated_revenue_usd_monthly DECIMAL(12,2),
    primary_income_sources   TEXT[] DEFAULT '{}',
    un_designation_date      TIMESTAMPTZ,
    ofac_designation         BOOLEAN DEFAULT FALSE,
    ofac_sdn_ref             VARCHAR(50),
    allied_gang_ids          UUID[] DEFAULT '{}',
    rival_gang_ids           UUID[] DEFAULT '{}',
    established_date         DATE,
    current_leader_id        UUID,
    intel_confidence         SMALLINT CHECK (intel_confidence BETWEEN 1 AND 10),
    last_intel_update        TIMESTAMPTZ,
    is_active                BOOLEAN DEFAULT TRUE,
    created_by               UUID NOT NULL,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE gang_incidents (
    incident_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gang_id             UUID NOT NULL REFERENCES gang_organizations(gang_id),
    incident_type       VARCHAR(50) NOT NULL,
    incident_date       TIMESTAMPTZ NOT NULL,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    casualties          SMALLINT DEFAULT 0,
    victim_ids          UUID[] DEFAULT '{}',
    sivc_alert_id       UUID,
    description         TEXT,
    intelligence_source VARCHAR(100),
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE gang_alliances (
    alliance_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gang_a_id           UUID NOT NULL REFERENCES gang_organizations(gang_id),
    gang_b_id           UUID NOT NULL REFERENCES gang_organizations(gang_id),
    alliance_type       VARCHAR(30) NOT NULL,
    start_date          DATE,
    end_date            DATE,
    confidence_level    SMALLINT,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT no_self_alliance CHECK (gang_a_id <> gang_b_id)
);

CREATE INDEX idx_gang_dept     ON gang_organizations(primary_dept_code) WHERE is_active = TRUE;
CREATE INDEX idx_gang_activity ON gang_organizations(activity_level) WHERE is_active = TRUE;
CREATE INDEX idx_gang_ofac     ON gang_organizations(ofac_designation) WHERE ofac_designation = TRUE;
CREATE INDEX idx_gang_incidents ON gang_incidents(gang_id, incident_date DESC);
CREATE INDEX idx_gang_incid_dept ON gang_incidents(dept_code, incident_date DESC);

COMMIT;
