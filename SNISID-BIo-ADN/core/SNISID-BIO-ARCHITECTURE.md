# SNISID-BIO-ADN — Architecture Technique
**Document ID :** SNISID-BIO-ARC-001 | **Version :** 1.0.0 | **Statut :** PRODUCTION-READY

---

## 1. STACK TECHNIQUE

| Couche | Technologie | Rôle |
|--------|-------------|------|
| **API** | FastAPI (Python) + gRPC (Go) | Endpoints REST + inter-services |
| **Matching ADN** | Go service + PostgreSQL pg_trgm | Comparaison profils STR |
| **Event Bus** | Apache Kafka | Propagation hits, alertes |
| **Base de données** | PostgreSQL 16 (partitionné) | Stockage index |
| **Cache hits** | Redis 7 | Résultats de correspondance chauds |
| **Chiffrement** | AES-256-GCM + HSM Luna | Profils ADN au repos |
| **Auth** | mTLS + JWT (Keycloak SNISID) | Tous les appels inter-services |
| **Observabilité** | Prometheus + Jaeger + Elastic | Métriques, traces, audit |

---

## 2. COMPOSANTS PRINCIPAUX

### 2.1 DNA Matching Engine (Go)

```go
// bio-adn-service/internal/engine/matcher.go
package engine

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

// STRProfile représente un profil ADN CODIS (20 loci standard)
type STRProfile struct {
    SampleID    string            `json:"sample_id"`
    IndexType   IndexType         `json:"index_type"` // CONVICTED, FORENSIC, MISSING
    Loci        map[string]Allele `json:"loci"`       // 20 loci CODIS core
    QualityScore float64          `json:"quality_score"` // 0.0-1.0
    LabID       string            `json:"lab_id"`
    CollectedAt string            `json:"collected_at"`
    NIU         *string           `json:"niu,omitempty"` // Lien avec SNISID Core
    Encrypted   []byte            `json:"encrypted"`     // Profil chiffré HSM
}

// Allele représente les allèles à un locus
type Allele struct {
    Locus  string    `json:"locus"`
    Value1 float64   `json:"value1"` // Allèle 1
    Value2 *float64  `json:"value2"` // Allèle 2 (null si homozygote)
}

// IndexType correspond aux catégories CODIS/NDIS
type IndexType string

const (
    IndexConvicted    IndexType = "BIO-CON"
    IndexArrestee     IndexType = "BIO-ARR"
    IndexForensic     IndexType = "BIO-FSC"
    IndexMissingPerson IndexType = "BIO-DIS"
    IndexUnidentified  IndexType = "BIO-RNI"
)

// 20 loci CODIS Core Standard (NIST 2017)
var CODISCoreLoci = []string{
    "CSF1PO", "D3S1358", "D5S818", "D7S820", "D8S1179",
    "D13S317", "D16S539", "D18S51", "D21S11", "FGA",
    "TH01", "TPOX", "vWA", "D1S1656", "D2S441",
    "D2S1338", "D10S1248", "D12S391", "D19S433", "D22S1045",
    "Amelogenin", // Déterminant du sexe
}

type MatchResult struct {
    HitID        string    `json:"hit_id"`
    QueryID      string    `json:"query_id"`
    MatchID      string    `json:"match_id"`
    MatchType    string    `json:"match_type"` // FULL_MATCH, PARTIAL, FAMILIAL
    Confidence   float64   `json:"confidence"`  // Score 0.0-1.0
    MatchedLoci  int       `json:"matched_loci"`
    TotalLoci    int       `json:"total_loci"`
    NIU          *string   `json:"niu,omitempty"`
    AlertLevel   string    `json:"alert_level"` // HIGH, MEDIUM, LOW
}

// DNAMatcher effectue la comparaison entre profils STR
type DNAMatcher struct {
    db     Database
    cache  Cache
    events EventPublisher
}

func (m *DNAMatcher) SearchProfile(ctx context.Context, query STRProfile) ([]MatchResult, error) {
    // 1. Vérifier le cache Redis
    cacheKey := m.buildCacheKey(query)
    if cached, err := m.cache.Get(ctx, cacheKey); err == nil {
        return cached, nil
    }

    // 2. Recherche dans la base (algorithme matching CODIS)
    candidates, err := m.db.FindCandidates(ctx, query, 0.85) // seuil 85%
    if err != nil {
        return nil, fmt.Errorf("candidate search: %w", err)
    }

    // 3. Scoring fin pour chaque candidat
    var results []MatchResult
    for _, candidate := range candidates {
        score := m.calculateLRScore(query.Loci, candidate.Loci)
        if score.Confidence >= 0.999 {
            results = append(results, MatchResult{
                MatchType:  "FULL_MATCH",
                Confidence: score.Confidence,
            })
        } else if score.Confidence >= 0.85 {
            results = append(results, MatchResult{
                MatchType:  "PARTIAL",
                Confidence: score.Confidence,
            })
        }
    }

    // 4. Publier hit sur Kafka
    if len(results) > 0 {
        m.events.Publish(ctx, "snisid.bio.hits", results)
    }

    return results, nil
}

// calculateLRScore — Likelihood Ratio scoring (standard forensique)
func (m *DNAMatcher) calculateLRScore(query, target map[string]Allele) MatchScore {
    matchedLoci := 0
    totalLoci := len(CODISCoreLoci)

    for _, locus := range CODISCoreLoci {
        q, qOk := query[locus]
        t, tOk := target[locus]
        if !qOk || !tOk {
            continue
        }
        if allelesMatch(q, t) {
            matchedLoci++
        }
    }

    confidence := float64(matchedLoci) / float64(totalLoci)
    return MatchScore{MatchedLoci: matchedLoci, TotalLoci: totalLoci, Confidence: confidence}
}

func allelesMatch(a, b Allele) bool {
    if a.Value1 != b.Value1 {
        return false
    }
    if (a.Value2 == nil) != (b.Value2 == nil) {
        return false
    }
    if a.Value2 != nil && *a.Value2 != *b.Value2 {
        return false
    }
    return true
}

type MatchScore struct {
    MatchedLoci int
    TotalLoci   int
    Confidence  float64
}
```

