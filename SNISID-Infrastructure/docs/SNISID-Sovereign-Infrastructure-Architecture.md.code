# SNISID — Sovereign National Infrastructure Architecture
**Classification:** RESTREINT DEFENSE / SENSIBLE SOUVERAIN  
**Version:** 4.0.0  
**Date:** 2026-05-25  
**Statut:** OFFICIEL — Phase 4 Infrastructure Bootable

---

## 1. Vue d'ensemble stratégique

Le Système National d'Identification et d'Inscription Digitale (SNISID) repose sur une infrastructure physique et cloud-native souveraine, conçue pour fonctionner **24/7/365** en totale autonomie nationale, y compris en situation de crise ou de rupture de connectivité internationale.

### Objectifs souverains
- **Zéro dépendance critique** à un fournisseur cloud étranger
- **Résilience active-active** multi-datacenter
- **Fonctionnement offline** des nœuds périphériques et mobiles
- **Observabilité totale** du système national
- **Sécurité runtime continue** (Zero Trust, mTLS, runtime detection)
- **Reprise d'activité** validée par drills réguliers

---

## 2. Topologie nationale officielle

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        SNISID NATIONAL TOPOLOGY                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   [National Core DC] ─────┬───── [Secondary DR DC]                          │
│        (Capitale)         │        (Region Securisee)                       │
│         ┌────────┐       │        ┌────────┐                              │
│         │K8s Core│◄───────┼───────►│K8s DR  │                              │
│         │K8s IAM │  Sync  │  Sync  │K8s DR  │                              │
│         │K8s Data│◄───────┼───────►│K8s DR  │                              │
│         │K8s Obs │◄───────┼───────►│K8s DR  │                              │
│         └────────┘       │        └────────┘                              │
│              ▲           │           ▲                                   │
│              │    Metro Mesh Link      │                                   │
│              └───────────┬─────────────┘                                   │
│                         │                                                   │
│     ┌───────────────────┼───────────────────┐                               │
│     │           Regional Edge Nodes         │                               │
│     │  ┌────────┐ ┌────────┐ ┌────────┐   │                               │
│     │  │Region 1│ │Region 2│ │Region N│   │  (Départements / Provinces)   │
│     │  │K8s Edge│ │K8s Edge│ │K8s Edge│   │                               │
│     │  └────────┘ └────────┘ └────────┘   │                               │
│     └───────────────────┬───────────────────┘                               │
│                         │                                                   │
│     ┌───────────────────┼───────────────────┐                               │
│     │         Mobile / Offline Nodes        │                               │
│     │  ┌────────┐ ┌────────┐ ┌────────┐   │                               │
│     │  │Mobile 1│ │Offline 1│ │Emergency│  │  (Terrain, zones isolées)   │
│     │  │K8s Edge│ │K8s Edge│ │ K8s Edge│  │                               │
│     │  └────────┘ └────────┘ └────────┘   │                               │
│     └─────────────────────────────────────┘                               │
│                                                                              │
│   [SOC National] ◄───────────────────────────── Monitoring & Alerting       │
│   [PKI Nationale] ◄─────────────────────────── Trust & Identity             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.1 Niveaux d'infrastructure

| Niveau | Description | Rôle | Résilience |
|--------|-------------|------|------------|
| **National Core DC** | Datacenter principal souverain | Workloads critiques nationaux, IAM, données primaires | Active-Active |
| **Secondary DR DC** | Datacenter de secours | Réplication synchrone/asynchrone, basculement automatique | Active-Active |
| **Regional Edge Nodes** | Nœuds régionaux | Services locaux, cache, traitement périphérique | Autonome 72h |
| **Mobile Offline Nodes** | Nœuds terrain mobiles | Enrôlement biométrique terrain, identité en zone isolée | Autonome offline |
| **SOC Infrastructure** | Security Operations Center national | Détection temps réel, réponse incidents, threat intel | Redondé géo |
| **PKI Infrastructure** | Infrastructure de confiance nationale | Certificats souverains, HSM, signatures légales | Air-gapped + online |

### 2.2 Topologie réseau — Metro Mesh National

```
National Backbone (DWDM/MPLS Souverain)
├── Core DC ◄─────────► DR DC          [Lien fibre dédié, <10ms latency]
├── Core DC ◄──┬──┬──► Edge R1..RN    [Réseau admin souverain / VPN MPLS]
├── DR DC  ◄───┴──┴──► Edge R1..RN
└── Edge Ri  ◄───────► Edge Rj        [Inter-régional maillé]
```

