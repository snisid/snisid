BEGIN;

CREATE TYPE biar_recovery_context AS ENUM (
    'POLICE_OPERATION','CHECKPOINT','PORT_SEIZURE','AIRPORT_SEIZURE',
    'COMMUNITY_SURRENDER','CRIME_SCENE','RAID','BORDER_SEIZURE','OTHER'
);

CREATE TYPE biar_weapon_disposition AS ENUM (
    'HELD_AS_EVIDENCE','DESTROYED','RETURNED_TO_OWNER',
    'TRANSFERRED_TO_POLICE','SENT_TO_INTERPOL','PENDING'
);

CREATE TABLE biar_illicit_weapons (
    weapon_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_biar_id    VARCHAR(25) UNIQUE NOT NULL,
    serial_number       VARCHAR(100),
    serial_obliterated  BOOLEAN DEFAULT FALSE,
    make                VARCHAR(100),
    model               VARCHAR(100),
    caliber             VARCHAR(30),
    weapon_type         VARCHAR(50) NOT NULL,
    manufacture_country CHAR(3),
    estimated_manufacture_year SMALLINT,
    recovery_date       TIMESTAMPTZ NOT NULL,
    recovery_context    biar_recovery_context NOT NULL,
    recovery_location   VARCHAR(300),
    recovery_dept_code  CHAR(2),
    recovery_commune    VARCHAR(100),
    recovery_lat        DECIMAL(10,7),
    recovery_lng        DECIMAL(10,7),
    seizing_unit        VARCHAR(50) NOT NULL,
    seizing_officer     UUID,
    case_reference      VARCHAR(100),
    from_person_id      UUID,
    gang_id             UUID,
    crime_category      VARCHAR(50),
    associated_cases    TEXT[] DEFAULT '{}',
    origin_country      CHAR(3),
    transit_countries   CHAR(3)[] DEFAULT '{}',
    trafficking_route   TEXT,
    import_method       TEXT,
    iarms_ref           VARCHAR(50),
    atf_etrace_ref      VARCHAR(50),
    reported_to_interpol BOOLEAN DEFAULT FALSE,
    interpol_reported_at TIMESTAMPTZ,
    disposition         biar_weapon_disposition DEFAULT 'HELD_AS_EVIDENCE',
    disposal_date       TIMESTAMPTZ,
    disposal_auth       UUID,
    quantity_ammunition INTEGER DEFAULT 0,
    ammunition_type     VARCHAR(50),
    photos_refs         TEXT[] DEFAULT '{}',
    notes               TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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

CREATE TABLE biar_iarms_sync_log (
    sync_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    weapon_id           UUID REFERENCES biar_illicit_weapons(weapon_id),
    direction           VARCHAR(10) NOT NULL,
    iarms_ref           VARCHAR(50),
    sync_status         VARCHAR(20) DEFAULT 'PENDING',
    synced_at           TIMESTAMPTZ,
    error_message       TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_biar_serial    ON biar_illicit_weapons(serial_number) WHERE serial_number IS NOT NULL;
CREATE INDEX idx_biar_gang      ON biar_illicit_weapons(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_biar_dept      ON biar_illicit_weapons(recovery_dept_code);
CREATE INDEX idx_biar_date      ON biar_illicit_weapons(recovery_date DESC);
CREATE INDEX idx_biar_iarms     ON biar_illicit_weapons(iarms_ref) WHERE iarms_ref IS NOT NULL;
CREATE INDEX idx_biar_origin    ON biar_illicit_weapons(origin_country);

COMMIT;
