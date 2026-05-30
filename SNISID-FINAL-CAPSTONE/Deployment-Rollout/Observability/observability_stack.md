# SNISID National Deployment Observability Stack
## Architecture de Supervision Technique et Opérationnelle en Temps Réel

---

## 1. Introduction & Architecture d'Observability

Pour piloter efficacement le déploiement national sur les 140 communes et 10 départements sans naviguer à l'aveugle, le SNISID déploie une **pile d'observabilité distribuée (Distributed Observability Stack)** de dernière génération. Cette infrastructure de supervision permet de visualiser instantanément l'avancement de la migration, les défaillances de synchronisation hors-ligne, les latences d'API, et l'état de préparation de chaque région.

```
                              OBSERVABILITY FLOW
                              
  [Local Edge Nodes (LEN)]                 [Central Datacenter Cluster]
    - Prometheus Node Exporter               - Prometheus Server (Metrics)
    - OpenTelemetry Agent                    - Loki Cluster (Log aggregation)
    - Local Edge DB Metrics                  - Tempo Cluster (Tracing)
               \                                    /
                \                                  /
                 v                                v
         [OpenTelemetry Collector Gateway (OTEL Collector Config)]
                                 |
                                 v
                     [GRAFANA CENTRAL DASHBOARDS]
                     - Rollout Progress Panel
                     - Biometric Latency Panel
                     - Alertmanager Integration
```

---

## 2. Les Quatre Piliers Technologiques de la Supervision

La pile repose sur des outils ouverts (Open Source) de référence, adaptés au contexte d'Haïti :

### 2.1 Prometheus (Collecte des Métriques)
*   **Rôle :** Extraction à intervalle régulier (scraping) des métriques techniques et métiers exposées par les Edge Nodes et le datacenter central.
*   **Métriques Métiers Clés :**
    *   `snisid_citizen_enrollments_total` : Nombre total d'enrôlements par commune.
    *   `snisid_migration_failures_total` : Taux d'échec d'importation de l'historique ONI.
    *   `snisid_offline_sync_delay_seconds` : Temps écoulé depuis la dernière synchronisation d'un nœud d'Edge.

### 2.2 Grafana (Visualisation et Dashboards)
*   **Rôle :** Interface de restitution unique pour la War Room nationale. Affiche la carte géographique interactive des départements et la courbe de progression de l'enrôlement national.
*   **Dashboards Déployés :** Dashboard National Rollout Master, Dashboard ABIS performance, Dashboard Connectivity & Solar Autonomy.

### 2.3 Grafana Loki (Centralisation des Journaux - Logs)
*   **Rôle :** Agrégation de tous les journaux d'événements techniques émis par les routeurs Starlink, les systèmes Linux des LEN et les applications applicatives de réconciliation.
*   **Bénéfice :** Permet à un ingénieur de la War Room de lire les logs d'une commune enclavée (ex: Pestel ou Abricots) sans avoir à s'y connecter physiquement en SSH (ce qui consommerait de la bande passante satellite critique).

### 2.4 OpenTelemetry (Traces Applicatives et Profiling)
*   **Rôle :** Instrumentation du code Go et Python de la *Migration Factory* et du *NIRE*.
*   **Bénéfice :** Visualisation du cheminement complet d'un dossier d'enrôlement sous forme de trace distribuée, de la capture biométrique en commune jusqu'à la validation définitive en base centrale. Identifie instantanément l'étape goulot d'étranglement (ex : latence réseau, temps de calcul d'indexation d'iris).

---

## 3. Stratégie d'Alerte et de Notification (Alertmanager Rules)

Les alertes de production sont réparties en 3 niveaux de gravité de routage (Critique, Avertissement, Info) et s'adossent à un système de notification multi-canal (PagerDuty pour les ingénieurs d'astreinte, SMS pour l'Incident Commander, et Slack pour la cellule de crise globale). Les règles précises de déclenchement sont définies dans le fichier de configuration de ce répertoire.