---

## 3. Sovereign Government Cloud

### 3.1 Stack souverain

| Domaine | Technologie | Justification souveraine |
|---------|-------------|--------------------------|
| **Compute** | Kubernetes (Kubeadm/RKE2) | Orchestration standard, portable, contrôlable |
| **Virtualization** | Proxmox VE + OpenStack (optionnel) | Hyperviseur on-premise, code ouvert, auditable |
| **Containers** | containerd | Runtime CNCF, pas de dépendance Docker Inc |
| **Registry** | Harbor | Registry souveraine, scanning intégré, signature Cosign |
| **Networking** | Cilium (eBPF) + Calico (policy) | Performance, observabilité réseau, policies L3-L7 |
| **DNS** | CoreDNS + BIND9 externe | Résolution interne souveraine, pas de dépendance DNS public |

### 3.2 Règles de souveraineté

- ❌ Aucun SaaS non souverain en production (pas de AWS/Azure/GCP control plane)
- ❌ Aucune dépendance à une licence propriétaire non révocable
- ❌ Aucun pipeline CI/CD dépendant de services externes non miroirés
- ✅ Tous les composants open-source avec support national ou contrat de maintenance
- ✅ Air-gap possible pour tous les systèmes critiques

---

## 4. National Kubernetes Platform

### 4.1 Architecture multi-cluster

```
National Kubernetes Platform
├── Core Cluster        [Workloads métier SNISID, API publique]
├── Identity Cluster    [IAM, OIDC, FIDO2, biométrie — isolation critique]
├── BPMN Cluster        [Moteurs de workflow, Camunda/Zeebe — isolation processus]
├── Data Cluster        [PostgreSQL, Kafka, Elasticsearch — données personnelles]
├── Cybersecurity Cluster [Falco, Wazuh, SIEM agents — workloads sécurité]
├── Observability Cluster [Prometheus, Loki, Jaeger — isolation métrologie]
└── Edge Clusters       [Regional, Mobile, Offline, Emergency]
```

### 4.2 Spécifications par cluster

| Propriété | Valeur |
|-----------|--------|
| Control plane | 3+ nœuds etcd + 3+ masters (HA stacked) |
| Worker segregation | Control plane tainté, workers par zone/rôle |
| CNI | Cilium (eBPF) + NetworkPolicies strictes |
| CSI | Ceph RBD / CephFS via Rook |
| Ingress | Istio Ingress Gateway (mTLS) + Cilium Gateway |
| Runtime | containerd + gVisor/Kata pour workloads sensibles |

### 4.3 Gouvernance namespace

| Tier | Description | Resource Quotas | Network Policy |
|------|-------------|-----------------|----------------|
| **Tier-0 Souverain** | IAM, PKI, secrets, payment, biométrie | Strictes (CPU/RAM/Pods) | Deny-all explicite |
| **Tier-1 Critique** | API publique, Kafka core, databases | Elevées mais limitées | Restrictif |
| **Tier-2 Métier** | BPMN, services métier, web apps | Standard | Par défaut restrictif |
| **Tier-3 Support** | Observabilité, logging, tooling | Généreuses | Ouvert intra-tier |
| **Tier-4 Edge** | Nœuds périphériques | Limitées par bande passante | Synchro contrôlée |

### 4.4 Autoscaling & résilience

- **Horizontal Pod Autoscaler (HPA)** : metrics CPU/mémoire + custom metrics (file Kafka, latence API)
- **Cluster Autoscaler** : scale-up nœuds workers via Proxmox/OpenStack provider
- **PodDisruptionBudget** : minAvailable=2 pour Tier-0/Tier-1
- **Topology Spread Constraints** : répartition跨 zones +跨 datacenters
- **Virtual Kubelet** : extension vers edge nodes

---

## 5. National GitOps Platform

### 5.1 Chaîne de déploiement souveraine

```
Developer ──► GitLab/GitHub (On-Prem) ──► CI Pipeline (Air-gapped runners)
                                               │
                                               ▼
                                       Harbor (Registry)
                                               │
                                               ▼
                                       ArgoCD (On-Prem) ──► Kyverno/OPA validation
                                               │
                                               ▼
                                    Kubernetes Clusters (GitOps only)
```

### 5.2 Outils & configuration

