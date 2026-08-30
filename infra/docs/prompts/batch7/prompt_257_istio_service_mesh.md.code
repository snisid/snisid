# PROMPT 257: ISTIO SERVICE MESH ARCHITECTURE

This architecture defines the service mesh strategy for SNISID, providing secure communication, traffic management, and deep observability across the national intelligence ecosystem.

---

## 1. Mesh Topology (Distributed Control Plane)

SNISID utilizes **Istio** in a multi-primary, multi-region configuration to ensure high availability and regional sovereignty.

- **Istiod (Primary)**: Deployed in each regional cluster to manage local proxy configurations and certificate issuance.
- **Sidecar Proxies (Envoy)**: Automatically injected into every application pod to intercept and secure all L4/L7 traffic.
- **Cross-Region Mesh**: Regional meshes are federated via **Istio East-West Gateways**, allowing secure cross-cluster service discovery.

---

## 2. Traffic Workflows

- **Internal (East-West)**: A service in `agency-intelligence` calling a service in `snisid-system` is automatically upgraded to **mTLS**.
- **Ingress (North-South)**: External traffic enters via the **Istio Ingress Gateway**, where it is authenticated (JWT) and authorized (RBAC) before reaching the target service.
- **Traffic Shifting**: Supports **Blue-Green** and **Canary** deployments by modifying Istio `VirtualService` weights in real-time.

---

## 3. Security Architecture (STRICT mTLS)

- **PeerAuthentication**: Configured in `STRICT` mode globally; any service attempting a plain-text connection is instantly rejected.
- **AuthorizationPolicies**: Fine-grained, identity-based access control (e.g., "Only the `fraud-engine` can call the `citizen-ledger`").
- **RequestAuthentication**: Validates JWT tokens from the National Identity Provider for every end-user request.

---

## 4. Observability Pipelines

Istio provides out-of-the-box observability for the entire SNISID network:

- **Metrics**: Standard L7 metrics (request count, error rate, latency) streamed to **Prometheus/Thanos**.
- **Distributed Tracing**: Service-to-service correlation IDs propagated to **Tempo/Jaeger** for deep request analysis.
- **Service Graph**: Visual representation of service dependencies and health via **Kiali**.

---

## 5. Runtime Governance Strategy

- **Policy Gating**: All mesh changes (VirtualServices, DestinationRules) must pass through a CI/CD validation pipeline before application.
- **Failover Logic**: Automatic circuit breaking and outlier detection to eject unhealthy pods from the load balancing pool.
- **Egress Control**: Strict egress gateways ensure that internal services can only communicate with pre-approved external national API endpoints.

---

**PROMPT 257 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 258 — INGRESS/EGRESS GATEWAY.**
