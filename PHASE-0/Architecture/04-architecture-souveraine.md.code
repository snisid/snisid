# 🏗️ SNISID — Sovereign Architecture Blueprint

**Document N° :** SNISID-ARC-004
**Étape Phase 0 :** 4/16
**Principe :** *Cloud-native, mais souverain.*

---

## 1. Principes Architecturaux (10)

1. **Cloud-native sur infrastructure souveraine** (pas de hyperscaler étranger pour données critiques)
2. **Microservices** — services autonomes, déployables indépendamment
3. **Kubernetes** — orchestration standard, portable
4. **API-first** — toute fonctionnalité exposée par API documentée OpenAPI 3.x
5. **Event-driven** — bus d'événements national (Kafka)
6. **Zero Trust** — aucune confiance implicite, vérification systématique
7. **Offline-first** — fonctionnement local + synchronisation
8. **Polyglot persistence** — bonne base pour bon usage (PostgreSQL, MongoDB, Redis, Elasticsearch, Object Storage)
9. **Observabilité native** — OpenTelemetry, logs/metrics/traces
10. **Infrastructure-as-Code** — Terraform, Ansible, GitOps

---

## 2. Vue Macro (5 couches)

```
┌────────────────────────────────────────────────────────────┐
│  L5 — EXPÉRIENCE              Apps citoyennes, portail,    │
│                                guichets, agents, mobile     │
├────────────────────────────────────────────────────────────┤
│  L4 — API GATEWAY              Kong/APISIX + WAF + OAuth2  │
├────────────────────────────────────────────────────────────┤
│  L3 — MICROSERVICES MÉTIER     Identité, État Civil,       │
│                                Biométrie, Justice, Police  │
├────────────────────────────────────────────────────────────┤
│  L2 — PLATEFORME               Kubernetes, Service Mesh,   │
│                                Kafka, PostgreSQL, MinIO    │
├────────────────────────────────────────────────────────────┤
│  L1 — INFRASTRUCTURE           Datacenter souverain        │
│                                + Edge nodes + Offline kits │
└────────────────────────────────────────────────────────────┘
```

---

## 3. Domaines Fonctionnels (Bounded Contexts DDD)

| Domaine | Microservices clés |
|---------|---------------------|
| **Identité** | Enrolment, KYC, Lifecycle, Federation, Consent |
| **Biométrie** | Capture, Quality, AFIS 1:N, Matching 1:1, Liveness |
| **État Civil** | Naissance, Mariage, Divorce, Décès, Adoption |
| **Justice** | Casier judiciaire, Jugements, Notifications |
| **Police** | Personnes recherchées, Plaintes, Interpellations |
| **Immigration** | Visas, Frontières, Passeports |
| **Notifications** | SMS, Email, Push, USSD |
| **Documents** | Génération, Signature, Archivage légal |
| **Audit** | Journalisation immuable, Forensics |

---

## 4. Stack Technique de Référence

| Couche | Technologie |
|--------|-------------|
| Orchestration | Kubernetes (vanilla ou RKE2/k3s pour edge) |
| Service Mesh | Istio ou Linkerd |
| API Gateway | Kong OSS ou APISIX |
| Event bus | Apache Kafka + Schema Registry |
| RDBMS | PostgreSQL 16+ (HA via Patroni) |
| NoSQL | MongoDB (documents), Redis (cache/session) |
| Recherche | Elasticsearch / OpenSearch |
| Object storage | MinIO (S3-compatible souverain) |
| Identity & Access | Keycloak (OIDC/SAML) + PKI nationale |
| Secrets | HashiCorp Vault |
| CI/CD | GitLab CE + ArgoCD (GitOps) |
| Observabilité | Prometheus + Grafana + Loki + Tempo (OTel) |
| SIEM | Wazuh ou Elastic Security |
| Backups | Velero + Restic vers stockage chiffré |
| IaC | Terraform + Ansible |
| Langages | Java/Spring (services critiques), Go (perf), Python (data/IA), TypeScript (front) |
| Front | React + Next.js + PWA pour offline |
| Mobile | Flutter (multi-plateforme) ou React Native |

> **Critère de choix :** open-source, communauté forte, pas de dépendance vendor critique, expertise mobilisable en Haïti.

---

## 5. Zero Trust Architecture

Principes appliqués :
- **Never trust, always verify** — chaque appel API authentifié + autorisé
- **mTLS partout** — service-to-service chiffré via service mesh
- **Identity-based access** — Keycloak comme IdP unique
- **Microsegmentation réseau** — NetworkPolicies Kubernetes
- **Least privilege** — RBAC strict, rotation secrets, JIT access
- **Continuous verification** — re-évaluation contextuelle (device, géo, comportement)

---

## 6. Patterns Architecturaux Adoptés

- **CQRS** — séparation lecture/écriture pour services à forte charge
- **Event Sourcing** — pour audit immuable (état civil, justice)
- **Saga pattern** — transactions distribuées (ex. enrôlement multi-étapes)
- **Circuit Breaker** — résilience (Resilience4j)
- **Backend-for-Frontend (BFF)** — agrégation API par canal
- **Outbox pattern** — fiabilité event publishing
- **Strangler Fig** — migration progressive des systèmes legacy (ONI, CNIGS)

---

## 7. Modèle de Déploiement

```
        ┌───────────────────────┐
        │  PRIMARY DATACENTER   │  Port-au-Prince
        │  (Production active)  │
        └──────────┬────────────┘
                   │ Replication synchrone (RPO ~0)
        ┌──────────▼────────────┐
        │  DR DATACENTER        │  Cap-Haïtien
        │  (Hot standby)        │
        └──────────┬────────────┘
                   │
   ┌───────────────┼───────────────┐
   │               │               │
┌──▼──┐         ┌──▼──┐         ┌──▼──┐
│Edge │         │Edge │         │Edge │   10 nodes
│Dépt │   ...   │Dépt │         │Dépt │   départementaux
└──┬──┘         └──┬──┘         └──┬──┘
   │               │               │
┌──▼─────────┐  ┌──▼─────────┐
│Offline Kits│  │Offline Kits│   Sections communales
│ + Mobiles  │  │ + Mobiles  │
└────────────┘  └────────────┘
```

---

## 8. Sécurité Architecturale

- **Chiffrement at-rest** : LUKS disques + chiffrement applicatif données sensibles (AES-256-GCM)
- **Chiffrement in-transit** : TLS 1.3 obligatoire, mTLS interne
- **HSM** (Hardware Security Module) pour PKI nationale (FIPS 140-2 niveau 3)
- **Tokenisation** des identifiants nationaux dans bases secondaires
- **Pseudonymisation** systématique en environnements hors prod

---

## 9. Performance & SLO Cibles

| Service | Latence p95 | Disponibilité | Throughput |
|---------|-------------|---------------|------------|
| Vérification identité 1:1 | < 500 ms | 99,9 % | 500 TPS |
| Match biométrique 1:N (10M) | < 3 s | 99,5 % | 50 TPS |
| Génération acte naissance | < 2 s | 99,9 % | 100 TPS |
| Enrôlement biométrique | < 10 s | 99 % | — |
| Sync offline kit | < 30 min / 5000 records | 99 % | — |

---

## 10. Évolutivité & Multi-tenant

- Chaque agence = tenant logique avec ses droits et données
- Quotas et throttling par tenant
- Possibilité d'isolation forte (namespace K8s dédié) pour agences sensibles (Renseignement, Justice)

---
*Fin du document — Étape 4/16*
