# MP-45 — ENFL-HT
## Registre National des Enfants Disparus, à Risque et Exploités
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-45 | Code : ENFL-HT
Dépendances      : DIPE-HT (MP-43), TRAIT-HT (MP-44), GANG-HT (MP-24), FIR-HT (MP-20)
Normes           : Convention des Droits de l'Enfant (CDE), INTERPOL ICSE, NCMEC
Acteurs          : IBESR, PNH, Parquet Mineurs, UNICEF, Save the Children, OIM
```

---

## 1. CONTEXTE

Les enfants haïtiens font face à des risques multiples documentés :
- **Restaveks** : 225,000+ enfants en servitude domestique (données UNICEF 2022)
- **Recrutement gang** : Enfants recrutés dès 10-12 ans par les gangs
- **Enlèvements ciblés** : Enfants de familles aisées pris en otage
- **Exploitation sexuelle** : Documentée dans camps de déplacés et zones de conflit
- **Séparation familiale** : Catastrophes naturelles et déplacements internes

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE enfl_risk_category AS ENUM (
    'MISSING_ABDUCTION',
    'GANG_RECRUITMENT',
    'DOMESTIC_SERVITUDE_RESTAVEK',
    'SEXUAL_EXPLOITATION',
    'TRAFFICKING',
    'UNACCOMPANIED_MIGRANT',
    'SEPARATED_DISASTER',
    'STREET_CHILD',
    'OTHER'
);

CREATE TYPE enfl_status AS ENUM (
    'AT_RISK', 'MISSING', 'LOCATED_SAFE', 'LOCATED_AT_RISK',
    'IN_CARE', 'REPATRIATED', 'DECEASED'
);

CREATE TABLE enfl_children (
    child_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_enfl_id    VARCHAR(25) UNIQUE NOT NULL,  -- ENFL-HT-AAAA-NNNNNN
    snisid_person_id    UUID,
    dipe_case_id        UUID,           -- Lien DIPE-HT si disparu
    trait_case_id       UUID,           -- Lien TRAIT-HT si traite
    risk_category       enfl_risk_category NOT NULL,
    status              enfl_status NOT NULL DEFAULT 'MISSING',
    full_name           VARCHAR(200) NOT NULL,
    dob                 DATE NOT NULL,
    age_at_registration SMALLINT,
    gender              VARCHAR(10),
    nationality         CHAR(3) DEFAULT 'HTI',
    photo_refs          TEXT[] DEFAULT '{}',
    distinguishing_marks TEXT,
    height_cm           SMALLINT,
    skin_tone           VARCHAR(30),
    guardian_name       VARCHAR(200),
    guardian_phone      VARCHAR(30),
    guardian_snisid_id  UUID,
    last_known_location VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    disappearance_date  TIMESTAMPTZ,
    gang_id             UUID,           -- Si gang implique
    recruiter_snisid_id UUID,           -- Si recruteur identifie
    afis_subject_id     UUID,
    dna_profile_id      UUID,
    interpol_icse_ref   VARCHAR(50),    -- INTERPOL Crimes Against Children
    ncmec_ref           VARCHAR(50),
    ibesr_ref           VARCHAR(50),
    assistance_type     TEXT[] DEFAULT '{}',
    current_shelter     VARCHAR(200),
    assigned_caseworker UUID,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE enfl_restaveks (
    restavek_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id            UUID NOT NULL REFERENCES enfl_children(child_id),
    employing_household VARCHAR(300),
    household_dept      CHAR(2),
    household_commune   VARCHAR(100),
    employing_person_id UUID,
    reported_conditions TEXT,
    school_attendance   BOOLEAN DEFAULT FALSE,
    ibesr_inspection    BOOLEAN DEFAULT FALSE,
    last_inspection_date DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_enfl_status      ON enfl_children(status, risk_category);
CREATE INDEX idx_enfl_dept        ON enfl_children(dept_code) WHERE status IN ('MISSING','AT_RISK');
CREATE INDEX idx_enfl_gang        ON enfl_children(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_enfl_name_fts    ON enfl_children USING gin(to_tsvector('simple', full_name));

COMMIT;
```

---

## 3. API REST

