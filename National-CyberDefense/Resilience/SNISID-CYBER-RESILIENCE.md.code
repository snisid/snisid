---
# ============================================================
# SNISID-Cyber — National Cyber Resilience Framework
# Continuité Cyber, Opérations Hors-Ligne (Offline)
# Document ID: SNISID-RESIL-001
# Version: 1.0.0
# ============================================================

## 1. CYBER SURVIVABILITY

La résilience cybernétique signifie que le gouvernement doit continuer à fonctionner *même pendant* une cyberattaque en cours.

## 2. DÉFENSE HORS-LIGNE (Offline Cyber Operations)

La topologie "Offline-First" (Phase 5) pose des défis uniques. Si l'Edge Node "Jacmel" est déconnecté du réseau WAN, le SOC Central (Port-au-Prince) est "aveugle" aux attaques s'y déroulant.

### 2.1 Edge SOC (K3s Local Defense)
Chaque Edge Node embarque une version "Edge" du SIEM et de l'EDR.
- L'agent Wazuh local envoie ses logs au serveur Wazuh Edge local.
- Les règles SOAR locales sont capables d'isoler automatiquement un poste infecté à Jacmel, **sans avoir besoin d'attendre la connexion à Port-au-Prince**.

### 2.2 Delayed Sync (Forensic Spooling)
Lorsque la connexion VSAT/Fibre est rétablie, les logs d'attaques mis en cache à Jacmel sont envoyés prioritairement (QoS) vers le SIEM Central pour analyse a posteriori par les analystes L3.

## 3. IMMUTABLE RECOVERY (Restauration Sécurisée)

En cas de chiffrement complet (Ransomware), l'équipe Cyber ne restaure pas les données "à l'aveugle".
- Les backups (S3 MinIO WORM) sont d'abord montés dans un environnement isolé (Sandbox Kubernetes).
- Un scan antiviral et de compromission complet est exécuté sur le backup pour s'assurer que l'attaquant n'y a pas caché de porte dérobée (Backdoor) avant de chiffrer.
- Seulement après validation, la donnée est réinjectée en production.

---
*Document ID: SNISID-RESIL-001 | Approuvé par: Architecte Souverain*
