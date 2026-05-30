# SNISID : L'USINE NATIONALE DES WORKFLOWS (NATIONAL WORKFLOW FACTORY)
## Index Principal et Architecture d'Orchestration
**République d'Haïti — Production Gouvernementale Critique H24/7/365**

---

## PRÉAMBULE
Ce document constitue le point d'entrée central (Master Index) de **l'Usine Nationale des Workflows (National Workflow Factory)** du SNISID. Ce système souverain orchestre l'intégralité des cycles de vie de l'état civil, de l'identité numérique et des processus de cyber-résilience de la République d'Haïti.

Le moteur de workflow est conçu pour une exécution 24/7 en production massive, s'appuyant sur des modèles de transactions distribuées (Saga) orchestrées par Temporal/Camunda, et des échanges inter-agences transitant par Apache Kafka et le bus d'interopérabilité souverain X-Road.

---

## 📑 SOMMAIRE DE LA DOCUMENTATION ARCHITECTURALE

L'architecture détaillée de l'Usine des Workflows est découpée en 5 volumes techniques distincts :

### 🏛️ [VOLUME 1 : Workflows d'État Civil (Civil Registration)](workflows/SNISID_WF_01_Civil_Registration.md)
Modélisation BPMN complète des processus d'état civil de l'ANH et du MJSP :
- **Naissances** (Simple, Reconnaissance, Tardive, Décret, Jugement)
- **Décès** (Standard, Judiciaire, Catastrophe, Disparition)
- **Mariages** (Civil, Judiciaire, Religieux)
- **Divorces** (Administratif, Judiciaire)
- **Adoptions** (Nationale, Internationale)

### 👤 [VOLUME 2 : Workflows d'Identité et de Sécurité (Identity & Security)](workflows/SNISID_WF_02_Identity_Security.md)
Modélisation des processus de gestion du cycle de vie de l'Identité Nationale (ONI/DCPJ) :
- Enrôlement, capture biométrique et Deduplication ABIS
- Correction d'identité, restauration après usurpation, révocation
- Adjudication, contestations et escalade judiciaire
- Processus anti-fraude et investigations DCPJ

### ⚙️ [VOLUME 3 : Ingénierie des Workflows (Engineering Standards)](workflows/SNISID_WF_03_Engineering_Standards.md)
Spécifications techniques obligatoires pour l'implémentation de tout BPMN :
- Accords de niveau de service (SLA) et objectifs (SLO)
- Stratégies de réessai (Retries), Dead Letter Queues (DLQ)
- Transactions distribuées : Modèle Saga et compensations automatiques (Rollbacks)
- Pistes d'audit immuables (WORM), validation humaine et signatures PKI
- Modèles mathématiques de *Fraud Scoring*

### 📡 [VOLUME 4 : Gouvernement Piloté par les Événements (Event-Driven Architecture)](workflows/SNISID_WF_04_Event_Driven_Government.md)
Architecture du bus asynchrone national :
- Registre des Topics Kafka souverains
- Schémas d'événements stricts (Protobuf / Avro)
- Tolérance aux pannes, idempotence absolue, et stratégies de rejeu
- Distribution inter-agences via X-Road

### 🛡️ [VOLUME 5 : Gouvernance et Cycle de Vie (Governance & Lifecycle)](workflows/SNISID_WF_05_Governance.md)
Cadre légal et opérationnel régissant la modification des processus d'État :
- Propriété fonctionnelle (Bureau de Gouvernance)
- Processus de validation légale stricte (Signature MJSP)
- Déploiement Blue-Green et versioning (SemVer)
- Catalogue des microservices et Matrice d'escalade d'urgence

---

## 🔐 PRINCIPE FONDAMENTAL DE CYBER-RÉSILIENCE
L'Usine Nationale des Workflows a été conçue pour surmonter les défis structurels d'Haïti (instabilité réseau, désastres naturels). Grâce à ses **événements asynchrones**, ses **sagas tolérantes aux coupures prolongées**, et son adhésion stricte au **Zero Trust cryptographique (AN-PKI)**, elle constitue le cœur technologique inaltérable de la République.

*Spécifié, modélisé et ratifié par l'Usine des Architectes de l'État Haïtien.*
*Classification : SOUVERAIN / INFRASTRUCTURE CRITIQUE NATIONALE DE SÉCURITÉ DE L'ÉTAT*
