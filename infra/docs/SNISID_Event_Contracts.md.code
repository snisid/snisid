# SNISID: Kafka Event Contracts

Event Contracts form the nervous system of the SNISID platform. By enforcing strict schemas over Apache Kafka, we guarantee asynchronous stability, historical replayability, and safe microservice decoupling.

## 1. Standard Event Envelope

Every event published to Kafka MUST follow this standard envelope. The envelope provides the necessary metadata for tracing, idempotency, and schema evolution, leaving the business logic inside the `payload` object.

```json
{
  "event_id": "UUID",              // Unique ID for idempotency (e.g., UUIDv7 for sortability)
  "event_type": "string",          // Follows <domain>.<entity>.<action> format
  "event_version": "v1",           // Schema version of the payload
  "timestamp": "ISO8601",          // Exact time of the event (UTC)
  "correlation_id": "UUID",        // Distributed tracing ID (OpenTelemetry)
  "producer": "service_name",      // e.g., "identity-service"
  "payload": { ... }               // The business data (schema defined in Registry)
}
```

---

## 2. Topic Taxonomy & Naming Conventions

**Convention:** `<domain>.<entity>.<action>`

### Core Topics
| Topic Name | Producer | Description |
| :--- | :--- | :--- |
| `identity.citizen.created` | Identity Service | Fired when a new identity is enrolled. |
| `identity.citizen.updated` | Identity Service | Fired on demographic updates. |
| `identity.citizen.deleted` | Identity Service | Fired on soft deletion. |
| `biometric.capture.verified` | AI Inference Service | Fired when a biometric match succeeds/fails. |
| `fraud.case.detected` | Fraud Service | Fired when the engine flags a risk. |
| `fraud.case.confirmed` | Fraud Service | Fired when an analyst confirms fraud. |
| `risk.profile.updated` | Fraud Service | Fired when an identity's aggregate risk changes. |
| `audit.log.recorded` | *All Services* | High-volume topic for immutable action logging. |
| `soc.alert.generated` | SOC Service | Fired to trigger SOAR playbooks. |
| `device.session.registered` | Auth Service | Fired when a device connects. |
| `device.session.flagged` | Fraud Service | Fired when a device is compromised. |

---

## 3. Schema Governance & Versioning

To ensure consumers do not break when producers update logic, we enforce strict Schema Registry controls (Avro/Protobuf).

*   **Rule 1:** Only **Backward Compatible** changes are permitted.
*   **Rule 2:** Additive fields are heavily preferred. Existing fields cannot be deleted or renamed.
*   **Rule 3:** Historical events are strictly immutable.
*   **Example Evolution:**
    *   `v1` → Additive only (e.g., adding an optional `middle_name`).
    *   `v2` → Breaking schema allowed, requires publishing to a new topic (e.g., `identity.citizen.v2.created`).

---

## 4. Replay Strategy

Historical replay is a critical requirement for AI model retraining, SOC forensics, and graph reconstruction.

*   **Persistence:** Full event persistence with immutable log retention (retention set to `compact` or absolute time thresholds based on data classification).
*   **Replay Engine Architecture:**
    1.  A dedicated `Replay Service` reads from the immutable Kafka topic.
    2.  It filters events based on requested parameters: `Time Window`, `Identity ID`, `Agency`, or `Fraud Case ID`.
    3.  It pushes filtered events to a temporary topic consumed by the requesting service (e.g., `Graph Intelligence Service` rebuilding a compromised node).

---

## 5. Dead Letter Queue (DLQ) Contracts

To prevent poison pills from blocking partitions, all consumers must implement DLQ logic.

### Standard DLQ Topics
*   `identity.dlq`
*   `fraud.dlq`
*   `soc.dlq`

### DLQ Rules
1.  **Max Retry Count:** Consumers must attempt to process an event with exponential backoff (e.g., 3 retries) before routing it to the DLQ.
2.  **Quarantine:** Malformed events (e.g., schema validation failure) are immediately quarantined in the DLQ without retry.
3.  **Forensic Retention:** DLQ topics have mandatory permanent forensic retention until manually resolved by a data engineer.
4.  **DLQ Envelope Addition:** When an event is moved to a DLQ, the consumer must wrap it with an `error_reason` and `failed_at` timestamp.
