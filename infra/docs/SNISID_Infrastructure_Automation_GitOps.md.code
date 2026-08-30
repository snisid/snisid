# SNISID: Infrastructure Automation & GitOps

This document defines the automation framework for SNISID, ensuring that the national infrastructure is treated as code (IaC) and managed through a high-assurance **GitOps** pipeline.

---

## 1. Enterprise-Grade Helm Architecture

All SNISID services are packaged as Helm charts to ensure consistency and modularity.

- **Modular Chart Design**: A "Base Chart" defines standard Kubernetes resources (Deployments, Services, Ingress), while "App Charts" inherit from the base and define specific parameters.
- **Environment Overrides**: Using `values.yaml` for shared defaults and environment-specific files (e.g., `values-prod-region-a.yaml`) for regional tuning.
- **Dependency Management**: Centralized Helm repository for shared dependencies (e.g., common sidecars, logging configurations).

---

## 2. GitOps Workflow (ArgoCD)

The cluster state is continuously reconciled with the Git repository.

- **App-of-Apps Pattern**: A master ArgoCD application manages a set of regional infrastructure and application groups.
- **Automated Sync & Pruning**: ArgoCD automatically applies Git changes and "Prunes" any resources that have been deleted from Git.
- **Sync Windows**: Restricting automated changes to specific "Governance Windows" approved by the National Security Council.
- **Self-Healing**: If a pod's configuration is manually altered via the CLI, ArgoCD detects the "Out of Sync" state and instantly reverts it.

---

## 3. Secret Management (HashiCorp Vault)

Secrets are never stored in Git, even in encrypted form.

- **Vault Integration**: Using the **External Secrets Operator (ESO)** or **Vault Agent Injector** to pull secrets into Kubernetes.
- **Dynamic Secrets**: Vault generates short-lived credentials for databases and Kafka, reducing the impact of a potential credential leak.
- **Identity-Based Access**: Pods use their **Kubernetes ServiceAccount** to authenticate with Vault, ensuring that only authorized workloads can access specific secrets.

---

## 4. Multi-Cluster Federation (Karmada)

To manage the national scale, regional clusters are federated into a single logical entity.

- **Global Policy Propagation**: Security policies (e.g., "Ban all root containers") are defined at the federation level and pushed to all regional clusters.
- **Cross-Cluster Failover**: If a regional cluster becomes unreachable, the federation layer automatically migrates workloads to a healthy region.
- **Unified Traffic Orchestration**: Coordinating Global Load Balancers (GSLB) to route citizen traffic to the nearest healthy cluster.

---

## 5. Auditability & Compliance

- **Infrastructure Audit Ledger**: Every Git commit and ArgoCD sync event is logged and cryptographically signed.
- **Compliance-as-Code**: Using **Kyverno** or **OPA Gatekeeper** to validate that all manifests in the Git repo conform to national security standards before they can be merged.
- **Drift Reporting**: Real-time alerts if the cluster state deviates from the Git source of truth for more than 60 seconds.
