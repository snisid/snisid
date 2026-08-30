BEGIN;

CREATE TABLE sivc_vehicle_sightings (
    sighting_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plate_number        VARCHAR(20) NOT NULL,
    source_type         VARCHAR(20) NOT NULL DEFAULT 'LAPI',
    lapi_unit_id        VARCHAR(100),
    reporting_agent_id  UUID,
    sighting_timestamp  TIMESTAMPTZ NOT NULL,
    location_lat        DECIMAL(10,7),
    location_lng        DECIMAL(10,7),
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    checkpoint_name     VARCHAR(150),
    matched_alert_id    UUID REFERENCES sivc_criminal_alerts(alert_id),
    matched_plate_id    UUID REFERENCES sivc_stolen_plates(plate_id),
    match_confidence    DECIMAL(5,2),
    alert_triggered     BOOLEAN DEFAULT FALSE,
    alert_level         sivc_alert_level,
    alert_sent_at       TIMESTAMPTZ,
    alert_recipients    TEXT[] DEFAULT '{}',
    image_ref           VARCHAR(500),
    video_clip_ref      VARCHAR(500),
    is_reviewed         BOOLEAN DEFAULT FALSE,
    reviewed_by         UUID,
    reviewed_at         TIMESTAMPTZ,
    review_notes        TEXT,
    false_positive      BOOLEAN DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sivc_sightings_plate     ON sivc_vehicle_sightings(plate_number);
CREATE INDEX idx_sivc_sightings_timestamp ON sivc_vehicle_sightings(sighting_timestamp DESC);
CREATE INDEX idx_sivc_sightings_dept      ON sivc_vehicle_sightings(dept_code, sighting_timestamp DESC);
CREATE INDEX idx_sivc_sightings_alert     ON sivc_vehicle_sightings(matched_alert_id) WHERE matched_alert_id IS NOT NULL;
CREATE INDEX idx_sivc_sightings_triggered ON sivc_vehicle_sightings(alert_triggered) WHERE alert_triggered = TRUE;

COMMIT;
