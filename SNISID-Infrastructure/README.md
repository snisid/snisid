# SNISID — Infrastructure Nationale Souveraine
**Système National d'Identification et d'Inscription Digitale**

**Classification:** RESTREINT DEFENSE  
**Statut:** Phase 4 — Infrastructure Bootable & Production Ready  
**Date:** 2026-05-25

---

## Objectif

Ce repository centralise l'intégralité de l'infrastructure physique et cloud-native souveraine qui fait tourner le SNISID 24/7/365.

**Règle absolue:**
- Survivre à une panne nationale
- Continuer offline (edge & emergency)
- Être observable en temps réel
- Être sécurisé (Zero Trust, mTLS, runtime security)
- Être scalable et bootable

---

## Structure du Repository

```
Infrastructure/
├── docs/                           # Architecture officielle nationale
│   └── SNISID-Sovereign-Infrastructure-Architecture.md
├── kubernetes/
│   ├── clusters/                   # kubeadm configs par cluster (Core, Identity, Data...)
│   ├── namespaces/                 # Définitions namespaces + Pod Security Standards
│   ├── policies/                   # Kyverno + OPA/Gatekeeper policies
│   └── network-policies/           # CiliumNetworkPolicies
├── terraform/
│   ├── modules/                    # Modules réutilisables (Proxmox, Ceph, Vault...)
│   ├── environments/               # Déploiements par DC (core, dr, edge...)
│   └── providers/                  # Configs providers souverains
├── gitops/
│   ├── argocd/                     # Applications ArgoCD nationales
│   └── apps/                       # App manifests pour tous les clusters
├── helm/
│   └── charts/                     # Charts SNISID officiels
├── networking/
│   ├── cilium/                     # CNI eBPF + policies L3-L7 + mesh
│   ├── coredns/                    # DNS souverain national
│   └── firewall/                   # Règles nftables/paloalto
├── storage/
│   └── ceph/                       # Rook Ceph — block, filesystem, object multi-site
├── observability/
│   ├── prometheus/                 # Rules, ServiceMonitors, Alertmanager
│   ├── grafana/                    # Dashboards nationaux
│   ├── loki/                       # Log aggregation pipelines
│   └── jaeger/                     # Distributed tracing OpenTelemetry
├── security/
│   ├── falco/                      # Runtime detection rules
│   ├── kyverno/                    # Admission control policies
│   ├── vault/                      # HA Raft + HSM configuration
│   └── cert-manager/               # PKI nationale + issuers
├── dr/
│   └── failover/                   # Automatisation basculement Core ↔ DR
├── edge/
│   ├── regional/                   # K3s edge départements
│   ├── mobile/                     # Nodes terrain mobiles
│   ├── offline/                    # Air-gap bundles zones isolées
│   └── emergency/                  # Nœuds catastrophe survie
├── runbooks/                       # SOPs infrastructure nationaux
│   ├── cluster-recovery.md
│   ├── kafka-recovery.md
│   ├── vault-recovery.md
│   ├── dr-failover.md
│   ├── certificate-rotation.md
│   └── edge-provisioning.md
└── standards/
    └── infrastructure-standards.md  # Standards obligatoires nationaux
```

---

## Topologie Nationale

| Niveau | Description | Technologie | Autonomie |
|--------|-------------|-------------|-----------|
| **National Core DC** | Datacenter principal | Proxmox + K8s HA | Active-Active avec DR |
| **Secondary DR DC** | Datacenter secours | Proxmox + K8s HA | Active-Active avec Core |
| **Regional Edge** | Départements/Provinces | K3s + Cilium | 72h offline |
| **Mobile Nodes** | Terrain mobile | K3s single-node | 24h opérationnel |
| **Offline Nodes** | Zones sans réseau | K3s air-gap bundle | 7j offline complet |
| **Emergency Nodes** | Catastrophes | K3s + solaire + satellite | 72h survie totale |
| **SOC National** | Security Operations | Wazuh + SIEM intégré | Redondé géo |
| **PKI Nationale** | Confiance souveraine | Vault + HSM Thales + cert-manager | Air-gapped + online |

---

## Clusters Kubernetes

| Cluster | Rôle | Tier | Nodes Masters | Nodes Workers | Isolation |
|---------|------|------|---------------|---------------|-----------|
| **Core** | API publique SNISID | Tier-1 | 3 | 5 | NetworkPolicies + Istio |
| **Identity** | IAM, Vault, PKI | Tier-0 | 3 | 3 | VLAN isolé + mTLS STRICT |
| **BPMN** | Workflows nationaux | Tier-1 | 3 | 4 | Namespace isolé |
| **Data** | PostgreSQL, Kafka, Search | Tier-1 | 3 | 6 | High I/O, Ceph dédié |
| **Cyber** | SOC, Falco, Wazuh | Tier-1 | 3 | 3 | Ingress limité SIEM |
| **Observability** | Prometheus, Grafana, Loki, Jaeger | Tier-3 | 3 | 4 | Read-only cross-cluster |
| **Edge Regional** | Services locaux départements | Tier-4 | 1-3 | 3+ | K3s lightweight |
| **Edge Mobile** | Terrain mobile | Tier-4 | 1 | 1-2 | Air-gap capable |

