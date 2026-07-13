# PHASE_11_IMPLEMENTATION.md

## Nom de la phase
Phase 11 - Usine à Workflows Nationale (Workflow-Factory)

## Objectif
Déployer le composant d'orchestration central du SNISID, l'**Usine à Workflows**. Ce module gère les règles métiers de l'État (Civil, Identité, Justice, Police), l'assignation de tâches humaines, l'évaluation des SLA et le fonctionnement hors-ligne (Offline).

## Fonctionnalités ajoutées
- Déploiement de l'Orchestrateur (Python) : `orchestrator.py`.
- Intégration du bus d'événements interne (`event_bus.py`) et du gestionnaire de dossiers (`case_manager.py`).
- Implémentation du moteur de SLA et d'escalade automatique (`sla_engine.py`).
- Ajout des flux spécialisés par domaine métier : `civil_registry_workflows.py`, `police_workflows.py`, etc.
- Démonstrateur du moteur de résolution de conflits hors-ligne (LWW).

## Fichiers créés / intégrés
L'ensemble de l'applicatif a été copié à la racine du projet sous `SNISID-Workflow-Factory/` :
- `SNISID-Workflow-Factory/orchestrator.py`
- `SNISID-Workflow-Factory/test_orchestrator.py`
- `SNISID-Workflow-Factory/BPMN/`
- `SNISID-Workflow-Factory/Case-Management/`
- `SNISID-Workflow-Factory/Civil-Registry/`
- `SNISID-Workflow-Factory/Event-Driven/`
- `SNISID-Workflow-Factory/Human-Tasks/`
- `SNISID-Workflow-Factory/Identity/`
- `SNISID-Workflow-Factory/Justice/`
- `SNISID-Workflow-Factory/Offline/`
- `SNISID-Workflow-Factory/Police/`
- `SNISID-Workflow-Factory/SLA/`

## Fichiers modifiés
Aucun. 

## Dépendances ajoutées
Le code s'exécute nativement en Python (>=3.8) via la bibliothèque standard (logging, uuid, datetime, threading).

## Variables d’environnement
- N/A pour la suite de tests. 

## Migrations ou changements de base de données
- N/A. Les données d'état sont maintenues en mémoire pour la logique d'orchestration avant persistance.

## Commandes de test / build / déploiement
L'intégrité de l'orchestrateur a été validée avec succès via la commande :
```bash
python test_orchestrator.py
```
*(Validation des 3 scénarios : Workflow normal, Conflit Offline LWW, et Violation SLA)*

## Procédure de rollback
Pour retirer le module de Workflow :
```bash
Remove-Item -Path "c:\Users\sopil\Desktop\snisid system\SNISID-Workflow-Factory" -Recurse -Force
```

## Risques connus
- L'orchestrateur centralise une logique lourde. Son exécution de production nécessitera une migration vers une architecture asynchrone hautement disponible (Kafka/RabbitMQ) afin d'éviter le "Single Point of Failure" inhérent à un exécutable Python unique.

## Points à valider manuellement
- Intégrer les schémas d'événements créés en Phase 9 (`standard_event_schema.json`) à la logique du fichier `event_bus.py`.
