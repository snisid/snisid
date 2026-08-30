# PHASE_03_IMPLEMENTATION.md

## Nom de la phase
Phase 3 - Moteur de Workflow (SNISID Workflow Engine)

## Objectif
Implémentation du moteur de workflow national orchestrant les processus métiers critiques (comme l'enregistrement civil) en utilisant Camunda Zeebe et Temporal, avec intégration à Kafka et configuration de l'observabilité.

## Fonctionnalités ajoutées
- Déploiement de l'architecture du Moteur de Workflow (Node.js/TypeScript).
- Intégration des définitions BPMN (`01_birth_registration.bpmn`).
- Configuration de l'observabilité avancée (Tempo, Loki, Prometheus, Grafana).

## Fichiers créés / intégrés
L'ensemble de l'architecture logicielle `SNISID-Phase3/` a été intégré au projet :
- `workflow-engine/src/*` (TypeScript)
- `bpmn/*`
- `observability/*`
- `runbooks/*`

## Fichiers modifiés
Aucun fichier externe aux modules de la Phase 3 n'a été modifié.

## Dépendances ajoutées
- Node.js (`@temporalio/client`, `@camunda8/sdk`, `kafkajs`, `@opentelemetry/sdk-node`, etc.)

## Commandes exécutées
```bash
npm install
npm run build
```
*(Remarque : L'installation a rencontré une erreur d'espace disque `ENOSPC`, nécessitant un nettoyage pour les phases suivantes).*

## Risques connus
- L'espace disque du poste de déploiement est limité, affectant le téléchargement des caches `npm`.
