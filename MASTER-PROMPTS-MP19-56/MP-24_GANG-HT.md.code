# MP-24 — GANG-HT
## Registre National des Organisations Criminelles et Gangs d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-24 | Code : GANG-HT | Version : 1.0.0
Dépendances      : CHEF-HT (MP-25), TERR-HT (MP-26), RDEP-HT (MP-22), RESO-HT (MP-28), SIVC-HT (MP-18)
Normes           : INTERPOL Criminal Organizations, Résolution CSNU 2653, OFAC SDN
Acteurs          : DCPJ BAC, Cellule Renseignement (CIR), DEA liaison, Panel experts ONU
```

---

## 1. CONTEXTE ET PANORAMA DES GANGS HAÏTIENS

| Groupe           | Territoire principal         | Effectif estimé | Capacité armée | Statut OFAC/ONU      |
|------------------|------------------------------|-----------------|----------------|----------------------|
| Viv Ansanm       | P-au-P métropole (coalition) | 3,000-5,000     | Très élevée    | ONU 2653 (partiel)   |
| G9 an Fanm       | Cité Soleil, Martissant      | 500-800         | Élevée         | Leader OFAC désigné  |
| G-Pep            | Nord P-au-P                  | 400-600         | Élevée         | Non désigné          |
| 400 Mawozo       | Croix-des-Bouquets, Est      | 400-600         | Élevée         | OFAC désigné         |
| Gran Grif        | Artibonite (Liancourt)       | 200-400         | Moyenne        | Non désigné          |
| Bèlè             | Sud (Les Cayes)              | 100-200         | Moyenne        | Non désigné          |
| Fantom 509       | Centre-ville P-au-P          | 150-300         | Élevée         | Non désigné          |

---

## 2. ARCHITECTURE

```
services/gang-svc/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── organization.go
│   │   ├── incident.go
│   │   ├── alliance.go
│   │   └── enums.go
│   ├── repository/
│   │   ├── postgres/organization_repo.go
│   │   ├── neo4j/graph_repo.go
│   │   └── clickhouse/analytics_repo.go
│   ├── service/
│   │   ├── organization_service.go
│   │   ├── incident_service.go
│   │   └── alliance_service.go
│   └── api/rest/
│       ├── organization_handler.go
│       ├── incident_handler.go
│       └── intel_handler.go
└── Dockerfile
```

---

## 3. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE gang_structure_type AS ENUM (
    'HIERARCHY','NETWORK','CELL','COALITION','FRANCHISE'
);

CREATE TYPE gang_activity_level AS ENUM (
    'DORMANT','LOW','MODERATE','HIGH','EXTREME'
);

CREATE TYPE gang_primary_activity AS ENUM (
    'KIDNAPPING','DRUG_TRAFFICKING','ARMS_TRAFFICKING',
    'EXTORTION','TERRITORY_CONTROL','CONTRACT_KILLING',
    'HUMAN_TRAFFICKING','MONEY_LAUNDERING','MIXED'
);

CREATE TABLE gang_organizations (
    gang_id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_gang_id         VARCHAR(25) UNIQUE NOT NULL,  -- GANG-HT-NNNNNN
    name                     VARCHAR(150) NOT NULL,
    aliases                  TEXT[] DEFAULT '{}',
    structure_type           gang_structure_type,
    primary_activity         gang_primary_activity NOT NULL,
    activity_level           gang_activity_level NOT NULL DEFAULT 'HIGH',
    estimated_members        INTEGER,
    armed_members_pct        SMALLINT,
    heavy_weapons            BOOLEAN DEFAULT FALSE,
    primary_dept_code        CHAR(2) NOT NULL,
    territory_communes       TEXT[] DEFAULT '{}',
    territory_geojson        JSONB,
    estimated_revenue_usd_monthly DECIMAL(12,2),
    primary_income_sources   TEXT[] DEFAULT '{}',
    un_designation_date      TIMESTAMPTZ,
    ofac_designation         BOOLEAN DEFAULT FALSE,
    ofac_sdn_ref             VARCHAR(50),
    allied_gang_ids          UUID[] DEFAULT '{}',
    rival_gang_ids           UUID[] DEFAULT '{}',
    established_date         DATE,
    current_leader_id        UUID,
    intel_confidence         SMALLINT CHECK (intel_confidence BETWEEN 1 AND 10),
    last_intel_update        TIMESTAMPTZ,
    is_active                BOOLEAN DEFAULT TRUE,
    created_by               UUID NOT NULL,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE gang_incidents (
    incident_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gang_id             UUID NOT NULL REFERENCES gang_organizations(gang_id),
    incident_type       VARCHAR(50) NOT NULL,
    incident_date       TIMESTAMPTZ NOT NULL,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    casualties          SMALLINT DEFAULT 0,
    victim_ids          UUID[] DEFAULT '{}',
    sivc_alert_id       UUID,
    description         TEXT,
    intelligence_source VARCHAR(100),
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE gang_alliances (
    alliance_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gang_a_id           UUID NOT NULL REFERENCES gang_organizations(gang_id),
    gang_b_id           UUID NOT NULL REFERENCES gang_organizations(gang_id),
    alliance_type       VARCHAR(30) NOT NULL,
    start_date          DATE,
    end_date            DATE,
    confidence_level    SMALLINT,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT no_self_alliance CHECK (gang_a_id <> gang_b_id)
);

CREATE INDEX idx_gang_dept     ON gang_organizations(primary_dept_code) WHERE is_active = TRUE;
CREATE INDEX idx_gang_activity ON gang_organizations(activity_level) WHERE is_active = TRUE;
CREATE INDEX idx_gang_ofac     ON gang_organizations(ofac_designation) WHERE ofac_designation = TRUE;
CREATE INDEX idx_gang_incidents ON gang_incidents(gang_id, incident_date DESC);
CREATE INDEX idx_gang_incid_dept ON gang_incidents(dept_code, incident_date DESC);

COMMIT;
```

