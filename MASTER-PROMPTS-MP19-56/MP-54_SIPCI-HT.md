# MP-54 — SIPCI-HT
## Système de Protection des Infrastructures Critiques d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-54 | Code : SIPCI-HT
Dépendances      : SIGEO-HT (MP-48), GANG-HT (MP-24), SIVC-HT (MP-18), SIGDC-HT (MP-49)
Normes           : NIPP (National Infrastructure Protection Plan USA), ISO 22301, IEC 62443
Acteurs          : MSP, ED (Électricité d'Haïti), MTPTC, APN, Télécoms, MJSP
```

---

## 1. CONTEXTE

Haïti a subi de nombreuses attaques contre ses infrastructures critiques :
- **Centrales électriques** : Sabotages répétés de l'EDH dans des zones disputées
- **Axes routiers** (RN1, RN2, RN3) : Occupés par des barrages de gangs
- **Ports** (Port-au-Prince, Cap-Haïtien) : Sous pression extorsion
- **Télécommunications** : Tours Digicel/Natcom ciblées dans des zones de conflit
- **Eau potable** : Adductions dans des zones contrôlées par gangs
- **Hôpitaux** : Attaques documentées contre personnels soignants

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE sipci_asset_category AS ENUM (
    'ENERGY',           -- Centrales, lignes HT, stations service
    'TRANSPORT',        -- Routes nationales, ponts, aéroports, ports
    'WATER',            -- Adductions, stations pompage, barrages
    'TELECOMS',         -- Tours, fibres, datacenters
    'HEALTH',           -- Hôpitaux, centres de santé critiques
    'FINANCE',          -- BRH, banques systémiques, MonCash/Digicel
    'GOVERNMENT',       -- Palais National, tribunaux, ministères
    'EDUCATION',        -- Grandes universités, centres formation PNH
    'FOOD_SUPPLY'       -- Marchés stratégiques, silos, entrepôts
);

CREATE TYPE sipci_threat_level AS ENUM (
    'NORMAL', 'ELEVATED', 'HIGH', 'SEVERE', 'CRITICAL'
);

CREATE TABLE sipci_assets (
    asset_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_sipci_id   VARCHAR(25) UNIQUE NOT NULL,  -- SIPCI-HT-NNNNNN
    asset_name          VARCHAR(200) NOT NULL,
    asset_category      sipci_asset_category NOT NULL,
    owner_entity        VARCHAR(200),
    operating_org       VARCHAR(200),
    location_desc       VARCHAR(300),
    dept_code           CHAR(2) NOT NULL,
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7) NOT NULL,
    lng                 DECIMAL(10,7) NOT NULL,

    -- Évaluation de criticité
    criticality_score   SMALLINT CHECK (criticality_score BETWEEN 1 AND 10),
    population_served   INTEGER,
    dependency_assets   UUID[] DEFAULT '{}',    -- Autres SIPCI assets qui en dépendent
    single_point_failure BOOLEAN DEFAULT FALSE, -- Défaillance = coupure nationale

    -- Statut sécuritaire
    current_threat_level sipci_threat_level NOT NULL DEFAULT 'NORMAL',
    is_in_gang_zone     BOOLEAN DEFAULT FALSE,
    controlling_gang_id UUID,
    under_extortion     BOOLEAN DEFAULT FALSE,
    extors_case_id      UUID,

    -- Incidents
    incident_count_12m  INTEGER DEFAULT 0,
    last_incident_date  TIMESTAMPTZ,

    -- Protection assignée
    protection_unit     VARCHAR(50),
    security_guards     INTEGER DEFAULT 0,
    has_cctv            BOOLEAN DEFAULT FALSE,
    cctv_count          SMALLINT DEFAULT 0,
    has_perimeter       BOOLEAN DEFAULT FALSE,
    has_backup_power    BOOLEAN DEFAULT FALSE,

    -- Contacts
    site_manager_name   VARCHAR(200),
    site_manager_phone  VARCHAR(30),
    emergency_contact   VARCHAR(200),

    notes               TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sipci_incidents (
    incident_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id            UUID NOT NULL REFERENCES sipci_assets(asset_id),
    incident_type       VARCHAR(50) NOT NULL,
    -- ATTACK, SABOTAGE, EXTORTION, OCCUPATION, THEFT, THREAT, DISRUPTION
    incident_date       TIMESTAMPTZ NOT NULL,
    perpetrator_type    VARCHAR(30),    -- GANG, UNKNOWN, CRIMINAL, POLITICAL
    gang_id             UUID,
    description         TEXT NOT NULL,
    impact_severity     SMALLINT CHECK (impact_severity BETWEEN 1 AND 10),
    population_affected INTEGER,
    service_disruption_hours DECIMAL(8,2),
    economic_loss_usd   DECIMAL(15,2),
    responding_units    TEXT[] DEFAULT '{}',
    sivc_alert_ids      UUID[] DEFAULT '{}',
    case_reference      VARCHAR(100),
    resolution_date     TIMESTAMPTZ,
    resolution_notes    TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sipci_protection_plans (
    plan_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id            UUID NOT NULL REFERENCES sipci_assets(asset_id),
    plan_name           VARCHAR(200),
    threat_scenarios    TEXT[] DEFAULT '{}',
    mitigation_measures TEXT NOT NULL,
    response_procedures TEXT,
    resources_required  TEXT,
    responsible_unit    VARCHAR(100),
    review_date         DATE,
    is_active           BOOLEAN DEFAULT TRUE,
    approved_by         UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sipci_assets_category ON sipci_assets(asset_category, current_threat_level);
CREATE INDEX idx_sipci_assets_dept     ON sipci_assets(dept_code);
CREATE INDEX idx_sipci_assets_gang     ON sipci_assets(is_in_gang_zone) WHERE is_in_gang_zone = TRUE;
CREATE INDEX idx_sipci_assets_critical ON sipci_assets(criticality_score DESC) WHERE single_point_failure = TRUE;
CREATE INDEX idx_sipci_incidents_date  ON sipci_incidents(incident_date DESC, asset_id);
CREATE INDEX idx_sipci_assets_geom     ON sipci_assets(lat, lng);

COMMIT;
```

---

## 3. SERVICE GO CLÉ — CALCUL RISQUE TEMPS RÉEL

```go
package service

import (
    "context"
    "github.com/snisid/sipci-svc/internal/domain"
)

// AssessAssetRisk evalue le risque actuel d une infrastructure critique
func (s *InfraProtectionService) AssessAssetRisk(
    ctx context.Context,
    assetID string,
) (*domain.RiskAssessment, error) {
    asset, err := s.repo.FindByID(ctx, assetID)
    if err != nil {
        return nil, err
    }

    assessment := &domain.RiskAssessment{
        AssetID:   assetID,
        BaseScore: float64(asset.CriticalityScore) * 10,
    }

    // 1. Verifier si dans zone gang (TERR-HT)
    terrCheck, _ := s.terrClient.CheckPointSafety(ctx, asset.Lat, asset.Lng)
    if terrCheck != nil && !terrCheck.IsSafe {
        assessment.AddFactor("IN_GANG_TERRITORY", 40.0)
        if terrCheck.HighRisk {
            assessment.AddFactor("FULL_GANG_CONTROL", 30.0)
        }
    }

    // 2. Incidents recents (12 derniers mois)
    recentIncidents, _ := s.repo.CountRecentIncidents(ctx, assetID, 12)
    if recentIncidents > 0 {
        assessment.AddFactor("RECENT_INCIDENTS", float64(recentIncidents) * 5.0)
    }

    // 3. Extorsion active
    if asset.UnderExtortion {
        assessment.AddFactor("ACTIVE_EXTORTION", 25.0)
    }

    // 4. Alertes vehiculaires recentes dans la zone
    vehicleAlerts, _ := s.sivcClient.GetAlertsNearLocation(ctx, asset.Lat, asset.Lng, 500)
    if len(vehicleAlerts) > 0 {
        assessment.AddFactor("NEARBY_VEHICLE_ALERTS", float64(len(vehicleAlerts)) * 3.0)
    }

    // Calculer niveau de menace final
    assessment.FinalScore = assessment.ComputeFinalScore()
    assessment.ThreatLevel = domain.ClassifyThreatLevel(assessment.FinalScore)

    // Mettre a jour le niveau de menace de l actif
    _ = s.repo.UpdateThreatLevel(ctx, assetID, assessment.ThreatLevel)

    return assessment, nil
}
```

---

## 4. API REST

| Méthode | Endpoint                              | Rôle             | Description                     |
|---------|---------------------------------------|------------------|---------------------------------|
| `GET`   | `/api/v1/sipci/assets`                | MSP, MTPTC       | Inventaire infrastructures      |
| `GET`   | `/api/v1/sipci/assets/:id`            | MSP, MTPTC       | Détail infrastructure           |
| `GET`   | `/api/v1/sipci/assets/critical`       | MSP              | Actifs criticité ≥ 8            |
| `GET`   | `/api/v1/sipci/assets/under-threat`   | MSP, DCPJ        | Actifs sous menace élevée       |
| `POST`  | `/api/v1/sipci/assets`                | MSP_ADMIN        | Enregistrer infrastructure      |
| `POST`  | `/api/v1/sipci/incidents`             | SIPCI_OFFICER    | Déclarer incident               |
| `GET`   | `/api/v1/sipci/incidents/recent`      | MSP, PNH         | Incidents des 30 derniers jours |
| `GET`   | `/api/v1/sipci/risk-map`              | MSP              | Carte de risque (GeoJSON)       |
| `POST`  | `/api/v1/sipci/assets/:id/assess`     | MSP_ANALYST      | Évaluer risque en temps réel    |

---

## 5. DASHBOARD EXÉCUTIF — VUE EN TEMPS RÉEL

```
Tableau de bord MSP (React + MapLibre):
┌─────────────────────────────────────────────────────────┐
│  STATUT INFRASTRUCTURES CRITIQUES HAÏTI — Temps Réel    │
├─────────────────┬───────────────┬───────────────────────┤
│  ÉNERGIE  [6🔴] │ TRANSPORT[4🟡]│ EAU       [2🟢]       │
│  SANTÉ    [3🔴] │ TÉLÉCOM  [1🔴]│ FINANCE   [5🟡]       │
└─────────────────┴───────────────┴───────────────────────┘
[CARTE MAPLIBRE — Points colorés par niveau menace]
[LISTE ALERTES RÉCENTES — Défilant temps réel via SSE]
```

---

## 6. VARIABLES D'ENVIRONNEMENT

```dotenv
SIPCI_DB_HOST=localhost
SIPCI_DB_NAME=snisid_sipci
SIPCI_TERR_SERVICE_URL=http://terr-svc:8098
SIPCI_SIVC_SERVICE_URL=http://sivc-svc:8090
SIPCI_SIGEO_SERVICE_URL=http://sigeo-svc:8125
SIPCI_GANG_SERVICE_URL=http://gang-svc:8095
SIPCI_THREAT_REASSESS_INTERVAL_MIN=60
SIPCI_CRITICAL_THRESHOLD=8
SIPCI_KAFKA_BROKERS=kafka:9092
SIPCI_SERVICE_PORT=8131
```

---
*MP-54 — SIPCI-HT — Protection Infrastructures Critiques — SNISID — République d'Haïti*
