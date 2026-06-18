BEGIN;

CREATE INDEX idx_gang_dept       ON gang_organizations(primary_dept_code) WHERE is_active = TRUE;
CREATE INDEX idx_gang_activity   ON gang_organizations(activity_level) WHERE is_active = TRUE;
CREATE INDEX idx_gang_ofac       ON gang_organizations(ofac_designation) WHERE ofac_designation = TRUE;
CREATE INDEX idx_gang_incidents  ON gang_incidents(gang_id, incident_date DESC);
CREATE INDEX idx_gang_incid_dept ON gang_incidents(dept_code, incident_date DESC);
CREATE INDEX idx_gang_members_gang ON gang_members(gang_id);
CREATE INDEX idx_gang_territory_gang ON gang_territories(gang_id);
CREATE INDEX idx_gang_territory_dept ON gang_territories(dept_code);

ALTER TABLE gang_organizations ENABLE ROW LEVEL SECURITY;
ALTER TABLE gang_members        ENABLE ROW LEVEL SECURITY;
ALTER TABLE gang_incidents      ENABLE ROW LEVEL SECURITY;
ALTER TABLE gang_territories    ENABLE ROW LEVEL SECURITY;

CREATE POLICY gang_org_select ON gang_organizations FOR SELECT USING (
    current_setting('app.user_role') IN ('DCPJ_INTEL','DCPJ_BAC','DCPJ','BRI','MJSP')
);

CREATE POLICY gang_org_insert ON gang_organizations FOR INSERT WITH CHECK (
    current_setting('app.user_role') = 'DCPJ_INTEL'
);

CREATE POLICY gang_members_select ON gang_members FOR SELECT USING (
    current_setting('app.user_role') IN ('DCPJ_INTEL','DCPJ_BAC','DCPJ','BRI')
);

CREATE POLICY gang_incidents_select ON gang_incidents FOR SELECT USING (
    current_setting('app.user_role') IN ('DCPJ_INTEL','DCPJ_BAC','DCPJ','BRI')
);

CREATE POLICY gang_territories_select ON gang_territories FOR SELECT USING (
    current_setting('app.user_role') IN ('DCPJ_INTEL','DCPJ_BAC','DCPJ','BRI','MJSP')
);

COMMIT;
