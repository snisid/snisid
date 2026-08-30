BEGIN;

CREATE INDEX idx_expl_type  ON expl_incidents(explosive_type, incident_date DESC);
CREATE INDEX idx_expl_dept  ON expl_incidents(dept_code);
CREATE INDEX idx_expl_gang  ON expl_incidents(gang_id) WHERE gang_id IS NOT NULL;

ALTER TABLE expl_incidents ENABLE ROW LEVEL SECURITY;
ALTER TABLE expl_legal_stocks ENABLE ROW LEVEL SECURITY;

CREATE POLICY expl_incidents_select ON expl_incidents
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('DCPJ','PNH_OFFICER','PNH_ADMIN','FAd_H')
    );

CREATE POLICY expl_incidents_insert ON expl_incidents
    FOR INSERT WITH CHECK (
        current_setting('snisid.user_role') IN ('DCPJ','PNH_OFFICER','FAd_H')
    );

CREATE POLICY expl_stocks_select ON expl_legal_stocks
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('DCPJ','DCPJ_ADMIN','PNH_ADMIN')
    );

CREATE POLICY expl_stocks_insert ON expl_legal_stocks
    FOR INSERT WITH CHECK (
        current_setting('snisid.user_role') = 'DCPJ_ADMIN'
    );

COMMIT;
