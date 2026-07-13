# MP-31 — EXPL-HT
## Registre National des Explosifs, IED et Munitions
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-31 | Code : EXPL-HT
Dépendances      : BIAR-HT (MP-30), SIAR-HT (MP-29), SNISID-BIO-ADN, PORT-HT (MP-38)
Normes           : INTERPOL EXPLOINT, Protocol Nairobi, Convention Ottawa mines
Acteurs          : DCPJ Déminage, PNH EOD, Douanes, Forces Armées d'Haïti (FAd'H)
```

---

## 1. CONTEXTE

Les gangs haïtiens utilisent des engins explosifs improvisés (IED), des grenades et
roquettes RPG. Ce module trace tous les explosifs : stocks légaux, matériaux volés,
IED récupérés, et incidents d'explosion pour analyses forensiques.

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE expl_type AS ENUM (
    'IED','GRENADE','RPG','MORTAR','LANDMINE','DYNAMITE',
    'BLASTING_CAP','AMMUNITION_BULK','MILITARY_ORDNANCE','UNKNOWN'
);

CREATE TYPE expl_status AS ENUM (
    'RECOVERED','DESTROYED','DETONATED','STORED_EVIDENCE','TRANSFERRED'
);

CREATE TABLE expl_incidents (
    incident_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_expl_id    VARCHAR(25) UNIQUE NOT NULL,
    incident_type       VARCHAR(30) NOT NULL,     -- FIND, DETONATION, SEIZURE, SURRENDER
    explosive_type      expl_type NOT NULL,
    status              expl_status NOT NULL DEFAULT 'RECOVERED',
    quantity            INTEGER DEFAULT 1,
    weight_kg           DECIMAL(10,3),
    manufacturer        VARCHAR(100),
    lot_number          VARCHAR(50),
    manufacture_country CHAR(3),
    estimated_date      DATE,
    incident_date       TIMESTAMPTZ NOT NULL,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    responding_unit     VARCHAR(50),
    eod_officer         UUID,
    casualties          SMALLINT DEFAULT 0,
    gang_id             UUID,
    from_person_id      UUID,
    case_reference      VARCHAR(100),
    dna_sample_taken    BOOLEAN DEFAULT FALSE,
    bio_sample_ref      VARCHAR(100),
    photo_refs          TEXT[] DEFAULT '{}',
    interpol_exploint_ref VARCHAR(50),
    notes               TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE expl_legal_stocks (
    stock_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    holder_entity       VARCHAR(200) NOT NULL,
    holder_type         VARCHAR(30),     -- MINING, CONSTRUCTION, MILITARY, POLICE
    explosive_type      expl_type NOT NULL,
    quantity_kg         DECIMAL(12,3) NOT NULL,
    storage_location    TEXT NOT NULL,
    dept_code           CHAR(2),
    license_ref         VARCHAR(50),
    last_audit_date     DATE,
    next_audit_date     DATE,
    is_secured          BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_expl_type  ON expl_incidents(explosive_type, incident_date DESC);
CREATE INDEX idx_expl_dept  ON expl_incidents(dept_code);
CREATE INDEX idx_expl_gang  ON expl_incidents(gang_id) WHERE gang_id IS NOT NULL;

COMMIT;
```

---

## 3. API REST

| Méthode | Endpoint                          | Rôle            | Description                   |
|---------|-----------------------------------|-----------------|-------------------------------|
| `POST`  | `/api/v1/expl/incidents`          | DCPJ, FAd'H EOD | Déclarer incident explosif    |
| `GET`   | `/api/v1/expl/incidents/:id`      | DCPJ, FAd'H     | Détail incident               |
| `GET`   | `/api/v1/expl/incidents/by-dept`  | DCPJ            | Incidents par département     |
| `GET`   | `/api/v1/expl/legal-stocks`       | DCPJ_ADMIN      | Stocks légaux enregistrés     |
| `POST`  | `/api/v1/expl/legal-stocks`       | DCPJ_ADMIN      | Déclarer stock légal          |

---

## 4. VARIABLES D'ENVIRONNEMENT

```dotenv
EXPL_DB_HOST=localhost
EXPL_DB_NAME=snisid_expl
EXPL_INTERPOL_GATEWAY=https://i247-gateway.pnh.gov.ht/exploint
EXPL_BIO_ADN_URL=http://bio-adn-svc:8080
EXPL_SERVICE_PORT=8104
```

---
*MP-31 — EXPL-HT — Explosifs et Munitions — SNISID — République d'Haïti*

---
---

