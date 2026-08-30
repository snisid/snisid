# MP-27 — SANC-HT
## Interface Nationale de Synchronisation avec les Listes de Sanctions Internationales
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-27 | Code : SANC-HT
Dépendances      : GANG-HT (MP-24), CHEF-HT (MP-25), RDEP-HT (MP-22), UCREF-INT (MP-39)
Normes           : Résolution CSNU 2653, GAFI/FATF R.6, OFAC SDN, UE 2022/2337
Acteurs          : MJSP, UCREF, BRH (Banque République Haïti), DGI Douanes
```

---

## 1. CONTEXTE

La Résolution 2653 du Conseil de Sécurité de l'ONU (2022) impose gel d'avoirs et
interdiction de voyager contre les dirigeants de gangs haïtiens. L'OFAC américain a
sanctionné Jimmy Chérizier, Gabriel Jean-Pierre et d'autres. Ce module synchronise
automatiquement toutes les listes de sanctions et croise avec les identités SNISID.

### Sources de données

| Source              | Fréq. MàJ   | Format | Volume       | Priorité  |
|---------------------|-------------|--------|--------------|-----------|
| CSNU Comité 2653    | Hebdomadaire| XML    | ~50 entrées  | Critique  |
| OFAC SDN (USA)      | Quotidien   | XML    | 14,000+      | Critique  |
| UE Consolidated     | Quotidien   | XML    | 2,200+       | Haute     |
| INTERPOL Notices    | Temps réel  | API    | Variable     | Critique  |
| Canada OSFI         | Mensuel     | XML    | 2,500+       | Normale   |
| Royaume-Uni OFSI    | Quotidien   | XML    | 4,000+       | Haute     |

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE sanc_source AS ENUM (
    'UN_2653','OFAC_SDN','EU_CONSOLIDATED',
    'INTERPOL','CANADA_OSFI','UK_OFSI','OTHER'
);

CREATE TYPE sanc_measure AS ENUM (
    'ASSETS_FREEZE','TRAVEL_BAN','ARMS_EMBARGO','ALL_MEASURES'
);

CREATE TYPE sanc_entity_type AS ENUM (
    'INDIVIDUAL','ORGANIZATION','VESSEL','AIRCRAFT'
);

CREATE TABLE sanc_entries (
    sanc_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source              sanc_source NOT NULL,
    source_ref_id       VARCHAR(100) NOT NULL,
    entity_type         sanc_entity_type NOT NULL,
    entity_name         VARCHAR(300) NOT NULL,
    aliases             TEXT[] DEFAULT '{}',
    nationality         TEXT[] DEFAULT '{}',
    date_of_birth       DATE,
    place_of_birth      TEXT,
    passport_numbers    TEXT[] DEFAULT '{}',
    national_id_numbers TEXT[] DEFAULT '{}',
    measure_types       sanc_measure[] DEFAULT '{}',
    listing_date        DATE NOT NULL,
    end_date            DATE,
    is_active           BOOLEAN DEFAULT TRUE,
    listing_reason      TEXT,
    committee_notes     TEXT,
    snisid_person_id    UUID,
    gang_id             UUID,
    chef_member_id      UUID,
    match_confidence    SMALLINT,
    source_updated_at   TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(source, source_ref_id)
);

CREATE TABLE sanc_sync_log (
    sync_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source              sanc_source NOT NULL,
    started_at          TIMESTAMPTZ NOT NULL,
    completed_at        TIMESTAMPTZ,
    entries_processed   INTEGER DEFAULT 0,
    entries_added       INTEGER DEFAULT 0,
    entries_updated     INTEGER DEFAULT 0,
    entries_removed     INTEGER DEFAULT 0,
    errors              INTEGER DEFAULT 0,
    status              VARCHAR(20) DEFAULT 'RUNNING',
    error_details       TEXT
);

CREATE TABLE sanc_identity_matches (
    match_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sanc_id             UUID NOT NULL REFERENCES sanc_entries(sanc_id),
    snisid_person_id    UUID NOT NULL,
    match_score         DECIMAL(5,2) NOT NULL,
    match_fields        TEXT[] DEFAULT '{}',
    confirmed_by        UUID,
    is_confirmed        BOOLEAN DEFAULT FALSE,
    is_false_positive   BOOLEAN DEFAULT FALSE,
    reviewed_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sanc_source    ON sanc_entries(source, is_active);
CREATE INDEX idx_sanc_name_fts  ON sanc_entries USING gin(to_tsvector('simple', entity_name));
CREATE INDEX idx_sanc_aliases   ON sanc_entries USING gin(aliases);
CREATE INDEX idx_sanc_matches   ON sanc_identity_matches(snisid_person_id);

COMMIT;
```

---

## 3. SERVICE GO — SYNCHRONISATION ET MATCHING

