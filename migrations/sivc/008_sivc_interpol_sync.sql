BEGIN;

CREATE TABLE sivc_interpol_sync_log (
    sync_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alert_id            UUID REFERENCES sivc_criminal_alerts(alert_id),
    stolen_plate_id     UUID REFERENCES sivc_stolen_plates(plate_id),
    interpol_smv_id     VARCHAR(50),
    interpol_sad_id     VARCHAR(50),
    sync_direction      sivc_sync_direction NOT NULL,
    sync_status         sivc_sync_status NOT NULL DEFAULT 'PENDING',
    sync_timestamp      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    retry_count         SMALLINT DEFAULT 0,
    request_payload     JSONB,
    response_payload    JSONB,
    error_code          VARCHAR(50),
    error_message       TEXT,
    processed_by        UUID,
    processed_at        TIMESTAMPTZ
);

CREATE INDEX idx_sivc_interpol_status ON sivc_interpol_sync_log(sync_status, sync_direction);
CREATE INDEX idx_sivc_interpol_smv_id ON sivc_interpol_sync_log(interpol_smv_id) WHERE interpol_smv_id IS NOT NULL;

COMMIT;
