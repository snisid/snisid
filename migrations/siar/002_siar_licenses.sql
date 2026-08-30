BEGIN;

CREATE TABLE siar_licenses (
    license_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    license_number      VARCHAR(50) UNIQUE NOT NULL,
    holder_snisid_id    UUID NOT NULL,
    holder_name         VARCHAR(200) NOT NULL,
    license_type        siar_license_type NOT NULL,
    firearms_authorized INTEGER DEFAULT 1,
    issue_date          DATE NOT NULL,
    expiry_date         DATE NOT NULL,
    issuing_authority   VARCHAR(100) NOT NULL,
    is_active           BOOLEAN DEFAULT TRUE,
    revocation_reason   TEXT,
    revoked_at          TIMESTAMPTZ,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
