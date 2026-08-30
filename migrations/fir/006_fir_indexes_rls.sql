BEGIN;

CREATE INDEX IF NOT EXISTS idx_fir_records_snisid    ON fir_criminal_records(snisid_person_id);
CREATE INDEX IF NOT EXISTS idx_fir_records_fir_id    ON fir_criminal_records(national_fir_id);
CREATE INDEX IF NOT EXISTS idx_fir_charges_record    ON fir_charges(record_id);
CREATE INDEX IF NOT EXISTS idx_fir_charges_arrest_dt ON fir_charges(arrest_date DESC);
CREATE INDEX IF NOT EXISTS idx_fir_charges_verdict_dt ON fir_charges(verdict_date DESC);
CREATE INDEX IF NOT EXISTS idx_fir_charges_dept      ON fir_charges(court_dept);
CREATE INDEX IF NOT EXISTS idx_fir_aliases_record    ON fir_aliases(record_id);
CREATE INDEX IF NOT EXISTS idx_fir_movements_record  ON fir_movements(record_id);
CREATE INDEX IF NOT EXISTS idx_fir_movements_type    ON fir_movements(movement_type);

ALTER TABLE fir_criminal_records ENABLE ROW LEVEL SECURITY;
ALTER TABLE fir_charges ENABLE ROW LEVEL SECURITY;
ALTER TABLE fir_aliases ENABLE ROW LEVEL SECURITY;

CREATE POLICY fir_records_select_policy ON fir_criminal_records
    FOR SELECT USING (
        current_setting('app.user_role', TRUE) IN
        ('DCPJ','PARQUET','TRIBUNAL','JUDGE','POLICE_OFFICER','SUPERADMIN')
    );

CREATE POLICY fir_records_all_policy ON fir_criminal_records
    USING (current_setting('app.user_role', TRUE) = 'SUPERADMIN')
    WITH CHECK (current_setting('app.user_role', TRUE) = 'SUPERADMIN');

CREATE POLICY fir_charges_select_policy ON fir_charges
    FOR SELECT USING (
        current_setting('app.user_role', TRUE) IN
        ('DCPJ','PARQUET','TRIBUNAL','JUDGE','POLICE_OFFICER','SUPERADMIN')
    );

CREATE POLICY fir_aliases_select_policy ON fir_aliases
    FOR SELECT USING (
        current_setting('app.user_role', TRUE) IN
        ('DCPJ','PARQUET','TRIBUNAL','JUDGE','POLICE_OFFICER','SUPERADMIN')
    );

COMMIT;
