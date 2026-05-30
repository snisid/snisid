# PROMPT 254: NAMESPACE ISOLATION ARCHITECTURE

This architecture defines the strict isolation boundaries for agencies and services within the SNISID Kubernetes environment.

---

## 1. Namespace Topology (Hierarchical)

SNISID uses **Hierarchical Namespaces (HNC)** to manage agency-level isolation while sharing common administrative overhead.

```
snisid-root/
├── agency-intelligence/
│   ├── auth-service/
│   ├── graph-engine/
│   └── data-ingest/
├── agency-interior/
│   ├── biometric-match/
│   └── census-portal/
└── snisid-system/ (Platform Ops)
    ├── monitoring/
    └── security-mesh/
```

---

## 2. Isolation Policies

### Network Micro-segmentation
- **Default Deny**: Every namespace starts with a `DefaultDeny` NetworkPolicy for all ingress and egress traffic.
- **Service-to-Service**: Only explicitly allowed cross-namespace traffic is permitted (e.g., `agency-intelligence` can talk to `snisid-system/monitoring`).
- **External Traffic**: All external egress must pass through the **Sovereign Egress Gateway**.

### Resource Quotas
- Every agency namespace is assigned a **ResourceQuota** (CPU, RAM, Storage, GPU).
- **LimitRanges** ensure that individual pods do not exceed allocated agency resources or starve the cluster.

---

## 3. Security Workflows

1.  **Agency Onboarding**: An agency sub-root namespace is created via GitOps.
2.  **Policy Injection**: Kyverno automatically applies the mandatory `NetworkPolicy`, `ResourceQuota`, and `RBAC` roles.
3.  **Credential Provisioning**: A dedicated HashiCorp Vault path is created for the agency, accessible only by its service accounts.

---

## 4. Governance Controls (RBAC Enforcement)

- **Least Privilege**: Agency administrators are granted `NamespaceAdmin` only within their sub-tree.
- **Platform Admin**: Only the central **Sovereign Infrastructure Team** has cluster-wide `admin` rights.
- **Audit**: All RBAC changes and cross-namespace access requests are logged to the forensic audit ledger.

---

## 5. Runtime Enforcement Strategy

- **Admission Control**: Kyverno blocks any pod that attempts to use `hostNetwork`, `hostPID`, or mount sensitive host paths.
- **Identity Isolation**: SPIFFE/SPIRE ensures that a pod in `agency-intelligence` cannot masquerade as a pod in `agency-interior`.
- **Egress Filtering**: FQDN-based egress filtering ensures that agencies can only connect to authorized national endpoints.

---

**PROMPT 254 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 255 — AUTOSCALING SYSTEM.**
