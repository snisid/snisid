# PROMPT 271: GIT REPOSITORY STRUCTURE

This architecture defines the multi-repo strategy and branching workflows for the SNISID platform, ensuring maximum isolation, security, and administrative control.

---

## 1. Repo Topology (Separation of Concerns)

SNISID uses a decoupled repository model to prevent "Monolith Exhaustion" and enforce least-privilege access.

- **`snisid-core`**: Primary backend services (Go), AI engines (Python), and core intelligence logic.
- **`snisid-infra`**: Infrastructure as Code (Terraform/OpenTofu) for cloud resources and Kubernetes clusters.
- **`snisid-manifests`**: GitOps repository containing Kubernetes YAML/Helm charts for all environments (ArgoCD target).
- **`snisid-security`**: Centralized repository for Kyverno policies, OPA rules, and security scanning configurations.
- **`snisid-docs`**: Architectural blueprints, operational playbooks, and national compliance documentation.

---

## 2. Branching Workflows (Trunk-Based + Promotion)

### Development Repos (`snisid-core`)
- **`main`**: The source of truth. All features are developed in short-lived feature branches and merged via Pull Request.
- **Release Tags**: Immutable version tags (e.g., `v1.2.3`) trigger the automated build and promotion pipeline.

### GitOps Repos (`snisid-manifests`)
- **`alpha`**: Automatically updated by CI when `main` is updated.
- **`beta`**: Promotion-only branch; requires manual sign-off after Alpha validation.
- **`prod`**: Final production branch; requires cryptographic approval from the Infrastructure Lead and CISO.

---

## 3. Governance Architecture

- **Protected Branches**: Direct pushes to `main`, `beta`, or `prod` are prohibited.
- **Pull Request Requirements**:
    - Minimum of 2 approved reviews.
    - All status checks (CI, Lint, Security) must be green.
    - No merge conflicts.
- **Commit Signing**: Every commit must be cryptographically signed using GPG/SSH keys tied to the officer's national identity.

---

## 4. Security Controls

- **Repository Isolation**: Developers from `agency-intelligence` cannot view or access the `agency-defense` repositories.
- **Secret Scanning**: Gitleaks is integrated into the pre-commit and CI stages to block any accidental secret exposure.
- **Audit Logging**: Every Git operation (clone, push, pull, PR comment) is recorded and streamed to the Forensic Audit Ledger.
- **Dependency Guard**: Automated monitoring of `go.mod` and `requirements.txt` for vulnerable or unauthorized libraries.

---

## 5. Integration Strategy

- **Webhooks**: High-security, encrypted webhooks trigger GitHub Actions and ArgoCD reconciliations.
- **Cross-Repo Sync**: Changes in `snisid-core` (app) automatically trigger updates in `snisid-manifests` (ops) via automated PRs.
- **Mirroring**: For air-gapped deployments, repositories are unidirectionally mirrored to internal Git servers within the sovereign data center.

---

**PROMPT 271 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 272 — ARGO CD CONFIGURATION.**
