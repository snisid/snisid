# PROMPT 266: FULL CI/CD PIPELINE ARCHITECTURE

This architecture defines the end-to-end CI/CD lifecycle for the SNISID platform, using a secure, GitOps-driven model with GitHub Actions and ArgoCD.

---

## 1. CI/CD Topology (GitOps-Centric)

SNISID utilizes a decoupled CI and CD model to ensure maximum security and auditability.

- **CI Layer (GitHub Actions)**: Handles code validation, testing, security scanning, and container image builds.
- **CD Layer (ArgoCD)**: Handles the declarative deployment of infrastructure and applications to Kubernetes.
- **Artifact Registry**: OCI-compliant, air-gap-ready repository for signed container images and Helm charts.
- **Config Repo**: A dedicated, encrypted Git repository containing the environment-specific manifests.

---

## 2. Deployment Workflows (Multi-Environment)

1.  **Code Commit**: Developer pushes code to a feature branch.
2.  **Continuous Integration**:
    - **Build**: Compiles Go/Python code.
    - **Test**: Executes unit and integration tests.
    - **Security**: Runs SAST (Semgrep) and Dependency Scanning.
    - **Artifact Generation**: Builds and signs the container image (Cosign).
3.  **Promotion**: Upon successful merge to `main`, the CI pipeline updates the image tag in the **Config Repo**.
4.  **Continuous Deployment**:
    - **Sync**: ArgoCD detects the manifest change.
    - **Validation**: Kyverno validates the new manifest against cluster security policies.
    - **Execution**: ArgoCD applies the change to the **Alpha (Capital)** region first.
    - **Global Rollout**: After 2 hours of healthy metrics, the change is propagated to Beta and Gamma regions.

---

## 3. Validation Pipelines (Automated Gates)

Every deployment is gated by **Analysis Runs**:

- **Smoke Tests**: Verifies basic API connectivity.
- **Integration Tests**: Executes cross-service workflows in a transient namespace.
- **Security Gating**: Fails the pipeline if any "High" or "Critical" vulnerabilities are detected in the image SBOM.
- **Compliance Gating**: Checks if the deployment violates any national data residency or resource quota rules.

---

## 4. Security Controls (Government-Grade)

- **Cryptographic Signing**: Every image must be signed; unsigned images are blocked by the cluster admission controller.
- **OIDC Identity**: GitHub Actions uses OIDC (Keyless) to authenticate with the cloud provider and Vault, eliminating static secrets.
- **Least Privilege**: The CI pipeline has "Push-only" access to the registry; the CD pipeline has "Read-only" access to the registry and "Admin" access only to its dedicated namespace.
- **Audit Ledger**: Every pipeline execution, including the specific commit hash and approval logs, is stored in the immutable forensic audit trail.

---

## 5. Recovery Mechanisms (Rollback)

- **Git-Based Rollback**: Reverting a commit in the Config Repo automatically triggers ArgoCD to roll back the cluster to the previous healthy state.
- **Automated Health Rollback**: If post-deployment monitoring detects a failure, the pipeline issues an automated `kubectl rollout undo`.
- **DR Restoration**: In the event of total regional failure, the CI/CD system can re-provision the entire national stack in a new region from the GitOps source of truth.

---

**PROMPT 266 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 267 — AUTOMATED BUILD SYSTEM.**
