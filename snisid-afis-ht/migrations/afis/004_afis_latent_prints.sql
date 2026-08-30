BEGIN;

CREATE TABLE afis_latent_prints (
    latent_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_reference      VARCHAR(100) NOT NULL,
    crime_scene_id      UUID,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    found_at            TIMESTAMPTZ NOT NULL,
    image_ref           VARCHAR(500) NOT NULL,
    nfiq2_score         SMALLINT,
    finger_position     afis_finger_position DEFAULT 'UNKNOWN',
    is_identified       BOOLEAN DEFAULT FALSE,
    matched_subject_id  UUID REFERENCES afis_subjects(subject_id),
    match_score         DECIMAL(5,2),
    examined_by         UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_afis_latent_case ON afis_latent_prints(case_reference);
CREATE INDEX idx_afis_latent_identified ON afis_latent_prints(is_identified);
CREATE INDEX idx_afis_latent_matched ON afis_latent_prints(matched_subject_id) WHERE matched_subject_id IS NOT NULL;
CREATE INDEX idx_afis_latent_found ON afis_latent_prints(found_at DESC);

COMMIT;