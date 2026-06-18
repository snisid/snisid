# MP-22 — RDEP-HT
## Registre des Déportés et Extradés d'Haïti
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL SEULEMENT
Module           : MP-22
Code SNISID      : RDEP-HT
Version          : 1.0.0
Dépendances      : FIR-HT (MP-20), AFIS-HT (MP-19), GANG-HT (MP-24), SIFR-HT (MP-33)
Normes           : IML/IOM, FBI IAFIS, Accord bilatéral HT-USA, Décret haïtien immigration
Acteurs          : DGI (Direction Générale Immigration), PNH, DCPJ, Frontières
```

---

## 1. CONTEXTE

Haïti reçoit environ 15,000-20,000 déportés par an des États-Unis, du Canada,
de la République Dominicaine et d'autres pays. Le gang 400 Mawozo est composé
majoritairement d'ex-déportés américains. Sans registre centralisé :
- Les déportés avec antécédents criminels graves sont libérés sans surveillance
- Impossible de lier un gang haïtien à un alias américain documenté
- Les accords bilatéraux de partage d'information ne sont pas exploités

### Flux de déportations documentés (2024)

| Pays d'origine  | Estimé/an | Avec antécédents criminels | Mode arrivée      |
|-----------------|-----------|---------------------------|-------------------|
| États-Unis      | 8,000+    | ~40%                      | Vols charter ICE  |
| Rép. Dominicaine| 50,000+   | ~15%                      | Postes frontières |
| Bahamas         | 800+      | ~25%                      | Traversée maritime|
| Canada          | 200+      | ~35%                      | Vols Air Canada   |
| Autres          | 500+      | Variable                  | Divers            |

---

## 2. ARCHITECTURE

```
services/rdep-svc/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── deportee.go
│   │   ├── foreign_record.go
│   │   ├── monitoring.go
│   │   └── enums.go
│   ├── repository/postgres/
│   │   ├── deportee_repo.go
│   │   └── foreign_record_repo.go
│   ├── service/
│   │   ├── intake_service.go
│   │   ├── screening_service.go
│   │   └── monitoring_service.go
│   └── api/rest/
│       ├── intake_handler.go
│       └── monitoring_handler.go
└── Dockerfile
```

---

## 3. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE rdep_deportation_country AS ENUM (
    'USA','CAN','DOM','BHS','CUB','JAM','TTO','MEX','BRA','FRA','OTHER'
);

CREATE TYPE rdep_criminal_risk AS ENUM (
    'NONE','LOW','MEDIUM','HIGH','VERY_HIGH'
);

CREATE TYPE rdep_monitoring_status AS ENUM (
    'ACTIVE','SUSPENDED','COMPLETED','FLED','DECEASED'
);

CREATE TABLE rdep_deportees (
    deportee_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_rdep_id    VARCHAR(25) UNIQUE NOT NULL,   -- Format: RDEP-HT-AAAA-NNNNNN
    snisid_person_id    UUID NOT NULL,
    fir_record_id       UUID,
    afis_subject_id     UUID,

    -- Informations déportation
    deportation_country rdep_deportation_country NOT NULL,
    deportation_date    TIMESTAMPTZ NOT NULL,
    arrival_port        VARCHAR(100) NOT NULL,          -- PAP, CAP, Frontière Malpasse, etc.
    arrival_dept_code   CHAR(2),
    deporting_agency    VARCHAR(100),                   -- ICE, CBSA, DNCD-DOM, etc.
    deportation_reason  TEXT,
    flight_number       VARCHAR(20),

    -- Identité étrangère
    foreign_name        VARCHAR(200),
    foreign_aliases     TEXT[] DEFAULT '{}',
    foreign_id_number   VARCHAR(100),                   -- SSN/SIN masqué, etc.
    foreign_country_id  VARCHAR(50),

    -- Antécédents criminels étrangers
    has_foreign_record  BOOLEAN DEFAULT FALSE,
    criminal_risk_level rdep_criminal_risk DEFAULT 'NONE',
    convicted_offenses  TEXT[] DEFAULT '{}',
    gang_affiliated     BOOLEAN DEFAULT FALSE,
    gang_name           VARCHAR(100),

    -- Surveillance
    monitoring_required BOOLEAN DEFAULT FALSE,
    monitoring_status   rdep_monitoring_status DEFAULT 'ACTIVE',
    monitoring_unit     VARCHAR(50),
    monitoring_officer  UUID,
    monitoring_end_date TIMESTAMPTZ,

    -- Localisation actuelle
    current_address     TEXT,
    current_commune     VARCHAR(100),
    current_dept_code   CHAR(2),

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rdep_foreign_records (
    foreign_record_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deportee_id         UUID NOT NULL REFERENCES rdep_deportees(deportee_id),
    country             rdep_deportation_country NOT NULL,
    court_name          VARCHAR(200),
    offense_description TEXT NOT NULL,
    offense_date        TIMESTAMPTZ,
    conviction_date     TIMESTAMPTZ,
    sentence            TEXT,
    prison_served       TEXT,
    fbi_number          VARCHAR(50),
    interpol_ref        VARCHAR(50),
    source_document     VARCHAR(500),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rdep_monitoring_events (
    event_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deportee_id         UUID NOT NULL REFERENCES rdep_deportees(deportee_id),
    event_type          VARCHAR(50) NOT NULL,  -- CHECK_IN, VIOLATION, ADDRESS_CHANGE
    event_date          TIMESTAMPTZ NOT NULL,
    location_lat        DECIMAL(10,7),
    location_lng        DECIMAL(10,7),
    notes               TEXT,
    reported_by         UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rdep_deportees_person      ON rdep_deportees(snisid_person_id);
CREATE INDEX idx_rdep_deportees_country     ON rdep_deportees(deportation_country);
CREATE INDEX idx_rdep_deportees_risk        ON rdep_deportees(criminal_risk_level) WHERE criminal_risk_level IN ('HIGH','VERY_HIGH');
CREATE INDEX idx_rdep_deportees_gang        ON rdep_deportees(gang_affiliated) WHERE gang_affiliated = TRUE;
CREATE INDEX idx_rdep_deportees_monitoring  ON rdep_deportees(monitoring_status) WHERE monitoring_required = TRUE;

COMMIT;
```

