# MP-30 — BIAR-HT
## Base Nationale des Armes Illicites et Interface INTERPOL iARMS
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-30 | Code : BIAR-HT
Dépendances      : SIAR-HT (MP-29), TRAF-AR (MP-32), GANG-HT (MP-24), PORT-HT (MP-38)
Normes           : INTERPOL iARMS (1.5M records), TCA, Programme Action onusien sur SALW
Acteurs          : DCPJ, Douanes, BCN INTERPOL Port-au-Prince, ATF/DEA liaison USA
```

---

## 1. CONTEXTE

BIAR-HT centralise toutes les armes illicites : saisies, récupérées, rapportées par les
partenaires internationaux. La base INTERPOL iARMS contient plus de 1.5 million d'armes
volées ou perdues. Ce module est le nœud national haïtien pour les échanges avec iARMS.

### Sources d'alimentation de BIAR-HT

| Source                | Type de données                                 |
|-----------------------|------------------------------------------------|
| SIAR-HT               | Armes saisies localement, vols déclarés         |
| INTERPOL iARMS        | Armes signalées par 194 pays membres            |
| ATF eTRACE (USA)      | Traçage armes d'origine américaine              |
| Douanes (PORT-HT)     | Interceptions de conteneurs et fret aérien      |
| MINUSTAH/BINUH        | Rapports de désarmement                         |

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE biar_recovery_context AS ENUM (
    'POLICE_OPERATION','CHECKPOINT','PORT_SEIZURE','AIRPORT_SEIZURE',
    'COMMUNITY_SURRENDER','CRIME_SCENE','RAID','BORDER_SEIZURE','OTHER'
);

CREATE TYPE biar_weapon_disposition AS ENUM (
    'HELD_AS_EVIDENCE','DESTROYED','RETURNED_TO_OWNER',
    'TRANSFERRED_TO_POLICE','SENT_TO_INTERPOL','PENDING'
);

CREATE TABLE biar_illicit_weapons (
    weapon_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_biar_id    VARCHAR(25) UNIQUE NOT NULL,     -- BIAR-HT-NNNNNN
    serial_number       VARCHAR(100),
    serial_obliterated  BOOLEAN DEFAULT FALSE,           -- Numero efface (crime)
    make                VARCHAR(100),
    model               VARCHAR(100),
    caliber             VARCHAR(30),
    weapon_type         VARCHAR(50) NOT NULL,
    manufacture_country CHAR(3),
    estimated_manufacture_year SMALLINT,

    -- Decouverte / saisie
    recovery_date       TIMESTAMPTZ NOT NULL,
    recovery_context    biar_recovery_context NOT NULL,
    recovery_location   VARCHAR(300),
    recovery_dept_code  CHAR(2),
    recovery_commune    VARCHAR(100),
    recovery_lat        DECIMAL(10,7),
    recovery_lng        DECIMAL(10,7),
    seizing_unit        VARCHAR(50) NOT NULL,
    seizing_officer     UUID,
    case_reference      VARCHAR(100),

    -- Liens criminels
    from_person_id      UUID,                            -- Personne chez qui saisie
    gang_id             UUID,
    crime_category      VARCHAR(50),
    associated_cases    TEXT[] DEFAULT '{}',

    -- Traçage international
    origin_country      CHAR(3),
    transit_countries   CHAR(3)[] DEFAULT '{}',
    trafficking_route   TEXT,
    import_method       TEXT,                            -- Conteneur, bagages, go-fast...

    -- INTERPOL / ATF
    iarms_ref           VARCHAR(50),
    atf_etrace_ref      VARCHAR(50),
    reported_to_interpol BOOLEAN DEFAULT FALSE,
    interpol_reported_at TIMESTAMPTZ,

    -- Sort de l arme
    disposition         biar_weapon_disposition DEFAULT 'HELD_AS_EVIDENCE',
    disposal_date       TIMESTAMPTZ,
    disposal_auth       UUID,

    quantity_ammunition INTEGER DEFAULT 0,
    ammunition_type     VARCHAR(50),
    photos_refs         TEXT[] DEFAULT '{}',
    notes               TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE biar_batch_seizures (
    batch_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_reference     VARCHAR(50) UNIQUE NOT NULL,
    operation_name      TEXT,
    seizure_date        TIMESTAMPTZ NOT NULL,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    total_weapons       INTEGER NOT NULL,
    weapon_ids          UUID[] DEFAULT '{}',
    seizing_unit        VARCHAR(50) NOT NULL,
    lead_officer        UUID,
    partnering_agencies TEXT[] DEFAULT '{}',
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE biar_iarms_sync_log (
    sync_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    weapon_id           UUID REFERENCES biar_illicit_weapons(weapon_id),
    direction           VARCHAR(10) NOT NULL,     -- OUTBOUND / INBOUND
    iarms_ref           VARCHAR(50),
    sync_status         VARCHAR(20) DEFAULT 'PENDING',
    synced_at           TIMESTAMPTZ,
    error_message       TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_biar_serial    ON biar_illicit_weapons(serial_number) WHERE serial_number IS NOT NULL;
CREATE INDEX idx_biar_gang      ON biar_illicit_weapons(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_biar_dept      ON biar_illicit_weapons(recovery_dept_code);
CREATE INDEX idx_biar_date      ON biar_illicit_weapons(recovery_date DESC);
CREATE INDEX idx_biar_iarms     ON biar_illicit_weapons(iarms_ref) WHERE iarms_ref IS NOT NULL;
CREATE INDEX idx_biar_origin    ON biar_illicit_weapons(origin_country);

COMMIT;
```

