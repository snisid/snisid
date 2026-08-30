# SNISID: Core Data Architecture & Governance

This document consolidates the official Data Models, Graph Structure, Storage Strategy, and Governance rules for the SNISID platform. It serves as the definitive source of truth for the database structures underlying the backend microservices.

## 1. Core Data Model Schemas

### 🪪 Identity Schema
Canonical representation of a national identity.
```json
{
  "identity_id": "UUID",
  "national_id": "string",
  "status": "ACTIVE | SUSPENDED | DECEASED | FRAUD_REVIEW",
  "citizen_profile_id": "UUID",
  "biometric_profile_id": "UUID",
  "risk_profile_id": "UUID",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "version": "integer"
}
```
*Constraints:* Immutable `identity_id`, soft delete only, mandatory audit on mutation.

### 👤 Citizen Profile Schema
```json
{
  "citizen_profile_id": "UUID",
  "first_name_encrypted": "string (Vault Ciphertext)",
  "first_name_bidx": "string (Blind Index)",
  "last_name_encrypted": "string (Vault Ciphertext)",
  "last_name_bidx": "string (Blind Index)",
  "date_of_birth_encrypted": "string (Vault Ciphertext)",
  "gender": "string",
  "nationality": "string",
  "addresses": [], // Encrypted at object level
  "phone_numbers": [], // Encrypted at object level
  "emails": [], // Encrypted at object level
  "documents": [],
  "created_at": "timestamp"
}
```

### 🧬 Biometric Schema
```json
{
  "biometric_profile_id": "UUID",
  "face_embedding_vector": "vector<float>",
  "fingerprint_hashes": [],
  "iris_hash": "string",
  "liveness_score": "float",
  "deepfake_score": "float",
  "capture_device_id": "UUID",
  "created_at": "timestamp"
}
```

### 🚨 Fraud Schema
```json
{
  "fraud_case_id": "UUID",
  "identity_id": "UUID",
  "fraud_type": "DUPLICATE_IDENTITY | SYNTHETIC_IDENTITY | DEEPFAKE",
  "severity": "LOW | MEDIUM | HIGH | CRITICAL",
  "confidence_score": "float",
  "status": "OPEN | INVESTIGATING | CONFIRMED | CLOSED",
  "detected_by": "AI | RULE_ENGINE | HUMAN",
  "created_at": "timestamp"
}
```

### 💻 Device Schema
```json
{
  "device_id": "UUID",
  "device_type": "MOBILE | DESKTOP | KIOSK",
  "os": "string",
  "browser": "string",
  "ip_address": "string",
  "geo_location": {},
  "risk_score": "float",
  "created_at": "timestamp"
}
```

### ⚠️ Risk Schema
```json
{
  "risk_profile_id": "UUID",
  "identity_id": "UUID",
  "risk_score": "float",
  "risk_level": "LOW | MEDIUM | HIGH | CRITICAL",
  "risk_factors": [],
  "last_computed_at": "timestamp"
}
```

### 📜 Audit Schema
```json
{
  "audit_event_id": "UUID",
  "actor_id": "UUID",
  "action": "string",
  "resource_type": "string",
  "resource_id": "UUID",
  "agency": "string",
  "timestamp": "timestamp",
  "ip_address": "string",
  "correlation_id": "UUID"
}
```

---

## 2. Graph Model (Neo4j)

### 🔵 Nodes
`Citizen`, `Identity`, `Device`, `Agency`, `FraudCase`, `Biometric`, `IPAddress`, `Location`, `Document`.

### 🔗 Relationships
*   `(Citizen)-[:OWNS]->(Identity)`
*   `(Identity)-[:USES]->(Device)`
*   `(Device)-[:CONNECTED_FROM]->(IPAddress)`
*   `(Identity)-[:LINKED_TO]->(FraudCase)`

**Temporal Relationships:** All relationships must support `valid_from`, `valid_to`, and `event_timestamp` properties to enable historical replay and Point-in-Time querying.

**Risk Propagation Edges:** 
If a Device is compromised, linked identities inherit risk, linked IPs inherit suspicion, and associated agencies are alerted. Enables detection of shared devices/addresses and repeated embeddings.

---

## 3. Storage Strategy

### 🐘 Relational Storage (PostgreSQL)
Primary transactional persistence.
*Tables:* `identities`, `citizen_profiles`, `biometric_profiles`, `risk_profiles`, `fraud_cases`, `devices`, `audit_events`, `agencies`.

### ⚡ Ephemeral / Cache (Redis)
*Data:* `session_tokens`, `risk_scores`, `hot_identity_profiles`, `rate_limit_counters`, `active_alerts`.

### 🗄️ Object Storage (MinIO / S3-compatible)
*Data:* Biometric images, scanned documents, forensic snapshots, AI model weights, database backups.

**Detailed Infrastructure**: See the [SNISID Secure Object Storage](file:///c:/Users/sopil/Desktop/SNISID/SNISID_Secure_Object_Storage.md) for WORM immutability, SSE-KMS encryption, and air-gapped archival strategies.

### 📈 Time-Series Storage (VictoriaMetrics / InfluxDB)
*Data:* Performance metrics, telemetry, SOC event frequency, infrastructure monitoring.

---

## 4. Data Governance & Security

### 🔒 Encryption & Sovereignty
*   **Transit:** TLS 1.3 externally, mTLS internally.
*   **Rest:** AES-256 with KMS/Vault integration.
*   **Sovereignty:** All national data remains in-country. Sovereign cloud usage is preferred, with external replication strictly prohibited.

### 🧾 PII Classification
| Classification | Example | Retention |
| :--- | :--- | :--- |
| **Public** | Agency names | Permanent |
| **Sensitive** | Address, Session Logs | 90 days (Logs) |
| **Critical** | Biometrics | Policy-controlled |
| **Restricted** | Intelligence Flags, Audits | 10 years (Audits) |

### 🏛️ Multi-Agency Access
Controlled via strict RBAC (Role-Based Access Control) + ABAC (Attribute-Based Access Control) anchored in Zero Trust principles for agencies such as ANH, DGI, DCPJ, Immigration, and Courts.
