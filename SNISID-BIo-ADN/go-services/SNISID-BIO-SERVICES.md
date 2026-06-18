# SNISID-BIO-ADN — Services Go
**Document ID :** SNISID-BIO-SVC-001 | **Version :** 1.0.0

---

## 1. STRUCTURE DU PROJET GO

```
bio-adn-service/
├── cmd/
│   ├── api/main.go              # Serveur REST FastAPI (bridgé via gRPC)
│   └── grpc/main.go             # Serveur gRPC interne
├── internal/
│   ├── engine/
│   │   ├── matcher.go           # Moteur matching ADN STR
│   │   ├── matcher_test.go
│   │   ├── familial.go          # Algorithme recherche familiale
│   │   └── scoring.go          # Likelihood Ratio scoring
│   ├── indexes/
│   │   ├── dna.go               # CRUD index ADN
│   │   ├── persons.go           # CRUD index personnes
│   │   └── property.go          # CRUD index biens
│   ├── sync/
│   │   ├── ldis_uploader.go     # Upload LDIS → SDIS
│   │   ├── sdis_uploader.go     # Upload SDIS → NDIS
│   │   └── scheduler.go         # Ordonnancement hebdomadaire
│   ├── crypto/
│   │   ├── hsm.go               # Interface PKCS#11 Luna HSM
│   │   └── envelope.go          # Envelope encryption
│   ├── audit/
│   │   └── logger.go            # Audit immuable + signature ECDSA
│   ├── kafka/
│   │   ├── producer.go
│   │   └── consumer.go
│   └── db/
│       ├── postgres.go          # Connexion PostgreSQL
│       └── redis.go             # Cache Redis
├── proto/
│   └── bio_adn.proto            # Définitions gRPC
├── pkg/
│   └── models/                  # Types partagés
└── migrations/
    ├── 001_indexes_adn.sql
    ├── 002_indexes_personnes.sql
    └── 003_indexes_biens.sql
```

---

## 2. SERVICE gRPC PRINCIPAL

```protobuf
// proto/bio_adn.proto
syntax = "proto3";
package ht.gov.snisid.bio;
option go_package = "github.com/snisid/bio-adn/proto";

import "google/protobuf/timestamp.proto";

// ──────────────────────────────────────────────
// DNA Matching Service
// ──────────────────────────────────────────────
service DNAMatchingService {
    rpc SubmitProfile   (SubmitProfileRequest)  returns (SubmitProfileResponse);
    rpc SearchProfile   (SearchProfileRequest)  returns (SearchProfileResponse);
    rpc GetHitDetails   (HitDetailsRequest)     returns (HitDetailsResponse);
    rpc ExpungeProfile  (ExpungeRequest)         returns (ExpungeResponse);
}

// ──────────────────────────────────────────────
// Wanted Persons Service
// ──────────────────────────────────────────────
service WantedPersonsService {
    rpc CreateRecord    (CreateWantedRequest)    returns (WantedRecord);
    rpc QueryRecord     (QueryWantedRequest)     returns (QueryWantedResponse);
    rpc UpdateStatus    (UpdateStatusRequest)    returns (WantedRecord);
    rpc VerifyHit       (VerifyHitRequest)       returns (VerifyHitResponse);
}

// ──────────────────────────────────────────────
// LAPI Query Service (temps réel < 200ms)
// ──────────────────────────────────────────────
service LAPIQueryService {
    rpc QueryPlate      (PlateQueryRequest)      returns (PlateQueryResponse);
    rpc QueryVIN        (VINQueryRequest)        returns (VINQueryResponse);
}

// ──────────────────────────────────────────────
// Messages
// ──────────────────────────────────────────────
message SubmitProfileRequest {
    string specimen_number    = 1;
    string index_type         = 2;
    bytes  loci_data          = 3; // JSON sérialisé + chiffré
    float  quality_score      = 4;
    string lab_id             = 5;
    string case_number        = 6;
    string collected_date     = 7;
    string correlation_id     = 8;
}

message SubmitProfileResponse {
    string sample_id          = 1;
    bool   accepted           = 2;
    string rejection_reason   = 3;
}

message SearchProfileRequest {
    bytes  loci_data          = 1;
    string index_type         = 2;
    string case_number        = 3;
    string requesting_agency  = 4;
    string officer_niu        = 5;
    string purpose            = 6;
    float  min_confidence     = 7; // défaut 0.85
    bool   include_familial   = 8; // défaut false
}

message SearchProfileResponse {
    repeated MatchResult hits = 1;
    int32    total_hits       = 2;
    int32    search_duration_ms = 3;
}

message MatchResult {
    string hit_id             = 1;
    string match_sample_id    = 2;
    string match_type         = 3;
    float  confidence         = 4;
    int32  matched_loci       = 5;
    int32  total_loci         = 6;
    string alert_level        = 7;
    string match_index_type   = 8;
    string case_number        = 9;
    string mco_contact        = 10;
}

message PlateQueryRequest {
    string plate_number       = 1;
    string camera_id          = 2;
    string location           = 3;
    string query_id           = 4;
}

message PlateQueryResponse {
    string query_id           = 1;
    bool   hit_found          = 2;
    string hit_type           = 3; // STOLEN_VEHICLE, WANTED_PERSON, etc.
    string record_number      = 4;
    string alert_level        = 5;
    string mco_contact        = 6;
    int32  response_ms        = 7;
}

message ExpungeRequest {
    string sample_id          = 1;
    string court_order_ref    = 2;
    string ordered_by         = 3;
    string reason             = 4;
    string officer_niu        = 5;
}

message ExpungeResponse {
    bool   success            = 1;
    string expunge_id         = 2;
    string timestamp          = 3;
}
```

