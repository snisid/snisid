BEGIN;

CREATE TABLE sivc_gang_associations (
    assoc_id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_id                UUID NOT NULL REFERENCES sivc_criminal_alerts(alert_id) ON DELETE CASCADE,
    gang_identifier         VARCHAR(100),
    gang_territory_dept     CHAR(2),
    gang_territory_communes TEXT[] DEFAULT '{}',
    gang_snisid_id          UUID,
    vehicle_role            VARCHAR(100),
    association_confidence  SMALLINT CHECK (association_confidence BETWEEN 1 AND 10),
    intelligence_source     VARCHAR(200),
    source_classification   VARCHAR(20) DEFAULT 'RESTRICTED',
    first_seen_date         TIMESTAMPTZ,
    last_confirmed_date     TIMESTAMPTZ,
    notes                   TEXT,
    created_by              UUID NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sivc_gang_assoc_gang     ON sivc_gang_associations(gang_identifier);
CREATE INDEX idx_sivc_gang_assoc_dept     ON sivc_gang_associations(gang_territory_dept);
CREATE INDEX idx_sivc_gang_assoc_alert    ON sivc_gang_associations(alert_id);

COMMIT;
