# SNISID: National Graph Intelligence Architecture

The Sovereign Graph Engine (Neo4j) provides a multi-dimensional view of relationships across the nation, enabling the detection of fraud rings, terror cells, and complex identity clusters.

---

## 1. Graph Schema & Node Relationship Model

SNISID models the nation as a **Hyper-Connected Property Graph**.

### 1.1. Core Node Types
- **Citizen**: `(name, nid_hash, trust_score, clearance)`
- **Device**: `(fingerprint, type, os, last_known_location)`
- **Location**: `(address, coordinates, jurisdiction)`
- **Agency**: `(name, code, type)`
- **Event**: `(type, timestamp, result)`

### 1.2. Key Relationships
- `(Citizen)-[:OWNS]->(Device)`
- `(Citizen)-[:LIVES_AT]->(Location)`
- `(Citizen)-[:WORKS_FOR]->(Agency)`
- `(Citizen)-[:AUTHENTICATED_BY]->(Agency)`
- `(Citizen)-[:LINKED_TO]->(Citizen)` (Family, Business, or Fraud correlation)

---

## 2. Distributed Storage & Scaling Architecture

To handle millions of nodes and billions of relationships:

- **Neo4j Causal Clustering**:
  - **Core Servers (3+)**: Handle write transactions and maintain the Raft-based cluster state.
  - **Read Replicas (Scalable)**: Handle high-volume read queries for the Fraud Engine and SOC.
- **Fabric Architecture**: Sharding the national graph across regional clusters (e.g., Ouest vs. Nord) while allowing cross-regional queries for national security investigations.

---

## 3. Real-Time Graph Intelligence Updates

The graph is updated in real-time by the **Flink Sync Sink**.

1.  **Event Ingestion**: `identity.updated` arrives in Kafka.
2.  **Cypher Update**: Flink executes an idempotent Cypher query:
    ```cypher
    MERGE (c:Citizen {nid: $nid})
    SET c.last_seen = $timestamp
    MERGE (l:Location {id: $loc_id})
    MERGE (c)-[r:SEEN_AT]->(l)
    SET r.timestamp = $timestamp
    ```
3.  **Trigger**: New relationships can trigger a **Graph Data Science (GDS)** algorithm (e.g., "PageRank" or "Community Detection") to re-evaluate the risk score of the entire cluster.

---

## 4. Query Optimization Strategy

- **Index Management**: Mandatory indexes on `nid_hash`, `device_fingerprint`, and `event_id`.
- **Query Depth Limiting**: To prevent "Supernode" traversal from freezing the cluster, all real-time queries are capped at 3-5 hops.
- **APOC Procedures**: Utilizing optimized procedures for complex path-finding and temporal graph analysis.

---

## 5. Graph Intelligence Capabilities

- **Fraud Ring Detection**: Identifying "Synthetic Identity" clusters where multiple citizens share the same phone number, address, or device fingerprint.
- **Lateral Movement Analysis**: Visualizing how an attacker moved from a compromised employee account to sensitive agency resources.
- **Chain of Trust**: Verifying the relationship path from a Root Authority to a specific identity verification event.

---

## 6. Operational Deployment Model

- **Neo4j on Kubernetes**: Deployed via Helm with **Persistent Volume Claims (PVC)** on high-speed NVMe storage.
- **Backup & Restore**: Integration with the **Sovereign Backup System** for daily consistent snapshots of the entire graph state.
- **Security**: 
  - **Role-Based Access**: SOC analysts see the graph; Developers only see anonymized/masked node properties.
  - **Encryption**: TLS 1.3 for all cluster communication and AES-256 for data-at-rest.
