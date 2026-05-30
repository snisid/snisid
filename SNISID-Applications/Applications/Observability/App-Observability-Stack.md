# APP OBSERVABILITY STACK — SNISID
## Stack d'Observabilité Applicative Nationale

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-OBS-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |

---

## 1. PRÉSENTATION

Stack complète d'observabilité pour monitorer en temps réel toutes les applications SNISID : crashes, échecs de synchronisation, latence API, anomalies de sécurité et expérience utilisateur.

---

## 2. DOMAINES MONITORÉS

### 2.1 App Crashes

| Métrique | Source | Seuil | Alerte |
|----------|--------|-------|--------|
| Crash Rate | Sentry / Firebase Crashlytics | < 0.1% | > 0.5% |
| ANR Rate | Android vitals | < 0.05% | > 0.1% |
| OOM Rate | Memory tracking | < 0.01% | > 0.05% |
| Crash-free Users | Session tracking | > 99.5% | < 99% |
| Time to Crash | Duration tracking | > 30 min | < 10 min |

### 2.2 Offline Sync Failures

| Métrique | Seuil | Alerte |
|----------|-------|--------|
| Sync Success Rate | > 98% | < 95% |
| Queue Backlog | < 1000 | > 5000 |
| Sync Latency (P0) | < 1 min | > 5 min |
| Conflict Rate | < 1% | > 5% |
| Data Loss Events | 0 | > 0 |

### 2.3 API Latency

| Métrique | p50 | p95 | p99 | Alerte |
|----------|-----|-----|-----|--------|
| Identity API | < 200ms | < 500ms | < 1s | > 2s p99 |
| Civil Registry API | < 300ms | < 800ms | < 2s | > 3s p99 |
| Sync API | < 500ms | < 2s | < 5s | > 10s p99 |
| Notification API | < 100ms | < 300ms | < 1s | > 2s p99 |
| Case Management | < 300ms | < 1s | < 2s | > 5s p99 |

### 2.4 Security Anomalies

| Détection | Source | Action |
|-----------|--------|--------|
| Brute Force | Login attempts | Auto-block IP |
| Unusual Geolocation | Geo-IP | MFA challenge |
| Device Spoofing | Attestation fails | Block + Alert |
| API Abuse | Rate limit exceeded | Throttle + Alert |
| Data Exfiltration | Volume anomaly | Block + Investigate |
| Tamper Attempts | Integrity check | Alert + Log |
| Anomalous Time | Clock drift > 5 min | Reject + Alert |

### 2.5 User Experience

| Métrique | Source | Cible |
|----------|--------|-------|
| App Load Time | RUM | < 2s |
| Screen Transitions | RUM | < 300ms |
| User Flow Completion | Analytics | > 90% |
| Error Rate | RUM | < 2% |
| Session Duration | Analytics | > 5 min |
| Daily Active Users | Analytics | Tracking |
| Feature Usage | Analytics | Per feature |

---

## 3. OUTILS

### 3.1 Metrics — Prometheus

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'snapid-apps'
    static_configs:
      - targets:
        - 'api-gateway:9090'
        - 'identity-service:9090'
        - 'sync-service:9090'
    metrics_path: '/metrics'
    scheme: 'https'
    tls_config:
      cert_file: /etc/prometheus/certs/client.crt
      key_file: /etc/prometheus/certs/client.key
  
  - job_name: 'snapid-mobile'
    metrics_path: '/mobile-metrics'
    scrape_interval: 5m  # Mobile push metrics less frequently
```

### 3.2 Mobile Analytics — OpenTelemetry

```yaml
# otel-collector-config.yaml
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
    send_batch_size: 1024
  attributes:
    actions:
      - key: environment
        value: production
        action: insert
  memory_limiter:
    check_interval: 1s
    limit_mib: 512

exporters:
  prometheus:
    endpoint: 0.0.0.0:8889
    namespace: snapid
  loki:
    endpoint: http://loki:3100/loki/api/v1/push
  tempo:
    endpoint: tempo:4317

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [tempo]
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch, attributes]
      exporters: [prometheus]
    logs:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [loki]
