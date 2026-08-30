# MP-34 — MAR-HT
## Système National de Surveillance Maritime d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-34 | Code : MAR-HT
Dépendances      : PORT-HT (MP-38), TRAF-AR (MP-32), TRAIT-HT (MP-44), SIVC-HT (MP-18)
Normes           : INTERPOL SVD (Stolen Vessels Database), Code ISPS, UNCLOS, CMB
Acteurs          : Garde-Côtes Haïtienne (GCH), USCG liaison, DEA maritime, JIATF-South
```

---

## 1. CONTEXTE

Les 1,771 km de côtes haïtiennes sont presque entièrement non surveillées. Le corridor
Windward Passage (entre Haïti et Cuba) est une route majeure de trafic de cocaïne
vers les USA. L'Île de la Tortue et les baies du Nord sont des points d'entrée
documentés pour armes et drogues. Des milliers d'Haïtiens tentent chaque année la
traversée vers les Bahamas et la Floride dans des embarcations de fortune.

### Zones maritimes prioritaires

| Zone                     | Menaces documentées                              | Priorité    |
|--------------------------|--------------------------------------------------|-------------|
| Windward Passage         | Trafic cocaïne (Colombia → USA via Haiti)        | Critique    |
| Île de la Tortue         | Armes, drogues, migrants clandestins             | Critique    |
| Golfe de la Gonâve       | Go-fasts, migration, trafic d'armes intérieur    | Haute       |
| Baie de Port-au-Prince   | Conteneurs suspects, go-fasts                    | Haute       |
| Côtes Nord (Cap-Haïtien) | Migration vers Bahamas, trafic régional          | Haute       |
| Canal du Vent (Sud)      | Entrée par Les Cayes et Jérémie                  | Moyenne     |

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE mar_vessel_type AS ENUM (
    'CARGO_SHIP','TANKER','FISHING_BOAT','GO_FAST',
    'SAILBOAT','YACHT','FERRY','PATROL_BOAT',
    'WOODEN_BOAT','CANOE','UNKNOWN'
);

CREATE TYPE mar_vessel_status AS ENUM (
    'REGISTERED','STOLEN','SUSPECTED','DETAINED',
    'SUNK','DESTROYED','MISSING','INTERPOL_ALERT'
);

CREATE TYPE mar_incident_type AS ENUM (
    'DRUG_SEIZURE','ARMS_SEIZURE','MIGRANT_INTERDICTION',
    'SMUGGLING','SUSPICIOUS_ACTIVITY','DISTRESS',
    'PIRACY','ILLEGAL_FISHING','HUMAN_TRAFFICKING'
);

CREATE TABLE mar_vessels (
    vessel_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_mar_id     VARCHAR(25) UNIQUE NOT NULL,   -- MAR-HT-NNNNNN
    vessel_name         VARCHAR(150),
    imo_number          VARCHAR(20),                   -- IMO unique ID
    mmsi                VARCHAR(15),                   -- AIS transponder ID
    call_sign           VARCHAR(15),
    vessel_type         mar_vessel_type NOT NULL,
    flag_country        CHAR(3),
    hull_color          VARCHAR(50),
    length_m            DECIMAL(8,2),
    tonnage_gt          INTEGER,
    engine_count        SMALLINT,
    horsepower          INTEGER,                       -- Critique pour go-fasts
    owner_name          VARCHAR(200),
    owner_snisid_id     UUID,
    registration_number VARCHAR(50),
    registration_port   VARCHAR(100),
    status              mar_vessel_status NOT NULL DEFAULT 'REGISTERED',
    gang_id             UUID,
    interpol_svd_ref    VARCHAR(50),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mar_ais_sightings (
    sighting_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vessel_id           UUID REFERENCES mar_vessels(vessel_id),
    mmsi                VARCHAR(15),
    vessel_name         VARCHAR(150),
    sighting_timestamp  TIMESTAMPTZ NOT NULL,
    lat                 DECIMAL(10,7) NOT NULL,
    lng                 DECIMAL(10,7) NOT NULL,
    speed_knots         DECIMAL(5,2),
    heading_degrees     SMALLINT,
    destination         VARCHAR(100),
    source_type         VARCHAR(30),     -- AIS_TERRESTRIAL, AIS_SATELLITE, RADAR, VISUAL
    zone_code           VARCHAR(20),     -- WINDWARD_PASS, TORTUE, GONAVE, etc.
    alert_triggered     BOOLEAN DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mar_incidents (
    incident_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vessel_id           UUID REFERENCES mar_vessels(vessel_id),
    incident_type       mar_incident_type NOT NULL,
    incident_date       TIMESTAMPTZ NOT NULL,
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    zone_desc           VARCHAR(100),
    responding_unit     VARCHAR(50),     -- GCH, USCG, JIATF-South
    outcome             TEXT,
    persons_involved    INTEGER DEFAULT 0,
    snisid_person_ids   UUID[] DEFAULT '{}',
    drug_types          TEXT[] DEFAULT '{}',
    drug_weight_kg      DECIMAL(12,3),
    weapons_found       BOOLEAN DEFAULT FALSE,
    weapons_count       INTEGER DEFAULT 0,
    migrants_count      INTEGER DEFAULT 0,
    biar_refs           UUID[] DEFAULT '{}',
    case_reference      VARCHAR(100),
    photo_refs          TEXT[] DEFAULT '{}',
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mar_watch_vessels (
    watch_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vessel_id           UUID REFERENCES mar_vessels(vessel_id),
    mmsi                VARCHAR(15),
    vessel_name         VARCHAR(150),
    watch_reason        TEXT NOT NULL,
    alert_level         VARCHAR(20) DEFAULT 'CAUTION',
    requesting_unit     VARCHAR(50),
    is_active           BOOLEAN DEFAULT TRUE,
    expiry_date         TIMESTAMPTZ,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_mar_ais_timestamp   ON mar_ais_sightings(sighting_timestamp DESC);
CREATE INDEX idx_mar_ais_mmsi        ON mar_ais_sightings(mmsi);
CREATE INDEX idx_mar_ais_coords      ON mar_ais_sightings(lat, lng);
CREATE INDEX idx_mar_incidents_date  ON mar_incidents(incident_date DESC);
CREATE INDEX idx_mar_incidents_type  ON mar_incidents(incident_type);
CREATE INDEX idx_mar_watch_active    ON mar_watch_vessels(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_mar_vessels_status  ON mar_vessels(status);

COMMIT;
```

