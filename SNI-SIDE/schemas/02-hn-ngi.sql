-- ============================================================
-- SNI-SIDE: HN-NGI (National Biometric Database)
-- PostgreSQL 16 + Milvus Vector Store
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_biometric_national;
SET search_path TO snisid_biometric_national;

-- ============ BIOMETRIC ENROLLMENTS ============
CREATE TABLE biometric_enrollments (
    enrollment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    niu VARCHAR(10) NOT NULL,
    person_name VARCHAR(255) NOT NULL,
    enrollment_date TIMESTAMPTZ DEFAULT NOW(),
    enrollment_agency VARCHAR(100),
    enrollment_location VARCHAR(500),
    enrollment_device_id VARCHAR(100),
    operator_id VARCHAR(100),
    quality_score DECIMAL(5,2) CHECK (quality_score >= 0 AND quality_score <= 100),
    status VARCHAR(20) CHECK (status IN ('ENROLLED','VERIFIED','FLAGGED','DUPLICATE','REJECTED')),
    CONSTRAINT fk_citizen FOREIGN KEY (niu) REFERENCES snisid_identity.citizens(niu)
);

-- ============ FINGERPRINTS ============
CREATE TABLE fingerprints (
    fingerprint_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    enrollment_id UUID NOT NULL REFERENCES biometric_enrollments(enrollment_id),
    niu VARCHAR(10) NOT NULL,
    finger_position VARCHAR(20) CHECK (finger_position IN (
        'RIGHT_THUMB','RIGHT_INDEX','RIGHT_MIDDLE','RIGHT_RING','RIGHT_LITTLE',
        'LEFT_THUMB','LEFT_INDEX','LEFT_MIDDLE','LEFT_RING','LEFT_LITTLE',
        'FULL_RIGHT_SLAP','FULL_LEFT_SLAP','FULL_THUMBS'
    )),
    image_hash VARCHAR(64) NOT NULL,
    image_format VARCHAR(10) DEFAULT 'WSQ',
    image_path VARCHAR(500) NOT NULL,
    minutiae_count INT,
    nfiq_score INT CHECK (nfiq_score >= 1 AND nfiq_score <= 5),
    template_milvus_id VARCHAR(100),
    quality_score DECIMAL(5,2),
    status VARCHAR(20) DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_finger_niu_pos UNIQUE (niu, finger_position)
);

CREATE INDEX idx_fp_niu ON fingerprints(niu);
CREATE INDEX idx_fp_position ON fingerprints(finger_position);
CREATE INDEX idx_fp_nfiq ON fingerprints(nfiq_score);
CREATE INDEX idx_fp_milvus ON fingerprints(template_milvus_id);

