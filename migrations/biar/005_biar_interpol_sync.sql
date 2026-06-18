BEGIN;

CREATE TABLE biar_iarms_sync_log (
    sync_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    weapon_id           UUID REFERENCES biar_illicit_weapons(weapon_id),
    direction           VARCHAR(10) NOT NULL,
    iarms_ref           VARCHAR(50),
    sync_status         VARCHAR(20) DEFAULT 'PENDING',
    synced_at           TIMESTAMPTZ,
    error_message       TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