| Méthode | Endpoint                           | Rôle          | Description                    |
|---------|------------------------------------|---------------|--------------------------------|
| `POST`  | `/api/v1/enfl/children`            | IBESR, PNH    | Enregistrer enfant à risque    |
| `GET`   | `/api/v1/enfl/children/:id`        | IBESR, PNH    | Profil complet enfant          |
| `GET`   | `/api/v1/enfl/missing`             | PNH, IBESR    | Enfants disparus actifs        |
| `GET`   | `/api/v1/enfl/restaveks`           | IBESR         | Restaveks enregistrés          |
| `POST`  | `/api/v1/enfl/children/:id/locate` | PNH_OFFICER   | Signaler localisation          |
| `GET`   | `/api/v1/enfl/gang-recruited`      | DCPJ, IBESR   | Enfants recrutés par gangs     |

## 4. VARIABLES D'ENVIRONNEMENT

```dotenv
ENFL_DB_HOST=localhost
ENFL_DB_NAME=snisid_enfl
ENFL_IBESR_API_URL=https://api.ibesr.gov.ht
ENFL_NCMEC_API_URL=https://api.missingkids.org
ENFL_INTERPOL_ICSE_URL=https://i247-gateway.pnh.gov.ht/icse
ENFL_DIPE_SERVICE_URL=http://dipe-svc:8118
ENFL_SERVICE_PORT=8119
```

---
*MP-45 — ENFL-HT — Enfants Disparus et à Risque — SNISID — République d'Haïti*

---
---

# MP-46 — DPIDE-HT
## Registre et Suivi des Déplacés Internes (IDPs)
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-46 | Code : DPIDE-HT
Dépendances      : SIGDC-HT (MP-49), SIGEO-HT (MP-48), SNISID Identité, ENFL-HT (MP-45)
Normes           : Principes directeurs ONU IDPs, Sphère Standards, IOM DTM
Acteurs          : OIM, CSPAN (Coordination Sécurité et Protection), MSPP, OCHA
```

---

## 1. CONTEXTE

Haïti compte 580,000+ déplacés internes (IDPs) selon l'OIM (données 2024) :
- 360,000+ liés aux violences de gangs (Port-au-Prince métropole)
- 170,000+ liés au séisme du 14 août 2021 (Grand-Sud)
- Camps de déplacés : Champ-de-Mars, Stade Sylvio Cator, provinces

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE dpide_displacement_cause AS ENUM (
    'GANG_VIOLENCE', 'EARTHQUAKE', 'HURRICANE',
    'FLOOD', 'FIRE', 'POLITICAL_VIOLENCE', 'OTHER'
);

CREATE TYPE dpide_idp_status AS ENUM (
    'DISPLACED', 'IN_CAMP', 'WITH_HOST_FAMILY',
    'RELOCATED', 'RETURNED_HOME', 'EMIGRATED', 'DECEASED'
);

CREATE TABLE dpide_idps (
    idp_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_dpide_id   VARCHAR(25) UNIQUE NOT NULL,  -- DPIDE-HT-AAAA-NNNNNN
    snisid_person_id    UUID,
    full_name           VARCHAR(200) NOT NULL,
    dob                 DATE,
    gender              VARCHAR(10),
    household_size      SMALLINT DEFAULT 1,
    minors_count        SMALLINT DEFAULT 0,
    displacement_cause  dpide_displacement_cause NOT NULL,
    displacement_date   TIMESTAMPTZ NOT NULL,
    origin_address      TEXT,
    origin_dept_code    CHAR(2) NOT NULL,
    origin_commune      VARCHAR(100),
    status              dpide_idp_status NOT NULL DEFAULT 'DISPLACED',
    current_location    TEXT,
    current_dept_code   CHAR(2),
    current_commune     VARCHAR(100),
    current_lat         DECIMAL(10,7),
    current_lng         DECIMAL(10,7),
    camp_id             UUID,
    shelter_type        VARCHAR(50),   -- CAMP, HOST_FAMILY, RENTED, SPONTANEOUS
    has_nfi             BOOLEAN DEFAULT FALSE,      -- Non-Food Items
    receives_food_aid   BOOLEAN DEFAULT FALSE,
    has_latrines        BOOLEAN DEFAULT FALSE,
    has_water_access    BOOLEAN DEFAULT FALSE,
    medical_needs       TEXT[] DEFAULT '{}',
    iom_dtm_ref         VARCHAR(50),
    ocha_ref            VARCHAR(50),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE dpide_camps (
    camp_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    camp_name           VARCHAR(150) NOT NULL,
    dept_code           CHAR(2) NOT NULL,
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    displacement_cause  dpide_displacement_cause,
    managing_org        VARCHAR(150),
    capacity            INTEGER,
    current_population  INTEGER DEFAULT 0,
    is_active           BOOLEAN DEFAULT TRUE,
    has_medical_post    BOOLEAN DEFAULT FALSE,
    has_school          BOOLEAN DEFAULT FALSE,
    water_source        TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dpide_status      ON dpide_idps(status, displacement_cause);
CREATE INDEX idx_dpide_dept        ON dpide_idps(current_dept_code) WHERE status IN ('DISPLACED','IN_CAMP');
CREATE INDEX idx_dpide_cause       ON dpide_idps(displacement_cause);
CREATE INDEX idx_dpide_camp        ON dpide_idps(camp_id) WHERE camp_id IS NOT NULL;
CREATE INDEX idx_dpide_snisid      ON dpide_idps(snisid_person_id) WHERE snisid_person_id IS NOT NULL;

COMMIT;
```

