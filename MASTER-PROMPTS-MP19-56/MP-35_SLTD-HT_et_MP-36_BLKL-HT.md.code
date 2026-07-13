# MP-35 — SLTD-HT
## Interface Nationale INTERPOL SLTD — Documents de Voyage Perdus et Volés
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-35 | Code : SLTD-HT
Dépendances      : SIFR-HT (MP-33), BLKL-HT (MP-36), SNISID Identité, AERO-HT (MP-37)
Normes           : INTERPOL SLTD (100+ millions de documents), OACI Doc 9303
Acteurs          : DGMN, Passeport Bureau, PNH POLIFRONT, ONI
```

---

## 1. CONTEXTE

La base INTERPOL SLTD contient plus de 100 millions de documents de voyage perdus ou
volés provenant de 196 pays. En Haïti, des milliers de passeports et cartes d'identité
sont perdus ou volés chaque année. Ce module crée le registre national et l'interface
avec INTERPOL SLTD pour vérification aux frontières.

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE sltd_doc_type AS ENUM (
    'PASSPORT','NATIONAL_ID','TRAVEL_DOCUMENT',
    'VISA','RESIDENCE_PERMIT','REFUGEE_DOCUMENT','LAISSEZ_PASSER'
);

CREATE TYPE sltd_doc_status AS ENUM (
    'LOST','STOLEN','REVOKED','EXPIRED','FOUND','RECOVERED','CANCELLED'
);

CREATE TABLE sltd_documents (
    doc_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_sltd_id    VARCHAR(25) UNIQUE NOT NULL,  -- SLTD-HT-NNNNNN
    doc_type            sltd_doc_type NOT NULL,
    document_number     VARCHAR(100) NOT NULL,
    issuing_country     CHAR(3) NOT NULL DEFAULT 'HTI',
    holder_name         VARCHAR(200),
    holder_snisid_id    UUID,
    holder_dob          DATE,
    holder_nationality  CHAR(3) DEFAULT 'HTI',
    issue_date          DATE,
    expiry_date         DATE,
    status              sltd_doc_status NOT NULL,
    reported_date       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reported_by         UUID NOT NULL,
    reporting_dept_code CHAR(2),
    theft_context       TEXT,
    found_date          TIMESTAMPTZ,
    found_location      VARCHAR(300),
    interpol_sltd_ref   VARCHAR(50),
    reported_to_interpol BOOLEAN DEFAULT FALSE,
    interpol_reported_at TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sltd_check_log (
    check_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_number     VARCHAR(100) NOT NULL,
    doc_type            sltd_doc_type,
    checked_by          UUID NOT NULL,
    check_location      VARCHAR(100),
    post_id             UUID,
    result              VARCHAR(20) NOT NULL,    -- CLEAR, LOST, STOLEN, REVOKED
    source              VARCHAR(20) NOT NULL,    -- LOCAL, INTERPOL_SLTD, BOTH
    sltd_doc_id         UUID REFERENCES sltd_documents(doc_id),
    checked_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_sltd_doc_number ON sltd_documents(document_number, issuing_country)
    WHERE status IN ('LOST','STOLEN','REVOKED');
CREATE INDEX idx_sltd_holder    ON sltd_documents(holder_snisid_id) WHERE holder_snisid_id IS NOT NULL;
CREATE INDEX idx_sltd_status    ON sltd_documents(status);
CREATE INDEX idx_sltd_check_log ON sltd_check_log(document_number, checked_at DESC);

COMMIT;
```

---

## 3. SERVICE GO CLÉ — VÉRIFICATION INSTANTANÉE

