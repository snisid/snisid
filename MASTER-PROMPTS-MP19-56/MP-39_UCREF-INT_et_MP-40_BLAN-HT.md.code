# MP-39 — UCREF-INT
## Interface SNISID / UCREF — Renseignement Financier National
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-39 | Code : UCREF-INT
Dépendances      : BLAN-HT (MP-40), GANG-HT (MP-24), SANC-HT (MP-27), FIR-HT (MP-20)
Normes           : GAFI/FATF 40 Recommandations, Egmont Group FIU Standards
Acteurs          : UCREF (Unité Centrale de Renseignement Financier), BRH, Banques, MSS
```

---

## 1. CONTEXTE

L'UCREF est l'Unité de Renseignement Financier d'Haïti (Financial Intelligence Unit).
Ses rapports documentent : fausse facturation d'importations, smurfing via MonCash,
acquisitions immobilières suspectes par des proches de chefs de gangs, transferts
de rançons via des comptes de prête-noms. Ce module crée l'interface bidirectionnelle
SNISID ↔ UCREF pour croiser les identités criminelles avec les flux financiers.

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE ucref_str_status AS ENUM (
    'RECEIVED','UNDER_ANALYSIS','DISSEMINATED','ARCHIVED','NO_ACTION'
);

CREATE TYPE ucref_report_type AS ENUM (
    'STR',          -- Suspicious Transaction Report
    'CTR',          -- Cash Transaction Report (> HTG 500,000)
    'INTERNATIONAL_WIRE',
    'REAL_ESTATE',
    'MONCASH_PATTERN',
    'CRYPTO_PATTERN'
);

CREATE TABLE ucref_str_reports (
    str_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_str_id     VARCHAR(25) UNIQUE NOT NULL,   -- STR-HT-AAAA-NNNNNN
    report_type         ucref_report_type NOT NULL,
    status              ucref_str_status NOT NULL DEFAULT 'RECEIVED',
    reporting_institution VARCHAR(200) NOT NULL,
    institution_type    VARCHAR(30),    -- BANK, MSB, MONCASH, INSURANCE, CASINO
    report_date         TIMESTAMPTZ NOT NULL,
    transaction_date    TIMESTAMPTZ,
    transaction_amount  DECIMAL(18,2),
    transaction_currency CHAR(3) DEFAULT 'HTG',
    transaction_amount_usd DECIMAL(18,2),

    -- Personnes impliquees
    subject_snisid_ids  UUID[] DEFAULT '{}',
    subject_names       TEXT[] DEFAULT '{}',
    subject_accounts    TEXT[] DEFAULT '{}',

    -- Contexte
    suspicious_activity TEXT NOT NULL,
    ml_typology         VARCHAR(100),   -- Smurfing, Trade-Based ML, Real Estate, etc.
    predicate_crime     VARCHAR(100),   -- Crime sous-jacent suspecte

    -- Liens criminels SNISID
    gang_id             UUID,           -- Si lien gang identifie
    fpr_person_ids      UUID[] DEFAULT '{}',
    sanc_match_ids      UUID[] DEFAULT '{}',

    -- Traitement
    analyst_id          UUID,
    analysis_notes      TEXT,
    disseminated_to     TEXT[] DEFAULT '{}',   -- DCPJ, MJSP, PARQUET, etc.
    disseminated_at     TIMESTAMPTZ,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE ucref_financial_profiles (
    profile_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snisid_person_id    UUID NOT NULL UNIQUE,
    total_str_count     INTEGER DEFAULT 0,
    total_ctr_count     INTEGER DEFAULT 0,
    estimated_illegal_assets_usd DECIMAL(18,2),
    known_accounts      JSONB,          -- [{institution, account_type, country}]
    known_properties    JSONB,          -- [{address, value, acquisition_date}]
    known_businesses    TEXT[] DEFAULT '{}',
    ml_risk_score       SMALLINT CHECK (ml_risk_score BETWEEN 0 AND 100),
    is_pep              BOOLEAN DEFAULT FALSE,   -- Politically Exposed Person
    last_updated        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE ucref_moncash_patterns (
    pattern_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    str_id              UUID REFERENCES ucref_str_reports(str_id),
    phone_number        VARCHAR(20) NOT NULL,    -- Numero MonCash
    snisid_person_id    UUID,
    pattern_type        VARCHAR(50),    -- STRUCTURING, RAPID_TRANSFERS, RANSOM_RECEIPT
    transaction_count   INTEGER,
    total_amount_htg    DECIMAL(18,2),
    period_start        TIMESTAMPTZ,
    period_end          TIMESTAMPTZ,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ucref_str_status  ON ucref_str_reports(status, report_date DESC);
CREATE INDEX idx_ucref_str_gang    ON ucref_str_reports(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_ucref_str_subjects ON ucref_str_reports USING gin(subject_snisid_ids);
CREATE INDEX idx_ucref_profiles    ON ucref_financial_profiles(snisid_person_id);
CREATE INDEX idx_ucref_moncash     ON ucref_moncash_patterns(phone_number);

COMMIT;
```

