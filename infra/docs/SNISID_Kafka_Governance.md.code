# SNISID: Kafka Topic Governance Framework

To ensure consistency, security, and scalability across the national event mesh, SNISID enforces a strict governance model for all Kafka topics.

---

## 1. Topic Naming Convention

SNISID uses a hierarchical, dot-notated naming standard:
`<agency>.<domain>.<entity>.<version>.<status>`

| Segment | Example | Description |
| :--- | :--- | :--- |
| **Agency** | `national` | The owning agency (e.g., `national`, `police`, `border`). |
| **Domain** | `identity` | The functional area (e.g., `identity`, `fraud`, `audit`). |
| **Entity** | `citizen` | The primary object (e.g., `citizen`, `match`, `login`). |
| **Version** | `v1` | Major schema version (to allow parallel evolution). |
| **Status** | `live` | Operational status (`live`, `retry`, `dlq`, `archive`). |

**Example**: `national.identity.citizen.v1.live`

---

## 2. Topic Classification & Retention Matrix

Topics are classified based on the sensitivity of the data and the required historical depth.

| Class | Retention | PII Status | Storage Tier |
| :--- | :--- | :--- | :--- |
| **Operational** | 7 Days | No | Hot (SSD) |
| **Audit/Security** | Forever | No | Cold (S3/WORM) |
| **Transactional** | 30 Days | Masked | Hot + Cold |
| **PII/Biometric** | 1 Year | **Encrypted** | Cold + Air-gap |

---

## 3. PII Event Isolation & Security

- **PII Topics**: Topics containing sensitive identifiers (National ID, Biometrics) are isolated in a dedicated **Encrypted Prefix**. 
- **E2EE**: Payloads in PII topics are encrypted by the producer using the **Vault Transit Engine** before ingestion. 
- **Scrubbing**: Operational topics (`.live`) are strictly monitored for PII leakage. If PII is detected, the topic is automatically purged and the service is quarantined.

---

## 4. Schema Evolution Rules

SNISID utilizes a centralized **Schema Registry** (Protobuf).
- **Compatibility**: Enforced `BACKWARD_TRANSITIVE`. This ensures that any consumer can read any previous version of a message.
- **Breaking Changes**: Requires a new topic name (increment the `<version>` segment).
- **Metadata**: Every schema must include standard metadata fields: `correlation_id`, `timestamp`, `originating_service`, and `trace_context`.

---

## 5. Multi-Agency Segregation (ACLs)

Access is governed by **Least Privilege** and **Sovereign Isolation**.
- **Agency Scoping**: A service from the `Police` agency can only `READ` from `national.*` topics and `READ/WRITE` to `police.*` topics.
- **RBAC for Kafka**:
  - `Event Producer`: Write-only to specific topics.
  - `Event Consumer`: Read-only + Offset Management.
  - `Security Auditor`: Read-only to all `audit.*` topics.

---

## 6. Operational Governance Workflows

### 6.1. Topic Creation Workflow
1.  **Request**: Developer requests a new topic via a GitOps (Yaml) definition.
2.  **Linting**: Automated CI check for naming convention and partition alignment.
3.  **Approval**: Security review for PII classification.
4.  **Deployment**: Strimzi Operator creates the topic and applies ACLs.

### 6.2. Schema Promotion Workflow
- Schemas are first registered in the `Dev` registry.
- Promotion to `Prod` requires passing automated compatibility tests against all active consumers.
