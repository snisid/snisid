# MP-23 — OPR-HT
## Ordonnances de Protection et Restrictions Judiciaires d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL SEULEMENT
Module           : MP-23 | Code : OPR-HT
Dépendances      : FIR-HT (MP-20), FPR-HT (MP-17), SIPEP-HT (MP-21), BLKL-HT (MP-36)
Normes           : Code de procédure pénale haïtien, Convention CEDAW, CIPD
Acteurs          : Tribunaux, Parquet, PNH, Victimes/Plaignants, IBESR
```

---

## 1. CONTEXTE

Haïti ne dispose d'aucun registre centralisé des ordonnances de protection.
Une ordonnance émise à Port-au-Prince n'est pas connue à Cap-Haïtien ou Gonaïves.
Ce vide juridico-informatique met en danger : les victimes de violence conjugale,
les témoins de crimes, les survivants de kidnapping ayant dénoncé des membres de gangs.

### Types d'ordonnances requises dans le contexte haïtien

| Type                      | Contexte typique                                          |
|---------------------------|-----------------------------------------------------------|
| No-contact / Stay-away    | Violence conjugale, harcèlement                           |
| Exclusion domicile        | Expulsion conjoint violent du logement familial           |
| Zone exclusion gang       | Membre de gang banni d'un quartier/commune post-arrestation|
| Protection témoin         | Survivant de kidnapping / témoin à charge dans procès gang|
| Restriction de sortie     | Prévenu attendant jugement avec risque de fuite           |

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE opr_order_type AS ENUM (
    'RESTRAINING_ORDER', 'NO_CONTACT', 'STAY_AWAY',
    'PROTECTIVE', 'WITNESS_PROTECTION', 'GANG_EXCLUSION_ZONE',
    'TRAVEL_RESTRICTION'
);

CREATE TYPE opr_status AS ENUM (
    'ACTIVE', 'EXPIRED', 'VIOLATED', 'DISMISSED', 'APPEALED'
);

CREATE TABLE opr_protection_orders (
    order_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number        VARCHAR(30) UNIQUE NOT NULL,     -- OPR-HT-AAAA-NNNNNN
    order_type          opr_order_type NOT NULL,
    status              opr_status NOT NULL DEFAULT 'ACTIVE',

    protected_person_id UUID NOT NULL,
    subject_person_id   UUID NOT NULL,
    subject_fir_id      UUID,

    exclusion_radius_m  INTEGER,
    exclusion_addresses TEXT[] DEFAULT '{}',
    no_contact_modes    TEXT[] DEFAULT '{}',
    geographic_ban_geojson JSONB,

    issuing_court       VARCHAR(150) NOT NULL,
    issuing_judge       VARCHAR(150),
    issue_date          TIMESTAMPTZ NOT NULL,
    expiry_date         TIMESTAMPTZ NOT NULL,
    is_renewable        BOOLEAN DEFAULT TRUE,

    violation_count     SMALLINT DEFAULT 0,
    last_violation_at   TIMESTAMPTZ,

    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE opr_violations (
    violation_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id            UUID NOT NULL REFERENCES opr_protection_orders(order_id),
    violation_date      TIMESTAMPTZ NOT NULL,
    violation_type      VARCHAR(100) NOT NULL,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    reported_by         UUID NOT NULL,
    arrest_made         BOOLEAN DEFAULT FALSE,
    arrest_case_ref     VARCHAR(100),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE opr_witness_protections (
    protection_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    protected_person_id UUID NOT NULL,
    threat_level        VARCHAR(20) NOT NULL,
    gang_id             UUID,
    alias_assigned      VARCHAR(150),
    assigned_unit       VARCHAR(50),
    is_active           BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_opr_subject   ON opr_protection_orders(subject_person_id) WHERE status = 'ACTIVE';
CREATE INDEX idx_opr_protected ON opr_protection_orders(protected_person_id) WHERE status = 'ACTIVE';
CREATE INDEX idx_opr_expiry    ON opr_protection_orders(expiry_date) WHERE status = 'ACTIVE';

COMMIT;
```

