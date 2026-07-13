# MP-53 — CORR-HT
## Système Anti-Corruption et Intégrité Institutionnelle de la PNH
### Master Prompt d'Implémentation — SNISID

```
Classification   : TOP SECRET / DIRECTEUR PNH / MJSP UNIQUEMENT
Module           : MP-53 | Code : CORR-HT
Dépendances      : FIR-HT (MP-20), GANG-HT (MP-24), UCREF-INT (MP-39), BLAN-HT (MP-40)
Normes           : Convention ONU Anti-Corruption (CNUCC), UNODC Integrity Standards
Acteurs          : Inspection Générale PNH (IGPNH), ULCC, MJSP, Parquet Anti-Corruption
```

---

## 1. CONTEXTE ET ENJEU CRITIQUE

La corruption interne constitue l'une des menaces les plus graves pour SNISID et pour
la sécurité nationale. Des agents PNH corrompus peuvent :

- **Vendre des informations** à des gangs (alerter avant une opération)
- **Cloner des plaques SE** ou falsifier des documents dans le système
- **Supprimer des mandats** ou des dossiers criminels contre paiement
- **Faciliter la libération** de détenus dangereux via fausses procédures
- **Accéder illégalement** à des données biométriques ou d'identité

### Indicateurs de corruption documentés en Haïti

| Indicateur                                  | Module source     | Poids risque |
|---------------------------------------------|-------------------|--------------|
| Accès à des données sensibles hors service  | Audit SNISID      | Critique     |
| Modification ou suppression de dossiers     | FIR-HT, FPR-HT    | Critique     |
| Transactions financières suspectes          | UCREF-INT/BLAN-HT | Haute        |
| Connexion en zone de gang contrôlée         | SIGEO-HT/LAPI     | Haute        |
| Réseau social avec membres de gangs         | RESO-HT/Neo4j     | Haute        |
| Libérations anormales de détenus            | SIPEP-HT          | Critique     |
| Plaques véhicules SE créées sans autorisation| SIVC-HT          | Critique     |

---

## 2. ARCHITECTURE DE SÉCURITÉ INTERNE

```
┌────────────────────────────────────────────────────────────────┐
│             CORR-HT — COUCHE SURVEILLANCE INTERNE              │
├────────────────────────────────────────────────────────────────┤
│  Sources d'alimentation (accès lecture seule sur autres modules)│
│  Audit logs → SNISID all modules                              │
│  Kafka topics: *.audit, *.access, *.modified, *.deleted       │
├────────────────────────────────────────────────────────────────┤
│  Moteur de détection (Règles + ML)                             │
│  - Règles comportementales (accès hors heures, masse data)    │
│  - Isolation Forest (détection anomalies comportement agent)  │
│  - Graphe social (Neo4j: liens agents ↔ suspects)            │
├────────────────────────────────────────────────────────────────┤
│  Base CORR-HT (PostgreSQL isolé — accès ultra-restreint)      │
│  - Fiches agents sous enquête                                  │
│  - Signalements whistleblowers                                 │
│  - Transactions suspectes liées à des agents                  │
│  - Résultats d'enquêtes IGPNH                                  │
└────────────────────────────────────────────────────────────────┘
```

---

## 3. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE corr_allegation_type AS ENUM (
    'DATA_LEAK_TO_GANG',
    'RECORD_TAMPERING',
    'UNAUTHORIZED_ACCESS',
    'BRIBERY',
    'EXTORTION_OF_CIVILIANS',
    'FACILITATED_PRISON_ESCAPE',
    'STOLEN_CREDENTIALS',
    'FINANCIAL_CORRUPTION',
    'GANG_AFFILIATION',
    'OTHER'
);

CREATE TYPE corr_severity AS ENUM (
    'LOW', 'MEDIUM', 'HIGH', 'CRITICAL'
);

CREATE TYPE corr_status AS ENUM (
    'REPORTED', 'UNDER_INVESTIGATION', 'SUBSTANTIATED',
    'UNSUBSTANTIATED', 'REFERRED_TO_PARQUET', 'CLOSED'
);

