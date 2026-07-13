# 👁 SNISID — Observability Stack

> Phase 3 / Étape 12 — Stack d'observabilité complet.
> CNCF-compatible, on-premise + multi-DC.

## Composants

| Composant | Version | Rôle |
|-----------|---------|------|
| **Prometheus** | 2.55+ | Collecte metrics + recording rules + alertes |
| **Thanos** | 0.36+ | Multi-DC storage, requêtes globales, longue rétention |
| **Alertmanager** | 0.27+ | Routage alertes (PagerDuty, Slack, SMS, Email) |
| **Loki** | 3.x | Logs centralisés |
| **Tempo** | 2.x | Traces (compatible Jaeger / Zipkin / OTLP) |
| **OpenTelemetry Collector** | 0.110+ | Hub metrics/logs/traces, tail sampling |
| **Grafana** | 11.x | Dashboards + Explore + SLO + Alerting UI |

## Arborescence

```
observability/
├── README.md
├── prometheus/
│   ├── prometheus.yml
│   ├── rules/
│   │   ├── slo.rules.yml             # Recording rules (SLO error budget)
│   │   ├── workflow.rules.yml        # Workflows metrics
│   │   └── kafka.rules.yml           # Kafka health
│   └── alerts/
│       ├── workflow.alerts.yml
│       ├── kafka.alerts.yml
│       ├── platform.alerts.yml
│       ├── security.alerts.yml
│       └── slo.alerts.yml
├── alertmanager/
│   └── alertmanager.yml
├── otel/
│   ├── collector.yaml
│   └── exporters.yaml
├── loki/
│   └── loki.yaml
├── tempo/
│   └── tempo.yaml
├── grafana/
│   ├── datasources/datasources.yaml
│   ├── provisioning/dashboards.yaml
│   └── dashboards/
│       ├── slo-national.json
│       ├── workflow-health.json
│       ├── kafka-event-mesh.json
│       └── citizen-experience.json
└── deploy/
    └── docker-compose.observability.yml
```

## Démarrage local (dev)

```bash
docker compose -f deploy/docker-compose.observability.yml up -d
open http://localhost:3000   # Grafana (admin/snisid-dev)
```

## Prod (k8s)

Tous les composants sont déployés via Helm + ArgoCD :
- `prometheus-stack` (kube-prometheus-stack)
- `grafana-loki`
- `grafana-tempo`
- `opentelemetry-collector`
- `thanos`
