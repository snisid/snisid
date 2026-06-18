BEGIN;

CREATE TYPE sanc_source AS ENUM (
    'UN_2653','OFAC_SDN','EU_CONSOLIDATED',
    'INTERPOL','CANADA_OSFI','UK_OFSI','OTHER'
);

CREATE TYPE sanc_measure AS ENUM (
    'ASSETS_FREEZE','TRAVEL_BAN','ARMS_EMBARGO','ALL_MEASURES'
);

CREATE TYPE sanc_entity_type AS ENUM (
    'INDIVIDUAL','ORGANIZATION','VESSEL','AIRCRAFT'
);

CREATE TABLE sanc_entries (
    sanc_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source              sanc_source NOT NULL,
    source_ref_id       VARCHAR(100) NOT NULL,
    entity_type         sanc_entity_type NOT NULL,
    entity_name         VARCHAR(300) NOT NULL,
    aliases             TEXT[] DEFAULT '{}',
    nationality         TEXT[] DEFAULT '{}',
    date_of_birth       DATE,
    place_of_birth      TEXT,
    passport_numbers    TEXT[] DEFAULT '{}',
    national_id_numbers TEXT[] DEFAULT '{}',
    measure_types       sanc_measure[] DEFAULT '{}',
    listing_date        DATE NOT NULL,
    end_date            DATE,
    is_active           BOOLEAN DEFAULT TRUE,
    listing_reason      TEXT,
    committee_notes     TEXT,
    snisid_person_id    UUID,
    gang_id             UUID,
    chef_member_id      UUID,
    match_confidence    SMALLINT,
    source_updated_at   TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(source, source_ref_id)
);

CREATE TABLE sanc_sync_log (
    sync_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source              sanc_source NOT NULL,
    started_at          TIMESTAMPTZ NOT NULL,
    completed_at        TIMESTAMPTZ,
    entries_processed   INTEGER DEFAULT 0,
    entries_added       INTEGER DEFAULT 0,
    entries_updated     INTEGER DEFAULT 0,
    entries_removed     INTEGER DEFAULT 0,
    errors              INTEGER DEFAULT 0,
    status              VARCHAR(20) DEFAULT 'RUNNING',
    error_details       TEXT
);

CREATE TABLE sanc_identity_matches (
    match_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sanc_id             UUID NOT NULL REFERENCES sanc_entries(sanc_id),
    snisid_person_id    UUID NOT NULL,
    match_score         DECIMAL(5,2) NOT NULL,
    match_fields        TEXT[] DEFAULT '{}',
    confirmed_by        UUID,
    is_confirmed        BOOLEAN DEFAULT FALSE,
    is_false_positive   BOOLEAN DEFAULT FALSE,
    reviewed_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sanc_source    ON sanc_entries(source, is_active);
CREATE INDEX idx_sanc_name_fts  ON sanc_entries USING gin(to_tsvector('simple', entity_name));
CREATE INDEX idx_sanc_aliases   ON sanc_entries USING gin(aliases);
CREATE INDEX idx_sanc_matches   ON sanc_identity_matches(snisid_person_id);

COMMIT;