---

## Stack Technologique Souverain

| Domaine | Technologie | Justification |
|---------|-------------|---------------|
| Compute | Kubernetes (kubeadm/RKE2/K3s) | Standard CNCF, portable |
| Virtualization | Proxmox VE | Hyperviseur on-premise, auditable |
| Containers | containerd | Runtime CNCF, pas de dépendance Docker |
| Registry | Harbor | Scanning CVE, Cosign signatures, replication |
| CNI | Cilium (eBPF) | Performance, observabilité réseau, encryption WireGuard |
| DNS | CoreDNS + BIND9 | Résolution interne souveraine |
| GitOps | ArgoCD | Déploiement 100% GitOps, auto-sync |
| CI/CD | GitLab CI / GitHub Actions On-Prem | Runners air-gapped |
| Secrets | HashiCorp Vault + HSM Thales | Auto-unseal PKCS#11, jamais hardcodé |
| Certificates | cert-manager + Vault PKI | Rotation auto, CA nationale |
| Service Mesh | Istio | mTLS STRICT mesh-wide, circuit breakers |
| Policies | Kyverno + OPA/Gatekeeper | Admission control, standards nationaux |
| Storage Block | Ceph RBD (Rook) | Réplication cross-rack, multi-site |
| Storage File | CephFS (Rook) | Partagé national |
| Storage Object | Ceph RGW (Rook) | S3 compatible, immutable backups, multi-site |
| Metrics | Prometheus + Thanos/Cortex | HA, long-term retention 2 ans |
| Logs | Loki | Multi-tenant, pipelines structurés |
| Traces | Jaeger + OpenTelemetry | Distributed tracing national |
| Dashboards | Grafana + Alertmanager | National dashboard, alerting SOC |
| Runtime Security | Falco + Tetragon | Détection eBPF temps réel |
| Scanning | Trivy (Harbor + Operator) | CVE + secrets scanning continu |
| Supply Chain | Cosign + Syft | Signatures images, SBOM obligatoires |

---

## Règles Absolues (Zero Exception)

1. **Aucun changement manuel en production.** Tout passe par Git → CI → ArgoCD.
2. **Jamais de secrets hardcodés.** Tous les secrets injectés via Vault Agent / External Secrets Operator.
3. **mTLS STRICT sur tout le mesh.** Aucune communication inter-service sans TLS mutuel.
4. **Default DENY réseau.** Cilium NetworkPolicies + Istio AuthorizationPolicy deny-all par défaut.
5. **Tout doit être observable.** Si ce n'est pas monitoré, ce n'est pas en production.
6. **Images signées + digest.** Pas de `latest`. Pas d'image non scannée.
7. **DR testé mensuellement.** RTO < 30 min, RPO < 15 min (Tier-0).
8. **Edge fonctionne offline.** L'identité nationale ne s'arrête pas quand le réseau tombe.

---

## Démarrage Rapide (Bootstrap National)

### 1. Infrastructure Physique
```bash
cd terraform/environments/core
terraform init
terraform plan
terraform apply
```

### 2. Bootstrap Kubernetes (premier master)
```bash
# Sur core-prod-t1-master-01
sudo kubeadm init --config=kubernetes/clusters/core/kubeadm-config.yaml
```

### 3. CNI + Service Mesh
```bash
kubectl apply -f networking/cilium/cilium-values.yaml
kubectl apply -f networking/cilium/cilium-clusterwide-policies.yaml
kubectl apply -f networking/cilium/istio-gateway.yaml
```

### 4. Stockage National
```bash
kubectl apply -f storage/ceph/rook-cluster.yaml
```

### 5. GitOps
```bash
kubectl apply -f gitops/argocd/projects/national-app-project.yaml
kubectl apply -f gitops/argocd/apps/
```

### 6. Sécurité
```bash
kubectl apply -f security/kyverno/
kubectl apply -f security/falco/snisid-custom-rules.yaml
kubectl apply -f security/vault/vault-ha-values.yaml
kubectl apply -f security/cert-manager/clusterissuers.yaml
```

### 7. Observabilité
```bash
kubectl apply -f observability/prometheus/prometheus-rules.yaml
kubectl apply -f observability/grafana/dashboards/
```

---

## Contacts & Gouvernance

| Rôle | Entité | Responsabilité |
|------|--------|--------------|
| IGC | Inspection Générale Cyber | Validation architecture, PKI, exceptions |
| SOC National | Security Operations Center | Monitoring 24/7, réponse incidents |
| Équipe Infra Nationale | Direction Technique | Ops, runbooks, provisioning |
| Équipe Edge | Terrain National | Mobile, offline, emergency nodes |
| Équipe Observabilité | Métrologie Nationale | Prometheus, Loki, Jaeger, SLOs |

---

*SNISID — Infrastructure souveraine bootable 24/7. Classification: RESTREINT DEFENSE.*
