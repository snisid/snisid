BEGIN;

CREATE TABLE sltd_documents (
    doc_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_sltd_id    VARCHAR(25) UNIQUE NOT NULL,
    doc_type            sltd_doc_type NOT NULL,
    document_number     VARCHAR(100) NOT NULL,
    issuing_country     CHAR(3) NOT NULL DEFAULT 'HTI',
    holder_name         VARCHAR(200),
    holder_snisid_id    UUID,
    holder_dob          DATE,
    holder_nationality  CHAR(3) DEFAULT 'HTI',
    issue_date          DATE,
    expiry_date         DATE,
    status              sltd_doc_status NOT NULL,
    reported_date       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reported_by         UUID NOT NULL,
    reporting_dept_code CHAR(2),
    theft_context       TEXT,
    found_date          TIMESTAMPTZ,
    found_location      VARCHAR(300),
    interpol_sltd_ref   VARCHAR(50),
    reported_to_interpol BOOLEAN DEFAULT FALSE,
    interpol_reported_at TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