---

## 3. SERVICE GO CLÉ — TRAITEMENT AIS EN TEMPS RÉEL

```go
package service

import (
    "context"
    "time"
    "github.com/snisid/mar-svc/internal/domain"
)

// ProcessAISSighting traite un message AIS en temps reel
func (s *MaritimeService) ProcessAISSighting(
    ctx context.Context,
    msg domain.AISMessage,
) error {
    // 1. Verifier si le vessel est sous surveillance
    watchVessel, _ := s.watchRepo.FindByMMSI(ctx, msg.MMSI)
    if watchVessel != nil && watchVessel.IsActive {
        alert := domain.MaritimeAlert{
            MMSI:        msg.MMSI,
            VesselName:  msg.VesselName,
            Lat:         msg.Lat,
            Lng:         msg.Lng,
            AlertLevel:  watchVessel.AlertLevel,
            WatchReason: watchVessel.WatchReason,
            DetectedAt:  time.Now(),
        }
        _ = s.kafka.Publish(ctx, "mar.vessel.alert", alert)
        _ = s.notifier.AlertGCH(ctx, alert)
    }

    // 2. Detection comportement suspect (AIS spoofing, zone interdite)
    if s.isSuspiciousBehavior(msg) {
        _ = s.kafka.Publish(ctx, "mar.suspicious.behavior", msg)
    }

    // 3. Enregistrer le sighting
    return s.aisRepo.SaveSighting(ctx, domain.AISSighting{
        MMSI:              msg.MMSI,
        VesselName:        msg.VesselName,
        SightingTimestamp: time.Now(),
        Lat:               msg.Lat,
        Lng:               msg.Lng,
        SpeedKnots:        msg.SpeedKnots,
        HeadingDegrees:    msg.Heading,
        SourceType:        msg.Source,
        ZoneCode:          s.detectZone(msg.Lat, msg.Lng),
        AlertTriggered:    watchVessel != nil,
    })
}

func (s *MaritimeService) isSuspiciousBehavior(msg domain.AISMessage) bool {
    // Speed > 30 knots = go-fast probable
    if msg.SpeedKnots > 30 {
        return true
    }
    // AIS signal coupe puis rallume = possible camouflage
    lastSighting, _ := s.aisRepo.GetLastSighting(context.Background(), msg.MMSI)
    if lastSighting != nil && time.Since(lastSighting.SightingTimestamp) > 4*time.Hour {
        return true
    }
    return false
}
```

---

## 4. API REST

| Méthode | Endpoint                              | Rôle            | Description                       |
|---------|---------------------------------------|-----------------|-----------------------------------|
| `GET`   | `/api/v1/mar/vessels/:id`             | GCH, DCPJ       | Profil complet vessel             |
| `POST`  | `/api/v1/mar/vessels`                 | GCH_ADMIN       | Enregistrer vessel                |
| `POST`  | `/api/v1/mar/incidents`               | GCH_OFFICER     | Déclarer incident maritime        |
| `GET`   | `/api/v1/mar/incidents/recent`        | GCH, DCPJ       | Incidents récents (24h)           |
| `POST`  | `/api/v1/mar/watch`                   | GCH_SUPERVISOR  | Mettre vessel sous surveillance   |
| `GET`   | `/api/v1/mar/watch/active`            | GCH             | Vessels sous surveillance active  |
| `GET`   | `/api/v1/mar/ais/live`                | GCH             | Positions AIS en direct (SSE)     |
| `GET`   | `/api/v1/mar/zones/:zone/activity`    | GCH, DCPJ       | Activité par zone maritime        |
| `GET`   | `/api/v1/mar/stats/incidents`         | GCH_ADMIN       | Stats incidents par type          |

---

## 5. INTÉGRATIONS

- **PORT-HT** : Arrivées suspectes → vérification cargo et manifestes
- **TRAF-AR** : Cargaisons d'armes interdites → lien BIAR-HT
- **TRAIT-HT** : Embarcations de migrants → identification, assistance
- **INTERPOL SVD** : Vessels volés → sync bidirectionnelle
- **JIATF-South** : Partage info opérationnel via canal sécurisé

---

## 6. VARIABLES D'ENVIRONNEMENT

```dotenv
MAR_DB_HOST=localhost
MAR_DB_NAME=snisid_mar
MAR_REDIS_ADDR=redis-master:6379
MAR_AIS_STREAM_URL=wss://ais-feed.pnh.gov.ht/stream
MAR_INTERPOL_SVD_URL=https://i247-gateway.pnh.gov.ht/svd
MAR_JIATF_API_URL=https://api.jiatfs.southcom.mil
MAR_KAFKA_BROKERS=kafka:9092
MAR_SUSPICIOUS_SPEED_KNOTS=30
MAR_SERVICE_PORT=8107
```

---
*MP-34 — MAR-HT — Surveillance Maritime — SNISID — République d'Haïti*
