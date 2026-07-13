# MP-43 — DIPE-HT
## Registre National des Personnes Disparues d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-43 | Code : DIPE-HT
Dépendances      : SNISID-BIO-ADN, AFIS-HT (MP-19), RVIN-HT (MP-50), SIVC-HT (MP-18)
Normes           : INTERPOL Missing Persons Notice, Convention CEDAW, CDE (enfants)
Acteurs          : PNH, Parquet, IBESR, Croix-Rouge Haïtienne, OIM, MSPP
```

---

## 1. CONTEXTE

En Haïti, les disparitions de personnes sont endémiques. Contextes principaux :
- **Kidnapping** : Disparition soudaine suivie d'une demande de rançon
- **Catastrophes naturelles** : Séismes 2010, 2021 — milliers de disparus non identifiés
- **Migration** : Personnes portées disparues lors de traversées clandestines
- **Violence de gang** : Exécutions extrajudiciaires — corps jamais retrouvés
- **Enfants** : Enlèvements, recrutements forcés par gangs, traite

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE dipe_case_type AS ENUM (
    'KIDNAPPING_SUSPECTED',
    'VOLUNTARY_DISAPPEARANCE',
    'DISASTER_RELATED',
    'GANG_VIOLENCE',
    'MIGRATION_RELATED',
    'CHILD_ABDUCTION',
    'TRAFFICKING_SUSPECTED',
    'UNKNOWN'
);

CREATE TYPE dipe_case_status AS ENUM (
    'OPEN', 'LOCATED_ALIVE', 'BODY_IDENTIFIED',
    'BODY_UNIDENTIFIED', 'CANCELLED', 'COLD_CASE'
);

CREATE TABLE dipe_missing_persons (
    case_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_dipe_id    VARCHAR(25) UNIQUE NOT NULL,  -- DIPE-HT-AAAA-NNNNNN
    case_type           dipe_case_type NOT NULL,
    status              dipe_case_status NOT NULL DEFAULT 'OPEN',

    -- Identite personne disparue
    snisid_person_id    UUID,
    full_name           VARCHAR(200) NOT NULL,
    aliases             TEXT[] DEFAULT '{}',
    dob                 DATE,
    gender              VARCHAR(10),
    nationality         CHAR(3) DEFAULT 'HTI',
    occupation          VARCHAR(100),
    photo_refs          TEXT[] DEFAULT '{}',

    -- Description physique
    height_cm           SMALLINT,
    weight_kg           SMALLINT,
    skin_tone           VARCHAR(30),
    eye_color           VARCHAR(30),
    hair_color          VARCHAR(30),
    distinguishing_marks TEXT,
    clothing_last_seen  TEXT,

    -- Circonstances de la disparition
    last_seen_date      TIMESTAMPTZ NOT NULL,
    last_seen_location  VARCHAR(300),
    last_seen_dept_code CHAR(2),
    last_seen_commune   VARCHAR(100),
    last_seen_lat       DECIMAL(10,7),
    last_seen_lng       DECIMAL(10,7),
    circumstances       TEXT,

    -- Liens criminels eventuels
    sivc_alert_id       UUID,            -- Vehicule kidnapping SIVC-HT
    gang_id             UUID,
    extors_case_id      UUID,            -- Lien rançon EXTORS-HT

    -- Famille / signalement
    reported_by_name    VARCHAR(200),
    reported_by_phone   VARCHAR(30),
    reported_by_snisid  UUID,
    report_date         TIMESTAMPTZ NOT NULL,
    reporting_unit      VARCHAR(50),

    -- Biometrie (si disponible)
    afis_subject_id     UUID,
    dna_sample_ref      VARCHAR(100),
    dna_profile_id      UUID,

    -- INTERPOL
    interpol_notice_ref VARCHAR(50),
    ncmec_ref           VARCHAR(50),     -- Pour enfants

    -- Resolution
    resolution_date     TIMESTAMPTZ,
    resolution_notes    TEXT,
    rvin_case_id        UUID,            -- Lien RVIN si corps non identifie

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE dipe_sightings (
    sighting_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id             UUID NOT NULL REFERENCES dipe_missing_persons(case_id),
    sighting_date       TIMESTAMPTZ NOT NULL,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    reported_by         UUID,
    report_method       VARCHAR(30),    -- TIP_LINE, LAPI, FIELD_OFFICER, PUBLIC
    confidence          SMALLINT,
    photo_ref           VARCHAR(500),
    verified            BOOLEAN DEFAULT FALSE,
    verified_by         UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE dipe_disaster_missing (
    disaster_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id             UUID NOT NULL REFERENCES dipe_missing_persons(case_id),
    disaster_type       VARCHAR(30) NOT NULL,  -- EARTHQUAKE, HURRICANE, FLOOD
    disaster_name       VARCHAR(100),
    disaster_date       DATE NOT NULL,
    last_known_address  TEXT,
    shelter_checked     TEXT[] DEFAULT '{}',
    hospital_checked    TEXT[] DEFAULT '{}',
    morgue_checked      TEXT[] DEFAULT '{}',
    rc_haiti_ref        VARCHAR(50),           -- Reference Croix-Rouge Haïtienne
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dipe_status      ON dipe_missing_persons(status, last_seen_date DESC);
CREATE INDEX idx_dipe_type        ON dipe_missing_persons(case_type) WHERE status = 'OPEN';
CREATE INDEX idx_dipe_dept        ON dipe_missing_persons(last_seen_dept_code) WHERE status = 'OPEN';
CREATE INDEX idx_dipe_person      ON dipe_missing_persons(snisid_person_id) WHERE snisid_person_id IS NOT NULL;
CREATE INDEX idx_dipe_sightings   ON dipe_sightings(case_id, sighting_date DESC);
CREATE INDEX idx_dipe_name_fts    ON dipe_missing_persons
    USING gin(to_tsvector('simple', full_name));

COMMIT;
```

