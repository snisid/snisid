# MP-55 — CYBRE-HT
## Système National de Lutte contre la Cybercriminalité d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-55 | Code : CYBRE-HT
Dépendances      : CRYPT-HT (MP-42), CORR-HT (MP-53), SIPCI-HT (MP-54), UCREF-INT (MP-39)
Normes           : Convention de Budapest (Cybercrime), ITU-T X.1205, NIST CSF
Acteurs          : DCPJ Cellule Cybercriminalité, CONATEL, Parquet, MTPTC
```

---

## 1. CONTEXTE

La cybercriminalité en Haïti touche principalement :
- **Fraudes MonCash** : SIM swapping, phishing, faux agents de transfert
- **Fraudes DigiCel/Natcom** : Clonage cartes SIM, numéros Premium surtaxés
- **Piratage systèmes état** : Tentatives d'intrusion sur SNISID, registres civils
- **Arnaque en ligne** : Romance scams, faux emplois, escroqueries diaspora
- **Désinformation** : Manipulation coordonnée des réseaux sociaux liée aux gangs

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE cybre_crime_type AS ENUM (
    'MONCASH_FRAUD',
    'SIM_SWAPPING',
    'PHISHING',
    'IDENTITY_THEFT_DIGITAL',
    'SYSTEM_INTRUSION',
    'RANSOMWARE',
    'SOCIAL_MEDIA_MANIPULATION',
    'ONLINE_SCAM',
    'DIGITAL_EXTORTION',
    'CRYPTO_FRAUD',
    'CHILD_EXPLOITATION_ONLINE',
    'STATE_SYSTEM_ATTACK',
    'OTHER'
);

CREATE TYPE cybre_severity AS ENUM (
    'LOW', 'MEDIUM', 'HIGH', 'CRITICAL'
);

CREATE TABLE cybre_incidents (
    incident_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_cybre_id   VARCHAR(25) UNIQUE NOT NULL,  -- CYBRE-HT-AAAA-NNNNNN
    crime_type          cybre_crime_type NOT NULL,
    severity            cybre_severity NOT NULL DEFAULT 'MEDIUM',
    status              VARCHAR(20) DEFAULT 'OPEN',

    -- Victimes
    victim_count        INTEGER DEFAULT 1,
    victim_snisid_ids   UUID[] DEFAULT '{}',
    victim_types        TEXT[] DEFAULT '{}',         -- INDIVIDUAL, BUSINESS, GOVERNMENT
    total_financial_loss_usd DECIMAL(15,2),

    -- Contexte
    incident_date       TIMESTAMPTZ NOT NULL,
    reported_date       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    attack_vector       TEXT,                        -- Comment l attaque s est produite
    attack_method       TEXT,
    targeted_platform   VARCHAR(100),                -- MonCash, Digicel, SNISIDapp, etc.
    targeted_system     VARCHAR(100),                -- Si attaque système étatique

    -- Suspects
    suspect_ids         UUID[] DEFAULT '{}',
    suspect_phone       TEXT[] DEFAULT '{}',
    suspect_email       TEXT[] DEFAULT '{}',
    suspect_ip_hashes   TEXT[] DEFAULT '{}',         -- Hash IPs pour analyse
    crypto_wallet_ids   UUID[] DEFAULT '{}',         -- Lien CRYPT-HT
    suspect_countries   CHAR(3)[] DEFAULT '{}',      -- Pays d origine estimés

    -- Preuves numériques
    digital_evidence_refs TEXT[] DEFAULT '{}',
    hash_evidence       TEXT[] DEFAULT '{}',
    chain_of_custody_ref VARCHAR(100),

    -- Investigation
    investigating_unit  VARCHAR(50) DEFAULT 'DCPJ_CYBER',
    conatel_ref         VARCHAR(50),                 -- Ref Commission télécom
    case_reference      VARCHAR(100),
    parquet_ref         VARCHAR(100),
    ucref_str_id        UUID,

    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cybre_moncash_fraud_patterns (
    pattern_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id         UUID REFERENCES cybre_incidents(incident_id),
    fraud_type          VARCHAR(50) NOT NULL,        -- SIM_SWAP, PHISHING, AGENT_FRAUD
    moncash_phone       VARCHAR(20),
    amount_stolen_htg   DECIMAL(14,2),
    victims_count       INTEGER DEFAULT 1,
    modus_operandi      TEXT,
    detected_by         VARCHAR(50),                 -- DIGICEL_FRAUD_TEAM, VICTIM_REPORT
    linked_phone_numbers TEXT[] DEFAULT '{}',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cybre_intrusion_attempts (
    attempt_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id         UUID REFERENCES cybre_incidents(incident_id),
    target_system       VARCHAR(100) NOT NULL,
    attack_timestamp    TIMESTAMPTZ NOT NULL,
    attack_type         VARCHAR(50),                 -- BRUTE_FORCE, SQL_INJECTION, PHISHING, API_ABUSE
    source_ip_hash      VARCHAR(64),
    source_country      CHAR(3),
    was_successful      BOOLEAN DEFAULT FALSE,
    data_potentially_accessed TEXT,
    snisid_module_targeted VARCHAR(30),
    detection_source    VARCHAR(50),
    mitigated_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cybre_threat_intelligence (
    threat_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    indicator_type      VARCHAR(30) NOT NULL,        -- IP, DOMAIN, HASH, EMAIL, PHONE
    indicator_value     VARCHAR(500) NOT NULL,
    threat_category     cybre_crime_type,
    confidence_score    SMALLINT CHECK (confidence_score BETWEEN 0 AND 100),
    source              VARCHAR(100),
    is_active           BOOLEAN DEFAULT TRUE,
    first_seen          TIMESTAMPTZ,
    last_seen           TIMESTAMPTZ,
    linked_incidents    UUID[] DEFAULT '{}',
    misp_ref            VARCHAR(100),                -- MISP threat sharing platform ref
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cybre_incidents_type   ON cybre_incidents(crime_type, severity, status);
CREATE INDEX idx_cybre_incidents_date   ON cybre_incidents(incident_date DESC);
CREATE INDEX idx_cybre_moncash_phone    ON cybre_moncash_fraud_patterns(moncash_phone);
CREATE INDEX idx_cybre_intrusions_target ON cybre_intrusion_attempts(target_system, attack_timestamp DESC);
CREATE INDEX idx_cybre_intel_indicator  ON cybre_threat_intelligence(indicator_type, indicator_value)
    WHERE is_active = TRUE;

COMMIT;
```

