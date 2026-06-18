BEGIN;

CREATE TABLE sltd_check_log (
    check_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_number     VARCHAR(100) NOT NULL,
    doc_type            sltd_doc_type,
    checked_by          UUID NOT NULL,
    check_location      VARCHAR(100),
    post_id             UUID,
    result              VARCHAR(20) NOT NULL,
    source              VARCHAR(20) NOT NULL,
    sltd_doc_id         UUID REFERENCES sltd_documents(doc_id),
    checked_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
