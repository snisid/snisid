# PROMPT 272: ARGO CD CONFIGURATION

This architecture defines the declarative GitOps engine for the SNISID platform, using Argo CD to ensure cluster state matches the national source of truth.

---

## 1. Argo CD Architecture (Hub-and-Spoke)

SNISID uses a centralized Argo CD instance to manage multiple regional clusters.

- **Argo CD Hub**: Deployed in the management cluster (`snisid-admin`).
- **Target Clusters**: Regional worker clusters (Alpha, Beta, Gamma) are registered as external destinations.
- **ApplicationSets**: Used to automatically deploy a service across all clusters based on labels or directory structure.

---

## 2. Synchronization Workflows

1.  **Detection**: Argo CD polls the `snisid-manifests` Git repository every 3 minutes (or receives a webhook).
2.  **Diffing**: Compares the live cluster state with the Git manifests.
3.  **Syncing**:
    - **Automated Sync**: Enabled for Alpha/Beta environments.
    - **Manual Sync**: Required for Production; changes are staged but not applied until approved.
4.  **Pruning**: Automatically deletes resources in the cluster that are no longer defined in Git.

---

## 3. Security Enforcement (Governance)

- **Namespace Restriction**: Argo CD "Projects" restrict which repositories can deploy to which namespaces (e.g., `agency-intelligence` can only deploy to `intelligence-*`).
- **Manifest Validation**: Integrated with **Kyverno** to ensure all resources comply with national security standards before they are applied.
- **RBAC**: Integration with the National SSO (Dex/Keycloak); users are assigned granular roles (Viewer, Sync-Only, Admin) based on their agency.
- **Resource Whitelist**: Argo CD is prohibited from managing sensitive resources like `ClusterRoles` or `Namespaces` directly; these are handled via the IaC pipeline.

---

## 4. Multi-Cluster Management

- **Cluster Secret Management**: Cluster credentials are encrypted using HashiCorp Vault and injected into Argo CD.
- **Regional Sharding**: For national-scale deployments, Argo CD controllers are sharded to handle thousands of applications across dozens of clusters without performance degradation.
- **Inter-Cluster Dependencies**: Sync waves ensure that shared services (e.g., Service Mesh, Vault) are fully operational before application microservices are deployed.

---

## 5. Resilience Model (Self-Healing)

- **Self-Healing**: Enabled globally; any manual "drift" from the Git state is automatically corrected within 60 seconds.
- **Sync Windows**: Prohibits synchronization during critical national holidays or maintenance periods to prevent accidental outages.
- **High Availability**: The Argo CD control plane is deployed with multiple replicas, persistent volume backends, and automated backups of the `argocd-cm` configuration.

---

**PROMPT 272 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 273 — AUTOMATED IMAGE PROMOTION.**
