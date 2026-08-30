-- OPR-HT Migration: Ordonnances de Protection et Restrictions Judiciaires

BEGIN;

CREATE TYPE opr_order_type AS ENUM (
    'RESTRAINING_ORDER', 'NO_CONTACT', 'STAY_AWAY',
    'PROTECTIVE', 'WITNESS_PROTECTION', 'GANG_EXCLUSION_ZONE',
    'TRAVEL_RESTRICTION'
);

CREATE TYPE opr_status AS ENUM (
    'ACTIVE', 'EXPIRED', 'VIOLATED', 'DISMISSED', 'APPEALED'
);

CREATE TABLE opr_protection_orders (
    order_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number        VARCHAR(30) UNIQUE NOT NULL,     -- OPR-HT-AAAA-NNNNNN
    order_type          opr_order_type NOT NULL,
    status              opr_status NOT NULL DEFAULT 'ACTIVE',

    protected_person_id UUID NOT NULL,
    subject_person_id   UUID NOT NULL,
    subject_fir_id      UUID,

    exclusion_radius_m  INTEGER,
    exclusion_addresses TEXT[] DEFAULT '{}',
    no_contact_modes    TEXT[] DEFAULT '{}',
    geographic_ban_geojson JSONB,

    issuing_court       VARCHAR(150) NOT NULL,
    issuing_judge       VARCHAR(150),
    issue_date          TIMESTAMPTZ NOT NULL,
    expiry_date         TIMESTAMPTZ NOT NULL,
    is_renewable        BOOLEAN DEFAULT TRUE,

    violation_count     SMALLINT DEFAULT 0,
    last_violation_at   TIMESTAMPTZ,

    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE opr_violations (
    violation_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id            UUID NOT NULL REFERENCES opr_protection_orders(order_id),
    violation_date      TIMESTAMPTZ NOT NULL,
    violation_type      VARCHAR(100) NOT NULL,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    reported_by         UUID NOT NULL,
    arrest_made         BOOLEAN DEFAULT FALSE,
    arrest_case_ref     VARCHAR(100),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE opr_witness_protections (
    protection_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    protected_person_id UUID NOT NULL,
    threat_level        VARCHAR(20) NOT NULL,
    gang_id             UUID,
    alias_assigned      VARCHAR(150),
    assigned_unit       VARCHAR(50),
    is_active           BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_opr_subject   ON opr_protection_orders(subject_person_id) WHERE status = 'ACTIVE';
CREATE INDEX idx_opr_protected ON opr_protection_orders(protected_person_id) WHERE status = 'ACTIVE';
CREATE INDEX idx_opr_expiry    ON opr_protection_orders(expiry_date) WHERE status = 'ACTIVE';

COMMIT;
