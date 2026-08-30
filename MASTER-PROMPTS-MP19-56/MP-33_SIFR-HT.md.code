# MP-33 — SIFR-HT
## Système d'Information Frontières et Routes Terrestres d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-33 | Code : SIFR-HT
Dépendances      : BLKL-HT (MP-36), SLTD-HT (MP-35), FPR-HT (MP-17), RDEP-HT (MP-22)
Normes           : INTERPOL PISCES, OIM Standards Border Management, OACI Doc 9303
Acteurs          : DGMN (Direction Générale Migration et Naturalisation), PNH POLIFRONT
```

---

## 1. CONTEXTE

Haïti partage 360 km de frontière terrestre avec la République Dominicaine. Les postes
officiels principaux (Malpasse/Jimaní, Ouanaminthe/Dajabon, Belladère/Comendador,
Anse-à-Pitres/Pedernales) coexistent avec des dizaines de passages clandestins utilisés
par les trafiquants et les gangs. Ce module crée le premier système de gestion
intégrée des frontières.

### Postes frontaliers officiels

| Poste             | Département | Statut        | Volume journalier |
|-------------------|-------------|---------------|-------------------|
| Malpasse / Jimaní | Ouest       | Principal     | 3,000-5,000       |
| Ouanaminthe / Dajabón | Nord-Est | Principal     | 2,000-4,000       |
| Belladère / Comendador | Centre  | Secondaire   | 500-1,000         |
| Anse-à-Pitres / Pedernales | Sud-Est | Secondaire | 200-500          |
| Fonds-Verrettes   | Ouest       | Informel      | Non contrôlé      |

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE sifr_crossing_direction AS ENUM ('ENTRY', 'EXIT');
CREATE TYPE sifr_doc_type AS ENUM (
    'PASSPORT', 'NATIONAL_ID', 'LAISSEZ_PASSER',
    'BIRTH_CERTIFICATE', 'TRAVEL_DOCUMENT', 'NONE'
);
CREATE TYPE sifr_alert_type AS ENUM (
    'WANTED_PERSON', 'STOLEN_DOCUMENT', 'BLACKLIST',
    'ACTIVE_WARRANT', 'SANCTIONS', 'CUSTOMS_ALERT'
);

CREATE TABLE sifr_border_posts (
    post_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_code           VARCHAR(10) UNIQUE NOT NULL,  -- MLP, OUN, BLD, AAP, etc.
    name                VARCHAR(150) NOT NULL,
    dept_code           CHAR(2) NOT NULL,
    border_country      CHAR(3) NOT NULL DEFAULT 'DOM',
    post_lat            DECIMAL(10,7),
    post_lng            DECIMAL(10,7),
    is_official         BOOLEAN DEFAULT TRUE,
    is_active           BOOLEAN DEFAULT TRUE,
    lanes_count         SMALLINT DEFAULT 2,
    has_biometric_scanner BOOLEAN DEFAULT FALSE,
    has_vehicle_scanner BOOLEAN DEFAULT FALSE,
    operating_hours     VARCHAR(50),
    commanding_officer  UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sifr_crossings (
    crossing_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id             UUID NOT NULL REFERENCES sifr_border_posts(post_id),
    direction           sifr_crossing_direction NOT NULL,
    crossing_datetime   TIMESTAMPTZ NOT NULL,
    snisid_person_id    UUID,
    document_type       sifr_doc_type NOT NULL DEFAULT 'PASSPORT',
    document_number     VARCHAR(100),
    document_country    CHAR(3),
    document_expiry     DATE,
    traveler_name       VARCHAR(200) NOT NULL,
    traveler_dob        DATE,
    traveler_nationality CHAR(3),
    vehicle_plate       VARCHAR(20),
    lane_number         SMALLINT,
    processing_officer  UUID NOT NULL,
    alert_triggered     BOOLEAN DEFAULT FALSE,
    alert_type          sifr_alert_type,
    alert_action_taken  TEXT,
    processing_time_sec INTEGER,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sifr_alerts_log (
    alert_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crossing_id         UUID REFERENCES sifr_crossings(crossing_id),
    post_id             UUID NOT NULL,
    alert_type          sifr_alert_type NOT NULL,
    snisid_person_id    UUID,
    document_number     VARCHAR(100),
    vehicle_plate       VARCHAR(20),
    alert_source        VARCHAR(50),     -- FPR, BLKL, SLTD, OPR, SANC
    source_record_id    UUID,
    notified_units      TEXT[] DEFAULT '{}',
    action_taken        TEXT,
    resolved            BOOLEAN DEFAULT FALSE,
    resolved_by         UUID,
    resolved_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sifr_clandestine_crossings (
    report_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    reported_date       TIMESTAMPTZ NOT NULL,
    crossing_type       VARCHAR(50),     -- FOOT, VEHICLE, BOAT
    estimated_persons   INTEGER,
    gang_related        BOOLEAN DEFAULT FALSE,
    gang_id             UUID,
    trafficking_type    TEXT,
    reported_by         UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sifr_crossings_datetime ON sifr_crossings(crossing_datetime DESC);
CREATE INDEX idx_sifr_crossings_person   ON sifr_crossings(snisid_person_id) WHERE snisid_person_id IS NOT NULL;
CREATE INDEX idx_sifr_crossings_doc      ON sifr_crossings(document_number);
CREATE INDEX idx_sifr_crossings_alert    ON sifr_crossings(alert_triggered) WHERE alert_triggered = TRUE;
CREATE INDEX idx_sifr_crossings_post     ON sifr_crossings(post_id, crossing_datetime DESC);

COMMIT;
```

