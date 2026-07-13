# SNISID: Kubernetes Production Architecture

The Kubernetes Production Architecture provides the "Sovereign Cloud" for SNISID, ensuring that all microservices, AI workloads, and data pipelines run on a resilient, secure, and horizontally scalable platform.

---

## 1. Multi-Region Federated Topology

SNISID is deployed across multiple sovereign geographic regions to ensure survival even during a regional disaster.

- **Federated Control Plane**: Uses **Karmada** or **Clusternet** to orchestrate workloads across independent regional clusters.
- **Node Pools**:
  - **General Purpose**: For standard microservices (API Gateway, Identity Engine).
  - **GPU Accelerated**: Optimized for AI inference and training (NVIDIA A100/H100).
  - **High-Memory**: For stateful stream processing (Flink, RocksDB).
  - **Hardened Edge**: Specialized nodes for Border Intelligence Kiosks.

---

## 2. Zero Trust Service Mesh (Istio)

All communication within the cluster is governed by a **Zero Trust** model.

- **Mandatory mTLS**: Every pod-to-pod connection is encrypted and authenticated via Istio.
- **Micro-Segmentation**: Strict **AuthorizationPolicies** ensure that only authorized services can communicate (e.g., the Border API can talk to the ABIS, but not to the National Pension DB).
- **Traffic Orchestration**: Canary rollouts and circuit breakers are managed at the mesh level to prevent cascading failures.

---

## 3. GitOps & Infra-as-Code

Infrastructure is entirely immutable and managed via **GitOps**.

- **Source of Truth**: All K8s manifests, Helm charts, and Istio policies are stored in a secure national Git repository.
- **Continuous Reconciliation**: **ArgoCD** or **Flux** monitors the Git repo and automatically applies changes to the clusters, reverting any manual `kubectl` modifications.
- **Enterprise Helm Charts**: Modular charts with strictly defined values for `Production`, `Staging`, and `DR` (Disaster Recovery) environments.

---

## 4. Self-Healing & Resilience

- **Pod Disruption Budgets (PDB)**: Ensuring that critical services (e.g., Auth) always have a minimum number of healthy replicas.
- **Auto-Healing Nodes**: Integration with cloud-native node repairers to automatically replace unhealthy worker nodes.
- **Liveness/Readiness Probes**: Granular health checks for every container to ensure traffic is only routed to ready pods.

---

## 5. Security & Isolation

- **Namespace Isolation**: Each agency or service category (e.g., `agency-border`, `core-identity`) is isolated in a dedicated namespace with strict ResourceQuotas and RBAC.
- **Pod Security Standards**: Enforcing the `Restricted` profile across all namespaces, prohibiting privileged containers and root access.
- **Runtime Security**: Integration with **Falco** to detect anomalous behavior (e.g., unexpected shell execution) within a running container.

---

## 6. Config & Secret Management

- **External Secrets Operator**: Syncing secrets from **HashiCorp Vault** directly into Kubernetes Secret objects.
- **Dynamic Configuration**: Using ConfigMaps for non-sensitive runtime parameters, with automated pod restarts upon update via Reloader.
- **Audit Logs**: 100% of Kubernetes API calls are logged to the **Sovereign Audit Ledger**.
