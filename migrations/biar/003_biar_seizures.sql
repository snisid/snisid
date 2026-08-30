BEGIN;

CREATE TABLE biar_batch_seizures (
    batch_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_reference     VARCHAR(50) UNIQUE NOT NULL,
    operation_name      TEXT,
    seizure_date        TIMESTAMPTZ NOT NULL,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    total_weapons       INTEGER NOT NULL,
    weapon_ids          UUID[] DEFAULT '{}',
    seizing_unit        VARCHAR(50) NOT NULL,
    lead_officer        UUID,
    partnering_agencies TEXT[] DEFAULT '{}',
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
