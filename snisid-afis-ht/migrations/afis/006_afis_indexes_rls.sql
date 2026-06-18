BEGIN;

ALTER TABLE afis_subjects ENABLE ROW LEVEL SECURITY;
ALTER TABLE afis_fingerprints ENABLE ROW LEVEL SECURITY;
ALTER TABLE afis_latent_prints ENABLE ROW LEVEL SECURITY;
ALTER TABLE afis_search_transactions ENABLE ROW LEVEL SECURITY;

CREATE POLICY afis_subjects_read_policy ON afis_subjects
    FOR SELECT TO pnh_officer, lab_forensic, dcpj_supervisor, parquet_proc
    USING (true);

CREATE POLICY afis_subjects_write_policy ON afis_subjects
    FOR INSERT TO pnh_officer, lab_forensic, dcpj_supervisor
    WITH CHECK (true);

CREATE POLICY afis_prints_read_policy ON afis_fingerprints
    FOR SELECT TO pnh_officer, lab_forensic, dcpj_supervisor, parquet_proc
    USING (true);

CREATE POLICY afis_prints_write_policy ON afis_fingerprints
    FOR INSERT TO pnh_officer, lab_forensic, dcpj_supervisor
    WITH CHECK (true);

CREATE POLICY afis_latent_read_policy ON afis_latent_prints
    FOR SELECT TO lab_forensic, dcpj_supervisor
    USING (true);

CREATE POLICY afis_latent_write_policy ON afis_latent_prints
    FOR INSERT TO lab_forensic, dcpj_supervisor
    WITH CHECK (true);

CREATE POLICY afis_search_read_policy ON afis_search_transactions
    FOR SELECT TO pnh_officer, lab_forensic, dcpj_supervisor
    USING (true);

CREATE POLICY afis_search_write_policy ON afis_search_transactions
    FOR INSERT TO pnh_officer, lab_forensic, dcpj_supervisor
    WITH CHECK (true);

CREATE OR REPLACE FUNCTION afis_update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER afis_subjects_updated_at
    BEFORE UPDATE ON afis_subjects
    FOR EACH ROW EXECUTE FUNCTION afis_update_updated_at();

COMMIT;