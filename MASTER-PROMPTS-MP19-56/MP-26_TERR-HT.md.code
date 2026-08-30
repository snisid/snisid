# MP-26 — TERR-HT
## Cartographie des Territoires Contrôlés par les Gangs d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-26 | Code : TERR-HT
Dépendances      : GANG-HT (MP-24), SIGEO-HT (MP-48), SIVC-HT (MP-18), CHEF-HT (MP-25)
Normes           : ISO 19115 GIS, OGC GeoJSON/GeoPackage, PostGIS 3.x, OGC WMTS
Acteurs          : DCPJ, BRI, GIPNH, MSP, ONGs humanitaires, MINUSTAH/missions ONU
```

---

## 1. CONTEXTE

En juin 2024, les gangs contrôlaient environ 80% de la zone métropolitaine de
Port-au-Prince selon l'ONU. La cartographie dynamique de ces zones est essentielle
pour : la planification opérationnelle PNH, la sécurité des convois humanitaires,
l'analyse de l'évolution des dynamiques criminelles, et les alertes aux citoyens.

### Degrés de contrôle territorial

| Niveau               | Définition                                  | Impact opérationnel          |
|----------------------|---------------------------------------------|------------------------------|
| FULL_CONTROL         | Gang contrôle 90-100% (checkpoints, taxes)  | Zone interdite accès civil   |
| STRONG_INFLUENCE     | 60-90% — présence régulière de combattants  | Escorte armée obligatoire    |
| CONTESTED            | Zone disputée entre gangs (front de guerre) | Risque extrême — éviter      |
| WEAK_INFLUENCE       | < 30% — présence sporadique                 | Risque modéré               |
| STATE_CONTROLLED     | PNH/MSS rétabli le contrôle                 | Zones d'opération normales   |
| NO_MAN_LAND          | Aucun contrôle stable — population fuie     | Accès humanitaire d'urgence  |

---

## 2. STACK TECHNIQUE

- **PostgreSQL + PostGIS 3.4** : Polygones territoriaux, calculs géospatiaux, SRID 4326
- **ClickHouse** : Historique des changements de territoires, tendances temporelles
- **Martin Tile Server** : Tuiles vectorielles (.mvt) pour MapLibre GL JS
- **Kafka** : Événements modification de territoire en temps réel
- **Frontend** : MapLibre GL JS + PMTiles pour consultation offline (agents terrain)
- **Sources data** : Rapports PNH terrain, ACLED, OCHA, analyse images satellites

---

## 3. BASE DE DONNÉES POSTGIS

```sql
BEGIN;

CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;

CREATE TYPE terr_control_level AS ENUM (
    'FULL_CONTROL', 'STRONG_INFLUENCE', 'CONTESTED',
    'WEAK_INFLUENCE', 'STATE_CONTROLLED', 'NO_MAN_LAND'
);

CREATE TYPE terr_source AS ENUM (
    'PNH_FIELD_REPORT', 'SATELLITE_ANALYSIS', 'INFORMANT',
    'NGO_REPORT', 'ACLED', 'MEDIA_CROSS_CHECK', 'LAPI_ANALYSIS'
);