```go
package service

import (
    "context"
    "encoding/xml"
    "fmt"
    "net/http"
    "strings"
    "time"
    "github.com/snisid/sanc-svc/internal/domain"
)

type SyncService struct {
    repo    domain.SanctionsRepository
    kafka   domain.EventPublisher
    snisid  domain.SNISIDClient
}

// SyncOFAC synchronise la liste SDN de l OFAC
func (s *SyncService) SyncOFAC(ctx context.Context) (*domain.SyncResult, error) {
    log := &domain.SyncLog{
        Source:    domain.SourceOFAC,
        StartedAt: time.Now(),
    }
    resp, err := http.Get("https://www.treasury.gov/ofac/downloads/sdn.xml")
    if err != nil {
        return nil, fmt.Errorf("telechargement OFAC SDN: %w", err)
    }
    defer resp.Body.Close()

    var sdnList OFACSDNList
    if err := xml.NewDecoder(resp.Body).Decode(&sdnList); err != nil {
        return nil, fmt.Errorf("parse XML OFAC: %w", err)
    }

    for _, entry := range sdnList.SDNEntries {
        sanc := s.convertOFACEntry(entry)
        if err := s.repo.UpsertEntry(ctx, sanc); err != nil {
            log.Errors++
            continue
        }
        log.EntriesProcessed++

        // Recherche correspondances SNISID (matching nom + DOB + nationalite)
        matches, _ := s.findSNISIDMatches(ctx, sanc)
        for _, m := range matches {
            if m.Score >= 0.85 {
                _ = s.repo.SaveMatch(ctx, m)
                _ = s.kafka.Publish(ctx, "sanc.match.found", m)
            }
        }
    }
    log.CompletedAt = time.Now()
    _ = s.repo.SaveSyncLog(ctx, log)
    return &domain.SyncResult{Log: log}, nil
}

// CheckPersonRealTime verifie une personne SNISID contre toutes listes actives
func (s *SyncService) CheckPersonRealTime(
    ctx context.Context,
    personID string,
) (*domain.PersonSanctionsResult, error) {
    person, err := s.snisid.GetPerson(ctx, personID)
    if err != nil {
        return nil, err
    }
    results, err := s.repo.SearchByNameAndDOB(ctx,
        person.FullName, person.DateOfBirth,
        person.Aliases, person.Nationalities)
    if err != nil {
        return nil, err
    }
    return &domain.PersonSanctionsResult{
        PersonID:  personID,
        IsSanctioned: len(results) > 0,
        Matches:   results,
        CheckedAt: time.Now(),
    }, nil
}

func (s *SyncService) fuzzyNameMatch(a, b string) float64 {
    a = strings.ToLower(strings.TrimSpace(a))
    b = strings.ToLower(strings.TrimSpace(b))
    if a == b {
        return 1.0
    }
    // Implementation Jaro-Winkler
    return jaroWinkler(a, b)
}
```

---

## 4. API REST

| Méthode | Endpoint                             | Rôle          | Description                         |
|---------|--------------------------------------|---------------|-------------------------------------|
| `GET`   | `/api/v1/sanc/check/:person_id`      | Tout SNISID   | Vérifier personne vs listes         |
| `POST`  | `/api/v1/sanc/check/name`            | UCREF, DGI    | Screening KYC par nom               |
| `GET`   | `/api/v1/sanc/entries`               | UCREF, MJSP   | Lister entrées actives (paginé)     |
| `GET`   | `/api/v1/sanc/entries/haiti`         | MJSP          | Entrées Haïti-spécifiques (CSNU)    |
| `GET`   | `/api/v1/sanc/matches/unconfirmed`   | MJSP, UCREF   | Correspondances à confirmer         |
| `POST`  | `/api/v1/sanc/matches/:id/confirm`   | MJSP_ADMIN    | Confirmer ou rejeter                |
| `POST`  | `/api/v1/sanc/sync/trigger`          | SUPERADMIN    | Synchronisation manuelle            |
| `GET`   | `/api/v1/sanc/sync/status`           | ADMIN         | Statut des synchronisations         |

---

## 5. SCHEDULER AUTOMATIQUE

```go
// Synchronisation automatique planifiée
func (s *SyncScheduler) Start(ctx context.Context) {
    // OFAC : toutes les 6 heures
    go s.scheduledSync(ctx, "OFAC_SDN", 6*time.Hour, s.syncSvc.SyncOFAC)
    // CSNU 2653 : toutes les 24 heures
    go s.scheduledSync(ctx, "UN_2653", 24*time.Hour, s.syncSvc.SyncUN2653)
    // UE : toutes les 12 heures
    go s.scheduledSync(ctx, "EU_CONSOLIDATED", 12*time.Hour, s.syncSvc.SyncEU)
    // Canada : toutes les 72 heures
    go s.scheduledSync(ctx, "CANADA_OSFI", 72*time.Hour, s.syncSvc.SyncCanadaOSFI)
}
```

---

## 6. VARIABLES D'ENVIRONNEMENT

```dotenv
SANC_DB_HOST=localhost
SANC_DB_NAME=snisid_sanc
SANC_OFAC_URL=https://www.treasury.gov/ofac/downloads/sdn.xml
SANC_EU_URL=https://data.europa.eu/api/hub/repo/datasets/eu-consolidated-sanctions
SANC_UN_URL=https://scsanctions.un.org/resources/xml/en/consolidated.xml
SANC_CANADA_URL=https://www.osfi-bsif.gc.ca/Eng/ToolsOutilsDataDonnees/Pages/RGL.aspx
SANC_SYNC_OFAC_HOURS=6
SANC_SYNC_UN_HOURS=24
SANC_MATCH_THRESHOLD=0.80
SANC_KAFKA_BROKERS=kafka:9092
SANC_SERVICE_PORT=8100
```

---
*MP-27 — SANC-HT — Listes Sanctions — SNISID — République d'Haïti*