```go
package service

import "context"

// CheckDocument verifie un document en < 100ms
func (s *SLTDService) CheckDocument(
    ctx context.Context,
    docNumber, issuingCountry string,
) (*CheckResult, error) {
    // 1. Cache Redis hotlist (< 5ms)
    if cached, _ := s.cache.Get(ctx, docNumber); cached != nil {
        return cached, nil
    }
    // 2. Base locale SLTD-HT
    doc, _ := s.repo.FindByNumber(ctx, docNumber, issuingCountry)
    if doc != nil {
        r := &CheckResult{
            DocumentNumber: docNumber,
            IsStolen: doc.Status == "STOLEN",
            IsLost:   doc.Status == "LOST",
            Status:   doc.Status,
            Source:   "LOCAL",
        }
        _ = s.cache.Set(ctx, docNumber, r, 15*60)
        return r, nil
    }
    // 3. Fallback INTERPOL SLTD (< 200ms via I-24/7)
    interpol, _ := s.interpolClient.CheckSLTD(ctx, docNumber, issuingCountry)
    if interpol != nil {
        return interpol, nil
    }
    return &CheckResult{DocumentNumber: docNumber, Status: "CLEAR"}, nil
}
```

---

## 4. API REST

| Méthode | Endpoint                             | Rôle        | Description                       |
|---------|--------------------------------------|-------------|-----------------------------------|
| `GET`   | `/api/v1/sltd/check/:num`            | SIFR, AERO  | Vérifier un document              |
| `POST`  | `/api/v1/sltd/report/lost`           | PUBLIC, PNH | Déclarer document perdu           |
| `POST`  | `/api/v1/sltd/report/stolen`         | PNH_OFFICER | Déclarer document volé            |
| `PATCH` | `/api/v1/sltd/:id/found`             | PNH_OFFICER | Marquer document retrouvé         |
| `GET`   | `/api/v1/sltd/stats`                 | DGMN_ADMIN  | Statistiques par type/département |

---

## 5. VARIABLES D'ENVIRONNEMENT

```dotenv
SLTD_DB_HOST=localhost
SLTD_DB_NAME=snisid_sltd
SLTD_REDIS_ADDR=redis-master:6379
SLTD_HOTLIST_TTL_MINUTES=15
SLTD_INTERPOL_SLTD_URL=https://i247-gateway.pnh.gov.ht/sltd
SLTD_SERVICE_PORT=8108
```

---
*MP-35 — SLTD-HT — Documents Perdus/Volés — SNISID — République d'Haïti*

---
---

# MP-36 — BLKL-HT
## Liste Noire Nationale — Personnes Interdites d'Entrée ou de Sortie du Territoire
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-36 | Code : BLKL-HT
Dépendances      : SIFR-HT (MP-33), FPR-HT (MP-17), OPR-HT (MP-23), SANC-HT (MP-27)
Normes           : Loi haïtienne sur l'immigration, Résolution CSNU 2653 (interdictions voyage)
Acteurs          : DGMN, MJSP, Parquet, Tribunaux
```

---

## 1. CONTEXTE

La liste noire nationale centralise toutes les personnes interdites de sortie ou d'entrée
sur le territoire haïtien. Sources : mandats d'arrêt actifs, sanctions ONU/OFAC avec
interdiction de voyage, ordonnances judiciaires restreignant les déplacements, personnes
expulsées avec interdiction de réentrée, membres de gangs sous enquête.

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE blkl_restriction_type AS ENUM (
    'ENTRY_BAN',            -- Interdit d entre sur le territoire
    'EXIT_BAN',             -- Interdit de sortir (suspect sous enquete)
    'BOTH_BAN',             -- Double interdiction (tres graves)
    'CONDITIONAL_BAN'       -- Avec conditions (caution, rapport regulier)
);

CREATE TYPE blkl_source AS ENUM (
    'JUDICIAL_ORDER', 'WANTED_WARRANT', 'UN_SANCTIONS',
    'OFAC_SANCTIONS', 'MINISTERIAL_ORDER', 'EXPULSION',
    'OPR_TRAVEL_RESTRICTION', 'INTERPOL_NOTICE'
);

CREATE TABLE blkl_blacklist (
    entry_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_blkl_id    VARCHAR(25) UNIQUE NOT NULL,  -- BLKL-HT-NNNNNN
    snisid_person_id    UUID NOT NULL,
    restriction_type    blkl_restriction_type NOT NULL,
    source              blkl_source NOT NULL,
    source_record_id    UUID,
    reason              TEXT NOT NULL,
    court_order_ref     VARCHAR(100),
    ordered_by          VARCHAR(150),
    effective_date      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expiry_date         TIMESTAMPTZ,
    is_permanent        BOOLEAN DEFAULT FALSE,
    is_active           BOOLEAN DEFAULT TRUE,
    alert_level         VARCHAR(20) DEFAULT 'WANTED',
    armed_dangerous     BOOLEAN DEFAULT FALSE,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE blkl_alerts_log (
    alert_log_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id            UUID NOT NULL REFERENCES blkl_blacklist(entry_id),
    triggered_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    post_code           VARCHAR(10),
    direction           VARCHAR(10),
    action_taken        TEXT,
    officer_id          UUID,
    outcome             TEXT
);

CREATE INDEX idx_blkl_person   ON blkl_blacklist(snisid_person_id) WHERE is_active = TRUE;
CREATE INDEX idx_blkl_type     ON blkl_blacklist(restriction_type) WHERE is_active = TRUE;
CREATE INDEX idx_blkl_expiry   ON blkl_blacklist(expiry_date) WHERE is_active = TRUE AND expiry_date IS NOT NULL;

-- Trigger: expirer automatiquement les entrees
CREATE OR REPLACE FUNCTION blkl_auto_expire()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.expiry_date IS NOT NULL AND NEW.expiry_date < NOW() AND NEW.is_active = TRUE THEN
        NEW.is_active := FALSE;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_blkl_expire
    BEFORE UPDATE ON blkl_blacklist
    FOR EACH ROW EXECUTE FUNCTION blkl_auto_expire();

COMMIT;
```

