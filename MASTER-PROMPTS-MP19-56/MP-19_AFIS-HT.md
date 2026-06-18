# MP-19 — AFIS-HT
## Automated Fingerprint Identification System d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL SEULEMENT
Module           : MP-19
Code SNISID      : AFIS-HT
Version          : 1.0.0
Dépendances      : FIR-HT (MP-20), FPR-HT (MP-17), SNISID-BIO-ADN
Normes           : NIST NFIQ2, ISO/IEC 19794-2, INTERPOL AFIS Standards, FBI IAFIS
Acteurs          : DCPJ, Labo Forensique PNH, Parquet, Bureaux Identification
```

---

## 1. CONTEXTE ET VISION STRATÉGIQUE

### 1.1 Problématique haïtienne

En Haïti, l'identification criminelle repose encore sur des fiches cartonnées et des
carnets manuscrits dans les commissariats. En l'absence d'un système AFIS national :
- Les récidivistes utilisent des alias multiples sans être détectés
- Les déportés avec casiers américains ne peuvent être liés à leurs antécédents
- Les scènes de crime livrent des empreintes latentes qui ne sont jamais exploitées
- L'identification des cadavres non identifiés (RVIN-HT) est impossible sans AFIS

### 1.2 Portée fonctionnelle

| Capacité                    | Description                                              |
|-----------------------------|----------------------------------------------------------|
| Enrôlement 10 doigts        | Capture livescanner + cartes papier numérisées           |
| Recherche tenprint-to-ten   | Identification rapide (< 2s) en base nationale           |
| Recherche latent-to-ten     | Empreintes scène de crime → identité                     |
| Recherche palm print        | Paume complète et partielle                              |
| Interface INTERPOL AFIS     | Échange via I-24/7 BCN Port-au-Prince                    |
| Déduplication biométrique   | Fusion doublons identités SNISID                         |
| NFIQ2 scoring               | Contrôle qualité obligatoire (score ≥ 60)                |

---

## 2. ARCHITECTURE

```
┌──────────────────────────────────────────────────────────────┐
│                       AFIS-HT STACK                          │
├──────────────────────────────────────────────────────────────┤
│  Clients                                                      │
│  ┌────────────┐ ┌──────────────┐ ┌──────────────────────┐   │
│  │ Console PNH│ │ Labo Forensi │ │  App Mobile Agents   │   │
│  └──────┬─────┘ └──────┬───────┘ └──────────┬───────────┘   │
├─────────┼──────────────┼──────────────────────────────────── │
│  API Gateway (Kong + mTLS/SPIFFE)                            │
├─────────────────────────────────────────────────────────────┤
│  Microservices Go                                             │
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐    │
│  │ afis-svc    │  │ latent-svc   │  │ quality-svc     │    │
│  │ (enrollment)│  │ (crime scene)│  │ (NFIQ2 scoring) │    │
│  └──────┬──────┘  └──────┬───────┘  └─────────┬───────┘    │
├─────────┼────────────────┼──────────────────────────────────┤
│  Couche données                                               │
│  ┌──────┐ ┌───────┐ ┌──────────┐ ┌───────┐ ┌────────────┐  │
│  │ PgSQL│ │ Redis │ │ Milvus   │ │ MinIO │ │ Kafka      │  │
│  │ Meta │ │ Cache │ │ Vectors  │ │Images │ │ Events     │  │
│  └──────┘ └───────┘ └──────────┘ └───────┘ └────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Structure dépôt

```
services/afis-svc/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── fingerprint.go
│   │   ├── latent.go
│   │   ├── search_result.go
│   │   └── enums.go
│   ├── repository/
│   │   ├── postgres/fingerprint_repo.go
│   │   ├── milvus/vector_repo.go
│   │   └── minio/image_repo.go
│   ├── service/
│   │   ├── enrollment_service.go
│   │   ├── search_service.go
│   │   ├── latent_service.go
│   │   └── quality_service.go
│   └── api/rest/
│       ├── enroll_handler.go
│       ├── search_handler.go
│       └── latent_handler.go
└── Dockerfile
```

---

## 3. BASE DE DONNÉES

### Migration 001 — Enrôlements

