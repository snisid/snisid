BEGIN;

CREATE TABLE siar_firearms (
    firearm_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_siar_id    VARCHAR(25) UNIQUE NOT NULL,
    serial_number       VARCHAR(100),
    make                VARCHAR(100) NOT NULL,
    model               VARCHAR(100) NOT NULL,
    caliber             VARCHAR(30) NOT NULL,
    weapon_type         siar_weapon_type NOT NULL,
    manufacture_year    SMALLINT,
    manufacture_country CHAR(3),
    status              siar_status NOT NULL DEFAULT 'REGISTERED',
    reg_type            siar_registration_type NOT NULL,

    owner_snisid_id     UUID,
    owner_entity_name   VARCHAR(200),
    license_number      VARCHAR(50),
    license_expiry      DATE,

    import_date         DATE,
    import_country      CHAR(3),
    import_permit_ref   VARCHAR(100),
    importer_name       VARCHAR(200),
    customs_entry_ref   VARCHAR(100),

    current_dept_code   CHAR(2),
    storage_location    TEXT,

    fir_record_id       UUID,
    gang_id             UUID,
    case_references     TEXT[] DEFAULT '{}',

    iarms_ref           VARCHAR(50),
    atf_etrace_ref      VARCHAR(50),

    notes               TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
