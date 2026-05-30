---
# ============================================================
# SNISID-Infra — Edge Nodes Infrastructure
# Infrastructure Régionale Décentralisée (K3s)
# Document ID: SNISID-EDGE-001
# Version: 1.0.0
# ============================================================

## 1. LE CONCEPT DE "L'EDGE COMPUTING" GOUVERNEMENTAL

Au lieu de faire voyager chaque requête de la province (ex: Jérémie) jusqu'à la capitale (Port-au-Prince), des mini-datacenters ("Edge Nodes") sont déployés dans les chefs-lieux départementaux et les commissariats importants.

## 2. ARCHITECTURE D'UN NOEUD EDGE

Un Noeud Edge standard est une "boîte durcie" (Ruggedized Server Appliance) contenant :
- 2x Serveurs physiques x86 basse consommation (NUC Enterprise ou Dell Edge).
- 1x Switch durci avec routage IPsec.
- 1x Batterie UPS intégrée (lithium-ion).
- Unité de stockage NVMe en RAID 1.

### 2.1 Stack Logicielle (K3s + NATS)
- L'infrastructure tourne sur **K3s** (Une version ultra-légère de Kubernetes conçue pour l'Edge).
- Un broker de messagerie ultra-léger **NATS JetStream** agit comme mini-Kafka local.
- Une base de données locale **SQLite / PostgreSQL Edge** (ex: KubeDB) pour le cache.

## 3. GESTION CENTRALISÉE (Fleet Management)

Gérer des centaines de Noeuds Edge manuellement est impossible.
- Le Datacenter Primaire utilise un système de gestion de flotte (ex: Rancher Fleet ou Azure Arc-like).
- Lorsqu'une mise à jour de l'application PNH Mobile est approuvée, elle est "Poussée" (GitOps) de manière asynchrone vers tous les Noeuds Edge. Les noeuds la téléchargent en arrière-plan et redémarrent leurs pods locaux.

---
*Document ID: SNISID-EDGE-001 | Approuvé par: Direction des Opérations Territoriales*