```sql
BEGIN;

CREATE TYPE afis_finger_position AS ENUM (
    'RIGHT_THUMB','RIGHT_INDEX','RIGHT_MIDDLE','RIGHT_RING','RIGHT_LITTLE',
    'LEFT_THUMB','LEFT_INDEX','LEFT_MIDDLE','LEFT_RING','LEFT_LITTLE',
    'RIGHT_PALM','LEFT_PALM','UNKNOWN'
);

CREATE TYPE afis_capture_method AS ENUM (
    'LIVESCANNER','INKROLL','LATENT_LIFT','PHOTO','UNKNOWN'
);

CREATE TYPE afis_subject_type AS ENUM (
    'SUSPECT','CRIMINAL','VICTIM','UNKNOWN_DECEASED','MISSING_PERSON','EMPLOYEE'
);

CREATE TABLE afis_subjects (
    subject_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snisid_person_id    UUID,                          -- Lien identité SNISID
    fir_record_id       UUID,                          -- Lien casier judiciaire
    subject_type        afis_subject_type NOT NULL,
    national_afis_id    VARCHAR(20) UNIQUE,            -- Format: AFIS-AAAA-NNNNNNN
    alias_ids           UUID[] DEFAULT '{}',
    enrolment_date      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    enrolling_unit      VARCHAR(50) NOT NULL,
    enrolling_officer   UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE afis_fingerprints (
    print_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_id          UUID NOT NULL REFERENCES afis_subjects(subject_id),
    finger_position     afis_finger_position NOT NULL,
    capture_method      afis_capture_method NOT NULL DEFAULT 'LIVESCANNER',
    nfiq2_score         SMALLINT CHECK (nfiq2_score BETWEEN 0 AND 100),
    quality_accepted    BOOLEAN GENERATED ALWAYS AS (nfiq2_score >= 60) STORED,
    image_ref           VARCHAR(500) NOT NULL,          -- Référence MinIO (WSQ format)
    minutiae_count      SMALLINT,
    milvus_vector_id    VARCHAR(100),                   -- ID vecteur dans Milvus
    template_version    VARCHAR(10) DEFAULT 'ISO_2011',
    is_primary          BOOLEAN DEFAULT FALSE,
    captured_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE afis_latent_prints (
    latent_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_reference      VARCHAR(100) NOT NULL,          -- Réf dossier PNH
    crime_scene_id      UUID,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    found_at            TIMESTAMPTZ NOT NULL,
    image_ref           VARCHAR(500) NOT NULL,
    nfiq2_score         SMALLINT,
    finger_position     afis_finger_position DEFAULT 'UNKNOWN',
    is_identified       BOOLEAN DEFAULT FALSE,
    matched_subject_id  UUID REFERENCES afis_subjects(subject_id),
    match_score         DECIMAL(5,2),
    examined_by         UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE afis_search_transactions (
    transaction_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_type    VARCHAR(20) NOT NULL,           -- TEN2TEN, LATENT2TEN, PALM
    query_subject_id    UUID,
    query_latent_id     UUID,
    hits_count          SMALLINT DEFAULT 0,
    top_score           DECIMAL(5,2),
    top_match_id        UUID,
    search_duration_ms  INTEGER,
    requested_by        UUID NOT NULL,
    requesting_unit     VARCHAR(50),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_afis_subjects_snisid    ON afis_subjects(snisid_person_id) WHERE snisid_person_id IS NOT NULL;
CREATE INDEX idx_afis_prints_subject     ON afis_fingerprints(subject_id, finger_position);
CREATE INDEX idx_afis_prints_quality     ON afis_fingerprints(quality_accepted) WHERE quality_accepted = TRUE;
CREATE INDEX idx_afis_latent_case        ON afis_latent_prints(case_reference);
CREATE INDEX idx_afis_latent_identified  ON afis_latent_prints(is_identified);

COMMIT;
```

---

## 4. MICROSERVICE GO — DOMAINE PRINCIPAL

