# 🔭 NATIONAL OBSERVABILITY ANALYTICS STACK

> **Objectif** : Voir la santé analytique nationale en permanence.

---

## 1. PÉRIMÈTRE

| Domaine | Monitoring |
|---------|------------|
| Data pipeline failures | ✅ |
| Model drift | ✅ |
| Analytics latency | ✅ |
| Dashboard availability | ✅ |

---

## 2. STACK

| Domaine | Outil |
|---------|-------|
| Metrics | **Prometheus** + Thanos (long-term) |
| Logs | **Loki** |
| Traces | **Tempo** |
| Visualization | **Grafana** |
| Alerting | Alertmanager + PagerDuty/OpsGenie souverain |
| Synthetic monitoring | Blackbox exporter + k6 |
| Profiling | Pyroscope |

---

## 3. MÉTRIQUES CLÉS

### Pipelines

| Métrique | Description |
|----------|-------------|
| `airflow_dag_failures_total` | Échecs DAG par nom |
| `spark_job_duration_seconds` | Durée jobs Spark |
| `flink_job_restart_count` | Redémarrages Flink |
| `kafka_consumer_lag_records` | Retard consommateurs |
| `lakehouse_ingest_rows_total` | Volume ingéré |

### Modèles

| Métrique | Description |
|----------|-------------|
| `model_inference_latency_seconds` | Latence inférence |
| `model_inference_errors_total` | Erreurs serving |
| `model_drift_score` | Score drift (Evidently) |
| `model_prediction_distribution` | Distribution outputs |
| `feature_freshness_seconds` | Fraîcheur features |

### Analytics

| Métrique | Description |
|----------|-------------|
| `superset_query_duration_seconds` | Latence dashboards |
| `trino_query_failures_total` | Échecs requêtes |
| `dashboard_availability_ratio` | Dispo cockpits (probe) |
| `dq_score{dataset}` | Score qualité courant |

---

## 4. ALERTES PROMETHEUS (extraits)

```yaml
groups:
- name: snisid_analytics_alerts
  rules:
  - alert: PipelineFailureHigh
    expr: increase(airflow_dag_failures_total[1h]) > 3
    for: 10m
    labels: {severity: warning, team: data-eng}
    annotations:
      summary: "DAG {{ $labels.dag_id }} en échec répété"
      runbook: "runbooks/pipeline-failure-recovery.md"

  - alert: KafkaConsumerLagCritical
    expr: kafka_consumer_lag_records > 1000000
    for: 5m
    labels: {severity: critical}
    annotations:
      summary: "Lag Kafka critique sur {{ $labels.topic }}"

  - alert: ModelDriftDetected
    expr: model_drift_score > 0.3
    for: 30m
    labels: {severity: warning, team: ds}
    annotations:
      summary: "Drift détecté sur {{ $labels.model }}"
      runbook: "runbooks/model-rollback.md"

  - alert: DashboardDown
    expr: probe_success{job="superset_probe"} == 0
    for: 2m
    labels: {severity: critical}
    annotations:
      summary: "Dashboard Superset KO"

  - alert: DataQualityLow
    expr: dq_score < 80
    for: 15m
    labels: {severity: warning}
    annotations:
      summary: "DQ score {{ $labels.dataset }} = {{ $value }}"
```

---

## 5. LOGS LOKI

Étiquettes standardisées :
```
{service, env, region, severity, trace_id, dataset, model}
```

Requête type :
```logql
{service="flink", env="prod"} |= "FATAL" | json | severity="critical"
```

Rétention : 30 jours hot, 1 an cold (S3 souverain).

---

## 6. TRACES TEMPO

Instrumentation OpenTelemetry sur :
- API analytics
- Pipelines Spark/Flink
- Inférence modèles KServe
- Requêtes Trino/Superset

Trace complète d'une décision : Frontend → API → Trino → Lakehouse → Cache.

---

## 7. DASHBOARDS OBSERVABILITY

| Dashboard | Audience |
|-----------|----------|
| Pipeline Health | Data Engineering |
| ML Models Health | Data Science |
| Lakehouse Storage | Plateforme |
| BI / Dashboards SLO | Ops |
| End-to-End Latency | Architecture |
| Cost / Capacity | FinOps |

---

## 8. SYNTHETIC MONITORING

Tests probes (k6) toutes les minutes :
- Login cockpit présidentiel
- Chargement dashboard ministériel
- Requête fraude scoring API
- Recherche identité service

Échec probe → alerte critique immédiate.

---

## 9. CAPACITÉ & RÉTENTION

| Donnée | Hot | Cold |
|--------|-----|------|
| Metrics | 30j (Prom) | 2 ans (Thanos S3) |
| Logs | 30j | 1 an |
| Traces | 14j | 90j |
| Synthetic | 90j | 1 an |

---

## 10. PRINCIPE

> **Sans observabilité, pas de souveraineté analytique.**
> Chaque pipeline / modèle / dashboard est observable ou n'est pas en production.
