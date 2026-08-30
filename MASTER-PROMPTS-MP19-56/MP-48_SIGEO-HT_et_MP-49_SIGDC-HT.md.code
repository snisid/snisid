# MP-48 — SIGEO-HT
## Système National de Géo-Intelligence Criminelle et Cartographie Sécuritaire
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-48 | Code : SIGEO-HT
Dépendances      : TERR-HT (MP-26), SIVC-HT (MP-18), GANG-HT (MP-24), SIGDC-HT (MP-49)
Normes           : ISO 19115 GIS, OGC WMS/WFS/WMTS, PostGIS 3.x, Open Street Map
Acteurs          : DCPJ, BRI, PNH toutes unités, MSP, ONGs sécurité, Humanitaires
```

---

## 1. CONTEXTE

SIGEO-HT est le SIG (Système d'Information Géographique) sécuritaire central de SNISID.
Il fusionne en temps réel les données spatiales de tous les modules (incidents criminels,
territoires gangs, sightings LAPI, alertes véhiculaires, IDP camps) pour produire une
image opérationnelle commune (IOC) pour les décideurs et les agents terrain.

---

## 2. ARCHITECTURE TECHNIQUE

```
┌──────────────────────────────────────────────────────────────┐
│                    SIGEO-HT STACK GIS                         │
├──────────────────────────────────────────────────────────────┤
│  Frontend                                                     │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  MapLibre GL JS + PMTiles                            │    │
│  │  Couches: Fond OSM + Administratif Haïti + Crime    │    │
│  └─────────────────────────────────────────────────────┘    │
├──────────────────────────────────────────────────────────────┤
│  Tile Server                   │  GeoServer / Martin         │
│  (WMTS Vectoriel)              │  (WMS Raster heatmaps)     │
├──────────────────────────────────────────────────────────────┤
│  PostgreSQL + PostGIS 3.4      │  ClickHouse Spatial         │
│  (Données vecteur temps réel) │  (Analytiques historiques)  │
├──────────────────────────────────────────────────────────────┤
│  Kafka Consumers (Alimentation depuis tous les modules)       │
│  GANG.incident → SIVC.sighting → EXTORS.toll → DIPE.case    │
└──────────────────────────────────────────────────────────────┘
```

---

## 3. BASE DE DONNÉES

```sql
BEGIN;

CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;
CREATE EXTENSION IF NOT EXISTS h3;                 -- Indexation hexagonale H3 Uber

