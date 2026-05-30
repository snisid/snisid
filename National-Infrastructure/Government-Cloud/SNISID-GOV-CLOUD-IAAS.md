---
# ============================================================
# SNISID-Infra — Government Sovereign Cloud (IaaS/PaaS)
# Cloud privé, Isolation Multi-Tenant (OpenStack / Kubernetes)
# Document ID: SNISID-GOV-CLOUD-001
# Version: 1.0.0
# ============================================================

## 1. VISION: THE SOVEREIGN CLOUD

L'État ne peut pas confier les données d'état civil, de justice et de police à un fournisseur Cloud public soumis au Cloud Act américain ou à d'autres juridictions étrangères. Le SNISID s'appuie sur un **Cloud Privé Souverain** hébergé dans les Datacenters de la Phase 5.

## 2. STACK TECHNOLOGIQUE (IaaS & PaaS)

| Couche | Technologie Standard | Rôle |
|--------|----------------------|------|
| **Compute (VMs)** | OpenStack (KVM) | Fournir des Machines Virtuelles aux Ministères legacy. |
| **Containers (PaaS)** | Kubernetes (RKE2) | Fournir des clusters sécurisés pour les applications modernes. |
| **Storage (Block/Obj)**| Ceph & MinIO | Fournir des disques persistants et du stockage S3-compatible. |
| **Networking (SDN)** | Cilium / OVN | Microsegmentation, eBPF routing. |

## 3. ISOLATION MULTI-TENANT (Siloing)

Même si le Ministère de la Santé (MSPP) et le SNISID partagent les mêmes serveurs physiques, ils sont strictement isolés (Multitenancy).

- **Network Isolation :** Les réseaux virtuels (VPC) du MSPP ne peuvent pas router vers les réseaux SNISID, sauf via l'API Gateway officielle (Phase 4).
- **Compute Isolation :** Les Pods Kubernetes de la DCPJ ne seront jamais schedulés sur le même serveur physique (Bare-Metal Node) que les Pods publics du portail citoyen (utilisation de `nodeAffinity` et `taints`).
- **Storage Isolation :** Les volumes chiffrés utilisent des clés KMS distinctes gérées dans HashiCorp Vault pour chaque ministère.

## 4. KUBERNETES-AS-A-SERVICE (KaaS)

Plutôt que d'avoir un "Cluster Géant", le Cloud Gouvernemental déploie des clusters dédiés à la volée via **Cluster API (CAPI)**.
Le Ministère des Finances (DGI) demande un cluster de 10 noeuds via Terraform -> L'infrastructure provisionne les VMs OpenStack et boot le cluster RKE2 automatiquement en 15 minutes.

---
*Document ID: SNISID-GOV-CLOUD-001 | Approuvé par: Direction de l'Infrastructure Numérique*