---

## 3. SERVICE GO CLÉ

```go
package service

import (
    "context"
    "github.com/snisid/dipe-svc/internal/domain"
)

// MatchWithRVIN tente de faire correspondre un disparu avec RVIN (corps non identifies)
func (s *MissingPersonService) MatchWithRVIN(
    ctx context.Context,
    caseID string,
) (*domain.MatchResult, error) {
    missing, err := s.repo.FindByID(ctx, caseID)
    if err != nil {
        return nil, err
    }

    // 1. Matching biometrique ADN si disponible
    if missing.DNAProfileID != "" {
        dnaMatches, _ := s.bioADNClient.SearchUnidentifiedBodies(ctx, missing.DNAProfileID)
        if len(dnaMatches) > 0 {
            return &domain.MatchResult{
                CaseID:      caseID,
                MatchType:   "DNA",
                Confidence:  99.9,
                RVINCaseID:  dnaMatches[0].RVINCaseID,
            }, nil
        }
    }

    // 2. Matching AFIS (empreintes)
    if missing.AFISSubjectID != "" {
        afisMatches, _ := s.afisClient.SearchLatentAgainstUnidentified(ctx, missing.AFISSubjectID)
        if len(afisMatches) > 0 && afisMatches[0].Score >= 0.90 {
            return &domain.MatchResult{
                CaseID:     caseID,
                MatchType:  "AFIS",
                Confidence: afisMatches[0].Score * 100,
                RVINCaseID: afisMatches[0].RVINCaseID,
            }, nil
        }
    }

    // 3. Matching morphologique (hauteur, poids, signes distinctifs)
    morphMatches, _ := s.rvinClient.SearchByMorphology(ctx, domain.MorphQuery{
        HeightCM: missing.HeightCM,
        WeightKG: missing.WeightKG,
        Gender:   missing.Gender,
        SkinTone: missing.SkinTone,
        DeptCode: missing.LastSeenDeptCode,
    })
    return &domain.MatchResult{
        CaseID:       caseID,
        MatchType:    "MORPHOLOGICAL",
        Candidates:   morphMatches,
    }, nil
}
```

---

## 4. API REST

