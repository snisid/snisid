# PROMPT 289: AUTOMATED PATCH MANAGEMENT

This architecture defines the zero-downtime, automated patching strategy for the SNISID platform's operating systems, container images, and core infrastructure components.

---

## 1. Patching Architecture (Immutable & Rolling)

SNISID treats patches as versioned infrastructure updates, avoiding "In-Place" mutations where possible.

- **OS Patching (Kured)**: The **KUbernetes REboot Daemon** monitors for the existence of a reboot sentinel (e.g., `/var/run/reboot-required`) and orchestrates safe, rolling node reboots.
- **Container Patching (CI/CD)**: Vulnerability scanners (Prompt 269) trigger automated rebuilds of base distroless images when a CVE is detected.
- **Infrastructure Patching (Terraform)**: Versioned modules ensure that updates to cloud providers or network configurations are applied via the standard GitOps pipeline.

---

## 2. Workflows (The Patching Lifecycle)

1.  **Scanning**: Continuous vulnerability scanning identifies outdated kernels or vulnerable libraries in the production worker pool.
2.  **Staging**: Patches are first applied to the "Shadow Infrastructure" (Prompt 252) where automated E2E tests verify that the patch doesn't break mission-critical services.
3.  **Approval**: High-severity patches (e.g., Zero-Days) bypass standard maintenance windows and trigger an automated "Emergency PR" in the infrastructure repository.
4.  **Rolling Application**: Kured drains a single node at a time, applies the OS update, reboots, and waits for all pods to be healthy before moving to the next node.

---

## 3. Integration Strategy (Zero-Downtime)

- **Pod Disruption Budgets (PDBs)**: Enforced for all services to ensure that the patching orchestrator never drains more nodes than the application can tolerate.
- **Istio Traffic Shifting**: During a node reboot, the service mesh automatically reroutes traffic to pods in other availability zones, ensuring zero user-facing downtime.
- **Stateful Workloads**: Automated pre-stop hooks ensure that databases and Kafka brokers perform a graceful handoff of leadership before the node is rebooted.

---

## 4. Security & Privacy

- **Cryptographic Attestation**: Every patch applied to the OS or a container image must be signed by the national build engine, proving it hasn't been tampered with.
- **Isolated Update Mirrors**: Patches are pulled from internal, air-gapped repository mirrors that have been pre-scanned and approved by the national security team.
- **Audit Ledger**: Every node reboot and container swap is recorded in the forensic ledger, including the CVE ID that triggered the patch.

---

## 5. Governance Model

- **SLA for Vulnerabilities**: Critical (24 hours), High (7 days), Medium (30 days). The system automatically escalates if a patch is not applied within these windows.
- **Patch Transparency**: Real-time dashboard showing the "Patch Level" of every cluster and node in the national federation.
- **Rollback Readiness**: The system maintains a "Last Known Good" version of the OS and container images for instant rollback if a patch causes unexpected performance degradation.

---

**PROMPT 289 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 290 — FINOPS & COST MANAGEMENT AUTOMATION.**
