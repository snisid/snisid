CREATE TABLE IF NOT EXISTS criminal_members (
    member_id           TEXT PRIMARY KEY,
    national_chef_id    TEXT,
    snisid_person_id    TEXT,
    fir_record_id       TEXT,
    afis_subject_id     TEXT,
    rdep_deportee_id    TEXT,
    primary_gang_id     TEXT NOT NULL,
    role_in_gang        TEXT NOT NULL,
    role_description    TEXT,
    joined_date         TIMESTAMPTZ,
    rank_level          INTEGER,
    aliases             JSONB DEFAULT '[]',
    known_languages     JSONB DEFAULT '[]',
    tattoo_description  TEXT,
    physical_description TEXT,
    photo_refs          JSONB DEFAULT '[]',
    territory_dept      TEXT,
    territory_communes  JSONB DEFAULT '[]',
    known_armed         BOOLEAN DEFAULT FALSE,
    weapon_types        JSONB DEFAULT '[]',
    trained_combatant   BOOLEAN DEFAULT FALSE,
    status              TEXT NOT NULL DEFAULT 'ACTIVE',
    un_designated       BOOLEAN DEFAULT FALSE,
    ofac_designated     BOOLEAN DEFAULT FALSE,
    ofac_sdn_ref        TEXT,
    interpol_notice_ref TEXT,
    last_known_address  TEXT,
    last_seen_at        TIMESTAMPTZ,
    intel_confidence    INTEGER,
    created_by          TEXT NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS intelligence_notes (
    note_id      TEXT PRIMARY KEY,
    member_id    TEXT NOT NULL REFERENCES criminal_members(member_id),
    source_id    TEXT,
    note_type    TEXT NOT NULL,
    content      TEXT NOT NULL,
    confidence   INTEGER,
    collected_at TIMESTAMPTZ,
    created_by   TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sightings (
    sighting_id TEXT PRIMARY KEY,
    member_id   TEXT NOT NULL REFERENCES criminal_members(member_id),
    source_id   TEXT,
    dept        TEXT NOT NULL,
    commune     TEXT,
    latitude    DOUBLE PRECISION,
    longitude   DOUBLE PRECISION,
    spotted_at  TIMESTAMPTZ NOT NULL,
    confidence  INTEGER,
    notes       TEXT,
    created_by  TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_members_gang ON criminal_members(primary_gang_id);
CREATE INDEX idx_members_status ON criminal_members(status);
CREATE INDEX idx_members_role ON criminal_members(role_in_gang);
CREATE INDEX idx_intel_member ON intelligence_notes(member_id);
CREATE INDEX idx_sightings_member ON sightings(member_id);
CREATE INDEX idx_sightings_dept ON sightings(dept);