| Méthode | Endpoint                             | Rôle           | Description                      |
|---------|--------------------------------------|----------------|----------------------------------|
| `POST`  | `/api/v1/dipe/cases`                 | PNH, PUBLIC    | Signaler disparition             |
| `GET`   | `/api/v1/dipe/cases/:id`             | PNH, PARQUET   | Détail dossier                   |
| `GET`   | `/api/v1/dipe/cases/open`            | PNH, DCPJ      | Tous dossiers ouverts            |
| `POST`  | `/api/v1/dipe/cases/:id/sightings`   | PNH, PUBLIC    | Signaler observation             |
| `GET`   | `/api/v1/dipe/match/rvin/:id`        | PNH, PARQUET   | Matching avec corps non identifiés|
| `PATCH` | `/api/v1/dipe/cases/:id/resolve`     | PNH_OFFICER    | Résoudre le dossier              |
| `GET`   | `/api/v1/dipe/stats/by-type`         | DCPJ_ADMIN     | Stats par type                   |
| `GET`   | `/api/v1/dipe/hotline/tips`          | PNH_OPERATOR   | Signalements ligne directe       |

---

## 5. VARIABLES D'ENVIRONNEMENT

```dotenv
DIPE_DB_HOST=localhost
DIPE_DB_NAME=snisid_dipe
DIPE_AFIS_SERVICE_URL=http://afis-svc:8091
DIPE_BIO_ADN_SERVICE_URL=http://bio-adn-svc:8080
DIPE_RVIN_SERVICE_URL=http://rvin-svc:8120
DIPE_INTERPOL_NOTICE_URL=https://i247-gateway.pnh.gov.ht/notices
DIPE_NCMEC_API_URL=https://api.missingkids.org
DIPE_HOTLINE_NUMBER=116
DIPE_SERVICE_PORT=8118
```

---
*MP-43 — DIPE-HT — Personnes Disparues — SNISID — République d'Haïti*

---
---

# MP-44 — TRAIT-HT
## Traite des Personnes, Migration Irrégulière et Passeurs
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-44 | Code : TRAIT-HT
Dépendances      : MAR-HT (MP-34), SIFR-HT (MP-33), ENFL-HT (MP-45), DIPE-HT (MP-43)
Normes           : Protocole de Palerme (ONU), IOM Counter-Trafficking Standards
Acteurs          : PNH, OIM, BID, IBESR, Parquet, BCPE (Brigade Contrôle Frontières)
```

---

## 1. CONTEXTE

Haïti est à la fois pays source, transit et destination de traite des personnes :
- **Restaveks** : Enfants placés comme domestiques (forme d'esclavage moderne)
- **Migration forcée vers DOM** : Travailleurs haïtiens exploités en République Dominicaine
- **Traversée maritime** : Embarcations de fortune vers Bahamas/USA — 200-600 USD/personne
- **Recrutement gang** : Jeunes enrôlés de force par les gangs
- **Prostitution forcée** : Réseaux documentés dans zones de conflit

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE trait_type AS ENUM (
    'LABOR_EXPLOITATION', 'SEXUAL_EXPLOITATION', 'FORCED_MARRIAGE',
    'CHILD_DOMESTIC_SERVITUDE', 'GANG_RECRUITMENT_FORCED',
    'IRREGULAR_MIGRATION_FACILITATION', 'ORGAN_TRAFFICKING', 'OTHER'
);

CREATE TYPE trait_victim_status AS ENUM (
    'IDENTIFIED_VICTIM', 'POTENTIAL_VICTIM', 'WITNESS',
    'RESCUED', 'REPATRIATED', 'DECEASED', 'MISSING'
);

CREATE TABLE trait_cases (
    case_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_trait_id   VARCHAR(25) UNIQUE NOT NULL,  -- TRAIT-HT-AAAA-NNNNNN
    trait_type          trait_type NOT NULL,
    status              VARCHAR(20) DEFAULT 'OPEN',
    victim_count        SMALLINT DEFAULT 1,
    minor_count         SMALLINT DEFAULT 0,
    origin_country      CHAR(3) DEFAULT 'HTI',
    transit_countries   CHAR(3)[] DEFAULT '{}',
    destination_country CHAR(3),
    route_description   TEXT,
    transport_mode      TEXT[] DEFAULT '{}',     -- BOAT, BUS, FOOT, AIR
    mar_incident_id     UUID,                    -- Lien MAR-HT si maritime
    sifr_crossing_ids   UUID[] DEFAULT '{}',     -- Postes frontiere impliques
    gang_id             UUID,                    -- Si gang facilite
    recruiter_ids       UUID[] DEFAULT '{}',     -- Passeurs / recruteurs SNISID
    total_amount_paid   DECIMAL(12,2),
    amount_per_person   DECIMAL(10,2),
    currency            CHAR(3) DEFAULT 'USD',
    investigating_unit  VARCHAR(50),
    case_reference      VARCHAR(100),
    iom_case_ref        VARCHAR(50),
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE trait_victims (
    victim_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id             UUID NOT NULL REFERENCES trait_cases(case_id),
    snisid_person_id    UUID,
    victim_status       trait_victim_status NOT NULL,
    full_name           VARCHAR(200),
    nationality         CHAR(3) DEFAULT 'HTI',
    dob                 DATE,
    gender              VARCHAR(10),
    is_minor            BOOLEAN DEFAULT FALSE,
    exploitation_type   TEXT,
    rescue_date         TIMESTAMPTZ,
    rescue_location     VARCHAR(300),
    current_location    TEXT,
    assistance_provided TEXT[] DEFAULT '{}',  -- SHELTER, LEGAL, MEDICAL, REPATRIATION
    dipe_case_id        UUID,                 -- Lien DIPE si disparu initialement
    afis_subject_id     UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE trait_networks (
    network_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    network_name        VARCHAR(150),
    primary_route       TEXT,
    origin_dept         CHAR(2),
    known_members       UUID[] DEFAULT '{}',   -- SNISID person IDs
    gang_affiliations   UUID[] DEFAULT '{}',
    monthly_volume_est  INTEGER,
    fee_per_person_usd  DECIMAL(10,2),
    is_active           BOOLEAN DEFAULT TRUE,
    intel_confidence    SMALLINT,
    linked_cases        UUID[] DEFAULT '{}',
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_trait_cases_type    ON trait_cases(trait_type, status);
CREATE INDEX idx_trait_cases_gang    ON trait_cases(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_trait_victims_case  ON trait_victims(case_id);
CREATE INDEX idx_trait_victims_minor ON trait_victims(is_minor) WHERE is_minor = TRUE;
CREATE INDEX idx_trait_victims_snisid ON trait_victims(snisid_person_id) WHERE snisid_person_id IS NOT NULL;
CREATE INDEX idx_trait_networks_active ON trait_networks(is_active) WHERE is_active = TRUE;

COMMIT;
```