---

## 3. SERVICE GO CLÉ

```go
package service

import (
    "context"
    "github.com/google/uuid"
    "github.com/snisid/blkl-svc/internal/domain"
)

func (s *BlacklistService) CheckPerson(
    ctx context.Context,
    personID uuid.UUID,
) (*domain.BlacklistCheckResult, error) {
    // Cache Redis ultra-rapide (hotlist frontiere)
    if cached, _ := s.cache.Get(ctx, personID.String()); cached != nil {
        return cached, nil
    }
    entries, err := s.repo.FindActiveByPerson(ctx, personID)
    if err != nil || len(entries) == 0 {
        return &domain.BlacklistCheckResult{IsBlacklisted: false}, nil
    }
    result := &domain.BlacklistCheckResult{
        IsBlacklisted:   true,
        PersonID:        personID,
        Restrictions:    entries,
        ArmedDangerous:  false,
    }
    for _, e := range entries {
        if e.ArmedDangerous {
            result.ArmedDangerous = true
        }
    }
    _ = s.cache.Set(ctx, personID.String(), result, 900) // 15 min TTL
    return result, nil
}
```

---

## 4. API REST

| Méthode | Endpoint                          | Rôle            | Description                        |
|---------|-----------------------------------|-----------------|------------------------------------|
| `GET`   | `/api/v1/blkl/check/:person_id`   | SIFR, AERO      | Vérifier personne en temps réel    |
| `POST`  | `/api/v1/blkl/entries`            | MJSP, TRIBUNAL  | Ajouter à la liste noire           |
| `PATCH` | `/api/v1/blkl/entries/:id/lift`   | MJSP_ADMIN      | Lever l'interdiction               |
| `GET`   | `/api/v1/blkl/entries/active`     | DGMN, DCPJ      | Liste complète active              |
| `GET`   | `/api/v1/blkl/expiring-soon`      | DGMN_ADMIN      | Entrées expirant dans 30 jours     |

---

## 5. VARIABLES D'ENVIRONNEMENT

```dotenv
BLKL_DB_HOST=localhost
BLKL_DB_NAME=snisid_blkl
BLKL_REDIS_ADDR=redis-master:6379
BLKL_HOTLIST_TTL_SECONDS=900
BLKL_FPR_SERVICE_URL=http://fpr-svc:8085
BLKL_SANC_SERVICE_URL=http://sanc-svc:8100
BLKL_SERVICE_PORT=8110
```

---
*MP-36 — BLKL-HT — Liste Noire Territoriale — SNISID — République d'Haïti*
