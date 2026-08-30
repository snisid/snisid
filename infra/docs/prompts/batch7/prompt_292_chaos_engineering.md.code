# PROMPT 292: CHAOS ENGINEERING & FAULT INJECTION PIPELINE

This architecture defines the proactive resilience testing strategy for the SNISID platform, ensuring that the platform's self-healing and auto-recovery mechanisms are effective under extreme, unpredictable failure conditions.

---

## 1. Chaos Architecture (Controlled Turbulence)

SNISID utilizes a sophisticated chaos engineering stack to inject controlled faults into the infrastructure.

- **Chaos Orchestrator**: **Chaos Mesh** and **LitmusChaos** handle the definition and execution of chaos experiments.
- **Fault Injectors**:
    - **Pod Chaos**: Randomly killing pods or containers.
    - **Network Chaos**: Injecting latency, packet loss, or partitioning between regions.
    - **IO Chaos**: Simulating slow or failing disk drives.
    - **Kernel Chaos**: Simulating kernel panics or resource exhaustion.
- **Observability Integration**: Chaos experiments are tightly coupled with Prometheus and Grafana to measure the impact on "Steady State" metrics.

---

## 2. Experiment Workflows (Continuous Testing)

1.  **Steady State Definition**: The system defines a "Normal" performance baseline for the target service (e.g., "99th percentile latency < 200ms").
2.  **Hypothesis Formulation**: "If we kill the primary database node, the standby should take over in < 15 seconds with zero data loss."
3.  **Experiment Execution**: The chaos engine injects the fault in a controlled, automated fashion.
4.  **Verification**: The system automatically checks if the "Steady State" was maintained or if the system recovered within the defined SLA.
5.  **Rollback**: If the experiment causes unintended cascade failures, the chaos engine instantly terminates the fault injection.

---

## 3. Integration Strategy (The "Chaos Pipeline")

- **Pre-Production Gating**: Chaos experiments are part of the standard CI/CD pipeline; a new service cannot be promoted to production until it passes a "Resilience Suite."
- **Game Day Orchestration**: Automated monthly events where multiple faults are injected simultaneously across a full regional cluster to test inter-service dependencies.
- **AI-Driven Blast Radius**: The system uses AI to predict the blast radius of an experiment and automatically adjusts the intensity to prevent actual user-facing outages.

---

## 4. Analysis & Reporting

- **Resilience Scorecard**: Every microservice is assigned a "Chaos Grade" based on how many fault scenarios it survived.
- **Cascade Failure Visualization**: Tracing data (Prompt 279) is used to visualize how a fault in one service impacted downstream consumers during a chaos experiment.
- **Auto-Issue Generation**: Failed chaos experiments automatically trigger Jira/GitLab issues with the full forensic context and logs for the engineering team.

---

## 5. Governance Model

- **Safe-to-Chaos Windows**: Experiments in production are only allowed during authorized "Chaos Windows" and require a "Chaos Sentinel" (automated watcher) to be active.
- **Audit Ledger**: Every chaos experiment, including the hypothesis, the raw results, and the identities of those who authorized it, is recorded in the forensic ledger.
- **National Resilience Report**: Monthly report for national security leaders detailing the platform's demonstrated ability to survive regional disasters and targeted infrastructure attacks.

---

**PROMPT 292 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 293 — AUTOMATED COMPLIANCE REPORTING & AUDITING.**