---

## 3. SERVICE GO CLÉ — DÉTECTION FRAUDE MONCASH

```go
package service

import (
    "context"
    "time"
    "github.com/snisid/cybre-svc/internal/domain"
)

// AnalyzeMoncashPattern detecte les patterns de fraude sur MonCash
func (s *CybreService) AnalyzeMoncashPattern(
    ctx context.Context,
    phoneNumber string,
    period time.Duration,
) (*domain.FraudPatternResult, error) {
    since := time.Now().Add(-period)

    // Verif patterns connus UCREF
    ucrefPatterns, _ := s.ucrefClient.GetMoncashPatterns(ctx, phoneNumber)

    // Verif transactions inhabituelles (heure, montant, frequence)
    stats, _ := s.digicelClient.GetTransactionStats(ctx, phoneNumber, since)

    result := &domain.FraudPatternResult{
        PhoneNumber: phoneNumber,
        AnalyzedFrom: since,
        AnalyzedTo:   time.Now(),
    }

    // Regle: > 10 transferts en < 1h = structuring/smurfing suspect
    if stats.TransactionCount > 10 && stats.PeriodHours < 1 {
        result.AddFlag("RAPID_STRUCTURING", 70)
    }

    // Regle: montant total > HTG 100,000 en 24h pour compte non professionnel
    if stats.TotalAmountHTG > 100000 && !stats.IsBusinessAccount {
        result.AddFlag("HIGH_VOLUME_PERSONAL", 60)
    }

    // Regle: multiple victimes signalent ce numero
    victimCount, _ := s.repo.CountVictimReports(ctx, phoneNumber)
    if victimCount > 3 {
        result.AddFlag("MULTIPLE_VICTIM_REPORTS", 90)
    }

    if len(ucrefPatterns) > 0 {
        result.AddFlag("UCREF_FLAGGED", 85)
    }

    result.RiskScore = result.ComputeRiskScore()
    if result.RiskScore >= 75 {
        _ = s.kafka.Publish(ctx, "cybre.moncash.fraud.alert", result)
        _ = s.ucrefClient.CreateSTR(ctx, domain.STRRequest{
            ReportType:   "MONCASH_PATTERN",
            PhoneNumber:  phoneNumber,
            SuspicionDesc: result.Summary(),
        })
    }
    return result, nil
}
```

