# SNISID National Rollout — Key Performance Indicators (KPIs)
## Cadre de Pilotage Métrique et Mesure de Succès du Déploiement

---

## 1. Introduction

Pour assurer un déploiement national rigoureux et basé sur des preuves concrètes, la phase 15 s'appuie sur 5 indicateurs de performance clés (KPI) majeurs. Ces KPI sont monitorés en continu par la War Room via Prometheus et Grafana et font l'objet d'un rapport quotidien au Comité Exécutif.

---

## 2. Les 5 Indicateurs de Succès Majeurs (National Rollout KPIs)

```
+---------------------------------------------------------------------------------+
|                                 SNISID KPIs                                     |
+---------------------------------------------------------------------------------+
|  KPI 1: Regional Readiness Score (Seuil requis >= 85%)                         |
|  KPI 2: Migration Success Rate (Seuil requis >= 99.98%)                         |
|  KPI 3: Identity Reconciliation Accuracy (Seuil requis >= 99.99%)               |
|  KPI 4: Offline Continuity Rate (Seuil requis >= 99.95% uptime)                  |
|  KPI 5: Citizen Onboarding Velocity (Cible: 200,000 enrôlements par mois)       |
+---------------------------------------------------------------------------------+
```

### KPI 1 : Regional Readiness Score (Préparation Régionale)
*   **Objectif :** Élever la maturité technique, énergétique, humaine et de connectivité de chaque département avant d'autoriser le Go-Live local.
*   **Formule de Calcul :**
    $$\text{Score Readiness} = (w_1 \times \text{Sec}) + (w_2 \times \text{Eng}) + (w_3 \times \text{Net}) + (w_4 \times \text{HR}) + (w_5 \times \text{Gov})$$
    *(Voir spécification détaillée dans `regional_activation.md`)*
*   **Seuil Target :** $\ge 85/100$ obligatoire pour déclencher l'autorisation d'ouverture physique.

### KPI 2 : Migration Success Rate (Taux de Réussite de Migration)
*   **Objectif :** Garantir que l'intégralité des dossiers historiques de l'ancien système de l'ONI et des registres civils d'Haïti sont migrés sans corruption ni déperdition d'information historique.
*   **Formule de Calcul :**
    $$\text{Migration Success Rate} = \frac{\text{Nombre de Dossiers Migrés avec Succès}}{\text{Nombre de Dossiers Ingestés Total}} \times 100$$
*   **Seuil Target :** $\ge 99.98\%$ de succès de traitement global. Les dossiers en anomalie doivent impérativement être identifiés et placés en quarantaine (Q1/Q2) sous moins de 24 heures sans ralentir le reste du lot.

### KPI 3 : Identity Reconciliation Accuracy (Exactitude de la Réconciliation)
*   **Objectif :** S'assurer de l'unicité stricte de l'identité en Haïti. Prévenir les erreurs de correspondance démographique (Faux Positifs de fusion d'identité) et piéger les tentatives de fraude ou d'usurpation d'identité biométrique.
*   **Formule de Calcul :**
    $$\text{Accuracy Rate} = 100 - (\text{Taux d'Erreurs de Type I (Faux Positifs)} + \text{Taux d'Erreurs de Type II (Faux Négatifs)})$$
*   **Seuil Target :** $\ge 99.99\%$ d'exactitude vérifiée par échantillonnage statistique régulier et contre-audit humain de la file de quarantaine du NIRE.

### KPI 4 : Offline Continuity Rate (Autonomie et Résilience Réseau)
*   **Objectif :** S'assurer que les régions isolées sans accès Internet ou subissant des pannes électriques prolongées continuent d'enrôler et de valider les identités locales sans interruption de service.
*   **Formule de Calcul :**
    $$\text{Offline Continuity Rate} = \frac{\text{Temps de Disponibilité Opérationnelle du LEN (en ligne + hors-ligne)}}{\text{Temps de Service Total de Bureau (BLC)}} \times 100$$
*   **Seuil Target :** $\ge 99.95\%$ d'uptime opérationnel global (garanti par l'architecture Offline-First et l'alimentation solaire autonome).

### KPI 5 : Citizen Onboarding Velocity (Vitesse d'Enrôlement Citoyen)
*   **Objectif :** Assurer la transition de l'ensemble de la population haïtienne éligible (environ 7 millions de personnes de plus de 18 ans) vers la nouvelle Carte Nationale d'Identification Biométrique Unique (CNIBU) dans un délai de 24 mois.
*   **Formule de Calcul :**
    $$\text{Onboarding Velocity} = \text{Nombre de nouveaux citoyens enrôlés par mois}$$
*   **Cible :** Moyenne minimale de $200,000$ citoyens enrôlés et validés par mois à l'échelle nationale après la bascule.