```go
// services/afis-svc/internal/domain/fingerprint.go
package domain

import (
    "errors"
    "time"
    "github.com/google/uuid"
)

type Fingerprint struct {
    PrintID          uuid.UUID       `json:"print_id" db:"print_id"`
    SubjectID        uuid.UUID       `json:"subject_id" db:"subject_id"`
    FingerPosition   FingerPosition  `json:"finger_position" db:"finger_position"`
    CaptureMethod    CaptureMethod   `json:"capture_method" db:"capture_method"`
    NFIQ2Score       int16           `json:"nfiq2_score" db:"nfiq2_score"`
    QualityAccepted  bool            `json:"quality_accepted" db:"quality_accepted"`
    ImageRef         string          `json:"image_ref" db:"image_ref"`
    MinutiaeCount    *int16          `json:"minutiae_count,omitempty" db:"minutiae_count"`
    MilvusVectorID   *string         `json:"milvus_vector_id,omitempty" db:"milvus_vector_id"`
    CapturedAt       time.Time       `json:"captured_at" db:"captured_at"`
    CreatedBy        uuid.UUID       `json:"created_by" db:"created_by"`
}

var ErrQualityTooLow = errors.New("qualité empreinte insuffisante: NFIQ2 score < 60")
var ErrMissingRequiredFingers = errors.New("empreintes obligatoires manquantes (pouces + index requis)")

func (f *Fingerprint) IsHighQuality() bool { return f.NFIQ2Score >= 80 }

type SearchResult struct {
    CandidateID   uuid.UUID `json:"candidate_id"`
    SubjectID     uuid.UUID `json:"subject_id"`
    Score         float64   `json:"score"`
    Rank          int       `json:"rank"`
    NationalAFISID string  `json:"national_afis_id"`
}

type EnrollmentRequest struct {
    SubjectType    SubjectType    `json:"subject_type" validate:"required"`
    SNISIDPersonID *uuid.UUID     `json:"snisid_person_id,omitempty"`
    FIRRecordID    *uuid.UUID     `json:"fir_record_id,omitempty"`
    EnrollingUnit  string         `json:"enrolling_unit" validate:"required"`
    Fingerprints   []FingerprintCapture `json:"fingerprints" validate:"required,min=2"`
}

type FingerprintCapture struct {
    Position    FingerPosition `json:"position" validate:"required"`
    Method      CaptureMethod  `json:"method"`
    ImageBase64 string         `json:"image_base64" validate:"required"` // WSQ/PNG
    NFIQ2Score  int16          `json:"nfiq2_score" validate:"required,min=0,max=100"`
}
```

```go
// services/afis-svc/internal/service/search_service.go
package service

import (
    "context"
    "fmt"
    "sort"
    "github.com/snisid/afis-svc/internal/domain"
)

const (
    MinMatchScore     = 0.85  // Seuil hit automatique
    CandidateListSize = 15    // Top-N candidats retournés
)

func (s *SearchService) SearchTenprint(
    ctx context.Context,
    req domain.EnrollmentRequest,
) ([]domain.SearchResult, error) {
    // 1. Vectoriser les empreintes soumises
    vectors, err := s.vectorizer.Vectorize(ctx, req.Fingerprints)
    if err != nil {
        return nil, fmt.Errorf("vectorisation échouée: %w", err)
    }

    // 2. Recherche ANN dans Milvus (< 500ms pour 10M empreintes)
    candidates, err := s.milvusRepo.SearchNearest(ctx, vectors, CandidateListSize)
    if err != nil {
        return nil, fmt.Errorf("recherche Milvus échouée: %w", err)
    }

    // 3. Vérification minutiae (matcher forensique)
    var results []domain.SearchResult
    for _, c := range candidates {
        score, _ := s.matcher.CompareMinutiae(ctx, req.Fingerprints, c.SubjectID)
        if score >= MinMatchScore {
            results = append(results, domain.SearchResult{
                CandidateID:    c.PrintID,
                SubjectID:      c.SubjectID,
                Score:          score,
                NationalAFISID: c.NationalAFISID,
            })
        }
    }

    sort.Slice(results, func(i, j int) bool {
        return results[i].Score > results[j].Score
    })
    for i := range results { results[i].Rank = i + 1 }
    return results, nil
}
```

---

## 5. API REST

| Méthode | Endpoint                             | Rôle requis         | Description                         |
|---------|--------------------------------------|---------------------|-------------------------------------|
| `POST`  | `/api/v1/afis/enroll`                | PNH_OFFICER, DCPJ   | Enrôlement nouveau sujet            |
| `POST`  | `/api/v1/afis/search/tenprint`       | PNH_OFFICER, DCPJ   | Identification par 10 doigts        |
| `POST`  | `/api/v1/afis/search/latent`         | LAB_FORENSIC, DCPJ  | Empreinte latente → identité        |
| `GET`   | `/api/v1/afis/subjects/:id`          | Toute unité PNH     | Profil AFIS complet                 |
| `GET`   | `/api/v1/afis/subjects/:id/history`  | DCPJ, LAB_FORENSIC  | Historique recherches du sujet      |
| `POST`  | `/api/v1/afis/latents`               | LAB_FORENSIC        | Soumettre empreinte latente         |
| `PATCH` | `/api/v1/afis/latents/:id/match`     | LAB_FORENSIC, DCPJ  | Confirmer correspondance latente    |
| `GET`   | `/api/v1/afis/quality/check`         | Tout agent          | Vérifier qualité image (NFIQ2)      |
| `GET`   | `/api/v1/afis/stats`                 | DCPJ_ADMIN          | Statistiques base nationale         |

---

## 6. INTÉGRATIONS INTER-MODULES

