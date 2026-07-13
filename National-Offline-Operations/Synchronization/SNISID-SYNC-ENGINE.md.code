---
# ============================================================
# SNISID-Edge — National Offline Synchronization Engine
# Bidirectional Sync et Résolution des Conflits
# Document ID: SNISID-SYNC-001
# Version: 1.0.0
# ============================================================

## 1. LE MOTEUR DE SYNCHRONISATION (Heartbeat de l'État)

Le Sync Engine est responsable de la cohérence de la donnée (Eventual Consistency) entre les milliers d'appareils de terrain et le Datacenter National.

## 2. MODÈLE DE RÉSOLUTION DE CONFLITS (CRDTs & Vector Clocks)

Que se passe-t-il si un citoyen fait une demande de passeport au bureau central de l'Immigration (Connecté), et que *simultanément*, un policier à Ouanaminthe (Déconnecté) met à jour son adresse sur une tablette ?

### 2.1 Approche par Horloge Logique (Vector Clocks)
Le système n'utilise pas l'heure du système (car l'horloge d'une tablette peut être fausse). Il utilise des compteurs d'état.
- Si le changement A (Port-au-Prince) et le changement B (Ouanaminthe) surviennent sur le même objet (le profil du citoyen), le Moteur de Sync fusionnera les champs non-conflictuels (ex: le passeport d'un côté, l'adresse de l'autre).
- S'ils modifient le *même* champ (Conflit dur), le modèle applique la règle gouvernementale **"Central Always Wins"** (Le Datacenter Primaire fait autorité), mais génère une alerte d'audit "Conflit Résolu" pour l'administrateur.

## 3. DELTA SYNCHRONIZATION (Optimisation Bande Passante)

La connexion satellite (VSAT) coûte cher et offre peu de bande passante.
Lors de la reconnexion d'un Edge Node après 3 jours de tempête, il ne télécharge pas les 15 millions d'identités. Il ne télécharge que le "Delta" (les modifications exactes qui ont eu lieu pendant ces 3 jours), compressé au format binaire (ex: Protocole gRPC).

---
*Document ID: SNISID-SYNC-001 | Approuvé par: Data Governance Board*
