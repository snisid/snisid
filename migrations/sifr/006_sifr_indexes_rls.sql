BEGIN;

CREATE INDEX idx_sifr_crossings_datetime ON sifr_crossings(crossing_datetime DESC);
CREATE INDEX idx_sifr_crossings_person   ON sifr_crossings(snisid_person_id) WHERE snisid_person_id IS NOT NULL;
CREATE INDEX idx_sifr_crossings_doc      ON sifr_crossings(document_number);
CREATE INDEX idx_sifr_crossings_alert    ON sifr_crossings(alert_triggered) WHERE alert_triggered = TRUE;
CREATE INDEX idx_sifr_crossings_post     ON sifr_crossings(post_id, crossing_datetime DESC);

ALTER TABLE sifr_border_posts ENABLE ROW LEVEL SECURITY;
ALTER TABLE sifr_crossings ENABLE ROW LEVEL SECURITY;
ALTER TABLE sifr_alerts_log ENABLE ROW LEVEL SECURITY;
ALTER TABLE sifr_clandestine_crossings ENABLE ROW LEVEL SECURITY;

CREATE POLICY sifr_posts_select ON sifr_border_posts
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('DGMN','PNH_OFFICER','DCPJ','PUBLIC_DGMN')
    );

CREATE POLICY sifr_crossings_select ON sifr_crossings
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('SIFR_AGENT','SIFR_SUPERVISOR','DGMN','DCPJ')
    );

CREATE POLICY sifr_crossings_insert ON sifr_crossings
    FOR INSERT WITH CHECK (
        current_setting('snisid.user_role') IN ('SIFR_AGENT','SIFR_SUPERVISOR')
    );

COMMIT;
