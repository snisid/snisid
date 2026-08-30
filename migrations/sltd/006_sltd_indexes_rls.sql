BEGIN;

CREATE UNIQUE INDEX idx_sltd_doc_number ON sltd_documents(document_number, issuing_country)
    WHERE status IN ('LOST','STOLEN','REVOKED');
CREATE INDEX idx_sltd_holder    ON sltd_documents(holder_snisid_id) WHERE holder_snisid_id IS NOT NULL;
CREATE INDEX idx_sltd_status    ON sltd_documents(status);
CREATE INDEX idx_sltd_check_log ON sltd_check_log(document_number, checked_at DESC);

ALTER TABLE sltd_documents ENABLE ROW LEVEL SECURITY;
ALTER TABLE sltd_check_log ENABLE ROW LEVEL SECURITY;

CREATE POLICY sltd_documents_select ON sltd_documents
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('DGMN','PNH_OFFICER','SIFR','AERO')
    );

CREATE POLICY sltd_documents_insert ON sltd_documents
    FOR INSERT WITH CHECK (
        current_setting('snisid.user_role') IN ('PUBLIC','PNH_OFFICER','DGMN')
    );

CREATE POLICY sltd_documents_update ON sltd_documents
    FOR UPDATE USING (
        current_setting('snisid.user_role') IN ('PNH_OFFICER','DGMN')
    );

COMMIT;
