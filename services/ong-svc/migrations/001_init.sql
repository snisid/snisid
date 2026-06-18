CREATE TYPE ong_registration_status AS ENUM (
    'REGISTERED', 'PENDING', 'SUSPENDED', 'REVOKED', 'OPERATING_WITHOUT_REGISTRATION'
);

CREATE TYPE ong_type AS ENUM (
    'HUMANITARIAN', 'DEVELOPMENT', 'ADVOCACY', 'FAITH_BASED',
    'DIASPORA', 'RESEARCH', 'MIXED', 'UNKNOWN'
);

CREATE TYPE ong_risk_flag AS ENUM (
    'NONE', 'FINANCIAL_IRREGULARITY', 'STAFF_SECURITY_CONCERN',
    'OPERATING_IN_RESTRICTED_ZONE', 'SANCTION_MATCH',
    'SUSPECTED_FRONT_ORGANIZATION', 'UNREGISTERED_ILLEGAL'
);

CREATE TABLE ong_organizations (
    org_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_ong_id     VARCHAR(25) UNIQUE NOT NULL,
    org_name            VARCHAR(200) NOT NULL,
    org_name_local      VARCHAR(200),
    acronym             VARCHAR(20),
    org_type            ong_type NOT NULL,
    registration_status ong_registration_status NOT NULL DEFAULT 'PENDING',
    mjsp_registration_number VARCHAR(50),
    registration_date   DATE,
    registration_expiry DATE,
    headquarter_country CHAR(3) NOT NULL,
    headquarter_city    VARCHAR(100),
    haiti_office_dept   CHAR(2),
    haiti_office_address TEXT,
    haiti_office_lat    DECIMAL(10,7),
    haiti_office_lng    DECIMAL(10,7),
    operating_depts     CHAR(2)[] DEFAULT '{}',
    operating_communes  TEXT[] DEFAULT '{}',
    sectors             TEXT[] DEFAULT '{}',
    annual_budget_usd   DECIMAL(15,2),
    funding_sources     TEXT[] DEFAULT '{}',
    major_donors        TEXT[] DEFAULT '{}',
    haiti_staff_count   INTEGER DEFAULT 0,
    expat_staff_count   INTEGER DEFAULT 0,
    director_name       VARCHAR(200),
    director_snisid_id  UUID,
    director_nationality CHAR(3),
    contact_email       VARCHAR(200),
    contact_phone       VARCHAR(30),
    risk_flag           ong_risk_flag NOT NULL DEFAULT 'NONE',
    risk_notes          TEXT,
    sanc_match_id       UUID,
    blan_case_id        UUID,
    ulcc_ref            VARCHAR(50),
    is_access_restricted BOOLEAN DEFAULT FALSE,
    access_restriction_reason TEXT,
    last_compliance_review DATE,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE ong_staff_registry (
    staff_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id              UUID NOT NULL REFERENCES ong_organizations(org_id),
    snisid_person_id    UUID,
    full_name           VARCHAR(200) NOT NULL,
    nationality         CHAR(3) NOT NULL,
    role                VARCHAR(100),
    is_expatriate       BOOLEAN DEFAULT FALSE,
    passport_number     VARCHAR(50),
    visa_type           VARCHAR(30),
    visa_expiry         DATE,
    entry_date          DATE,
    haiti_address       TEXT,
    dept_code           CHAR(2),
    sltd_check_passed   BOOLEAN DEFAULT FALSE,
    blkl_check_passed   BOOLEAN DEFAULT FALSE,
    sanc_check_passed   BOOLEAN DEFAULT FALSE,
    last_security_check TIMESTAMPTZ,
    is_active           BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE ong_field_access_requests (
    request_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id              UUID NOT NULL REFERENCES ong_organizations(org_id),
    access_type         VARCHAR(30),
    requested_zones     TEXT[] DEFAULT '{}',
    requested_depts     CHAR(2)[] DEFAULT '{}',
    access_date         DATE NOT NULL,
    access_date_end     DATE,
    purpose             TEXT NOT NULL,
    vehicle_count       SMALLINT DEFAULT 1,
    staff_count         SMALLINT DEFAULT 1,
    status              VARCHAR(20) DEFAULT 'PENDING',
    pnh_escort_required BOOLEAN DEFAULT FALSE,
    approved_by         UUID,
    approval_notes      TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ong_orgs_status ON ong_organizations(registration_status, risk_flag);
CREATE INDEX idx_ong_orgs_dept ON ong_organizations USING gin(operating_depts);
CREATE INDEX idx_ong_orgs_country ON ong_organizations(headquarter_country);
CREATE INDEX idx_ong_orgs_risk ON ong_organizations(risk_flag) WHERE risk_flag != 'NONE';
CREATE INDEX idx_ong_staff_org ON ong_staff_registry(org_id, is_active);
CREATE INDEX idx_ong_staff_passport ON ong_staff_registry(passport_number);
CREATE INDEX idx_ong_access_date ON ong_field_access_requests(access_date, status);
