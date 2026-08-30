BEGIN;

ALTER TABLE sivc_criminal_alerts ENABLE ROW LEVEL SECURITY;
ALTER TABLE sivc_stolen_plates   ENABLE ROW LEVEL SECURITY;
ALTER TABLE sivc_vehicle_sightings ENABLE ROW LEVEL SECURITY;

CREATE POLICY sivc_alerts_read_policy ON sivc_criminal_alerts
    FOR SELECT
    USING (
        current_setting('app.user_role', TRUE) IN
        ('BLVV', 'BLTS', 'BAC', 'DCPJ', 'CAE', 'BRI', 'GIPNH', 'SUPERADMIN')
    );

CREATE POLICY sivc_alerts_write_policy ON sivc_criminal_alerts
    FOR ALL
    USING (
        reporting_unit = current_setting('app.user_unit', TRUE)
        OR current_setting('app.user_role', TRUE) IN ('DCPJ', 'SUPERADMIN')
    );

CREATE OR REPLACE FUNCTION sivc_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    NEW.version = OLD.version + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_sivc_alerts_updated
    BEFORE UPDATE ON sivc_criminal_alerts
    FOR EACH ROW EXECUTE FUNCTION sivc_update_timestamp();

CREATE INDEX idx_sivc_alerts_fts ON sivc_criminal_alerts
    USING gin(to_tsvector('french',
        COALESCE(plate_number, '') || ' ' ||
        COALESCE(make, '') || ' ' ||
        COALESCE(model, '') || ' ' ||
        COALESCE(color_primary, '') || ' ' ||
        COALESCE(officer_safety_notes, '')
    ));

CREATE INDEX idx_sivc_alerts_plate_active
    ON sivc_criminal_alerts(plate_number, alert_level, crime_category)
    WHERE status = 'ACTIVE';

CREATE INDEX idx_sivc_alerts_wanted_dept
    ON sivc_criminal_alerts(last_seen_dept_code, alert_level)
    WHERE status = 'ACTIVE' AND alert_level IN ('WANTED', 'CRITICAL');

COMMIT;
