# Real-Time Analytics Platform

## Objectif
Fournir une intelligence nationale temps réel pour opérations, fraude, risque, identité et état national.

## Capacités

| Fonction | Support |
|---|---:|
| Streaming analytics | Oui |
| Fraud analytics | Oui |
| Operational dashboards | Oui |
| Risk scoring | Oui |
| Identity analytics | Oui |
| Alerting temps réel | Oui |

## Architecture

```text
Kafka Topics -> Stream Processing -> Feature/Metric Stores -> Dashboards/Alerts/APIs
                         \-> Lakehouse Silver/Gold
```

## Tableaux de bord nationaux

| Dashboard | Utilisateurs | Données |
|---|---|---|
| État opérationnel SNISID | Centre national opérations | Disponibilité services, files, incidents |
| Identité nationale | Agence identité | inscriptions, doublons, validations |
| Fraude et risque | Inspection/audit | anomalies, scores, patterns |
| Interopérabilité agences | DGO/IT | flux, latence, échecs APIs |
| Sécurité data | SOC/Security | accès anormaux, exports, violations |

## Principes

- Agrégation par défaut pour vues nationales.
- Drill-down soumis à autorisation ABAC.
- Alertes critiques vers SOC/NOC/DGO.
- Toute consultation de données sensibles est auditée.
- Les métriques officielles proviennent de datasets Gold certifiés.
