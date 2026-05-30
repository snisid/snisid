# SNISID Security Observability Stack

## 1. Objective
To provide complete visibility into the security state of the national infrastructure, moving beyond simple logs to a full observability model (Metrics, Logs, Traces).

## 2. Monitoring Domains

| Domain | Focus Area | Critical Metrics |
| :--- | :--- | :--- |
| **Threat Detections** | SIEM/EDR performance | Alert volume, False positive rate, Mean time to detect. |
| **Failed Authentications** | IAM Health | Spike in 401/403 errors, account lockout rates. |
| **Privilege Escalations** | Admin activity | Unexpected `sudo` usage, creation of new admin users. |
| **API Anomalies** | Gateway Traffic | Unusual request patterns, spikes in 5xx errors, payload size anomalies. |
| **Network Anomalies** | Traffic flow | Unexpected egress to foreign IPs, internal scanning patterns. |

## 3. Tooling Architecture

| Component | Recommended Tool | Purpose |
| :--- | :--- | :--- |
| **Metrics** | Prometheus / VictoriaMetrics | Time-series data for performance and anomaly spikes. |
| **Logs** | Loki / OpenSearch | Distributed log aggregation and searching. |
| **Tracing** | OpenTelemetry / Jaeger | Tracking requests across microservices to find the "patient zero" of a breach. |
| **Visualization** | Grafana | Unified dashboards for the entire observability stack. |

## 4. The Observability Pipeline
`Telemetry (App/OS/Net)` $\rightarrow$ `OpenTelemetry Collector` $\rightarrow$ `Storage (Prometheus/Loki/OpenSearch)` $\rightarrow$ `Grafana Dashboard`

## 5. Key Dashboards
- **SOC Health Dashboard:** SIEM ingestion rates, active analyst count, alert queue depth.
- **Infrastructure Security Dashboard:** Patch levels across all clusters, open port status.
- **Identity Dashboard:** Global login map, high-risk account activity.
- **Attack Surface Dashboard:** Number of public-facing APIs, DNS changes, certificate expiry.
