# PROMPT 258: INGRESS/EGRESS GATEWAY ARCHITECTURE

This architecture defines the secure boundary controls for all North-South and outbound traffic within the SNISID sovereign environment.

---

## 1. Gateway Topology (Regional High-Availability)

SNISID uses dedicated **Istio Ingress and Egress Gateways** in each regional cluster, fronted by a national-level **GSLB**.

- **Sovereign Ingress Gateway**: Entry point for all external requests (citizens, other agencies).
- **Sovereign Egress Gateway**: Unified exit point for all internal services connecting to external national APIs.
- **WAF Layer**: Every ingress gateway is fronted by an **Envoy-based WAF** (ModSecurity/Coraza) for real-time L7 protection.

---

## 2. Routing Workflows

### Ingress Path
1.  **L3/L4 Protection**: National Anti-DDoS scrubs traffic before it hits the gateway IP.
2.  **TLS Termination**: Gateway terminates TLS using **HSM-backed certificates**.
3.  **Authentication**: Requests must contain a valid JWT; otherwise, they are redirected to the National Auth Portal.
4.  **Routing**: Istio `VirtualService` routes the request to the target agency namespace based on the Host header or URI path.

### Egress Path
1.  **Identification**: Internal pods connect to an internal FQDN (e.g., `bank-api.national.internal`).
2.  **Redirection**: Traffic is transparently captured and routed to the **Egress Gateway**.
3.  **Filtering**: Gateway checks the destination against a whitelist; unauthorized outbound connections are dropped and alerted as "Data Exfiltration Attempt".

---

## 3. Security Enforcement Architecture

- **Identity-Aware Routing**: The gateway validates the user's `agency_claim` in the JWT to ensure they can only reach authorized endpoints.
- **Rate Limiting**: Per-agency and per-user rate limits enforced at the gateway to prevent resource exhaustion.
- **DDoS Mitigation**: Automated circuit breaking at the gateway level if abnormal traffic patterns are detected.

---

## 4. Observability Integration

- **Access Logs**: Every request (including headers and TLS version) is logged and streamed to the **Sovereign SOC**.
- **Traffic Analytics**: Real-time dashboards showing geolocation, status codes (4xx/5xx), and bandwidth usage per agency.
- **Anomaly Detection**: SIEM integration to correlate gateway logs with internal pod events for threat hunting.

---

## 5. Resilience Model

- **Multi-AZ Gateways**: Gateway pods are spread across multiple availability zones within a region.
- **Health-Based Failover**: GSLB automatically reroutes traffic to the backup region if the regional Ingress Gateway fails its health check.
- **Graceful Termination**: Envoy proxies ensure that in-flight requests are completed before a gateway pod is rotated during an update.

---

**PROMPT 258 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 259 — CLUSTER NETWORK SECURITY.**
