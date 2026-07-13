---
# ============================================================
# SNISID API Ecosystem — Master Architecture
# Plateforme Nationale d'Interopérabilité B2B & G2G
# Document ID: SNISID-API-ECO-001
# Version: 1.0.0
# ============================================================

## 1. VISION: L'ÉTAT COMME PLATEFORME (Government-as-a-Platform)

La Phase 10 transforme le SNISID d'un système d'information fermé (G2G - Government to Government) en une plateforme ouverte et sécurisée (B2B - Business to Business). L'écosystème API permet au secteur privé (Banques, Télécoms, Assurances) d'interagir avec les registres de l'État sans jamais compromettre la souveraineté ou la confidentialité des données.

## 2. ARCHITECTURE LOGIQUE DE L'ÉCOSYSTÈME

L'écosystème s'appuie sur le composant API Gateway (Kong) déployé en Phase 4, mais y ajoute les briques de gestion (API Management) :

1. **Le Developer Portal :** Interface web où les développeurs du secteur privé s'inscrivent, consultent la documentation (Swagger/OpenAPI) et demandent des clés d'accès.
2. **Le Plan de Contrôle (Control Plane) :** Gère le cycle de vie des APIs, le routage, et la distribution des certificats mTLS.
3. **Le Plan de Données (Data Plane) :** Les nœuds Kong Ingress qui traitent les millions de requêtes par seconde, exécutent les plugins de sécurité (WAF, Rate Limiting) et transmettent au backend (Istio Service Mesh).
4. **L'API Registry :** Le catalogue centralisé des contrats d'interface.

## 3. LE MODÈLE OPEN BANKING (e-KYC)

Le cas d'usage principal de la Phase 10 est la vérification d'identité (e-KYC) pour lutter contre le blanchiment d'argent et le financement du terrorisme (Loi Sanction).

**Flux d'une transaction Open Banking :**
1. Un citoyen se présente à la Sogebank pour ouvrir un compte.
2. Le guichetier scanne l'empreinte digitale du citoyen et tape son NIU (Numéro d'Identification Unique).
3. Le système de la banque appelle l'API SNISID `POST /b2b/v1/identity/verify`.
4. Le SNISID Gateway (Kong) vérifie le certificat client mTLS de la banque et son token OAuth2 (Keycloak).
5. La requête est routée vers l'ABIS GPU (Phase 2).
6. Le SNISID répond :
   ```json
   {
     "status": "success",
     "match": true,
     "confidence_score": 99.8,
     "timestamp": "2026-05-25T14:30:00Z",
     "transaction_id": "tx_987654321"
   }
   ```
> [!IMPORTANT]
> **Privacy by Design :** L'API ne renvoie *aucune* donnée démographique (Nom, Prénom, Date de Naissance, Adresse). Elle confirme uniquement que l'empreinte correspond au NIU fourni. Si la banque veut le nom, elle doit le demander au citoyen. L'État ne devient pas un courtier en données.

## 4. TOPOLOGIE RÉSEAU & DMZ

Les requêtes B2B ne pénètrent jamais directement dans le cœur du Datacenter de l'État.
- **DMZ Externe :** WAF (Web Application Firewall) + DDoS Protection.
- **DMZ Interne :** API Gateway (Kong) configuré en terminaison TLS mutuelle (mTLS).
- **Core Network :** Microservices métiers (Identity, Police) protégés par Cilium (eBPF) et Istio. Seul le trafic provenant de la DMZ Interne et validé cryptographiquement est autorisé à passer.

---
*Document ID: SNISID-API-ECO-001 | Approuvé par: Chief API Architect*
