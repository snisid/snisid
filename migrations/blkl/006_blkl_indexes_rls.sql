BEGIN;

CREATE INDEX idx_blkl_person   ON blkl_blacklist(snisid_person_id) WHERE is_active = TRUE;
CREATE INDEX idx_blkl_type     ON blkl_blacklist(restriction_type) WHERE is_active = TRUE;
CREATE INDEX idx_blkl_expiry   ON blkl_blacklist(expiry_date) WHERE is_active = TRUE AND expiry_date IS NOT NULL;

ALTER TABLE blkl_blacklist ENABLE ROW LEVEL SECURITY;
ALTER TABLE blkl_alerts_log ENABLE ROW LEVEL SECURITY;

CREATE POLICY blkl_blacklist_select ON blkl_blacklist
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('DGMN','DCPJ','MJSP','TRIBUNAL','SIFR','AERO')
    );

CREATE POLICY blkl_blacklist_insert ON blkl_blacklist
    FOR INSERT WITH CHECK (
        current_setting('snisid.user_role') IN ('MJSP','TRIBUNAL','DCPJ')
    );

CREATE POLICY blkl_blacklist_update ON blkl_blacklist
    FOR UPDATE USING (
        current_setting('snisid.user_role') IN ('MJSP_ADMIN')
    );

COMMIT;
