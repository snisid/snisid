---
# ============================================================
# SNISID-Infra — National Telemetry & Observability
# Monitoring, Alertes Prédictives et Tableaux de Bord
# Document ID: SNISID-OBS-001
# Version: 1.0.0
# ============================================================

## 1. STRATÉGIE D'OBSERVABILITÉ COMPLÈTE

Le NOC (National Operations Center) ne peut pas gérer un incident sans visibilité ("You can't fix what you can't see"). La stack d'observabilité collecte les métriques (Metrics), les journaux (Logs) et les traces (Traces) de l'ensemble de l'infrastructure.

## 2. LA STACK "PLG" (Prometheus, Loki, Grafana)

| Composant | Rôle | Source des données |
|-----------|------|--------------------|
| **Prometheus** | Stockage Séries Temporelles (Metrics) | Node Exporters (Serveurs), kube-state-metrics (K8s), PDU/UPS (SNMP) |
| **Loki** | Agrégation de Logs | Promtail (Logs système), FluentBit (Logs applicatifs) |
| **Tempo** | Traces distribuées | OpenTelemetry (Istio, API Gateway) |
| **Grafana** | Visualisation et Alerting | Prometheus, Loki, Tempo |

## 3. ALERTE PRÉDICTIVE (Predictive Alerting)

Plutôt que d'attendre qu'un disque dur soit plein à 100% (causant un arrêt de production), Prometheus utilise des régressions linéaires pour prédire la panne.
```yaml
# Prometheus Alerting Rule: Predictive Disk Fill
alert: PredictiveDiskFull
expr: predict_linear(node_filesystem_free_bytes{mountpoint="/"}[1h], 4 * 3600) < 0
for: 5m
labels:
  severity: warning
annotations:
  summary: "Disk will be full in 4 hours on {{ $labels.instance }}"
```

## 4. INFRASTRUCTURE TELEMETRY (Couche Physique)

Le monitoring ne s'arrête pas au logiciel. Le système surveille via SNMP :
- Température des salles serveurs (Hot Aisles / Cold Aisles).
- Niveau de carburant des génératrices.
- Charge des batteries des UPS.
- Vitesse de rotation des ventilateurs (HVAC).

---
*Document ID: SNISID-OBS-001 | Approuvé par: Head of SRE*
