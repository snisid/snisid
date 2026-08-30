# ⚡ REAL-TIME EVENT ANALYTICS ENGINE

> **Objectif** : Analyser les événements nationaux en temps réel pour détection immédiate.

---

## 1. CAPACITÉS

| Fonction | Support | Latence cible |
|----------|:-------:|---------------|
| Streaming analytics | ✅ | < 1 s |
| Fraud detection | ✅ | < 500 ms |
| Identity anomalies | ✅ | < 1 s |
| Operational monitoring | ✅ | temps réel |

---

## 2. STACK TECHNIQUE

| Domaine | Technologie | Rôle |
|---------|-------------|------|
| **Streaming bus** | Apache Kafka | Backbone événementiel national |
| **Stream processing** | Apache Flink | CEP, windowing, joins streaming |
| **Event analytics OLAP** | Apache Druid | Requêtes sub-seconde sur événements |
| **Schema Registry** | Confluent / Apicurio | Avro / Protobuf gouvernance schémas |
| **Visualisation** | Grafana + Superset | Dashboards temps réel |

---

## 3. TOPOLOGIE

```
       Applications SNISID
              │
              ▼
   ┌──────────────────────┐
   │       KAFKA          │  topics: enrollments, biometrics,
   │   (cluster 3+ AZ)    │          fraud_signals, security_events,
   └──────────┬───────────┘          workflow_events, gis_events
              │
   ┌──────────┴───────────┐
   ▼                      ▼
┌─────────┐         ┌────────────┐
│ FLINK   │         │ DRUID      │
│ jobs    │────────▶│ ingestion  │
│ (CEP)   │         │ (Kafka in) │
└────┬────┘         └─────┬──────┘
     │                    │
     ▼                    ▼
┌──────────┐       ┌────────────┐
│ Lakehouse│       │ Superset / │
│ (Delta)  │       │  Grafana   │
└──────────┘       └────────────┘
```

---

## 4. CAS D'USAGE — FRAUD DETECTION

### 4.1 Pattern : enrôlements anormalement rapides par un agent

```sql
-- Flink SQL
INSERT INTO fraud_alerts
SELECT
  agent_id,
  region,
  TUMBLE_START(event_ts, INTERVAL '10' MINUTE) AS win_start,
  COUNT(*) AS enrollments,
  'AGENT_RATE_ANOMALY' AS rule_id
FROM enrollments_kafka
GROUP BY agent_id, region, TUMBLE(event_ts, INTERVAL '10' MINUTE)
HAVING COUNT(*) > 50;
```

### 4.2 Pattern CEP : tentative multi-identité

```java
Pattern<EnrollmentEvent, ?> multiIdentity = Pattern
    .<EnrollmentEvent>begin("first")
    .where(e -> e.getBiometricMatchScore() > 0.95)
    .followedBy("second")
    .where(e -> e.getBiometricMatchScore() > 0.95)
    .within(Time.minutes(60));

CEP.pattern(enrollmentStream.keyBy(e -> e.getBiometricFingerprintHash()),
            multiIdentity)
   .select(new MultiIdentityAlertSelector())
   .addSink(new KafkaSink<>("fraud_alerts"));
```

---

## 5. CAS D'USAGE — IDENTITY ANOMALIES

| Règle | Description |
|-------|-------------|
| Geo-impossible | Même NIN utilisé dans 2 régions distantes < 1h |
| Replay biometric | Empreinte identique > N fois en 24h |
| Off-hours | Activité agent hors plage horaire |
| Burst | Pic d'enregistrements 5x supérieur à la moyenne |
| Insider | Agent consulte > X dossiers sans raison |

---

## 6. CAS D'USAGE — OPERATIONAL MONITORING

Druid permet requêtes sub-seconde :

```sql
-- Druid SQL : performance services dernière heure
SELECT
  service_name,
  COUNT(*) AS calls,
  APPROX_QUANTILE(latency_ms, 0.95) AS p95_latency,
  SUM(CASE WHEN status >= 500 THEN 1 ELSE 0 END) AS errors
FROM "snisid_service_events"
WHERE __time > CURRENT_TIMESTAMP - INTERVAL '1' HOUR
GROUP BY service_name
ORDER BY p95_latency DESC;
```

---

## 7. ALERTING TEMPS RÉEL

Flink → Kafka topic `alerts` → consommateurs :
- **Alertmanager** (équipes ops)
- **Webhook Risk Center** (analystes renseignement)
- **Mobile push** (décideurs niveau ministère)
- **War-Room dashboard** (en cas de crise)

Niveaux : `INFO`, `WARNING`, `CRITICAL`, `NATIONAL_EMERGENCY`.

---

## 8. PERFORMANCE CIBLE

| Métrique | Cible |
|----------|-------|
| Throughput Kafka | 500 000 msg/s |
| Flink end-to-end latency | < 1 s P95 |
| Druid query latency | < 500 ms P95 |
| Rétention hot (Druid) | 30 jours |
| Rétention warm (Lakehouse) | 10 ans |

---

## 9. GOUVERNANCE ÉVÉNEMENTIELLE

- Schémas Avro versionnés et obligatoires
- Validation à l'ingestion (rejet si non conforme)
- Catalogue d'événements documenté
- Politique de PII : pseudonymisation systématique des NIN
