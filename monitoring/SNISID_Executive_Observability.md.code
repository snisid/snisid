# SNISID Executive Observability: Configurations & Scripts

**Classification:** RESTRICTED / SOVEREIGN OBSERVABILITY
**Compliance:** NIST SP 800-137 / ISO 27001 / SLA Governance

This operational playbook defines the OpenTelemetry (OTel) configurations, Prometheus PromQL metrics mapping, Grafana panels, and Python prediction analytics scripts for the SNISID Executive Observability Platform.

---

## 1. OpenTelemetry Collector configuration

This configuration processes incoming telemetry streams, stripping out sensitive citizen identifiers (PII) at the edge before metrics are committed to Prometheus storage.

```yaml
# File: /opt/snisid/monitoring/otel-collector-config.yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    timeout: 1s
    send_batch_size: 256
    
  # Redaction Filter: Remove names, biometrics, or plain-text national IDs
  attributes/redact_pii:
    actions:
      - key: citizen.name
        action: delete
      - key: citizen.biometric.vector
        action: delete
      - key: citizen.nni
        action: hash # Convert to cryptographic hash hash-id for tracking without disclosure
        
  memory_limiter:
    check_interval: 1s
    limit_percentage: 85
    spike_limit_percentage: 15

exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"
    namespace: "snisid_executive"
  loki:
    endpoint: "http://loki.monitoring.svc.cluster.local:3100/loki/api/v1/push"

service:
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, attributes/redact_pii, batch]
      exporters: [prometheus]
    logs:
      receivers: [otlp]
      processors: [memory_limiter, attributes/redact_pii, batch]
      exporters: [loki]
```

---

## 2. Prometheus PromQL KPI Query Map

These PromQL queries calculate the operational, citizen, fraud, infrastructure, and security SLA indexes displayed on the executive console.

### 2.1. Daily Registration Throughput Rate
```promql
# Measures citizen registration counts aggregated over a 24-hour rate
sum(increase(snisid_identity_registrations_total[24h]))
```

### 2.2. Database Replication Lag (SLA Control)
```promql
# Monitors synchronous CockroachDB replica lag. Threshold: Alert if lag >= 2.0 seconds
max(cockroach_range_replication_latency_seconds{job="cockroachdb"})
```

### 2.3. Active Portal Citizen User Sessions
```promql
# Active sessions in Keycloak auth provider
sum(keycloak_active_user_sessions{realm="citizen"})
```

### 2.4. Biometric Duplication Hits (Fraud Alert)
```promql
# Rate of biometric duplicate locks triggered by the ABIS engine per minute
sum(rate(snisid_abis_matches_rejected_total{reason="duplicate"}[5m])) * 60
```

### 2.5. API HTTP 5xx Error Rate (SLA Control)
```promql
# Calculates the percentage of failing inter-agency API calls. Threshold: Alert if > 0.5%
(sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))) * 100
```

### 2.6. Inter-Agency Request Latency (p95 percentile)
```promql
# Measures inter-agency API response latency. Threshold: Alert if p95 latency >= 250ms
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))
```

---

## 3. Predictive Capacity & AI Anomaly Script

This Python script runs as a cron job, analyzing historical Prometheus disk space and load parameters to predict capacity limits, automatically generating high-priority SLA warnings.

```python
# File: /opt/snisid/monitoring/predict_capacity.py
import sys
import numpy as np

def predict_exhaustion(time_series_data, threshold_percent=90.0):
    """
    Analyzes historical data points using linear regression to estimate 
    when the metrics will breach the defined warning threshold.
    Time series input is a list of tuples: (timestamp_seconds, value)
    """
    print(f"[*] Analyzing {len(time_series_data)} historical data points...")
    
    if len(time_series_data) < 5:
        print("[-] Insufficient data points for regression modeling.")
        return None
        
    times = np.array([point[0] for point in time_series_data])
    values = np.array([point[1] for point in time_series_data])
    
    # Normalize times to offset from start
    start_time = times[0]
    relative_times = times - start_time
    
    # Fit Linear Regression line: y = m*x + c
    slope, intercept = np.polyfit(relative_times, values, 1)
    
    print(f"[*] Computed Trajectory: Slope={slope:.6f}, Intercept={intercept:.2f}")
    
    # Check if capacity is growing (positive slope)
    if slope <= 0:
        print("[+] Resource usage is stable or declining. No exhaust danger.")
        return -1
        
    # Calculate relative time to cross threshold
    relative_time_to_exhaust = (threshold_percent - intercept) / slope
    exhaust_timestamp = start_time + relative_time_to_exhaust
    
    time_remaining_hours = relative_time_to_exhaust / 3600.0
    print(f"[*] Projected Exhaustion Time remaining: {time_remaining_hours:.2f} hours")
    
    # Trigger alert if resource exhaustion is projected in less than 6 hours
    if time_remaining_hours <= 6.0:
        print(f"[!] WARNING: Resource exhaust predicted in {time_remaining_hours:.2f} hours!")
        return int(exhaust_timestamp)
        
    return -1

if __name__ == "__main__":
    # Test execution with mock storage capacity increase: 80% to 88% over 8 hours
    mock_data = [
        (0, 80.0),
        (3600, 81.0),
        (7200, 82.0),
        (10800, 83.0),
        (14400, 84.0),
        (18000, 85.0),
        (21600, 86.0),
        (25200, 87.0),
        (28800, 88.0)
    ]
    # Execute prediction check
    predict_exhaustion(mock_data)
```

---

## 4. Grafana Panel Alert Integration

This panel JSON snippet is deployed inside the Grafana configuration to alert the SOC when the inter-agency API latency violates SLA thresholds:

```json
{
  "title": "Inter-Agency API Latency SLA Alert",
  "type": "graph",
  "targets": [
    {
      "expr": "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))",
      "legendFormat": "p95 Latency",
      "refId": "A"
    }
  ],
  "alert": {
    "conditions": [
      {
        "evaluator": {
          "params": [0.25],
          "type": "gt"
        },
        "operator": {
          "type": "and"
        },
        "query": {
          "params": ["A", "5m", "now"]
        },
        "reducer": {
          "type": "avg"
        },
        "type": "query"
      }
    ],
    "executionErrorState": "alerting",
    "frequency": "60s",
    "handler": 1,
    "name": "API p95 Latency SLA Breach Warning",
    "noDataState": "no_data",
    "notifications": [
      {
        "uid": "snisid-soc-slack"
      }
    ]
  }
}
```

---

*Verified and signed by the SNISID Observability Operations Board.*
