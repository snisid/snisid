# MP-21 — SIPEP-HT
## Système d'Information Pénitentiaire d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL SEULEMENT
Module           : MP-21
Code SNISID      : SIPEP-HT
Version          : 1.0.0
Dépendances      : FIR-HT (MP-20), AFIS-HT (MP-19), FPR-HT (MP-17)
Normes           : UN Standard Minimum Rules (Nelson Mandela Rules), OEA/CIDH
Acteurs          : DAP (Direction Administration Pénitentiaire), Parquet, Tribunaux
```

---

## 1. CONTEXTE

Haïti compte environ 12,000 détenus pour une capacité de 2,500 places. Le Pénitencier
National de Port-au-Prince est à 600% de sa capacité. Les registres carcéraux sont
manuels ou inexistants. Des détenus sont emprisonnés sans dossier, d'autres libérés
par erreur. Ce module crée le premier registre pénitentiaire numérique national.

### Établissements pénitentiaires couverts

| Établissement               | Département | Capacité réelle | Type          |
|-----------------------------|-------------|-----------------|---------------|
| Pénitencier National P-au-P | Ouest       | 3,500+          | National      |
| Prison Civile Cap-Haïtien   | Nord        | 800+            | Départemental |
| Prison Civile Gonaïves      | Artibonite  | 400+            | Départemental |
| Prison Civile Les Cayes     | Sud         | 300+            | Départemental |
| CERMICOL (Mineurs)          | Ouest       | 100             | Spécialisé    |
| Établissement femmes (RESEK)| Ouest       | 150             | Spécialisé    |
| 9 autres prisons civiles    | Tout pays   | Variable        | Local         |

---

## 2. ARCHITECTURE

```
services/sipep-svc/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── detention.go
│   │   ├── transfer.go
│   │   ├── release.go
│   │   ├── health_event.go
│   │   └── enums.go
│   ├── repository/postgres/
│   │   ├── detention_repo.go
│   │   ├── inmate_repo.go
│   │   └── transfer_repo.go
│   ├── service/
│   │   ├── intake_service.go
│   │   ├── transfer_service.go
│   │   ├── release_service.go
│   │   └── overcrowding_service.go
│   └── api/rest/
│       ├── intake_handler.go
│       ├── inmate_handler.go
│       └── stats_handler.go
└── Dockerfile
```

---

## 3. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE sipep_detention_basis AS ENUM (
    'PREVENTIVE',        -- Détention provisoire (avant jugement)
    'SENTENCED',         -- Condamné purgeant peine
    'ADMINISTRATIVE',    -- Détention administrative (immigration, etc.)
    'CONTEMPT'           -- Outrage au tribunal
);

CREATE TYPE sipep_legal_status AS ENUM (
    'AWAITING_TRIAL',
    'ON_TRIAL',
    'SENTENCED',
    'APPEAL_PENDING',
    'CONDEMNED'
);

CREATE TYPE sipep_release_type AS ENUM (
    'SENTENCE_SERVED',
    'CONDITIONAL_RELEASE',
    'BAIL',
    'JUDICIAL_ORDER',
    'DEATH',
    'ESCAPE',
    'TRANSFER_OUT'
);

CREATE TABLE sipep_inmates (
    inmate_id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_inmate_id      VARCHAR(25) UNIQUE NOT NULL, -- Format: SIPEP-HT-AAAA-NNNNNN
    snisid_person_id        UUID NOT NULL,
    fir_record_id           UUID,
    afis_subject_id         UUID,
    current_facility        VARCHAR(100) NOT NULL,
    current_dept_code       CHAR(2),
    cell_block              VARCHAR(20),
    is_currently_detained   BOOLEAN DEFAULT TRUE,
    is_minor                BOOLEAN DEFAULT FALSE,
    is_female               BOOLEAN DEFAULT FALSE,
    has_special_needs       BOOLEAN DEFAULT FALSE,
    special_needs_notes     TEXT,
    intake_date             TIMESTAMPTZ NOT NULL,
    expected_release_date   TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sipep_detentions (
    detention_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inmate_id               UUID NOT NULL REFERENCES sipep_inmates(inmate_id),
    facility                VARCHAR(100) NOT NULL,
    detention_basis         sipep_detention_basis NOT NULL,
    legal_status            sipep_legal_status NOT NULL DEFAULT 'AWAITING_TRIAL',
    case_reference          VARCHAR(100),
    court_name              VARCHAR(150),
    arresting_authority     VARCHAR(100),
    warrant_number          VARCHAR(100),
    intake_date             TIMESTAMPTZ NOT NULL,
    intake_officer          UUID NOT NULL,
    sentence_duration_days  INTEGER,
    time_served_days        INTEGER GENERATED ALWAYS AS (
                                EXTRACT(DAY FROM NOW() - intake_date)::INTEGER
                            ) STORED,
    release_date            TIMESTAMPTZ,
    release_type            sipep_release_type,
    releasing_authority     VARCHAR(100),
    notes                   TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sipep_transfers (
    transfer_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inmate_id               UUID NOT NULL REFERENCES sipep_inmates(inmate_id),
    from_facility           VARCHAR(100) NOT NULL,
    to_facility             VARCHAR(100) NOT NULL,
    transfer_date           TIMESTAMPTZ NOT NULL,
    transfer_reason         TEXT,
    authorized_by           UUID NOT NULL,
    transport_unit          VARCHAR(50),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sipep_health_events (
    event_id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inmate_id               UUID NOT NULL REFERENCES sipep_inmates(inmate_id),
    event_type              VARCHAR(50) NOT NULL, -- INJURY, ILLNESS, DEATH, PSYCHIATRIC
    event_date              TIMESTAMPTZ NOT NULL,
    description             TEXT,
    treating_facility       VARCHAR(150),
    outcome                 TEXT,
    reported_by             UUID,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Alerte surpopulation: vue matérialisée actualisée toutes les heures
CREATE MATERIALIZED VIEW sipep_facility_occupancy AS
SELECT
    current_facility,
    COUNT(*) AS current_count,
    current_dept_code
FROM sipep_inmates
WHERE is_currently_detained = TRUE
GROUP BY current_facility, current_dept_code;

CREATE INDEX idx_sipep_inmates_facility ON sipep_inmates(current_facility) WHERE is_currently_detained = TRUE;
CREATE INDEX idx_sipep_inmates_person   ON sipep_inmates(snisid_person_id);
CREATE INDEX idx_sipep_detentions_case  ON sipep_detentions(case_reference);
CREATE INDEX idx_sipep_detentions_basis ON sipep_detentions(detention_basis, legal_status);

COMMIT;
```

