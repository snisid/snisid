# PROMPT 270: INFRASTRUCTURE AS CODE (IAC) MANAGEMENT

This architecture defines the modular, secure, and declarative Infrastructure as Code (IaC) strategy for the SNISID sovereign cloud.

---

## 1. IaC Topology (Modular Architecture)

SNISID utilizes **Terraform/OpenTofu** organized into reusable, high-level modules to ensure consistency across regions.

- **Layer 0 (Bootstrap)**: Identity (IAM), State Storage (S3/GCS), and Networking (VPC/VNET).
- **Layer 1 (Kubernetes)**: Cluster provisioning (EKS/GKE/Hardened K8s), Node Pools, and OIDC configuration.
- **Layer 2 (System Services)**: Shared infrastructure components like HashiCorp Vault, Kafka clusters, and Postgres instances.
- **Layer 3 (Workloads)**: Application-specific infrastructure (S3 buckets, dedicated IAM roles).

---

## 2. Provisioning Workflows (GitOps Driven)

1.  **Code Change**: Infrastructure engineer modifies a Terraform module in the `snisid-iac` repository.
2.  **Plan**: CI pipeline executes `terraform plan` and posts the output as a PR comment for review.
3.  **Validation**: Automated checks (Checkov/TFSec) ensure the code adheres to national security baselines.
4.  **Approval**: Cryptographic sign-off by two authorized senior architects.
5.  **Apply**: CI pipeline executes `terraform apply` using a short-lived OIDC token.

---

## 3. State Management Strategy

- **Remote State**: Terraform state is stored in a centralized, encrypted, and versioned bucket with object locking.
- **State Locking**: **DynamoDB/Redis** is used for state locking to prevent concurrent modifications and potential corruption.
- **Regional Isolation**: Each regional cluster has its own independent state file to minimize the "Blast Radius" of infrastructure changes.

---

## 4. Governance Model

- **Immutable Infrastructure**: Changes are never made manually via the cloud console; the GitOps repository is the single source of truth.
- **Drift Detection**: Automated jobs run `terraform plan` every 4 hours; any drift is alerted to the SOC.
- **Resource Tagging**: Mandatory tagging policy enforced by code (e.g., `agency`, `environment`, `security_tier`).

---

## 5. Compliance Automation (Policy-as-Code)

- **Sentinel/OPA**: Fine-grained policy enforcement (e.g., "Prohibit public S3 buckets" or "Require GPU nodes to use specific AMIs").
- **Cost Controls**: Automated estimation of cloud costs for every PR, with mandatory approval if the monthly budget is exceeded.
- **Disaster Recovery**: IaC modules are designed to be "Region-Agnostic" to allow rapid re-provisioning of the entire national stack in a different provider or on-premise hardware.

---

**PROMPT 270 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 271 — GIT REPOSITORY STRUCTURE.**
