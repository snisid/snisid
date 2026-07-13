# PROMPT 256: POD SECURITY ARCHITECTURE

This architecture defines the multi-layered security controls for all workloads running within the SNISID Kubernetes environment.

---

## 1. Security Architecture (Workload Hardening)

SNISID implements a **"Defense-in-Depth"** model for pod security, moving beyond basic Pod Security Standards.

- **Immutable Filesystem**: All pods are configured with `readOnlyRootFilesystem: true`.
- **Non-Root Execution**: `runAsNonRoot: true` is enforced globally; no container may run with UID 0.
- **Capability Restriction**: `allowPrivilegeEscalation: false` and removal of all default Linux capabilities (e.g., `CAP_NET_RAW`).

---

## 2. Admission Workflows (Policy Gating)

Every pod deployment must pass through the **Kyverno Admission Controller** before being admitted to the cluster.

1.  **Validation**: Kyverno checks the manifest against the "Sovereign Workload Policy".
2.  **Mutation**: Automatically injects security sidecars (Istio, SPIRE) and security context defaults if missing.
3.  **Verification**: Validates that the container image is signed by the SNISID CI system and has a valid SBOM.
4.  **Rejection**: Any manifest that violates security boundaries (e.g., attempting to mount `/etc`) is blocked with a forensic alert.

---

## 3. Policy Enforcement Model

- **Standard Policies**: Baseline and Restricted profiles from Kubernetes Pod Security Standards.
- **Sovereign Policies**: Custom Rego/Kyverno rules specific to national intelligence (e.g., "Agency A cannot share a node with Agency B").
- **Drift Detection**: Automated scanning of running pods to ensure they haven't been mutated at runtime to bypass security controls.

---

## 4. Runtime Protection Strategy (eBPF-Based)

SNISID uses **Tetragon** (eBPF-based) for real-time runtime enforcement.

- **Process Filtering**: Blocks execution of unauthorized binaries (e.g., `curl` or `nc`) within production pods.
- **System Call Monitoring**: Detects and prevents privilege escalation attempts at the kernel level.
- **Network Observability**: Real-time correlation of network traffic to specific process IDs and pod identities.

---

## 5. Governance Framework

- **Vulnerability Management**: Continuous scanning of all running images; pods are automatically evicted if a "Critical" CVE is discovered.
- **Identity Integrity**: SPIFFE identities are rotated hourly to minimize the window for credential theft.
- **Compliance Reporting**: Weekly automated reports summarizing policy violations, blocked admissions, and runtime security events for the National Security Oversight.

---

**PROMPT 256 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 257 — ISTIO SERVICE MESH.**
