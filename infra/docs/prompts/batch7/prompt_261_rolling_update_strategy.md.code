# PROMPT 261: ROLLING UPDATE STRATEGY

This architecture defines the zero-downtime rolling update strategy for the SNISID platform, ensuring service continuity even during major version upgrades.

---

## 1. Update Workflows (Standardized)

SNISID uses a **Native Kubernetes Rolling Update** strategy combined with custom readiness gates.

1.  **Preparation**: The new container image is pushed to the sovereign registry and scanned for vulnerabilities.
2.  **Trigger**: GitOps (ArgoCD) detects the version change in the manifest.
3.  **Surge Execution**:
    - `maxSurge: 25%`: Kubernetes spins up new pods before terminating old ones.
    - `maxUnavailable: 0%`: Ensures that the full capacity is maintained throughout the update.
4.  **Health Verification**: New pods must pass `Liveness`, `Readiness`, and `Startup` probes before receiving traffic.

---

## 2. Validation Pipelines (Pre/Post Update)

Updates are gated by an automated **Validation Pipeline**:

- **Pre-Update**: Smoke tests in a "Pre-Prod" namespace that mirrors the production environment.
- **Canary Readiness**: Automated verification that the new version does not increase error rates or latency in a small subset of traffic.
- **Service Mesh Gate**: Istio waits for the new pods to be "Ready" before injecting them into the service discovery pool.

---

## 3. Traffic Orchestration

- **Graceful Termination**: Envoy proxies (Istio) ensure that existing connections are drained gracefully before a pod is shut down.
- **Readiness-Based Routing**: Istio only routes traffic to pods that have passed the application-level health check (e.g., connection to Kafka/DB is verified).
- **Circuit Breaking**: If new pods show high failure rates, Istio automatically halts traffic to the new version to prevent cascading failures.

---

## 4. Rollback Strategy (Automated Readiness)

SNISID maintains **Rollback Readiness** for every deployment.

- **Automated Rollback**: If the error rate exceeds 1% or latency increases by 50ms during the update, ArgoCD triggers an automatic `kubectl rollout undo`.
- **Stateful Rollback**: Database migrations are designed to be forward and backward compatible (Expand-Contract pattern) to allow code rollbacks without data loss.
- **Snapshot Recovery**: If a critical failure occurs, the cluster can be restored to its last known healthy state using **Velero** snapshots.

---

## 5. Runtime Governance

- **Update Windows**: Production updates are scheduled during low-traffic national windows unless it's a security emergency.
- **Audit Logs**: Every update step (who triggered it, which pods were replaced, duration) is recorded in the forensic audit ledger.
- **Approval Gates**: Production updates for "Tier-0" services require cryptographic approval from two authorized infrastructure officers.

---

**PROMPT 261 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 262 — BLUE-GREEN DEPLOYMENT.**
