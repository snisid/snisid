BEGIN;

CREATE INDEX idx_mar_ais_timestamp   ON mar_ais_sightings(sighting_timestamp DESC);
CREATE INDEX idx_mar_ais_mmsi        ON mar_ais_sightings(mmsi);
CREATE INDEX idx_mar_ais_coords      ON mar_ais_sightings(lat, lng);
CREATE INDEX idx_mar_incidents_date  ON mar_incidents(incident_date DESC);
CREATE INDEX idx_mar_incidents_type  ON mar_incidents(incident_type);
CREATE INDEX idx_mar_watch_active    ON mar_watch_vessels(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_mar_vessels_status  ON mar_vessels(status);

ALTER TABLE mar_vessels ENABLE ROW LEVEL SECURITY;
ALTER TABLE mar_ais_sightings ENABLE ROW LEVEL SECURITY;
ALTER TABLE mar_incidents ENABLE ROW LEVEL SECURITY;
ALTER TABLE mar_watch_vessels ENABLE ROW LEVEL SECURITY;

CREATE POLICY mar_vessels_select ON mar_vessels
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('GCH','DCPJ','GCH_ADMIN')
    );

CREATE POLICY mar_vessels_insert ON mar_vessels
    FOR INSERT WITH CHECK (
        current_setting('snisid.user_role') = 'GCH_ADMIN'
    );

CREATE POLICY mar_incidents_select ON mar_incidents
    FOR SELECT USING (
        current_setting('snisid.user_role') IN ('GCH','DCPJ','GCH_OFFICER')
    );

CREATE POLICY mar_incidents_insert ON mar_incidents
    FOR INSERT WITH CHECK (
        current_setting('snisid.user_role') IN ('GCH','GCH_OFFICER')
    );

COMMIT;
