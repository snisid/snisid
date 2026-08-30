# 📊 KPI D'INTELLIGENCE STRATÉGIQUE

> **Objectif** : Rendre la maturité analytique nationale **mesurable**.

---

## 1. KPI PRINCIPAUX

| KPI | Objectif | Cible | Calcul |
|-----|----------|-------|--------|
| **Decision latency** | Faible | < 5 s pour critiques | T(alerte → décision) P95 |
| **Data quality score** | Élevé | ≥ 95 (Gold) | Composite GE |
| **Fraud detection accuracy** | Très élevée | Précision > 90 % top 1% | Confusion matrix |
| **Predictive accuracy** | Stable | MAPE < 8 % demande, < 3 % pop | Backtesting |
| **Dashboard availability** | > 99.9 % | 99.9 % minimum | Synthetic probes |

---

## 2. KPI ÉTENDUS PAR DOMAINE

### 2.1 Plateforme Analytique

| KPI | Cible |
|-----|-------|
| Lakehouse availability | > 99.95 % |
| Ingestion latency Bronze (P95) | < 60 s |
| Silver freshness | < 15 min |
| Gold freshness | < 30 min |
| Trino query P95 | < 3 s |
| Stockage utilisé / capacité | < 75 % |

### 2.2 IA / ML

| KPI | Cible |
|-----|-------|
| Modèles en production audités | 100 % |
| Modèles avec drift monitoring | 100 % |
| Modèles avec explainability | 100 % |
| Time-to-deploy modèle | < 30 j |
| Couverture tests biais | 100 % |
| Modèles avec rollback testé | 100 % |

### 2.3 BI & Dashboards

| KPI | Cible |
|-----|-------|
| Dispo cockpit présidentiel | > 99.99 % |
| Dispo cockpits ministériels | > 99.9 % |
| Adoption (utilisateurs actifs / inscrits) | > 80 % |
| Satisfaction décideurs (NPS) | > 50 |
| Délai mise en prod dashboard | < 1 j |

### 2.4 Fraude & Risque

| KPI | Cible |
|-----|-------|
| Time-to-alert fraude (P95) | < 60 s |
| Précision fraud scoring | > 90 % |
| Recall fraudes confirmées | > 75 % |
| Doublons biométriques détectés | > 99 % |
| Cas NRIC traités SLA 48h | > 95 % |

### 2.5 Crise

| KPI | Cible |
|-----|-------|
| Délai activation war-room | < 10 min |
| Couverture services régaliens monitoring | 100 % |
| Mise à jour carte impact | < 5 min |
| Disponibilité cockpit crise | 100 % |

### 2.6 GEOINT

| KPI | Cible |
|-----|-------|
| Latence tuile MVT P95 | < 200 ms |
| Couverture cartographique | 100 % communes |
| Fraîcheur position agents | < 30 s |
| Disponibilité GeoServer | > 99.9 % |

### 2.7 Data Governance

| KPI | Cible |
|-----|-------|
| Datasets avec owner | 100 % |
| Datasets avec lineage | > 95 % |
| DQ score moyen Gold | > 95 |
| Incidents qualité / mois | < 5 |

### 2.8 Observability

| KPI | Cible |
|-----|-------|
| MTTD (mean time to detect) | < 5 min |
| MTTR (mean time to resolve) | < 30 min HIGH |
| Alertes vraies / total | > 80 % |
| Couverture instrumentation services | > 95 % |

---

## 3. INDICE DE MATURITÉ ANALYTIQUE NATIONALE (IMAN)

```
IMAN = 0.20 · Plateforme
     + 0.20 · BI
     + 0.20 · IA/ML
     + 0.15 · Fraude/Risque
     + 0.10 · Crise
     + 0.10 · Gouvernance
     + 0.05 · GEOINT
```

Échelle 0–100. Cible Phase 18 fin : **IMAN ≥ 85**.

---

## 4. COCKPIT KPI INTELLIGENCE

Dashboard dédié `Maturité Analytique Nationale` :
- IMAN actuel + historique 12 mois
- Heatmap KPI par domaine
- Alertes KPI hors cible
- Tendances et objectifs trimestriels

---

## 5. PROCESSUS DE REVUE

| Cadence | Audience | Contenu |
|---------|----------|---------|
| Hebdo | Équipes plateforme | KPI techniques |
| Mensuel | Direction SNISID | Tous KPI |
| Trimestriel | Comité gouvernance | IMAN + roadmap |
| Annuel | Présidence / Ministères | Bilan stratégique |

---

## 6. PRINCIPE

> Ce qui ne se mesure pas ne se pilote pas.
> L'intelligence nationale est **mesurable, comparable, améliorable**.
