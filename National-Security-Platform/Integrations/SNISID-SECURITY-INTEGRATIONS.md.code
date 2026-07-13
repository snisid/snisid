---
# ============================================================
# SNISID-Security — National Security Integrations
# Fédération, Interopérabilité et API Souveraines
# Document ID: SNISID-SEC-INTEG-001
# Version: 1.0.0
# ============================================================

## 1. ÉCOSYSTÈME D'INTÉGRATION

La plateforme SNISID expose des APIs sécurisées via Kong Gateway (Phase 1) pour permettre aux différentes agences de consommer et de produire des événements.

## 2. CATALOGUE DES APIS DE SÉCURITÉ

Toutes les requêtes exigent un JWT signé par Keycloak (mTLS obligatoire).

| API Endpoint | Agence Consommatrice | Action | Protocole |
|--------------|----------------------|--------|-----------|
| `POST /v1/justice/cases` | PNH / Parquet | Ouvrir un dossier pénal | REST |
| `GET /v1/intelligence/graph` | DCPJ | Analyser les liens criminels | GraphQL |
| `POST /v1/border/crossings` | DGIE (Immigration) | Enregistrer une entrée/sortie | REST |
| `GET /v1/prison/inmates` | Prison / Justice | Liste des détenus | REST |
| `POST /v1/evidence/upload` | PNH Scientifique | Uploader une preuve numérique | REST/S3 |
| `wss://alerts.snisid.gov.ht` | Toutes agences | Recevoir les alertes Push | WebSockets |

## 3. IDENTITY FEDERATION (OIDC/SAML)

Les systèmes legacy de l'État (ex: base de données du Ministère de la Justice si existante) n'ont pas besoin de gérer des mots de passe. Ils s'intègrent au SNISID via **Keycloak (Single Sign-On)**.
Un juge se connecte avec sa carte à puce gouvernementale ; Keycloak valide le certificat et émet un jeton OIDC.

---
*Document ID: SNISID-SEC-INTEG-001 | Approuvé par: AND (Autorité Nationale Numérique)*
