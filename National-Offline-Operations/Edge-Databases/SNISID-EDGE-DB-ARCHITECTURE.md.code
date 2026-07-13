---
# ============================================================
# SNISID-Edge — National Edge Database Platform
# Stockage Local Chiffré et Cohérence Éventuelle
# Document ID: SNISID-EDGE-DB-001
# Version: 1.0.0
# ============================================================

## 1. LE DÉFI DU STOCKAGE DÉCONNECTÉ

Faire tourner une base de données sur un serveur au milieu de nulle part comporte deux risques :
1. **Risque Physique :** Le serveur peut être volé. La base doit donc être chiffrée (Encryption at Rest) et indéchiffrable sans la clé stockée dans un TPM ou fournie par le KMS central (qui refusera la clé s'il détecte un vol).
2. **Risque Logique (Split Brain) :** Les données locales modifiées vont entrer en conflit avec les données du serveur central.

## 2. CHOIX ARCHITECTURAL (SQLite Distribué / Turso)

Pour les Edge Nodes légers (tablettes, commissariats), une instance CockroachDB complète est trop lourde.
L'architecture utilise **LibSQL/SQLite** avec synchronisation native (type Turso/LiteFS) vers une base de données miroir centralisée.

### 2.1 Delayed Consistency (Cohérence différée)
Les bases Edge sont des "Read-Replicas" partielles.
- Le commissariat de Jacmel ne télécharge *que* la liste des citoyens résidant dans le département du Sud-Est, et la liste nationale des criminels recherchés. Cela réduit la base de 15 millions d'entrées à 500 000.
- Les requêtes d'écriture locales (ex: mise à jour d'adresse) sont asynchrones.

### 2.2 Replication Queues
Chaque Edge DB maintient un journal (Write-Ahead Log) des transactions locales. Dès que le réseau remonte, ce journal est envoyé au **National Sync Engine** (Phase 7, Etape 3) pour intégration dans la base CockroachDB centrale.

---
*Document ID: SNISID-EDGE-DB-001 | Approuvé par: Database Reliability Engineer*
