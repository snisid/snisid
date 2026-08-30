---
# ============================================================
# SNISID-Data — Sovereign Data Security Model
# Chiffrement Avancé et Audit Fabric
# Document ID: SNISID-DATA-SEC-001
# Version: 1.0.0
# ============================================================

## 1. SÉCURITÉ INTRINSÈQUE DES DONNÉES (Data At Rest)

Protéger le réseau (Phase 6) ne suffit pas. Si un disque dur du Datacenter (Phase 5) est volé, les données doivent être indéchiffrables.

## 2. ARCHITECTURE DE CHIFFREMENT

- **Chiffrement au Repos (AES-256) :** Tous les buckets MinIO du Lakehouse sont chiffrés.
- **KMS (Key Management Service) :** Géré par HashiCorp Vault. Rotation automatique des clés toutes les 24 heures.
- **Chiffrement au niveau de la Colonne :** Les champs critiques (Ex: `Template_Biometrique`) sont chiffrés avec une clé distincte de celle du disque dur.

## 3. DATA AUDIT FABRIC

Chaque requête SQL exécutée sur le Lakehouse (via Trino) génère un log immuable :
`[2026-05-25 10:00:00] User: agent_police_09 | Query: SELECT * FROM citizens WHERE nom='Dorsainvil' | Row Count: 5 | Policy: Allowed`

Ces logs d'audit sont envoyés au SOC (Wazuh/OpenSearch) pour détecter les exfiltrations massives ("Pourquoi cet agent fait un SELECT * sur 5 millions de citoyens ?").

---
*Document ID: SNISID-DATA-SEC-001 | Approuvé par: CISO National*
