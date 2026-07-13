# SNISID — Modèles d'Architecture C4
# SNISID — C4 Architecture Models

---

| Métadonnée | Valeur |
|---|---|
| **Document ID** | SNISID-ARC-C4-001 |
| **Version** | 1.0.0 |
| **Date** | 2026-05-25 |
| **Statut** | APPROUVÉ — Production |
| **Classification** | RESTREINT / RESTRICTED |
| **Auteur** | Architecture Team — SNISID Programme |
| **Révisé par** | Chief Architect, Security Architect, Infrastructure Lead |
| **Approuvé par** | Directeur Général, SNISID Programme Office |
| **Standard** | C4 Model (Simon Brown), Arc42, ISO/IEC 42010 |

---

## Table des Matières / Table of Contents

1. [Introduction et Contexte](#1-introduction-et-contexte)
2. [C4 Niveau 1 — Contexte Système](#2-c4-niveau-1--contexte-système)
3. [C4 Niveau 2 — Conteneurs](#3-c4-niveau-2--conteneurs)
4. [C4 Niveau 3 — Composants](#4-c4-niveau-3--composants)
5. [C4 Niveau 4 — Code et Chemins Critiques](#5-c4-niveau-4--code-et-chemins-critiques)
6. [Décisions Technologiques](#6-décisions-technologiques)
7. [Exigences Non-Fonctionnelles par Composant](#7-exigences-non-fonctionnelles-par-composant)
8. [Contraintes Architecturales](#8-contraintes-architecturales)
9. [Perspectives d'Évolution](#9-perspectives-dévolution)

---

## 1. Introduction et Contexte

Le **Système National d'Identification Souverain et Intégré de la République d'Haïti (SNISID)** constitue l'infrastructure numérique de souveraineté la plus critique du pays. Ce document présente l'architecture complète selon le modèle C4 (Context, Containers, Components, Code), offrant quatre niveaux de zoom progressifs pour chaque audience technique.

### 1.1 Portée du Document

Ce document couvre :
- L'intégration du SNISID dans l'écosystème numérique haïtien (C4 L1)
- Tous les conteneurs applicatifs et infrastructurels majeurs (C4 L2)
- Les composants internes des services critiques (C4 L3)
- Les chemins de code critiques pour les opérations sensibles (C4 L4)

### 1.2 Principes Architecturaux Fondamentaux

| Principe | Description | Implication |
|---|---|---|
| **Souveraineté numérique** | Toutes données biométriques restent en territoire haïtien | Pas de cloud étranger pour les données PII/biométriques |
| **Résilience by design** | Tolérance aux pannes réseau et électriques | Mode offline obligatoire pour tous les services de terrain |
| **Zero Trust** | Jamais confiance, toujours vérifier | mTLS, SPIFFE/SPIRE sur tous les services |
| **Privacy by Design** | RGPD-compatible, minimisation des données | Chiffrement de bout en bout, pseudonymisation |
| **Interopérabilité** | OpenAPI 3.1, FHIR R4, GovStack | Standards ouverts uniquement |
| **Haute disponibilité** | RTO < 15 min, RPO < 5 min | Active-Active multi-DC |
| **Auditabilité totale** | Traçabilité de chaque accès | Audit immutable sur blockchain privée |

---

## 2. C4 Niveau 1 — Contexte Système

### 2.1 Acteurs Primaires et Systèmes Externes

```mermaid
C4Context
    title Contexte Système SNISID — Écosystème Numérique Haïtien

    Person(citizen, "Citoyen Haïtien", "Personne physique résidant en Haïti ou diaspora, demandant des services d'identification nationale")
    Person(agent_oec, "Agent OEC", "Agent de l'Office de l'État Civil enregistrant les naissances, mariages, décès")
    Person(agent_mel, "Agent MEL", "Agent du Ministère Electoral et des Libertés enregistrant les électeurs")
    Person(admin_snisid, "Administrateur SNISID", "Agent technique gérant le système d'identification souverain")
    Person(agent_police, "Agent PNH/DCPJ", "Agent de la Police Nationale d'Haïti accédant aux données d'identité")

    System(snisid, "SNISID", "Système National d'Identification Souverain et Intégré — Gère l'identité numérique souveraine de tous les citoyens haïtiens (biométrie, documents, attributs civils)")

    System_Ext(oec_sys, "Système OEC", "Office de l'État Civil — Registre des naissances, mariages, décès")
    System_Ext(mel_sys, "Système MEL", "Ministère Electoral — Registre électoral, listes électorales")
    System_Ext(mfa_sys, "Système MAECI", "Ministère des Affaires Étrangères — Passeports, visas, légalisations")
    System_Ext(msp_sys, "Système MSP", "Ministère de la Santé Publique — Dossiers médicaux, certificats de santé")
    System_Ext(dgif_sys, "Système DGIF", "Direction Générale des Impôts et Finances — NIF, déclarations fiscales")
    System_Ext(minjustice, "Système Justice", "Ministère de la Justice — Casier judiciaire, actes notariés")
    System_Ext(interpol, "INTERPOL", "Base internationale — Vérification des documents volés/perdus")
    System_Ext(icao_pki, "ICAO PKI Directory", "Infrastructure PKI ICAO — Validation eMRTD, passeports biométriques")
    System_Ext(un_unhcr, "UNHCR / OIM", "Agences ONU — Coordination réfugiés, apatrides")
    System_Ext(diaspora_portal, "Portail Diaspora", "Services en ligne pour les Haïtiens de l'étranger")
    System_Ext(banking_sys, "Système Bancaire BRH", "Banque de la République d'Haïti — KYC, identification financière")
    System_Ext(caricom_id, "CARICOM Digital ID", "Système d'identification régionale CARICOM")
    System_Ext(sms_gateway, "Passerelle SMS/OTP", "Services de notification OTP — Digicel, Natcom")

    Rel(citizen, snisid, "Enrôlement biométrique, demande de documents, vérification d'identité", "HTTPS/TLS 1.3, API REST")
    Rel(agent_oec, snisid, "Enregistrement actes d'état civil, consultation identité", "HTTPS mTLS, VPN IPsec")
    Rel(agent_mel, snisid, "Inscription électorale, mise à jour dossier citoyen", "HTTPS mTLS, VPN IPsec")
    Rel(admin_snisid, snisid, "Administration, supervision, audit", "HTTPS mTLS, Zero Trust VPN")
    Rel(agent_police, snisid, "Vérification identité, recherche biométrique", "HTTPS mTLS, RBAC/ABAC strict")

    Rel(snisid, oec_sys, "Synchronisation état civil, validation données", "REST/gRPC mTLS, Kafka streaming")
    Rel(snisid, mel_sys, "Partage registre électoral, mise à jour inscriptions", "REST mTLS, fichiers sécurisés SFTP")
    Rel(snisid, mfa_sys, "Validation identité pour passeports, demandes consulaires", "REST mTLS, SAML 2.0")
    Rel(snisid, msp_sys, "Partage identifiant patient, coordination vaccination", "FHIR R4, REST mTLS")
    Rel(snisid, dgif_sys, "Partage NIF, identification fiscale", "REST mTLS, XML chiffré")
    Rel(snisid, minjustice, "Consultation casier judiciaire, actes notariés", "REST mTLS, audit trail")
    Rel(snisid, interpol, "Vérification documents volés, avis de recherche", "HTTPS REST, XML/JSON")
    Rel(snisid, icao_pki, "Validation passeports eMRTD, vérification BAC/EAC", "LDAPS, HTTPS")
    Rel(snisid, un_unhcr, "Coordination identité réfugiés, apatrides", "REST API, échange sécurisé")
    Rel(snisid, banking_sys, "Fourniture API KYC, identité pour services financiers", "OAuth 2.1, REST mTLS")
    Rel(snisid, caricom_id, "Interopérabilité régionale, reconnaissance d'identité", "REST API, OpenID Connect")
    Rel(snisid, sms_gateway, "Envoi OTP, notifications citoyennes", "HTTPS REST, chiffré")
    Rel(diaspora_portal, snisid, "Demandes à distance, renouvellement documents", "HTTPS, OIDC, mTLS")
```

### 2.2 Périmètre de Responsabilité SNISID

```mermaid
graph TB
    subgraph "Dans le périmètre SNISID"
        A[Registre d'Identité Nationale]
        B[Base Biométrique Souveraine]
        C[Infrastructure PKI Nationale]
        D[Moteur de Déduplication]
        E[Portail Citoyen / eID]
        F[API Gateway National]
        G[Audit Trail Immutable]
        H[Services d'Enrôlement Terrain]
    end

    subgraph "Hors périmètre — Systèmes Fédérés"
        I[OEC — Actes d'État Civil]
        J[MEL — Registre Électoral]
        K[MSP — Dossiers Médicaux]
        L[BRH — Système Bancaire]
        M[MAECI — Passeports]
    end

    subgraph "Infrastructure Partagée"
        N[Datacenters Nationaux]
        O[Backbone Réseau Haïti]
        P[Passerelles SMS/Email]
    end

    A --> B
    A --> C
    B --> D
    F --> A
    F --> B
    G --> A
    H --> A

    A <-->|"APIs Gouvernées"| I
    A <-->|"APIs Gouvernées"| J
    A <-->|"APIs Gouvernées"| K
    A <-->|"APIs Gouvernées"| L
    A <-->|"APIs Gouvernées"| M
```

---

## 3. C4 Niveau 2 — Conteneurs

### 3.1 Diagramme des Conteneurs SNISID

```mermaid
C4Container
    title Conteneurs SNISID — Vue d'Architecture Complète

    Person(citizen, "Citoyen", "Utilisateur final")
    Person(agent, "Agent Terrain", "Agent d'enrôlement, OEC, MEL")
    Person(admin, "Administrateur", "Ops, sécurité, audit")

    System_Boundary(snisid_boundary, "SNISID — Système National d'Identification Souverain") {

        Container(api_gateway, "API Gateway", "Kong Enterprise + Envoy Proxy", "Point d'entrée unique: authentification, rate limiting, routage, transformation. Applique OAuth 2.1 et mTLS sur toutes les routes")
        Container(identity_svc, "Identity Service", "Go 1.22 + gRPC + REST", "Cœur du système: gestion du cycle de vie des identités nationales (NIN). CRUD identité, recherche, déduplication civile")
        Container(biometric_svc, "Biometric Service", "Python 3.12 + FastAPI + C++ libs", "Capture, traitement, déduplication biométrique (empreintes 10 doigts, iris, facial). Moteur AFIS intégré")
        Container(enrollment_svc, "Enrollment Service", "Go 1.22 + gRPC", "Orchestration du processus d'enrôlement complet: collecte données, validation, soumission, notification")
        Container(document_svc, "Document Service", "Java 21 + Spring Boot", "Génération de documents officiels: carte nationale, passeport, certificats. Intégration imprimerie sécurisée")
        Container(auth_svc, "Authentication Service", "Go 1.22 + Keycloak", "OpenID Connect / OAuth 2.1, gestion sessions, MFA, émission tokens JWT/PASETO")
        Container(notification_svc, "Notification Service", "Node.js 22 + Bull Queue", "SMS, email, push notifications. Gestion des templates, retry, delivery tracking")
        Container(audit_svc, "Audit Service", "Go 1.22 + Kafka consumer", "Collecte immutable de tous les événements systèmes. Écriture sur ledger blockchain privé Hyperledger Fabric")
        Container(admin_portal, "Admin Portal", "React 18 + TypeScript + Vite", "Interface d'administration: supervision système, gestion utilisateurs, tableaux de bord opérationnels")
        Container(citizen_portal, "Portail Citoyen", "Next.js 14 + TypeScript", "Interface web et mobile citoyenne: demandes, suivi, accès dossier personnel, téléchargement documents")
        Container(field_agent_app, "Application Terrain", "Flutter 3.x (iOS/Android/Desktop)", "Application hors-ligne pour agents d'enrôlement. Synchronisation différée, capture biométrique locale")
        Container(interop_gateway, "Interoperability Gateway", "Apache Camel + Quarkus", "Transformation de protocoles, médiation, mappage de données entre SNISID et systèmes gouvernementaux")
        Container(search_svc, "Search Service", "Elasticsearch 8.x + Go", "Recherche full-text et sémantique sur les identités, documents, logs d'audit")
        Container(analytics_svc, "Analytics Service", "Apache Spark + Trino + Superset", "BI et analytics: statistiques démographiques, tableaux de bord opérationnels, rapports gouvernementaux")

        ContainerDb(identity_db, "Identity Database", "PostgreSQL 16 (Patroni HA)", "Données d'identité civile, attributs personnels, historique modifications. Chiffré AES-256-GCM")
        ContainerDb(biometric_db, "Biometric Vault", "PostgreSQL 16 + pgvector", "Templates biométriques chiffrés, index AFIS. Isolation réseau maximale")
        ContainerDb(document_db, "Document Store", "PostgreSQL 16 + MinIO", "Métadonnées documents, images haute résolution, modèles de documents")
        ContainerDb(audit_db, "Audit Ledger", "Hyperledger Fabric + CockroachDB", "Registre immuable d'audit. Multi-DC, consensus PBFT")
        ContainerDb(cache, "Cache Distribué", "Redis Cluster 7.x", "Sessions, tokens, résultats de recherche fréquents, rate limiting counters")
        ContainerDb(message_bus, "Message Bus", "Apache Kafka 3.x (KRaft)", "Event streaming: events d'enrôlement, synchronisation inter-services, audit events")
        ContainerDb(secret_store, "Coffre-fort Secrets", "HashiCorp Vault Enterprise", "Secrets, clés PKI, certificats TLS, credentials bases de données. HSM-backed")
    }

    System_Ext(oec, "OEC", "État Civil")
    System_Ext(mel, "MEL", "Electoral")
    System_Ext(icao, "ICAO PKI", "Passeports")
    System_Ext(interpol_ext, "INTERPOL", "Vérification")

    Rel(citizen, citizen_portal, "Accès portail", "HTTPS TLS 1.3")
    Rel(agent, field_agent_app, "Enrôlement terrain", "App offline/online")
    Rel(admin, admin_portal, "Administration", "HTTPS mTLS, MFA obligatoire")

    Rel(citizen_portal, api_gateway, "Requêtes API", "HTTPS + OAuth 2.1")
    Rel(field_agent_app, api_gateway, "Sync données", "HTTPS mTLS")
    Rel(admin_portal, api_gateway, "API admin", "HTTPS mTLS + MFA")

    Rel(api_gateway, auth_svc, "Validation token", "gRPC mTLS")
    Rel(api_gateway, identity_svc, "Opérations identité", "gRPC mTLS")
    Rel(api_gateway, biometric_svc, "Opérations biométriques", "gRPC mTLS")
    Rel(api_gateway, enrollment_svc, "Orchestration enrôlement", "gRPC mTLS")
    Rel(api_gateway, document_svc, "Génération documents", "gRPC mTLS")

    Rel(identity_svc, identity_db, "Lecture/écriture identités", "PostgreSQL TLS")
    Rel(identity_svc, message_bus, "Publication events", "Kafka TLS SASL")
    Rel(biometric_svc, biometric_db, "Stockage templates", "PostgreSQL TLS")
    Rel(enrollment_svc, message_bus, "Events d'enrôlement", "Kafka TLS SASL")
    Rel(audit_svc, audit_db, "Écriture audit", "gRPC mTLS")
    Rel(auth_svc, cache, "Sessions/tokens", "Redis TLS")
    Rel(interop_gateway, oec, "Synchronisation", "mTLS REST/gRPC")
    Rel(interop_gateway, mel, "Partage données", "mTLS REST")
    Rel(document_svc, icao, "Validation eMRTD", "LDAPS")
    Rel(identity_svc, interpol_ext, "Vérification", "HTTPS REST")
    Rel(identity_svc, secret_store, "Récupération secrets", "Vault API mTLS")
```

### 3.2 Description des Conteneurs

| Conteneur | Technologie | Rôle | Criticité |
|---|---|---|---|
| **API Gateway** | Kong Enterprise 3.x + Envoy | Point d'entrée unique, sécurité périmétrique | CRITIQUE |
| **Identity Service** | Go 1.22, gRPC, REST | Gestionnaire cycle de vie identités NIN | CRITIQUE |
| **Biometric Service** | Python 3.12, FastAPI, C++ | AFIS, capture biométrique, déduplication | CRITIQUE |
| **Enrollment Service** | Go 1.22, gRPC | Orchestration processus d'enrôlement | ÉLEVÉ |
| **Document Service** | Java 21, Spring Boot | Génération documents officiels | ÉLEVÉ |
| **Auth Service** | Go 1.22, Keycloak 23 | OIDC/OAuth 2.1, MFA, sessions | CRITIQUE |
| **Audit Service** | Go 1.22, Kafka consumer | Traçabilité immuable | CRITIQUE |
| **Interop Gateway** | Apache Camel, Quarkus | Médiation inter-systèmes | ÉLEVÉ |
| **Identity DB** | PostgreSQL 16 (Patroni) | Stockage principal identités | CRITIQUE |
| **Biometric Vault** | PostgreSQL 16 + pgvector | Templates biométriques isolés | CRITIQUE |
| **Message Bus** | Apache Kafka 3.x KRaft | Streaming événements | ÉLEVÉ |
| **Secret Store** | HashiCorp Vault Enterprise | PKI, secrets, credentials | CRITIQUE |

---

## 4. C4 Niveau 3 — Composants

### 4.1 Identity Service — Composants Internes

```mermaid
C4Component
    title Identity Service — Composants Internes

    Container_Boundary(identity_svc, "Identity Service (Go 1.22)") {

        Component(nin_manager, "NIN Manager", "Go — internal package", "Génération et validation des Numéros d'Identification Nationale (NIN). Algorithme de Luhn étendu, unicité garantie, vérification de doublons")
        Component(civil_validator, "Civil Data Validator", "Go — internal package", "Validation des données civiles: nom, prénom, date de naissance, lieu de naissance. Règles métier haïtiennes")
        Component(identity_repo, "Identity Repository", "Go — infrastructure layer", "Accès base de données PostgreSQL via pgx/v5. CRUD identités, requêtes complexes, transactions ACID")
        Component(dedup_engine, "Deduplication Engine", "Go — domain service", "Détection des doublons sur données civiles: phonétique haïtienne (Metaphone créole), distance Levenshtein, score de confiance")
        Component(event_publisher, "Event Publisher", "Go — infrastructure", "Publication d'événements sur Kafka: IdentityCreated, IdentityUpdated, IdentityRevoked")
        Component(grpc_handler, "gRPC Handler", "Go — presentation layer", "Endpoints gRPC: CreateIdentity, GetIdentity, UpdateIdentity, RevokeIdentity, SearchIdentities")
        Component(rest_handler, "REST Handler", "Go — presentation layer", "Endpoints REST/JSON pour les clients non-gRPC. Génération OpenAPI 3.1 automatique")
        Component(audit_emitter, "Audit Event Emitter", "Go — cross-cutting", "Émission systématique des événements d'audit pour toute opération sur une identité")
        Component(cache_manager, "Cache Manager", "Go — infrastructure", "Cache Redis pour les identités fréquemment accédées. TTL: 15 min, invalidation sur modification")
        Component(crypto_service, "Cryptography Service", "Go — security", "Chiffrement/déchiffrement données sensibles: AES-256-GCM, signatures ECDSA P-384, gestion clés via Vault")
    }

    ContainerDb(identity_db, "Identity Database", "PostgreSQL 16")
    ContainerDb(cache, "Redis Cluster", "Cache")
    ContainerDb(message_bus, "Kafka", "Message Bus")
    Container(audit_svc, "Audit Service", "Audit")
    Container(biometric_svc, "Biometric Service", "Biométrie")
    Container(vault, "HashiCorp Vault", "Secrets")

    Rel(grpc_handler, nin_manager, "Demande NIN", "function call")
    Rel(grpc_handler, civil_validator, "Validation données", "function call")
    Rel(grpc_handler, dedup_engine, "Vérification doublons", "function call")
    Rel(grpc_handler, identity_repo, "Persistance", "function call")
    Rel(grpc_handler, audit_emitter, "Emission audit", "function call")
    Rel(rest_handler, grpc_handler, "Délègue vers gRPC", "function call")

    Rel(nin_manager, crypto_service, "Signature NIN", "function call")
    Rel(identity_repo, identity_db, "SQL ACID", "pgx/v5 TLS")
    Rel(identity_repo, cache_manager, "Cache lookup", "function call")
    Rel(cache_manager, cache, "Get/Set", "Redis TLS")
    Rel(event_publisher, message_bus, "Publish", "Kafka SASL/TLS")
    Rel(audit_emitter, message_bus, "Audit events", "Kafka SASL/TLS")
    Rel(crypto_service, vault, "Key retrieval", "Vault API mTLS")
    Rel(dedup_engine, biometric_svc, "Vérification biométrique", "gRPC mTLS")
```

### 4.2 Biometric Service — Composants Internes

```mermaid
C4Component
    title Biometric Service — Composants Internes

    Container_Boundary(bio_svc, "Biometric Service (Python 3.12 + FastAPI)") {

        Component(capture_api, "Capture API", "FastAPI + WebSocket", "API de capture en temps réel: empreintes digitales (10 doigts ISO 19794-4), iris (ISO 19794-6), photo faciale (ISO 19794-5)")
        Component(quality_checker, "Quality Assessment Engine", "Python + OpenCV + NFIQ2", "Évaluation qualité biométrique: NFIQ2 pour empreintes (score min 40/100), IREX pour iris, ICAO LDS pour facial")
        Component(template_extractor, "Feature Extraction Engine", "C++ extension (ctypes)", "Extraction de templates biométriques propriétaires. Integration SDK AFIS. ISO 19794-2 minutiae extraction")
        Component(dedup_afis, "AFIS Deduplication", "C++ SDK + Python wrapper", "Moteur AFIS (Automated Fingerprint Identification System): 1:N matching, seuil de correspondance configurable, ranking des candidats")
        Component(liveness_detector, "Liveness Detection", "PyTorch + ONNX", "Anti-spoofing: détection présence vivante pour empreintes (coupures, artefacts), iris (réflexion, dilatation), face (deepfake detection)")
        Component(template_encryptor, "Template Encryption", "Python + cryptography lib", "Chiffrement AES-256-GCM des templates avant stockage. Clés via HSM, rotation annuelle")
        Component(biometric_repo, "Biometric Repository", "Python — SQLAlchemy async", "Accès base biométrique PostgreSQL + pgvector. Stockage templates, index vectoriels pour matching rapide")
        Component(matching_svc, "1:1 Matching Service", "C++ + Python wrapper", "Vérification biométrique 1:1: comparaison template présenté vs template stocké. Score > 0.85 requis")
        Component(iso_validator, "ISO Standards Validator", "Python", "Validation conformité ISO/IEC: 19794-2/4/5/6. Rejet si non-conforme")
        Component(bio_audit, "Biometric Audit Logger", "Python", "Journal spécifique biométrie: qui, quand, quel doigt/iris, score matching, résultat")
    }

    ContainerDb(biometric_db, "Biometric Vault", "PostgreSQL + pgvector")
    Container(vault, "HashiCorp Vault", "HSM-backed keys")
    Container(audit_svc, "Audit Service", "Immutable log")

    Rel(capture_api, quality_checker, "Vérification qualité", "sync call")
    Rel(capture_api, iso_validator, "Validation ISO", "sync call")
    Rel(capture_api, liveness_detector, "Anti-spoofing", "async ML inference")
    Rel(quality_checker, template_extractor, "Extraction si qualité OK", "conditional call")
    Rel(template_extractor, dedup_afis, "Déduplication 1:N", "C++ interop")
    Rel(template_extractor, template_encryptor, "Chiffrement avant stockage", "function call")
    Rel(template_encryptor, vault, "Clé de chiffrement", "Vault API mTLS")
    Rel(template_encryptor, biometric_repo, "Stockage chiffré", "function call")
    Rel(biometric_repo, biometric_db, "Persistance", "SQLAlchemy async TLS")
    Rel(matching_svc, biometric_repo, "Récupération template", "function call")
    Rel(matching_svc, dedup_afis, "Score matching", "C++ interop")
    Rel(bio_audit, audit_svc, "Emission logs", "gRPC mTLS")
```

### 4.3 API Gateway — Composants Internes

```mermaid
C4Component
    title API Gateway — Composants Internes (Kong Enterprise + Envoy)

    Container_Boundary(api_gw, "API Gateway Layer") {

        Component(tls_terminator, "TLS Terminator", "Envoy Proxy", "Terminaison TLS 1.3, inspection certificat client (mTLS), validation SNI, HSTS enforcement")
        Component(rate_limiter, "Rate Limiter", "Kong Plugin (Redis-backed)", "Limitation taux: par IP, par client OAuth, par agence. Algorithme token bucket. Redis Cluster pour état distribué")
        Component(auth_filter, "Auth/AuthZ Filter", "Kong Plugin + OPA", "Validation tokens JWT/PASETO, introspection OIDC, délégation AuthZ à OPA pour ABAC")
        Component(request_validator, "Request Validator", "Kong Plugin", "Validation schéma OpenAPI 3.1, sanitisation inputs, rejet requêtes malformées, Content-Type enforcement")
        Component(router, "Smart Router", "Kong — declarative config", "Routage vers services backend: version-aware, canary deployment, blue/green routing, health-check based")
        Component(transformer, "Request/Response Transformer", "Kong Plugin", "Transformation headers, masquage données sensibles en réponse, ajout correlation-id, enrichissement contexte")
        Component(circuit_breaker, "Circuit Breaker", "Envoy — Outlier Detection", "Protection contre cascading failures: threshold 50% errors/5s, half-open après 30s, métriques exposées")
        Component(logging_plugin, "Access Log / Tracing", "Kong + OpenTelemetry", "Génération logs structurés JSON, traces distribuées (W3C TraceContext), export vers Jaeger/Tempo")
        Component(waf_plugin, "WAF — Web App Firewall", "Kong + ModSecurity (OWASP CRS)", "Protection OWASP Top 10, règles SNISID custom, blocage SQLi/XSS/Path Traversal")
        Component(cache_plugin, "Response Cache", "Kong Plugin (Redis)", "Cache des réponses GET non-sensibles: TTL 60s, invalidation sur modification, Vary headers respectés")
    }

    Container(auth_svc, "Auth Service", "OIDC/OAuth 2.1")
    Container(opa, "OPA Policy Engine", "ABAC Policies")
    Container(identity_svc, "Identity Service", "Backend")
    Container(biometric_svc, "Biometric Service", "Backend")

    Rel(tls_terminator, auth_filter, "Requête authentifiée", "internal")
    Rel(auth_filter, auth_svc, "Token introspection", "gRPC mTLS")
    Rel(auth_filter, opa, "Authorization check", "HTTP REST")
    Rel(auth_filter, waf_plugin, "Requête autorisée", "internal")
    Rel(waf_plugin, rate_limiter, "Requête filtrée", "internal")
    Rel(rate_limiter, request_validator, "Non throttlée", "internal")
    Rel(request_validator, transformer, "Requête valide", "internal")
    Rel(transformer, router, "Requête enrichie", "internal")
    Rel(router, circuit_breaker, "Route sélectionnée", "internal")
    Rel(circuit_breaker, identity_svc, "Forwarding", "gRPC mTLS")
    Rel(circuit_breaker, biometric_svc, "Forwarding", "gRPC mTLS")
    Rel(logging_plugin, transformer, "Intercept all", "sidecar pattern")
```

---

## 5. C4 Niveau 4 — Code et Chemins Critiques

### 5.1 Chemin Critique 1 — Enrôlement Biométrique Complet

```mermaid
sequenceDiagram
    participant Agent as Agent Terrain (App Flutter)
    participant GW as API Gateway (Kong+Envoy)
    participant Auth as Auth Service
    participant Enroll as Enrollment Service
    participant Identity as Identity Service
    participant Bio as Biometric Service
    participant Vault as HashiCorp Vault
    participant DB as Identity DB (PostgreSQL)
    participant BioDB as Biometric Vault
    participant Kafka as Apache Kafka
    participant Audit as Audit Service
    participant Notif as Notification Service

    Note over Agent,Notif: FLUX: Enrôlement citoyen complet — Nouveau NIN

    Agent->>GW: POST /v1/enrollment/start (mTLS, JWT agent)
    GW->>Auth: Introspection token agent
    Auth-->>GW: Token valide, rôles=[AGENT_ENROLL]
    GW->>GW: WAF check, rate limit, schema validation
    GW->>Enroll: gRPC EnrollmentStart(agentId, stationId)
    Enroll->>Enroll: Génération enrollmentSessionId (UUID v7)
    Enroll->>Kafka: Publish EnrollmentStarted{sessionId, agentId, timestamp}
    Enroll-->>GW: EnrollmentSession{sessionId, expiresAt}
    GW-->>Agent: 201 Created {sessionId, expiresAt: +30min}

    Note over Agent,Bio: Phase 1: Données Civiles

    Agent->>GW: POST /v1/enrollment/{sessionId}/civil-data
    GW->>Identity: ValidateCivilData(nom, prénom, ddn, lieu)
    Identity->>Identity: Règles métier: nom non-vide, ddn passé, lieu haïtien
    Identity->>Identity: Déduplication civile: Levenshtein + Metaphone créole
    Identity-->>GW: ValidationResult{isValid, duplicateCandidates[]}
    GW-->>Agent: 200 {valid: true, warnings: []}

    Note over Agent,Bio: Phase 2: Capture Biométrique

    loop Pour chaque doigt (10 tentatives max par doigt)
        Agent->>GW: POST /v1/biometric/{sessionId}/fingerprint (ISO 19794-4)
        GW->>Bio: CaptureFingerprint(fingerData, position, sessionId)
        Bio->>Bio: NFIQ2 Quality Check (min score: 40)
        alt Qualité insuffisante
            Bio-->>GW: 400 {quality: 35, retry: true}
            GW-->>Agent: 400 Qualité insuffisante, réessayer
        else Qualité acceptée
            Bio->>Bio: Liveness detection (anti-spoofing)
            Bio->>Bio: ISO 19794-2 minutiae extraction
            Bio->>Bio: 1:N AFIS deduplication search
            alt Doublon biométrique détecté
                Bio-->>GW: 409 {duplicateNIN: "HTI-XXXX", confidence: 0.98}
                GW-->>Agent: 409 Identité existante détectée
            else Unique
                Bio->>Vault: GetEncryptionKey(keyId="biometric-2026")
                Vault-->>Bio: AES-256-GCM key (ephemeral)
                Bio->>Bio: Encrypt(template, key)
                Bio->>BioDB: INSERT encrypted_template
                Bio-->>GW: 200 {fingerprintId, quality: 87}
                GW-->>Agent: 200 OK
            end
        end
    end

    Note over Agent,Bio: Phase 3: Capture Iris et Photo

    Agent->>GW: POST /v1/biometric/{sessionId}/iris (ISO 19794-6)
    GW->>Bio: CaptureIris(irisData, eye=BOTH, sessionId)
    Bio->>Bio: IREX quality check + liveness
    Bio->>BioDB: INSERT iris_template_encrypted
    Bio-->>GW: 200 {irisId, quality: 92}

    Agent->>GW: POST /v1/biometric/{sessionId}/photo (ICAO LDS)
    GW->>Bio: CapturePhoto(imageData, sessionId)
    Bio->>Bio: ICAO compliance check (frontal, neutral, quality)
    Bio->>BioDB: INSERT photo_template_encrypted
    Bio-->>GW: 200 {photoId, icaoCompliant: true}

    Note over Enroll,DB: Phase 4: Finalisation et Attribution NIN

    Agent->>GW: POST /v1/enrollment/{sessionId}/finalize
    GW->>Enroll: FinalizeEnrollment(sessionId)
    Enroll->>Identity: CreateIdentity(civilData, biometricRefs)
    Identity->>Identity: Génération NIN (algo Luhn étendu + checksum)
    Identity->>Vault: SignNIN(nin, privateKey="identity-signing-2026")
    Vault-->>Identity: Signature ECDSA P-384
    Identity->>DB: BEGIN TRANSACTION
    Identity->>DB: INSERT identity{nin, civilData, biometricRefs, signature}
    Identity->>DB: INSERT identity_history{action=CREATED, actor=agentId}
    Identity->>DB: COMMIT
    Identity->>Kafka: Publish IdentityCreated{nin, sessionId, timestamp}
    Identity-->>Enroll: Identity{nin, createdAt}
    Enroll->>Kafka: Publish EnrollmentCompleted{sessionId, nin}
    Enroll->>Notif: NotifyCitizen{nin, channel=SMS, phone}
    Notif->>Notif: Enqueue SMS via Digicel/Natcom gateway
    Enroll-->>GW: EnrollmentResult{nin, status=COMPLETE}
    GW-->>Agent: 201 Created {nin: "HTI-2026-XXXXXXXX", status: "ACTIVE"}

    Note over Audit: Audit Continu
    Kafka->>Audit: Consume all events
    Audit->>Audit: Write to Hyperledger Fabric ledger
```

### 5.2 Chemin Critique 2 — Authentification et Vérification d'Identité 1:1

```mermaid
sequenceDiagram
    participant Client as Application Consommatrice
    participant GW as API Gateway
    participant Auth as Auth Service (Keycloak)
    participant OPA as OPA Policy Engine
    participant Identity as Identity Service
    participant Bio as Biometric Service
    participant BioDB as Biometric Vault
    participant Audit as Audit Service

    Note over Client,Audit: FLUX: Vérification identité 1:1 par agence gouvernementale

    Client->>Auth: POST /auth/token (client_credentials, client_id, client_secret)
    Auth->>Auth: Validation client credentials OAuth 2.1
    Auth->>Auth: Chargement rôles/scopes: [identity:verify, biometric:match:1to1]
    Auth-->>Client: access_token (JWT, exp: 900s), token_type=Bearer

    Client->>GW: POST /v1/identity/verify {nin, biometricData, purpose}
    Note right of GW: mTLS: client certificate validé
    GW->>GW: TLS termination, cert validation
    GW->>GW: JWT signature verification (JWKS endpoint)
    GW->>OPA: AuthZCheck{subject, action=identity:verify, resource=nin, purpose, clientId}

    OPA->>OPA: Evaluate policy identity_verification.rego
    Note right of OPA: Policy: purpose must be in allowed_purposes[clientId]
    Note right of OPA: Policy: time must be within working_hours OR emergency_override
    Note right of OPA: Policy: data_minimization — only return authorized fields

    alt Accès refusé par OPA
        OPA-->>GW: {allow: false, reason: "purpose_not_authorized"}
        GW-->>Client: 403 Forbidden {code: "POLICY_VIOLATION", reason: "purpose_not_authorized"}
    else Accès autorisé
        OPA-->>GW: {allow: true, allowedFields: ["nin", "name", "photo"], maskFields: ["address"]}

        GW->>Identity: VerifyIdentity(nin)
        Identity->>Identity: Récupération identité par NIN (cache-first)
        alt Identité non trouvée
            Identity-->>GW: NotFoundError
            GW-->>Client: 404 {code: "NIN_NOT_FOUND"}
        else Identité trouvée
            Identity->>Identity: Check statut: ACTIVE ? Sinon REVOKED/SUSPENDED
            alt Identité révoquée
                Identity-->>GW: RevokedError{reason, revokedAt}
                GW-->>Client: 410 Gone {code: "IDENTITY_REVOKED", since: "date"}
            else Identité active
                Identity-->>GW: IdentityResult{civilData, biometricRefs}

                GW->>Bio: MatchBiometric1to1(biometricRefs.fingerprintId, presentedTemplate)
                Bio->>BioDB: GetEncryptedTemplate(fingerprintId)
                Bio->>Bio: Decrypt template (Vault key retrieval)
                Bio->>Bio: 1:1 AFIS matching
                Bio->>Bio: Liveness check sur template présenté

                alt Score matching < 0.85
                    Bio-->>GW: MatchResult{matched: false, score: 0.72}
                    GW-->>Client: 200 {verified: false, matchScore: null}
                    Note right of GW: Score masqué si non-match pour sécurité
                else Score >= 0.85
                    Bio-->>GW: MatchResult{matched: true, score: 0.97}
                    GW->>GW: Appliquer OPA field masking
                    GW-->>Client: 200 {verified: true, identity: {nin, name, photo}, matchScore: null}
                    Note right of GW: Score toujours masqué (anti-gaming)
                end
            end
        end
    end

    GW->>Audit: LogVerification{clientId, nin, result, purpose, timestamp, ip}
    Audit->>Audit: Write immutable record to Hyperledger
```

### 5.3 Structure de Code — Identity Service (Architecture Hexagonale)

```
identity-service/
├── cmd/
│   └── server/
│       └── main.go                    # Point d'entrée: config, wiring DI, graceful shutdown
├── internal/
│   ├── domain/                        # Domaine métier pur (aucune dépendance externe)
│   │   ├── identity/
│   │   │   ├── entity.go              # Identity struct, NIN type, statuts
│   │   │   ├── repository.go          # Interface IdentityRepository (port)
│   │   │   ├── service.go             # IdentityDomainService: règles métier
│   │   │   ├── nin.go                 # Génération/validation NIN haïtien
│   │   │   ├── events.go              # Domain events: IdentityCreated, etc.
│   │   │   └── errors.go              # Domain errors: DuplicateNIN, etc.
│   │   └── deduplication/
│   │       ├── civil.go               # Algorithme déduplication civile
│   │       └── metaphone_creole.go    # Phonétique créole haïtienne
│   ├── application/                   # Use cases (orchestration)
│   │   ├── create_identity.go         # Use case: CreateIdentity
│   │   ├── verify_identity.go         # Use case: VerifyIdentity
│   │   ├── revoke_identity.go         # Use case: RevokeIdentity
│   │   └── search_identities.go       # Use case: SearchIdentities
│   ├── infrastructure/                # Adaptateurs externes
│   │   ├── persistence/
│   │   │   ├── postgres_repo.go       # Implémentation PostgreSQL du repository
│   │   │   ├── migrations/            # Fichiers Flyway SQL
│   │   │   └── queries/               # SQL queries (sqlc generated)
│   │   ├── cache/
│   │   │   └── redis_cache.go         # Cache Redis avec chiffrement côté client
│   │   ├── messaging/
│   │   │   └── kafka_publisher.go     # Publication events Kafka
│   │   ├── vault/
│   │   │   └── vault_client.go        # Intégration HashiCorp Vault
│   │   └── biometric/
│   │       └── bio_client.go          # Client gRPC Biometric Service
│   └── interfaces/                    # Interfaces entrantes
│       ├── grpc/
│       │   ├── handler.go             # gRPC server handlers
│       │   └── interceptors/          # Auth, logging, tracing interceptors
│       └── rest/
│           ├── handler.go             # REST handlers (chi router)
│           └── middleware/            # Auth, CORS, rate limit middleware
├── api/
│   └── proto/
│       └── identity/v1/               # Protobuf definitions
│           └── identity.proto
├── configs/
│   ├── config.yaml                    # Configuration base
│   └── config.production.yaml        # Overrides production
└── deployments/
    ├── kubernetes/                    # Manifestes K8s
    └── helm/                          # Chart Helm
```

---

## 6. Décisions Technologiques

### 6.1 Matrice de Décisions Architecturales (ADR)

| ADR-ID | Composant | Décision | Alternative Rejetée | Raison |
|---|---|---|---|---|
| ADR-001 | API Gateway | Kong Enterprise + Envoy | Nginx, HAProxy, AWS API GW | Kong: plugins riches, Envoy: service mesh natif, souveraineté (on-prem) |
| ADR-002 | Identity Service | Go 1.22 | Java Spring, Python | Performance, faible empreinte mémoire, concurrence native (goroutines) |
| ADR-003 | Biometric Service | Python 3.12 + C++ | Go, Java | Python: écosystème ML (PyTorch, OpenCV), C++: performance AFIS critique |
| ADR-004 | Identity DB | PostgreSQL 16 (Patroni) | MySQL, Oracle, MongoDB | ACID strict, JSON natif, PostGIS, Patroni HA, souveraineté (open source) |
| ADR-005 | Message Bus | Apache Kafka 3.x KRaft | RabbitMQ, Pulsar | Volume élevé, rétention configurable, KRaft (sans ZooKeeper), replay events |
| ADR-006 | Auth | Keycloak 23 + Go service | Auth0, Okta, custom | Keycloak: OIDC complet, on-prem, souveraineté, auditable |
| ADR-007 | Secrets | HashiCorp Vault Enterprise | AWS Secrets Manager | On-prem, HSM integration, PKI native, audit complet |
| ADR-008 | Container Orch. | Kubernetes 1.30 (RKE2) | OpenShift, Nomad | RKE2: souveraineté (Rancher/SUSE), FIPS mode, gouvernement-ready |
| ADR-009 | Service Mesh | Istio 1.22 | Linkerd, Consul Connect | mTLS automatique, politique réseau fine, télémétrie complète |
| ADR-010 | Observabilité | Prometheus + Grafana + Loki + Tempo | Datadog, New Relic | Stack open source, on-prem, pas d'exfiltration télémétrie vers l'étranger |
| ADR-011 | CI/CD | GitLab CI + ArgoCD | Jenkins, GitHub Actions | Self-hosted, GitOps natif, RBAC pipeline, souveraineté |
| ADR-012 | Biometric DB | PostgreSQL + pgvector | Cassandra, Elasticsearch | pgvector: recherche vectorielle native, cohérence transactionnelle |

### 6.2 Stack Technologique Complet

```yaml
# SNISID Technology Stack — Version 1.0 (2026)
stack:
  languages:
    primary:
      - language: Go
        version: "1.22"
        usage: [identity-service, auth-service, audit-service, enrollment-service]
        justification: "Performance, concurrence, faible latence, déploiement statique"
      - language: Python
        version: "3.12"
        usage: [biometric-service, analytics-service, ML models]
        justification: "Écosystème ML/AI, bibliothèques biométriques, data science"
      - language: Java
        version: "21 LTS"
        usage: [document-service, interop-gateway]
        justification: "Spring Boot mature, Camel intégration, JVM stable"
      - language: TypeScript
        version: "5.x"
        usage: [citizen-portal-nextjs, admin-portal-react, notification-service-nodejs]

  frameworks:
    backend:
      - {name: gRPC, version: "1.64", usage: "Communication inter-services"}
      - {name: FastAPI, version: "0.111", usage: "REST API Python (Biometric Service)"}
      - {name: Spring Boot, version: "3.3", usage: "Document & Interop Services"}
      - {name: Apache Camel, version: "4.x", usage: "Intégration & médiation"}
      - {name: Quarkus, version: "3.x", usage: "Native image Interop Gateway"}
    frontend:
      - {name: Next.js, version: "14", usage: "Portail citoyen SSR/SSG"}
      - {name: React, version: "18", usage: "Admin portal SPA"}
      - {name: Flutter, version: "3.22", usage: "App terrain iOS/Android/Desktop"}

  databases:
    primary:
      - {name: PostgreSQL, version: "16", ha: "Patroni + etcd", usage: "Identity, Document, Auth"}
      - {name: CockroachDB, version: "23.x", usage: "Audit ledger distribué multi-DC"}
    cache:
      - {name: Redis, version: "7.x", mode: "Cluster", usage: "Sessions, cache, rate limiting"}
    search:
      - {name: Elasticsearch, version: "8.x", usage: "Full-text search, logs"}
    analytics:
      - {name: Apache Spark, version: "3.5", usage: "Batch analytics"}
      - {name: Trino, version: "450", usage: "Requêtes analytiques SQL"}

  infrastructure:
    orchestration: "Kubernetes 1.30 (RKE2 — SUSE Rancher, FIPS-enabled)"
    service_mesh: "Istio 1.22 (mTLS strict mode)"
    api_gateway: "Kong Enterprise 3.7 + Envoy 1.30"
    secrets: "HashiCorp Vault Enterprise 1.16 (HSM-backed)"
    ci_cd: "GitLab CE 17.x + ArgoCD 2.11"
    storage: "Ceph 18 (Reef) + MinIO"
    messaging: "Apache Kafka 3.7 (KRaft mode)"
    pki: "EJBCA Enterprise 8.x + HashiCorp Vault PKI"
    monitoring: "Prometheus 2.52 + Grafana 11 + Loki 3 + Tempo 2"
    tracing: "Jaeger 1.57 + OpenTelemetry 0.10"
    logging: "Fluentbit + Loki + Elasticsearch"
    policy: "OPA (Open Policy Agent) 0.65 + Gatekeeper 3.16"
```

---

## 7. Exigences Non-Fonctionnelles par Composant

### 7.1 Tableau NFR Complet

| Composant | Disponibilité | Latence P99 | Débit | RTO | RPO | Criticité |
|---|---|---|---|---|---|---|
| **API Gateway** | 99.99% | < 50 ms | 10 000 req/s | 2 min | N/A | CRITIQUE |
| **Identity Service** | 99.95% | < 200 ms | 2 000 req/s | 5 min | 1 min | CRITIQUE |
| **Biometric Service** | 99.95% | < 3 s (capture), < 500 ms (verify) | 500 req/s | 10 min | 5 min | CRITIQUE |
| **Auth Service** | 99.99% | < 100 ms | 5 000 req/s | 2 min | 30 s | CRITIQUE |
| **Enrollment Service** | 99.90% | < 5 s (end-to-end) | 200 req/s | 15 min | 5 min | ÉLEVÉ |
| **Document Service** | 99.90% | < 30 s (génération) | 100 req/s | 30 min | 15 min | ÉLEVÉ |
| **Notification Service** | 99.50% | < 10 s (delivery) | 1 000 notif/s | 60 min | 30 min | MOYEN |
| **Audit Service** | 99.99% | < 100 ms (write async) | 5 000 events/s | 2 min | 0 (sync) | CRITIQUE |
| **Identity DB** | 99.99% | < 10 ms (read), < 50 ms (write) | 5 000 TPS | 5 min | 1 min | CRITIQUE |
| **Biometric Vault** | 99.95% | < 20 ms (read) | 1 000 TPS | 10 min | 5 min | CRITIQUE |
| **Kafka Cluster** | 99.95% | < 10 ms (produce) | 100 000 msg/s | 5 min | 0 (sync replica) | ÉLEVÉ |
| **Search Service** | 99.90% | < 500 ms | 500 req/s | 30 min | 60 min | MOYEN |

### 7.2 NFR Sécurité — Tous Composants

```yaml
security_requirements:
  encryption:
    in_transit:
      protocol: "TLS 1.3 minimum (TLS 1.2 toléré pour legacy)"
      cipher_suites:
        - TLS_AES_256_GCM_SHA384
        - TLS_CHACHA20_POLY1305_SHA256
      certificate_validity: "90 jours maximum (rotation automatique)"
      mtls: "Obligatoire pour communication inter-services"
    at_rest:
      algorithm: "AES-256-GCM"
      key_management: "HashiCorp Vault (HSM-backed)"
      key_rotation: "Annuelle (biométrie), trimestrielle (identité)"

  authentication:
    human_users:
      method: "OIDC/OAuth 2.1 via Keycloak"
      mfa: "Obligatoire pour tous les accès (TOTP/FIDO2)"
      session_duration: "8h maximum (agents), 30min (admin)"
    services:
      method: "SPIFFE/SPIRE workload identity + mTLS"
      rotation: "Token SPIFFE: 1h, certificats mTLS: 24h"

  authorization:
    model: "ABAC (Attribute-Based Access Control) via OPA"
    principles: "Least privilege, need-to-know"
    review: "Revue trimestrielle des politiques"

  audit:
    requirement: "100% des accès aux données personnelles tracés"
    retention: "7 ans minimum (données civiles), 3 ans (logs opérationnels)"
    tamper_proof: "Hyperledger Fabric ledger immutable"
    fields_required: [who, when, what, where, why, result]
```

### 7.3 NFR Performance — Scalabilité

```yaml
scalability_requirements:
  current_phase_0:
    enrolled_identities: 5_000_000  # Phase 0: ~50% population
    daily_transactions: 50_000
    concurrent_users: 500
    peak_multiplier: 5x             # Période électorale

  target_phase_2:
    enrolled_identities: 12_000_000  # Population totale
    daily_transactions: 500_000
    concurrent_users: 5_000
    peak_multiplier: 10x

  auto_scaling:
    identity_service:
      min_replicas: 3
      max_replicas: 20
      scale_trigger: "CPU > 70% OU latency P95 > 200ms"
    biometric_service:
      min_replicas: 3
      max_replicas: 10
      scale_trigger: "Queue depth > 100 OU CPU > 80%"
    api_gateway:
      min_replicas: 3
      max_replicas: 15
      scale_trigger: "RPS > 7000"
```

---

## 8. Contraintes Architecturales

### 8.1 Contraintes Souveraineté

| Contrainte | Description | Impact Architectural |
|---|---|---|
| **Données sur territoire haïtien** | Toutes données biométriques et PII doivent rester en Haïti | Pas de cloud public pour données sensibles. Datacenters souverains obligatoires |
| **Clés cryptographiques en Haïti** | HSM physiquement en territoire haïtien | HSM on-premise FIPS 140-2 L3 dans datacenters PAP et CAP |
| **Code source auditable** | Tout code doit être auditable par l'État haïtien | Open source préféré, escrow source code pour solutions propriétaires |
| **Fournisseurs diversifiés** | Pas de lock-in fournisseur unique >40% | Architecture multi-fournisseurs, standards ouverts |

### 8.2 Contraintes Infrastructure Haïtienne

| Contrainte | Impact | Mitigation |
|---|---|---|
| **Électricité instable** | Pannes fréquentes (délestage) | UPS 4h + générateurs diesel, mode offline obligatoire |
| **Connectivité limitée** | Internet < 10 Mbps dans certaines zones | Synchronisation différée, compression, delta sync |
| **Canicule tropicale** | T° datacenter difficile à maintenir | Climatisation redondante N+1, monitoring thermique |
| **Risques sismiques** | Haïti zone sismique active | Datacenters sur fondations parasismiques, DR offshore |
| **Capacités techniques locales** | Ressources humaines limitées | Formation obligatoire, documentation en créole/français |

---

## 9. Perspectives d'Évolution

### 9.1 Roadmap Architecture

```mermaid
timeline
    title SNISID — Évolution Architecture
    section Phase 0 (2026)
        Q1 : Infrastructure de base
           : Identity + Biometric Services
           : 5 agences pilotes
        Q2 : Portail citoyen v1
           : PKI opérationnelle
           : SOC actif
        Q3 : 15 agences intégrées
           : DR offshore actif
           : Audit blockchain
        Q4 : 5M identités
           : Full HA Active-Active
           : API publique v1

    section Phase 1 (2027)
        Q1 : eID Card NFC
           : Passeport biométrique
           : Mobile ID
        Q2 : Interopérabilité CARICOM
           : UNHCR integration
           : diaspora services
        Q3 : 10M identités
           : Analytics avancés
           : BI gouvernemental
        Q4 : Certification ISO 27001
           : Audit international
```

---

## Bloc d'Approbation / Approval Block

| Rôle | Nom | Signature | Date |
|---|---|---|---|
| **Architecte en Chef** | [À compléter] | [Signature] | 2026-05-25 |
| **Architecte Sécurité** | [À compléter] | [Signature] | 2026-05-25 |
| **Directeur Infrastructure** | [À compléter] | [Signature] | 2026-05-25 |
| **Directeur Général SNISID** | [À compléter] | [Signature] | 2026-05-25 |
| **CISO** | [À compléter] | [Signature] | 2026-05-25 |

---

*Document SNISID-ARC-C4-001 v1.0.0 — RESTREINT — © République d'Haïti, Programme SNISID, 2026*
*Toute reproduction ou diffusion non autorisée est interdite par la loi haïtienne.*
