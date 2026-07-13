# SNISID-BIO-ADN — Schémas PostgreSQL Complets
**Document ID :** SNISID-BIO-SQL-001 | **Version :** 1.0.0

---

## MIGRATION 001 — Index ADN (Catégorie CODIS)

```sql
-- ============================================================
-- SNISID-BIO-ADN : Schémas PostgreSQL
-- Inspiré CODIS (FBI) adapté à Haïti
-- Classification : SOUVERAIN / JUDICIAIRE
-- ============================================================

-- Extension nécessaire
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Schéma dédié
CREATE SCHEMA IF NOT EXISTS bio_adn;
SET search_path TO bio_adn, public;

-- ──────────────────────────────────────────────────
-- 1. LABORATOIRES ACCRÉDITÉS (LDIS/SDIS/NDIS)
-- ──────────────────────────────────────────────────
CREATE TABLE bio_laboratories (
    lab_id          UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    lab_code        VARCHAR(20) UNIQUE NOT NULL,     -- ex: LDIS-PAP-001
    lab_name        VARCHAR(200) NOT NULL,
    lab_level       VARCHAR(10) NOT NULL             -- LDIS, SDIS, NDIS
                    CHECK (lab_level IN ('LDIS','SDIS','NDIS')),
    department      VARCHAR(50),                     -- Département haïtien
    institution     VARCHAR(100),                    -- PNH, DCPJ, ANH, MSPP
    accreditation   VARCHAR(100),                    -- Numéro accréditation
    contact_email   VARCHAR(200),
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- ──────────────────────────────────────────────────
-- 2. PROFILS ADN STR — 20 LOCI CODIS CORE
-- ──────────────────────────────────────────────────
-- PRINCIPE : PAS de données nominales ici (comme CODIS)
CREATE TABLE bio_str_profiles (
    sample_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    specimen_number VARCHAR(100) UNIQUE NOT NULL,    -- Numéro de scellé labo
    index_type      VARCHAR(10) NOT NULL
                    CHECK (index_type IN ('BIO-CON','BIO-ARR','BIO-FSC','BIO-DIS','BIO-RNI')),
    -- 20 Loci CODIS Core (valeurs chiffrées AES-256-GCM via HSM)
    loci_encrypted  BYTEA NOT NULL,                  -- JSON chiffré des 20 loci
    loci_hash       VARCHAR(64) NOT NULL,            -- SHA-256 pour déduplication
    amelogenin      CHAR(2),                         -- XX ou XY (sexe)
    quality_score   DECIMAL(4,3) CHECK (quality_score BETWEEN 0 AND 1),
    loci_count      SMALLINT DEFAULT 20,             -- Nombre de loci valides
    -- Métadonnées labo (pas d'identité)
    lab_id          UUID REFERENCES bio_laboratories(lab_id),
    case_number     VARCHAR(100),                    -- Numéro dossier judiciaire
    collected_date  DATE NOT NULL,
    analysis_date   DATE,
    -- Niveaux d'index
    uploaded_ldis   BOOLEAN DEFAULT FALSE,
    uploaded_sdis   BOOLEAN DEFAULT FALSE,
    uploaded_ndis   BOOLEAN DEFAULT FALSE,
    ndis_upload_date TIMESTAMPTZ,
    -- Statut
    is_expunged     BOOLEAN DEFAULT FALSE,           -- Effacement légal
    expunge_date    TIMESTAMPTZ,
    expunge_order   VARCHAR(200),
    -- Audit
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Index pour matching rapide sur le hash
CREATE INDEX idx_bio_str_hash ON bio_str_profiles(loci_hash);
CREATE INDEX idx_bio_str_index_type ON bio_str_profiles(index_type);
CREATE INDEX idx_bio_str_case ON bio_str_profiles(case_number);

-- ──────────────────────────────────────────────────
-- 3. LIAISON ADN ↔ IDENTITÉ (accès restreint DCPJ-DIR)
-- ──────────────────────────────────────────────────
CREATE TABLE bio_identity_links (
    link_id         UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    sample_id       UUID UNIQUE REFERENCES bio_str_profiles(sample_id),
    niu             VARCHAR(20),                     -- NIU SNISID Core
    linked_by_agent VARCHAR(100) NOT NULL,
    linked_at       TIMESTAMPTZ DEFAULT NOW(),
    court_order_ref VARCHAR(200) NOT NULL,           -- Ordonnance tribunal
    purpose         VARCHAR(100) NOT NULL,
    -- Revue légale obligatoire
    reviewed_by     VARCHAR(100),
    reviewed_at     TIMESTAMPTZ,
    review_outcome  VARCHAR(20)                      -- APPROVED, REJECTED
);

-- ──────────────────────────────────────────────────
-- 4. HITS ADN — Correspondances trouvées
-- ──────────────────────────────────────────────────
CREATE TABLE bio_hits (
    hit_id          UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    query_sample_id UUID REFERENCES bio_str_profiles(sample_id),
    match_sample_id UUID REFERENCES bio_str_profiles(sample_id),
    match_type      VARCHAR(20) NOT NULL             -- FULL_MATCH, PARTIAL, FAMILIAL
                    CHECK (match_type IN ('FULL_MATCH','PARTIAL','FAMILIAL')),
    confidence      DECIMAL(5,4) NOT NULL,
    matched_loci    SMALLINT NOT NULL,
    total_loci      SMALLINT NOT NULL,
    hit_level       VARCHAR(10) NOT NULL             -- LDIS, SDIS, NDIS
                    CHECK (hit_level IN ('LDIS','SDIS','NDIS')),
    alert_sent      BOOLEAN DEFAULT FALSE,
    alert_sent_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

---

## MIGRATION 002 — Index Personnes (Catégorie NCIC)

```sql
-- ──────────────────────────────────────────────────
-- 5. PERSONNES RECHERCHÉES (PER-REC) — Equivalent NCIC Wanted Person
-- ──────────────────────────────────────────────────
CREATE TABLE per_wanted_persons (
    record_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_number   VARCHAR(50) UNIQUE NOT NULL,     -- Format: PRE-2026-000001
    niu             VARCHAR(20),                     -- NIU si connu
    -- Identité (peut être partielle pour fugitifs non identifiés)
    last_name       VARCHAR(100),
    first_name      VARCHAR(100),
    aliases         JSONB DEFAULT '[]',              -- Alias connus
    date_of_birth   DATE,
    gender          CHAR(1),
    nationality     CHAR(3),                         -- HTI, etc.
    -- Mandat
    warrant_type    VARCHAR(50) NOT NULL,            -- ARREST, EXTRADITION, etc.
    warrant_number  VARCHAR(100),
    issuing_court   VARCHAR(200),
    issuing_date    DATE NOT NULL,
    charges         JSONB NOT NULL DEFAULT '[]',     -- Chef d'accusations
    danger_level    VARCHAR(10) DEFAULT 'MEDIUM'
                    CHECK (danger_level IN ('LOW','MEDIUM','HIGH','CRITICAL')),
    armed_dangerous BOOLEAN DEFAULT FALSE,
    -- Physique
    height_cm       SMALLINT,
    weight_kg       SMALLINT,
    eye_color       VARCHAR(30),
    hair_color      VARCHAR(30),
    distinguishing_marks TEXT,
    -- Opérationnel
    entering_agency VARCHAR(100) NOT NULL,
    entering_officer VARCHAR(100) NOT NULL,
    mco_contact     VARCHAR(200),                   -- Contact agence entrant
    -- Géo / statut
    last_known_location VARCHAR(200),
    status          VARCHAR(20) DEFAULT 'ACTIVE'
                    CHECK (status IN ('ACTIVE','CLEARED','EXPIRED','SUSPENDED')),
    expiry_date     DATE,
    -- Biométrie
    fingerprint_ref UUID,                           -- Référence NGI-HT (empreintes)
    photo_refs      JSONB DEFAULT '[]',             -- URLs photos chiffrées
    bio_sample_ref  UUID REFERENCES bio_str_profiles(sample_id),
    -- Interopérabilité
    interpol_notice VARCHAR(50),                    -- Notice INTERPOL si applicable
    -- Audit
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    last_hit_at     TIMESTAMPTZ
);

