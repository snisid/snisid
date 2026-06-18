-- ============================================================
-- SNI-SIDE: Document Fraud Database
-- PostgreSQL 16
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_doc_fraud;
SET search_path TO snisid_doc_fraud;

-- ============ DOCUMENTS ============
CREATE TABLE fraud_documents (
    document_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_type VARCHAR(30) CHECK (document_type IN (
        'PASSPORT','CIN','DRIVER_LICENSE','BIRTH_CERTIFICATE','MARRIAGE_CERTIFICATE',
        'DEATH_CERTIFICATE','RESIDENCE_PERMIT','VISA','REFUGEE_CARD','DIPLOMATIC_PASSPORT'
    )),
    document_number VARCHAR(100) NOT NULL,
    document_number_hash VARCHAR(64) GENERATED ALWAYS AS (encode(sha256(document_number::bytea), 'hex')) STORED,
    issuing_country VARCHAR(100),
    issuing_authority VARCHAR(255),
    issue_date DATE,
    expiry_date DATE,
    holder_niu VARCHAR(10),
    holder_name VARCHAR(255) NOT NULL,
    holder_dob DATE,
    holder_gender VARCHAR(1),
    holder_nationality VARCHAR(100),
    mrz_line_1 VARCHAR(100),
    mrz_line_2 VARCHAR(100),
    mrz_line_3 VARCHAR(100),
    mrz_hash VARCHAR(64),
    document_image_front TEXT,
    document_image_back TEXT,
    chip_data JSONB DEFAULT '{}',
    biometric_chip_present BOOLEAN DEFAULT FALSE,
    security_features_verified JSONB DEFAULT '{}',
    status VARCHAR(20) CHECK (status IN ('VALID','EXPIRED','REVOKED','REPORTED_STOLEN','REPORTED_LOST','SUSPENDED','FRAUDULENT')),
    risk_score DECIMAL(5,2),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_citizen FOREIGN KEY (holder_niu) REFERENCES snisid_identity.citizens(niu)
);

CREATE INDEX idx_fdoc_number ON fraud_documents(document_number);
CREATE INDEX idx_fdoc_type ON fraud_documents(document_type);
CREATE INDEX idx_fdoc_holder ON fraud_documents(holder_niu);
CREATE INDEX idx_fdoc_holder_name ON fraud_documents USING gin(to_tsvector('french', holder_name));
CREATE INDEX idx_fdoc_status ON fraud_documents(status);
CREATE INDEX idx_fdoc_risk ON fraud_documents(risk_score DESC);
CREATE INDEX idx_fdoc_mrz ON fraud_documents(mrz_hash);
CREATE INDEX idx_fdoc_country ON fraud_documents(issuing_country);
CREATE INDEX idx_fdoc_expiry ON fraud_documents(expiry_date);

-- ============ FRAUD REPORTS ============
CREATE TABLE document_fraud_reports (
    fraud_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES fraud_documents(document_id),
    fraud_type VARCHAR(50) CHECK (fraud_type IN (
        'FORGERY','ALTERATION','COUNTERFEIT','IMPERSONATION','STOLEN_BLANK',
        'FALSE_IDENTITY','QR_CODE_TAMPERING','MRZ_MANIPULATION','CHIP_CLONING',
        'PHOTO_SUBSTITUTION','DATA_PAGE_REPLACEMENT','VISA_FORGERY','STAMP_FORGERY'
    )),
    detection_date TIMESTAMPTZ NOT NULL,
    detection_location VARCHAR(500),
    detection_method VARCHAR(100) CHECK (detection_method IN (
        'VISUAL_INSPECTION','UV_LIGHT','MRZ_SCANNER','QR_SCANNER','BIOMETRIC_VERIFICATION',
        'CHIP_VERIFICATION','DATABASE_CROSSCHECK','AI_ANALYSIS','TIP_REFERRAL'
    )),
    detected_by_agency VARCHAR(100),
    detected_by_officer VARCHAR(255),
    confidence_score DECIMAL(5,2),
    fraudulent_elements TEXT[],
    images JSONB DEFAULT '[]',
    analysis_notes TEXT,
    case_id UUID,
    status VARCHAR(20) CHECK (status IN ('PENDING_VERIFICATION','CONFIRMED','FALSE_POSITIVE','PROSECUTED')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_fraud_document ON document_fraud_reports(document_id);
CREATE INDEX idx_fraud_type ON document_fraud_reports(fraud_type);
CREATE INDEX idx_fraud_method ON document_fraud_reports(detection_method);
CREATE INDEX idx_fraud_status ON document_fraud_reports(status);
CREATE INDEX idx_fraud_date ON document_fraud_reports(detection_date DESC);

-- ============ STOLEN DOCUMENTS ============
CREATE TABLE stolen_documents (
    stolen_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES fraud_documents(document_id),
    report_date TIMESTAMPTZ NOT NULL,
    report_location VARCHAR(500),
    reported_by VARCHAR(255),
    reporting_agency VARCHAR(100),
    theft_circumstances TEXT,
    police_report_ref VARCHAR(100),
    status VARCHAR(20) CHECK (status IN ('STOLEN','RECOVERED','CANCELLED')),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE fraud_documents ENABLE ROW LEVEL SECURITY;
CREATE POLICY doc_fraud_pnh_select ON fraud_documents FOR SELECT USING (
    current_setting('snisid.agency') IN ('PNH','DCPJ','IMMIGRATION','ONI','SNISID_ADMIN')
);
