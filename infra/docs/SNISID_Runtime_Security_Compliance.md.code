# SNISID: Runtime Security & Compliance

To defend against advanced persistent threats (APTs) and internal misconfigurations, SNISID implements a multi-layered runtime security and compliance-as-code framework.

---

## 1. Kernel-Level Runtime Monitoring (Falco/Tetragon)

We monitor the "Behavioral Integrity" of every container at the syscall level.

- **Falco Rules**: A set of national security rules that detect suspicious events (e.g., a web server spawning a shell, sensitive files like `/etc/shadow` being read, or outbound connections to unauthorized IPs).
- **Tetragon Enforcement**: Using eBPF to not only detect but **instantly kill** any process that violates security boundaries (e.g., an unauthorized attempt to modify the container's filesystem).
- **Regional Security Aggregation**: Security events are collected locally in each region and streamed to the **National SOC** for real-time correlation.

---

## 2. Policy Enforcement as Code (Kyverno/OPA)

Ensuring that no non-compliant resource ever enters the cluster.

- **Admission Control**: **Kyverno** or **OPA Gatekeeper** intercepts all Kubernetes API requests.
- **Mandatory Security Policies**:
  - **No Root Containers**: Every container must run as a non-privileged user.
  - **Immutable Filesystem**: Forcing root filesystems to be read-only for most microservices.
  - **Signed Images Only**: Verifying the **Cosign** signature before allowing a pod to start.
- **Resource Guardrails**: Enforcing memory/CPU limits and quotas for every agency namespace to prevent resource exhaustion attacks.

---

## 3. Compliance & Auditability

- **Continuous Compliance Scans**: Automated tools like **Kubescape** or **Checkov** run daily scans against the running cluster to identify drift from the national security baseline.
- **Sovereign Audit Ledger**: Every policy violation, runtime security alert, and admission decision is cryptographically signed and stored in the immutable ledger.
- **Compliance Dashboards**: Real-time visualization of the platform's compliance posture for the **National Security Council**.

---

## 4. Workload Identity (SPIRE)

- **Identity-Based Authorization**: Moving beyond static secrets to dynamic, short-lived identities for every workload.
- **SPIFFE ID**: Every pod is issued a unique SPIFFE ID that it uses to authenticate with other services, databases, and the **Vault** secret engine.
- **Hardware Root of Trust**: Integrating SPIRE with **TPM/HSM** on the worker nodes to ensure that the workload identity is tied to the physical hardware.

---

## 5. Vulnerability Management

- **Runtime Vulnerability Scanning**: Integration with **Trivy Operator** to continuously scan running pods for newly discovered CVEs.
- **Automated Quarantine**: If a "Critical" vulnerability is discovered in a running pod, the **Autonomous SOC** can automatically trigger a rolling update or isolate the pod until a patch is applied.
