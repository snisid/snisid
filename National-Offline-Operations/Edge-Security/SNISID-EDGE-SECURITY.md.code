---
# ============================================================
# SNISID-Edge — National Edge Security Framework
# Sécurisation des Equipements Non-Surveillés
# Document ID: SNISID-EDGE-SEC-001
# Version: 1.0.0
# ============================================================

## 1. LE MODÈLE DE MENACE EDGE (Physical Tampering)

Dans un Datacenter (Phase 5), des gardes armés et des caméras protègent les serveurs. Un Edge Node (ex: Commissariat de quartier) peut être physiquement volé ou saboté. L'approche de sécurité est donc radicalement différente.

## 2. SÉCURITÉ MATÉRIELLE (Hardware Trust)

- **Boîtier Anti-Sabotage (Tamper Detection) :** Si le boîtier physique du serveur Edge est ouvert par un tournevis non autorisé, un interrupteur matériel déclenche un court-circuit cryptographique qui efface instantanément les clés en RAM (Zeroization).
- **Secure Boot & TPM :** Le système de fichiers complet est chiffré (LUKS). La clé de déchiffrement est scellée dans la puce TPM 2.0. Si le disque dur est extrait et branché sur un autre ordinateur, il est totalement illisible.

## 3. OFFLINE IAM & RBAC

Comment un policier peut-il s'authentifier si le serveur Keycloak de la capitale est injoignable ?
- L'Edge Node maintient un mini-Keycloak (ou cache LDAP) en lecture seule, synchronisé quotidiennement.
- **Smartcards (PKI) :** Le badge du fonctionnaire contient un certificat cryptographique. Le serveur Edge valide mathématiquement le certificat localement, sans aucun besoin d'Internet.

---
*Document ID: SNISID-EDGE-SEC-001 | Approuvé par: CISO National*