# MP-32 — TRAF-AR
## Analyse des Routes de Trafic d'Armes vers Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-32 | Code : TRAF-AR
Dépendances      : BIAR-HT (MP-30), SIAR-HT (MP-29), MAR-HT (MP-34), PORT-HT (MP-38)
Normes           : Traité sur le Commerce des Armes (TCA), Résolution CSNU 2653
Acteurs          : DCPJ, ATF/DEA liaison, UNODC, Douanes, BCN INTERPOL HTI
```

---

## 1. CONTEXTE

Corridors documentés du trafic d'armes vers Haïti :
- **Miami → Haïti** : Achats de façade (straw purchases) en Floride — voie principale
- **Géorgie → Haïti** : Gun shops Atlanta/Savannah, expédition maritime
- **Jamaïque ↔ Haïti** : Échanges armes-drogues (guns-for-drugs)
- **République Dominicaine → Haïti** : Passages terrestres Malpasse/Ouanaminthe
- **Bateau → Côtes** : Go-fasts depuis Bahamas et côtes USA

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE trafar_route_type AS ENUM (
    'MARITIME_DIRECT','MARITIME_VIA_BAHAMAS',
    'AIR_CARGO','AIR_PASSENGER',
    'LAND_BORDER_DOM','LAND_BORDER_OTHER',
    'POSTAL','MIXED'
);

CREATE TYPE trafar_method AS ENUM (
    'STRAW_PURCHASE','STOLEN_DIVERTED','CORRUPT_OFFICIAL',
    'FALSE_END_USER','DARK_WEB','DIPLOMATIC_POUCH',
    'CONCEALED_CARGO','DRUGS_FOR_GUNS_SWAP','UNKNOWN'
);

CREATE TABLE trafar_routes (
    route_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_name          VARCHAR(150) NOT NULL,
    route_type          trafar_route_type NOT NULL,
    trafficking_method  trafar_method NOT NULL,
    origin_country      CHAR(3) NOT NULL,
    origin_city         VARCHAR(100),
    transit_points      JSONB,       -- [{country, city, transport_mode}]
    entry_point_haiti   VARCHAR(100),
    entry_dept_code     CHAR(2),
    associated_gang_ids UUID[] DEFAULT '{}',
    known_suppliers     TEXT[] DEFAULT '{}',
    activity_level      VARCHAR(20) DEFAULT 'ACTIVE',
    estimated_volume_monthly INTEGER,
    weapon_types        TEXT[] DEFAULT '{}',
    intel_confidence    SMALLINT CHECK (intel_confidence BETWEEN 1 AND 10),
    first_detected      DATE,
    last_confirmed      DATE,
    linked_case_refs    TEXT[] DEFAULT '{}',
    biar_weapon_ids     UUID[] DEFAULT '{}',
    atf_case_refs       TEXT[] DEFAULT '{}',
    unodc_ref           VARCHAR(50),
    analyst_notes       TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE trafar_shipments (
    shipment_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id            UUID REFERENCES trafar_routes(route_id),
    shipment_date       TIMESTAMPTZ NOT NULL,
    intercepted         BOOLEAN DEFAULT FALSE,
    interception_date   TIMESTAMPTZ,
    interception_location VARCHAR(300),
    interception_unit   VARCHAR(50),
    weapons_count       INTEGER,
    weapons_types       TEXT[] DEFAULT '{}',
    estimated_value_usd DECIMAL(12,2),
    linked_persons      UUID[] DEFAULT '{}',
    port_ht_ref         UUID,       -- Lien PORT-HT si interception portuaire
    mar_ht_ref          UUID,       -- Lien MAR-HT si interception maritime
    case_reference      VARCHAR(100),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE trafar_suppliers (
    supplier_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supplier_name       VARCHAR(200),
    supplier_type       VARCHAR(50),   -- DEALER, CARTEL, CORRUPT_OFFICIAL, UNKNOWN
    country             CHAR(3) NOT NULL,
    city                VARCHAR(100),
    snisid_person_id    UUID,          -- Si identifie dans SNISID
    linked_routes       UUID[] DEFAULT '{}',
    atf_subject_ref     VARCHAR(50),
    interpol_notice_ref VARCHAR(50),
    is_active           BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_trafar_routes_type     ON trafar_routes(route_type, activity_level);
CREATE INDEX idx_trafar_routes_origin   ON trafar_routes(origin_country);
CREATE INDEX idx_trafar_routes_entry    ON trafar_routes(entry_dept_code);
CREATE INDEX idx_trafar_shipments_route ON trafar_shipments(route_id);
CREATE INDEX idx_trafar_shipments_date  ON trafar_shipments(shipment_date DESC);

COMMIT;
```

---

## 3. API REST

| Méthode | Endpoint                              | Rôle         | Description                     |
|---------|---------------------------------------|--------------|---------------------------------|
| `GET`   | `/api/v1/trafar/routes`               | DCPJ_INTEL   | Toutes routes actives           |
| `GET`   | `/api/v1/trafar/routes/:id`           | DCPJ_INTEL   | Détail route                    |
| `POST`  | `/api/v1/trafar/routes`               | DCPJ_INTEL   | Documenter nouvelle route       |
| `POST`  | `/api/v1/trafar/shipments`            | DCPJ, DOUANES| Enregistrer envoi intercepté    |
| `GET`   | `/api/v1/trafar/map`                  | DCPJ_INTEL   | Carte GeoJSON des routes        |
| `GET`   | `/api/v1/trafar/stats/by-origin`      | DCPJ, ATF    | Stats par pays d'origine        |
| `GET`   | `/api/v1/trafar/suppliers`            | DCPJ_INTEL   | Fournisseurs identifiés         |

---

## 4. VARIABLES D'ENVIRONNEMENT

```dotenv
TRAFAR_DB_HOST=localhost
TRAFAR_DB_NAME=snisid_trafar
TRAFAR_ATF_URL=https://etrace.atf.gov/api
TRAFAR_UNODC_URL=https://www.unodc.org/api
TRAFAR_MAR_SERVICE_URL=http://mar-svc:8107
TRAFAR_PORT_SERVICE_URL=http://port-svc:8111
TRAFAR_SERVICE_PORT=8105
```

---
*MP-32 — TRAF-AR — Routes Trafic d'Armes — SNISID — République d'Haïti*
