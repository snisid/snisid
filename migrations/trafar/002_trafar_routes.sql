BEGIN;

CREATE TABLE trafar_routes (
    route_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_name          VARCHAR(150) NOT NULL,
    route_type          trafar_route_type NOT NULL,
    trafficking_method  trafar_method NOT NULL,
    origin_country      CHAR(3) NOT NULL,
    origin_city         VARCHAR(100),
    transit_points      JSONB,
    entry_point_haiti   VARCHAR(100),
    entry_dept_code     CHAR(2),
    associated_gang_ids UUID[] DEFAULT '{}',
    known_suppliers     TEXT[] DEFAULT '{}',
    activity_level      VARCHAR(20) DEFAULT 'ACTIVE',
    estimated_volume_monthly INTEGER,
    weapon_types        TEXT[] DEFAULT '{}',
    intel_confidence    SMALLINT CHECK (intel_confidence BETWEEN 1 AND 10),
    first_detected      DATE,
    last_confirmed      DATE,
    linked_case_refs    TEXT[] DEFAULT '{}',
    biar_weapon_ids     UUID[] DEFAULT '{}',
    atf_case_refs       TEXT[] DEFAULT '{}',
    unodc_ref           VARCHAR(50),
    analyst_notes       TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