| Domaine | Technologie | Configuration critique |
|---------|-------------|------------------------|
| GitOps Engine | ArgoCD (HA) | 3 replicas, auto-sync ON, prune ON, selfHeal ON |
| CI/CD | GitLab CI (On-Prem) + GitHub Actions (On-Prem) | Runners air-gapped, pas de sortie internet |
| Registry | Harbor | Trivy scanning, Notary/Cosign sign, replication multi-site |
| Secrets | HashiCorp Vault (auto-unseal Shamir/HSM) | PKI interne, dynamic secrets, encryption transit |
| Policy-as-Code | Kyverno + OPA/Gatekeeper (redondant) | Validation admission obligatoire |
| Packaging | Helm v3 + Kustomize | Charts versionnés, values par environnement |

### 5.3 Règle d'or

> **AUCUN changement manuel en production.**
> Toute modification passe par : commit → CI → scan → validation policy → merge → ArgoCD sync.

---

## 6. National Service Mesh

### 6.1 Istio — Configuration souveraine

| Fonction | Configuration | Détail |
|----------|---------------|--------|
| mTLS | **OBLIGATOIRE** (STRICT) | PeerAuthentication mesh-wide STRICT |
| Authorization | L4 + L7 | AuthorizationPolicy par défaut DENY |
| Ingress | Istio Gateway + SDS | Certificats via cert-manager + Vault PKI |
| Egress | Egress Gateway contrôlé | Whitelist des destinations autorisées |
| Traffic management | Circuit breakers, retries, timeouts | DestinationRule configurée par service |
| Observabilité | Envoy metrics + access logs | EnvoyFilter pour traçage forcé |

### 6.2 Politiques obligatoires