---

## 4. SERVICE CLÉ — GESTION LIBÉRATION

```go
// services/sipep-svc/internal/service/release_service.go
package service

import (
    "context"
    "fmt"
    "time"
    "github.com/google/uuid"
    "github.com/snisid/sipep-svc/internal/domain"
)

func (s *ReleaseService) ProcessRelease(
    ctx context.Context,
    inmateID uuid.UUID,
    req domain.ReleaseRequest,
    authorizedBy uuid.UUID,
) (*domain.Detention, error) {
    inmate, err := s.repo.FindInmate(ctx, inmateID)
    if err != nil {
        return nil, fmt.Errorf("détenu introuvable: %w", err)
    }
    if !inmate.IsCurrentlyDetained {
        return nil, fmt.Errorf("détenu non actuellement incarcéré")
    }

    detention, err := s.repo.GetActiveDetention(ctx, inmateID)
    if err != nil {
        return nil, fmt.Errorf("dossier détention actif introuvable: %w", err)
    }

    now := time.Now()
    detention.ReleaseDate  = &now
    detention.ReleaseType  = req.ReleaseType
    detention.ReleasingAuthority = req.Authority

    if err := s.repo.UpdateDetention(ctx, detention); err != nil {
        return nil, fmt.Errorf("mise à jour détention: %w", err)
    }

    inmate.IsCurrentlyDetained = false
    _ = s.repo.UpdateInmate(ctx, inmate)

    // Notifier FIR-HT et FPR-HT via Kafka
    _ = s.kafka.Publish(ctx, "sipep.inmate.released", domain.InmateReleasedEvent{
        InmateID:      inmateID,
        PersonID:      inmate.SNISIDPersonID,
        FacilityCode:  inmate.CurrentFacility,
        ReleaseType:   req.ReleaseType,
        ReleasedAt:    now,
        AuthorizedBy:  authorizedBy,
    })

    // Si libération = ESCAPE → déclencher alerte FPR immédiatement
    if req.ReleaseType == domain.ReleaseTypeEscape {
        _ = s.kafka.Publish(ctx, "sipep.escape.alert", domain.EscapeAlertEvent{
            InmateID:   inmateID,
            PersonID:   inmate.SNISIDPersonID,
            FacilityCode: inmate.CurrentFacility,
            EscapedAt:  now,
        })
    }

    return detention, nil
}
```

---

## 5. API REST

| Méthode  | Endpoint                                   | Rôle requis        | Description                      |
|----------|--------------------------------------------|--------------------|----------------------------------|
| `POST`   | `/api/v1/sipep/intake`                     | DAP_OFFICER        | Enregistrer entrée détenu        |
| `GET`    | `/api/v1/sipep/inmates/:id`                | DAP, PARQUET       | Fiche détenu complète            |
| `GET`    | `/api/v1/sipep/inmates/search`             | DAP, DCPJ          | Recherche multi-critères         |
| `POST`   | `/api/v1/sipep/release`                    | DAP_SUPERVISOR     | Traiter une libération           |
| `POST`   | `/api/v1/sipep/transfers`                  | DAP_ADMIN          | Transfert entre établissements   |
| `POST`   | `/api/v1/sipep/health-events`              | DAP_MEDICAL        | Enregistrer événement médical    |
| `GET`    | `/api/v1/sipep/facilities/occupancy`       | DAP_ADMIN          | Taux d'occupation en temps réel |
| `GET`    | `/api/v1/sipep/alerts/overcrowding`        | DAP_ADMIN, MJSP    | Alertes surpopulation            |
| `GET`    | `/api/v1/sipep/stats/preventive-detention` | PARQUET, MJSP      | Stats détention provisoire       |

---

## 6. INTÉGRATIONS

- **FIR-HT** : Kafka consumer `sipep.inmate.released` → mise à jour casier
- **FPR-HT** : Alerte immédiate sur `sipep.escape.alert` → mandat actif
- **AFIS-HT** : Enrôlement biométrique à chaque entrée en détention
- **DIPE-HT** : Si détenu non identifié → création automatique missing person case

---

## 7. VARIABLES D'ENVIRONNEMENT

```dotenv
SIPEP_DB_HOST=localhost
SIPEP_DB_NAME=snisid_sipep
SIPEP_REDIS_ADDR=redis-master:6379
SIPEP_FIR_SERVICE_URL=http://fir-svc:8093
SIPEP_FPR_SERVICE_URL=http://fpr-svc:8085
SIPEP_KAFKA_BROKERS=kafka:9092
SIPEP_OVERCROWDING_THRESHOLD=1.5
SIPEP_SERVICE_PORT=8092
```

---
*MP-21 — SIPEP-HT — Système Pénitentiaire — SNISID — République d'Haïti*
