# PROMPT 278: AUTOMATED METRICS COLLECTION

This architecture defines the high-precision metrics collection and performance monitoring strategy for the SNISID platform, enabling real-time operational awareness and predictive scaling.

---

## 1. Metrics Architecture (Federated & Scalable)

SNISID utilizes a multi-tier metrics stack to handle millions of time-series across a distributed federation.

- **Local Collector (Prometheus)**: Regional instances that scrape metrics from microservices, Kubernetes nodes, and hardware sensors.
- **Global Aggregator (Thanos)**: Provides a unified, long-term query view across all regional Prometheus instances.
- **Agent (OpenTelemetry)**: Standardized instrumentation for application-level metrics (e.g., request count, latency, custom AI performance counters).
- **TSDB Backend**: Compressed, object-storage based backend (S3/GCS) for efficient historical data retention.

---

## 2. Collection Workflows (Pull & Push)

1.  **Discovery**: Prometheus automatically discovers new pods and services via the Kubernetes API.
2.  **Scraping**: High-frequency scraping (every 10–15 seconds) for mission-critical core services.
3.  **Pushgateway**: Used for short-lived batch jobs (e.g., daily intelligence indexing) that cannot be scraped.
4.  **Enrichment**: Every metric is automatically labeled with its `region`, `agency`, `security_tier`, and `hardware_generation`.

---

## 3. Aggregation Strategy (Intelligent Downsampling)

- **Raw Data (2 Weeks)**: Retained at full resolution for deep forensic debugging.
- **Downsampled (1 Year)**: Aggregated at 1-minute intervals for trend analysis and capacity planning.
- **Historical (Infinite)**: Critical performance KPIs are archived for national multi-year growth reporting.

---

## 4. Security & Privacy

- **mTLS Scrapes**: Prometheus uses mTLS (via Istio) to securely pull metrics from pods, preventing unauthorized data exposure.
- **RBAC Querying**: Grafana users are restricted to viewing metrics for their specific agency and environment.
- **Metric Anonymization**: AI-driven filters redact any sensitive identifiers accidentally exposed in metric labels (e.g., User IDs or IP addresses).

---

## 5. Governance Model

- **Alerting SLIs/SLOs**: Automated alerts are generated if Service Level Indicators (SLIs) deviate from the National Service Level Objectives (SLOs).
- **Audit Ledger**: All changes to alerting rules and recording rules are managed via GitOps and recorded in the forensic ledger.
- **Capacity Forecasting**: AI analyzes historical metrics to predict hardware exhaustion 30 days in advance, triggering automated node pool expansions.

---

**PROMPT 278 IS FULLY ARCHITECTED.**
**READY FOR PROMPT 279 — DISTRIBUTED TRACING SYSTEM.**
