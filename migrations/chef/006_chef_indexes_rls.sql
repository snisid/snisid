BEGIN;

CREATE INDEX IF NOT EXISTS idx_chef_gang        ON chef_criminal_members(primary_gang_id, role_in_gang);
CREATE INDEX IF NOT EXISTS idx_chef_status      ON chef_criminal_members(status) WHERE status = 'ACTIVE';
CREATE INDEX IF NOT EXISTS idx_chef_ofac        ON chef_criminal_members(ofac_designated) WHERE ofac_designated = TRUE;
CREATE INDEX IF NOT EXISTS idx_chef_un          ON chef_criminal_members(un_designated) WHERE un_designated = TRUE;
CREATE INDEX IF NOT EXISTS idx_chef_dept        ON chef_criminal_members(territory_dept) WHERE status = 'ACTIVE';
CREATE INDEX IF NOT EXISTS idx_chef_sightings   ON chef_sightings(member_id, sighted_at DESC);

ALTER TABLE chef_criminal_members    ENABLE ROW LEVEL SECURITY;
ALTER TABLE chef_intelligence_notes  ENABLE ROW LEVEL SECURITY;
ALTER TABLE chef_cross_gang_links    ENABLE ROW LEVEL SECURITY;
ALTER TABLE chef_sightings           ENABLE ROW LEVEL SECURITY;

CREATE POLICY chef_members_policy ON chef_criminal_members
    USING (true)
    WITH CHECK (true);

CREATE POLICY chef_intel_notes_policy ON chef_intelligence_notes
    USING (true)
    WITH CHECK (true);

CREATE POLICY chef_links_policy ON chef_cross_gang_links
    USING (true)
    WITH CHECK (true);

CREATE POLICY chef_sightings_policy ON chef_sightings
    USING (true)
    WITH CHECK (true);

COMMIT;
