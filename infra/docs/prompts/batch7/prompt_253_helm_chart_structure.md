# PROMPT 253: HELM CHART STRUCTURE

This architecture defines the standardized Helm chart structure for all SNISID microservices, ensuring modularity, security, and GitOps compatibility.

---

## 1. Helm Repository Structure

SNISID maintains a centralized, air-gap-ready OCI-compliant Helm repository.

```
charts/
├── library/             # Reusable helper templates
│   └── snisid-common/
├── base/                # Base charts for service types
│   ├── snisid-go-app/
│   ├── snisid-python-ml/
│   └── snisid-stateful-db/
└── services/            # Application-specific charts
    ├── core-orchestrator/
    ├── identity-verify/
    └── fraud-engine/
```

---

## 2. Template Architecture (Modular & Reusable)

Every SNISID chart inherits from the **SNISID Common Library** to ensure consistency across security and observability configurations.

- **`deployment.yaml`**: Standardized resource limits, affinity rules, and priority classes.
- **`security.yaml`**: Automatically injects Istio Sidecars, SPIFFE identities, and NetworkPolicies.
- **`service.yaml`**: Standardized service discovery and mTLS port definitions.
- **`vault-secret.yaml`**: Templates for the **Secrets Store CSI Driver** to fetch keys from HashiCorp Vault.

---

## 3. Deployment Workflows (GitOps Driven)

1.  **Chart Packaging**: The CI pipeline packages the chart and pushes it to the OCI registry (e.g., Harbor).
2.  **Manifest Definition**: The GitOps repository contains a `HelmRelease` (or ArgoCD `Application`) referencing the chart version.
3.  **Environment Overrides**: Use `values.yaml` for defaults, and `values-alpha.yaml`, `values-prod.yaml` for regional/environment-specific configuration.
4.  **Reconciliation**: ArgoCD detects the change and applies the Helm chart to the cluster.

---

## 4. Versioning Strategy

- **Semantic Versioning (SemVer)**: `MAJOR.MINOR.PATCH`.
- **Immutable Versions**: Once a chart version is pushed (e.g., `1.2.3`), it cannot be overwritten.
- **Dependency Pinning**: All sub-charts and library charts must be pinned to exact versions in `Chart.yaml`.

---

## 5. Governance Standards

- **Resource Hardening**: All charts MUST define `resources.requests` and `resources.limits`.
- **Security Context**: Default `runAsNonRoot: true` and `readOnlyRootFilesystem: true` enforced via template validation.
- **Metadata Tagging**: Mandatory labels for `agency-id`, `service-tier`, and `data-classification`.
- **Linting**: Every chart must pass `helm lint` and `kube-linter` before being accepted into the registry.

---

**PROMPT 253 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 254 — NAMESPACE ISOLATION.**