---

## 3. API REST

| Méthode | Endpoint                              | Rôle         | Description                     |
|---------|---------------------------------------|--------------|---------------------------------|
| `POST`  | `/api/v1/trait/cases`                 | PNH, OIM     | Ouvrir dossier traite           |
| `POST`  | `/api/v1/trait/cases/:id/victims`     | PNH, OIM     | Enregistrer victime             |
| `GET`   | `/api/v1/trait/cases/:id`             | PNH, PARQUET | Détail dossier                  |
| `GET`   | `/api/v1/trait/victims/minors`        | IBESR, PNH   | Victimes mineures               |
| `POST`  | `/api/v1/trait/networks`              | DCPJ_INTEL   | Documenter réseau passeurs      |
| `GET`   | `/api/v1/trait/stats/by-type`         | PNH_ADMIN    | Stats par type traite           |
| `GET`   | `/api/v1/trait/cases/maritime`        | GCH, PNH     | Cas liés à migration maritime   |

---

## 4. VARIABLES D'ENVIRONNEMENT

```dotenv
TRAIT_DB_HOST=localhost
TRAIT_DB_NAME=snisid_trait
TRAIT_IOM_API_URL=https://api.iom.int/counter-trafficking
TRAIT_IOM_API_KEY=<VAULT:trait/iom_api_key>
TRAIT_MAR_SERVICE_URL=http://mar-svc:8107
TRAIT_DIPE_SERVICE_URL=http://dipe-svc:8118
TRAIT_ENFL_SERVICE_URL=http://enfl-svc:8119
TRAIT_SERVICE_PORT=8122
```

---
*MP-44 — TRAIT-HT — Traite des Personnes — SNISID — République d'Haïti*
