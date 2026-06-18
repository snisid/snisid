BEGIN;

CREATE TABLE blkl_blacklist (
    entry_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_blkl_id    VARCHAR(25) UNIQUE NOT NULL,
    snisid_person_id    UUID NOT NULL,
    restriction_type    blkl_restriction_type NOT NULL,
    source              blkl_source NOT NULL,
    source_record_id    UUID,
    reason              TEXT NOT NULL,
    court_order_ref     VARCHAR(100),
    ordered_by          VARCHAR(150),
    effective_date      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expiry_date         TIMESTAMPTZ,
    is_permanent        BOOLEAN DEFAULT FALSE,
    is_active           BOOLEAN DEFAULT TRUE,
    alert_level         VARCHAR(20) DEFAULT 'WANTED',
    armed_dangerous     BOOLEAN DEFAULT FALSE,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
