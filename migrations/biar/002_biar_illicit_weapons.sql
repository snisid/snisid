BEGIN;

CREATE TABLE biar_illicit_weapons (
    weapon_id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_biar_id        VARCHAR(25) UNIQUE NOT NULL,
    serial_number           VARCHAR(100),
    serial_obliterated      BOOLEAN DEFAULT FALSE,
    make                    VARCHAR(100),
    model                   VARCHAR(100),
    caliber                 VARCHAR(30),
    weapon_type             VARCHAR(50) NOT NULL,
    manufacture_country     CHAR(3),
    estimated_manufacture_year SMALLINT,

    recovery_date           TIMESTAMPTZ NOT NULL,
    recovery_context        biar_recovery_context NOT NULL,
    recovery_location       VARCHAR(300),
    recovery_dept_code      CHAR(2),
    recovery_commune        VARCHAR(100),
    recovery_lat            DECIMAL(10,7),
    recovery_lng            DECIMAL(10,7),
    seizing_unit            VARCHAR(50) NOT NULL,
    seizing_officer         UUID,
    case_reference          VARCHAR(100),

    from_person_id          UUID,
    gang_id                 UUID,
    crime_category          VARCHAR(50),
    associated_cases        TEXT[] DEFAULT '{}',

    origin_country          CHAR(3),
    transit_countries       CHAR(3)[] DEFAULT '{}',
    trafficking_route       TEXT,
    import_method           TEXT,

    iarms_ref               VARCHAR(50),
    atf_etrace_ref          VARCHAR(50),
    reported_to_interpol    BOOLEAN DEFAULT FALSE,
    interpol_reported_at    TIMESTAMPTZ,

    disposition             biar_weapon_disposition DEFAULT 'HELD_AS_EVIDENCE',
    disposal_date           TIMESTAMPTZ,
    disposal_auth           UUID,

    quantity_ammunition     INTEGER DEFAULT 0,
    ammunition_type         VARCHAR(50),
    photos_refs             TEXT[] DEFAULT '{}',
    notes                   TEXT,
    created_by              UUID NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
