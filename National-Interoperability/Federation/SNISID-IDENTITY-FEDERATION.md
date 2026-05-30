---
# ============================================================
# SNISID-Interop — National Identity Federation Platform
# SSO, OIDC, et SAML inter-agences
# Document ID: SNISID-FEDERATION-001
# Version: 1.0.0
# ============================================================

## 1. CONCEPT : SINGLE SIGN-ON GOUVERNEMENTAL

Un fonctionnaire (policier, médecin, agent fiscal) ne doit pas avoir 15 mots de passe différents pour accéder aux systèmes de l'État. 
L'authentification est centralisée via Keycloak (Phase 1), qui agit comme **Identity Provider (IdP) Souverain**.

## 2. MODÈLE DE FÉDÉRATION (TRUST FEDERATION)

### 2.1 Agences avec leur propre Annuaire (Active Directory)
Exemple: Le Ministère des Finances (DGI) a déjà un Active Directory (AD).
- Keycloak SNISID configure une **Fédération SAML ou LDAP/AD** avec le serveur DGI.
- L'agent DGI se connecte avec ses identifiants habituels.
- Keycloak valide via AD, génère un JWT SNISID avec les rôles globaux, et l'agent accède à la plateforme d'interopérabilité.

### 2.2 Agences sans Annuaire
Exemple: Un commissariat rural.
- Keycloak agit comme l'annuaire primaire (Identity Broker).
- L'agent se connecte avec sa Carte Puce (PKCS#11) ou YubiKey (FIDO2/WebAuthn).

## 3. DELEGATED AUTHORIZATION (OAuth 2.0)

Lorsqu'un citoyen utilise le portail gouvernemental pour demander un casier judiciaire :
1. Le citoyen s'authentifie via SNISID (OIDC).
2. L'application "Casier Judiciaire" demande l'autorisation d'accéder au dossier (Scope: `justice:read_record`).
3. Un jeton (Access Token) à courte durée de vie est généré et passé dans chaque requête API.

---
*Document ID: SNISID-FEDERATION-001 | Approuvé par: CISO National*
