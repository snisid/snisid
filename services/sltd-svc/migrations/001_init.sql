BEGIN;
CREATE TYPE sltd_doc_type AS ENUM ('PASSPORT','NATIONAL_ID','TRAVEL_DOCUMENT','VISA','RESIDENCE_PERMIT','REFUGEE_DOCUMENT','LAISSEZ_PASSER');
CREATE TYPE sltd_doc_status AS ENUM ('LOST','STOLEN','REVOKED','EXPIRED','FOUND','RECOVERED','CANCELLED');
CREATE TABLE sltd_documents (
    doc_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), national_sltd_id VARCHAR(25) UNIQUE NOT NULL,
    doc_type sltd_doc_type NOT NULL, document_number VARCHAR(100) NOT NULL,
    issuing_country CHAR(3) NOT NULL DEFAULT 'HTI', holder_name VARCHAR(200),
    holder_snisid_id UUID, holder_dob DATE, holder_nationality CHAR(3) DEFAULT 'HTI',
    issue_date DATE, expiry_date DATE, status sltd_doc_status NOT NULL,
    reported_date TIMESTAMPTZ NOT NULL DEFAULT NOW(), reported_by UUID NOT NULL,
    reporting_dept_code CHAR(2), theft_context TEXT, found_date TIMESTAMPTZ,
    found_location VARCHAR(300), interpol_sltd_ref VARCHAR(50),
    reported_to_interpol BOOLEAN DEFAULT FALSE, interpol_reported_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE sltd_check_log (
    check_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), document_number VARCHAR(100) NOT NULL,
    doc_type sltd_doc_type, checked_by UUID NOT NULL, check_location VARCHAR(100),
    post_id UUID, result VARCHAR(20) NOT NULL, source VARCHAR(20) NOT NULL,
    sltd_doc_id UUID REFERENCES sltd_documents(doc_id),
    checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_sltd_doc_number ON sltd_documents(document_number, issuing_country) WHERE status IN ('LOST','STOLEN','REVOKED');
CREATE INDEX idx_sltd_holder ON sltd_documents(holder_snisid_id) WHERE holder_snisid_id IS NOT NULL;
CREATE INDEX idx_sltd_status ON sltd_documents(status);
CREATE INDEX idx_sltd_check_log ON sltd_check_log(document_number, checked_at DESC);
COMMIT;
