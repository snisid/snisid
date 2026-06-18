# MP-37 — AERO-HT
## Registre des Aéronefs Illicites, Pistes Clandestines et Sécurité Aérienne
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-37 | Code : AERO-HT
Dépendances      : SLTD-HT (MP-35), BLKL-HT (MP-36), TRAF-AR (MP-32), MAR-HT (MP-34)
Normes           : OACI Annexe 17, Convention Chicago, FAA registry, INTERPOL Aviation
Acteurs          : Autorité Aéronautique Haïtienne (AAH), DEA Air Wing, DHS CBP Air
```

---

## 1. CONTEXTE

Haïti dispose de l'Aéroport International Toussaint Louverture (PAP) et d'une dizaine
d'aérodromes régionaux. Des dizaines de pistes clandestines ont été identifiées dans
l'Artibonite et le Nord-Ouest pour le trafic de drogues et armes. Des aéronefs non
immatriculés effectuent des vols nocturnes vers les côtes haïtiennes.

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE aero_aircraft_type AS ENUM (
    'COMMERCIAL_JET','TURBOPROP','PISTON_SINGLE','PISTON_TWIN',
    'HELICOPTER','ULTRALIGHT','DRONE_LARGE','UNKNOWN'
);

CREATE TYPE aero_strip_status AS ENUM (
    'ACTIVE','INACTIVE','DESTROYED','LEGALIZED','UNDER_SURVEILLANCE'
);

CREATE TABLE aero_aircraft_registry (
    aircraft_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_mark   VARCHAR(20),
    icao_hex_code       VARCHAR(10),
    aircraft_type       aero_aircraft_type NOT NULL,
    make                VARCHAR(100),
    model               VARCHAR(100),
    manufacture_year    SMALLINT,
    flag_country        CHAR(3),
    owner_name          VARCHAR(200),
    owner_snisid_id     UUID,
    operator_name       VARCHAR(200),
    is_registered       BOOLEAN DEFAULT FALSE,
    is_suspected        BOOLEAN DEFAULT FALSE,
    is_stolen           BOOLEAN DEFAULT FALSE,
    gang_id             UUID,
    drug_trafficking    BOOLEAN DEFAULT FALSE,
    interpol_ref        VARCHAR(50),
    faa_registry_ref    VARCHAR(50),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE aero_clandestine_strips (
    strip_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    strip_name          VARCHAR(150),
    dept_code           CHAR(2) NOT NULL,
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7) NOT NULL,
    lng                 DECIMAL(10,7) NOT NULL,
    length_m            INTEGER,
    surface_type        VARCHAR(30),   -- GRASS, DIRT, ASPHALT, GRAVEL
    status              aero_strip_status NOT NULL DEFAULT 'ACTIVE',
    capable_aircraft    TEXT[] DEFAULT '{}',
    gang_id             UUID,
    first_detected      DATE,
    last_activity_date  DATE,
    source_intel        TEXT,
    satellite_image_ref VARCHAR(500),
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE aero_suspicious_flights (
    flight_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aircraft_id         UUID REFERENCES aero_aircraft_registry(aircraft_id),
    registration_mark   VARCHAR(20),
    flight_date         TIMESTAMPTZ NOT NULL,
    origin_airport      VARCHAR(10),
    destination_airport VARCHAR(10),
    origin_country      CHAR(3),
    destination_country CHAR(3) DEFAULT 'HTI',
    landing_strip_id    UUID REFERENCES aero_clandestine_strips(strip_id),
    landing_location    VARCHAR(300),
    flight_type         VARCHAR(30),   -- DRUG_RUN, ARMS_DELIVERY, UNKNOWN
    cargo_suspected     TEXT,
    source_radar        VARCHAR(50),
    source_informant    BOOLEAN DEFAULT FALSE,
    case_reference      VARCHAR(100),
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_aero_registry_mark    ON aero_aircraft_registry(registration_mark);
CREATE INDEX idx_aero_registry_gang    ON aero_aircraft_registry(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_aero_strips_dept      ON aero_clandestine_strips(dept_code) WHERE status = 'ACTIVE';
CREATE INDEX idx_aero_strips_coords    ON aero_clandestine_strips(lat, lng);
CREATE INDEX idx_aero_flights_date     ON aero_suspicious_flights(flight_date DESC);

COMMIT;
```

---

## 3. API REST

| Méthode | Endpoint                           | Rôle         | Description                      |
|---------|------------------------------------|--------------|----------------------------------|
| `GET`   | `/api/v1/aero/check/:reg`          | AAH, DCPJ    | Vérifier immatriculation         |
| `POST`  | `/api/v1/aero/strips`              | DCPJ_INTEL   | Documenter piste clandestine     |
| `GET`   | `/api/v1/aero/strips/map`          | DCPJ         | Carte GeoJSON pistes             |
| `POST`  | `/api/v1/aero/flights/suspicious`  | DCPJ, AAH    | Signaler vol suspect             |
| `GET`   | `/api/v1/aero/stats/strips`        | DCPJ_ADMIN   | Stats pistes par département     |

