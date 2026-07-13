# National AI/ML Platform

## Objectif
Préparer une IA souveraine pour détection fraude, doublons, risque, optimisation workflows et détection menaces, tout en garantissant gouvernance, explicabilité et audit.

## Domaines IA

| Domaine | Usage |
|---|---:|
| Fraud detection | Oui |
| Duplicate detection | Oui |
| Risk analytics | Oui |
| Workflow optimization | Oui |
| Threat detection | Oui |

## Composants

| Composant | Fonction |
|---|---|
| Feature Store | Variables certifiées et historisées |
| MLflow/Registry | Enregistrement modèles et versions |
| Model Governance | Approbation, documentation, risques |
| Explainability | SHAP/LIME, règles explicatives |
| Monitoring | Drift, performance, biais |
| Human Review | Décisions sensibles validées humainement |

## Règles obligatoires

- Aucun modèle sans dataset d'entraînement catalogué.
- Aucun modèle sans fiche modèle approuvée.
- Aucun score critique sans explication.
- Aucun modèle en production sans monitoring drift/biais.
- Aucune décision administrative définitive uniquement automatisée pour cas sensibles.
- Tous les scores et décisions sont envoyés à l'Audit Data Fabric.

## Model Card minimale

| Champ | Description |
|---|---|
| model_id | Identifiant unique |
| owner | Propriétaire métier |
| purpose | Finalité autorisée |
| training_data | Datasets et versions |
| features | Variables utilisées |
| performance | Métriques validation |
| bias_tests | Résultats équité |
| explainability | Méthode d'explication |
| approval | Comité/date |
| monitoring | Drift, seuils, alertes |
