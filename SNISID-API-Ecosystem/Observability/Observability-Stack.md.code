# SNISID API Observability Stack

## Architecture
Comprehensive monitoring using the LGTM stack (Loki, Grafana, Tempo, Mimir/Prometheus).

## 1. Metrics (Prometheus)
- `api_request_duration_seconds`: Latency tracking.
- `api_requests_total`: Traffic volume and error rates (4xx, 5xx).
- `api_active_connections`: Current load.

## 2. Dashboards (Grafana)
- **National Overview**: Total traffic across all agencies.
- **Agency Deep-dive**: Specific health metrics for ONI, DGI, etc.
- **Security Dashboard**: WAF blocks, failed auth attempts.

## 3. Logs (Loki)
- Centralized collection of all Gateway and Mesh logs.
- Structured logging (JSON) for easy querying.

## 4. Tracing (OpenTelemetry / Tempo)
- Distributed tracing to visualize calls across multiple agencies.
- Correlation IDs injected at the Gateway.

## 5. Alerts
- Slack/Email notifications for:
  - Latency > 2s for critical APIs.
  - Error rate > 1%.
  - Unauthorized access spikes.
