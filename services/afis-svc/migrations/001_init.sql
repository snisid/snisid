CREATE TYPE finger_position AS ENUM ('RIGHT_THUMB','RIGHT_INDEX','RIGHT_MIDDLE','RIGHT_RING','RIGHT_LITTLE','LEFT_THUMB','LEFT_INDEX','LEFT_MIDDLE','LEFT_RING','LEFT_LITTLE','RIGHT_PALM','LEFT_PALM','UNKNOWN');
CREATE TYPE capture_method AS ENUM ('LIVESCANNER','INKROLL','LATENT_LIFT','PHOTO','UNKNOWN');
CREATE TYPE subject_type AS ENUM ('SUSPECT','CRIMINAL','VICTIM','UNKNOWN_DECEASED','MISSING_PERSON','EMPLOYEE');
CREATE TYPE transaction_type AS ENUM ('TEN2TEN','LATENT2TEN','PALM');

CREATE TABLE subjects (
    subject_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snisid_person_id UUID UNIQUE,
    fir_record_id UUID UNIQUE,
    subject_type subject_type NOT NULL,
    national_afis_id VARCHAR(50) UNIQUE,
    enrolling_unit VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE fingerprints (
    print_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_id UUID NOT NULL REFERENCES subjects(subject_id),
    finger_position finger_position NOT NULL,
    capture_method capture_method NOT NULL DEFAULT 'UNKNOWN',
    nfiq2_score SMALLINT NOT NULL CHECK (nfiq2_score >= 0 AND nfiq2_score <= 100),
    quality_accepted BOOLEAN NOT NULL DEFAULT false,
    image_ref VARCHAR(500) NOT NULL,
    minutiae_count SMALLINT,
    milvus_vector_id VARCHAR(100),
    captured_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID NOT NULL
);

CREATE TABLE latent_prints (
    latent_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_reference VARCHAR(100) NOT NULL,
    crime_scene_id UUID,
    location_desc TEXT,
    dept_code VARCHAR(10),
    found_at TIMESTAMPTZ NOT NULL,
    image_ref VARCHAR(500) NOT NULL,
    nfiq2_score SMALLINT CHECK (nfiq2_score >= 0 AND nfiq2_score <= 100),
    finger_position finger_position NOT NULL DEFAULT 'UNKNOWN',
    is_identified BOOLEAN NOT NULL DEFAULT false,
    matched_subject_id UUID REFERENCES subjects(subject_id),
    match_score DOUBLE PRECISION,
    examined_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE search_transactions (
    transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_type transaction_type NOT NULL,
    query_subject_id UUID REFERENCES subjects(subject_id),
    query_latent_id UUID REFERENCES latent_prints(latent_id),
    hits_count SMALLINT NOT NULL DEFAULT 0,
    top_score DOUBLE PRECISION,
    top_match_id UUID,
    search_duration_ms INTEGER,
    requested_by UUID NOT NULL,
    requesting_unit VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_fingerprints_subject ON fingerprints(subject_id);
CREATE INDEX idx_fingerprints_position ON fingerprints(finger_position);
CREATE INDEX idx_latent_prints_case ON latent_prints(case_reference);
CREATE INDEX idx_latent_prints_dept ON latent_prints(dept_code);
CREATE INDEX idx_search_transactions_type ON search_transactions(transaction_type);
CREATE INDEX idx_subjects_national_afis ON subjects(national_afis_id);
