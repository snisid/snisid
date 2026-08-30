# 👁 WORKFLOW OBSERVABILITY MODEL

> **Phase 3 / Étape 12** — Visibilité temps réel sur l'administration nationale.
> Version : 1.0.0

---

## 1. Objectif

> **On ne gouverne que ce que l'on mesure.**

Le modèle d'observabilité SNISID permet de :
- **Voir** chaque workflow, chaque événement, chaque transition.
- **Mesurer** les SLI vs SLO en temps réel.
- **Alerter** avant que le citoyen ne subisse l'incident.
- **Tracer** une demande de bout en bout (citoyen → mairie → ONI → tribunal → carte).
- **Auditer** sans rien pouvoir falsifier.

---

## 2. Les 4 Piliers (Three Pillars + Audit)

| Pilier | Question | Outil |
|--------|----------|-------|
| **Metrics** | « Quoi se passe ? » | **Prometheus + Thanos** |
| **Logs** | « Pourquoi ça se passe ? » | **Loki** |
| **Traces** | « Où ça se passe (chemin) ? » | **Tempo + OpenTelemetry** |
| **Audit** | « Qui a fait quoi, quand, avec quelle preuve ? » | **Kafka audit.* + WORM** |

Toutes les 4 sources convergent dans **Grafana** (cockpit unique).

---

## 3. Domaines à monitorer

| Domaine | Indicateurs principaux |
|---------|------------------------|
| **BPMN failures** | `workflow_failures_total`, `incident_count` |
| **Kafka lag** | `kafka_consumergroup_lag`, `kafka_under_replicated_partitions` |
| **SLA breaches** | `workflow_sla_breach_total`, `workflow_duration_seconds` |
| **Fraud anomalies** | `fraud_score`, `fraud_detected_total`, `fraud_cases_open` |
| **Workflow bottlenecks** | `user_task_age_seconds`, `job_backoff_seconds` |
| **Platform** | CPU/RAM/disk, latence DC1↔DC2, état PKI/TSA |
| **Sécurité** | `auth_failed_total`, `policy_violations_total` |
| **Citoyen** | `citizen_request_duration_seconds`, `citizen_appeal_count` |

---

## 4. Architecture Observabilité

```
┌──────────────────────────────────────────────────────────────────────┐
│                            APPLICATIONS                              │
│  (Workflow Engine, Zeebe, Temporal, Kafka, Microservices, Edge)      │
└─────────────┬───────────────────────┬──────────────────────┬────────┘
              │ metrics (Prometheus)  │ logs (OTLP)          │ traces (OTLP)
              ▼                       ▼                      ▼
        ┌──────────┐         ┌─────────────────────────────────────┐
        │Prometheus│         │     OpenTelemetry Collector         │
        │  + Thanos│         │   (filter, batch, tail-sampling)    │
        └────┬─────┘         └──────────┬──────────────┬──────────┘
             │                          │              │
             │                       ┌──▼──┐        ┌──▼──┐
             │                       │ Loki│        │Tempo│
             │                       └──┬──┘        └──┬──┘
             ▼                          │              │
        ┌──────────┐                    │              │
        │Alertmgr  │                    │              │
        └────┬─────┘                    │              │
             │                          │              │
             ▼                          │              │
  PagerDuty / Slack / SMS / Email       │              │
                                        │              │
                                        ▼              ▼
                                  ┌─────────────────────────┐
                                  │       GRAFANA           │
                                  │ (dashboards + explore)  │
                                  └─────────────────────────┘
                                        ▲
                                        │
                                  ┌─────┴─────┐
                                  │ Kafka     │
                                  │ audit.*   │  (lake + WORM)
                                  └───────────┘
```

---

## 5. Conventions de Nommage (Metrics)

Préfixe : `snisid_`

| Métrique | Type | Labels |
|----------|------|--------|
| `snisid_workflow_started_total` | counter | `workflow_id`, `version`, `region` |
| `snisid_workflow_completed_total` | counter | `workflow_id`, `version`, `region`, `result` |
| `snisid_workflow_failures_total` | counter | `workflow_id`, `version`, `region`, `error_type` |
| `snisid_workflow_duration_seconds` | histogram | `workflow_id`, `version`, `region` |
| `snisid_workflow_sla_breach_total` | counter | `workflow_id`, `severity` |
| `snisid_workflow_active` | gauge | `workflow_id`, `version` |
| `snisid_user_task_age_seconds` | histogram | `workflow_id`, `task_id`, `assignee_group` |
| `snisid_fraud_score` | histogram | `workflow_id`, `decision` |
| `snisid_fraud_detected_total` | counter | `workflow_id`, `action` |
| `snisid_pki_sign_duration_seconds` | histogram | `purpose` |
| `snisid_kafka_publish_total` | counter | `topic` |
| `snisid_kafka_publish_failures_total` | counter | `topic` |
| `snisid_citizen_request_duration_seconds` | histogram | `channel`, `domain` |

