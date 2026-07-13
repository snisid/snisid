# SNISID: OPA Policy Decision Engine Architecture

Open Policy Agent (OPA) is the centralized "brain" of the SNISID platform. It decouples security and business logic from application code, providing a unified, high-performance engine for making all authorization, trust, and risk decisions.

---

## 1. OPA Architecture: The PDP/PEP/PAP Model

SNISID implements a distributed architecture to ensure sub-5ms latency and high availability.

```mermaid
graph TD
    subgraph "Policy Administration Point (PAP)"
        Git[(Git Repo: Rego Policies)]
        CI[CI/CD: OPA Test & Sign]
        BundleServer[OPA Bundle Server]
        Git --> CI --> BundleServer
    end

    subgraph "Policy Decision Point (PDP)"
        OPA_Sidecar[OPA Sidecar (Per Pod)]
        OPA_Cluster[Regional OPA Cluster]
        BundleServer -->|Pull Bundles| OPA_Sidecar
        BundleServer -->|Pull Bundles| OPA_Cluster
    end

    subgraph "Policy Enforcement Point (PEP)"
        Envoy[Istio Envoy Proxy]
        Gateway[API Gateway]
        K8s_Adm[K8s Admission Controller]
        
        Envoy -->|Query| OPA_Sidecar
        Gateway -->|Query| OPA_Cluster
        K8s_Adm -->|Query| OPA_Cluster
    end
```

---

## 2. Policy Lifecycle & GitOps Workflow

All security policies follow a strict software development lifecycle (SDLC).

1. **Authoring**: Security Architects write policies in **Rego**.
2. **Testing**: Every commit triggers `opa test` to validate logic against baseline mocks.
3. **Verification**: Policies must be cryptographically signed by an authorized security principal.
4. **Distribution**: Policies are packaged into **OPA Bundles** and published to the central Bundle Server.
5. **Deployment**: OPA sidecars automatically pull the latest bundles, ensuring zero-downtime policy updates across the national mesh.

---

## 3. Decision Pipelines & Audit Logging

Every OPA decision is a high-fidelity audit event.

- **Decision Logger**: OPA is configured to push raw decision logs (JSON) to a local sidecar agent.
- **Kafka Stream**: The agent forwards logs to the `audit.authz.decisions` topic in Kafka.
- **Masking**: Sensitive PII attributes (e.g., specific citizen names) are masked at the OPA level before being logged to ensure privacy compliance.

---

## 4. Multi-Tenant Policy Isolation

To support multiple agencies (Tax, Police, Immigration), OPA enforces strict logical isolation.

- **Scoped Bundles**: Each agency has its own Git repository and OPA bundle.
- **Namespace Isolation**: OPA sidecars in the `tax-ns` only pull the Tax Authority bundle.
- **Rego Namespacing**: Policies are structured by agency: `package snisid.authz.dgi`.

---

## 5. Enforcement Workflows

### 5.1. Istio Integration (Envoy ExtAuthz)
Envoy proxies are configured with the `ext_authz` filter, which intercepts every L7 request and performs a gRPC call to the local OPA sidecar.

### 5.2. Kubernetes Integration (Gatekeeper)
OPA Gatekeeper enforces policies at the infrastructure layer (e.g., "All pods must have a `security-zone` label" or "Egress is blocked for the AI-service namespace").

---

## 6. Security Hardening Strategy

- **Unix Domain Sockets (UDS)**: PEP-to-PDP communication uses UDS to eliminate network overhead and prevent man-in-the-middle attacks within the pod.
- **Resource Constraints**: OPA sidecars are strictly limited in CPU/Memory to prevent "Policy Bomb" denial-of-service attacks.
- **Read-Only Storage**: Policy bundles are loaded into memory; the OPA sidecar has no write access to the underlying filesystem.
- **mTLS**: Communication with the central Bundle Server is mandated over TLS 1.3 with hardware-attested certificates.

---

## 7. Performance & Scalability

- **Sub-5ms Latency**: Local sidecar deployment ensures that authorization checks do not impact user experience.
- **Scale-Out**: Regional OPA clusters handle heavy ingress traffic for the API Gateway, scaling horizontally based on request volume.
- **Partial Evaluation**: OPA is optimized to skip non-relevant policy blocks, ensuring constant-time evaluation even as the policy library grows.