CREATE INDEX idx_per_wanted_niu ON per_wanted_persons(niu);
CREATE INDEX idx_per_wanted_status ON per_wanted_persons(status);
CREATE INDEX idx_per_wanted_name ON per_wanted_persons USING gin(to_tsvector('french', COALESCE(last_name,'') || ' ' || COALESCE(first_name,'')));

-- ──────────────────────────────────────────────────
-- 6. PERSONNES DISPARUES (PER-DIS) — Equivalent NCIC Missing Person
-- ──────────────────────────────────────────────────
CREATE TABLE per_missing_persons (
    record_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_number   VARCHAR(50) UNIQUE NOT NULL,     -- Format: DIS-2026-000001
    niu             VARCHAR(20),
    -- Identité
    last_name       VARCHAR(100) NOT NULL,
    first_name      VARCHAR(100) NOT NULL,
    date_of_birth   DATE,
    age_at_missing  SMALLINT,
    gender          CHAR(1),
    nationality     CHAR(3),
    -- Catégorie (inspiré NCIC)
    category        VARCHAR(20) NOT NULL
                    CHECK (category IN (
                        'CHILD',            -- Moins de 18 ans
                        'ENDANGERED',       -- En danger physique
                        'INVOLUNTARY',      -- Disparition involontaire
                        'CATASTROPHE',      -- Catastrophe naturelle
                        'UNEMANCIPATED',    -- Mineur non émancipé
                        'OTHER'
                    )),
    -- Circonstances
    missing_date    TIMESTAMPTZ NOT NULL,
    missing_location VARCHAR(200) NOT NULL,
    circumstances   TEXT,
    last_seen_clothing VARCHAR(500),
    -- Physique
    height_cm       SMALLINT,
    weight_kg       SMALLINT,
    distinctive_features TEXT,
    -- Famille
    family_contact  VARCHAR(200),
    family_phone    VARCHAR(50),
    -- Biométrie
    photo_refs      JSONB DEFAULT '[]',
    bio_sample_ref  UUID REFERENCES bio_str_profiles(sample_id),    -- ADN famille
    family_bio_refs JSONB DEFAULT '[]',             -- ADN proches (BIO-DIS)
    -- Médical
    medical_conditions TEXT,
    medications     TEXT,
    -- Statut
    status          VARCHAR(20) DEFAULT 'ACTIVE'
                    CHECK (status IN ('ACTIVE','LOCATED','DECEASED','CANCELLED')),
    located_date    TIMESTAMPTZ,
    -- Entrant
    entering_agency VARCHAR(100) NOT NULL,
    ncmec_notified  BOOLEAN DEFAULT FALSE,          -- Équivalent NCMEC pour enfants
    -- Audit
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- ──────────────────────────────────────────────────
-- 7. REGISTRE DÉLINQUANTS SEXUELS (PER-SEX) — Equivalent NSSOR
-- ──────────────────────────────────────────────────
CREATE TABLE per_sex_offenders (
    record_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    niu             VARCHAR(20) NOT NULL,
    conviction_date DATE NOT NULL,
    conviction_court VARCHAR(200) NOT NULL,
    offenses        JSONB NOT NULL DEFAULT '[]',
    risk_level      VARCHAR(10)
                    CHECK (risk_level IN ('LOW','MEDIUM','HIGH')),
    registration_expiry DATE,
    current_address TEXT,                           -- Chiffré en base
    employer        VARCHAR(200),                   -- Chiffré en base
    restrictions    TEXT,
    last_verified   DATE,
    status          VARCHAR(20) DEFAULT 'ACTIVE',
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- ──────────────────────────────────────────────────
-- 8. MEMBRES DE GANGS (PER-GNG) — Equivalent NCIC Gang File
-- ──────────────────────────────────────────────────
CREATE TABLE per_gang_members (
    record_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    niu             VARCHAR(20),
    -- Identité
    last_name       VARCHAR(100),
    first_name      VARCHAR(100),
    aliases         JSONB DEFAULT '[]',
    -- Gang
    gang_name       VARCHAR(200) NOT NULL,
    gang_code       VARCHAR(50),                    -- Code interne DCPJ
    membership_type VARCHAR(30)                     -- LEADER, MEMBER, ASSOCIATE
                    CHECK (membership_type IN ('LEADER','MEMBER','ASSOCIATE','PROSPECT')),
    territory       VARCHAR(200),                   -- Zone opérationnelle
    -- Données opérationnelles
    known_weapons   JSONB DEFAULT '[]',
    criminal_activities JSONB DEFAULT '[]',
    -- Classification
    threat_level    VARCHAR(10) DEFAULT 'HIGH'
                    CHECK (threat_level IN ('LOW','MEDIUM','HIGH','CRITICAL')),
    -- Accès restreint DCPJ uniquement
    intelligence_notes TEXT,                        -- Chiffré
    source_reliability VARCHAR(10),                 -- A,B,C,D,E,F (OTAN)
    -- Statut
    status          VARCHAR(20) DEFAULT 'ACTIVE',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);
```

---

## MIGRATION 003 — Index Biens (Catégorie NCIC)

```sql
-- ──────────────────────────────────────────────────
-- 9. VÉHICULES VOLÉS (BIE-VEH) — Intégration FOVeS/SIV
-- ──────────────────────────────────────────────────
CREATE TABLE bie_stolen_vehicles (
    record_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_number   VARCHAR(50) UNIQUE NOT NULL,     -- Format: VEH-2026-000001
    -- Identification véhicule
    vin             VARCHAR(17) UNIQUE,              -- Vehicle Identification Number
    plate_number    VARCHAR(20),
    plate_dept      VARCHAR(50),
    vehicle_make    VARCHAR(100),
    vehicle_model   VARCHAR(100),
    vehicle_year    SMALLINT,
    vehicle_color   VARCHAR(50),
    vehicle_type    VARCHAR(50),
    -- Vol
    theft_date      DATE NOT NULL,
    theft_location  VARCHAR(200) NOT NULL,
    theft_department VARCHAR(50),
    -- Propriétaire
    owner_niu       VARCHAR(20),
    owner_name      VARCHAR(200),
    owner_phone     VARCHAR(50),
    -- Lien FOVeS/SIV
    foves_record_id UUID,
    -- Statut
    status          VARCHAR(20) DEFAULT 'STOLEN'
                    CHECK (status IN ('STOLEN','RECOVERED','CANCELLED')),
    recovered_date  DATE,
    recovered_location VARCHAR(200),
    -- Entrant
    entering_agency VARCHAR(100) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- ──────────────────────────────────────────────────
-- 10. ARMES À FEU VOLÉES (BIE-ARM)
-- ──────────────────────────────────────────────────
CREATE TABLE bie_stolen_firearms (
    record_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_number   VARCHAR(50) UNIQUE NOT NULL,
    serial_number   VARCHAR(100) UNIQUE NOT NULL,
    make            VARCHAR(100),
    model           VARCHAR(100),
    caliber         VARCHAR(50),
    firearm_type    VARCHAR(50),                    -- PISTOL, RIFLE, SHOTGUN, etc.
    barrel_length   DECIMAL(5,2),
    -- Vol
    theft_date      DATE NOT NULL,
    theft_location  VARCHAR(200),
    owner_niu       VARCHAR(20),
    -- Statut
    status          VARCHAR(20) DEFAULT 'STOLEN'
                    CHECK (status IN ('STOLEN','RECOVERED','CANCELLED')),
    recovered_date  DATE,
    entering_agency VARCHAR(100) NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- ──────────────────────────────────────────────────
-- 11. DOCUMENTS VOLÉS (BIE-DOC)
-- ──────────────────────────────────────────────────
CREATE TABLE bie_stolen_documents (
    record_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_number   VARCHAR(50) UNIQUE NOT NULL,
    document_type   VARCHAR(50) NOT NULL
                    CHECK (document_type IN ('PASSPORT','CIN','ACTE_NAISSANCE','PERMIS_CONDUIRE','TITRE_FONCIER','AUTRE')),
    document_number VARCHAR(100),
    issuing_agency  VARCHAR(100),
    issue_date      DATE,
    expiry_date     DATE,
    -- Propriétaire
    owner_niu       VARCHAR(20),
    owner_name      VARCHAR(200),
    -- Vol / Perte
    report_date     DATE NOT NULL,
    report_location VARCHAR(200),
    theft_type      VARCHAR(20)                     -- STOLEN, LOST, FORGED
                    CHECK (theft_type IN ('STOLEN','LOST','FORGED')),
    -- Statut
    status          VARCHAR(20) DEFAULT 'ACTIVE'
                    CHECK (status IN ('ACTIVE','RECOVERED','CANCELLED')),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- ──────────────────────────────────────────────────
-- 12. EMBARCATIONS VOLÉES (BIE-EMB) — Critique pour Haïti
-- ──────────────────────────────────────────────────
CREATE TABLE bie_stolen_vessels (
    record_id       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    record_number   VARCHAR(50) UNIQUE NOT NULL,
    vessel_name     VARCHAR(200),
    registration_number VARCHAR(100),
    hull_id_number  VARCHAR(50),                    -- HIN
    vessel_type     VARCHAR(50),                    -- FISHING, PLEASURE, FERRY, etc.
    vessel_make     VARCHAR(100),
    vessel_length_m DECIMAL(6,2),
    hull_color      VARCHAR(50),
    -- Port / Zone
    home_port       VARCHAR(200),
    theft_location  VARCHAR(200) NOT NULL,
    theft_date      DATE NOT NULL,
    -- Propriétaire
    owner_niu       VARCHAR(20),
    owner_name      VARCHAR(200),
    -- Statut
    status          VARCHAR(20) DEFAULT 'STOLEN',
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- ──────────────────────────────────────────────────
-- 13. AUDIT LOG FORENSIQUE (immuable)
-- ──────────────────────────────────────────────────
CREATE TABLE bio_audit_log (
    log_id          UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    event_type      VARCHAR(100) NOT NULL,
    table_name      VARCHAR(100),
    record_id       UUID,
    officer_niu     VARCHAR(20) NOT NULL,
    agency_code     VARCHAR(50) NOT NULL,
    purpose         VARCHAR(200) NOT NULL,
    case_number     VARCHAR(100),
    ip_hash         VARCHAR(64),                    -- SHA-256 de l'IP
    action          VARCHAR(20)                     -- CREATE, READ, UPDATE, DELETE
                    CHECK (action IN ('CREATE','READ','UPDATE','DELETE','SEARCH','HIT')),
    details         JSONB,
    signature       TEXT,                           -- ECDSA-P256 du log
    created_at      TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Partitionnement par mois
CREATE TABLE bio_audit_log_2026_06 PARTITION OF bio_audit_log
    FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');

-- Row-Level Security
ALTER TABLE bio_str_profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE per_wanted_persons ENABLE ROW LEVEL SECURITY;
ALTER TABLE per_gang_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE bio_identity_links ENABLE ROW LEVEL SECURITY;

-- Politique : seuls les agents DCPJ-DIR peuvent accéder aux liens identité
CREATE POLICY bio_identity_links_policy ON bio_identity_links
    USING (current_user = 'snisid_dcpj_director' OR current_user = 'snisid_admin');
```
