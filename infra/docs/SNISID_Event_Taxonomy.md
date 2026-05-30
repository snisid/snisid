# SNISID: Official Event Taxonomy

This document establishes the official Kafka Event Taxonomy for the SNISID platform. It guarantees that as the system scales to handle national-level throughput across multiple government agencies, all events remain discoverable, strictly governed, and forensically traceable.

---

## 1. Naming Convention Rules

To ensure clarity and compatibility with the Kafka Schema Registry, all topics and event types MUST follow a strict dot-delimited hierarchical naming convention:

**Format:** `[domain].[entity].[action]`

*   **`[domain]`**: The bounded context owning the data (e.g., `identity`, `fraud`, `soc`).
*   **`[entity]`**: The specific aggregate or resource being affected (e.g., `citizen`, `case`, `session`).
*   **`[action]`**: A past-tense verb describing the immutable state change (e.g., `created`, `updated`, `flagged`).

*Example:* `identity.citizen.created`

---

## 2. Standard Event Envelope Schema

Every event published within SNISID MUST be wrapped in the following standardized JSON/Protobuf envelope. This isolates business logic from the metadata required for distributed tracing, replay, and routing.

```json
{
  "metadata": {
    "event_id": "UUID",               // Unique UUIDv7 for idempotency and sortability
    "event_type": "string",           // Must match the [domain].[entity].[action] convention
    "schema_version": "integer",      // E.g., 1. Tied to Confluent Schema Registry
    "timestamp": "ISO8601",           // Exact UTC time the event occurred
    "correlation_id": "UUID",         // For OpenTelemetry distributed tracing across microservices
    "producer_service": "string",     // Identifier of the microservice (e.g., "auth-service-v1.4")
    "agency_context": "string"        // The agency that triggered the flow (e.g., "DCPJ")
  },
  "payload": {
    // The domain-specific business data (immutable state change)
    // E.g., { "national_id": "12345", "status": "ACTIVE" }
  }
}
```

---

## 3. Event Taxonomy Hierarchy & Ownership

The following matrix defines the official taxonomy, the owning service responsible for producing the events, and the immutable semantics of each action.

### 3.1. Identity Domain
**Owner:** Identity Service

| Event Type | Trigger Condition |
| :--- | :--- |
| `identity.citizen.created` | A new citizen profile is successfully enrolled. |
| `identity.citizen.updated` | Demographic or contact information is modified. |
| `identity.citizen.suspended` | An identity is placed under manual review or frozen. |
| `identity.citizen.deceased` | Vital records confirms the death of a citizen. |

### 3.2. Fraud Domain
**Owner:** Fraud Detection Service

| Event Type | Trigger Condition |
| :--- | :--- |
| `fraud.case.opened` | An anomaly triggers the creation of a new investigation case. |
| `fraud.case.updated` | AI or a human analyst adds evidence to the case. |
| `fraud.case.closed` | The case is resolved (marked as false positive or confirmed). |

### 3.3. Risk Domain
**Owner:** Fraud / Risk Service

| Event Type | Trigger Condition |
| :--- | :--- |
| `risk.profile.computed` | The aggregate risk score for an identity changes based on velocity rules or graph intelligence. |

### 3.4. Audit Domain
**Owner:** All Services (Ingested via Kafka Connect to SIEM)

| Event Type | Trigger Condition |
| :--- | :--- |
| `audit.record.logged` | High-throughput logging of every state-mutating action or high-clearance data access. |

**Technical Implementation**: See the [SNISID Sovereign Audit Ledger](file:///c:/Users/sopil/Desktop/SNISID/SNISID_Sovereign_Audit_Ledger.md) for immutable storage, Merkle-tree integrity, and forensic retention models.

### 3.5. SOC Alert Events
**Owner:** SOC Alert Service

| Event Type | Trigger Condition |
| :--- | :--- |
| `soc.alert.generated` | Standard priority alert mapped to a MITRE ATT&CK tactic. |
| `soc.alert.critical` | Severe anomaly detected (e.g., mass biometric spoofing attempt). Triggers immediate SOAR containment. |

### 3.6. Device Events
**Owner:** Fraud / Intelligence Services

| Event Type | Trigger Condition |
| :--- | :--- |
| `device.fingerprint.registered` | A new hardware/browser fingerprint is seen. |
| `device.fingerprint.flagged` | A device is associated with a known fraud ring. |

### 3.7. Authentication Events
**Owner:** Authentication Service

| Event Type | Trigger Condition |
| :--- | :--- |
| `auth.session.started` | A user or agency agent successfully logs in and receives a JWT. |
| `auth.session.failed` | Invalid credentials, OTP failure, or biometric rejection. |
| `auth.session.revoked` | A token is explicitly destroyed (e.g., by SOC containment). |

### 3.8. System Health Events
**Owner:** Infrastructure / K8s Operators

| Event Type | Trigger Condition |
| :--- | :--- |
| `system.node.degraded` | A K8s node or DB replica is failing readiness probes. |
| `system.circuit.opened` | A microservice circuit breaker trips due to high error rates. |

---

## 4. Governance & Versioning Rules

### Immutable Event Principles
An event represents a fact that *already happened*. It cannot be un-happened, deleted, or modified.
*   **Rule:** Consumers must treat all events as read-only.
*   **Rule:** If a mistake is made (e.g., a citizen's name was spelled wrong in `identity.citizen.created`), the producer must emit a compensating event (`identity.citizen.updated`) rather than modifying the original event in Kafka.

### Event Versioning Strategy
SNISID uses Confluent Schema Registry (Protobuf/Avro) with strict `BACKWARD_TRANSITIVE` compatibility mode.
*   **Minor Evolution (v1 -> v1.1):** Adding new, optional fields to the `payload`. The event name remains the same. Downstream consumers will simply ignore the new fields until they are updated.
*   **Major Evolution (v1 -> v2):** Renaming fields, changing data types, or deleting fields. **This is strictly prohibited on existing topics.** If a breaking change is absolutely necessary, the producer must publish to a entirely new event name (e.g., `identity.citizen.v2.created`), and both versions must be maintained during the 18-24 month deprecation window.

### Replay and Forensic Traceability
*   The `event_id` and `timestamp` fields enable the Replay Service to seek back in time and stream historical events in exact chronological order.
*   The `correlation_id` ensures that if `auth.session.failed` triggers a `soc.alert.critical`, the SOC analyst can use the correlation ID in Jaeger/Splunk to trace the exact network hop and payload that started the chain.
