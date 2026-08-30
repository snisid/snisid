# SNISID Master Citizen Registry Architecture
## CQRS Read Model & Fast-Search Microservice

This document outlines the architectural design for the **Master Citizen Registry Service**. Operating as the "Query" component in the system's Command Query Responsibility Segregation (CQRS) architecture, this microservice is hyper-optimized for lightning-fast reads, fuzzy searching, and inter-agency identity verification.

---

## 1. Master Record & Schema Architecture

The Citizen Registry does not possess the authority to *alter* citizen data. It builds a **Materialized View** (a read-only snapshot) by consuming immutable events from the Kafka Event Store (populated by the Identity Service).

### Document Schema (Elasticsearch / OpenSearch JSON)
The registry utilizes a NoSQL document structure optimized for full-text and fuzzy search against the Haitian populace.

```json
{
  "mappings": {
    "properties": {
      "niu": { "type": "keyword" },
      "demographics": {
        "properties": {
          "nom": { "type": "text", "analyzer": "haitian_creole_analyzer" },
          "prenom": { "type": "text", "analyzer": "haitian_creole_analyzer" },
          "date_naissance": { "type": "date" },
          "lieu_naissance": { "type": "keyword" },
          "sexe": { "type": "keyword" }
        }
      },
      "biometric_references": {
        "properties": {
          "abis_gallery_id": { "type": "keyword" },
          "fingerprint_enrolled": { "type": "boolean" },
          "iris_enrolled": { "type": "boolean" },
          "facial_enrolled": { "type": "boolean" }
        }
      },
      "status": {
        "type": "keyword"  // e.g., "ACTIVE", "DECEASED", "SUSPENDED"
      },
      "version": { "type": "long" } // Used for optimistic locking & event ordering
    }
  }
}
```

---

## 2. API & Event Contracts

### Consume: Event Contract (AsyncAPI)
The Registry consumes the `CitizenRegistered` and `IdentityUpdated` events from Kafka to update its Elasticsearch indices.

```yaml
channels:
  identity.events:
    subscribe:
      summary: Consume identity mutation events to rebuild the read model.
      message:
        $ref: '#/components/messages/IdentityMutatedEvent'
```

### Expose: API Contract (OpenAPI)
The API exposes highly restricted REST endpoints to the API Gateway for authorized agencies to query.

```yaml
paths:
  /v1/registry/search:
    get:
      summary: Fuzzy search citizen registry
      parameters:
        - name: query
          in: query
          required: true
          schema:
            type: string
            example: "Jean Baptiste"
      responses:
        '200':
          description: Array of matching minimal citizen profiles.
```

---

## 3. Resilience, Encryption & Replication

### High Availability (HA) Replication
- **Cross-Region Clusters:** The OpenSearch cluster is deployed across both Port-au-Prince (DC1) and Cap-Haïtien (DC2) using Cross-Cluster Replication (CCR). If DC1 burns down, queries instantly failover to DC2.
- **Index Sharding:** Indices are sharded across multiple data nodes to distribute load, with at least 2 replica shards per primary shard.

### Encryption (At Rest & In Transit)
- **At Rest:** The underlying Kubernetes persistent volumes (provided by Ceph) are encrypted via **LUKS** combined with physical server Hardware TPMs. Additionally, Elasticsearch is configured with native TLS index encryption.
- **In Transit:** Every query to the Registry Service passes through the Istio Envoy proxy, guaranteeing **STRICT mTLS**.

### Deduplication & Fraud Prevention
- The Registry exposes the aggregated `status` field. If the underlying Identity Service suspends a citizen (e.g., due to an ABIS deduplication hit catching a citizen attempting to register twice under different names), the Registry instantly reflects `status: SUSPENDED`, blocking that NIU from being used at any government agency (like DGI or DGIE) instantly.

---

## 4. Architecture Diagrams (Mermaid)

### 1. CQRS Synchronization & Flow Diagram
This flowchart demonstrates how the Citizen Registry builds its data model autonomously without querying the primary CockroachDB database.

```mermaid
graph TD
    classDef source fill:#ffebee,stroke:#c62828,stroke-width:2px;
    classDef broker fill:#fff3e0,stroke:#e65100,stroke-width:2px;
    classDef registry fill:#e3f2fd,stroke:#1565c0,stroke-width:2px;
    classDef client fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;

    IdSvc[Identity Service <br/> Command/Write Path]:::source
    Kafka[(Kafka Topic: <br/> snisid.identity.events)]:::broker
    
    subgraph Read_Model [Master Citizen Registry (Read Path)]
        RegApp[Registry Microservice <br/> Golang]:::registry
        ES[(OpenSearch Cluster)]:::registry
        RegApp -->|Bulk Upsert Documents| ES
    end

    GW[SNISID API Gateway]:::client

    IdSvc -->|Publish Event| Kafka
    Kafka -->|Consume (Consumer Group: registry)| RegApp
    
    GW -->|GET /v1/registry/search?q=...| RegApp
    RegApp -->|Fuzzy Query| ES
```

### 2. Multi-Region OpenSearch Deployment Model
To guarantee that police and hospitals can always query the registry, it is heavily replicated.

```mermaid
graph LR
    classDef node fill:#e3f2fd,stroke:#1565c0,stroke-width:2px;
    classDef master fill:#fff3e0,stroke:#e65100,stroke-width:2px;

    subgraph DC1 [Port-au-Prince Cluster]
        M1[Master Node 1]:::master
        D1[Data Node A <br/> Shard 0 (Primary)]:::node
        D2[Data Node B <br/> Shard 1 (Replica)]:::node
    end

    subgraph DC2 [Cap-Haïtien Cluster]
        M2[Master Node 2]:::master
        D3[Data Node C <br/> Shard 0 (Replica)]:::node
        D4[Data Node D <br/> Shard 1 (Primary)]:::node
    end

    M1 <.->|Cross-Cluster Sync| M2
    D1 <==>|Async Index Replication| D3
    D2 <==>|Async Index Replication| D4
```

---
*Prepared by the SNISID Cloud Infrastructure & Resilience Board.*
