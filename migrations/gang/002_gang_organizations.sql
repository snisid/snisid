BEGIN;

CREATE TABLE gang_organizations (
    gang_id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_gang_id         VARCHAR(25) UNIQUE NOT NULL,
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

COMMIT;
