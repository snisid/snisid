BEGIN;

CREATE TABLE sivc_stolen_plates (
    plate_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plate_number        VARCHAR(20) NOT NULL,
    plate_category      sivc_plate_category NOT NULL,
    original_vehicle_id UUID,
    original_make       VARCHAR(100),
    original_model      VARCHAR(100),
    original_vin        VARCHAR(17),
    theft_date          TIMESTAMPTZ NOT NULL,
    theft_location      VARCHAR(300),
    theft_dept_code     CHAR(2) NOT NULL,
    theft_commune       VARCHAR(100),
    theft_context       TEXT,
    reporting_unit      VARCHAR(50) NOT NULL,
    reporting_officer_id UUID,
    blvv_case_number    VARCHAR(50),
    status              sivc_stolen_plate_status NOT NULL DEFAULT 'STOLEN',
    recovered_date      TIMESTAMPTZ,
    recovery_location   VARCHAR(300),
    recovery_dept_code  CHAR(2),
    used_in_crime       BOOLEAN DEFAULT FALSE,
    crime_categories    sivc_crime_category[] DEFAULT '{}',
    crime_alert_ids     UUID[] DEFAULT '{}',
    is_state_plate_clone BOOLEAN DEFAULT FALSE,
    impersonated_agency  VARCHAR(100),
    interpol_sad_id     VARCHAR(50),
    interpol_reported   BOOLEAN DEFAULT FALSE,
    notes               TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_sivc_stolen_plates_unique
    ON sivc_stolen_plates(plate_number)
    WHERE status = 'STOLEN';

CREATE INDEX idx_sivc_stolen_plates_category
    ON sivc_stolen_plates(plate_category, status);

CREATE INDEX idx_sivc_stolen_plates_state
    ON sivc_stolen_plates(is_state_plate_clone)
    WHERE is_state_plate_clone = TRUE AND status = 'STOLEN';

COMMIT;
