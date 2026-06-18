BEGIN;

CREATE TABLE siar_transfers (
    transfer_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firearm_id          UUID NOT NULL REFERENCES siar_firearms(firearm_id),
    from_owner_id       UUID,
    from_owner_name     VARCHAR(200),
    to_owner_id         UUID,
    to_owner_name       VARCHAR(200),
    transfer_type       siar_transfer_type NOT NULL,
    transfer_date       DATE NOT NULL,
    permit_ref          VARCHAR(100),
    authorized_by       UUID,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE siar_seizures (
    seizure_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firearm_id          UUID REFERENCES siar_firearms(firearm_id),
    serial_number       VARCHAR(100),
    make                VARCHAR(100),
    model               VARCHAR(100),
    seizure_date        TIMESTAMPTZ NOT NULL,
    seizing_unit        VARCHAR(50) NOT NULL,
    seizing_officer     UUID,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    context             TEXT,
    from_person_id      UUID,
    from_person_name    VARCHAR(200),
    gang_id             UUID,
    case_reference      VARCHAR(100),
    disposed_of         BOOLEAN DEFAULT FALSE,
    disposal_method     VARCHAR(50),
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