---

## 3. SYNC SCHEDULER (LDIS → SDIS → NDIS)

```go
// internal/sync/scheduler.go
package sync

import (
    "context"
    "log"
    "time"
)

// SyncScheduler orchestre la synchronisation LDIS→SDIS→NDIS
// Même cadence que CODIS (hebdomadaire pour NDIS, quotidien pour SDIS)
type SyncScheduler struct {
    ldisUploader *LDISUploader
    sdisUploader *SDISUploader
    level        string // LDIS, SDIS
}

func NewSyncScheduler(level string) *SyncScheduler {
    return &SyncScheduler{level: level}
}

func (s *SyncScheduler) Start(ctx context.Context) {
    switch s.level {
    case "LDIS":
        // Upload LDIS → SDIS : quotidien à 02h00 (heure Haiti, UTC-4)
        go s.runSchedule(ctx, "LDIS→SDIS", "02:00", s.ldisUploader.Upload)

    case "SDIS":
        // Upload SDIS → NDIS : hebdomadaire le dimanche à 03h00
        go s.runSchedule(ctx, "SDIS→NDIS", "sunday-03:00", s.sdisUploader.Upload)
    }
}

func (s *SyncScheduler) runSchedule(ctx context.Context, name, schedule string, fn func(context.Context) error) {
    for {
        next := s.nextRun(schedule)
        select {
        case <-ctx.Done():
            return
        case <-time.After(time.Until(next)):
            log.Printf("[BIO-SYNC] Démarrage synchronisation %s", name)
            if err := fn(ctx); err != nil {
                log.Printf("[BIO-SYNC] ERREUR %s: %v", name, err)
                // Kafka : publier event d'erreur pour monitoring
            } else {
                log.Printf("[BIO-SYNC] Succès %s", name)
            }
        }
    }
}

func (s *SyncScheduler) nextRun(schedule string) time.Time {
    now := time.Now()
    // Calcul de la prochaine exécution selon le schedule
    // ... (implémentation selon cron-like parsing)
    return now.Add(24 * time.Hour) // placeholder
}
```

---

## 4. SERVICE LAPI (< 200ms SLA)

```go
// internal/indexes/lapi_query.go
package indexes

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

const LAPIMaxResponseMs = 200

type LAPIQueryService struct {
    redis *redis.Client
    db    Database
}

// QueryPlate — Interrogation plaque temps réel pour LAPI (MP-16)
// SLA : < 200ms P99
func (s *LAPIQueryService) QueryPlate(ctx context.Context, plate string) (*PlateHitResult, error) {
    start := time.Now()

    // 1. Cache Redis L1 (< 1ms si présent)
    cacheKey := fmt.Sprintf("lapi:plate:%s", plate)
    if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil {
        return parseCachedHit(cached), nil
    }

    // 2. Base de données (BIE-VEH + BIE-PLQ)
    ctx, cancel := context.WithTimeout(ctx, 150*time.Millisecond)
    defer cancel()

    hit, err := s.db.QueryPlateIndex(ctx, plate)
    if err != nil {
        return &PlateHitResult{HitFound: false, ResponseMs: int(time.Since(start).Milliseconds())}, nil
    }

    // 3. Mettre en cache si pas de hit (TTL 5 minutes)
    if hit == nil {
        s.redis.Set(ctx, cacheKey, "NO_HIT", 5*time.Minute)
        return &PlateHitResult{HitFound: false, ResponseMs: int(time.Since(start).Milliseconds())}, nil
    }

    // 4. Hit trouvé — cache 30 secondes (données critiques, TTL court)
    s.redis.Set(ctx, cacheKey, marshalHit(hit), 30*time.Second)

    responseMs := int(time.Since(start).Milliseconds())
    if responseMs > LAPIMaxResponseMs {
        // Alerte SLA breach sur Kafka
        go s.publishSLABreach(plate, responseMs)
    }

    return &PlateHitResult{
        HitFound:     true,
        HitType:      hit.HitType,
        RecordNumber: hit.RecordNumber,
        AlertLevel:   hit.AlertLevel,
        MCOContact:   hit.MCOContact,
        ResponseMs:   responseMs,
    }, nil
}
```
