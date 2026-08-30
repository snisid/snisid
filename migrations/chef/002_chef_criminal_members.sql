BEGIN;

CREATE TABLE chef_criminal_members (
    member_id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_chef_id        VARCHAR(25) UNIQUE NOT NULL,
    snisid_person_id        UUID NOT NULL,
    fir_record_id           UUID,
    afis_subject_id         UUID,
    rdep_deportee_id        UUID,

    primary_gang_id         UUID NOT NULL,
    role_in_gang            chef_role_type NOT NULL,
    role_description        TEXT,
    joined_date             DATE,
    rank_level              SMALLINT,

    aliases                 TEXT[] DEFAULT '{}',
    known_languages         TEXT[] DEFAULT '{}',
    tattoo_description      TEXT,
    physical_description    TEXT,
    photo_refs              TEXT[] DEFAULT '{}',

    territory_dept          CHAR(2),
    territory_communes      TEXT[] DEFAULT '{}',

    known_armed             BOOLEAN DEFAULT FALSE,
    weapon_types            TEXT[] DEFAULT '{}',
    trained_combatant       BOOLEAN DEFAULT FALSE,

    status                  chef_status NOT NULL DEFAULT 'ACTIVE',
    un_designated           BOOLEAN DEFAULT FALSE,
    un_designation_date     TIMESTAMPTZ,
    ofac_designated         BOOLEAN DEFAULT FALSE,
    ofac_sdn_ref            VARCHAR(50),
    interpol_notice_ref     VARCHAR(50),

    last_known_address      TEXT,
    last_known_dept         CHAR(2),
    last_seen_at            TIMESTAMPTZ,
    last_seen_location      VARCHAR(300),

    estimated_wealth_usd    DECIMAL(15,2),
    known_assets            TEXT[] DEFAULT '{}',

    intel_classification    VARCHAR(20) DEFAULT 'SECRET',
    intel_confidence        SMALLINT CHECK (intel_confidence BETWEEN 1 AND 10),
    created_by              UUID NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
