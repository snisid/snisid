# Architecture: National Data Sovereignty

## 1. Physical Sovereignty
- **Data Locality**: All SNISID primary and replica nodes must be hosted on hardware owned or strictly leased by the sovereign government within national territory.
- **Air-Gap Capability**: Infrastructure must be capable of operating without external internet connectivity for a minimum of 30 days.

## 2. Logical Segmentation (Agency Isolation)
The system uses **Namespace-Level Isolation** for multi-agency data:
- `snisid-anh`: Restricted civil registry data.
- `snisid-dgi`: Restricted fiscal and tax data.
- `snisid-dcpj`: Restricted security and criminal intelligence.

Communication between segments is permitted *only* via the **Signed Event Bus** (Kafka), with all messages validated by the **Policy Intelligence Optimizer**.

## 3. Cryptographic Backbone
- **Identity**: SPIFFE/SPIRE for workload identity and mTLS.
- **Signing**: All agency events are signed using Ed25519 hardware-backed keys (HSM).
- **At-Rest**: LUKS-level disk encryption + application-level field encryption for sensitive biometrics.

## 4. Sovereignty Stack (Non-Cloud Dependent)
- **Database**: PostgreSQL (Relational) + Neo4j (Graph).
- **Streaming**: Apache Kafka.
- **Orchestration**: K3s/Kubernetes (Sovereign Cluster).
- **Object Store**: MinIO (S3-Compatible).
