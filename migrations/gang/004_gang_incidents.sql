BEGIN;

CREATE TABLE gang_incidents (
    incident_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gang_id             UUID NOT NULL REFERENCES gang_organizations(gang_id) ON DELETE CASCADE,
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

COMMIT;