---

## 3. API REST

| Méthode | Endpoint                             | Rôle         | Description                       |
|---------|--------------------------------------|--------------|-----------------------------------|
| `POST`  | `/api/v1/ucref/str`                  | UCREF_AGENT  | Soumettre déclaration de soupçon  |
| `GET`   | `/api/v1/ucref/str/:id`              | UCREF_ANALYST| Détail STR                        |
| `GET`   | `/api/v1/ucref/profile/:person_id`   | UCREF, DCPJ  | Profil financier d'une personne   |
| `POST`  | `/api/v1/ucref/moncash/pattern`      | UCREF_ANALYST| Enregistrer pattern MonCash       |
| `GET`   | `/api/v1/ucref/str/unanalyzed`       | UCREF_ANALYST| STRs en attente d'analyse         |
| `POST`  | `/api/v1/ucref/str/:id/disseminate`  | UCREF_CHIEF  | Transmettre aux autorités         |
| `GET`   | `/api/v1/ucref/gang-finances/:id`    | UCREF, DCPJ  | Profil financier d'un gang        |

---

## 4. VARIABLES D'ENVIRONNEMENT

```dotenv
UCREF_DB_HOST=localhost
UCREF_DB_NAME=snisid_ucref
UCREF_GANG_SERVICE_URL=http://gang-svc:8095
UCREF_SANC_SERVICE_URL=http://sanc-svc:8100
UCREF_FPR_SERVICE_URL=http://fpr-svc:8085
UCREF_EGMONT_API_URL=https://goaml.egmont.org/api
UCREF_FATF_COMPLIANCE_MODE=STRICT
UCREF_SERVICE_PORT=8112
```

---
*MP-39 — UCREF-INT — Renseignement Financier — SNISID — République d'Haïti*

---
---

