# SNISID National Platform: Full End-to-End System Architecture

This document defines the production-grade architecture for **SNISID**, a sovereign-scale cyber intelligence and identity platform.

## 📐 System Overview

SNISID is built on a **Modular Microservice Architecture** utilizing a **Zero Trust** security model. The system is designed to scale from a single-node field deployment to a national multi-cluster infrastructure.

```mermaid
graph TD
    subgraph "External World"
        Agencies["Agency Clients / Field Units"]
        SIEM["External SIEM / Interop"]
    end

    subgraph "Trust Zone: Ingress (WAF/WAP)"
        GW["API Gateway (Go/Gin)"]
    }

    subgraph "Trust Zone: Application (Istio Mesh)"
        ID["Identity API"]
        FR["Fraud Engine"]
        SOC["SOC API"]
        SW["AI Agent Swarm"]
        OPS["AIOps Engine"]
        TW["Digital Twin Sync"]
    end

    subgraph "Trust Zone: Intelligence (GPU/HBM)"
        AF["ArcFace Biometrics (Python)"]
        GNN["GNN Fraud Analysis (Python)"]
        DP["Deepfake Detection (Python)"]
    end

    subgraph "Data Tier (Encrypted)"
        PG[(PostgreSQL)]
        NEO[(Neo4j Graph)]
        KAF{Kafka Event Bus}
        RED[(Redis Streams)]
    end

    subgraph "SOC Command Center"
        UI["React Dashboard (3D/Graph)"]
        LGTM["Loki/Grafana/Tempo/Mimir"]
    end

    Agencies --> GW
    GW --> ID
    ID --> KAF
    KAF --> FR
    KAF --> AF
    KAF --> GNN
    FR --> NEO
    ID --> PG
    FR --> RED
    RED --> UI
    SW --> KAF
    TW --> NEO
    ID -.-> LGTM
```

---

## 🏗️ Layered Architecture Breakdown

### 1. Infrastructure Layer
- **Orchestration**: Kubernetes (K8s) as the universal abstraction.
- **Scaling**: Horizontal Pod Autoscaling (HPA) driven by AIOps predictive metrics.
- **Resilience**: Multi-AZ/Multi-Region deployments with Kafka MirrorMaker 2 for state replication.

### 2. Networking & Zero Trust Layer
- **Service Mesh**: Istio provides transparent mTLS, fine-grained `AuthorizationPolicies`, and L7 observability.
- **Workload Identity**: SPIRE manages X.509 SVID issuance and rotation across clusters.
- **Encryption**: AES-256 for data-at-rest; TLS 1.3 for data-in-transit.

### 3. Core Subsystems
- **Identity Management**: Handles lifecycle, enrollment, and verification of national identities.
- **Biometric Intelligence**: Real-time facial recognition and deepfake analysis using ResNet-50 and CNN backbones.
- **Graph Intelligence**: Analyzes identity relationships and synthetic clusters using Neo4j and GAT (Graph Attention Networks).
- **Autonomous Swarm**: A collaborative layer of AI agents (Threat Hunter, Fraud Investigator) orchestrating complex security investigations.

### 4. Data Flow Overview
1.  **Ingestion**: Identities are submitted via the API Gateway.
2.  **Enrichment**: The Identity API triggers biometric and fraud analysis via the Kafka backbone.
3.  **Intelligence**: AI services process embeddings and graph relationships asynchronously.
4.  **Consensus**: The Fraud Engine correlates AI results and updates the Neo4j graph.
5.  **Persistence**: Final identity state is committed to PostgreSQL with full audit trails.
6.  **Visualization**: Real-time alerts are pushed to the React Dashboard via Redis Streams and WebSockets.

---

## 📡 Communication Model
- **Synchronous (REST/gRPC)**: Used for critical, immediate requests (e.g., UI to API Gateway).
- **Asynchronous (Event-Driven)**: Used for the majority of inter-service traffic via Kafka to ensure decoupled scalability and fault tolerance.
- **Streaming**: Used for sub-millisecond alert propagation via Redis Streams and WebSockets.

## 🔒 Security Boundaries
- **DMZ**: API Gateway acts as the sole entry point, performing initial JWT validation and WAF filtering.
- **Internal Mesh**: All service-to-service communication is encrypted and requires a valid SPIFFE ID.
- **Data Vault**: Databases are isolated in a restricted trust zone, accessible only by authenticated service principals.
- **Air-Gap Capability**: The entire stack can be deployed in a fully disconnected environment for high-security field operations.
