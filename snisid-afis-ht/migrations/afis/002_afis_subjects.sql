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

CREATE INDEX idx_afis_subjects_snisid ON afis_subjects(snisid_person_id) WHERE snisid_person_id IS NOT NULL;
CREATE INDEX idx_afis_subjects_fir ON afis_subjects(fir_record_id) WHERE fir_record_id IS NOT NULL;
CREATE INDEX idx_afis_subjects_national_id ON afis_subjects(national_afis_id);
CREATE INDEX idx_afis_subjects_type ON afis_subjects(subject_type);

COMMIT;