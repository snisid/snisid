# SNISID: Sovereign CI/CD Pipelines

The Sovereign CI/CD Pipelines provide the "Automated Trust" layer for SNISID, ensuring that every line of code is verified, scanned, and cryptographically signed before it enters the national production environment.

---

## 1. Air-Gap Ready Architecture (GitLab/Tekton)

The pipeline is designed to operate within high-security, sovereign network zones.

- **Self-Hosted GitLab**: All source code, issue tracking, and CI runners are hosted within the national data centers.
- **Tekton CD**: Kubernetes-native CI/CD that runs as pods within the cluster, allowing for direct integration with the **Zero Trust Mesh**.
- **Binary Mirroring**: External dependencies (libraries, base images) are mirrored into a local **Sovereign Artifact Registry (Harbor)** and scanned before being made available to developers.

---

## 2. Multi-Stage Security Gating

Every build must pass through a series of "Hard Gates" to proceed.

- **SAST (Static Analysis)**: Scanning source code for vulnerabilities and hardcoded secrets.
- **SCA (Software Composition Analysis)**: Identifying vulnerable open-source dependencies (e.g., using **Trivy** or **Snyk**).
- **DAST (Dynamic Analysis)**: Automated security testing against a running staging environment.
- **Policy Check (Kyverno)**: Verifying that the generated Kubernetes manifests conform to national security standards (e.g., no root access, restricted capabilities).

---

## 3. Binary Provenance & Image Signing

- **Cosign Integration**: Every container image built by the pipeline is cryptographically signed using a key stored in the **National HSM**.
- **Admission Control Enforcement**: The production cluster is configured with a **Sigstore Admission Controller** that blocks the execution of any image not signed by the official SNISID build key.
- **Sovereign SBOM**: Generating a Software Bill of Materials (SBOM) for every release, providing a complete inventory of all components in the production environment.

---

## 4. Automated Deployment & Rollback

- **ArgoCD Integration**: The CI pipeline updates the GitOps repository, triggering an automated reconciliation by ArgoCD.
- **Canary Deployments**: Using **Istio/Flagger** to automatically shift a small percentage of traffic (e.g., 5%) to a new version and monitoring health metrics before full rollout.
- **Automated Rollback**: If the error rate increases during a canary rollout, Flagger automatically reverts the deployment to the previous stable version.

---

## 5. Audit Ledger Integration

- **Pipeline Telemetry**: Every build, scan result, and deployment event is logged to the **Sovereign Audit Ledger**.
- **Cryptographic Evidence**: The audit ledger stores the hashes of the source code, the build environment, and the final artifact, ensuring full traceability of the software supply chain.
