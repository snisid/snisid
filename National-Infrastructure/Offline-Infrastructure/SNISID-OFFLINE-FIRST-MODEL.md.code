---
# ============================================================
# SNISID-Infra — National Offline-First Infrastructure Model
# Opérations Déconnectées et Autonomie Régionale
# Document ID: SNISID-OFFLINE-001
# Version: 1.0.0
# ============================================================

## 1. PRINCIPE D'AUTONOMIE RÉGIONALE (SURVIVABILITY)

Si un ouragan coupe la fibre optique et détruit la parabole VSAT d'un département, l'État local (Mairie, Police, Hôpital) doit continuer de fonctionner. C'est le principe de "Disconnected Operations".

## 2. MODE DÉGRADÉ (OFFLINE)

Lorsqu'un Edge Node perd sa connexion au Datacenter Primaire :

1. **Local Authentication :** L'instance locale (mini-Keycloak ou cache LDAP) permet aux policiers et médecins locaux de s'authentifier.
2. **Local Read Cache :** Les vérifications d'identité se font sur le cache synchronisé la nuit précédente (ex: Liste des personnes recherchées du département).
3. **Local Write Spooling :** Les créations de données (ex: Enregistrement d'une naissance ou Arrestation) sont écrites sur la base locale et mises en attente dans la file NATS JetStream (Spooling).

## 3. SYNCHRONIZATION RECOVERY (Retour à la normale)

Dès que le lien WAN remonte (même par intermittence, ex: 3G instable) :

1. L'Edge Node initie un tunnel mTLS avec le Datacenter.
2. Le Spooler NATS local "vide" sa file d'attente vers le Kafka central.
3. Le Moteur de Synchronisation (Phase 4) traite les événements par ordre chronologique.
4. L'Edge Node télécharge (Pull) les mises à jour nationales (nouveaux mandats d'arrêts, etc.) sous forme de patch binaire compressé pour économiser la bande passante.

---
*Document ID: SNISID-OFFLINE-001 | Approuvé par: Architecte Souverain*