---

## 3. SERVICE GO CLÉ — SOUMISSION INTERPOL iARMS

```go
package service

import (
    "context"
    "fmt"
    "time"
    "github.com/snisid/biar-svc/internal/domain"
)

type IARMSClient struct {
    gatewayURL string
    apiKey     string
    ncbCode    string   // HTI
}

func (c *IARMSClient) SubmitIllicitWeapon(
    ctx context.Context,
    w *domain.IllicitWeapon,
) (string, error) {
    record := IARMSRecord{
        NCBRef:         w.NationalBIARID,
        OriginCountry:  "HTI",
        SerialNumber:   w.SerialNumber,
        Make:           w.Make,
        Model:          w.Model,
        Caliber:        w.Caliber,
        WeaponType:     w.WeaponType,
        RecoveryDate:   w.RecoveryDate.Format("2006-01-02"),
        RecoveryCountry: "HTI",
        Notes:          fmt.Sprintf("Recovery context: %s, Unit: %s", w.RecoveryContext, w.SeizingUnit),
    }
    // POST vers gateway I-24/7 national (BCN Port-au-Prince)
    iarmRef, err := c.postToGateway(ctx, record)
    if err != nil {
        return "", fmt.Errorf("soumission iARMS: %w", err)
    }
    return iarmRef, nil
}

func (s *WeaponService) SyncFromIARMS(ctx context.Context) (*domain.SyncResult, error) {
    // Recuperer nouvelles entrees iARMS pertinentes pour Haiti (transit, origine HTI)
    entries, err := s.iarmsClient.FetchRecentEntries(ctx, "HTI")
    if err != nil {
        return nil, err
    }
    result := &domain.SyncResult{StartedAt: time.Now()}
    for _, e := range entries {
        w := convertIARMSEntry(e)
        _ = s.repo.UpsertFromIARMS(ctx, w)
        result.EntriesProcessed++
    }
    result.CompletedAt = time.Now()
    return result, nil
}
```

---

## 4. API REST

| Méthode | Endpoint                              | Rôle            | Description                      |
|---------|---------------------------------------|-----------------|----------------------------------|
| `POST`  | `/api/v1/biar/weapons`                | PNH, DOUANES    | Déclarer arme illicite saisie    |
| `GET`   | `/api/v1/biar/weapons/:id`            | DCPJ, PNH       | Détail arme illicite             |
| `GET`   | `/api/v1/biar/check/serial/:sn`       | PNH_OFFICER     | Vérifier numéro de série         |
| `POST`  | `/api/v1/biar/batches`                | DCPJ_SUPERVISOR | Créer saisie en lot              |
| `GET`   | `/api/v1/biar/stats/by-gang`          | DCPJ_INTEL      | Armes par gang                   |
| `GET`   | `/api/v1/biar/stats/by-origin`        | DCPJ, ATF       | Armes par pays d'origine         |
| `GET`   | `/api/v1/biar/stats/routes`           | DCPJ, ATF       | Routes de trafic identifiées     |
| `POST`  | `/api/v1/biar/iarms/sync`             | SUPERADMIN      | Synchronisation iARMS manuelle   |

---

## 5. ANALYTIQUES CLICKHOUSE

```sql
CREATE TABLE biar_weapon_events (
    weapon_id UUID, recovery_date Date,
    weapon_type String, caliber String,
    origin_country FixedString(3),
    recovery_dept FixedString(2),
    gang_id Nullable(UUID)
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(recovery_date)
ORDER BY (recovery_dept, weapon_type, recovery_date);

-- Tendance mensuelle par type et departement
CREATE MATERIALIZED VIEW biar_monthly_stats
ENGINE = SummingMergeTree()
ORDER BY (month, recovery_dept, weapon_type)
AS SELECT toStartOfMonth(recovery_date) AS month,
    recovery_dept, weapon_type,
    count() AS seizure_count
FROM biar_weapon_events
GROUP BY month, recovery_dept, weapon_type;
```

---

## 6. VARIABLES D'ENVIRONNEMENT

```dotenv
BIAR_DB_HOST=localhost
BIAR_DB_NAME=snisid_biar
BIAR_IARMS_GATEWAY=https://i247-gateway.pnh.gov.ht/iarms
BIAR_ATF_ETRACE_URL=https://etrace.atf.gov/api
BIAR_ATF_API_KEY=<VAULT:biar/atf_api_key>
BIAR_CLICKHOUSE_ADDR=clickhouse:9000
BIAR_IARMS_SYNC_HOURS=12
BIAR_SERVICE_PORT=8103
```

---
*MP-30 — BIAR-HT — Armes Illicites — SNISID — République d'Haïti*
