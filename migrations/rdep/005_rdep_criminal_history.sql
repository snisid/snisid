BEGIN;

CREATE TABLE rdep_foreign_records (
    foreign_record_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deportee_id           UUID NOT NULL REFERENCES rdep_deportees(deportee_id),
    country               rdep_deportation_country NOT NULL,
    court_name            VARCHAR(200),
    offense_description   TEXT NOT NULL,
    offense_date          TIMESTAMPTZ,
    conviction_date       TIMESTAMPTZ,
    sentence              TEXT,
    prison_served         TEXT,
    fbi_number            VARCHAR(50),
    interpol_ref          VARCHAR(50),
    source_document       VARCHAR(500),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rdep_monitoring_events (
    event_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deportee_id           UUID NOT NULL REFERENCES rdep_deportees(deportee_id),
    event_type            rdep_event_type NOT NULL,
    event_date            TIMESTAMPTZ NOT NULL,
    location_lat          DECIMAL(10,7),
    location_lng          DECIMAL(10,7),
    notes                 TEXT,
    reported_by           UUID NOT NULL,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
