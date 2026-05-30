---
# ============================================================
# SNISID API Ecosystem — National API Governance & Monetization
# Monétisation, SLA et Régulation
# Document ID: SNISID-API-GOV-001
# Version: 1.0.0
# ============================================================

## 1. LE MODÈLE DE GOUVERNANCE

L'accès aux APIs de l'État n'est pas un droit absolu pour le secteur privé, c'est un privilège contractuel. La gouvernance des APIs SNISID définit les règles d'engagement, les limites techniques (Quotas) et le modèle financier.

## 2. MODÈLE ÉCONOMIQUE SOUVERAIN (Monétisation B2B)

Pour garantir l'autonomie financière du SNISID (Maintenance des serveurs, Licences support, Salaires des ingénieurs d'élite), l'accès à l'API de vérification e-KYC est facturé.

**Pricing Model (Tiered) :**
- **Tier 1 (Startup / Fintech) :** Jusqu'à 10 000 req/mois -> 15 Gourdes / req.
- **Tier 2 (Banques Commerciales) :** De 10 001 à 500 000 req/mois -> 10 Gourdes / req.
- **Tier 3 (Télécoms - Digicel/Natcom) :** > 500 000 req/mois -> 7 Gourdes / req.
- **G2G (Ministères) :** Gratuit (Budget de l'État).

Cette monétisation génère des revenus massifs (estimés à 300 millions de Gourdes par an) qui sont sanctuarisés sur un compte de la Banque Centrale (BRH) dédié au "Fonds de Résilience Numérique".

## 3. SLA (Service Level Agreement) ET RATE LIMITING

- **SLA :** L'État garantit une disponibilité de 99.99% pour l'API B2B.
- **Rate Limiting (Protection contre l'épuisement) :** Le plugin `rate-limiting-advanced` de Kong est configuré par "Consumer" (Client).
  - Si la Sogebank achète un forfait de 100 req/sec, toute requête supplémentaire est rejetée avec un code HTTP `429 Too Many Requests`.
- **Spike Arrest :** Lissage des pics de trafic imprévus pour éviter de surcharger les GPU de l'ABIS biométrique.

## 4. CONDITIONS JURIDIQUES DE RÉVOCATION

L'État se réserve le droit de révoquer instantanément la clé API (via Keycloak) d'un acteur privé si le SOC (Phase 6) détecte :
- Un ratio anormal de requêtes de type "Brute Force" (Tentative de deviner des NIU).
- Une requête provenant d'une adresse IP non déclarée dans la liste blanche (IP Whitelisting stricte imposée).
- Une faille de sécurité documentée chez le partenaire privé.

---
*Document ID: SNISID-API-GOV-001 | Approuvé par: Direction de la Régulation Économique*
