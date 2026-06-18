BEGIN;

CREATE TABLE sifr_alerts_log (
    alert_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crossing_id         UUID REFERENCES sifr_crossings(crossing_id),
    post_id             UUID NOT NULL,
    alert_type          sifr_alert_type NOT NULL,
    snisid_person_id    UUID,
    document_number     VARCHAR(100),
    vehicle_plate       VARCHAR(20),
    alert_source        VARCHAR(50),
    source_record_id    UUID,
    notified_units      TEXT[] DEFAULT '{}',
    action_taken        TEXT,
    resolved            BOOLEAN DEFAULT FALSE,
    resolved_by         UUID,
    resolved_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
