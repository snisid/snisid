# SNISID: Multi-Agency Cloud Orchestration & Workload Identity

This document defines the framework for managing multiple government agencies within the shared SNISID infrastructure and securing their communications using hardware-backed workload identities.

---

## 1. Multi-Agency Tenant Isolation

To ensure sovereignty and security, each agency is treated as a first-class tenant with cryptographic and resource boundaries.

- **Hierarchical Namespaces**: Each agency is assigned a parent namespace (e.g., `agency-dhs`) containing sub-namespaces for specific projects (e.g., `agency-dhs-biometrics`).
- **Cryptographic Isolation**: Each agency namespace has its own dedicated **Vault Transit Engine** for data encryption and its own **KMS key alias**.
- **Resource Guardrails**: Strict **ResourceQuotas** and **LimitRanges** prevent one agency from exhausting the cluster's GPU or memory resources.
- **Agency-Specific Ingress**: Dedicated ingress gateways per agency, allowing for localized WAF rules and audit trails.

---

## 2. Workload Identity (SPIFFE/SPIRE)

We move beyond static Kubernetes secrets to dynamic, short-lived, and verifiable identities.

- **SPIFFE ID Issuance**: Every pod is automatically issued a SPIFFE identity (e.g., `spiffe://snisid.gov/ns/snisid-core/sa/identity-matcher`).
- **Hardware Root of Trust**: The **SPIRE Agent** running on each node validates the pod's identity using node attestation (TPM/HSM) and workload attestation (Unix UID, K8s namespace).
- **Identity-Aware Communication**: Services use their SPIFFE identities to establish mTLS connections, ensuring that identity is verified at the kernel/hardware level, not just the network level.
- **Seamless Vault Integration**: Pods use their SPIRE-issued SVID (SVID) to authenticate with Vault, eliminating the need for long-lived bootstrap tokens.

---

## 3. National Edge Computing (Border Kiosks)

Managing infrastructure at the "National Edge" requires a different orchestration model.

- **K3s/KubeEdge**: Lightweight Kubernetes distributions deployed at border crossings, airports, and regional offices.
- **Autonomous Edge Operation**: Edge clusters are designed to operate offline for up to 48 hours, using local biometric caches and local AI inference if the central connection is lost.
- **Secure Tunneling**: Edge clusters connect to the central national cluster via encrypted **WireGuard** or **Istio Multi-Cluster** tunnels.
- **Centralized Policy Push**: Security policies and AI model updates are pushed from the central GitOps repo to all edge clusters simultaneously.

---

## 4. Infrastructure Audit Ledger

- **System-Level Traceability**: Every orchestration event (e.g., namespace creation, SPIFFE ID issuance, resource quota change) is logged.
- **Immutable Evidence**: These logs are cryptographically hashed and streamed to the **Sovereign Audit Ledger**, providing forensic evidence of the infrastructure's state at any point in time.
- **Tamper Detection**: Real-time monitoring of the audit ledger to detect unauthorized modifications to the cluster's control plane.
