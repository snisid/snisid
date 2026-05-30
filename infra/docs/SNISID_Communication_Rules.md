# SNISID: Service Communication & Resilience Rules

To prevent cascading failures in a highly distributed national architecture, SNISID enforces strict rules on how services communicate. 

## 1. Sync vs. Async Rules

### When to use Synchronous (REST / gRPC)
Synchronous communication is blocking. The caller waits for the response. It should ONLY be used when:
1.  **Immediate Response is Mandatory:** A citizen waiting at a border checkpoint requires a synchronous pass/fail (e.g., API Gateway calling Auth Service, Fraud calling AI Inference).
2.  **Read Operations:** Fetching data to display on a UI or returning a query to an agency.

### When to use Asynchronous (Kafka)
Asynchronous communication is non-blocking. It is the default for SNISID state changes.
1.  **State Mutations:** Creating, updating, or deleting data (e.g., enrolling a new identity).
2.  **Cross-Domain Triggers:** Identity creation triggering a Fraud check.
3.  **High-Latency Operations:** Document parsing, extensive graph clustering algorithms, sending emails/SMS.

---

## 2. Timeout Policies

No service is allowed to wait indefinitely. All synchronous network calls must have strict, enforceable timeouts.

*   **P99 Edge Timeout (API Gateway):** Maximum 10 seconds. If internal services haven't responded, the Gateway returns a `504 Gateway Timeout`.
*   **Internal Service-to-Service (gRPC):** Maximum 2 seconds. Internal calls must fail fast.
*   **Database Query Timeout (PostgreSQL/Neo4j):** Maximum 5 seconds for OLTP transactions. Long-running analytical queries must be routed to read-replicas.

---

## 3. Retry Logic & Exponential Backoff

If a synchronous call or database query fails due to a transient error (e.g., network blip, lock timeout), the caller must retry intelligently.

*   **Idempotent Methods Only:** Retries are ONLY allowed on `GET`, `PUT`, and `DELETE` requests, or on `POST` requests if they include a unique Idempotency Key (e.g., `X-Idempotency-Key: UUID`).
*   **Exponential Backoff:** Retries must not bombard a struggling service.
    *   Attempt 1: Wait 100ms
    *   Attempt 2: Wait 500ms
    *   Attempt 3: Wait 2000ms
    *   Max Retries: 3. If the 3rd fails, fallback or return an error.
*   **Jitter:** A random timing variance (+/- 10%) must be added to the backoff to prevent "thundering herd" problems where hundreds of failing pods retry at the exact same millisecond.

---

## 4. Circuit Breakers

To protect degraded services from being overwhelmed by retries, all inter-service synchronous calls must be wrapped in a Circuit Breaker (managed via Istio Service Mesh or code-level libraries like resilience4j / go-breaker).

*   **CLOSED State (Normal):** Requests flow freely.
*   **OPEN State (Tripped):** If the failure rate exceeds 50% over a 10-second sliding window, the circuit "opens". All subsequent requests immediately fail with a `503 Service Unavailable` without attempting the network call, giving the downstream service time to recover.
*   **HALF-OPEN State (Testing):** After a 30-second cooldown, the circuit allows a single test request through. If it succeeds, the circuit CLOSES. If it fails, the circuit OPENS for another 30 seconds.

---

## 5. Service Discovery Contracts

*   **No Hardcoded IPs:** Services must never reference hardcoded IPs. 
*   **DNS Resolution:** Internal routing uses Kubernetes CoreDNS (e.g., `fraud-service.snisid.svc.cluster.local`).
*   **Service Mesh Abstraction:** Developers do not implement service discovery in application code. The Envoy sidecar proxy transparently handles routing, load balancing, and discovering healthy pod endpoints.
