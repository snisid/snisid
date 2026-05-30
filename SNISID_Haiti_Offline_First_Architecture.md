# SNISID : ARCHITECTURE DE RÉSILIENCE NATIONALE "OFFLINE-FIRST"
## Index Principal et Topologie de Survie Gouvernementale
**République d'Haïti — Infrastructure Souveraine de Continuité de l'État**

---

## PRÉAMBULE
Ce document constitue l'index central de l'architecture de continuité et de survie numérique de la République d'Haïti. Le SNISID est conçu pour garantir l'exercice de la souveraineté étatique (identité, état civil, frontières) y compris lors de scénarios extrêmes : rupture totale des câbles sous-marins, destructions de datacenters par séisme, cyberattaques étatiques massives, ou instabilité civile.

La conception repose sur un paradigme **"Offline-First"** (Priorité au mode déconnecté) couplé à une topologie **Edge Computing** (Informatique en périphérie).

---

## 📑 SOMMAIRE DE L'ARCHITECTURE DE RÉSILIENCE

La documentation est fragmentée en 5 volumes d'ingénierie tactique et stratégique :

### 🌍 [VOLUME 1 : Architecture Edge Nationale](resilience/SNISID_RES_01_National_Edge.md)
Topologie de décentralisation de l'infrastructure étatique :
- Clusters Edge Régionaux (Départements).
- Nœuds de validation déconnectés (Hôpitaux, Tribunaux).
- Hubs de synchronisation intermédiaires.

### 🎒 [VOLUME 2 : Enrôlement Déconnecté (Offline Enrollment)](resilience/SNISID_RES_02_Offline_Enrollment.md)
Opérations de terrain en zones blanches :
- Kits d'enrôlement biométrique mobiles autonomes (Solaire / VSAT).
- Caches biométriques locaux chiffrés (Zero Trust).
- PKI portable et validation citoyenne hors-ligne.

### 🔄 [VOLUME 3 : Moteur de Synchronisation (Sync Engine)](resilience/SNISID_RES_03_Synchronization.md)
Résilience de la donnée asynchrone :
- Protocoles de synchronisation reprenables (Resumable Sync).
- Algorithmes de résolution de conflits (CRDTs).
- Rejeu d'événements (Replay) et journaux d'état immuables.

### 🛡️ [VOLUME 4 : Résilience et Continuité Nationale (Disaster Recovery)](resilience/SNISID_RES_04_National_Resilience_DR.md)
Stratégies de survie des cœurs de réseau :
- Topologie Active-Active inter-datacenters (Port-au-Prince / Cap-Haïtien).
- Systèmes de sauvegarde "Air-Gapped" (Isolés).
- Routage de secours et architecture de bascule régionale (Failover).

### 🚨 [VOLUME 5 : Playbooks Opérationnels de Crise](resilience/SNISID_RES_05_Operational_Playbooks.md)
Procédures d'urgence pour le SOC et le commandement de l'État :
- Workflows d'intervention d'urgence.
- Diagrammes d'états de dégradation (Graceful Degradation).
- Procédures de rétablissement (Post-Incident Recovery).

---

## 🔐 DOCTRINE "SURVIVE AND OPERATE"
L'État ne s'arrête jamais. Si le lien vers la capitale est coupé, le département devient autonome. Si le département est coupé, la commune prend le relais. La sécurité cryptographique (Zero Trust) s'applique jusqu'au terminal mobile alimenté par batterie solaire dans les mornes.

*Spécifié, modélisé et ratifié par la Direction de l'Ingénierie de la Continuité de l'État.*
*Classification : SOUVERAIN / INFRASTRUCTURE CRITIQUE NATIONALE / HAUTE DISPONIBILITÉ*