```

### 3.3 Logs — Loki

```
┌──────────────────────────────────────────────┐
│              LOG LEVELS                      │
├──────────────────────────────────────────────┤
│  INFO    ─── Normal operations               │
│  WARN    ─── Non-critical issues            │
│  ERROR   ─── Application errors            │
│  CRITICAL ── System failures                │
│  SECURITY ── Security events (immutable)    │
│  AUDIT   ─── Audit trail (immutable)        │
└──────────────────────────────────────────────┘
```

### 3.4 Dashboards — Grafana

| Dashboard | Description | Refresh |
|-----------|-------------|---------|
| **App Health** | Global app health, uptime, crashes | 30s |
| **API Performance** | Latency, throughput, errors | 30s |
| **Sync Status** | Offline sync metrics | 1m |
| **Security** | Anomalies, threats, attacks | 10s |
| **User Experience** | RUM, satisfaction, adoption | 5m |
| **Business KPIs** | Usage, adoption, performance | 15m |

---

## 4. ALERTING

### 4.1 Alert Rules

```yaml
groups:
  - name: app_alerts
    rules:
      - alert: HighCrashRate
        expr: rate(app_crashes_total[5m]) > 0.005
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "App crash rate above 0.5%"
          
      - alert: SyncFailureHigh
        expr: sync_failure_rate > 0.05
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Sync failure rate above 5%"
          
      - alert: APILatencyHigh
        expr: histogram_quantile(0.95, rate(api_latency_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "API p95 latency above 1s"
          
      - alert: SecurityAnomaly
        expr: security_anomaly_count > 10
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Multiple security anomalies detected"
```

### 4.2 Notification Channels

| Priorité | Canal | Destinataire |
|----------|-------|-------------|
| **P0** | SMS + Phone + Slack | On-call engineer |
| **P1** | Phone + Slack | Engineering team |
| **P2** | Slack + Email | Team lead |
| **P3** | Email | Report |

---

## 5. DASHBOARD GRAFANA PRINCIPAL

```
┌──────────────────────────────────────────────────────────┐
│  🇭🇹 SNISID App Observability  │  🟢 All Systems Go      │
├──────────────────────────────────────────────────────────┤
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │
│  │  Total   │ │  Crash   │ │  Sync    │ │  API     │   │
│  │  Users   │ │   Rate   │ │  Health  │ │  Latency │   │
│  │  234,567 │ │  0.08%   │ │  99.2%   │ │  245ms   │   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘   │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │  API Latency (p95) — Last 24h                    │   │
│  │  ╱╲    ╱╲                                        │   │
│  │ ╱  ╲  ╱  ╲    ╱╲                                │   │
│  │╱    ╲╱    ╲  ╱  ╲                               │   │
│  │           ╲╱    ╲                               │   │
│  └──────────────────────────────────────────────────┘   │
│                                                          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                │
│  │  Crashing │ │   Sync   │ │ Security │                │
│  │  Apps    │ │  Queues  │ │  Events  │                │
│  │  Android:3│ │  891 ✓   │ │   2 ⚠   │                │
│  │  iOS:1   │ │  12 ✗    │ │   0 🔴  │                │
│  └──────────┘ └──────────┘ └──────────┘                │
├──────────────────────────────────────────────────────────┤
│  Time: 24h │ 7d │ 30d  │  Refresh: 30s                  │
└──────────────────────────────────────────────────────────┘
```

---

## 6. SLAs

| Métrique | SLA | Mesure |
|----------|-----|--------|
| App Uptime | > 99.9% | Crash-free sessions |
| API Availability | > 99.95% | HTTP 200/5xx |
| Sync Success | > 99% | Successful/total syncs |
| Alert Response (P0) | < 15 min | Time to acknowledge |
| Alert Response (P1) | < 1h | Time to acknowledge |
| MTTR (Critical) | < 2h | Time to resolve |

---
*Fin du document — App Observability Stack v1.0*