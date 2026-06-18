-- FIR-HT Migration: Fichier Individuel des Renseignements
-- Casier Judiciaire National d'Haïti

BEGIN;

CREATE TYPE fir_offense_class AS ENUM (
    'CONTRAVENTION',    -- Infraction mineure
    'DELIT',            -- Infraction correctionnelle
    'CRIME',            -- Crime grave
    'FELONY_FOREIGN'    -- Infraction étrangère (déportés)
);

CREATE TYPE fir_case_status AS ENUM (
    'OPEN','PENDING_TRIAL','CONVICTED','ACQUITTED',
    'DISMISSED','APPEAL_PENDING','EXPUNGED'
);

CREATE TYPE fir_sentence_type AS ENUM (
    'PRISON','SUSPENDED','FINE','COMMUNITY_SERVICE',
    'DEATH_PENALTY','ACQUITTAL','PROBATION'
);

CREATE SEQUENCE fir_record_seq START 1;

CREATE TABLE fir_criminal_records (
    record_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_fir_id     VARCHAR(25) UNIQUE NOT NULL,   -- Format: FIR-HT-AAAA-NNNNNN
    snisid_person_id    UUID NOT NULL,                  -- Lien identité SNISID
    afis_subject_id     UUID,                           -- Lien AFIS-HT
    is_haitian_national BOOLEAN DEFAULT TRUE,
    aliases             TEXT[] DEFAULT '{}',
    is_active           BOOLEAN DEFAULT TRUE,
    is_expunged         BOOLEAN DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fir_arrests (
    arrest_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id           UUID NOT NULL REFERENCES fir_criminal_records(record_id),
    arrest_date         TIMESTAMPTZ NOT NULL,
    arresting_unit      VARCHAR(50) NOT NULL,
    arresting_officer   UUID,
    arrest_location     VARCHAR(300),
    dept_code           CHAR(2),
    charges_text        TEXT NOT NULL,
    offense_class       fir_offense_class NOT NULL,
    case_reference      VARCHAR(100),
    release_date        TIMESTAMPTZ,
    release_reason      TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fir_convictions (
    conviction_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id           UUID NOT NULL REFERENCES fir_criminal_records(record_id),
    case_reference      VARCHAR(100) NOT NULL,
    court_name          VARCHAR(150) NOT NULL,
    court_dept          CHAR(2),
    offense_class       fir_offense_class NOT NULL,
    offense_description TEXT NOT NULL,
    ipc_code            VARCHAR(30),                    -- Code pénal haïtien
    verdict_date        TIMESTAMPTZ NOT NULL,
    case_status         fir_case_status NOT NULL,
    sentence_type       fir_sentence_type,
    sentence_duration_days INTEGER,
    fine_amount_gdes    DECIMAL(12,2),
    sentence_start      TIMESTAMPTZ,
    sentence_end        TIMESTAMPTZ,
    is_foreign_record   BOOLEAN DEFAULT FALSE,
    foreign_country     CHAR(3),                        -- HTI, USA, DOM, etc.
    interpol_ccc_ref    VARCHAR(50),
    judge_name          VARCHAR(150),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fir_certificates (
    cert_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id           UUID REFERENCES fir_criminal_records(record_id),
    snisid_person_id    UUID NOT NULL,
    certificate_number  VARCHAR(30) UNIQUE NOT NULL,
    issued_for          VARCHAR(200),                   -- Motif: emploi, visa, etc.
    result              VARCHAR(20) NOT NULL,            -- CLEAN, HAS_RECORD
    issued_by           UUID NOT NULL,
    issuing_office      VARCHAR(100),
    issued_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at          TIMESTAMPTZ,
    qr_code_ref         VARCHAR(200)                    -- Vérification authenticité
);

CREATE INDEX idx_fir_records_snisid   ON fir_criminal_records(snisid_person_id);
CREATE INDEX idx_fir_records_fir_id   ON fir_criminal_records(national_fir_id);
CREATE INDEX idx_fir_arrests_date     ON fir_arrests(arrest_date DESC);
CREATE INDEX idx_fir_convictions_date ON fir_convictions(verdict_date DESC);
CREATE INDEX idx_fir_convictions_dept ON fir_convictions(court_dept);

ALTER TABLE fir_criminal_records ENABLE ROW LEVEL SECURITY;
CREATE POLICY fir_read_policy ON fir_criminal_records
    FOR SELECT USING (
        current_setting('app.user_role', TRUE) IN
        ('DCPJ','PARQUET','TRIBUNAL','JUDGE','POLICE_OFFICER','SUPERADMIN')
    );

COMMIT;
