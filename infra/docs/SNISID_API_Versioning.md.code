# SNISID: API Versioning & Backward Compatibility Strategy

Because SNISID is a national infrastructure platform consumed by multiple independent government agencies (Tax, Police, Immigration), breaking changes can disrupt critical state operations. This document defines the strict API versioning, schema evolution, and lifecycle management rules required to guarantee backward compatibility.

---

## 1. API Lifecycle Model

Every API endpoint and Kafka event schema in SNISID follows a strict 4-stage lifecycle:

1.  **DRAFT / EXPERIMENTAL:** New endpoints being tested in the Staging environment. Not guaranteed to be stable.
2.  **ACTIVE:** The current, fully supported production version of the API.
3.  **DEPRECATED:** The API is still fully functional in production, but a newer version exists. New integrations are prohibited. A sunset date is officially published.
4.  **RETIRED:** The endpoint is permanently disabled and returns an `HTTP 410 Gone`.

---

## 2. Versioning Rules

### URI-Based Major Versioning
For maximum explicit clarity to external agencies, SNISID uses **URI-based versioning** for major, breaking changes.
*   **Format:** `https://api.snisid.gov/identity/v1/citizens/{id}`
*   Major versions (v1, v2) imply incompatible, breaking changes.

### Non-Breaking Minor Evolution
We do **not** bump the URI version for minor changes. APIs are designed to evolve in place. Clients must be written robustly to ignore unrecognized fields.

**What is a BREAKING change? (Requires a new Major Version)**
*   Renaming or removing an existing field.
*   Changing the data type of an existing field (e.g., Integer to String).
*   Adding a new *required* request parameter or payload field.
*   Changing or removing an existing error code format.

**What is a NON-BREAKING change? (Evolves in the current version)**
*   Adding new endpoints to the API.
*   Adding new *optional* request parameters.
*   Adding new fields to the response payload.

---

## 3. Migration and Deprecation Strategy

Government IT systems move slowly. Therefore, the deprecation window is intentionally long.

1.  **Simultaneous Support:** The SNISID API Gateway must support `v1` and `v2` simultaneously in production for a minimum overlapping period.
2.  **Sunset Window:** A deprecated API must remain fully operational for a minimum of **18 to 24 months** after the `v2` release.
3.  **Sunset Headers:** Once `v1` is marked Deprecated, the API Gateway automatically injects IETF standard headers into all `v1` responses to warn consumers:
    *   `Deprecation: @1672531199` (Timestamp when it was deprecated)
    *   `Sunset: @1704067199` (Timestamp when it will be turned off completely)
    *   `Link: <https://docs.snisid.gov/migration/v2>; rel="deprecation"`
4.  **Backend Pattern:** Where possible, the API Gateway uses the "Tolerant Reader" pattern, mapping `v1` requests internally to the `v2` backend microservice, reducing the need to maintain two separate codebases.

---

## 4. Event Schema Evolution Rules (Kafka)

Asynchronous events (Kafka) are the backbone of SNISID. Event schemas are even harder to evolve than APIs because data is stored persistently.

1.  **Schema Registry:** All Kafka topics require strict schemas defined in **Protobuf** or **Avro** and managed via a centralized Schema Registry.
2.  **Compatibility Mode:** The Schema Registry is strictly locked to **BACKWARD_TRANSITIVE** compatibility mode.
    *   This means a new schema can be used to read data produced with all previous schemas.
3.  **Deployment Order (Two-Phase Upgrade):**
    *   *Step 1:* Upgrade the **Consumers** first to recognize the new schema (they will simply ignore the missing new fields from old events).
    *   *Step 2:* Once all consumers are deployed, upgrade the **Producers** to start publishing events using the new schema.
4.  **Tombstoning vs Deletion:** Fields in an event schema cannot be deleted. If a field is no longer used, it is marked as `@deprecated` in the Protobuf definition, and producers begin sending it as `null`.

---

## 5. Compatibility Matrix Strategy

To prevent a cascading update nightmare (where updating the Identity Service forces an immediate update of the Fraud Service), services must adhere to **Consumer-Driven Contracts (CDC)**.

*   **Pact Testing:** Service A defines a contract (a "Pact") detailing exactly what it expects from Service B's API.
*   Before Service B can deploy a new version to production, the CI/CD pipeline runs Service A's contract tests against Service B.
*   If Service B accidentally breaks the contract (e.g., changes a field name Service A relies on), the deployment is automatically blocked.
*   This creates an automated, living compatibility matrix where independent teams can deploy confidently, knowing mathematically that they are not breaking downstream government or internal consumers.
