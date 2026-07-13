# BATCH 2: BACKEND CORE — GO MICROSERVICES + EVENT BACKBONE

## 🎯 OBJECTIF
Construire le noyau backend distribué de SNISID, capable de supporter une charge nationale avec une résilience et une sécurité souveraines.

---

## 🧱 ARCHITECTURE & SERVICES

### 1. CORE SERVICES
- **API Gateway**: Point d'entrée unique, authentification, rate limiting, routing.
- **Identity Service**: Gestion du cycle de vie des identités souveraines.
- **Citizen Service**: Registre national des citoyens et documents.
- **Agency Service**: Gestion des accès et permissions inter-organisations.
- **Fraud Service**: Orchestration de la détection de fraude.
- **Risk Engine**: Calcul de score de risque en temps réel.
- **Notification Service**: Routage multicanal (Email, SMS, Push).
- **Audit Service**: Forensic ledger et traçabilité absolue.
- **Investigation Service**: Outils pour les analystes de sécurité.

### 2. DESIGN PATTERNS
- **Go Microservices**: Utilisation de Go 1.22+ pour la performance et la concurrence.
- **Communication**: gRPC pour l'inter-service, REST/OpenAPI pour l'externe.
- **Event-Driven**: Backbone Kafka pour la propagation asynchrone des états.
- **CQRS & DDD**: Séparation des lectures/écritures et modélisation métier stricte.
- **Hexagonal Architecture**: Indépendance vis-à-vis des bases de données et des drivers.

---

## 🛠️ INFRASTRUCTURE & OPS

### 1. DATABASES & CACHE
- **PostgreSQL**: Stockage persistant, consistant et relationnel.
- **Redis**: Cache distribué, rate limiting et gestion de session.
- **Kafka Integration**: Intégration native pour le streaming d'événements.

### 2. OBSERVABILITÉ & RÉSILIENCE
- **OpenTelemetry**: Tracing distribué et métriques métier.
- **Service Discovery**: Enregistrement dynamique des services via Kubernetes/Istio.
- **Self-Healing**: Health checks gRPC/HTTP et politiques de retry.

---

## 🔒 SÉCURITÉ & DÉPLOIEMENT

### 1. SECURITY-BY-DESIGN
- **mTLS**: Chiffrement systématique des flux inter-services.
- **Audit Logging**: Chaque mutation est enregistrée dans l'Audit Service.
- **Validation**: Schémas Protobuf/JSON rigoureux.

### 2. CI/CD & SCALING
- **Docker/K8s**: Conteneurisation et orchestration native.
- **Horizontal Scaling**: Mise à l'échelle automatique basée sur le CPU/RAM et les messages Kafka.
- **Blue/Green & Canary**: Stratégies de déploiement progressives (Prompts 262-263).

---

## 📜 APIs & WORKFLOWS
- **API Spec**: `api/proto/` (Protobuf) et `api/openapi/` (Swagger).
- **Enrollment Workflow**: Citizen Service -> Kafka -> Risk Engine -> Fraud Service -> Identity Service.
- **Audit Workflow**: Middleware -> Audit Service -> Kafka -> Forensic Store.

---

**BATCH 2 IS ARCHITECTURALLY DEFINED.**
**READY FOR IMPLEMENTATION SCAFFOLDING.**