---

## 4. SERVICE CLÉ — SCREENING À L'ARRIVÉE

```go
package service

import (
    "context"
    "github.com/snisid/rdep-svc/internal/domain"
)

type ScreeningService struct {
    repo      domain.DeporteeRepository
    fbiClient domain.FBIRecordClient     // API FBI IAFIS
    interp    domain.InterpolClient
    gangRepo  domain.GangRepository
    kafka     domain.EventPublisher
}

// ScreenDeportee effectue le criblage complet à l'arrivée
func (s *ScreeningService) ScreenDeportee(
    ctx context.Context,
    intake domain.DeporteeIntakeRequest,
) (*domain.ScreeningResult, error) {
    result := &domain.ScreeningResult{
        PersonID: intake.SNISIDPersonID,
        RiskLevel: domain.RiskNone,
    }

    // 1. Vérification AFIS local (empreintes prises à l'arrivée)
    afisHit, _ := s.afisClient.CheckPrint(ctx, intake.FingerprintData)
    if afisHit != nil {
        result.HasLocalRecord = true
        result.LocalRecordID = afisHit.SubjectID
        result.RiskLevel = elevateRisk(result.RiskLevel, domain.RiskMedium)
    }

    // 2. Vérification records FBI (si déporté USA)
    if intake.DeportationCountry == "USA" && intake.FBINumber != "" {
        fbiRecord, err := s.fbiClient.GetRecord(ctx, intake.FBINumber)
        if err == nil && fbiRecord != nil {
            result.HasForeignRecord = true
            result.ForeignRecords = append(result.ForeignRecords, *fbiRecord)
            if fbiRecord.HasViolentOffenses() {
                result.RiskLevel = elevateRisk(result.RiskLevel, domain.RiskHigh)
            }
        }
    }

    // 3. Vérification affiliation gang
    if intake.GangName != "" {
        gangMatch, _ := s.gangRepo.FindByName(ctx, intake.GangName)
        if gangMatch != nil {
            result.GangAffiliated = true
            result.GangID = gangMatch.GangID
            result.RiskLevel = elevateRisk(result.RiskLevel, domain.RiskVeryHigh)
        }
    }

    // 4. Vérification INTERPOL notices
    notices, _ := s.interp.CheckNotices(ctx, intake.SNISIDPersonID)
    if len(notices) > 0 {
        result.InterpolNotices = notices
        result.RiskLevel = elevateRisk(result.RiskLevel, domain.RiskVeryHigh)
    }

    // Publier résultat pour décision monitoring
    _ = s.kafka.Publish(ctx, "rdep.screening.completed", result)
    return result, nil
}

func elevateRisk(current, proposed domain.CriminalRisk) domain.CriminalRisk {
    if proposed > current { return proposed }
    return current
}
```

---

## 5. API REST

| Méthode | Endpoint                              | Rôle requis        | Description                         |
|---------|---------------------------------------|--------------------|-------------------------------------|
| `POST`  | `/api/v1/rdep/intake`                 | DGI, SIFR_AGENT    | Enregistrer arrivée déporté         |
| `POST`  | `/api/v1/rdep/:id/screen`             | DCPJ, DGI          | Lancer screening complet            |
| `GET`   | `/api/v1/rdep/:id`                    | DCPJ, DGI          | Profil complet déporté              |
| `POST`  | `/api/v1/rdep/:id/monitoring/events`  | PNH_OFFICER        | Enregistrer événement surveillance  |
| `GET`   | `/api/v1/rdep/high-risk`              | DCPJ, PNH_ADMIN    | Liste déportés haut risque          |
| `GET`   | `/api/v1/rdep/gang-affiliated`        | DCPJ, GANG_UNIT    | Déportés à affiliation gang         |
| `GET`   | `/api/v1/rdep/stats/by-country`       | DGI, MJSP          | Statistiques par pays d'origine     |

---

## 6. INTÉGRATIONS

- **GANG-HT** : Si gang_affiliated → lier automatiquement dans registre des gangs
- **SIFR-HT** : Alerte frontière si déporté à risque élevé tente de re-entrer
- **SIGEO-HT** : Cartographier les communes avec forte concentration de déportés à risque
- **SANC-HT** : Vérification automatique contre listes OFAC/ONU à l'arrivée

---

## 7. VARIABLES D'ENVIRONNEMENT

```dotenv
RDEP_DB_HOST=localhost
RDEP_DB_NAME=snisid_rdep
RDEP_FBI_API_URL=https://api.fbi.gov/records
RDEP_FBI_API_KEY=<VAULT:rdep/fbi_api_key>
RDEP_INTERPOL_GATEWAY=https://i247-gateway.pnh.gov.ht
RDEP_AFIS_SERVICE_URL=http://afis-svc:8091
RDEP_GANG_SERVICE_URL=http://gang-svc:8095
RDEP_MONITORING_PERIOD_DAYS=365
RDEP_SERVICE_PORT=8094
```

---
*MP-22 — RDEP-HT — Registre Déportés — SNISID — République d'Haïti*