---

## 4. API REST

| Méthode | Endpoint                               | Rôle             | Description                    |
|---------|----------------------------------------|------------------|--------------------------------|
| `POST`  | `/api/v1/cybre/incidents`              | DCPJ_CYBER       | Déclarer incident cybercrime   |
| `GET`   | `/api/v1/cybre/incidents/:id`          | DCPJ_CYBER       | Détail incident                |
| `POST`  | `/api/v1/cybre/moncash/analyze`        | DCPJ_CYBER, UCREF| Analyser numéro MonCash        |
| `GET`   | `/api/v1/cybre/intrusions/recent`      | DCPJ_CYBER       | Tentatives intrusion récentes  |
| `POST`  | `/api/v1/cybre/threat-intel`           | DCPJ_CYBER       | Ajouter indicateur menace      |
| `GET`   | `/api/v1/cybre/threat-intel/check`     | DCPJ_CYBER       | Vérifier un indicateur         |
| `GET`   | `/api/v1/cybre/stats/by-type`          | DCPJ_ADMIN       | Stats par type cybercrime      |
| `GET`   | `/api/v1/cybre/snisid/attack-surface`  | SIPCI, IT_ADMIN  | Surface d'attaque SNISID       |

---

## 5. VARIABLES D'ENVIRONNEMENT

```dotenv
CYBRE_DB_HOST=localhost
CYBRE_DB_NAME=snisid_cybre
CYBRE_UCREF_SERVICE_URL=http://ucref-svc:8112
CYBRE_CRYPT_SERVICE_URL=http://crypt-svc:8117
CYBRE_MISP_URL=https://misp.pnh.gov.ht
CYBRE_MISP_API_KEY=<VAULT:cybre/misp_api_key>
CYBRE_DIGICEL_FRAUD_API=https://api.digicel.ht/fraud
CYBRE_MONCASH_FRAUD_THRESHOLD_HTG=100000
CYBRE_RAPID_TX_THRESHOLD=10
CYBRE_VICTIM_REPORTS_THRESHOLD=3
CYBRE_KAFKA_BROKERS=kafka:9092
CYBRE_SERVICE_PORT=8132
```

---
*MP-55 — CYBRE-HT — Cybercriminalité Nationale — SNISID — République d'Haïti*

---
---

