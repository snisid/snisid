BEGIN;

CREATE TYPE siar_weapon_type AS ENUM (
    'HANDGUN','RIFLE','SHOTGUN','SUBMACHINE_GUN','ASSAULT_RIFLE',
    'MACHINE_GUN','SNIPER','RPG','GRENADE','HOMEMADE','OTHER'
);

CREATE TYPE siar_status AS ENUM (
    'REGISTERED','REPORTED_STOLEN','SEIZED','DESTROYED',
    'REPORTED_LOST','TRANSFERRED','DEACTIVATED'
);

CREATE TYPE siar_registration_type AS ENUM (
    'CIVILIAN','POLICE','MILITARY','SECURITY_COMPANY',
    'EMBASSY','ILLEGAL_FOUND','HISTORICAL'
);

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

CREATE TABLE siar_licenses (
    license_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    license_number      VARCHAR(50) UNIQUE NOT NULL,
    holder_snisid_id    UUID NOT NULL,
    license_type        VARCHAR(50) NOT NULL,
    firearms_authorized INTEGER DEFAULT 1,
    issue_date          DATE NOT NULL,
    expiry_date         DATE NOT NULL,
    issuing_authority   VARCHAR(100) NOT NULL,
    is_active           BOOLEAN DEFAULT TRUE,
    revocation_reason   TEXT,
    revoked_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE siar_transfers (
    transfer_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firearm_id          UUID NOT NULL REFERENCES siar_firearms(firearm_id),
    from_owner_id       UUID,
    to_owner_id         UUID,
    transfer_type       VARCHAR(50),
    transfer_date       DATE NOT NULL,
    permit_ref          VARCHAR(100),
    authorized_by       UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE siar_seizures (
    seizure_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firearm_id          UUID REFERENCES siar_firearms(firearm_id),
    seizure_date        TIMESTAMPTZ NOT NULL,
    seizing_unit        VARCHAR(50) NOT NULL,
    seizing_officer     UUID,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    context             TEXT,
    from_person_id      UUID,
    gang_id             UUID,
    case_reference      VARCHAR(100),
    disposed_of         BOOLEAN DEFAULT FALSE,
    disposal_method     VARCHAR(50),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_siar_serial    ON siar_firearms(serial_number) WHERE serial_number IS NOT NULL;
CREATE INDEX idx_siar_status    ON siar_firearms(status);
CREATE INDEX idx_siar_gang      ON siar_firearms(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_siar_owner     ON siar_firearms(owner_snisid_id) WHERE owner_snisid_id IS NOT NULL;
CREATE INDEX idx_siar_iarms     ON siar_firearms(iarms_ref) WHERE iarms_ref IS NOT NULL;
CREATE INDEX idx_siar_licenses  ON siar_licenses(holder_snisid_id, is_active);

COMMIT;
