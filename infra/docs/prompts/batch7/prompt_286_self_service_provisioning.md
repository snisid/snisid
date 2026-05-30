# PROMPT 286: SELF-SERVICE RESOURCE PROVISIONING

This architecture defines the automated, self-service strategy for infrastructure provisioning within the SNISID platform, empowering agency teams to deploy compliant resources without manual intervention from the central infrastructure team.

---

## 1. Provisioning Architecture (Abstraction Layer)

SNISID utilizes a "Product-Centric" infrastructure model where complex resources are exposed as simple, declarative templates.

- **Frontend**: The **Developer Portal** (Prompt 282) provides a UI for requesting resources.
- **Orchestrator**: **Crossplane** (Kubernetes-native) or **Terraform Cloud/Enterprise** handles the lifecycle of cloud and on-premise resources.
- **Provider Layer**: Standardized providers for the national sovereign cloud (OpenStack/vSphere) and public cloud extensions.
- **Policy Gate**: **Kyverno/OPA** validates every request against national quotas and security baselines.

---

## 2. Provisioning Workflows (The "Golden Path")

1.  **Selection**: A developer selects a "Resource Template" (e.g., "Isolated PostgreSQL Database") from the portal.
2.  **Configuration**: The developer provides minimal inputs (e.g., `db_size`, `owning_agency`, `environment`).
3.  **Validation**: The system automatically checks the requesting user's permissions and the project's remaining budget/quota.
4.  **Execution**: Crossplane creates the resource in the specified cloud provider and automatically injects the connection credentials into the developer's Kubernetes namespace via HashiCorp Vault.

---

## 3. Governance Model (Controlled Autonomy)

- **Quota Management**: Agencies are assigned "National Resource Quotas" (CPU, RAM, Storage, Budget); self-service requests that exceed these quotas are automatically routed for secondary approval.
- **Standardized Tags**: Every provisioned resource is automatically tagged with mandatory metadata for billing and audit (e.g., `snisid-id`, `created-by`, `compliance-tier`).
- **TTL (Time-To-Live)**: Development and Sandbox resources are created with a mandatory TTL and are automatically deleted after 30 days unless an extension is requested.

---

## 4. Lifecycle Management (Continuous Reconciliation)

- **Drift Detection**: The orchestrator continuously monitors the live state of the resource; if a manual change occurs in the cloud console, the system automatically reverts it to the Git-defined state.
- **In-Place Upgrades**: Infrastructure updates (e.g., database version upgrades) are managed via the same self-service interface, ensuring zero-downtime migrations.
- **Retirement**: When a project is marked as "Archived," all associated self-service resources are automatically snapshotted and then decommissioned to save national resources.

---

## 5. Integration Strategy

- **Secret Management**: Connection strings and API keys are never exposed to the user; they are automatically written to the service's Vault path and mounted as K8s secrets.
- **Monitoring Integration**: Provisioning a new resource automatically creates the corresponding dashboards in Grafana and alerting rules in Alertmanager.
- **Audit Ledger**: Every provisioning, modification, and deletion event is cryptographically signed and recorded in the forensic ledger.

---

**PROMPT 286 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 287 — AUTOMATED CAPACITY PLANNING & RIGHTSIZING.**
