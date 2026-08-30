# PROMPT 296: AUTOMATED INFRASTRUCTURE SECURITY HARDENING

This architecture defines the continuous, automated strategy for hardening the SNISID platform's infrastructure, ensuring that all components (OS, Kubernetes, Network) adhere to the highest national security standards and "Defense-in-Depth" principles.

---

## 1. Hardening Architecture (Multi-Layer Shielding)

SNISID utilizes a proactive hardening stack that enforces security at every layer of the infrastructure.

- **OS Hardening (CIS Benchmarks)**: Automated application of CIS (Center for Internet Security) L1/L2 benchmarks to all physical and virtual nodes via Ansible/Salt.
- **Kernel Shielding (eBPF/LSM)**: Utilizing eBPF and Linux Security Modules (AppArmor/SELinux) to restrict system calls and file access at runtime.
- **Kubernetes Hardening (Kyverno)**: Policy-driven enforcement of Pod Security Standards (Privileged, Baseline, Restricted) across all namespaces.
- **Network Hardening (Cilium)**: Default-deny egress and ingress policies with L7 API filtering to prevent unauthorized lateral movement.

---

## 2. Hardening Workflows (The Compliance Loop)

1.  **Baseline Application**: Every new node or cluster is automatically provisioned with the "Sovereign Hardened Image."
2.  **Continuous Scanning**: Tools like **Kube-bench** and **Kube-hunter** run scheduled jobs to identify deviations from the hardened baseline.
3.  **Auto-Remediation**: If a hardening control is disabled (e.g., "SSH root login enabled"), the system automatically reverts the change and logs it in the forensic ledger.
4.  **Vulnerability Shielding**: When a new kernel or K8s vulnerability is announced, the system automatically deploys temporary eBPF "Virtual Patches" to block the exploit path until a permanent patch is applied (Prompt 289).

---

## 3. Integration Strategy (Hardening-as-Code)

- **Hardened Gold Images**: The automated image promotion pipeline (Prompt 273) only accepts images that have passed a 100% hardening check.
- **Infrastructure-as-Code (IaC) Linting**: Terraform/Crossplane manifests are linted for security misconfigurations (e.g., "Public S3 Bucket", "Open Security Group") before being applied.
- **Dynamic Policy Updates**: Hardening policies are version-controlled in the central security repository and synchronized globally via ArgoCD.

---

## 4. Security & Privacy

- **Minimalist Footprint**: Hardening policies enforce the removal of all unnecessary packages, services, and users from the OS and container images.
- **Forensic Integrity**: Every hardening event (application, scan, remediation) is cryptographically signed and stored in the forensic ledger for long-term verification.
- **Confidential Computing**: Where hardware supports it, hardening policies enforce the use of TEEs (Trusted Execution Environments) for sensitive cryptographic operations.

---

## 5. Governance Model

- **Hardening Drift SLA**: Deviations from the hardened baseline must be remediated in < 15 minutes for Tier-0 infrastructure.
- **Security Exceptions**: Hardening controls can only be waived for specific workloads with a multi-signature approval and a defined expiration date.
- **National Hardening Report**: Real-time dashboard for national security officers showing the "Hardening Coverage" and "Control Effectiveness" across all agency clusters.

---

**PROMPT 296 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 297 — AUTOMATED IDENTITY & ACCESS GOVERNANCE.**
