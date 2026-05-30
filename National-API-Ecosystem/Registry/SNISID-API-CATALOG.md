---
# ============================================================
# SNISID API Ecosystem — National API Registry & Catalog
# Inventaire Officiel des Interfaces de l'État
# Document ID: SNISID-API-CAT-001
# Version: 1.0.0
# ============================================================

## 1. LE CATALOGUE OFFICIEL (API REGISTRY)

Le Registre des APIs recense l'ensemble des contrats d'interface exposés par l'API Gateway (Kong). Il sert d'annuaire unique (Single Source of Truth) pour éviter la duplication des efforts entre ministères.

## 2. INVENTAIRE DES APIs MAJEURES

### B2B (Exposées au secteur privé payant)
1. **Identity Verification API (`/b2b/v1/identity/verify`) :** Vérification KYC via Empreinte ou Code PIN.
2. **Digital Address API (`/b2b/v1/address/validate`) :** Vérifie l'existence légale d'une adresse géographique (utile pour les livraisons ou les assurances). Ne révèle pas qui y habite.

### G2G (Internes au Gouvernement, Gratuites)
1. **Civil Registry API (`/g2g/v1/civil/acts`) :** Recherche d'actes de naissance/décès. Consommée par le Ministère des Affaires Étrangères (Passeports).
2. **Police Clearance API (`/g2g/v1/justice/clearance`) :** Vérification de l'existence d'un casier judiciaire vierge. Consommée par les Ressources Humaines de l'État lors de recrutements.
3. **Vehicle Registration API (`/g2g/v1/vehicles/{plate}`) :** Recherche du propriétaire d'un véhicule. Consommée par les patrouilles mobiles de la PNH (Phase 8).

## 3. VERSIONING ET DÉPRÉCIATION

Les APIs de l'État respectent la sémantique de versioning dans l'URL (ex: `/v1/`, `/v2/`).
**Politique de dépréciation (Deprecation Policy) :**
- Lorsqu'une `/v2/` est publiée, la `/v1/` est maintenue pendant exactement 18 mois (Sunset Period).
- Au-delà, elle est coupée. Le portail développeur envoie des alertes automatisées aux consommateurs 6 mois, 3 mois et 1 semaine avant la coupure.

---
*Document ID: SNISID-API-CAT-001 | Approuvé par: Enterprise Architecture Board*
