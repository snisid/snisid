BEGIN;

CREATE TABLE sivc_kidnapping_incidents (
    incident_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_id                UUID NOT NULL REFERENCES sivc_criminal_alerts(alert_id) ON DELETE CASCADE,
    victim_count            SMALLINT NOT NULL DEFAULT 1,
    victim_snisid_ids       UUID[] DEFAULT '{}',
    victims_nationality     CHAR(3)[] DEFAULT '{}',
    victims_description     TEXT,
    abduction_date          TIMESTAMPTZ NOT NULL,
    abduction_location      VARCHAR(300),
    abduction_dept_code     CHAR(2),
    abduction_commune       VARCHAR(100),
    abduction_context       TEXT,
    ransom_demanded         BOOLEAN DEFAULT FALSE,
    ransom_amount           DECIMAL(15,2),
    ransom_currency         CHAR(3) DEFAULT 'USD',
    ransom_channel          VARCHAR(100),
    incident_status         sivc_kidnapping_status NOT NULL DEFAULT 'IN_PROGRESS',
    resolution_date         TIMESTAMPTZ,
    resolution_location     VARCHAR(300),
    resolution_notes        TEXT,
    cae_case_number         VARCHAR(50),
    dcpj_case_number        VARCHAR(50),
    created_by              UUID NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sivc_kidnap_alert  ON sivc_kidnapping_incidents(alert_id);
CREATE INDEX idx_sivc_kidnap_status ON sivc_kidnapping_incidents(incident_status);
CREATE INDEX idx_sivc_kidnap_dept   ON sivc_kidnapping_incidents(abduction_dept_code);
CREATE INDEX idx_sivc_kidnap_date   ON sivc_kidnapping_incidents(abduction_date DESC);

COMMIT;
