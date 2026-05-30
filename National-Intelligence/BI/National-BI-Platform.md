# 📊 NATIONAL BI PLATFORM

> **Objectif** : Business Intelligence gouvernementale offrant aux dirigeants une visibilité nationale temps réel.

---

## 1. CAPACITÉS

| Fonction | Support | Outil principal |
|----------|:-------:|-----------------|
| Dashboards | ✅ | Apache Superset |
| KPI tracking | ✅ | Grafana |
| Executive reports | ✅ | Metabase + export PDF signé |
| Cross-agency analytics | ✅ | Superset (datasets fédérés via Trino) |

---

## 2. STACK BI

| Composant | Rôle |
|-----------|------|
| **Apache Superset** | Plateforme BI principale, dashboards self-service |
| **Grafana** | KPI temps réel, alerting opérationnel |
| **Metabase** | Reporting exécutif simplifié, ad-hoc |
| **Trino** | Couche SQL fédérée Lakehouse / PostgreSQL / Druid |
| **dbt** | Modélisation analytique versionnée Git |

---

## 3. ORGANISATION DES DASHBOARDS

```
Superset Workspaces
├── 01_Presidence/                # Cockpit présidentiel
├── 02_Ministeres/
│   ├── Interieur/
│   ├── Justice/
│   ├── Sante/
│   ├── Education/
│   └── Finances/
├── 03_Regions/                   # 10 départements haïtiens
├── 04_Operations/                # Opérations SNISID
├── 05_Securite_Fraude/
├── 06_Crise/
└── 07_Analytics_Avance/          # Datascience exploratoire
```

---

## 4. EXEMPLES DE DASHBOARDS CLÉS

### 4.1 Dashboard Présidentiel
- Population identifiée (cumul + croissance)
- Couverture nationale par région (heatmap)
- Disponibilité services régaliens (SLO)
- Alertes critiques (sécurité, crise, fraude)
- Indicateurs économiques liés à l'identité

### 4.2 Dashboard Ministère Intérieur
- Enrôlements en cours par département
- Documents délivrés / en attente
- Performance bureaux locaux
- Backlogs et files d'attente

### 4.3 Dashboard Opérations
- Throughput biométrique
- Latence services API
- Taux de succès workflows
- Charge agents nationaux

### 4.4 Dashboard Cross-Agency
- Échanges inter-ministères
- Conformité partage données
- KPI interopérabilité

---

## 5. SÉCURITÉ BI

| Contrôle | Mise en œuvre |
|----------|---------------|
| Authentification | OIDC SSO avec IAM SNISID |
| Autorisation | RBAC + row-level security Superset |
| Audit | Toutes requêtes loguées (Loki) |
| Watermarking | Exports PDF marqués utilisateur + timestamp |
| Classification | Étiquetage SECRET / CONFIDENTIEL / PUBLIC |

---

## 6. EXEMPLE DE DÉFINITION DASHBOARD (Superset YAML)

```yaml
dashboard_title: "Cockpit Présidentiel SNISID"
slug: presidential-cockpit
owners: [presidence_snisid]
position: |
  CHART-population: {row: 0, col: 0, size_x: 6, size_y: 4}
  CHART-coverage-map: {row: 0, col: 6, size_x: 6, size_y: 8}
  CHART-services-slo: {row: 4, col: 0, size_x: 6, size_y: 4}
  CHART-critical-alerts: {row: 8, col: 0, size_x: 12, size_y: 4}
metadata:
  refresh_frequency: 60          # secondes
  default_filters:
    region: "ALL"
  color_scheme: "haiti_official"
  rls_filter: "user_region = current_user_region()"
```

---

## 7. KPI TEMPS RÉEL (GRAFANA)

```yaml
# grafana-dashboard-national-kpi.json (extrait)
panels:
  - title: "Enrôlements / minute"
    type: timeseries
    datasource: prometheus
    targets:
      - expr: rate(snisid_enrollments_total[1m])

  - title: "Taux fraude détectée"
    type: stat
    targets:
      - expr: |
          sum(rate(snisid_fraud_alerts_total[5m]))
          / sum(rate(snisid_enrollments_total[5m]))
    thresholds:
      - {value: 0.01, color: green}
      - {value: 0.05, color: orange}
      - {value: 0.1,  color: red}

  - title: "Disponibilité services (SLO 99.9%)"
    type: gauge
    targets:
      - expr: avg_over_time(snisid_service_availability[24h])
```

---

## 8. EXPORTS EXÉCUTIFS (METABASE)

- Rapports hebdomadaires PDF signés cryptographiquement
- Distribution sécurisée via portail ministériel
- Archivage WORM 10 ans
- Traçabilité complète des consultations

---

## 9. SLA

| Indicateur | Cible |
|------------|-------|
| Disponibilité plateforme | > 99.9 % |
| Latence dashboard | < 2 s (P95) |
| Fraîcheur données | < 5 min (gold) |
| Temps mise en prod dashboard | < 1 jour |
