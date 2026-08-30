# PROMPT 260: CONFIG MANAGEMENT SYSTEM

This architecture defines the unified configuration and secrets management strategy for the SNISID ecosystem, ensuring security, consistency, and GitOps integration.

---

## 1. Config Architecture (Layered)

SNISID uses a decoupled configuration model to separate application logic from environmental state.

- **Immutable Base Configs**: Hardcoded in container images (e.g., internal directory structures).
- **Environmental ConfigMaps**: Managed via GitOps for non-sensitive settings (e.g., log levels, API endpoints).
- **Dynamic Secrets**: Managed via **HashiCorp Vault** for sensitive data (e.g., DB credentials, API keys, TLS certs).

---

## 2. Secret Workflows (Vault Integration)

SNISID implements a **"No-Static-Secrets"** policy.

1.  **Secret Creation**: Security officers define secrets in HashiCorp Vault via Terraform.
2.  **Access Control**: Kubernetes ServiceAccounts are granted specific Vault roles using **Kubernetes Auth Method**.
3.  **Secret Injection**: The **Secrets Store CSI Driver** (or Vault Agent Sidecar) mounts the secrets as a volume or injects them as environment variables into the pod.
4.  **Rotation**: Vault automatically rotates credentials (e.g., database passwords) every 24 hours without service restart.

---

## 3. Governance Model

- **Config Integrity**: All ConfigMaps and Secret definitions are stored in an encrypted GitOps repository.
- **Peer Review**: Every configuration change requires a mandatory "Double-Blind" review process.
- **Version Control**: Every config change is tagged and can be rolled back instantly via GitOps.
- **Audit Ledger**: All secret access and config mutations are logged to the forensic audit trail.

---

## 4. Synchronization Strategy (GitOps)

- **Source of Truth**: The `snisid-infra` Git repository.
- **Reconciliation**: **ArgoCD** continuously monitors the repository and synchronizes the cluster state.
- **Manual Override Protection**: Any manual `kubectl edit configmap` is automatically reverted by ArgoCD within 3 minutes to prevent configuration drift.

---

## 5. Security Controls

- **Encryption at Rest**: All `ConfigMap` and `Secret` data in ETCD is encrypted using **HSM-backed keys**.
- **Namespace Isolation**: Secrets in `agency-intelligence` are cryptographically isolated and inaccessible to `agency-interior`.
- **EBPF Monitoring**: Tetragon monitors for any unauthorized attempts to read `/etc/secrets` or access the Vault API endpoint from unintended processes.

---

**PROMPT 260 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 261 — ROLLING UPDATE STRATEGY.**
