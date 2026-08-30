# PROMPT 259: CLUSTER NETWORK SECURITY MODEL

This architecture defines the internal network security and micro-segmentation strategy for the SNISID Kubernetes clusters, ensuring that lateral movement is impossible even in the event of a pod compromise.

---

## 1. Network Topology (Zero Trust Micro-segmentation)

SNISID uses **Cilium** as the CNI (Container Network Interface) to provide high-performance, eBPF-powered network security.

- **Micro-segmentation**: Traffic is restricted at the pod level based on identity, not just IP address.
- **East-West Inspection**: All traffic between services is inspected and logged using eBPF, providing deep visibility without the overhead of sidecars for basic L4 filtering.
- **Global Trust Domain**: Unified network identity across federated clusters.

---

## 2. Segmentation Workflows

1.  **Identity Creation**: Every pod is assigned a **Cilium Identity** based on its Kubernetes labels and namespace.
2.  **Policy Definition**: Developers define `CiliumNetworkPolicy` or standard `NetworkPolicy` objects as part of their service's GitOps manifest.
3.  **Automatic Enforcement**: The Cilium agent on each node intercepts traffic and applies policies before the packet even leaves the pod's network namespace.

---

## 3. Enforcement Architecture

- **Default Deny All**: A global policy ensures that pods cannot communicate unless an explicit allow-list rule exists.
- **Namespace Boundaries**: Cross-agency traffic is blocked by default; only authorized "Shared Services" (like Vault or Kafka) can receive traffic from multiple agency namespaces.
- **Identity-Based Filtering**: Rules are written as `fromEndpoints: [{ "agency": "intelligence" }]`, making policies resilient to pod IP changes.

---

## 4. Monitoring Pipelines (Real-time Analytics)

- **Cilium Hubble**: Provides real-time observability into every network connection, flow, and dropped packet.
- **Flow Logging**: All network flows are streamed to the **Sovereign SIEM** for forensic analysis.
- **Threat Detection**: Integrated with **Falco** to detect abnormal network behavior (e.g., a web service suddenly attempting to scan the internal network).

---

## 5. Recovery Mechanisms (Automated Containment)

- **Quarantine Policy**: If a pod is flagged as compromised, an automated workflow applies a `QuarantineNetworkPolicy` that severs all its network connections.
- **Dynamic Policy Update**: Security teams can push global "Kill-switch" policies across the entire national federation in seconds using Karmada.
- **Self-Healing Connectivity**: If a CNI agent fails, the node is automatically marked as `Unhealthy` and drained to prevent a security bypass.

---

**PROMPT 259 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 260 — CONFIG MANAGEMENT SYSTEM.**