---

## 3. API REST

| Méthode | Endpoint                        | Rôle          | Description                     |
|---------|---------------------------------|---------------|---------------------------------|
| `POST`  | `/api/v1/dpide/idps`            | OIM, CSPAN    | Enregistrer déplacé             |
| `GET`   | `/api/v1/dpide/idps/:id`        | OIM, PNH      | Profil IDP complet              |
| `GET`   | `/api/v1/dpide/camps`           | OIM, OCHA     | Liste camps actifs              |
| `GET`   | `/api/v1/dpide/stats/overview`  | OCHA, MSPP    | Vue d'ensemble (counts)         |
| `PATCH` | `/api/v1/dpide/idps/:id/status` | OIM_AGENT     | Mettre à jour statut IDP        |
| `GET`   | `/api/v1/dpide/heatmap`         | OCHA          | Carte densité déplacés          |

## 4. VARIABLES D'ENVIRONNEMENT

```dotenv
DPIDE_DB_HOST=localhost
DPIDE_DB_NAME=snisid_dpide
DPIDE_IOM_DTM_API_URL=https://dtm.iom.int/api
DPIDE_OCHA_API_URL=https://api.reliefweb.int/v1
DPIDE_SIGEO_SERVICE_URL=http://sigeo-svc:8125
DPIDE_SERVICE_PORT=8121
```

---
*MP-46 — DPIDE-HT — Déplacés Internes — SNISID — République d'Haïti*

---
---

