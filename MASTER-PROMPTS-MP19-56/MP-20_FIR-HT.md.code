# MP-20 — FIR-HT
## Fichier Individuel des Renseignements — Casier Judiciaire National d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL SEULEMENT
Module           : MP-20
Code SNISID      : FIR-HT
Version          : 1.0.0
Dépendances      : AFIS-HT (MP-19), FPR-HT (MP-17), SIPEP-HT (MP-21), RDEP-HT (MP-22)
Normes           : INTERPOL CCC, Loi haïtienne sur le casier judiciaire (1974 révisée)
Acteurs          : Parquet, DCPJ, Tribunaux, BRI, Douanes, Employeurs publics
```

---

## 1. CONTEXTE

Haïti ne dispose d'aucun casier judiciaire centralisé et informatisé. Les tribunaux
gèrent des dossiers papier par ressort territorial. Les récidivistes changent de commune
pour effacer leur historique. Ce module constitue le registre national consolidé de
toutes les condamnations, arrestations et décisions judiciaires.

### Entités alimentant FIR-HT

| Source                | Données apportées                                    |
|-----------------------|------------------------------------------------------|
| DCPJ                  | Arrestations, gardes à vue, inculpations             |
| Tribunaux de Paix     | Jugements contraventionnels                          |
| Tribunaux Correctionnels | Délits, peines correctionnelles                   |
| Cours d'Assises       | Crimes, réclusion criminelle                         |
| Cours d'Appel         | Décisions d'appel, relaxes                           |
| SIPEP-HT (MP-21)      | Entrées/sorties pénitentiaires, libérations          |
| RDEP-HT (MP-22)       | Condamnations étrangères des déportés                |
| INTERPOL CCC          | Casiers d'autres pays via I-24/7                     |

---

## 2. ARCHITECTURE

```
services/fir-svc/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── criminal_record.go
│   │   ├── conviction.go
│   │   ├── arrest.go
│   │   ├── judicial_decision.go
│   │   └── enums.go
│   ├── repository/postgres/
│   │   ├── criminal_record_repo.go
│   │   ├── conviction_repo.go
│   │   └── arrest_repo.go
│   ├── service/
│   │   ├── record_service.go
│   │   ├── certificate_service.go    ← Extrait de casier judiciaire
│   │   └── expungement_service.go    ← Réhabilitation judiciaire
│   └── api/rest/
│       ├── record_handler.go
│       ├── certificate_handler.go
│       └── search_handler.go
└── Dockerfile
```

---

## 3. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE fir_offense_class AS ENUM (
    'CONTRAVENTION',    -- Infraction mineure
    'DELIT',            -- Infraction correctionnelle
    'CRIME',            -- Crime grave
    'FELONY_FOREIGN'    -- Infraction étrangère (déportés)
);

CREATE TYPE fir_case_status AS ENUM (
    'OPEN','PENDING_TRIAL','CONVICTED','ACQUITTED',
    'DISMISSED','APPEAL_PENDING','EXPUNGED'
);

CREATE TYPE fir_sentence_type AS ENUM (
    'PRISON','SUSPENDED','FINE','COMMUNITY_SERVICE',
    'DEATH_PENALTY','ACQUITTAL','PROBATION'
);

CREATE TABLE fir_criminal_records (
    record_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_fir_id     VARCHAR(25) UNIQUE NOT NULL,   -- Format: FIR-HT-AAAA-NNNNNN
    snisid_person_id    UUID NOT NULL,                  -- Lien identité SNISID
    afis_subject_id     UUID,                           -- Lien AFIS-HT
    is_haitian_national BOOLEAN DEFAULT TRUE,
    aliases             TEXT[] DEFAULT '{}',
    is_active           BOOLEAN DEFAULT TRUE,
    is_expunged         BOOLEAN DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fir_arrests (
    arrest_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id           UUID NOT NULL REFERENCES fir_criminal_records(record_id),
    arrest_date         TIMESTAMPTZ NOT NULL,
    arresting_unit      VARCHAR(50) NOT NULL,
    arresting_officer   UUID,
    arrest_location     VARCHAR(300),
    dept_code           CHAR(2),
    charges_text        TEXT NOT NULL,
    offense_class       fir_offense_class NOT NULL,
    case_reference      VARCHAR(100),
    release_date        TIMESTAMPTZ,
    release_reason      TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fir_convictions (
    conviction_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id           UUID NOT NULL REFERENCES fir_criminal_records(record_id),
    case_reference      VARCHAR(100) NOT NULL,
    court_name          VARCHAR(150) NOT NULL,
    court_dept          CHAR(2),
    offense_class       fir_offense_class NOT NULL,
    offense_description TEXT NOT NULL,
    ipc_code            VARCHAR(30),                    -- Code pénal haïtien
    verdict_date        TIMESTAMPTZ NOT NULL,
    case_status         fir_case_status NOT NULL,
    sentence_type       fir_sentence_type,
    sentence_duration_days INTEGER,
    fine_amount_gdes    DECIMAL(12,2),
    sentence_start      TIMESTAMPTZ,
    sentence_end        TIMESTAMPTZ,
    is_foreign_record   BOOLEAN DEFAULT FALSE,
    foreign_country     CHAR(3),                        -- HTI, USA, DOM, etc.
    interpol_ccc_ref    VARCHAR(50),
    judge_name          VARCHAR(150),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE fir_certificates (
    cert_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id           UUID REFERENCES fir_criminal_records(record_id),
    snisid_person_id    UUID NOT NULL,
    certificate_number  VARCHAR(30) UNIQUE NOT NULL,
    issued_for          VARCHAR(200),                   -- Motif: emploi, visa, etc.
    result              VARCHAR(20) NOT NULL,            -- CLEAN, HAS_RECORD
    issued_by           UUID NOT NULL,
    issuing_office      VARCHAR(100),
    issued_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at          TIMESTAMPTZ,
    qr_code_ref         VARCHAR(200)                    -- Vérification authenticité
);

CREATE INDEX idx_fir_records_snisid   ON fir_criminal_records(snisid_person_id);
CREATE INDEX idx_fir_records_fir_id   ON fir_criminal_records(national_fir_id);
CREATE INDEX idx_fir_arrests_date     ON fir_arrests(arrest_date DESC);
CREATE INDEX idx_fir_convictions_date ON fir_convictions(verdict_date DESC);
CREATE INDEX idx_fir_convictions_dept ON fir_convictions(court_dept);

ALTER TABLE fir_criminal_records ENABLE ROW LEVEL SECURITY;
CREATE POLICY fir_read_policy ON fir_criminal_records
    FOR SELECT USING (
        current_setting('app.user_role', TRUE) IN
        ('DCPJ','PARQUET','TRIBUNAL','JUDGE','POLICE_OFFICER','SUPERADMIN')
    );

COMMIT;
```

