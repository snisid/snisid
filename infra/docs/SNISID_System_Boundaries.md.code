# SNISID: System Boundaries and Trust Architecture

This document defines the strict system boundaries, trust levels, and security enforcement mechanisms for the National Secure Identity & Intelligence System (SNISID). It ensures that the highly sensitive internal core is mathematically and architecturally isolated from external entities.

## 1. Trust Zone Diagram

```text
+-----------------------------------------------------------------------------------------------+
|                                     [ZONE 3: UNTRUSTED]                                       |
|                               (Public Citizens, 3rd Party APIs)                               |
|                                    Trust Level: 0 (None)                                      |
+-----------------------------------------------------------------------------------------------+
                                               ||
                                  [Internet / Public Network]
                                               ||
+-----------------------------------------------------------------------------------------------+
|                             [ZONE 1: DMZ / ENFORCEMENT BOUNDARY]                              |
|                              (API Gateway, WAF, Auth Brokers)                                 |
|                             Trust Level: 2 (Verification Layer)                               |
+-----------------------------------------------------------------------------------------------+
                 ||                                                         ||
       [mTLS / Govt Intranet]                                      [Internal mTLS Network]
                 ||                                                         ||
+----------------------------------+                   +----------------------------------------+
|      [ZONE 2: SEMI-TRUSTED]      |                   |          [ZONE 0: CORE TRUST]          |
|     (Government Agencies:        |                   |          (SNISID Core System)          |
|    Police, Tax, Immigration)     |                   |                                        |
|   Trust Level: 3 (Authorized)    |                   |   - Identity DB (PostgreSQL)           |
+----------------------------------+                   |   - Fraud/Graph Intelligence (Neo4j)   |
                                                       |   - Event Bus (Kafka)                  |
                                                       |   - SOC / Data Warehouses              |
                                                       |                                        |
                                                       |    Trust Level: 5 (Absolute Trust)     |
                                                       +----------------------------------------+
```

## 2. Trust Levels and Boundary Rules

### Zone 0: SNISID Core System (Absolute Trust)
*   **Definition:** The highly secure backend infrastructure where raw data, biometrics, and core business logic reside.
*   **Trust Level:** 5 (Absolute Trust). All internal components still use SPIFFE/SPIRE for service-to-service authentication (Zero Trust).
*   **Boundary Rules:** 
    *   No inbound connections from Zone 3 (Public Internet). 
    *   No direct inbound connections from Zone 2 (Govt Agencies). 
    *   All inbound traffic *must* route exclusively through Zone 1 (DMZ/API Gateway).
    *   No outbound internet access for any databases or core worker nodes.

### Zone 1: DMZ / Enforcement Boundary (Verification Layer)
*   **Definition:** The perimeter layer containing Web Application Firewalls (WAF), Load Balancers, and the API Gateway.
*   **Trust Level:** 2.
*   **Boundary Rules:**
    *   Terminates public TLS and mTLS connections.
    *   Validates JWTs and scopes before proxying requests to Zone 0.
    *   Enforces strict rate limiting, DDoS mitigation, and payload inspection.

### Zone 2: Government Agencies (External Trusted Entities)
*   **Definition:** Approved state actors (Tax Authority, Police, Immigration) operating on isolated networks.
*   **Trust Level:** 3.
*   **Boundary Rules:**
    *   Agencies are trusted to *request* data, but their requests are never inherently trusted as safe.
    *   Cannot execute SQL queries or access the event bus directly.
    *   Must connect via dedicated leased lines or Gov-VPN.
    *   Must authenticate via Mutual TLS (mTLS) using certificates issued by the National CA.

### Zone 3: External Systems (Untrusted)
*   **Definition:** The public internet, third-party systems, citizens' mobile apps, and non-state verifiers.
*   **Trust Level:** 0.
*   **Boundary Rules:**
    *   Can only reach specific public-facing endpoints in Zone 1.
    *   Strict CAPTCHA, behavioral analysis, and rate-limiting applied.

---

## 3. Data Flow Restrictions

1.  **Read Operations (External -> Core):**
    *   **Zone 3 (Public) to Zone 0:** Denied. Must pass through Zone 1. Can only read own data (Citizen Portal) using strict MFA tokens.
    *   **Zone 2 (Agencies) to Zone 0:** Denied direct access. Must pass through Zone 1. Data is minimized (e.g., Police query returns a binary "Valid/Invalid" rather than full biometric templates).
2.  **Write Operations (Core Modification):**
    *   **Zone 3:** Can only initiate "Requests" (e.g., address change request), never direct writes.
    *   **Zone 2:** Agencies cannot modify core identity data unless specifically authorized via cryptographic multi-signature flows (e.g., Civil Registry registering a death).
3.  **Data Egress (Core -> External):**
    *   Raw biometric data (faces, fingerprints) **never** leaves Zone 0. It is only used internally for matching. External queries receive cryptographic proofs or boolean results.
    *   No bulk data extraction is permitted over the network.

---

## 4. Security Enforcement Architecture

*   **Ingress Controller & WAF (Zone 1):** Intercepts all traffic. Drops malformed packets, blocks known malicious IPs, and defends against OWASP Top 10 vulnerabilities.
*   **API Gateway Policy Engine (Zone 1):** Executes Open Policy Agent (OPA) rules. Validates that the requesting entity (e.g., Agency X) has the specific `scope` required to access the requested endpoint (e.g., `/api/v1/identity/verify`).
*   **Service Mesh (Zone 0):** Uses Istio/Linkerd. Ensures that even if the API Gateway is breached, internal services require cryptographic identity (mTLS) to talk to one another (e.g., `IdentityService` cannot talk to `GraphDB` unless explicitly allowed by the mesh policy).
*   **Data-at-Rest Encryption (Zone 0):** All data volumes are encrypted (AES-256) utilizing hardware security modules (HSM) to store the root keys.

---

## 5. Audit Requirements Between Zones

To ensure accountability and trace potential breaches, strict auditing is enforced at all boundary crossings:

1.  **Boundary Crossing Logs:** Every request passing from Zone 3/2 -> Zone 1, and Zone 1 -> Zone 0 is logged.
2.  **Required Audit Fields:**
    *   Timestamp (NTP synchronized).
    *   Source IP / Source Agency ID.
    *   Target Endpoint / Resource.
    *   Cryptographic Identity (Client Certificate Thumbprint or JWT ID).
    *   Action Taken (Allow/Deny/Rate Limited).
3.  **Immutability:** Audit logs are written asynchronously to a write-once, read-many (WORM) storage system.
4.  **Anomaly Detection:** The SOC (Zone 0) continually analyzes the audit stream (via Kafka) using Machine Learning to detect "impossible travel" (e.g., a citizen logging in from two different countries) or "agency abuse" (e.g., Police querying 10,000 identities per second).
