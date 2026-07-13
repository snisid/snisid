# PHASE_15_IMPLEMENTATION.md

## Nom de la phase
Phase 15 - National Deployment Rollout (Opérationnalisation)

## Objectif
Préparer et automatiser le déploiement national du système SNISID en structurant les stratégies de "Cutover" (bascule), la gestion des incidents ("War-Room" et "Hypercare"), et surtout en opérationnalisant la stack d'observabilité (Prometheus/Grafana) et le moteur de réconciliation d'identités (NIRE) via des API REST.

## Fonctionnalités ajoutées
1. **API NIRE (National Identity Reconciliation Engine)** : Création d'un service web (FastAPI) exposant le moteur d'IA de réconciliation (Levenshtein + Cosine Similarity biométrique). L'API Gateway peut désormais interroger `/api/v1/reconcile` pour valider en temps réel les identités migrées ou créées.
2. **Observability Stack (Infrastructure-as-Code)** : Création du fichier d'orchestration Docker Compose pour lancer de manière coordonnée Prometheus, Grafana et l'OpenTelemetry Collector.

## Fichiers créés
- `SNISID-FINAL-CAPSTONE/Deployment-Rollout/Identity-Reconciliation/api_nire.py`
- `SNISID-FINAL-CAPSTONE/Deployment-Rollout/Identity-Reconciliation/requirements.txt`
- `SNISID-FINAL-CAPSTONE/Deployment-Rollout/Observability/docker-compose.observability.yml`

## Fichiers modifiés
- `PHASE_15_IMPLEMENTATION.md` (Mise à jour pour inclure l'opérationnalisation technique).

## Dépendances ajoutées
- Python (FastAPI, Uvicorn, Pydantic) pour l'API NIRE.
- Docker et Docker Compose pour la stack d'Observabilité.

## Variables d’environnement
- `GF_SECURITY_ADMIN_PASSWORD` (Provisionné dans le docker-compose pour Grafana).

## Changements de base de données
- Aucun changement structurel de base de données (le NIRE fonctionne en in-memory pour l'évaluation).

## Commandes exécutées
Création de l'arborescence et génération de code source :
```bash
# Pour lancer l'API NIRE (Nécessite Python et libération de l'espace disque) :
cd SNISID-FINAL-CAPSTONE/Deployment-Rollout/Identity-Reconciliation
pip install -r requirements.txt
python api_nire.py

# Pour lancer la Stack d'Observabilité (Nécessite Docker et libération de l'espace disque) :
cd SNISID-FINAL-CAPSTONE/Deployment-Rollout/Observability
docker-compose -f docker-compose.observability.yml up -d
```

## Instructions de déploiement
L'API Gateway (Phase 10) devra router le trafic d'enregistrement d'identité vers l'API NIRE sur le port `8000`. L'Observability stack sera déployée en production sur le port `3000` (Grafana) et `9090` (Prometheus).

## Procédure de rollback
Pour annuler l'implémentation logicielle de la Phase 15 :
```bash
rm SNISID-FINAL-CAPSTONE/Deployment-Rollout/Identity-Reconciliation/api_nire.py
rm SNISID-FINAL-CAPSTONE/Deployment-Rollout/Identity-Reconciliation/requirements.txt
rm SNISID-FINAL-CAPSTONE/Deployment-Rollout/Observability/docker-compose.observability.yml
```

## Risques connus
- **Surcharge de l'API NIRE** : Les calculs Levenshtein peuvent consommer du CPU si le volume de migration (Cutover) est très élevé. Prévoir un scaling horizontal (Uvicorn workers).
- **ENOSPC** : Espace disque insuffisant pour télécharger les images Docker Grafana/Prometheus pour le moment.

## Points à valider manuellement
- S'assurer que les images Docker se téléchargent avec succès lorsque l'espace disque du serveur sera libéré.
- Vérifier la connexion entre Prometheus et l'OpenTelemetry Collector via le port 8889.
