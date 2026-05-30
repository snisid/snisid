# PROMPT 279: DISTRIBUTED TRACING SYSTEM

This architecture defines the end-to-end distributed tracing strategy for the SNISID platform, enabling deep visibility into cross-service communication and bottleneck identification within the national intelligence mesh.

---

## 1. Tracing Architecture (OpenTelemetry Native)

SNISID utilizes a standardized tracing stack based on the OpenTelemetry (OTel) ecosystem.

- **SDK (OpenTelemetry)**: Integrated into all microservices (Go, Python, Java) to generate spans and propagate trace contexts via B3 or W3C headers.
- **Collector (OTel Collector)**: Deployed as sidecars or DaemonSets to receive, process, and export spans to the regional backend.
- **Backend (Tempo/Jaeger)**: High-scale, cost-effective trace storage integrated with the regional object storage.
- **Visualization (Grafana)**: Integrated with Loki (logs) and Prometheus (metrics) to provide "exemplars" (links from metrics/logs directly to the relevant trace).

---

## 2. Instrumentation Workflows (Auto & Manual)

1.  **Service Mesh (Istio)**: Automatically generates spans for all L7 (HTTP/gRPC) traffic between microservices without requiring code changes.
2.  **SDK Instrumentation**: Developers use OTel SDKs to add "Manual Spans" for complex internal logic (e.g., AI model inference time, database transaction duration).
3.  **Context Propagation**: Ensures that the `trace_id` is passed through Kafka messages, ensuring that a single user request can be followed across asynchronous boundaries.

---

## 3. Sampling Strategy (Intelligent & Adaptive)

To handle national-scale traffic without overwhelming storage, SNISID uses a **Tail-Based Sampling** model:

- **Success Responses**: 1% sampling rate for standard successful requests.
- **Error Responses**: 100% sampling rate for any request that results in a 5xx error.
- **High-Latency**: 100% sampling for any request that exceeds the P95 latency threshold.
- **Security-Tiered**: 100% sampling for all requests involving the `national-vault` or `identity-core`.

---

## 4. Security & Privacy

- **Span Scrubbing**: The OTel Collector automatically redacts sensitive data (e.g., authentication tokens, PII, query parameters) from span attributes.
- **Encrypted Export**: Traces are transmitted via TLS to the regional backend.
- **Trace RBAC**: Only authorized SREs and Security Officers can view traces; visibility is restricted by the `agency` label on the spans.

---

## 5. Governance Model

- **Standardized Span Names**: Enforced naming convention for all spans to ensure consistent searching and aggregation.
- **Audit Ledger**: All trace queries are logged in the forensic ledger to prevent unauthorized "User Journey" mining.
- **Performance Budgeting**: AI analyzes trace data to identify services that are consistently over their "Latency Budget" and automatically creates Jira/GitLab issues for the development team.

---

**PROMPT 279 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 280 — ALERTING & NOTIFICATION SYSTEM.**