---

## 4. VARIABLES D'ENVIRONNEMENT

```dotenv
AERO_DB_HOST=localhost
AERO_DB_NAME=snisid_aero
AERO_FAA_REGISTRY_URL=https://registry.faa.gov/api
AERO_ICAO_RADAR_URL=https://opensky-network.org/api
AERO_GANG_SERVICE_URL=http://gang-svc:8095
AERO_SERVICE_PORT=8109
```

---
*MP-37 — AERO-HT — Aéronefs Illicites — SNISID — République d'Haïti*

---
---

# MP-38 — PORT-HT
## Système de Sécurité Portuaire et Ciblage de Conteneurs
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-38 | Code : PORT-HT
Dépendances      : MAR-HT (MP-34), TRAF-AR (MP-32), BIAR-HT (MP-30), BLAN-HT (MP-40)
Normes           : Code ISPS, OMI résolution MSC.104(73), Programme C-TPAT, WCO SAFE
Acteurs          : APN (Autorité Portuaire Nationale), Douanes, CBP (USA), DCPJ
```

---

## 1. CONTEXTE

Port-au-Prince et Cap-Haïtien sont les deux ports principaux d'Haïti. Ils sont également
des points d'entrée documentés pour la cocaïne (Cap-Haïtien → Anvers, 1,156 kg saisis en
2025). Le système de ciblage des conteneurs vise à identifier les cargaisons suspectes
avant déchargement, en croisant les manifestes avec les renseignements criminels.

### Ports couverts

| Port               | Capacité TEU/an | Risque documenté          |
|--------------------|-----------------|---------------------------|
| Port-au-Prince     | 400,000+        | Drogue, armes, contrebande|
| Cap-Haïtien        | 50,000+         | Transit cocaïne Europe    |
| Gonaïves           | 20,000+         | Contrebande                |
| Les Cayes          | 10,000+         | Petite contrebande         |

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE port_risk_level AS ENUM ('LOW','MEDIUM','HIGH','CRITICAL');
CREATE TYPE port_container_status AS ENUM (
    'PENDING_INSPECTION','CLEARED','HELD_FOR_INSPECTION',
    'SEIZED','RELEASED_AFTER_INSPECTION'
);

CREATE TABLE port_vessels_arrivals (
    arrival_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    port_code           VARCHAR(10) NOT NULL,  -- PAP, CAP, GON, CAY
    vessel_imo          VARCHAR(20),
    vessel_name         VARCHAR(150) NOT NULL,
    flag_country        CHAR(3),
    shipping_company    VARCHAR(200),
    arrival_date        TIMESTAMPTZ NOT NULL,
    origin_port         VARCHAR(100),
    origin_country      CHAR(3),
    container_count     INTEGER DEFAULT 0,
    manifest_ref        VARCHAR(100),
    mar_vessel_id       UUID,              -- Lien MAR-HT si vessel suspecte
    risk_score          SMALLINT DEFAULT 0,
    risk_level          port_risk_level DEFAULT 'LOW',
    cbp_targeting_ref   VARCHAR(50),       -- US Customs Pre-Targeting
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE port_containers (
    container_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    arrival_id          UUID NOT NULL REFERENCES port_vessels_arrivals(arrival_id),
    container_number    VARCHAR(20) NOT NULL,
    container_type      VARCHAR(10),       -- 20GP, 40HC, REEFER, etc.
    declared_content    TEXT NOT NULL,
    declared_weight_kg  DECIMAL(12,3),
    declared_value_usd  DECIMAL(15,2),
    shipper_name        VARCHAR(200),
    shipper_country     CHAR(3),
    consignee_name      VARCHAR(200),
    consignee_snisid_id UUID,
    status              port_container_status NOT NULL DEFAULT 'PENDING_INSPECTION',
    risk_score          SMALLINT DEFAULT 0,
    risk_level          port_risk_level DEFAULT 'LOW',
    risk_flags          TEXT[] DEFAULT '{}',
    selected_for_scan   BOOLEAN DEFAULT FALSE,
    scan_date           TIMESTAMPTZ,
    scan_result         TEXT,
    seized              BOOLEAN DEFAULT FALSE,
    seizure_description TEXT,
    case_reference      VARCHAR(100),
    cbp_targeting_match BOOLEAN DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE port_risk_factors (
    factor_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    container_id        UUID NOT NULL REFERENCES port_containers(container_id),
    factor_type         VARCHAR(50) NOT NULL,
    description         TEXT NOT NULL,
    weight_score        SMALLINT NOT NULL,  -- Points ajoutes au risk_score
    source              VARCHAR(50),        -- BLKL, BLAN, GANG, CBP, INTEL
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Fonction de calcul automatique du score de risque
CREATE OR REPLACE FUNCTION port_compute_risk(
    p_consignee_id UUID,
    p_shipper_country CHAR(3),
    p_declared_content TEXT
) RETURNS INTEGER AS $$
DECLARE
    score INTEGER := 0;
BEGIN
    -- Pays d origine a haut risque (Colombie, Mexique, Venezuela)
    IF p_shipper_country IN ('COL','MEX','VEN','ECU') THEN score := score + 30; END IF;
    -- Consignataire dans BLKL ou SANC
    -- (verif via vue materialisee cross-module)
    IF p_consignee_id IS NOT NULL AND EXISTS (
        SELECT 1 FROM blkl_blacklist
        WHERE snisid_person_id = p_consignee_id AND is_active = TRUE
    ) THEN score := score + 50; END IF;
    -- Contenu generique suspect
    IF lower(p_declared_content) ~ 'general cargo|mixed goods|used items' THEN score := score + 15; END IF;
    RETURN LEAST(score, 100);
END;
$$ LANGUAGE plpgsql;

CREATE INDEX idx_port_containers_risk   ON port_containers(risk_level, status);
CREATE INDEX idx_port_containers_arrival ON port_containers(arrival_id);
CREATE INDEX idx_port_arrivals_date     ON port_vessels_arrivals(arrival_date DESC);
CREATE INDEX idx_port_arrivals_port     ON port_vessels_arrivals(port_code, arrival_date DESC);

COMMIT;
```

