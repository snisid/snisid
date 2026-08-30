# SNI-SIDE

## SNISID National Intelligence, Security, Investigation and Sovereign Data Ecosystem

Extension souveraine du SNISID (Système National d'Identification Sécurisée et d'Interopérabilité Digitale) pour l'intelligence nationale, la sécurité, l'investigation criminelle, et la gestion souveraine des données.

### Architecture

```
SNI-SIDE/
├── README.md                           ← Ce fichier
├── SNI-SIDE-ARCHITECTURE.md            ← Architecture globale + Mermaid diagrams
├── ROADMAP-2026-2035.md                ← Plan d'implémentation décennal
│
├── schemas/                            ← Schémas de base de données
│   ├── 01-ncid.sql                     ← National Criminal Intelligence Database
│   ├── 02-hn-ngi.sql                   ← National Biometric Database (HN-NGI)
│   ├── 03-hn-codis.sql                 ← Combined DNA Index System (CODIS)
│   ├── 04-missing-persons.sql          ← Missing Persons Database
│   ├── 05-vehicle-intelligence.sql     ← Vehicle Intelligence Database
│   ├── 06-alpr.sql                     ← National ALPR Database (CockroachDB)
│   ├── 07-firearms.sql                 ← Firearms Intelligence Database
│   ├── 08-border-intelligence.sql      ← Border Intelligence Database
│   ├── 09-counter-narcotics.sql        ← Counter Narcotics Database
│   ├── 10-financial-crime.sql          ← Financial Crime Database
│   ├── 11-cybercrime.sql               ← Cybercrime Intelligence Database
│   ├── 12-watchlist.sql                ← National Watchlist Database
│   ├── 13-document-fraud.sql           ← Document Fraud Database
│   ├── 14-geoint.sql                   ← GEOINT Database (PostGIS)
│   ├── 15-digital-evidence.sql         ← Digital Evidence Repository (CockroachDB)
│   └── neo4j-graph.cypher             ← Neo4j Sovereign Intelligence Graph
│
├── api/
│   └── openapi-sniside.yaml           ← REST API spec (OpenAPI 3.0)
│
├── events/
│   ├── kafka-topics.yaml              ← Topics Kafka (100+ topics)
│   └── avro/
│       ├── wanted_person.v1.avsc      ← Wanted person event schema
│       ├── biometric_match.v1.avsc    ← Biometric match event schema
│       ├── alpr_read.v1.avsc          ← ALPR read event schema
│       ├── watchlist_match.v1.avsc    ← Watchlist match event schema
│       └── fusion_alert.v1.avsc       ← AI Fusion alert event schema
│
├── ai/
│   └── ai-fusion-center.py           ← National AI Fusion Center models
│                                      (Fraud GNN, ArcFace, DNA AI, AML,
│                                       Cyber Threat, Predictive Crime,
│                                       Deepfake Detection, GraphRAG)
│
├── k8s/
│   └── sniside-deployment.yaml       ← Kubernetes manifests (PostgreSQL 16,
│                                       CockroachDB, Neo4j Enterprise, Milvus,
│                                       Kafka/Strimzi, ClickHouse, MinIO, Redis,
│                                       microservices, HPA, Istio, Cilium)
│
├── security/
│   └── zero-trust-policies.yaml      ← OPA, Istio, Cilium, Vault, SPIFFE/SPIRE
│
├── advanced-systems/
│   └── README.md                     ← Fusion Center, RTCC, NTIP, Child Protection,
│                                       DVI, Data Lake, Counter Terrorism, Maritime,
│                                       Aviation, Critical Infrastructure, etc.
│
└── integration/
    └── snisi-integration-map.md      ← Mappage complet SNI-SIDE → SNISID Core
```

### 15 Bases Nationales

| # | Base | Moteur | Type de Données | Volume Estimé |
|:--|:--|:--:|:--|:--:|
| 1 | NCID | PostgreSQL 16 | Personnes recherchées, mandats, cas, gangs | 500K personnes |
| 2 | HN-NGI | PostgreSQL + Milvus | Empreintes, visage, iris, voix | 10M templates |
| 3 | HN-CODIS | PostgreSQL 16 | ADN criminel, scène de crime, familial | 200K profils |
| 4 | Missing Persons | PostgreSQL 16 | Enfants/adultes disparus, kidnappings | 50K cas |
| 5 | Vehicle Intelligence | PostgreSQL 16 | Véhicules, propriétaires, historique | 2M véhicules |
| 6 | National ALPR | CockroachDB | Lectures plaques, caméras, routes | 200M lectures/mois |
| 7 | Firearms | PostgreSQL 16 | Armes, balistique, scènes de crime | 500K armes |
| 8 | Border Intelligence | PostgreSQL 16 | Entrées/sorties, visas, déportations | 50M passages |
| 9 | Counter Narcotics | PostgreSQL 16 | Cartels, routes, saisies | 100K événements |
| 10 | Financial Crime | PostgreSQL 16 | Transactions, AML, PEP | 10M transactions |
| 11 | Cybercrime | PostgreSQL 16 | IOC, malware, wallets | 2M IOCs |
| 12 | Watchlist | PostgreSQL 16 | Personnes, véhicules, documents | 200K entrées |
| 13 | Document Fraud | PostgreSQL 16 | Passeports, CIN, permis fraudés | 1M documents |
| 14 | GEOINT | PostgreSQL + PostGIS | Drones, satellites, hotspots | 500K couches |
| 15 | Evidence | CockroachDB + MinIO | Photos, vidéos, audio, forensique | 500TB |

### Technologies

| Technologie | Usage |
|:--|:--|
| PostgreSQL 16 | Bases relationnelles principales |
| CockroachDB | Bases géo-distribuées (ALPR, Evidence) |
| Neo4j Enterprise 5.x | Graph Intelligence Souverain |
| Milvus | Vector DB pour biométrie |
| ClickHouse | Analytics temps réel |
| Kafka (Strimzi) | Event Streaming |
| MinIO | Object Store souverain |
| Redis | Cache Search Engine |
| Kubernetes (RKE2) | Orchestration |
| Istio | Service Mesh (mTLS) |
| Cilium | eBPF Network Security |
| OPA | Policy Engine |
| Vault | Secrets Management |
| SPIFFE/SPIRE | Workload Identity |
| PyTorch + PyG | AI/ML Models |
| MLflow | Model Registry |
| OpenTelemetry | Observabilité |

### Déploiement

```bash
# 1. Namespace
kubectl create namespace sniside

# 2. Databases
kubectl apply -f k8s/sniside-deployment.yaml

# 3. Kafka Topics
kubectl apply -f events/kafka-topics.yaml

# 4. Security Policies
kubectl apply -f security/zero-trust-policies.yaml

# 5. APIs
kubectl apply -f api/openapi-sniside.yaml
```

### Intégration SNISID

Voir [integration/snisi-integration-map.md](integration/snisi-integration-map.md) pour le mapping complet avec SNISID Core, Identity Registry, ABIS, Fraud Engine, Kafka, Neo4j, MLOps, SOC, Cyber Defense, Interop, et API Ecosystem.

### Roadmap

Voir [ROADMAP-2026-2035.md](ROADMAP-2026-2035.md) — 8 phases sur 10 ans, budget total $58.4M.