---

## 3. SERVICE GO CLÉ

```go
package service

import (
    "context"
    "github.com/google/uuid"
    "github.com/snisid/opr-svc/internal/domain"
)

func (s *OPRService) CheckSubject(
    ctx context.Context,
    personID uuid.UUID,
) (*domain.OPRCheckResult, error) {
    orders, err := s.repo.FindActiveBySubject(ctx, personID)
    if err != nil || len(orders) == 0 {
        return &domain.OPRCheckResult{HasActiveOrder: false}, nil
    }
    return &domain.OPRCheckResult{
        HasActiveOrder: true,
        Orders:         orders,
        HighestType:    s.getHighestSeverity(orders),
    }, nil
}

func (s *OPRService) RecordViolation(
    ctx context.Context,
    orderID uuid.UUID,
    req domain.ViolationRequest,
    reportedBy uuid.UUID,
) error {
    order, err := s.repo.FindByID(ctx, orderID)
    if err != nil {
        return err
    }
    order.ViolationCount++
    order.Status = domain.StatusViolated
    _ = s.repo.Update(ctx, order)

    _ = s.kafka.Publish(ctx, "opr.violation.reported", domain.ViolationEvent{
        OrderID:    orderID,
        PersonID:   order.SubjectPersonID,
        ViolType:   req.ViolationType,
        ReportedBy: reportedBy,
    })
    // Si violations repetees => signal mandat FPR
    if order.ViolationCount >= 3 {
        _ = s.kafka.Publish(ctx, "opr.warrant.request", domain.WarrantRequestEvent{
            PersonID: order.SubjectPersonID,
            Reason:   "OPR violations repetees >= 3",
        })
    }
    return nil
}
```

---

## 4. API REST

| Méthode | Endpoint                           | Rôle           | Description                      |
|---------|------------------------------------|----------------|----------------------------------|
| `POST`  | `/api/v1/opr/orders`               | TRIBUNAL       | Créer ordonnance                 |
| `GET`   | `/api/v1/opr/check/:person_id`     | PNH, SIFR      | Vérifier ordonnances d'un sujet  |
| `POST`  | `/api/v1/opr/violations`           | PNH_OFFICER    | Signaler violation               |
| `PATCH` | `/api/v1/opr/orders/:id/renew`     | TRIBUNAL       | Renouveler ordonnance            |
| `GET`   | `/api/v1/opr/expiring-soon`        | TRIBUNAL_ADMIN | Ordonnances expirant dans 30j    |
| `POST`  | `/api/v1/opr/witness-protection`   | DCPJ_ADMIN     | Créer protection témoin          |
| `GET`   | `/api/v1/opr/orders/by-gang/:id`   | DCPJ           | Ordonnances vs membres d'un gang |

---

## 5. INTÉGRATIONS

- **FPR-HT** : Violation ≥ 3 → émission mandat d'arrêt automatique via Kafka
- **BLKL-HT** : Sujets sous OPR active avec restriction voyage → liste noire frontière
- **SIGEO-HT** : Zones d'exclusion géographiques → couche cartographique
- **GANG-HT** : Ordonnances zones de gangs → `gang_id` lié au dossier OPR

---

## 6. VARIABLES D'ENVIRONNEMENT

```dotenv
OPR_DB_HOST=localhost
OPR_DB_NAME=snisid_opr
OPR_FPR_SERVICE_URL=http://fpr-svc:8085
OPR_BLKL_SERVICE_URL=http://blkl-svc:8110
OPR_EXPIRY_ALERT_DAYS=30
OPR_SERVICE_PORT=8096
```

---
*MP-23 — OPR-HT — Ordonnances de Protection — SNISID — République d'Haïti*
