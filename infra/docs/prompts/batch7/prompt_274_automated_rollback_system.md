# PROMPT 274: AUTOMATED ROLLBACK SYSTEM

This architecture defines the real-time automated rollback system for SNISID, ensuring that any deployment anomaly is neutralized instantly without manual intervention.

---

## 1. Rollback Architecture (Integrated Feedback Loop)

SNISID utilizes a closed-loop rollback system where observability signals directly drive the deployment state.

- **Detection Engine**: **Prometheus** and **Istio** provide real-time telemetry on the new version's health.
- **Rollback Controller**: **ArgoCD** and **Argo Rollouts** manage the state transition back to the last known good configuration.
- **Snapshot Manager**: **Velero** provides the ability to restore persistent volume states if a rollback requires data reversion.

---

## 2. Trigger Workflows (Auto-Abort)

Rollbacks are triggered automatically if any of the following "Kill-Switch" conditions are met:

1.  **Error Rate Spike**: 5xx error rate increases by >1% compared to the stable baseline.
2.  **Latency Degradation**: P99 latency increases by >100ms.
3.  **CrashLoopBackOff**: Any new pod fails to reach a "Ready" state within 5 minutes.
4.  **Security Violation**: Kyverno or Falco detects a policy breach in the new version.
5.  **Manual Veto**: An authorized officer can trigger an instant global rollback via a single authenticated API call.

---

## 3. Validation Mechanisms (Post-Rollback)

After a rollback is initiated, the system performs a **Recovery Validation**:

- **Stable Restoration**: Verifies that the previous version pods are back to 100% capacity and "Ready."
- **Traffic Re-convergence**: Istio confirms that 100% of traffic has been shifted away from the failed version.
- **State Integrity**: Checks that the database and message queues are synchronized with the rolled-back application state.

---

## 4. Recovery Orchestration (State-Aware)

- **Database Compatibility**: SNISID enforces a "Two-Version Compatibility" rule for all DB schemas, ensuring the code can be rolled back even after a migration has run.
- **Event Replay**: If the failed version corrupted any transient state, the **Kafka Recovery Worker** replays events from the last healthy checkpoint.
- **Cache Invalidation**: Automated clearing of Redis caches for the affected service to prevent "Ghost Data" from the failed version.

---

## 5. Governance Model

- **Incident Record**: Every automated rollback generates a "Severity-1" incident report in the forensic ledger.
- **Post-Mortem Requirement**: The failed version is "Frozen" in a debug namespace, and the development team must provide a root cause analysis (RCA) before the next deployment attempt is allowed.
- **Global Halt**: If three consecutive rollbacks occur for the same service, the entire CI/CD pipeline for that agency is locked until manually cleared by the Infrastructure Lead.

---

**PROMPT 274 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 275 — COMPLIANCE AS CODE.**
