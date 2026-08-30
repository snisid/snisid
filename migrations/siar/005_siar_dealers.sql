BEGIN;

CREATE TABLE siar_dealers (
    dealer_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dealer_license_no   VARCHAR(50) UNIQUE NOT NULL,
    business_name       VARCHAR(200) NOT NULL,
    business_reg_no     VARCHAR(100),
    owner_snisid_id     UUID NOT NULL,
    owner_name          VARCHAR(200) NOT NULL,
    address             TEXT,
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    phone               VARCHAR(30),
    email               VARCHAR(100),
    license_type        siar_license_type NOT NULL DEFAULT 'DEALER',
    status              siar_dealer_status NOT NULL DEFAULT 'ACTIVE',
    license_issue_date  DATE NOT NULL,
    license_expiry_date DATE NOT NULL,
    premises_inspected  BOOLEAN DEFAULT FALSE,
    last_inspection_date DATE,
    notes               TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
