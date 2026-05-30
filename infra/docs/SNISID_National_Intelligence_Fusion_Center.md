# SNISID: National Intelligence Fusion Center (NIFC)

The National Intelligence Fusion Center (NIFC) acts as the "Data Hub" for SNISID, integrating multi-agency telemetry to build a comprehensive, 360-degree context for every national identity.

---

## 1. Multi-Source Ingestion & Handshake

NIFC integrates data from external "Sovereign Partners" using a high-assurance **Ingestion Gateway**.

### 1.1. Ingestion Sources
- **Telecom Providers**: Real-time geolocation, SIM swap alerts, and call metadata.
- **Banking Infrastructure**: High-value transaction alerts and credit/fraud scores.
- **Border Control**: Entry/exit records and passport verification events.
- **Law Enforcement**: Wanted person alerts and criminal history updates.

### 1.2. Secure Handshake Protocol
All external data providers must connect via:
- **Mutual TLS (mTLS)**: Certificate-based identity verification for the partner agency.
- **SPIRE SVID**: Each partner is assigned a unique workload identity within the SNISID mesh.
- **Data Sovereignty Contract**: A cryptographically signed agreement that defines which fields can be shared and for what purpose.

---

## 2. Unified Identity Context & Data Fusion

NIFC uses the **National Graph (Neo4j)** to fuse disparate data points into a single "Citizen Context."

- **Identity Stitching**: Automatically linking a Telecom SIM, a Bank Account, and a National ID number to the same physical person based on high-confidence identifiers.
- **The Golden Record**: A real-time, aggregated view of an identity's current state (e.g., *"Citizen X: Currently at Location Y (via Telecom), just performed Transaction Z (via Bank), is currently 'Verified' (via SNISID Core)"*).
- **Enrichment Pipelines**: Flink jobs that aggregate these signals to calculate the **ISTS (Internal Service Trust Score)** for the identity.

---

## 3. Data Sovereignty & PII Protection

- **Tenant Isolation**: External data is stored in **Sovereign Object Storage** with strict tenant-level encryption. Data from the "Banking" tenant cannot be seen by the "Telecom" tenant without a legal warrant.
- **Field-Level Access Control**: Access to specific fused fields (e.g., precise geolocation) requires a signed **Access Token** specifying the investigation ID and duration.
- **Automated Anonymization**: All non-critical fused data is automatically masked or anonymized after 90 days to comply with data minimization laws.

---

## 4. NIFC Operational Workflows

- **Intelligence Alerting**: A high-risk event from a Bank (e.g., "Account Compromise") instantly triggers an "Identity Quarantine" in the SNISID Core to prevent document forgery using that account.
- **National Security Queries**: Analysts can perform "Cross-Domain Queries" (e.g., *"Show all identities that performed a login from IP X and have a recent SIM swap alert"*).

---

## 5. Scalability & Resilience

- **Distributed Ingestion Gateways**: Deployed at regional peering points to minimize latency for partner agencies.
- **Exactly-Once Delivery**: Kafka-based pipelines ensure that no intelligence event is lost or processed twice during the fusion process.
