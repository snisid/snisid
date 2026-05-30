# Data Observability Stack

## Objectif
Surveiller la santé des données nationales, des pipelines, du stockage, de la qualité et des accès.

## À monitorer

| Domaine | Monitoring |
|---|---:|
| Data quality | Oui |
| Pipeline failures | Oui |
| Data drift | Oui |
| Storage health | Oui |
| Access anomalies | Oui |
| Freshness | Oui |
| Lineage gaps | Oui |
| Metadata completeness | Oui |

## Outils

| Domaine | Outil |
|---|---|
| Metrics | Prometheus |
| Logs | Loki |
| Data observability | OpenMetadata/OpenLineage |
| Analytics | Grafana |
| Alerting | Alertmanager/SIEM |

## SLO data

| Service | SLO |
|---|---:|
| Ingestion critique | 99.9% jobs réussis |
| Fraîcheur identité | < 5 minutes pour événements critiques |
| Catalogue metadata | 100% datasets avec owner/classification |
| Lineage critique | 100% complet |
| DQ Gold | >= seuil par classification |
| Audit ingestion | 99.99% événements capturés |

## Alertes critiques

- Pipeline critique échoué > 2 tentatives.
- DQ score sous seuil.
- Dataset sans owner/classification.
- Accès massif inhabituel.
- Écart de volume > 3 sigma.
- Drift modèle IA supérieur au seuil.
- Retard audit events > 1 minute.