# MP-56 — ONG-HT
## Registre National des ONGs et Acteurs Humanitaires en Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-56 | Code : ONG-HT
Dépendances      : SIGEO-HT (MP-48), DPIDE-HT (MP-46), SIGDC-HT (MP-49), SANC-HT (MP-27)
Normes           : OCDE DAC, IOM Humanitarian Principles, Accord Cluster OCHA, Loi 89 Haiti ONG
Acteurs          : MJSP (accréditation), MICT, OCHA, Ministère Plan, ULCC
```

---

## 1. CONTEXTE

Haïti héberge entre 3,000 et 10,000 ONGs enregistrées ou opérant sans enregistrement —
la plus haute densité d'ONGs per capita au monde. Ce contexte crée des risques
sécuritaires : ONGs servant de couverture au blanchiment, personnel étranger non
vérifié, accès à des zones sensibles, détournement de fonds humanitaires.

Ce module ne vise pas à surveiller l'aide humanitaire légitime, mais à :
1. Identifier les organisations non enregistrées opérant illégalement
2. Vérifier le personnel étranger entrant (lien SIFR-HT)
3. Détecter des ONGs servant de façade financière (lien BLAN-HT)
4. Coordonner l'accès aux zones sensibles avec la PNH et le GIPNH

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

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
    national_ong_id     VARCHAR(25) UNIQUE NOT NULL,  -- ONG-HT-AAAA-NNNNNN
    org_name            VARCHAR(200) NOT NULL,
    org_name_local      VARCHAR(200),                 -- Nom en français/créole
    acronym             VARCHAR(20),
    org_type            ong_type NOT NULL,
    registration_status ong_registration_status NOT NULL DEFAULT 'PENDING',
    mjsp_registration_number VARCHAR(50),
    registration_date   DATE,
    registration_expiry DATE,

    -- Siège et présence
    headquarter_country CHAR(3) NOT NULL,
    headquarter_city    VARCHAR(100),
    haiti_office_dept   CHAR(2),
    haiti_office_address TEXT,
    haiti_office_lat    DECIMAL(10,7),
    haiti_office_lng    DECIMAL(10,7),
    operating_depts     CHAR(2)[] DEFAULT '{}',
    operating_communes  TEXT[] DEFAULT '{}',

    -- Secteurs d'intervention
    sectors             TEXT[] DEFAULT '{}',          -- HEALTH, FOOD, SHELTER, WASH, etc.
    annual_budget_usd   DECIMAL(15,2),
    funding_sources     TEXT[] DEFAULT '{}',
    major_donors        TEXT[] DEFAULT '{}',

    -- Personnel
    haiti_staff_count   INTEGER DEFAULT 0,
    expat_staff_count   INTEGER DEFAULT 0,
    director_name       VARCHAR(200),
    director_snisid_id  UUID,
    director_nationality CHAR(3),
    contact_email       VARCHAR(200),
    contact_phone       VARCHAR(30),

    -- Sécurité et conformité
    risk_flag           ong_risk_flag NOT NULL DEFAULT 'NONE',
    risk_notes          TEXT,
    sanc_match_id       UUID,                         -- Lien SANC-HT si match
    blan_case_id        UUID,                         -- Lien BLAN-HT si suspect
    ulcc_ref            VARCHAR(50),                  -- Ref ULCC (anti-corruption)
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
    access_type         VARCHAR(30),         -- CONVOY, DISTRIBUTION, ASSESSMENT
    requested_zones     TEXT[] DEFAULT '{}',
    requested_depts     CHAR(2)[] DEFAULT '{}',
    access_date         DATE NOT NULL,
    access_date_end     DATE,
    purpose             TEXT NOT NULL,
    vehicle_count       SMALLINT DEFAULT 1,
    staff_count         SMALLINT DEFAULT 1,
    status              VARCHAR(20) DEFAULT 'PENDING',  -- APPROVED, DENIED, PENDING
    pnh_escort_required BOOLEAN DEFAULT FALSE,
    approved_by         UUID,
    approval_notes      TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ong_orgs_status   ON ong_organizations(registration_status, risk_flag);
CREATE INDEX idx_ong_orgs_dept     ON ong_organizations USING gin(operating_depts);
CREATE INDEX idx_ong_orgs_country  ON ong_organizations(headquarter_country);
CREATE INDEX idx_ong_orgs_risk     ON ong_organizations(risk_flag) WHERE risk_flag != 'NONE';
CREATE INDEX idx_ong_staff_org     ON ong_staff_registry(org_id, is_active);
CREATE INDEX idx_ong_staff_passport ON ong_staff_registry(passport_number);
CREATE INDEX idx_ong_access_date   ON ong_field_access_requests(access_date, status);

COMMIT;
```

---

## 3. SERVICE GO CLÉ — SCREENING ORGANISATION