-- Référentiel administratif haïtien (couche de base)
CREATE TABLE sigeo_admin_boundaries (
    boundary_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    level               SMALLINT NOT NULL,         -- 1=Pays, 2=Dept, 3=Arrondissement, 4=Commune, 5=Section
    code                VARCHAR(20) UNIQUE NOT NULL,
    name                VARCHAR(150) NOT NULL,
    parent_code         VARCHAR(20),
    geom                GEOMETRY(MultiPolygon, 4326) NOT NULL,
    population_est      INTEGER,
    area_km2            DECIMAL(10,3),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Couche centralisée des incidents géolocalisés (vue matérialisée multi-sources)
CREATE TABLE sigeo_incidents_unified (
    event_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_module       VARCHAR(20) NOT NULL,      -- GANG, SIVC, DIPE, VICT, EXPL, etc.
    source_record_id    UUID NOT NULL,
    event_type          VARCHAR(50) NOT NULL,
    event_date          TIMESTAMPTZ NOT NULL,
    lat                 DECIMAL(10,7) NOT NULL,
    lng                 DECIMAL(10,7) NOT NULL,
    geom                GEOMETRY(Point, 4326) GENERATED ALWAYS AS (
                            ST_SetSRID(ST_MakePoint(lng, lat), 4326)
                        ) STORED,
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    h3_index_8          VARCHAR(20),               -- H3 resolution 8 (~460m cellules)
    h3_index_10         VARCHAR(20),               -- H3 resolution 10 (~65m cellules)
    severity            SMALLINT CHECK (severity BETWEEN 1 AND 10),
    gang_id             UUID,
    description         TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Trigger: calculer H3 index automatiquement
CREATE OR REPLACE FUNCTION sigeo_compute_h3()
RETURNS TRIGGER AS $$
BEGIN
    NEW.h3_index_8  := h3_lat_lng_to_cell(NEW.lat, NEW.lng, 8)::text;
    NEW.h3_index_10 := h3_lat_lng_to_cell(NEW.lat, NEW.lng, 10)::text;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_sigeo_h3
    BEFORE INSERT ON sigeo_incidents_unified
    FOR EACH ROW EXECUTE FUNCTION sigeo_compute_h3();

-- Couche checkpoints et barrages
CREATE TABLE sigeo_checkpoints (
    cp_id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cp_type             VARCHAR(30) NOT NULL,      -- POLICE, GANG_TOLL, MILITARY, HUMANITARIAN
    location            GEOMETRY(Point, 4326) NOT NULL,
    dept_code           CHAR(2),
    road_number         VARCHAR(10),
    description         VARCHAR(300),
    controlling_gang_id UUID,
    is_active           BOOLEAN DEFAULT TRUE,
    source_module       VARCHAR(20),
    source_record_id    UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Couche camps humanitaires
CREATE TABLE sigeo_humanitarian_sites (
    site_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_type           VARCHAR(30) NOT NULL,      -- IDP_CAMP, HOSPITAL, SHELTER, DISTRIBUTION
    site_name           VARCHAR(150),
    location            GEOMETRY(Point, 4326) NOT NULL,
    dept_code           CHAR(2),
    managing_org        VARCHAR(150),
    capacity            INTEGER,
    current_population  INTEGER,
    is_active           BOOLEAN DEFAULT TRUE,
    dpide_camp_id       UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index spatiaux
CREATE INDEX idx_sigeo_incidents_geom ON sigeo_incidents_unified USING GIST(geom);
CREATE INDEX idx_sigeo_incidents_h3_8 ON sigeo_incidents_unified(h3_index_8);
CREATE INDEX idx_sigeo_incidents_date ON sigeo_incidents_unified(event_date DESC);
CREATE INDEX idx_sigeo_incidents_dept ON sigeo_incidents_unified(dept_code, event_date DESC);
CREATE INDEX idx_sigeo_incidents_gang ON sigeo_incidents_unified(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_sigeo_cp_geom        ON sigeo_checkpoints USING GIST(location) WHERE is_active = TRUE;

COMMIT;
```

---

## 4. SERVICE GO CLÉ — ANALYSE DE ZONE

```go
package service

import (
    "context"
    "time"
    "github.com/snisid/sigeo-svc/internal/domain"
)

// GetZoneReport genere un rapport de securite pour une zone geographique
func (s *GeoIntelService) GetZoneReport(
    ctx context.Context,
    deptCode string,
    commune string,
    period time.Duration,
) (*domain.ZoneSecurityReport, error) {
    since := time.Now().Add(-period)

    // 1. Incidents par type dans la zone
    incidents, _ := s.repo.CountIncidentsByZone(ctx, deptCode, commune, since)

    // 2. Gangs actifs dans la zone
    gangs, _ := s.terrClient.GetActiveGangsByDept(ctx, deptCode)

    // 3. Alertes SIVC actives dans la zone
    vehicles, _ := s.sivcClient.GetActiveAlertsByDept(ctx, deptCode)

    // 4. Personnes recherchees signalees dans la zone
    wanted, _ := s.fprClient.GetRecentSightingsByDept(ctx, deptCode, since)

    // 5. Score de risque composite
    riskScore := s.computeRiskScore(incidents, len(gangs), len(vehicles))

    return &domain.ZoneSecurityReport{
        DeptCode:       deptCode,
        Commune:        commune,
        Period:         period,
        GeneratedAt:    time.Now(),
        IncidentCount:  incidents,
        ActiveGangs:    gangs,
        ActiveVehicleAlerts: vehicles,
        RecentWantedSightings: wanted,
        OverallRiskScore: riskScore,
        RiskLevel:      domain.ClassifyRisk(riskScore),
    }, nil
}
```

---

## 5. API REST

| Méthode | Endpoint                             | Rôle         | Description                       |
|---------|--------------------------------------|--------------|-----------------------------------|
| `GET`   | `/api/v1/sigeo/incidents/unified`    | PNH          | Incidents unifiés (GeoJSON)       |
| `GET`   | `/api/v1/sigeo/heatmap`              | DCPJ         | Carte de chaleur incidents        |
| `GET`   | `/api/v1/sigeo/zone-report`          | PNH, BRI     | Rapport sécurité d'une zone       |
| `GET`   | `/api/v1/sigeo/checkpoints/active`   | PNH          | Checkpoints actifs (GeoJSON)      |
| `GET`   | `/api/v1/sigeo/tiles/{z}/{x}/{y}.mvt`| FRONTEND     | Tuiles vectorielles MapLibre      |
| `GET`   | `/api/v1/sigeo/h3/clusters`          | DCPJ_INTEL   | Clusters H3 d'incidents           |
| `POST`  | `/api/v1/sigeo/incidents/ingest`     | KAFKA_WORKER | Ingestion événement module source |

---

## 6. VARIABLES D'ENVIRONNEMENT

```dotenv
SIGEO_DB_HOST=localhost
SIGEO_DB_NAME=snisid_sigeo
SIGEO_CLICKHOUSE_ADDR=clickhouse:9000
SIGEO_KAFKA_BROKERS=kafka:9092
SIGEO_MARTIN_TILE_URL=http://martin:3000
SIGEO_OSM_TILE_URL=https://tile.openstreetmap.org
SIGEO_H3_RESOLUTION_HEATMAP=8
SIGEO_SERVICE_PORT=8125
```

---
*MP-48 — SIGEO-HT — Géo-Intelligence Criminelle — SNISID — République d'Haïti*

---
---

# MP-49 — SIGDC-HT
## Système National de Gestion des Désastres Civils et Crises d'Urgence
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-49 | Code : SIGDC-HT
Dépendances      : SIGEO-HT (MP-48), DPIDE-HT (MP-46), RVIN-HT (MP-50), SNISID-BIO-ADN
Normes           : Cadre de Sendai 2015-2030, Sphère Standards, OSOCC, OIM DTM
Acteurs          : CSPAN, SNGRD, OCHA, Croix-Rouge Haïtienne, MSF, OIM, UNDP
```

---

## 1. CONTEXTE

Haïti est classé parmi les pays les plus vulnérables aux catastrophes naturelles :
- Séismes majeurs : 2010 (220,000 morts), 2021 (2,248 morts, Grand-Sud)
- Ouragans : Matthew 2016, Irma 2017, Laura 2020
- Inondations : Artibonite, Grand-Anse — cycliques
- Risque tsunami : Faille Enriquillo-Plaintain Garden

SIGDC-HT intègre les capacités d'alerte précoce, gestion des victimes et
coordination des secours dans SNISID — permettant l'identification biométrique
des victimes et des morts lors des catastrophes.

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE sigdc_disaster_type AS ENUM (
    'EARTHQUAKE', 'HURRICANE', 'TSUNAMI', 'FLOOD',
    'LANDSLIDE', 'FIRE_MASS', 'INDUSTRIAL_ACCIDENT',
    'EPIDEMIC', 'SECURITY_MASS_CASUALTY'
);

CREATE TYPE sigdc_alert_level AS ENUM (
    'WATCH', 'WARNING', 'EMERGENCY', 'CATASTROPHE'
);

CREATE TABLE sigdc_disasters (
    disaster_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_sigdc_id   VARCHAR(25) UNIQUE NOT NULL,  -- SIGDC-HT-AAAA-NNNNNN
    disaster_type       sigdc_disaster_type NOT NULL,
    disaster_name       VARCHAR(200),
    alert_level         sigdc_alert_level NOT NULL,
    status              VARCHAR(20) DEFAULT 'ACTIVE',
    onset_date          TIMESTAMPTZ NOT NULL,
    affected_depts      CHAR(2)[] DEFAULT '{}',
    affected_communes   TEXT[] DEFAULT '{}',
    epicenter_lat       DECIMAL(10,7),
    epicenter_lng       DECIMAL(10,7),
    magnitude           DECIMAL(4,2),                 -- Pour seismes
    wind_speed_kmh      INTEGER,                      -- Pour ouragans
    estimated_affected  INTEGER,
    confirmed_dead      INTEGER DEFAULT 0,
    confirmed_injured   INTEGER DEFAULT 0,
    confirmed_missing   INTEGER DEFAULT 0,
    confirmed_displaced INTEGER DEFAULT 0,
    response_agencies   TEXT[] DEFAULT '{}',
    coordination_center VARCHAR(200),
    ocha_flash_ref      VARCHAR(100),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sigdc_victim_registrations (
    registration_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    disaster_id         UUID NOT NULL REFERENCES sigdc_disasters(disaster_id),
    snisid_person_id    UUID,
    full_name           VARCHAR(200),
    dob                 DATE,
    gender              VARCHAR(10),
    status              VARCHAR(30) NOT NULL,   -- ALIVE, INJURED, DECEASED, MISSING
    injury_description  TEXT,
    location_found      VARCHAR(300),
    dept_code           CHAR(2),
    hospital_sent_to    VARCHAR(150),
    morgue_location     VARCHAR(150),
    afis_subject_id     UUID,
    dna_sample_taken    BOOLEAN DEFAULT FALSE,
    dna_sample_ref      VARCHAR(100),
    rvin_case_id        UUID,                   -- Si corps non identifie
    dpide_idp_id        UUID,                   -- Si personne deplacee
    registration_date   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    registered_by       UUID NOT NULL,
    org_registering     VARCHAR(100)
);

CREATE TABLE sigdc_resources (
    resource_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    disaster_id         UUID NOT NULL REFERENCES sigdc_disasters(disaster_id),
    resource_type       VARCHAR(50) NOT NULL,   -- SEARCH_RESCUE, MEDICAL, FOOD, SHELTER, WATER
    provider_org        VARCHAR(150),
    quantity            INTEGER,
    unit                VARCHAR(30),
    location_lat        DECIMAL(10,7),
    location_lng        DECIMAL(10,7),
    dept_code           CHAR(2),
    available_from      TIMESTAMPTZ,
    available_until     TIMESTAMPTZ,
    status              VARCHAR(20) DEFAULT 'AVAILABLE',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sigdc_early_warnings (
    warning_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    disaster_type       sigdc_disaster_type NOT NULL,
    alert_level         sigdc_alert_level NOT NULL,
    source_agency       VARCHAR(100),           -- UHM (seismes), NHC (ouragans)
    message_text        TEXT NOT NULL,
    affected_depts      CHAR(2)[] DEFAULT '{}',
    issued_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at          TIMESTAMPTZ,
    channels_sent       TEXT[] DEFAULT '{}',    -- SMS, RADIO, APP, SIRENE
    population_reached  INTEGER
);

CREATE INDEX idx_sigdc_disasters_status ON sigdc_disasters(status, disaster_type);
CREATE INDEX idx_sigdc_disasters_dept   ON sigdc_disasters USING gin(affected_depts);
CREATE INDEX idx_sigdc_victims_disaster ON sigdc_victim_registrations(disaster_id, status);
CREATE INDEX idx_sigdc_victims_snisid   ON sigdc_victim_registrations(snisid_person_id) WHERE snisid_person_id IS NOT NULL;
CREATE INDEX idx_sigdc_resources        ON sigdc_resources(disaster_id, resource_type, status);
CREATE INDEX idx_sigdc_warnings_date    ON sigdc_early_warnings(issued_at DESC, alert_level);

COMMIT;
```

---

## 3. SERVICE GO CLÉ — ALERTE PRÉCOCE

```go
package service

import (
    "context"
    "fmt"
    "time"
    "github.com/snisid/sigdc-svc/internal/domain"
)

func (s *DisasterService) IssueEarlyWarning(
    ctx context.Context,
    warning domain.EarlyWarningRequest,
) error {
    w := domain.EarlyWarning{
        DisasterType:  warning.DisasterType,
        AlertLevel:    warning.AlertLevel,
        SourceAgency:  warning.SourceAgency,
        MessageText:   warning.MessageText,
        AffectedDepts: warning.AffectedDepts,
        IssuedAt:      time.Now(),
        ExpiresAt:     time.Now().Add(warning.ValidDuration),
    }

    if err := s.repo.SaveWarning(ctx, &w); err != nil {
        return fmt.Errorf("sauvegarde alerte: %w", err)
    }

    // Diffusion multi-canal (parallele)
    go s.smsGateway.BroadcastToDepts(ctx, warning.AffectedDepts, warning.MessageText)
    go s.radioClient.BroadcastAlert(ctx, w)
    go s.kafka.Publish(ctx, "sigdc.early.warning.issued", w)

    // Si CATASTROPHE -> activer plan d urgence complet
    if warning.AlertLevel == domain.AlertLevelCatastrophe {
        go s.activateEmergencyPlan(ctx, warning.AffectedDepts)
    }
    return nil
}

// IdentifyDisasterVictim identifie une victime par biometrie en conditions de terrain
func (s *DisasterService) IdentifyDisasterVictim(
    ctx context.Context,
    req domain.VictimIdentificationRequest,
) (*domain.VictimIdentificationResult, error) {
    result := &domain.VictimIdentificationResult{}

    // 1. Empreintes digitales (si disponibles)
    if req.FingerprintData != nil {
        afisResult, _ := s.afisClient.SearchTenprint(ctx, *req.FingerprintData)
        if afisResult != nil && len(afisResult) > 0 && afisResult[0].Score >= 0.90 {
            result.Identified = true
            result.SNISIDPersonID = afisResult[0].SubjectID
            result.ConfidenceScore = afisResult[0].Score * 100
            result.IdentificationMethod = "AFIS"
        }
    }

    // 2. Reconnaissance faciale si non identifie par AFIS
    if !result.Identified && req.FacePhoto != nil {
        faceResult, _ := s.faceClient.SearchFace(ctx, *req.FacePhoto)
        if faceResult != nil && faceResult.Score >= 0.92 {
            result.Identified = true
            result.SNISIDPersonID = faceResult.PersonID
            result.IdentificationMethod = "FACE_RECOGNITION"
        }
    }

    // 3. ADN (si echantillon preleve)
    if !result.Identified && req.DNASampleRef != "" {
        dnaResult, _ := s.bioADNClient.SearchProfile(ctx, req.DNASampleRef)
        if dnaResult != nil {
            result.Identified = true
            result.SNISIDPersonID = dnaResult.PersonID
            result.IdentificationMethod = "DNA"
            result.ConfidenceScore = 99.9
        }
    }

    return result, nil
}
```

---

## 4. API REST

| Méthode | Endpoint                                | Rôle          | Description                      |
|---------|-----------------------------------------|---------------|----------------------------------|
| `POST`  | `/api/v1/sigdc/disasters`               | CSPAN, SNGRD  | Déclarer catastrophe             |
| `GET`   | `/api/v1/sigdc/disasters/active`        | Tout SNISID   | Catastrophes en cours            |
| `POST`  | `/api/v1/sigdc/warnings`                | CSPAN, UHM    | Émettre alerte précoce           |
| `POST`  | `/api/v1/sigdc/victims`                 | CSPAN, OIM    | Enregistrer victime              |
| `POST`  | `/api/v1/sigdc/victims/:id/identify`    | CSPAN, AFIS   | Identifier biométriquement       |
| `GET`   | `/api/v1/sigdc/resources/available`     | CSPAN, OCHA   | Ressources disponibles           |
| `GET`   | `/api/v1/sigdc/disasters/:id/dashboard` | CSPAN         | Tableau de bord catastrophe      |

---

## 5. VARIABLES D'ENVIRONNEMENT

```dotenv
SIGDC_DB_HOST=localhost
SIGDC_DB_NAME=snisid_sigdc
SIGDC_SMS_GATEWAY_URL=http://sms-gw:8080
SIGDC_RADIO_API_URL=http://radio-broadcast:8081
SIGDC_AFIS_SERVICE_URL=http://afis-svc:8091
SIGDC_BIO_ADN_SERVICE_URL=http://bio-adn-svc:8080
SIGDC_OCHA_RELIEF_WEB_URL=https://api.reliefweb.int
SIGDC_NHC_STORM_URL=https://www.nhc.noaa.gov/CurrentStorms.json
SIGDC_USGS_EARTHQUAKE_URL=https://earthquake.usgs.gov/fdsnws/event/1
SIGDC_KAFKA_BROKERS=kafka:9092
SIGDC_SERVICE_PORT=8126
```

---
*MP-49 — SIGDC-HT — Gestion Désastres Civils — SNISID — République d'Haïti*
