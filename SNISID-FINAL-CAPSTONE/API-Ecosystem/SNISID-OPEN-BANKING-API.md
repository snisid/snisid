---
# ============================================================
# SNISID Capstone — National API Ecosystem (Phase 10)
# Open Banking & Monétisation des Données
# Document ID: SNISID-CAP-API-001
# Version: 1.0.0
# ============================================================

## 1. OUVERTURE AU SECTEUR PRIVÉ (B2B)

Jusqu'à la Phase 9, le SNISID était un écosystème fermé (Intra-Gouvernemental). La Phase 10 expose l'API SNISID (via l'API Gateway de la Phase 4) aux acteurs privés certifiés (Banques, Télécoms, Assurances).

## 2. CAS D'USAGE "OPEN BANKING" (KYC)

La loi haïtienne exige que les banques (Unibank, Sogebank) vérifient l'identité de leurs clients (Know Your Customer - KYC) pour éviter le blanchiment d'argent.
- La banque appelle l'endpoint : `POST /v1/b2b/verify-identity`
- Elle envoie le numéro NIU et l'empreinte digitale scannée du client en agence.
- Le SNISID répond `Match: True` ou `Match: False`. **Aucune donnée démographique ou biométrique n'est renvoyée à la banque**, uniquement une confirmation cryptographique.

## 3. MODÈLE ÉCONOMIQUE (Monétisation)

Pour assurer la survie financière du Datacenter (Phase 5) sans dépendre exclusivement des impôts, l'API est monétisée.
- Chaque requête de vérification biométrique coûte 10 Gourdes (~0.07 USD) à la banque.
- Avec 10 millions de transactions par an, cela génère un revenu de 100 millions de Gourdes pour l'entretien des serveurs et le salaire des ingénieurs (Phase 16).

---
*Document ID: SNISID-CAP-API-001 | Approuvé par: Banque de la République d'Haïti (BRH)*
