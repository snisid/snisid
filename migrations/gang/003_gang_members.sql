BEGIN;

CREATE TABLE gang_members (
    member_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gang_id            UUID NOT NULL REFERENCES gang_organizations(gang_id) ON DELETE CASCADE,
    national_member_id VARCHAR(25) UNIQUE NOT NULL,
    full_name          VARCHAR(200) NOT NULL,
    aliases            TEXT[] DEFAULT '{}',
    role               VARCHAR(100),
    date_of_birth      DATE,
    place_of_birth     VARCHAR(200),
    nationality        CHAR(3) DEFAULT 'HTI',
    id_type            VARCHAR(30),
    id_number          VARCHAR(50),
    photo_ref          TEXT,
    fingerprint_hash   VARCHAR(128),
    last_known_address TEXT,
    dept_code          CHAR(2),
    commune            VARCHAR(100),
    lat                DECIMAL(10,7),
    lng                DECIMAL(10,7),
    is_leader          BOOLEAN DEFAULT FALSE,
    is_arrested        BOOLEAN DEFAULT FALSE,
    arrest_date        TIMESTAMPTZ,
    arrest_ref         VARCHAR(50),
    is_deceased        BOOLEAN DEFAULT FALSE,
    death_date         TIMESTAMPTZ,
    ofac_designated    BOOLEAN DEFAULT FALSE,
    ofac_sdn_ref       VARCHAR(50),
    intel_confidence   SMALLINT CHECK (intel_confidence BETWEEN 1 AND 10),
    notes              TEXT,
    created_by         UUID NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