CREATE TABLE terr_zones (
    zone_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gang_id             UUID NOT NULL,
    zone_name           VARCHAR(150),
    dept_code           CHAR(2) NOT NULL,
    commune             VARCHAR(100),
    section_communale   VARCHAR(100),

    geom                GEOMETRY(MultiPolygon, 4326) NOT NULL,
    area_km2            DECIMAL(10,3),
    centroid_lat        DECIMAL(10,7),
    centroid_lng        DECIMAL(10,7),

    control_level       terr_control_level NOT NULL,
    estimated_population INTEGER,
    strategic_importance SMALLINT CHECK (strategic_importance BETWEEN 1 AND 10),

    controls_national_road BOOLEAN DEFAULT FALSE,
    road_numbers        TEXT[] DEFAULT '{}',
    controls_port       BOOLEAN DEFAULT FALSE,
    controls_airport    BOOLEAN DEFAULT FALSE,
    controls_market     BOOLEAN DEFAULT FALSE,

    valid_from          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_to            TIMESTAMPTZ,
    is_current          BOOLEAN DEFAULT TRUE,
    intelligence_source terr_source NOT NULL,
    confidence_level    SMALLINT CHECK (confidence_level BETWEEN 1 AND 10),
    analyst_notes       TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Trigger pour calculer area et centroid
CREATE OR REPLACE FUNCTION terr_compute_geometry_props()
RETURNS TRIGGER AS $$
BEGIN
    NEW.area_km2 := ST_Area(NEW.geom::geography) / 1000000;
    NEW.centroid_lat := ST_Y(ST_Centroid(NEW.geom));
    NEW.centroid_lng := ST_X(ST_Centroid(NEW.geom));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_terr_geometry
    BEFORE INSERT OR UPDATE ON terr_zones
    FOR EACH ROW EXECUTE FUNCTION terr_compute_geometry_props();

-- Historique des changements
CREATE TABLE terr_zone_history (
    history_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_id             UUID NOT NULL REFERENCES terr_zones(zone_id),
    change_type         VARCHAR(30) NOT NULL,
    previous_control    terr_control_level,
    new_control         terr_control_level,
    change_date         TIMESTAMPTZ NOT NULL,
    trigger_event       TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Checkpoints et barrages routiers connus
CREATE TABLE terr_checkpoints (
    checkpoint_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gang_id             UUID NOT NULL,
    location            GEOMETRY(Point, 4326) NOT NULL,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    road_number         VARCHAR(20),
    is_armed            BOOLEAN DEFAULT TRUE,
    extortion_type      VARCHAR(100),
    reported_at         TIMESTAMPTZ NOT NULL,
    is_active           BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index spatiaux critiques
CREATE INDEX idx_terr_zones_geom        ON terr_zones USING GIST(geom) WHERE is_current = TRUE;
CREATE INDEX idx_terr_zones_dept        ON terr_zones(dept_code) WHERE is_current = TRUE;
CREATE INDEX idx_terr_zones_gang        ON terr_zones(gang_id) WHERE is_current = TRUE;
CREATE INDEX idx_terr_zones_control     ON terr_zones(control_level) WHERE is_current = TRUE;
CREATE INDEX idx_terr_checkpoints_geom  ON terr_checkpoints USING GIST(location) WHERE is_active = TRUE;

COMMIT;
```

---

## 4. SERVICE GO CLÉ — VÉRIFICATION POINT GÉOGRAPHIQUE

```go
package service

import (
    "context"
    "github.com/snisid/terr-svc/internal/domain"
)

// CheckPointSafety verifie si une coord est en zone de gang
func (s *TerritoryService) CheckPointSafety(
    ctx context.Context,
    lat, lng float64,
) (*domain.SafetyCheckResult, error) {
    zones, err := s.repo.FindZonesContainingPoint(ctx, lat, lng)
    result := &domain.SafetyCheckResult{
        Lat: lat, Lng: lng, IsSafe: true,
    }
    if err != nil || len(zones) == 0 {
        return result, nil
    }
    result.IsSafe = false
    for _, z := range zones {
        result.Zones = append(result.Zones, domain.ZoneInfo{
            GangName:     z.GangName,
            ControlLevel: z.ControlLevel,
            RiskScore:    z.StrategicImportance,
        })
        if z.ControlLevel == "FULL_CONTROL" || z.ControlLevel == "CONTESTED" {
            result.HighRisk = true
        }
    }
    // Verifier checkpoints proches (rayon 500m)
    checkpoints, _ := s.repo.FindNearbyCheckpoints(ctx, lat, lng, 500)
    result.NearbyCheckpoints = len(checkpoints)
    return result, nil
}

// GetRouteSafety evalue la securite d un itineraire complet
func (s *TerritoryService) GetRouteSafety(
    ctx context.Context,
    waypoints []domain.Point,
) (*domain.RouteSafetyResult, error) {
    result := &domain.RouteSafetyResult{
        Waypoints:  waypoints,
        SafeSegments: 0,
        DangerSegments: 0,
    }
    for _, pt := range waypoints {
        check, _ := s.CheckPointSafety(ctx, pt.Lat, pt.Lng)
        if check.HighRisk {
            result.DangerSegments++
        } else {
            result.SafeSegments++
        }
    }
    result.OverallRisk = s.computeOverallRisk(result)
    return result, nil
}
```

---

## 5. API REST

| Méthode | Endpoint                            | Rôle         | Description                         |
|---------|-------------------------------------|--------------|-------------------------------------|
| `GET`   | `/api/v1/terr/check`                | PNH, ONG     | Vérifier sécurité d'un point (lat,lng)|
| `GET`   | `/api/v1/terr/route-safety`         | PNH, ONG     | Sécurité d'un itinéraire            |
| `GET`   | `/api/v1/terr/zones`                | DCPJ         | Toutes zones actives (GeoJSON)      |
| `GET`   | `/api/v1/terr/zones/dept/:code`     | PNH          | Zones par département               |
| `POST`  | `/api/v1/terr/zones`                | DCPJ_INTEL   | Créer/mettre à jour zone            |
| `GET`   | `/api/v1/terr/zones/:gang_id`       | DCPJ         | Territoire d'un gang                |
| `POST`  | `/api/v1/terr/checkpoints`          | PNH_OFFICER  | Signaler checkpoint gang            |
| `GET`   | `/api/v1/terr/tiles/{z}/{x}/{y}.mvt`| FRONTEND     | Tuiles vectorielles MapLibre        |
| `GET`   | `/api/v1/terr/history/:zone_id`     | DCPJ_INTEL   | Historique changements d'une zone   |

---

## 6. VARIABLES D'ENVIRONNEMENT

```dotenv
TERR_DB_HOST=localhost
TERR_DB_NAME=snisid_terr
TERR_POSTGIS_SRID=4326
TERR_CLICKHOUSE_ADDR=clickhouse:9000
TERR_KAFKA_BROKERS=kafka:9092
TERR_GANG_SERVICE_URL=http://gang-svc:8095
TERR_TILE_CACHE_TTL_MIN=30
TERR_SERVICE_PORT=8098
```

---
*MP-26 — TERR-HT — Cartographie Territoriale Gangs — SNISID — République d'Haïti*