# MP-47 — VICT-HT
## Registre National des Victimes de Crimes Graves
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-47 | Code : VICT-HT
Dépendances      : FIR-HT (MP-20), SNISID-BIO-ADN, AFIS-HT (MP-19), RVIN-HT (MP-50)
Normes           : Déclaration ONU Victimes 1985, Règles Mandela, CIPD droits victimes
Acteurs          : Parquet, PNH, MSP, MSPP, ONG droits humains, Médecins Sans Frontières
```

---

## 1. CONTEXTE

Ce module documente toutes les victimes de crimes graves : homicides, viols collectifs,
massacres (Bel-Air, Lizon, Cité Soleil), tortures et disparitions forcées. Il constitue
la base de données pour les réparations futures, les poursuites judiciaires et les
rapports aux organisations internationales de droits humains (IACHR, HRC).

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE vict_crime_type AS ENUM (
    'HOMICIDE', 'MASS_KILLING', 'RAPE', 'GANG_RAPE',
    'TORTURE', 'FORCED_DISAPPEARANCE', 'EXTRAJUDICIAL_KILLING',
    'KIDNAPPING_VICTIM', 'MUTILATION', 'OTHER_GRAVE'
);

CREATE TYPE vict_victim_status AS ENUM (
    'ALIVE_SURVIVOR', 'DECEASED_IDENTIFIED',
    'DECEASED_UNIDENTIFIED', 'MISSING_PRESUMED_DEAD'
);

CREATE TABLE vict_victims (
    victim_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_vict_id    VARCHAR(25) UNIQUE NOT NULL,  -- VICT-HT-AAAA-NNNNNN
    snisid_person_id    UUID,
    crime_type          vict_crime_type NOT NULL,
    victim_status       vict_victim_status NOT NULL,
    full_name           VARCHAR(200),
    dob                 DATE,
    gender              VARCHAR(10),
    nationality         CHAR(3) DEFAULT 'HTI',
    occupation          VARCHAR(100),
    incident_date       TIMESTAMPTZ NOT NULL,
    incident_location   VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    perpetrator_ids     UUID[] DEFAULT '{}',
    gang_id             UUID,
    case_reference      VARCHAR(100),
    parquet_ref         VARCHAR(100),
    medical_report_ref  VARCHAR(200),
    autopsy_ref         VARCHAR(200),
    dna_sample_ref      VARCHAR(100),
    afis_subject_id     UUID,
    rvin_case_id        UUID,
    iachr_ref           VARCHAR(50),       -- CIDH/IACHR reference
    un_special_rap_ref  VARCHAR(50),
    needs_reparation    BOOLEAN DEFAULT FALSE,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE vict_mass_incidents (
    mass_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_name       VARCHAR(200) NOT NULL,   -- Ex: Massacre de Cité Soleil 2022
    crime_type          vict_crime_type NOT NULL,
    incident_date       TIMESTAMPTZ NOT NULL,
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    victim_count        INTEGER NOT NULL,
    survivor_count      INTEGER DEFAULT 0,
    perpetrator_gang_id UUID,
    description         TEXT,
    documented_by       TEXT[] DEFAULT '{}',    -- RNDDH, HRW, MSF, ONU, etc.
    iachr_case_ref      VARCHAR(50),
    linked_victim_ids   UUID[] DEFAULT '{}',
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vict_crime_type ON vict_victims(crime_type, victim_status);
CREATE INDEX idx_vict_dept       ON vict_victims(dept_code, incident_date DESC);
CREATE INDEX idx_vict_gang       ON vict_victims(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_vict_snisid     ON vict_victims(snisid_person_id) WHERE snisid_person_id IS NOT NULL;
CREATE INDEX idx_vict_mass       ON vict_mass_incidents(dept_code, incident_date DESC);

COMMIT;
```

---

## 3. API REST

| Méthode | Endpoint                           | Rôle          | Description                    |
|---------|------------------------------------|---------------|--------------------------------|
| `POST`  | `/api/v1/vict/victims`             | PNH, MSF      | Enregistrer victime            |
| `GET`   | `/api/v1/vict/victims/:id`         | PARQUET, PNH  | Profil victime complet         |
| `POST`  | `/api/v1/vict/mass-incidents`      | DCPJ, PARQUET | Documenter massacre            |
| `GET`   | `/api/v1/vict/mass-incidents`      | PARQUET, IACHR| Liste massacres documentés     |
| `GET`   | `/api/v1/vict/by-gang/:id`         | DCPJ, PARQUET | Victimes d'un gang             |
| `GET`   | `/api/v1/vict/stats/by-type`       | DCPJ_ADMIN    | Stats par type de crime        |
| `GET`   | `/api/v1/vict/reparation-list`     | MJSP, PARQUET | Victimes nécessitant réparation|

## 4. VARIABLES D'ENVIRONNEMENT

```dotenv
VICT_DB_HOST=localhost
VICT_DB_NAME=snisid_vict
VICT_RVIN_SERVICE_URL=http://rvin-svc:8120
VICT_BIO_ADN_SERVICE_URL=http://bio-adn-svc:8080
VICT_AFIS_SERVICE_URL=http://afis-svc:8091
VICT_KAFKA_BROKERS=kafka:9092
VICT_SERVICE_PORT=8123
```

---
*MP-47 — VICT-HT — Registre Victimes — SNISID — République d'Haïti*
