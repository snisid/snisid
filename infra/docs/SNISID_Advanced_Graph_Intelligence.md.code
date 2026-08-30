# SNISID: Advanced Graph Intelligence Architecture

The Advanced Graph Intelligence layer leverages Neo4j Graph Data Science (GDS) to perform national-scale relationship analysis, uncovering hidden patterns and community structures within the identity ecosystem.

---

## 1. Real-Time Graph Ingestion & Sync (Prompt 122)

The graph is a "Live Mirror" of the national event mesh.

- **Ingestion Pipeline**: Kafka -> Flink GDS Sink -> Neo4j.
- **Normalization**: High-velocity events are pre-aggregated in Flink (e.g., "100 login events between Device X and Citizen Y") before being emitted as a single `LAST_SEEN` relationship update in the graph.
- **Consistency**: Distributed transactions ensure that a node update in Neo4j is atomically linked to its originating Kafka offset.

---

## 2. Fraud Ring & Community Detection (Prompts 123, 125)

SNISID identifies malicious clusters through **Structural Relationship Analysis**.

### 2.1. Community Detection Algorithms
- **Louvain Modularity**: Identifies dense clusters of identities that share suspicious commonalities (e.g., shared addresses, phone numbers, or device fingerprints).
- **Weakly Connected Components (WCC)**: Rapidly identifies isolated sub-graphs that may represent autonomous fraud networks.

### 2.2. Fraud Ring Detection Workflow
1.  **Trigger**: A new "High-Risk" event (e.g., a biometric mismatch) occurs.
2.  **Expansion**: GDS performs a 3-hop expansion to identify all linked entities.
3.  **Scoring**: The engine calculates a "Cluster Risk Score" based on the density of suspicious relationships.
4.  **Action**: If the cluster exceeds a threshold, the entire group is flagged for **Sovereign Review**.

---

## 3. Risk Propagation & Centrality (Prompt 127)

Threats are not isolated; they propagate through relationships.

- **PageRank (Risk Centrality)**: Identifies "Influence Hubs" in the graph. An identity linked to 10 high-risk individuals receives a higher risk score than one linked to 10 low-risk individuals.
- **Personalized PageRank**: Calculates the probability of risk spreading from a "Patient Zero" (compromised account) to other connected nodes.

---

## 4. Graph Machine Learning (GDS ML) (Prompt 126)

- **Node Embeddings (Node2Vec / FastRP)**: Transforms the complex graph structure of a citizen's relationships into a multi-dimensional vector (Embedding).
- **Classification**: These embeddings are fed into the **AI Inference Pipeline** to classify an identity as "Normal," "Synthetic," or "Compromised" based on their position in the national graph.

---

## 5. Distributed Graph Computation & Scaling (Prompt 128)

- **Fabric Sharding**: The national graph is partitioned by jurisdiction (e.g., Ouest, Nord, Sud). Fabric allows for "Global Federated Queries" across regional shards for national security investigations.
- **Read-Scaling**: Deployment of **Read Replicas** in a causal cluster to handle the massive query load from the real-time AI inference engines.

---

## 6. Security, Governance & AI Integration

- **Security Hardening**:
  - **Property-Level Security**: Only investigators with "Level 4" clearance can see biometric relationship properties.
  - **Cypher Guard**: An OPA-backed proxy that inspects and filters Cypher queries to prevent unauthorized data exfiltration.
- **AI Integration**: The Graph Engine provides the "Structural Context" to the **Explainable AI (XAI)** layer, allowing the system to say: *"I flagged this user because they are 2 hops away from a known fraud ring."*