# MP-40 — BLAN-HT
## Transactions Suspectes, Blanchiment et Économie Criminelle
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-40 | Code : BLAN-HT
Dépendances      : UCREF-INT (MP-39), GANG-HT (MP-24), EXTORS-HT (MP-41), CRYPT-HT (MP-42)
Normes           : GAFI/FATF, Recommandation 3 (infractions sous-jacentes), Groupe Egmont
Acteurs          : UCREF, BRH, MJSP, Parquet financier, Douanes
```

---

## 1. CONTEXTE

Les typologies documentées de blanchiment en Haïti selon l'UCREF :
- **Smurfing MonCash** : Fractionnement de rançons en petits virages sous le seuil déclaratif
- **Fausse facturation import** : Surévaluation de marchandises importées (matériaux construction)
- **Immobilier criminel** : Achats de propriétés par proches de chefs de gangs
- **Entreprises écran** : Sociétés de construction, commerces alimentaires, stations service
- **Transferts diaspora** : Utilisation de canaux diaspora (Western Union, CAM Transfer) pour injection

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE blan_typology AS ENUM (
    'SMURFING', 'TRADE_BASED_ML', 'REAL_ESTATE',
    'SHELL_COMPANY', 'CASH_INTENSIVE_BUSINESS',
    'CRYPTO_MIXING', 'DIASPORA_TRANSFER',
    'RANSOM_LAUNDERING', 'CORRUPTION_PROCEEDS'
);

CREATE TYPE blan_asset_type AS ENUM (
    'REAL_ESTATE', 'VEHICLE', 'BUSINESS', 'BANK_ACCOUNT',
    'CRYPTO_WALLET', 'CASH', 'JEWELRY', 'LIVESTOCK', 'OTHER'
);

CREATE TABLE blan_cases (
    case_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_blan_id    VARCHAR(25) UNIQUE NOT NULL,  -- BLAN-HT-AAAA-NNNNNN
    case_title          TEXT NOT NULL,
    typology            blan_typology NOT NULL,
    status              VARCHAR(20) DEFAULT 'OPEN',
    total_amount_usd    DECIMAL(18,2),
    predicate_crime     VARCHAR(100),
    subject_ids         UUID[] DEFAULT '{}',       -- SNISID persons
    gang_id             UUID,
    str_ids             UUID[] DEFAULT '{}',        -- Liens UCREF STRs
    opened_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    analyst_id          UUID,
    parquet_ref         VARCHAR(100),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE blan_suspicious_assets (
    asset_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id             UUID NOT NULL REFERENCES blan_cases(case_id),
    asset_type          blan_asset_type NOT NULL,
    description         TEXT NOT NULL,
    address             TEXT,
    dept_code           CHAR(2),
    estimated_value_usd DECIMAL(15,2),
    acquisition_date    DATE,
    owner_snisid_id     UUID,
    owner_name          VARCHAR(200),
    registered_in       CHAR(3),        -- Pays d enregistrement
    is_frozen           BOOLEAN DEFAULT FALSE,
    freeze_order_ref    VARCHAR(100),
    is_seized           BOOLEAN DEFAULT FALSE,
    seizure_date        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE blan_transaction_chains (
    chain_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id             UUID NOT NULL REFERENCES blan_cases(case_id),
    step_number         SMALLINT NOT NULL,
    transaction_type    VARCHAR(50),
    from_account        VARCHAR(200),
    from_institution    VARCHAR(150),
    to_account          VARCHAR(200),
    to_institution      VARCHAR(150),
    amount              DECIMAL(18,2),
    currency            CHAR(3),
    amount_usd          DECIMAL(18,2),
    transaction_date    TIMESTAMPTZ,
    is_suspicious_step  BOOLEAN DEFAULT TRUE,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE blan_real_estate_flagged (
    property_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id             UUID REFERENCES blan_cases(case_id),
    address             TEXT NOT NULL,
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    property_type       VARCHAR(50),
    purchase_price_usd  DECIMAL(15,2),
    purchase_date       DATE,
    declared_owner      VARCHAR(200),
    beneficial_owner_id UUID,          -- SNISID personne reelle
    suspicious_reasons  TEXT[] DEFAULT '{}',
    is_frozen           BOOLEAN DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_blan_cases_status  ON blan_cases(status, typology);
CREATE INDEX idx_blan_cases_gang    ON blan_cases(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_blan_cases_persons ON blan_cases USING gin(subject_ids);
CREATE INDEX idx_blan_assets_type   ON blan_suspicious_assets(asset_type, is_frozen);
CREATE INDEX idx_blan_real_estate   ON blan_real_estate_flagged(dept_code);

COMMIT;
```

---

## 3. API REST

| Méthode | Endpoint                              | Rôle          | Description                     |
|---------|---------------------------------------|---------------|---------------------------------|
| `POST`  | `/api/v1/blan/cases`                  | UCREF_ANALYST | Ouvrir dossier blanchiment      |
| `GET`   | `/api/v1/blan/cases/:id`              | UCREF, DCPJ   | Détail dossier                  |
| `POST`  | `/api/v1/blan/cases/:id/assets`       | UCREF_ANALYST | Ajouter actif suspect           |
| `POST`  | `/api/v1/blan/cases/:id/chain`        | UCREF_ANALYST | Documenter chaîne transactions  |
| `GET`   | `/api/v1/blan/real-estate/flagged`    | UCREF, DCPJ   | Propriétés immobilières suspectes|
| `GET`   | `/api/v1/blan/assets/frozen`          | UCREF, PARQUET| Actifs gelés                    |
| `GET`   | `/api/v1/blan/stats/by-typology`      | UCREF_ADMIN   | Stats par typologie             |

---

## 4. VARIABLES D'ENVIRONNEMENT

```dotenv
BLAN_DB_HOST=localhost
BLAN_DB_NAME=snisid_blan
BLAN_UCREF_SERVICE_URL=http://ucref-svc:8112
BLAN_GANG_SERVICE_URL=http://gang-svc:8095
BLAN_KAFKA_BROKERS=kafka:9092
BLAN_SERVICE_PORT=8115
```

---
*MP-40 — BLAN-HT — Blanchiment et Économie Criminelle — SNISID — République d'Haïti*
