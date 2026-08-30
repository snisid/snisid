BEGIN;

CREATE TABLE afis_fingerprints (
    print_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_id          UUID NOT NULL REFERENCES afis_subjects(subject_id),
    finger_position     afis_finger_position NOT NULL,
    capture_method      afis_capture_method NOT NULL DEFAULT 'LIVESCANNER',
    nfiq2_score         SMALLINT CHECK (nfiq2_score BETWEEN 0 AND 100),
    quality_accepted    BOOLEAN GENERATED ALWAYS AS (nfiq2_score >= 60) STORED,
    image_ref           VARCHAR(500) NOT NULL,
    minutiae_count      SMALLINT,
    milvus_vector_id    VARCHAR(100),
    template_version    VARCHAR(10) DEFAULT 'ISO_2011',
    is_primary          BOOLEAN DEFAULT FALSE,
    captured_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
