---
# ============================================================
# SNISID API Ecosystem — National Developer Portal
# La Vitrine Technique de l'État pour le Secteur Privé
# Document ID: SNISID-API-PORTAL-001
# Version: 1.0.0
# ============================================================

## 1. LE PORTAIL DÉVELOPPEUR (DEV PORTAL)

Le Portail Développeur (`developer.snisid.gouv.ht`) est l'interface en libre-service ("Self-Service") permettant aux ingénieurs du secteur privé de s'intégrer au SNISID sans intervention manuelle excessive de l'État.
Il est construit avec **Backstage** (Framework de portail développé par Spotify et CNCF).

## 2. FONCTIONNALITÉS DU PORTAIL

- **Onboarding Automatisé :** Une banque soumet sa "Demande d'Accès B2B" via le portail. Elle y joint son certificat fiscal et l'approbation de la Banque Centrale (BRH). Un workflow administratif (Phase 11) valide la demande.
- **Documentation Interactive :** Toutes les APIs sont documentées au format OpenAPI 3.0 (Swagger UI). Les développeurs peuvent tester des requêtes mockées directement depuis leur navigateur.
- **Gestion des Clés :** Une fois validée, la banque génère ses propres tokens OAuth2 et télécharge son certificat mTLS via une interface sécurisée.
- **Billing Dashboard :** Tableau de bord affichant en temps réel la consommation d'API de la banque et la facturation mensuelle estimée (selon le modèle de tarification de la Gouvernance).

## 3. ENVIRONNEMENTS (Sandbox & Production)

L'État met à disposition deux environnements :
1. **Environnement Sandbox (Bac à sable) :** Données synthétiques fictives. Gratuit. Permet aux développeurs des banques de coder et tester leur intégration. Rate limit bas (10 req/s).
2. **Environnement Production :** Accès aux vraies données. Payant. Nécessite une adresse IP Whitelistée et un certificat mTLS validé par l'Autorité de Certification de l'État (Phase 6). Rate limit selon le contrat.

---
*Document ID: SNISID-API-PORTAL-001 | Approuvé par: Chief Developer Advocate*
