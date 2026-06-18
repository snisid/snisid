BEGIN;

CREATE TABLE sifr_clandestine_crossings (
    report_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    reported_date       TIMESTAMPTZ NOT NULL,
    crossing_type       VARCHAR(50),
    estimated_persons   INTEGER,
    gang_related        BOOLEAN DEFAULT FALSE,
    gang_id             UUID,
    trafficking_type    TEXT,
    reported_by         UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