---

## 3. SERVICE GO CLÉ — CIBLAGE AUTOMATIQUE

```go
package service

import (
    "context"
    "github.com/snisid/port-svc/internal/domain"
)

func (s *TargetingService) ScoreContainer(
    ctx context.Context,
    container *domain.Container,
) (*domain.RiskAssessment, error) {
    assessment := &domain.RiskAssessment{ContainerID: container.ContainerID}

    // 1. Consignataire dans BLKL ou BLAN
    if container.ConsigneeSnisidID != nil {
        blkl, _ := s.blklClient.CheckPerson(ctx, *container.ConsigneeSnisidID)
        if blkl != nil && blkl.IsBlacklisted {
            assessment.AddFlag("CONSIGNEE_BLACKLISTED", 50)
        }
        blan, _ := s.blanClient.CheckPerson(ctx, *container.ConsigneeSnisidID)
        if blan != nil && blan.HasSuspiciousTransactions {
            assessment.AddFlag("CONSIGNEE_MONEY_LAUNDERING", 30)
        }
    }

    // 2. Route maritime suspecte (pays d origine)
    if s.isHighRiskOrigin(container.ShipperCountry) {
        assessment.AddFlag("HIGH_RISK_ORIGIN_COUNTRY", 30)
    }

    // 3. Correspondance CBP targeting
    cbp, _ := s.cbpClient.CheckContainer(ctx, container.ContainerNumber)
    if cbp != nil && cbp.IsTargeted {
        assessment.AddFlag("CBP_TARGETING_MATCH", 60)
    }

    // 4. Poids vs valeur suspect
    if container.IsWeightValueAnomalous() {
        assessment.AddFlag("WEIGHT_VALUE_ANOMALY", 20)
    }

    assessment.ComputeFinalRisk()
    return assessment, nil
}
```

---

## 4. API REST

| Méthode | Endpoint                           | Rôle          | Description                      |
|---------|------------------------------------|---------------|----------------------------------|
| `POST`  | `/api/v1/port/arrivals`            | APN, DOUANES  | Enregistrer arrivée navire       |
| `GET`   | `/api/v1/port/arrivals/:id`        | DOUANES       | Détail arrivée et conteneurs     |
| `GET`   | `/api/v1/port/containers/high-risk`| DOUANES, DCPJ | Conteneurs à haut risque         |
| `POST`  | `/api/v1/port/containers/:id/scan` | DOUANES       | Enregistrer résultat scan        |
| `POST`  | `/api/v1/port/containers/:id/seize`| DCPJ, DOUANES | Saisir conteneur                 |
| `GET`   | `/api/v1/port/stats/seizures`      | DOUANES_ADMIN | Statistiques saisies             |

---

## 5. VARIABLES D'ENVIRONNEMENT

```dotenv
PORT_DB_HOST=localhost
PORT_DB_NAME=snisid_port
PORT_BLKL_SERVICE_URL=http://blkl-svc:8110
PORT_BLAN_SERVICE_URL=http://blan-svc:8115
PORT_CBP_API_URL=https://api.cbp.dhs.gov/targeting
PORT_CBP_API_KEY=<VAULT:port/cbp_api_key>
PORT_MAR_SERVICE_URL=http://mar-svc:8107
PORT_RISK_SCAN_THRESHOLD=60
PORT_SERVICE_PORT=8111
```

---
*MP-38 — PORT-HT — Sécurité Portuaire — SNISID — République d'Haïti*
