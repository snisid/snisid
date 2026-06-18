-- ============================================================
-- SNI-SIDE: National Digital Evidence Repository
-- CockroachDB + MinIO Object Store
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_evidence;
SET search_path TO snisid_evidence;

-- ============ EVIDENCE ITEMS ============
CREATE TABLE evidence_items (
    evidence_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    evidence_number VARCHAR(50) UNIQUE NOT NULL,
    evidence_type VARCHAR(30) CHECK (evidence_type IN (
        'PHOTO','VIDEO','AUDIO','CCTV','DRONE_FOOTAGE','PHONE_EXTRACT',
        'COMPUTER_FORENSIC','DOCUMENT_SCAN','SOCIAL_MEDIA','BODY_CAM','DASH_CAM'
    )),
    case_id UUID,
    case_number VARCHAR(50),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    file_hash VARCHAR(64) NOT NULL,
    file_size_bytes BIGINT,
    mime_type VARCHAR(100),
    original_filename VARCHAR(500),
    file_path_minio VARCHAR(500) NOT NULL,
    bucket_name VARCHAR(100) NOT NULL,
    captured_date TIMESTAMPTZ,
    captured_location VARCHAR(500),
    captured_lat DECIMAL(10,7),
    captured_lng DECIMAL(10,7),
    captured_by VARCHAR(255),
    captured_by_agency VARCHAR(100),
    device_id VARCHAR(100),
    device_type VARCHAR(100),
    chain_of_custody JSONB DEFAULT '[]',
    processing_status VARCHAR(30) CHECK (processing_status IN (
        'UPLOADED','PROCESSING','ANALYZED','REVIEWED','ADMITTED','REJECTED'
    )),
    classification VARCHAR(20) CHECK (classification IN (
        'UNCLASSIFIED','RESTRICTED','CONFIDENTIAL','SECRET','TOP_SECRET'
    )),
    retention_date DATE,
    status VARCHAR(20) CHECK (status IN ('ACTIVE','ARCHIVED','PURGED','EXPORTED')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT uq_file_hash UNIQUE (file_hash)
) PARTITION BY LIST (evidence_type);

CREATE INDEX idx_evidence_case ON evidence_items(case_id);
CREATE INDEX idx_evidence_type ON evidence_items(evidence_type);
CREATE INDEX idx_evidence_hash ON evidence_items(file_hash);
CREATE INDEX idx_evidence_date ON evidence_items(captured_date DESC);
CREATE INDEX idx_evidence_agency ON evidence_items(captured_by_agency);
CREATE INDEX idx_evidence_status ON evidence_items(status);
CREATE INDEX idx_evidence_processing ON evidence_items(processing_status);

-- ============ EVIDENCE ANALYSIS ============
CREATE TABLE evidence_analysis (
    analysis_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    evidence_id UUID NOT NULL REFERENCES evidence_items(evidence_id),
    analysis_type VARCHAR(50) CHECK (analysis_type IN (
        'FACE_DETECTION','FACE_RECOGNITION','VOICE_ANALYSIS','OBJECT_DETECTION',
        'OCR','SPEECH_TO_TEXT','MOTION_ANALYSIS','FORENSIC_ANALYSIS',
        'TIMELINE_RECONSTRUCTION','PERSON_REIDENTIFICATION','LICENSE_PLATE'
    )),
    analysis_engine VARCHAR(255),
    analysis_version VARCHAR(50),
    analysis_result JSONB NOT NULL DEFAULT '{}',
    confidence_score DECIMAL(5,2),
    processing_duration_ms INT,
    performed_by VARCHAR(100),
    performed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_analysis_evidence ON evidence_analysis(evidence_id);
CREATE INDEX idx_analysis_type ON evidence_analysis(analysis_type);
CREATE INDEX idx_analysis_time ON evidence_analysis(performed_at DESC);

-- ============ CHAIN OF CUSTODY ============
CREATE TABLE evidence_custody (
    custody_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    evidence_id UUID NOT NULL REFERENCES evidence_items(evidence_id),
    action VARCHAR(50) NOT NULL,
    actor_id VARCHAR(100) NOT NULL,
    actor_name VARCHAR(255) NOT NULL,
    actor_agency VARCHAR(100) NOT NULL,
    from_location VARCHAR(500),
    to_location VARCHAR(500),
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    digital_signature TEXT,
    notes TEXT
);

CREATE INDEX idx_custody_evidence ON evidence_custody(evidence_id);
CREATE INDEX idx_custody_time ON evidence_custody(timestamp);
CREATE INDEX idx_custody_actor ON evidence_custody(actor_id);

-- ============ MULTIMODAL SEARCH INDEX ============
CREATE TABLE evidence_multimodal_index (
    index_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    evidence_id UUID NOT NULL REFERENCES evidence_items(evidence_id),
    modality VARCHAR(20) CHECK (modality IN ('FACE','VOICE','OBJECT','TEXT','SCENE')),
    embedding_id VARCHAR(100) NOT NULL,
    embedding_dim INT,
    algorithm VARCHAR(100),
    vector_db_collection VARCHAR(100),
    indexed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_mm_evidence ON evidence_multimodal_index(evidence_id);
CREATE INDEX idx_mm_modality ON evidence_multimodal_index(modality);

ALTER TABLE evidence_items ENABLE ROW LEVEL SECURITY;
CREATE POLICY evidence_pnh_select ON evidence_items FOR SELECT USING (
    current_setting('snisid.agency') IN ('PNH','DCPJ','SOC','SNISID_ADMIN')
);
