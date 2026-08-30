BEGIN;

CREATE TABLE sivc_drug_incidents (
    incident_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_id            UUID NOT NULL REFERENCES sivc_criminal_alerts(alert_id) ON DELETE CASCADE,
    drug_types          TEXT[] NOT NULL,
    seizure_weight_kg   DECIMAL(12,4),
    estimated_value_usd DECIMAL(15,2),
    seizure_date        TIMESTAMPTZ,
    seizure_location    VARCHAR(300),
    seizure_dept_code   CHAR(2),
    seizure_commune     VARCHAR(100),
    route_type          sivc_route_type,
    origin_country      CHAR(3),
    transit_points      TEXT[] DEFAULT '{}',
    destination         VARCHAR(200),
    suspected_cartel    VARCHAR(200),
    blts_case_number    VARCHAR(50),
    interpol_ref        VARCHAR(50),
    concealment_method  TEXT,
    notes               TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sivc_drug_alert   ON sivc_drug_incidents(alert_id);
CREATE INDEX idx_sivc_drug_dept    ON sivc_drug_incidents(seizure_dept_code);
CREATE INDEX idx_sivc_drug_types   ON sivc_drug_incidents USING gin(drug_types);

COMMIT;
