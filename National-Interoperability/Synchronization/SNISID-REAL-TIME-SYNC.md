---
# ============================================================
# SNISID-Interop — National Real-Time Synchronization Engine
# Évènements, Cohérence Éventuelle et Mode Offline
# Document ID: SNISID-SYNC-ENGINE-001
# Version: 1.0.0
# ============================================================

## 1. EVENTUAL CONSISTENCY (Cohérence Éventuelle)

Avec l'architecture distribuée de l'État, il est impossible de garantir que toutes les bases de données (SNISID, DGI, MSPP) soient parfaitement synchronisées à la milliseconde près. Le modèle adopté est "l'Eventual Consistency".
- Si l'adresse d'un citoyen change à 10:00:00.
- La DGI la recevra via Kafka à 10:00:01. 
- S'il y a une panne réseau à la DGI, le message attendra dans Kafka et sera consommé à 14:00:00 lors du retour réseau.

## 2. OFFLINE SYNCHRONIZATION (Reconciliation)

Le modèle s'appuie sur la logique développée en Phase 3 pour les kits PNH (K3s Edge + NATS).

### 2.1 Conflits de Données (Conflict Resolution)
Si deux agences mettent à jour le même dossier hors ligne, la résolution s'opère selon des règles strictes définies dans le *Data Governance Model* :
- Règle 1 : L'agence *System of Record* (Propriétaire) gagne toujours.
- Règle 2 : À droits égaux, la version avec le Timestamp d'événement le plus récent gagne (Last-Write-Wins).

### 2.2 Retry Orchestration
Le *Sync Engine* (un service Go) gère les échecs de livraison HTTP (webhook) :
- Retry 1 : 10 secondes.
- Retry 2 : 1 minute.
- Retry 3 : 5 minutes.
- Retry N : Backoff Exponentiel.
- Après 24h : Message envoyé en Dead-Letter Queue (DLQ) pour intervention manuelle.

---
*Document ID: SNISID-SYNC-ENGINE-001 | Approuvé par: Data Architecture Board*