-- ============ PALM PRINTS ============
CREATE TABLE palmprints (
    palm_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    enrollment_id UUID NOT NULL REFERENCES biometric_enrollments(enrollment_id),
    niu VARCHAR(10) NOT NULL,
    palm_position VARCHAR(20) CHECK (palm_position IN ('RIGHT_FULL','LEFT_FULL','RIGHT_WRITER','LEFT_WRITER','RIGHT_LOWER','LEFT_LOWER')),
    image_hash VARCHAR(64) NOT NULL,
    image_path VARCHAR(500) NOT NULL,
    template_milvus_id VARCHAR(100),
    quality_score DECIMAL(5,2),
    status VARCHAR(20) DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_palm_niu ON palmprints(niu);

-- ============ IRIS ============
CREATE TABLE irises (
    iris_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    enrollment_id UUID NOT NULL REFERENCES biometric_enrollments(enrollment_id),
    niu VARCHAR(10) NOT NULL,
    iris_position VARCHAR(10) CHECK (iris_position IN ('RIGHT','LEFT')),
    image_hash VARCHAR(64) NOT NULL,
    image_path VARCHAR(500) NOT NULL,
    template_milvus_id VARCHAR(100),
    quality_score DECIMAL(5,2),
    status VARCHAR(20) DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_iris_niu_pos UNIQUE (niu, iris_position)
);

CREATE INDEX idx_iris_niu ON irises(niu);

-- ============ FACES ============
CREATE TABLE faces (
    face_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    enrollment_id UUID NOT NULL REFERENCES biometric_enrollments(enrollment_id),
    niu VARCHAR(10) NOT NULL,
    image_hash VARCHAR(64) NOT NULL,
    image_path VARCHAR(500) NOT NULL,
    embedding_milvus_id VARCHAR(100),
    embedding_dim INT DEFAULT 512,
    quality_score DECIMAL(5,2),
    pose_yaw DECIMAL(5,2),
    pose_pitch DECIMAL(5,2),
    pose_roll DECIMAL(5,2),
    face_bbox JSONB,
    landmarks JSONB,
    status VARCHAR(20) DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_face_niu ON faces(niu);
CREATE INDEX idx_face_milvus ON faces(embedding_milvus_id);
CREATE INDEX idx_face_quality ON faces(quality_score DESC);

-- ============ VOICE ============
CREATE TABLE voice_samples (
    voice_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    enrollment_id UUID NOT NULL REFERENCES biometric_enrollments(enrollment_id),
    niu VARCHAR(10) NOT NULL,
    audio_hash VARCHAR(64) NOT NULL,
    audio_path VARCHAR(500) NOT NULL,
    audio_duration_sec DECIMAL(10,2),
    sample_rate INT DEFAULT 16000,
    embedding_milvus_id VARCHAR(100),
    language VARCHAR(50),
    quality_score DECIMAL(5,2),
    status VARCHAR(20) DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_voice_niu ON voice_samples(niu);

-- ============ MULTIMODAL REFERENCES ============
CREATE TABLE multimodal_references (
    ref_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    niu VARCHAR(10) NOT NULL,
    fingerprints_present INT DEFAULT 0,
    palms_present INT DEFAULT 0,
    irises_present INT DEFAULT 0,
    faces_present INT DEFAULT 0,
    voices_present INT DEFAULT 0,
    overall_quality DECIMAL(5,2),
    gallery_id VARCHAR(100),
    status VARCHAR(20) DEFAULT 'COMPLETE',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_citizen FOREIGN KEY (niu) REFERENCES snisid_identity.citizens(niu),
    CONSTRAINT uq_mm_niu UNIQUE (niu)
);

-- ============ DUPLICATE DETECTION LOG ============
CREATE TABLE duplicate_detection_log (
    detection_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    niu_primary VARCHAR(10) NOT NULL,
    niu_candidate VARCHAR(10) NOT NULL,
    biometric_type VARCHAR(20) NOT NULL,
    match_score DECIMAL(8,4) NOT NULL,
    threshold DECIMAL(8,4) NOT NULL,
    algorithm VARCHAR(100),
    status VARCHAR(20) CHECK (status IN ('PENDING_REVIEW','CONFIRMED_DUPLICATE','FALSE_POSITIVE','RESOLVED')),
    reviewed_by VARCHAR(100),
    reviewed_at TIMESTAMPTZ,
    resolution_notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_dup_primary ON duplicate_detection_log(niu_primary);
CREATE INDEX idx_dup_candidate ON duplicate_detection_log(niu_candidate);
CREATE INDEX idx_dup_status ON duplicate_detection_log(status);
CREATE INDEX idx_dup_score ON duplicate_detection_log(match_score DESC);
CREATE INDEX idx_dup_created ON duplicate_detection_log(created_at DESC);

-- ============ BIOMETRIC SEARCH LOG ============
CREATE TABLE biometric_search_log (
    search_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    search_type VARCHAR(30) CHECK (search_type IN (
        'VERIFY_1_1','IDENTIFY_1_N','FACE_SEARCH','FINGERPRINT_SEARCH',
        'PALM_SEARCH','IRIS_SEARCH','VOICE_SEARCH','MULTIMODAL'
    )),
    query_biometric_type VARCHAR(20),
    query_hash VARCHAR(64),
    result_count INT,
    top_match_niu VARCHAR(10),
    top_match_score DECIMAL(8,4),
    search_duration_ms INT,
    user_id VARCHAR(100),
    agency VARCHAR(100),
    purpose VARCHAR(500),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_bio_search_type ON biometric_search_log(search_type);
CREATE INDEX idx_bio_search_agency ON biometric_search_log(agency);
CREATE INDEX idx_bio_search_time ON biometric_search_log(created_at DESC);

-- ============ RLS ============
ALTER TABLE biometric_enrollments ENABLE ROW LEVEL SECURITY;
ALTER TABLE faces ENABLE ROW LEVEL SECURITY;
ALTER TABLE fingerprints ENABLE ROW LEVEL SECURITY;

CREATE POLICY bio_pnh_select ON faces FOR SELECT USING (
    current_setting('snisid.agency') IN ('PNH','DCPJ','ONI','SNISID_ADMIN')
);