> Toutes les métriques exposées sur `/metrics` (port 9464) et scrappées par Prometheus.

---

## 6. SLO Definitions (Error Budget)

Implémentés via recording rules (cf. `prometheus/rules/slo.rules.yml`).

Exemples :
- `slo:civil_birth_simple:success_ratio:30d ≥ 0.999`
- `slo:identity_verification_online:p99_seconds:30d ≤ 30`
- `slo:kafka_mesh:availability:30d ≥ 0.9995`

---

## 7. Tracing (W3C Trace Context)

- Chaque requête citoyen génère un `traceId`.
- Propagé sur :
  - REST/gRPC headers `traceparent`
  - Kafka headers `trace-id`
  - Camunda variables `traceId`
  - Audit events
- **Sampling** :
  - 100 % pour workflows CRITIQUES
  - Tail-based pour workflows MEDIUM (erreurs/lenteurs 100 %, succès rapides 1 %)

---

## 8. Logs

| Source | Format | Destination |
|--------|--------|-------------|
| Workflow Engine | JSON structuré (pino) | Loki via OTel |
| Zeebe | JSON | Loki |
| Temporal | JSON | Loki |
| Kafka brokers | JSON | Loki |
| Microservices | JSON | Loki |

Champs obligatoires :
- `timestamp` (ISO 8601 + ms)
- `level` (`debug|info|warn|error|fatal`)
- `service`, `version`, `region`
- `traceId`, `spanId`
- `workflowId`, `workflowInstanceId` (si applicable)
- `subjectKey` (NIN ou autre, **toujours pseudonymisé** dans Loki)

Rétention :
- INFO : 30 jours
- WARN/ERROR : 90 jours
- Audit (Kafka) : 30 ans (WORM séparé)

---

## 9. Dashboards Grafana (4 dashboards de référence)

| Dashboard | Audience | Usage |
|-----------|----------|-------|
| **SNISID — SLO National** | Direction / WGO | KPI nationaux, error budget |
| **SNISID — Workflow Health** | Astreinte / SRE | Workflows, incidents, taux d'échec |
| **SNISID — Kafka Event Mesh** | Platform | Brokers, topics, consumers, lag |
| **SNISID — Citizen Experience** | Direction / Communication | Délais perçus, satisfaction, appels |

Voir `observability/grafana/dashboards/*.json`.

---

## 10. Alertes (Severity)

| Severity | Latence cible notification | Canal |
|----------|---------------------------|-------|
| **Sev0** (crise nationale) | < 30 s | PagerDuty + SMS direction + Slack `#snisid-crisis` |
| **Sev1** (service critique HS) | < 1 min | PagerDuty + Slack `#snisid-incidents` |
| **Sev2** (dégradation) | < 5 min | Slack + Email astreinte |
| **Sev3** (information) | < 1 h | Email digest |

Routes définies dans `observability/alertmanager/alertmanager.yml`.

---

## 11. Audit (au-delà des 3 piliers)

- Kafka `audit.*` (3 topics) → ingéré par :
  - **Loki** (consultation + corrélation)
  - **WORM S3-compliant** (preuve légale, immuable, 30 ans)
  - **Cold storage** chiffré (DC3) pour réquisitions judiciaires
- Hash chaîné Merkle → toute altération détectable
- Audit cross-référencé avec traces (`traceId` commun)

---

## 12. Pratiques d'utilisation (RACI)

| Acteur | Action |
|--------|--------|
| **SRE / Astreinte** | Dashboards Workflow Health + alertes Sev1/2 |
| **WGO** | SLO National + Workflow Health quotidien |
| **Direction** | Citizen Experience + SLO mensuel |
| **Cyber** | Security + audit anomalies |
| **DPO** | Vérification non-export PII hors Loki/Tempo |
| **Cour des Comptes** | Audit annuel WORM |

---

## 13. Standards & Conformité

- **OpenTelemetry** (CNCF) pour metrics/logs/traces
- **Prometheus** (CNCF) pour metrics
- **Loki / Tempo** (Grafana Labs)
- **W3C Trace Context** pour propagation
- **OpenMetrics** pour exposition
- **ISO 27001** logs sécurité 1 an minimum
- **Code de la santé / Code civil HT** pour rétention audit (30 ans)

---

**Maintenu par :** Platform Engineering + WGO + DPO
