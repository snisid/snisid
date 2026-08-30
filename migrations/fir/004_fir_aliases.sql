BEGIN;

CREATE TABLE fir_aliases (
    alias_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id   UUID NOT NULL REFERENCES fir_criminal_records(record_id),
    first_name  VARCHAR(100),
    last_name   VARCHAR(100),
    birth_date  DATE,
    id_document VARCHAR(50),
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
