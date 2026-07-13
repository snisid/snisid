# SNISID: Zero Trust Network Security

This document defines the network security model for SNISID, ensuring that every packet within the national infrastructure is authenticated, authorized, and inspected according to **Zero Trust** principles.

---

## 1. Pod-Level Micro-Segmentation

We implement strict network isolation using **Kubernetes NetworkPolicies** and **Istio AuthorizationPolicies**.

- **Default Deny**: All namespaces are configured with a default "Deny All" ingress and egress policy.
- **Service-to-Service Whitelisting**: Connections are only permitted if explicitly defined in a policy (e.g., `identity-engine` can only egress to `vault` and `kafka`).
- **Namespace Boundaries**: Enforcing strict boundaries between agencies. The `agency-finance` namespace cannot reach pods in `agency-intelligence` unless through a secured and audited gateway.

---

## 2. Secure Ingress Gateway

The Ingress Gateway is the "Digital Border" for all incoming traffic to the national cluster.

- **WAF & DDoS Protection**: All ingress traffic passes through an integrated **Web Application Firewall (WAF)** and DDoS mitigation layer.
- **mTLS Termination**: The gateway handles mTLS handshake with external agencies and then initiates a new mTLS connection to the internal service mesh.
- **Identity-Aware Routing**: Routing decisions are made based on the cryptographically verified identity of the requester, not just the IP address.

---

## 3. Secure Egress Controls

Egress traffic is strictly controlled to prevent data exfiltration and unauthorized communication with external services.

- **Egress Gateway**: All outbound traffic must pass through a dedicated **Egress Gateway** pod.
- **FQDN Whitelisting**: Only traffic to specific, pre-approved Fully Qualified Domain Names (e.g., `api.interpol.int`) is permitted.
- **Deep Packet Inspection (DPI)**: Egress traffic is inspected for anomalous patterns or PII-like data before leaving the sovereign network.

---

## 4. Istio Mesh Security

- **PeerAuthentication**: Enforcing `STRICT` mTLS across the entire mesh, ensuring that no unencrypted traffic is permitted.
- **RequestAuthentication**: Validating JWT tokens for every incoming request to ensure the user/agent is authenticated.
- **Service Mesh Observability**: Real-time visualization of the network graph (using **Kiali**) to detect unauthorized communication attempts.

---

## 5. Runtime Threat Detection

- **Network Anomaly Detection**: Integration with **Falco** or **Tetragon** to detect suspicious network behavior at the kernel level (e.g., a pod attempting to scan the local subnet).
- **Automated Containment**: If a network threat is detected, the **Autonomous SOC** instantly injects a "Quarantine Policy" to isolate the compromised pod from the rest of the mesh.
- **Audit Logs**: 100% of network flow logs (accepted and denied) are streamed to the **Sovereign Audit Ledger**.
