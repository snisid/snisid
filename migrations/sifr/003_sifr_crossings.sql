BEGIN;

CREATE TABLE sifr_crossings (
    crossing_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id             UUID NOT NULL REFERENCES sifr_border_posts(post_id),
    direction           sifr_crossing_direction NOT NULL,
    crossing_datetime   TIMESTAMPTZ NOT NULL,
    snisid_person_id    UUID,
    document_type       sifr_doc_type NOT NULL DEFAULT 'PASSPORT',
    document_number     VARCHAR(100),
    document_country    CHAR(3),
    document_expiry     DATE,
    traveler_name       VARCHAR(200) NOT NULL,
    traveler_dob        DATE,
    traveler_nationality CHAR(3),
    vehicle_plate       VARCHAR(20),
    lane_number         SMALLINT,
    processing_officer  UUID NOT NULL,
    alert_triggered     BOOLEAN DEFAULT FALSE,
    alert_type          sifr_alert_type,
    alert_action_taken  TEXT,
    processing_time_sec INTEGER,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
