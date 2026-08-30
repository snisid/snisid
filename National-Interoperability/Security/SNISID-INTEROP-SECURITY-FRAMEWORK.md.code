---
# ============================================================
# SNISID-Interop — Security Framework
# Sécurité des échanges API & Zero Trust
# Document ID: SNISID-INT-SEC-001
# Version: 1.0.0
# ============================================================

## 1. MODÈLE DE SÉCURITÉ DE L'INTEROPÉRABILITÉ

Le framework de sécurité de la Phase 4 empêche la compromission d'une agence (ex: Ministère de la Santé hacké) d'affecter le reste de l'État (Isolation des "Blast Radius").

## 2. PILIERS DE LA SÉCURITÉ DES ÉCHANGES

1. **Authentication (OIDC/mTLS) :** Chaque appel API sortant d'une agence vers le SNISID doit utiliser un jeton OIDC obtenu via un certificat mTLS client valide. Les clés d'API (API Keys) statiques sont strictement **bannies** en production.
2. **Authorization (OPA) :** L'Open Policy Agent vérifie que l'agence a le droit légal de faire cette requête.
3. **Data Loss Prevention (DLP) :** Le Service Mesh inspecte les payloads. Si le numéro de carte de crédit ou des données médicales non autorisées tentent de sortir, la connexion est coupée.
4. **Transaction Signing :** Pour les actions d'écriture (ex: `POST /v1/civil/deaths`), le payload JSON doit inclure une signature numérique (JWS - JSON Web Signature) par la clé privée de l'agent effectuant l'action.
5. **Rate Limiting Hard-Coded :** Empêche l'exfiltration massive de données. (Ex: Max 10 000 requêtes/jour pour une banque).

---
*Document ID: SNISID-INT-SEC-001 | Approuvé par: CISO National*
