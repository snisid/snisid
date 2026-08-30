BEGIN;

CREATE TABLE fir_movements (
    movement_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id     UUID NOT NULL REFERENCES fir_criminal_records(record_id),
    charge_id     UUID REFERENCES fir_charges(charge_id),
    movement_type fir_movement_type NOT NULL,
    description   TEXT,
    changed_by    UUID,
    metadata      JSONB DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
