BEGIN;

CREATE TABLE mar_watch_vessels (
    watch_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vessel_id           UUID REFERENCES mar_vessels(vessel_id),
    mmsi                VARCHAR(15),
    vessel_name         VARCHAR(150),
    watch_reason        TEXT NOT NULL,
    alert_level         VARCHAR(20) DEFAULT 'CAUTION',
    requesting_unit     VARCHAR(50),
    is_active           BOOLEAN DEFAULT TRUE,
    expiry_date         TIMESTAMPTZ,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