---

## 3. SERVICE GO CLÉ — VÉRIFICATION EN TEMPS RÉEL

```go
package service

import (
    "context"
    "time"
    "github.com/snisid/sifr-svc/internal/domain"
)

// ProcessCrossing verifie et enregistre un passage frontiere
func (s *BorderService) ProcessCrossing(
    ctx context.Context,
    req domain.CrossingRequest,
    officerID string,
) (*domain.CrossingResult, error) {
    result := &domain.CrossingResult{
        CrossingID: generateID(),
        ProcessedAt: time.Now(),
        Clearance:   domain.ClearanceGranted,
    }

    // 1. Verif document voyage (SLTD-HT) - sous 200ms
    docCheck, _ := s.sltdClient.CheckDocument(ctx, req.DocumentNumber, req.DocumentCountry)
    if docCheck != nil && docCheck.IsStolen {
        result.Clearance = domain.ClearanceDenied
        result.AlertType = domain.AlertStolenDocument
        result.AlertSource = "SLTD-HT"
    }

    // 2. Verif liste noire (BLKL-HT)
    if req.SNISIDPersonID != "" {
        blkl, _ := s.blklClient.CheckPerson(ctx, req.SNISIDPersonID)
        if blkl != nil && blkl.IsBlacklisted {
            result.Clearance = domain.ClearanceDenied
            result.AlertType = domain.AlertBlacklist
        }
    }

    // 3. Verif mandat actif (FPR-HT)
    warrant, _ := s.fprClient.CheckWarrant(ctx, req.DocumentNumber)
    if warrant != nil && warrant.IsActive {
        result.Clearance = domain.ClearanceDenied
        result.AlertType = domain.AlertWanted
        result.IsDangerous = warrant.ArmedAndDangerous
    }

    // 4. Verif ordonnance de protection voyage (OPR-HT)
    opr, _ := s.oprClient.CheckTravelRestriction(ctx, req.SNISIDPersonID)
    if opr != nil && opr.HasTravelRestriction {
        result.Clearance = domain.ClearanceDenied
    }

    // 5. Verif sanctions internationales (SANC-HT)
    sanc, _ := s.sancClient.CheckPerson(ctx, req.SNISIDPersonID)
    if sanc != nil && sanc.IsSanctioned {
        result.Clearance = domain.ClearanceDenied
        result.AlertType = domain.AlertSanctions
    }

    // Enregistrement passage (autorise ou refuse)
    _ = s.repo.CreateCrossing(ctx, domain.Crossing{
        PostID:          req.PostID,
        Direction:       req.Direction,
        CrossingDatetime: result.ProcessedAt,
        AlertTriggered:  result.Clearance == domain.ClearanceDenied,
        AlertType:       result.AlertType,
    })

    return result, nil
}
```

---

## 4. API REST

| Méthode | Endpoint                             | Rôle            | Description                      |
|---------|--------------------------------------|-----------------|----------------------------------|
| `POST`  | `/api/v1/sifr/crossings`             | SIFR_AGENT      | Enregistrer passage frontière    |
| `GET`   | `/api/v1/sifr/crossings/search`      | DGMN, DCPJ      | Recherche historique passages    |
| `GET`   | `/api/v1/sifr/crossings/person/:id`  | DCPJ            | Historique passages d'une personne|
| `GET`   | `/api/v1/sifr/alerts/active`         | SIFR_SUPERVISOR | Alertes actives                  |
| `GET`   | `/api/v1/sifr/posts`                 | PUBLIC_DGMN     | Liste postes frontaliers         |
| `GET`   | `/api/v1/sifr/stats/daily`           | DGMN_ADMIN      | Stats journalières par poste     |
| `POST`  | `/api/v1/sifr/clandestine`           | PNH_OFFICER     | Signaler passage clandestin      |

---

## 5. VARIABLES D'ENVIRONNEMENT

```dotenv
SIFR_DB_HOST=localhost
SIFR_DB_NAME=snisid_sifr
SIFR_REDIS_ADDR=redis-master:6379
SIFR_SLTD_SERVICE_URL=http://sltd-svc:8108
SIFR_BLKL_SERVICE_URL=http://blkl-svc:8110
SIFR_FPR_SERVICE_URL=http://fpr-svc:8085
SIFR_SANC_SERVICE_URL=http://sanc-svc:8100
SIFR_OPR_SERVICE_URL=http://opr-svc:8096
SIFR_PROCESSING_TIMEOUT_MS=500
SIFR_SERVICE_PORT=8106
```

---
*MP-33 — SIFR-HT — Frontières Terrestres — SNISID — République d'Haïti*
