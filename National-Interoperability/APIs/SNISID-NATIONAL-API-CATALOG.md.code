---
# ============================================================
# SNISID-Interop — National API Catalog & Developer Ecosystem
# Portail développeur, SDKs et Versioning
# Document ID: SNISID-API-CATALOG-001
# Version: 1.0.0
# ============================================================

## 1. DEVELOPER PORTAL SOUVERAIN

L'État fournit un portail centralisé (ex: `https://developer.snisid.gov.ht`) basé sur Backstage (Spotify) ou Kong Dev Portal. 
Toute agence gouvernementale (ou entité privée autorisée) souhaitant s'intégrer doit s'y inscrire.

## 2. FONCTIONNALITÉS DU CATALOGUE

1. **API Discovery :** Moteur de recherche pour toutes les APIs de l'État (Identité, Santé, Justice, Impôts).
2. **OpenAPI 3.1 Specs :** Documentation interactive (Swagger UI) pour tester les endpoints.
3. **Sandbox Environment :** Environnement de test avec des données fictives (Dummy Data) pour permettre aux développeurs externes de tester leur code sans toucher la production.
4. **SDK Generation :** Génération automatique de clients en Java, Python, Go et Node.js.

## 3. GOUVERNANCE ET VERSIONING

### 3.1 Versioning Stricte (URI)
Toutes les APIs doivent inclure la version majeure dans l'URI : `https://api.snisid.gov.ht/v1/identity`

### 3.2 Dépréciation (Sunset Policy)
Lorsqu'une `v2` est publiée, la `v1` doit fonctionner pendant au moins 12 mois. Le header HTTP `Deprecation: @Date` est ajouté à toutes les réponses de la `v1` 6 mois avant l'arrêt.

### 3.3 Semantic Versioning (SemVer)
- **MAJEURE** (v2) : Changement cassant (Breaking change - suppression d'un champ).
- **MINEURE** (v1.1) : Ajout d'un champ optionnel (Rétro-compatible).
- **PATCH** (v1.1.1) : Correction de bug interne invisible pour le client.

---
*Document ID: SNISID-API-CATALOG-001 | Approuvé par: Direction de l'Ingénierie Logicielle (AND)*