- `PeerAuthentication` : `mtls.mode: STRICT` dans tous les namespaces
- `AuthorizationPolicy` : `action: ALLOW` uniquement sur sources/services explicitement listés
- `DestinationRule` : circuit breaker (maxConnections, maxPendingRequests, outlierDetection)
- `Sidecar` : scope réduit aux services nécessaires (minimise surface d'attaque)

---

## 7. National Storage Platform

### 7.1 Ceph — Cluster distribué national

| Usage | Pool Ceph | Type | Réplication |
|-------|-----------|------|-------------|
| Block volumes K8s | `k8s-rbd` | RBD | 3x cross-rack (Core/DR/Edge) |
| Fichiers partagés | `k8s-fs` | CephFS | 3x cross-rack |
| Object souverain | `national-s3` | RGW | Erasure coding 4+2, multi-site sync |
| Backups immutables | `backup-gold` | RGW | 3x + WORM locking |

### 7.2 Stratégie backup nationale

| Propriété | Implémentation |
|-----------|----------------|
| Chiffrement | AES-256-GCM au repos + TLS 1.3 en transit |
| Immutabilité | Object Lock WORM (Write Once Read Many), 7 ans minimum |
| Air-gapped | Backup quotidien sur bande LTO + vault physique sécurisé |
| Réplication | Ceph multi-site sync Core ↔ DR ; async vers Edge (filtré) |
| Snapshots | Ceph snapshots automatiques toutes les 4h, rétention 30j |
| RPO | < 15 minutes pour données Tier-0 ; < 4h pour Tier-1/2 |
| RTO | < 30 minutes Core/DR ; < 4h Regional ; < 24h Edge mobile |

---

## 8. National Secrets & PKI Platform

### 8.1 HashiCorp Vault — Architecture HSM

```
Vault Raft Cluster (5 nœuds)
├── auto-unseal via HSM Thales Luna 7 (PKCS#11)
├── Shamir seal backup (6 parts, 4 threshold) — coffre physique
├── PKI Secret Engine (Root CA nationale SNISID + Intermediate CA)
├── Transit Engine (chiffrement données Tier-0)
├── KV v2 (secrets applicatifs, rotation automatique)
├── Database Engine (credentials dynamiques PostgreSQL)
└── SSH Engine (credentials dynamiques bastions)
```

### 8.2 cert-manager + CA nationale

| Rôle | Configuration |
|------|---------------|
| Émetteur interne | `ClusterIssuer` : Vault PKI (intermediate CA) |
| Émetteur externe | `ClusterIssuer` : ACME interne (pas Let's Encrypt — souveraineté) |
| Rotation | Certificats 90j, renouvellement auto à 30j |
| HSM | Clés privées root/intermediate stockées HSM, jamais exportées |

### 8.3 Règles de gestion des secrets

- ❌ **JAMAIS** de secret hardcodé dans un manifeste, une image, ou un repo Git
- ✅ Tous les secrets injectés via Vault Agent Injector ou External Secrets Operator
- ✅ Rotation automatique : credentials DB toutes les 24h, certificats auto-renouvelés
- ✅ Audit log Vault : forwarding temps réel vers SIEM national

---

## 9. Network Security Model — Zero Trust National

### 9.1 Segmentation obligatoire

| Zone | Description | Contrôle |
|------|-------------|----------|
| **Internet DMZ** | Reverse proxy, WAF | Ingress contrôlé uniquement, pas d'egress libre |
| **Management** | Bastions, VPN admin, jump hosts | Accès MFA + certificate-based, audité |
| **Tier-0** | IAM, PKI, HSM | Aucun accès direct, uniquement via service mesh mTLS |
| **Tier-1** | APIs, Kafka, DB core | Accès intra-mesh uniquement, egress filtré |
| **Tier-2** | Apps métier, BPMN | Accès mesh + API Gateway |
| **Observability** | Métriques, logs, traces | Read-only depuis tiers inférieurs, push vers tiers supérieurs |
| **Edge** | Nœuds régionaux/mobiles | VPN site-to-site ou WireGuard, policies restrictives |

### 9.2 East-West / North-South

- **North-South** : Istio Ingress Gateway + Cilium L7 policies + WAF nationaux
- **East-West** : mTLS Istio + Cilium NetworkPolicies L3-L4 + AuthorizationPolicies L7
- **Règle** : AUCUNE communication inter-service sans policy explicite (default deny)

### 9.3 Firewall policies

| Direction | Règle |
|-----------|-------|
| Intra-Tier-0 | Deny all, allow explicit par service ID |
| Tier-1 → Tier-0 | Allow par endpoints API publiés uniquement |
| Tier-2 → Tier-1 | Allow par service account identifié |
| Edge → Core | Allow sync contrôlé (Kafka mirror, API REST) |
| Egress global | Deny all, allow par Egress Gateway whitelist |

---

## 10. National Observability Stack

### 10.1 Stack — Observability Cluster dédié

| Signal | Technologie | Stockage | Rétention |
|--------|-------------|----------|-----------|
| **Metrics** | Prometheus (HA, Thanos/Cortex) | Ceph/Cassandra | 2 ans |
| **Logs** | Loki (HA, multi-tenant) | S3 Ceph | 1 an chaud, 7 ans froid |
| **Traces** | Jaeger + OpenTelemetry Collector | Elasticsearch/Ceph | 90 jours |
| **Dashboards** | Grafana (HA, auth OIDC/Vault) | — | — |
| **Alerting** | Alertmanager (HA) + PagerDuty souverain | — | — |
| **Profilage** | Pyroscope/Parca (optionnel) | Ceph | 30 jours |

### 10.2 Coverage obligatoire

| Système | Métriques | Logs | Traces | Alertes |
|---------|-----------|------|--------|---------|
| Kubernetes | ✅ node_exporter, kube-state-metrics | ✅ Fluent Bit | ✅ OpenTelemetry | ✅ critères SLO/SLA |
| APIs SNISID | ✅ latence, taux erreur, débit | ✅ access logs | ✅ Istio spans | ✅ p95 > 200ms, 5xx > 0.1% |
| Kafka | ✅ brokers, consumers, lag | ✅ logs brokers | ✅ — | ✅ consumer lag, offline partitions |
| IAM/Vault | ✅ auth rates, seal status | ✅ audit logs | ✅ — | ✅ seal, fail auth |
| BPMN | ✅ job counts, durations | ✅ workflow logs | ✅ Zeebe traces | ✅ incident BPMN |
| Infrastructure | ✅ CPU/ram/disk/réseau | ✅ syslog | ✅ — | ✅ disk >80%, load |

### 10.3 Règle d'or

> **Tout doit être observable.** Si un composant n'est pas monitoré, il n'est pas en production.

---

## 11. Runtime Security Platform

### 11.1 Détection & prévention

| Couche | Outil | Fonction |
|--------|-------|----------|
| **Runtime detection** | Falco (eBPF) | Détection comportements anormaux noyau/conteneurs |
| **Container scanning** | Trivy (Harbor CI) + Trivy Operator (cluster) | CVE, secrets dans images, config issues |
| **Admission control** | Kyverno + Gatekeeper | Validation/refus des manifests non conformes |
| **Supply chain** | Cosign (signature images) + Syft (SBOM) | Images signées, SBOM stockés, pas d'image non signée |
| **Network runtime** | Cilium Tetragon | Appareillage temps réel processus/réseau |

### 11.2 Règles Falco nationales

- Execution shell dans un conteneur (sauf images debug autorisées)
- Privilege escalation détectée (sudo, setuid)
- Accès fichiers sensibles (/etc/shadow, HSM paths)
- Connexion réseau sortante non attendue depuis Tier-0
- Modification binaire en cours d'exécution
- Contact vers IPs non whitelistées (C2 potentiel)

### 11.3 Règle d'or

> **La sécurité est continue.** Pas de scan ponctuel. Détection temps réel, réponse automatique (isolate pod → alert SOC → forensics).

---

## 12. Disaster Recovery Strategy

### 12.1 Topologie active-active

```
Core DC (Actif)         DR DC (Actif)
     │                       │
     ├── Ceph RBD sync ◄────┤  [Sync réplication block synchronisée]
     ├── Kafka MirrorMaker2─┤  [Async topics critiques]
     ├── PostgreSQL Patroni─┤  [Sync replication quorum]
     ├── Vault Raft ◄───────┤  [Cluster étendu 5 nœuds : 3 Core + 2 DR]
     └── ArgoCD ◄───────────┤  [Même repo Git, sync simultanée]
```

### 12.2 Objectifs

| Indicateur | Objectif | Méthode |
|------------|----------|---------|
| **RPO** | ≤ 15 min (Tier-0), ≤ 4h (Tier-1/2) | Sync Ceph, Patroni sync, Kafka MM2 |
| **RTO** | ≤ 30 min (Core/DR), ≤ 4h (Regional) | DNS failover, Istio locality lb, runbooks |
| **Backup** | Multi-site + air-gapped | Ceph RGW + LTO bande + vault physique |
| **Drills** | Mensuel (tabletop), trimestriel (full restore) | Runbooks validés, post-mortem systématique |

### 12.3 Basculement automatique

| Scénario | Déclencheur | Action |
|----------|-------------|--------|
| Panne Core DC | Health check failed x3 | DNS failover → DR DC ; ArgoCD continue sur DR |
| Corruption données | Ceph checksum failed | Restauration snapshot Ceph + replay Kafka |
| Ransomware | Chiffrement anormal détecté | Isolation réseau (Cilium policy emergency) + restore bande |
| Panne régionale | Edge offline > 72h | Activation Emergency Nodes, sync différée |

---

## 13. Edge Infrastructure

### 13.1 Types de nœuds edge

| Type | Fonction | Connectivité | Autonomie |
|------|----------|--------------|-----------|
| **Regional nodes** | Départements/Provinces | Fibre / MPLS admin | 72h cache + DB locale |
| **Mobile nodes** | Unités mobiles terrain | 4G/5G souverain / SATCOM | 24h opérationnel |
| **Offline nodes** | Zones isolées (rural, montagne) | Sync périodique (USB/SD sécurisés) | 7j offline complet |
| **Emergency nodes** | Catastrophes naturelles, crises | Satellite / radio / mesh ad-hoc | 72h sur batteries/solaire |

### 13.2 Architecture edge K8s

- **K3s** (lightweight K8s) sur matériel ruggedisé
- **etcd externe** remplacé par SQLite (K3s single-node) ou etcd 3-nœuds (regional)
- **Local storage** : SSD NVMe + Ceph edge (mini-cluster 3 nœuds si regional)
- **Sync Core** : Kafka MirrorMaker2 edge → core (async) ; API REST batch (offline)
- **Vault Edge** : Vault Agent en mode cache (token limité, pas de master key)

### 13.3 Règle d'or

> **Les edge nodes doivent fonctionner offline.** L'identité nationale ne s'arrête pas quand le réseau tombe.

---

## 14. Runbooks infrastructure (SOPs)

| Runbook | Fréquence test | Temps cible |
|---------|----------------|-------------|
| Kubernetes cluster recovery | Trimestriel | < 45 min |
| Kafka recovery (perte quorum) | Trimestriel | < 30 min |
| Vault recovery (seal/raft) | Mensuel | < 15 min |
| DR failover (Core → DR) | Mensuel | < 30 min |
| Certificate rotation emergency | Semestriel | < 2h |
| Ceph recovery (perte OSD) | Trimestriel | < 1h |
| Edge node provisioning | À chaque déploiement | < 30 min |
| SOC incident response (isolé pod) | Continu (automatisé) | < 5 min |

---

## 15. Infrastructure Standards

| Domaine | Standard |
|---------|----------|
| **Kubernetes** | Version N-2 max ; CNIs/CSIs compatibles souverains ; pas de CRD alpha en prod |
| **Naming** | `{region}-{env}-{tier}-{service}-{ordinal}` (ex: `core-prod-t0-vault-01`) |
| **GitOps** | Mono-repo `infrastructure-national` ; branche `main` protégée ; PR obligatoire + 2 reviews |
| **Terraform** | Modules versionnés ; state backend Vault/Consul ; `terraform plan` obligatoire en CI |
| **Helm** | Charts lintés (ct lint) ; values chiffrées (SOPS) ; pas de defaults dangereux |
| **Security** | CIS Kubernetes Benchmark v1.8+ ; SOC 2 Type II mappé ; politiques Kyverno imposées |
| **Logging** | Format JSON structuré (ECS) ; timestamp UTC ; champ `tenant=national` obligatoire |

---

## 16. Repository centralisé

```
Infrastructure/
├── Kubernetes/
│   ├── clusters/{core,identity,bpmn,data,cyber,obs,edge}
│   ├── namespaces/{tier-0,tier-1,tier-2,tier-3,tier-4}
│   ├── policies/kyverno/
│   └── network-policies/cilium/
├── Terraform/
│   ├── modules/{proxmox,ceph,cilium,vault,observability}
│   └── environments/{core,dr,edge-regional,edge-mobile}
├── GitOps/
│   ├── argocd/apps/{core,identity,bpmn,data,cyber,obs,edge}
│   └── argocd/projects/
├── Helm/
│   └── charts/{snisid-api,snisid-bpmn,vault-raft,ceph-cluster,istio-base}
├── Networking/
│   ├── cilium/{cni-policies,egress-gateways,l7-rules}
│   ├── coredns/{zones,forwarders}
│   └── firewall/{nftables,iptables,cloud-firewall}
├── Storage/
│   └── ceph/{cluster, pools, rgw-policies, backup-jobs}
├── Observability/
│   ├── prometheus/{rules,scrapes,alerts}
│   ├── grafana/{dashboards, datasources}
│   ├── loki/{pipelines,retention}
│   └── jaeger/{sampling,collectors}
├── Security/
│   ├── falco/{rules,custom-macros,outputs}
│   ├── kyverno/{policies,exceptions}
│   ├── vault/{policies,auth-methods,pki-roles}
│   └── cert-manager/{issuers,certificates}
├── DR/
│   └── failover/{dns-failover,argocd-apps,ceph-mirror,kafka-mm2}
├── Edge/
│   ├── regional/{k3s-config,manifests,sync-scripts}
│   ├── mobile/{k3s-config,manifests,sync-scripts}
│   └── offline/{k3s-config,manifests,airgap-bundle}
├── Runbooks/
│   ├── cluster-recovery.md
│   ├── kafka-recovery.md
│   ├── vault-recovery.md
│   ├── dr-failover.md
│   ├── certificate-rotation.md
│   └── edge-provisioning.md
└── Standards/
    └── infrastructure-standards.md
```

---

## 17. Validation de la Phase 4

| Élément | Statut | Preuve |
|---------|--------|--------|
| Kubernetes National Platform | ✅ | Manifests clusters, kubeadm configs, join tokens sécurisés |
| Sovereign Government Cloud | ✅ | Specs Proxmox, plans réseau, topologie Ceph |
| Multi-Cluster Architecture | ✅ | Définitions 7 clusters, federation, Istio multi-primary |
| GitOps Platform | ✅ | ArgoCD apps, repo structure, CI pipelines, Vault integration |
| Infrastructure as Code | ✅ | Terraform modules, Helm charts, state backends |
| Service Mesh | ✅ | Istio manifests, PeerAuthentication STRICT, policies mTLS |
| National Observability Stack | ✅ | Prometheus rules, Grafana dashboards, Loki pipelines |
| Storage Platform | ✅ | Ceph cluster CR, pools, RGW multi-site, backup jobs |
| Runtime Security | ✅ | Falco rules, Kyverno policies, Trivy scans, Cosign workflows |
| Production Infrastructure | ✅ | Runbooks, SOPs, DR plans, edge nodes, PKI HSM |

---

**Document approuvé pour mise en production nationale.**  
*Classification: RESTREINT DEFENSE — SNISID Infrastructure Nationale*
