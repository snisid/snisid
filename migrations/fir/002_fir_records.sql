BEGIN;

CREATE TABLE fir_criminal_records (
    record_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_fir_id     VARCHAR(25) UNIQUE NOT NULL,
    snisid_person_id    UUID NOT NULL,
    afis_subject_id     UUID,
    is_haitian_national BOOLEAN DEFAULT TRUE,
    is_active           BOOLEAN DEFAULT TRUE,
    is_expunged         BOOLEAN DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