---

## 4. SERVICE GO CLÉS

```go
// services/fir-svc/internal/service/record_service.go
package service

import (
    "context"
    "fmt"
    "time"
    "github.com/google/uuid"
    "github.com/snisid/fir-svc/internal/domain"
)

type RecordService struct {
    repo     domain.CriminalRecordRepository
    kafka    domain.EventPublisher
    snisid   domain.SNISIDClient
}

// GetOrCreateRecord retourne le casier existant ou en crée un nouveau
func (s *RecordService) GetOrCreateRecord(
    ctx context.Context,
    snisidPersonID uuid.UUID,
) (*domain.CriminalRecord, error) {
    existing, err := s.repo.FindByPersonID(ctx, snisidPersonID)
    if err == nil && existing != nil {
        return existing, nil
    }
    person, err := s.snisid.GetPerson(ctx, snisidPersonID)
    if err != nil {
        return nil, fmt.Errorf("personne SNISID introuvable: %w", err)
    }
    record := &domain.CriminalRecord{
        RecordID:        uuid.New(),
        NationalFIRID:   s.generateFIRID(),
        SNISIDPersonID:  snisidPersonID,
        IsHaitianNational: person.Nationality == "HTI",
        CreatedAt:       time.Now(),
    }
    if err := s.repo.Create(ctx, record); err != nil {
        return nil, fmt.Errorf("création casier: %w", err)
    }
    _ = s.kafka.Publish(ctx, "fir.record.created", record)
    return record, nil
}

func (s *RecordService) generateFIRID() string {
    return fmt.Sprintf("FIR-HT-%d-%06d", time.Now().Year(),
        s.repo.NextSequence())
}

// IssueCertificate génère un extrait de casier judiciaire officiel
func (s *RecordService) IssueCertificate(
    ctx context.Context,
    req domain.CertificateRequest,
    issuedBy uuid.UUID,
) (*domain.Certificate, error) {
    record, _ := s.repo.FindByPersonID(ctx, req.PersonID)
    result := "CLEAN"
    if record != nil && record.HasActiveConvictions() {
        result = "HAS_RECORD"
    }
    cert := &domain.Certificate{
        CertID:            uuid.New(),
        RecordID:          func() *uuid.UUID { if record != nil { return &record.RecordID }; return nil }(),
        SNISIDPersonID:    req.PersonID,
        CertificateNumber: s.generateCertNumber(),
        IssuedFor:         req.Purpose,
        Result:            result,
        IssuedBy:          issuedBy,
        IssuingOffice:     req.Office,
        IssuedAt:          time.Now(),
        ExpiresAt:         func() *time.Time { t := time.Now().AddDate(0, 3, 0); return &t }(),
    }
    _ = s.repo.SaveCertificate(ctx, cert)
    return cert, nil
}
```