---

## 4. MODÈLE NEO4J — GRAPHE ORGANISATIONNEL

```cypher
// Noeud Gang avec membres et alliances
(:Gang {gang_id, name, primary_activity, activity_level, territory})-
    [:LED_BY {since}]->(:Member {member_id, name, role, aliases: []})
(:Gang)-[:ALLIED_WITH {type, since, confidence}]->(:Gang)
(:Gang)-[:RIVALS_WITH {since}]->(:Gang)
(:Gang)-[:OPERATES_IN]->(:Location {dept_code, commune})
(:Gang)-[:INVOLVED_IN {role, date}]->(:Crime {type, date, casualties})
(:Gang)-[:USES_VEHICLE {role}]->(:Vehicle {plate_number})

-- Requete: coalitions entre gangs (partenaires de Viv Ansanm)
MATCH (g1:Gang)-[:ALLIED_WITH]-(g2:Gang)
WHERE g1.name = 'Viv Ansanm'
RETURN g2.name, g2.territory, g2.activity_level
ORDER BY g2.activity_level DESC
```

---

## 5. API REST

| Méthode | Endpoint                         | Rôle            | Description                      |
|---------|----------------------------------|-----------------|----------------------------------|
| `POST`  | `/api/v1/gangs`                  | DCPJ_INTEL      | Créer fiche organisation         |
| `GET`   | `/api/v1/gangs`                  | DCPJ, BAC       | Lister organisations actives     |
| `GET`   | `/api/v1/gangs/:id`              | DCPJ, BAC       | Détail complet organisation      |
| `POST`  | `/api/v1/gangs/:id/incidents`    | DCPJ, BAC       | Enregistrer incident             |
| `GET`   | `/api/v1/gangs/:id/members`      | DCPJ_INTEL      | Membres identifiés               |
| `GET`   | `/api/v1/gangs/by-dept/:code`    | DCPJ, BRI       | Gangs actifs par département     |
| `GET`   | `/api/v1/gangs/alliances/map`    | DCPJ_INTEL      | Carte des alliances              |
| `GET`   | `/api/v1/gangs/sanctioned`       | DCPJ, MJSP      | Gangs sous sanctions ONU/OFAC    |

---

## 6. ANALYTIQUES CLICKHOUSE

```sql
-- Table de faits gang-incidents dans ClickHouse
CREATE TABLE gang_incident_facts (
    incident_id UUID, gang_id UUID, gang_name String,
    incident_type String, dept_code FixedString(2),
    commune String, incident_date Date,
    casualties Int16, lat Float64, lng Float64
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(incident_date)
ORDER BY (dept_code, incident_type, incident_date);

-- Vue: incidents par dept et type (mensuel)
CREATE MATERIALIZED VIEW gang_incidents_monthly
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMM(month)
ORDER BY (month, dept_code, incident_type)
AS SELECT
    toStartOfMonth(incident_date) AS month,
    dept_code, incident_type,
    count() AS count,
    sum(casualties) AS total_casualties
FROM gang_incident_facts
GROUP BY month, dept_code, incident_type;
```

---

## 7. INTÉGRATIONS

- **CHEF-HT** : `current_leader_id` → profil chef dans CHEF-HT
- **TERR-HT** : `territory_geojson` → couche cartographique GIS temps réel
- **SIVC-HT** : Véhicules liés au gang → alertes LAPI automatiques
- **SANC-HT** : `ofac_sdn_ref` → vérification croisée OFAC/ONU
- **RESO-HT** : Graphe Neo4j → analyse réseau inter-gangs

---

## 8. VARIABLES D'ENVIRONNEMENT

```dotenv
GANG_DB_HOST=localhost
GANG_DB_NAME=snisid_gang
GANG_NEO4J_URI=bolt://neo4j:7687
GANG_CLICKHOUSE_ADDR=clickhouse:9000
GANG_KAFKA_BROKERS=kafka:9092
GANG_SANC_SERVICE_URL=http://sanc-svc:8100
GANG_SERVICE_PORT=8095
```

---
*MP-24 — GANG-HT — Registre Organisations Criminelles — SNISID — République d'Haïti*
