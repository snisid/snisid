# SNISID-Core — Plateforme Souveraine Nationale
## Phase 1 : National Core Architecture & Sovereign Platform Engineering

**Document ID :** SNISID-CORE-README-001  
**Version :** 1.0.0  
**Date :** Mai 2026  
**Standard :** Production-grade | Kubernetes-native | GitOps-managed | Zero Trust | Offline-first

---

## 🎯 Mission

SNISID-Core est le **noyau technique souverain** qui supporte l'ensemble des services gouvernementaux d'Haïti. C'est une plateforme cloud-native, entièrement bootable, conçue pour :

- **Haute disponibilité 24/7** — Active-Active Port-au-Prince ↔ Cap-Haïtien
- **Offline-first** — Fonctionnement 30+ jours sans connectivité centrale
- **Zero Trust** — Aucune confiance implicite, vérification continue
- **GitOps** — ArgoCD comme seule source de vérité

---

## 🏗️ Stack Technique Souveraine

| Composant | Technologie | Justification |
|-----------|------------|---------------|
| **Orchestration** | RKE2 (Rancher) | Distribution K8s souveraine, FIPS 140-2 |
| **Service Mesh** | Istio | mTLS strict, observabilité, trafic shaping |
| **CNI** | Cilium (eBPF) | Network policies avancées, BPF performance |
| **API Gateway** | Kong | Plugin ecosystem, OAuth 2.1, rate limiting |
| **IAM** | Keycloak + HashiCorp Vault | OIDC/SAML + secrets dynamiques |
| **PKI** | EJBCA + cert-manager | CA nationale souveraine |
| **Event Bus** | Apache Kafka (Strimzi) | Streaming national inter-agences |
| **Databases** | CockroachDB + PostgreSQL | HA distribuée + souveraineté |
| **Search** | OpenSearch | Fork souverain d'Elasticsearch |
| **GitOps** | ArgoCD | Déclaratif, auditable, rollback |
| **Registry** | Harbor | Registry OCI souverain + scanning |
| **Monitoring** | Prometheus + Thanos + Grafana | Métrique long-terme multi-cluster |
| **Logging** | Loki + Promtail | Agrégation logs sans Elasticsearch |
| **Tracing** | Tempo | Tracing distribué OpenTelemetry |
| **Storage** | Rook/Ceph + MinIO | Stockage objet/bloc souverain |
| **Backup** | Velero + pgBackRest | DR automatisé |
| **Security** | Kyverno + Falco + OPA | Admission + runtime + ABAC |
| **Offline** | K3s + NATS JetStream | Edge nodes souverains |

---

## 📁 Structure du Repository

```
SNISID-Core/
├── README.md                          ← Ce document
├── Makefile                           ← Bootstrap + deploy + test
├── PLATFORM_ENGINEERING_FRAMEWORK.md ← Standards souverains
├── DEPLOYMENT_READINESS_REPORT.md    ← Checklist 80+ items
├── SECURITY_AUDIT_REPORT.md          ← Audit CIS benchmarks
│
├── Kubernetes/                        ← Platform Core
│   ├── base/                          Namespaces, NetworkPolicies, ResourceQuotas
│   ├── control-plane/                 RKE2, ETCD, HA config
│   ├── node-pools/                    system, app, data, bio pools
│   ├── storage/                       Rook-Ceph, MinIO
│   ├── networking/                    Cilium, Istio, Ingress
│   ├── autoscaling/                   HPA, VPA, KEDA, CA
│   ├── security/                      Kyverno, Falco, OPA
│   └── overlays/                      prod, staging, dev (Kustomize)
│
├── API-Gateway/                       ← Kong + WAF + Policies
├── Identity/                          ← Identity Core Service
├── EventBus/                          ← Kafka + Schema Registry
├── PKI/                               ← PKI nationale + cert-manager
├── IAM/                               ← Keycloak + Vault + Teleport + OPA
├── Offline-Sync/                      ← K3s Edge + NATS + Delta Sync
├── GitOps/                            ← ArgoCD App-of-Apps
├── DevSecOps/                         ← GitLab CI + Harbor + Cosign
├── Observability/                     ← Prometheus + Grafana + Loki + Tempo
├── Security/                          ← Vault + Falco + Hardening
├── Databases/                         ← PostgreSQL HA + Redis + OpenSearch
├── Testing/                           ← k6 + Chaos Mesh + ZAP + DR
├── Infrastructure/                    ← Terraform + Ansible
└── Diagrams/                          ← C4, K8s, IAM, PKI diagrams
```

---

## 🚀 Bootstrap Rapide

```bash
# 1. Cloner le repo
git clone https://gitlab.snisid.gov.ht/snisid/snisid-core.git
cd snisid-core

# 2. Vérifier les prérequis
make check-prereqs

# 3. Bootstrap cluster RKE2
make bootstrap-cluster ENV=staging

# 4. Déployer l'infrastructure core
make deploy-core ENV=staging

# 5. Déployer les services
make deploy-services ENV=staging

# 6. Vérifier la santé
make health-check ENV=staging
```

---

## 🔒 Principes Zero Trust

1. **Never Trust, Always Verify** — Chaque requête vérifiée même intra-cluster
2. **Least Privilege** — RBAC + ABAC + Network Policies
3. **Assume Breach** — Segmentation micro, alertes runtime
4. **mTLS Everywhere** — Istio Strict mode sur tous les namespaces
5. **Signed Workloads** — Cosign + Kyverno policy enforcement

---

## 📊 SLOs Plateforme

| Service | Availability | P99 Latency | RPO | RTO |
|---------|-------------|-------------|-----|-----|
| API Gateway | 99.99% | < 50ms | 0 | 2 min |
| Identity Service | 99.95% | < 200ms | 0 | 2 min |
| Event Bus (Kafka) | 99.95% | < 10ms | 0 | 5 min |
| PKI/OCSP | 99.9% | < 100ms | 0 | 1h |
| Databases | 99.99% | < 5ms | 0 | 2 min |

---

*SNISID-Core v1.0.0 — République d'Haïti — Mai 2026*
