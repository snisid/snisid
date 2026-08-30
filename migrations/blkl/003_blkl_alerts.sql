BEGIN;

CREATE TABLE blkl_alerts_log (
    alert_log_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id            UUID NOT NULL REFERENCES blkl_blacklist(entry_id),
    triggered_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    post_code           VARCHAR(10),
    direction           VARCHAR(10),
    action_taken        TEXT,
    officer_id          UUID,
    outcome             TEXT
);

COMMIT;
