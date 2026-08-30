BEGIN;

CREATE TABLE sivc_criminal_alerts (
    alert_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plate_number        VARCHAR(20) NOT NULL,
    plate_category      sivc_plate_category,
    vin                 VARCHAR(17),
    chassis_number      VARCHAR(50),
    vehicle_type        sivc_vehicle_type,
    make                VARCHAR(100),
    model               VARCHAR(100),
    year                SMALLINT,
    color_primary       VARCHAR(50),
    color_secondary     VARCHAR(50),
    distinguishing_marks TEXT,
    foves_vehicle_id    UUID,
    crime_category      sivc_crime_category NOT NULL,
    crime_subcategory   VARCHAR(100),
    alert_level         sivc_alert_level NOT NULL DEFAULT 'CAUTION',
    status              sivc_alert_status NOT NULL DEFAULT 'ACTIVE',
    armed_and_dangerous BOOLEAN DEFAULT FALSE,
    do_not_stop_alone   BOOLEAN DEFAULT FALSE,
    officer_safety_notes TEXT,
    reporting_unit      VARCHAR(50) NOT NULL,
    reporting_officer_id UUID,
    incident_reference  VARCHAR(100),
    incident_date       TIMESTAMPTZ NOT NULL,
    expiry_date         TIMESTAMPTZ,
    associated_person_ids   UUID[] DEFAULT '{}',
    associated_case_ids     UUID[] DEFAULT '{}',
    associated_alert_ids    UUID[] DEFAULT '{}',
    interpol_smv_id         VARCHAR(50),
    interpol_reported       BOOLEAN DEFAULT FALSE,
    interpol_reported_at    TIMESTAMPTZ,
    last_seen_lat           DECIMAL(10,7),
    last_seen_lng           DECIMAL(10,7),
    last_seen_location      VARCHAR(300),
    last_seen_dept_code     CHAR(2),
    last_seen_commune       VARCHAR(100),
    last_seen_at            TIMESTAMPTZ,
    photo_refs          TEXT[] DEFAULT '{}',
    document_refs       TEXT[] DEFAULT '{}',
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by          UUID,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version             INTEGER NOT NULL DEFAULT 1,
    CONSTRAINT chk_plate_format CHECK (plate_number ~ '^[A-Z]{1,3}[-\s]?[0-9]{3,6}[A-Z]?$')
);

CREATE INDEX idx_sivc_alerts_plate    ON sivc_criminal_alerts(plate_number) WHERE status = 'ACTIVE';
CREATE INDEX idx_sivc_alerts_vin      ON sivc_criminal_alerts(vin) WHERE vin IS NOT NULL AND status = 'ACTIVE';
CREATE INDEX idx_sivc_alerts_category ON sivc_criminal_alerts(crime_category, status);
CREATE INDEX idx_sivc_alerts_dept     ON sivc_criminal_alerts(last_seen_dept_code) WHERE status = 'ACTIVE';
CREATE INDEX idx_sivc_alerts_level    ON sivc_criminal_alerts(alert_level) WHERE status = 'ACTIVE';
CREATE INDEX idx_sivc_alerts_persons  ON sivc_criminal_alerts USING gin(associated_person_ids);

COMMIT;