---

## 5. API REST

| Méthode | Endpoint                                 | Rôle requis              | Description                        |
|---------|------------------------------------------|--------------------------|------------------------------------|
| `POST`  | `/api/v1/fir/records`                    | DCPJ, PARQUET            | Créer fiche casier                 |
| `GET`   | `/api/v1/fir/records/:person_id`         | PNH, PARQUET, TRIBUNAL   | Consulter casier par personne      |
| `POST`  | `/api/v1/fir/records/:id/arrests`        | PNH_OFFICER, DCPJ        | Ajouter arrestation                |
| `POST`  | `/api/v1/fir/records/:id/convictions`    | TRIBUNAL, PARQUET        | Enregistrer condamnation           |
| `GET`   | `/api/v1/fir/certificates/issue`         | BUREAU_IDENTIFICATION    | Émettre extrait de casier          |
| `GET`   | `/api/v1/fir/certificates/verify/:num`   | PUBLIC                   | Vérifier authenticité extrait      |
| `POST`  | `/api/v1/fir/expunge`                    | TRIBUNAL_ADMIN           | Réhabilitation judiciaire          |
| `GET`   | `/api/v1/fir/search`                     | DCPJ, PARQUET            | Recherche multi-critères           |

---

## 6. INTÉGRATIONS

- **AFIS-HT (MP-19)** — à la création d'un casier, lier automatiquement avec le profil AFIS
- **SIPEP-HT (MP-21)** — synchronisation bidirectionnelle entrées/sorties pénitentiaires
- **RDEP-HT (MP-22)** — import automatique des antécédents USA/Canada des déportés
- **FPR-HT (MP-17)** — personnes avec mandats actifs → flag dans casier
- **INTERPOL CCC** — Circulation Criminal Charges: export/import via I-24/7

---

## 7. SÉCURITÉ & CONFORMITÉ

- Les données d'un casier sont à usage strictement officiel — violation = crime pénal
- Accès civil (employeurs, universités) uniquement via certificat officiel signé numériquement
- Droit de rectification : toute erreur peut être contestée via procédure tribunal
- Conservation : 10 ans après purge de peine, 30 ans pour crimes graves
- QR code sur chaque extrait → vérification authenticité publique (sans données sensibles)

---

## 8. VARIABLES D'ENVIRONNEMENT

```dotenv
FIR_DB_HOST=localhost
FIR_DB_NAME=snisid_fir
FIR_REDIS_ADDR=redis-master:6379
FIR_AFIS_SERVICE_URL=http://afis-svc:8091
FIR_SIPEP_SERVICE_URL=http://sipep-svc:8092
FIR_SNISID_IDENTITY_URL=http://identity-svc:8080
FIR_CERT_VALIDITY_MONTHS=3
FIR_SERVICE_PORT=8093
FIR_QR_SECRET=<VAULT:fir/qr_signing_secret>
```

---
*MP-20 — FIR-HT — Casier Judiciaire National — SNISID — République d'Haïti*