- **FIR-HT (MP-20)** — à l'enrôlement, créer/lier le casier judiciaire automatiquement
- **FPR-HT (MP-17)** — lier `afis_subjects.snisid_person_id` aux personnes recherchées
- **SNISID-BIO-ADN** — combiner profil AFIS + ADN pour identification cadavres (RVIN-HT)
- **RDEP-HT (MP-22)** — à l'arrivée d'un déporté, vérification AFIS contre base nationale
- **DIPE-HT (MP-43)** — empreintes latentes de scènes d'enlèvement → identité
- **INTERPOL AFIS** — transmission empreintes non identifiées via I-24/7 (BCN Haiti)

---

## 7. SÉCURITÉ ET RBAC

| Rôle              | Enrôlement | Recherche 10P | Latent | Admin |
|-------------------|-----------|---------------|--------|-------|
| `PNH_OFFICER`     | ✅         | ✅             | ❌      | ❌     |
| `LAB_FORENSIC`    | ✅         | ✅             | ✅      | ❌     |
| `DCPJ_SUPERVISOR` | ✅         | ✅             | ✅      | ✅     |
| `PARQUET_PROC`    | ❌         | ✅ (lecture)   | ❌      | ❌     |

- Chiffrement biométrique : AES-256 au repos, TLS 1.3 en transit
- Images WSQ stockées dans MinIO chiffré (bucket: `afis-biometric`)
- Vecteurs Milvus : non réversibles (hash + sel cryptographique)
- Audit : toute recherche loguée dans Kafka topic `afis.audit`

---

## 8. TESTS REQUIS

```
tests/
├── unit/
│   ├── enrollment_service_test.go      # TestEnroll_10Fingers_Valid
│   ├── quality_service_test.go         # TestNFIQ2_Score_Threshold
│   └── search_service_test.go          # TestSearch_HitAbove85percent
├── integration/
│   ├── milvus_search_test.go           # TestMilvus_SearchLatency_Sub500ms
│   └── latent_match_test.go            # TestLatent_SceneOfCrime_Match
```

---

## 9. ORDRE D'IMPLÉMENTATION

| Semaine | Livrable                              | Priorité  |
|---------|---------------------------------------|-----------|
| S1      | SQL migrations + domaine Go           | 🔴 Critique |
| S1      | API enrollment + qualité NFIQ2        | 🔴 Critique |
| S2      | Intégration Milvus + vectorisation    | 🔴 Critique |
| S2      | Recherche tenprint (ANN + minutiae)   | 🔴 Critique |
| S3      | Module latent (scènes de crime)       | 🟠 Haute    |
| S3      | Interface INTERPOL AFIS               | 🟠 Haute    |
| S4      | Intégration FIR-HT, RDEP-HT          | 🟠 Haute    |
| S5      | Tests couverture ≥ 85%                | 🟡 Normale  |

---

## 10. VARIABLES D'ENVIRONNEMENT

```dotenv
AFIS_DB_HOST=localhost
AFIS_DB_NAME=snisid_afis
AFIS_REDIS_ADDR=redis-master:6379
AFIS_MILVUS_ADDR=milvus:19530
AFIS_MILVUS_COLLECTION=afis_fingerprints
AFIS_MINIO_BUCKET=afis-biometric
AFIS_NFIQ2_MIN_SCORE=60
AFIS_MATCH_THRESHOLD=0.85
AFIS_CANDIDATE_LIST_SIZE=15
INTERPOL_AFIS_GATEWAY=https://i247-afis.pnh.gov.ht
AFIS_SERVICE_PORT=8091
```

---

## 11. FICHIERS À CRÉER (47 fichiers)

```
migrations/afis/
  001_afis_enums.sql
  002_afis_subjects.sql
  003_afis_fingerprints.sql
  004_afis_latent_prints.sql
  005_afis_search_transactions.sql
  006_afis_indexes_rls.sql

services/afis-svc/
  cmd/server/main.go
  internal/domain/{fingerprint,latent,search_result,enums}.go
  internal/repository/postgres/{fingerprint,latent,subject}_repo.go
  internal/repository/milvus/vector_repo.go
  internal/repository/minio/image_repo.go
  internal/service/{enrollment,search,latent,quality}_service.go
  internal/api/rest/{enroll,search,latent,quality}_handler.go
  internal/nfiq2/scorer.go
  internal/matcher/minutiae_matcher.go
  tests/unit/{enrollment,search,quality}_service_test.go
  tests/integration/{milvus_search,latent_match}_test.go
  Dockerfile
```

---
*MP-19 — AFIS-HT — SNISID — République d'Haïti — Usage officiel*