-- Fiches des agents sous surveillance ou enquête
CREATE TABLE corr_integrity_cases (
    case_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_corr_id    VARCHAR(25) UNIQUE NOT NULL,  -- CORR-HT-AAAA-NNNNNN
    officer_snisid_id   UUID NOT NULL,
    officer_badge       VARCHAR(30),
    officer_unit        VARCHAR(50),
    officer_rank        VARCHAR(50),
    allegation_type     corr_allegation_type NOT NULL,
    severity            corr_severity NOT NULL,
    status              corr_status NOT NULL DEFAULT 'REPORTED',

    -- Circonstances
    allegation_summary  TEXT NOT NULL,
    incident_date_from  TIMESTAMPTZ,
    incident_date_to    TIMESTAMPTZ,
    evidence_refs       TEXT[] DEFAULT '{}',

    -- Liens criminels (si agent lié à un gang)
    gang_id             UUID,
    gang_member_ids     UUID[] DEFAULT '{}',
    financial_gain_usd  DECIMAL(15,2),
    blan_case_id        UUID,               -- Lien blanchiment si corruption financière

    -- Source du signalement
    reported_by_type    VARCHAR(30),        -- WHISTLEBLOWER, AUDIT, IGPNH, KAFKA_ALERT
    reported_by_id      UUID,
    reporting_date      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_whistleblower    BOOLEAN DEFAULT FALSE,
    whistleblower_protected BOOLEAN DEFAULT FALSE,

    -- Enquête
    igpnh_investigator  UUID,
    investigation_start TIMESTAMPTZ,
    investigation_end   TIMESTAMPTZ,
    investigation_notes TEXT,

    -- Résolution
    sanctions_applied   TEXT,
    referred_to_parquet BOOLEAN DEFAULT FALSE,
    parquet_ref         VARCHAR(100),
    ulcc_ref            VARCHAR(100),

    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Signalements anonymes (whistleblower)
CREATE TABLE corr_whistleblower_reports (
    report_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_token        VARCHAR(64) UNIQUE NOT NULL,  -- Token anonyme pour suivi
    allegation_type     corr_allegation_type NOT NULL,
    severity_estimate   corr_severity,
    officer_unit_hint   VARCHAR(50),
    officer_rank_hint   VARCHAR(50),
    description         TEXT NOT NULL,
    evidence_description TEXT,
    submission_date     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_hash             VARCHAR(64),             -- Hash IP pour anti-abus (pas stocké en clair)
    processed           BOOLEAN DEFAULT FALSE,
    processed_by        UUID,
    integrity_case_id   UUID REFERENCES corr_integrity_cases(case_id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Anomalies comportementales détectées par audit SNISID
CREATE TABLE corr_behavioral_alerts (
    alert_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    officer_snisid_id   UUID NOT NULL,
    alert_type          VARCHAR(50) NOT NULL,
    -- MASS_DOWNLOAD: > 500 records en 1h
    -- OFF_HOURS_ACCESS: accès entre 00h-05h sans autorisation
    -- SENSITIVE_RECORD_ACCESS: accès fiches TOP SECRET hors mission
    -- RECORD_DELETION: suppression d'un dossier criminel
    -- GEOLOCATION_GANG_ZONE: agent localisé dans zone contrôlée gang
    description         TEXT NOT NULL,
    module_source       VARCHAR(30),
    risk_score          SMALLINT CHECK (risk_score BETWEEN 0 AND 100),
    auto_generated      BOOLEAN DEFAULT TRUE,
    reviewed            BOOLEAN DEFAULT FALSE,
    reviewed_by         UUID,
    is_false_positive   BOOLEAN DEFAULT FALSE,
    corr_case_id        UUID REFERENCES corr_integrity_cases(case_id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Déclaration de patrimoine des officiers de rang supérieur
CREATE TABLE corr_asset_declarations (
    declaration_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    officer_snisid_id   UUID NOT NULL,
    declaration_year    SMALLINT NOT NULL,
    real_estate_usd     DECIMAL(15,2) DEFAULT 0,
    vehicles_usd        DECIMAL(15,2) DEFAULT 0,
    bank_accounts_usd   DECIMAL(15,2) DEFAULT 0,
    other_assets_usd    DECIMAL(15,2) DEFAULT 0,
    total_assets_usd    DECIMAL(15,2) GENERATED ALWAYS AS (
                            real_estate_usd + vehicles_usd +
                            bank_accounts_usd + other_assets_usd
                        ) STORED,
    known_salary_annual_usd DECIMAL(12,2),
    unexplained_wealth_usd DECIMAL(15,2),
    is_flagged          BOOLEAN DEFAULT FALSE,
    verified_by         UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Row-Level Security ultra-strict: seulement IGPNH et MJSP
ALTER TABLE corr_integrity_cases      ENABLE ROW LEVEL SECURITY;
ALTER TABLE corr_behavioral_alerts    ENABLE ROW LEVEL SECURITY;
ALTER TABLE corr_asset_declarations   ENABLE ROW LEVEL SECURITY;

CREATE POLICY corr_strict_access ON corr_integrity_cases
    FOR ALL USING (
        current_setting('app.user_role', TRUE) IN ('IGPNH','MJSP_MINISTER','SUPERADMIN')
    );

CREATE POLICY corr_alerts_access ON corr_behavioral_alerts
    FOR ALL USING (
        current_setting('app.user_role', TRUE) IN ('IGPNH','MJSP_MINISTER','SUPERADMIN')
    );

CREATE INDEX idx_corr_cases_officer   ON corr_integrity_cases(officer_snisid_id);
CREATE INDEX idx_corr_cases_status    ON corr_integrity_cases(status, severity);
CREATE INDEX idx_corr_alerts_officer  ON corr_behavioral_alerts(officer_snisid_id, created_at DESC);
CREATE INDEX idx_corr_alerts_unrev    ON corr_behavioral_alerts(reviewed) WHERE reviewed = FALSE;

COMMIT;
```

---

## 4. SERVICE GO — DÉTECTION AUTOMATIQUE D'ANOMALIES

```go
package service

import (
    "context"
    "time"
    "github.com/snisid/corr-svc/internal/domain"
)

type BehavioralAnalyzer struct {
    repo    domain.CORRRepository
    kafka   domain.EventPublisher
    igpnh   domain.IGPNHNotifier
}

// AnalyzeAuditEvent analyse chaque evenement d audit en temps reel
func (s *BehavioralAnalyzer) AnalyzeAuditEvent(
    ctx context.Context,
    event domain.AuditEvent,
) error {
    alerts := []domain.BehavioralAlert{}

    // Regle 1: Telechargement massif de donnees
    if event.RecordCount > 500 && event.Duration < time.Hour {
        alerts = append(alerts, domain.BehavioralAlert{
            OfficerID:   event.OfficerID,
            AlertType:   "MASS_DOWNLOAD",
            Description: "Telechargement de plus de 500 dossiers en moins d une heure",
            RiskScore:   75,
            ModuleSource: event.Module,
        })
    }

    // Regle 2: Acces hors heures de service (00h-05h)
    if event.Timestamp.Hour() >= 0 && event.Timestamp.Hour() < 5 {
        alerts = append(alerts, domain.BehavioralAlert{
            OfficerID:   event.OfficerID,
            AlertType:   "OFF_HOURS_ACCESS",
            Description: "Acces au systeme entre 00h et 05h sans autorisation speciale",
            RiskScore:   60,
            ModuleSource: event.Module,
        })
    }

    // Regle 3: Suppression de dossier criminel
    if event.Action == "DELETE" && event.Module == "FIR-HT" {
        alerts = append(alerts, domain.BehavioralAlert{
            OfficerID:   event.OfficerID,
            AlertType:   "RECORD_DELETION",
            Description: "Suppression d un dossier FIR-HT — verification requise",
            RiskScore:   95,
            ModuleSource: "FIR-HT",
        })
    }

    // Regle 4: Acces dossier TOP SECRET hors mission assignee
    if event.Classification == "TOP_SECRET" && !event.IsAuthorized {
        alerts = append(alerts, domain.BehavioralAlert{
            OfficerID:   event.OfficerID,
            AlertType:   "SENSITIVE_RECORD_ACCESS",
            Description: "Acces non autorise a un dossier TOP SECRET",
            RiskScore:   90,
            ModuleSource: event.Module,
        })
    }

    for _, alert := range alerts {
        _ = s.repo.SaveBehavioralAlert(ctx, alert)

        // Si score > 80 -> notifier IGPNH immediatement
        if alert.RiskScore >= 80 {
            _ = s.igpnh.NotifyUrgent(ctx, alert)
        }

        _ = s.kafka.Publish(ctx, "corr.behavioral.alert", alert)
    }
    return nil
}
```

---

## 5. API REST

| Méthode | Endpoint                               | Rôle            | Description                       |
|---------|----------------------------------------|-----------------|-----------------------------------|
| `POST`  | `/api/v1/corr/cases`                   | IGPNH           | Ouvrir enquête intégrité          |
| `GET`   | `/api/v1/corr/cases/:id`               | IGPNH, MJSP     | Détail enquête                    |
| `GET`   | `/api/v1/corr/cases/active`            | IGPNH           | Enquêtes actives                  |
| `POST`  | `/api/v1/corr/whistleblower`           | ANONYMOUS       | Signalement anonyme               |
| `GET`   | `/api/v1/corr/whistleblower/:token`    | WHISTLEBLOWER   | Suivi signalement par token       |
| `GET`   | `/api/v1/corr/alerts/behavioral`       | IGPNH           | Alertes comportementales          |
| `POST`  | `/api/v1/corr/declarations`            | OFFICER         | Soumettre déclaration patrimoine  |
| `GET`   | `/api/v1/corr/declarations/flagged`    | IGPNH, ULCC     | Déclarations anormales            |
| `GET`   | `/api/v1/corr/risk-scores`             | IGPNH           | Scores de risque par agent        |

---

## 6. SÉCURITÉ SPÉCIALE

```
ISOLEMENT TOTAL:
- Base de données CORR-HT sur serveur physique séparé
- Réseau VLAN isolé — aucune connexion depuis internet
- Accès uniquement depuis terminaux IGPNH dédiés (hardware-bound)
- Chiffrement colonne-par-colonne (AES-256 + HSM Luna)
- Log de tout accès à CORR-HT — immutable (Kafka + write-once S3)
- Sauvegardes chiffrées hors-site (coffre BCN INTERPOL)
- Audit externe trimestriel (OIG ONU ou organisme indépendant)
```

---

## 7. VARIABLES D'ENVIRONNEMENT

```dotenv
CORR_DB_HOST=corr-isolated-db.pnh.internal
CORR_DB_PORT=5432
CORR_DB_NAME=snisid_corr_secret
CORR_DB_USER=corr_igpnh_svc
CORR_DB_PASSWORD=<HSM:corr/db_password>
CORR_NEO4J_URI=bolt://corr-neo4j.pnh.internal:7687
CORR_KAFKA_BROKERS=kafka:9092
CORR_KAFKA_TOPIC_AUDIT=*.audit
CORR_IGPNH_NOTIFY_URL=https://igpnh.pnh.gov.ht/api/urgent
CORR_BEHAVIORAL_DOWNLOAD_THRESHOLD=500
CORR_BEHAVIORAL_OFHOURS_START=0
CORR_BEHAVIORAL_OFHOURS_END=5
CORR_HIGH_RISK_ALERT_THRESHOLD=80
CORR_WHISTLEBLOWER_TOKEN_LENGTH=64
CORR_SERVICE_PORT=8130
```

---
*MP-53 — CORR-HT — Anti-Corruption et Intégrité — SNISID — République d'Haïti*
*ACCÈS RESTREINT : IGPNH / MJSP / DIRECTEUR PNH UNIQUEMENT*