```go
package service

import (
    "context"
    "github.com/snisid/ong-svc/internal/domain"
)

// ScreenOrganization effectue le criblage de securite d une ONG
func (s *ONGService) ScreenOrganization(
    ctx context.Context,
    orgID string,
) (*domain.ONGScreeningResult, error) {
    org, err := s.repo.FindByID(ctx, orgID)
    if err != nil {
        return nil, err
    }

    result := &domain.ONGScreeningResult{
        OrgID:     orgID,
        OrgName:   org.Name,
        RiskLevel: "NONE",
    }

    // 1. Verif directeur contre SANC-HT
    if org.DirectorSnisidID != "" {
        sanc, _ := s.sancClient.CheckPerson(ctx, org.DirectorSnisidID)
        if sanc != nil && sanc.IsSanctioned {
            result.AddFlag("DIRECTOR_SANCTIONED", domain.RiskCritical)
        }
    }

    // 2. Verif donateurs contre SANC-HT (si connus)
    for _, donor := range org.MajorDonors {
        sancOrg, _ := s.sancClient.CheckEntityName(ctx, donor)
        if sancOrg != nil && sancOrg.IsMatch {
            result.AddFlag("SANCTIONED_DONOR: " + donor, domain.RiskHigh)
        }
    }

    // 3. Verif UCREF / BLAN-HT (transactions suspectes liees)
    blanCheck, _ := s.blanClient.CheckEntityName(ctx, org.Name)
    if blanCheck != nil && blanCheck.HasActiveCases {
        result.AddFlag("ACTIVE_MONEY_LAUNDERING_CASE", domain.RiskCritical)
    }

    // 4. Verif si ONG opere dans zone gang controllee
    for _, deptCode := range org.OperatingDepts {
        terrCheck, _ := s.terrClient.GetDeptRisk(ctx, deptCode)
        if terrCheck != nil && terrCheck.HasFullControlZones {
            result.AddFlag("OPERATES_IN_GANG_TERRITORY: " + deptCode, domain.RiskMedium)
        }
    }

    result.FinalRiskLevel = result.ComputeRiskLevel()
    _ = s.repo.UpdateRiskFlag(ctx, orgID, result.FinalRiskLevel)

    return result, nil
}
```

---

## 4. API REST

| Méthode | Endpoint                               | Rôle          | Description                       |
|---------|----------------------------------------|---------------|-----------------------------------|
| `POST`  | `/api/v1/ong/organizations`            | MJSP_ADMIN    | Enregistrer organisation          |
| `GET`   | `/api/v1/ong/organizations`            | MJSP, PNH     | Lister ONGs                       |
| `GET`   | `/api/v1/ong/organizations/:id`        | MJSP, PNH     | Détail organisation               |
| `POST`  | `/api/v1/ong/organizations/:id/screen` | MJSP_SECURITY | Criblage sécurité organisation    |
| `POST`  | `/api/v1/ong/staff`                    | MJSP_ADMIN    | Enregistrer membre du personnel   |
| `POST`  | `/api/v1/ong/access-requests`          | ONG_MANAGER   | Demander accès terrain            |
| `PATCH` | `/api/v1/ong/access-requests/:id`      | PNH_ADMIN     | Approuver/refuser accès           |
| `GET`   | `/api/v1/ong/flagged`                  | MJSP, ULCC    | ONGs avec signalements            |
| `GET`   | `/api/v1/ong/unregistered`             | MJSP          | ONGs opérant sans enregistrement  |

---

## 5. VARIABLES D'ENVIRONNEMENT

```dotenv
ONG_DB_HOST=localhost
ONG_DB_NAME=snisid_ong
ONG_SANC_SERVICE_URL=http://sanc-svc:8100
ONG_BLAN_SERVICE_URL=http://blan-svc:8115
ONG_TERR_SERVICE_URL=http://terr-svc:8098
ONG_SLTD_SERVICE_URL=http://sltd-svc:8108
ONG_BLKL_SERVICE_URL=http://blkl-svc:8110
ONG_ULCC_INTEGRATION_URL=https://api.ulcc.gov.ht
ONG_OCHA_API_URL=https://api.reliefweb.int/v1
ONG_RISK_REASSESS_DAYS=90
ONG_SERVICE_PORT=8133
```

---
*MP-56 — ONG-HT — Registre ONGs et Acteurs Humanitaires — SNISID — République d'Haïti*
