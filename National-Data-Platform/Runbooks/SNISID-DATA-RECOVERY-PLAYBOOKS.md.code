---
# ============================================================
# SNISID-Data — Disaster Recovery & Breach Containment Playbooks
# Procédures de Sauvetage de la Donnée
# Document ID: SNISID-DATA-RUN-001
# Version: 1.0.0
# ============================================================

## 1. OBJECTIF DES RUNBOOKS

Un "Runbook" est une procédure d'urgence étape par étape. En cas de crise sur le Data Lakehouse (Corruption, Attaque Ransomware, Panne Matérielle), les ingénieurs ne doivent pas improviser, mais suivre la doctrine d'État.

## 2. PLAYBOOK 1 : DATA CORRUPTION RECOVERY (Iceberg Time Travel)

**Scénario :** Un bug dans un script de traitement (Spark) a remplacé toutes les adresses des citoyens par la valeur "NULL" dans la table `Silver_Citoyens`.

**Procédure :**
1. **Ne pas paniquer :** L'architecture Iceberg utilise le concept de "Snapshots" (instantanés).
2. **Identification :** Identifier l'ID du dernier Snapshot valide avant le bug.
3. **Rollback (Time Travel) :** Exécuter la commande Trino : 
   `CALL system.rollback_to_snapshot('Silver_Citoyens', <ID_VALIDE>)`
4. L'intégralité de la base revient à son état exact de la veille, en quelques secondes, sans restaurer de backups massifs depuis des bandes magnétiques.

## 3. PLAYBOOK 2 : DATA BREACH CONTAINMENT (Fuite de Données)

**Scénario :** Le SOC détecte qu'un administrateur système corrompu télécharge massivement le registre d'État Civil.

**Procédure :**
1. **Kill Switch :** Le CISO révoque instantanément l'accès du compte administrateur via le système IAM/Zero Trust (Phase 6).
2. **Key Revocation :** Vault tourne la clé de chiffrement KMS de la base de données. Même si l'administrateur a réussi à copier les fichiers physiques, il ne possède pas la nouvelle clé pour les déchiffrer.
3. **Forensic Audit :** Utilisation des logs de la *Data Audit Fabric* pour déterminer exactement quelles lignes (quels citoyens) ont été exportées avant le blocage, afin d'avertir la population concernée selon la loi (Phase 0).

---
*Document ID: SNISID-DATA-RUN-001 | Approuvé par: Chief Information Security Officer (CISO)*
