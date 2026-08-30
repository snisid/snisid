# SNISID National BPMN & Workflow Architecture
## Spécification d'Architecture d'Orchestration Nationale — SNISID v4.0

---

## 1. VISION & DIRECTIVES STRATÉGIQUES
Le Système National d'Identité et de Statut Individuel et Civil Digital (**SNISID**) repose sur un principe absolu : **aucun processus gouvernemental critique ne doit rester manuel.** 

La Plateforme Nationale d'Orchestration et la "Workflow Factory" forment le cœur décisionnel et exécutif de l'État Digital. Ce système transforme les administrations cloisonnées en un réseau synchronisé d'acteurs publics, géré par un moteur d'exécution de processus métier (BPMN 2.0) asynchrone, distribué, ultra-sécurisé, et piloté par les événements (Event-Driven Architecture - EDA).

---

## 2. ARCHITECTURE TECHNIQUE GLOBALE

L'architecture s'articule autour de 6 piliers majeurs intégrés via un bus de messages Kafka hautement disponible :

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                              COUCHES PORTAILS & CANAUX                                 │
│         Portail Citoyen (Web/Mobile)  •  Portails Agents  •  Terminaux Mobiles (Offline)  │
└──────────────────────────────────────────┬─────────────────────────────────────────────┘
                                           │ API Gateway
┌──────────────────────────────────────────▼─────────────────────────────────────────────┐
│                       NATIONAL BPMN ENGINE & WORKFLOW FACTORY                          │
│ ┌───────────────────────────┐ ┌───────────────────────────┐ ┌────────────────────────┐ │
│ │  Camunda BPMN Core Engine  │ │  Workflow Registry (v1)   │ │  SLA & Timer Engine    │ │
│ └─────────────┬─────────────┘ └─────────────┬─────────────┘ └───────────┬────────────┘ │
│               │                             │                           │              │
│ ┌─────────────▼─────────────┐ ┌─────────────▼─────────────┐ ┌───────────▼────────────┐ │
│ │  Human Task Manager (HTM)  │ │  National Case Manager    │ │  Offline Sync Service  │ │
│ └─────────────┬─────────────┘ └─────────────┬─────────────┘ └───────────┬────────────┘ │
└───────────────┼─────────────────────────────┼───────────────────────────┼──────────────┘
                │                             │                           │
┌───────────────▼─────────────────────────────▼───────────────────────────▼──────────────┐
│                    KAFKA EVENT BUS (NATIONAL SECURITY ORCHESTRATOR)                    │
│   Topics: birth.events  •  identity.events  •  justice.events  •  workflow.audit       │
└──────────────────────────────────────────┬─────────────────────────────────────────────┘
                                           │
┌──────────────────────────────────────────▼─────────────────────────────────────────────┐
│                          CADRE D'OBSERVABILITÉ & D'AUDIT                              │
│   Prometheus (Metrics)  •  Loki (Logs)  •  OpenTelemetry (Tracing)  •  Grafana Portal   │
└────────────────────────────────────────────────────────────────────────────────────────┘
```

### 2.1 Description des Composants Systèmes

| Domaine | Composant Technique | Fonctionnalité Clé |
| :--- | :--- | :--- |
| **BPMN Engine** | Camunda BPMN Core / Orchestrateur Custom | Exécution transactionnelle de graphes BPMN 2.0, gestion des états persistants, tokens d'exécution. |
| **Workflow Registry** | National Registry DB (PostgreSQL) | Catalogue consolidé des définitions, versions (SemVer), schémas d'événements, niveaux de classification (Secret-Défense à Public). |
| **Case Management** | National Case Management (NCM) | Gestion de dossiers complexes de cycle de vie (Dossier d'État Civil, Dossier Judiciaire, Fiche d'Immigration). |
| **Human Tasks** | Human Task Manager (HTM) | Allocation intelligente de tâches, file d'attente sécurisée, escalades automatiques, double-contrôle (règle des 4 yeux). |
| **SLA Engine** | National SLA Engine | Surveillance microseconde des délais légaux. Déclenchement automatique de runbooks d'alerte et de réaffectation d'urgence. |
| **Event Bus** | Apache Kafka Cluster | Messagerie asynchrone, persistance des journaux d'événements, découplage des microservices et livraison de messages garantie (At-Least-Once). |

---

## 3. PRINCIPES D'ORCHESTRATION PILOTÉE PAR LES ÉVÉNEMENTS (EVENT-DRIVEN)

L'orchestration au sein de SNISID est entièrement réactive. L'orchestrateur central ne réalise pas d'appels RPC bloquants. Il écoute des événements d'état et publie des commandes sur Kafka, permettant un découplage absolu et une résilience maximale en cas de panne d'un système tiers.

### 3.1 Cycle de Vie de l'Événement Standardisé
1. **Production** : Un service de terrain (ex: Maternité) publie l'événement `birth.created` contenant le payload standardisé et signé.
2. **Ingestion & Validation** : Le Event Bus valide le schéma via le Registry.
3. **Trigger BPMN** : Le moteur BPMN consomme l'événement, résout la corrélation (par ex. `National-Correlation-ID`) et démarre le workflow "Naissance Simple" ou met à jour une instance existante.
4. **Commandes de Tâche** : Le moteur génère une commande `task.assign` pour validation humaine par l'Officier de l'État Civil.
5. **Résolution** : Une fois validé, l'événement `workflow.approved` est diffusé, ce qui déclenche automatiquement la création de l'identité nationale dans le registre d'identification.

---

## 4. CAS DE GESTION ADMINISTRATIVE ET ARCHITECTURE DE DOSSIERS (CASE MANAGEMENT)

Le **National Case Management System (NCMS)** traite chaque dossier comme une entité dynamique (le "Dossier National") :

- **Dossier d'État Civil** : Regroupe les événements de naissance, mariage, divorce, décès, adoptions.
- **Dossier d'Identité** : Historique d'enrôlement biométrique, modifications de données civiles, renouvellements, révocations, et alertes de fraude.
- **Dossier Judiciaire & Police** : Mandats de perquisition, cas pénitentiaires, fiches de recherche, dossiers d'enquêtes DCPJ. Un chiffrement asymétrique garantit que seul le personnel de justice dument habilité peut lire le contenu sensible.
- **Dossier de Contrôle d'Immigration** : Demandes de visas, permis de séjour, rapports de passage frontières.

Chaque dossier dispose d'un journal immuable d'audits stocké en base et répliqué sur un registre cryptographique pour garantir une traçabilité inviolable.

---

## 5. ROBUSTESSE, HAUTE DISPONIBILITÉ ET SÉCURITÉ

- **Zéro Point Unique de Défaillance (SPOF)** : Moteurs BPMN conteneurisés en grappe active-active sur Kubernetes (K8s), répartis sur des zones de disponibilité physiques distinctes.
- **Haute performance** : Persistance optimisée avec des bases de données relationnelles partitionnées pour le runtime BPMN, et des bases de documents NoSQL pour le Case Management.
- **Sécurité Globale** :
  - Signature cryptographique (Ed25519) de chaque transition de statut.
  - Chiffrement au repos (AES-256-GCM) des données sensibles de cas.
  - Contrôle d'accès basé sur les rôles (RBAC) renforcé par des règles de sécurité basées sur les attributs (ABAC - Attribute-Based Access Control).