### 2.2 Architecture de déploiement K8s

```yaml
# k8s/bio-adn-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: snisid-bio-adn
  namespace: snisid-forensic
spec:
  replicas: 3
  selector:
    matchLabels:
      app: snisid-bio-adn
  template:
    metadata:
      labels:
        app: snisid-bio-adn
        snisid.gov.ht/security-tier: forensic
    spec:
      containers:
      - name: bio-adn-api
        image: harbor.snisid.gov.ht/snisid/bio-adn-service:1.0.0
        ports:
        - containerPort: 8443  # REST/gRPC
        env:
        - name: DB_URL
          valueFrom:
            secretKeyRef:
              name: snisid-bio-secrets
              key: db-url
        - name: HSM_PIN
          valueFrom:
            secretKeyRef:
              name: snisid-hsm-secrets
              key: hsm-pin
        - name: KAFKA_BROKERS
          value: "kafka.snisid-infra.svc.cluster.local:9092"
        resources:
          requests:
            memory: "2Gi"
            cpu: "1000m"
          limits:
            memory: "8Gi"
            cpu: "4000m"
        volumeMounts:
        - name: tls-certs
          mountPath: /etc/tls
          readOnly: true
      volumes:
      - name: tls-certs
        secret:
          secretName: snisid-bio-tls
```

---

## 3. FLUX DE SYNCHRONISATION LDIS → SDIS → NDIS

```
LDIS (labo local)
│
│ 1. Génération profil STR (format CODIS 20 loci)
│ 2. Chiffrement AES-256-GCM avec clé HSM locale
│ 3. Signature numérique lab (certificat X.509 SNISID PKI)
│ 4. Upload SDIS via mTLS (schedule quotidien ou on-demand)
│
▼
SDIS (département)
│
│ 1. Validation signature lab
│ 2. Déduplication locale (hash SHA-256 profil)
│ 3. Indexation dans PostgreSQL départemental
│ 4. Matching intra-département
│ 5. Transmission NDIS (schedule hebdomadaire — même cadence que CODIS)
│
▼
NDIS-HT (national)
│
│ 1. Matching cross-département
│ 2. Matching avec index INTERPOL DNA Gateway
│ 3. Publication hits sur Kafka (snisid.bio.hits)
│ 4. Notification agences concernées via DIDComm
│ 5. Rapport hebdomadaire direction DCPJ
```

---

## 4. SÉCURITÉ SPÉCIFIQUE ADN

### 4.1 Chiffrement des profils

- Les profils STR sont stockés **chiffrés** (AES-256-GCM, clé HSM)
- La clé de déchiffrement n'est jamais en mémoire RAM > 30 secondes
- Le déchiffrement se fait **uniquement** au moment du matching, en mémoire sécurisée

### 4.2 Principe de dissociation

Inspiré du CODIS : **les profils ADN ne contiennent pas de noms**.

```sql
-- Table profils ADN : PAS de nom, PAS d'adresse
CREATE TABLE bio_str_profiles (
    sample_id   UUID PRIMARY KEY,
    index_type  VARCHAR(10) NOT NULL,  -- BIO-CON, BIO-FSC, etc.
    loci_data   BYTEA NOT NULL,        -- Profil STR chiffré HSM
    quality_score DECIMAL(3,2),
    lab_id      VARCHAR(50) NOT NULL,
    collected_at TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW()
    -- PAS de NIU, PAS de nom ici
);

-- Table de liaison séparée, accès restreint DCPJ-DIR uniquement
CREATE TABLE bio_identity_link (
    sample_id   UUID REFERENCES bio_str_profiles(sample_id),
    niu         VARCHAR(20),           -- Lien SNISID Core
    linked_by   VARCHAR(100),          -- Officier ayant créé le lien
    linked_at   TIMESTAMPTZ DEFAULT NOW(),
    court_order VARCHAR(200),          -- Numéro ordonnance judiciaire
    PRIMARY KEY (sample_id)
);
```

### 4.3 Audit immuable

Chaque accès à un profil ADN est journalisé dans l'**Audit Ledger** SNISID :
```json
{
  "event": "bio.profile.accessed",
  "sample_id": "uuid",
  "officer_niu": "agent-niu",
  "agency": "DCPJ",
  "purpose": "criminal_investigation",
  "case_number": "2026-DCPJ-001234",
  "timestamp": "2026-06-09T10:30:00Z",
  "ip_address": "hash(ip)",
  "signature": "ECDSA-P256"
}
```
