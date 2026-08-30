BEGIN;

CREATE TABLE sipep_movements (
    movement_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inmate_id     UUID NOT NULL REFERENCES sipep_inmates(inmate_id),
    from_block    VARCHAR(20),
    to_block      VARCHAR(20) NOT NULL,
    from_facility VARCHAR(100),
    to_facility   VARCHAR(100) NOT NULL,
    movement_type VARCHAR(30) NOT NULL,
    reason        TEXT,
    authorized_by UUID NOT NULL,
    moved_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
