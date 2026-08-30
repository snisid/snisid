BEGIN;

CREATE TABLE afis_subjects (
    subject_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snisid_person_id    UUID,
    fir_record_id       UUID,
    subject_type        afis_subject_type NOT NULL,
    national_afis_id    VARCHAR(20) UNIQUE,
    alias_ids           UUID[] DEFAULT '{}',
    enrolment_date      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    enrolling_unit      VARCHAR(50) NOT NULL,
    enrolling_officer   UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
