# PROMPT 273: AUTOMATED IMAGE PROMOTION

This architecture defines the secure lifecycle and promotion strategy for container images within the SNISID ecosystem, ensuring that only validated and authorized artifacts reach the production workers.

---

## 1. Promotion Architecture (Registry-to-Registry)

SNISID utilizes a multi-stage registry model to enforce environmental isolation.

- **Dev Registry (Local)**: Untrusted, volatile images used for feature branch development.
- **Stage Registry (Alpha/Beta)**: Signed images that have passed all CI/CD unit and integration tests.
- **Prod Registry (National)**: Highly restricted, immutable images authorized for national-scale production.
- **Image Mirror**: Automated replication of authorized images across all 251–265 regional clusters to ensure high availability during a national network partition.

---

## 2. Registry Workflows (Secure Handoff)

1.  **Build**: The CI pipeline builds and signs an image with a unique SHA.
2.  **Tagging**: Images are tagged with their source Git hash and the environment (e.g., `agency-v1.2.3-alpha`).
3.  **Promotion Trigger**: Successful completion of the "Alpha Validation" pipeline triggers an automated PR to move the image to the `beta` branch of the manifest repo.
4.  **National Promotion**: After passing the "Beta Load Test," the image is promoted to the `prod` registry. This move involves:
    - **Vulnerability Re-scan**: Final check for any CVEs discovered since the original build.
    - **Cryptographic Re-sign**: Image is re-signed with the "National Production Key."

---

## 3. Validation Pipelines (Promotion Gates)

Every promotion step is guarded by a **Promotion Gate**:

- **Security Gate**: Any image with a "High" CVE discovered after the build is blocked from promotion.
- **Integrity Gate**: Verification that the image signature matches the build-time provenance.
- **Documentation Gate**: Checks for the existence of an updated `CHANGELOG.md` and `README.md` in the source repository.

---

## 4. Governance Model (Approval & Audit)

- **Least Privilege**: Only the "Promotion Service Account" has write access to the Prod Registry.
- **Approval Logic**:
    - **Alpha -> Beta**: Automated based on test results.
    - **Beta -> Prod**: Requires explicit manual approval in the Config Repo PR.
- **Audit Ledger**: Every promotion event (who, when, which SHA, from where to where) is recorded in the forensic audit trail.

---

## 5. Integration Strategy (GitOps)

- **ArgoCD Image Updater**: Automatically monitors the registries and updates the Kubernetes manifests in Git when a new authorized image is available for a specific environment.
- **SBOM Sovereignty**: The SBOM for the promoted image is automatically uploaded to the national security database for real-time dependency tracking across the federation.
- **Registry Cleanup**: Automated policy to delete images from the Dev registry that have not been promoted to Stage within 14 days.

---

**PROMPT 273 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 274 — AUTOMATED ROLLBACK SYSTEM.**
