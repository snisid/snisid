---
# ============================================================
# SNISID-Interop — National API Gateway Ecosystem
# Exposition des Services de l'État (Kong Enterprise)
# Document ID: SNISID-API-GATEWAY-001
# Version: 1.0.0
# ============================================================

## 1. ARCHITECTURE DE L'API GATEWAY

L'API Gateway Nationale (basée sur Kong) est le point d'entrée et de sortie exclusif pour tout trafic synchrone (REST/gRPC/GraphQL) entre les agences de l'État, et vers l'extérieur (Banques, Télécoms).

```yaml
# Topology logic
Client (App PNH / Banque) 
  --> WAF (Web Application Firewall) 
    --> Kong Gateway (Rate Limiting, JWT Verification) 
      --> Istio Ingress (mTLS) 
        --> SNISID Services
```

## 2. POLITIQUES DE CONSOMMATION (PLANS D'USAGE)

Pour assurer la stabilité et permettre la **monétisation** future de la vérification d'identité (KYC), des Tiers (Niveaux) sont définis.

| Tier | Public Cible | Limite (Rate Limit) | Fonctionnalités | Coût |
|------|--------------|---------------------|-----------------|------|
| **GOV-INTERNAL** | Agences de l'État | 10,000 req/min | Accès complet ABAC | Gratuit |
| **PUBLIC-OPEN** | Citoyens / Apps | 100 req/min | Données publiques, formulaires | Gratuit |
| **PRIVATE-STANDARD** | Banques, Télécoms | 500 req/min | Vérification KYC binaire (Oui/Non) | Payant par req |
| **PRIVATE-PREMIUM** | Grandes institutions | 5,000 req/min | KYC, Biométrie (si mandaté) | Abonnement |

## 3. SÉCURITÉ DE L'API (PLUGINS KONG)

Chaque route exposée est protégée par un ensemble de plugins stricts :

1. **OIDC / OAuth 2.0 Introspection :** Valide le token JWT contre Keycloak.
2. **Correlation ID :** Injecte un `X-Request-ID` pour le traçage distribué (Tempo/Zipkin).
3. **CORS :** Restreint aux origines `.gov.ht` (sauf API publiques).
4. **Rate Limiting Advanced :** Utilise Redis pour bloquer les pics de charge et les attaques DDoS (Layer 7).
5. **OpenTelemetry :** Envoie les métriques de latence et les codes HTTP vers Prometheus.
6. **Data Loss Prevention (DLP) :** Scanne les payloads sortants pour masquer (masking) les données sensibles si le consommateur n'a pas le rôle adéquat.

## 4. WORKFLOW D'INTÉGRATION POUR UN TIERS (Ex: Digicel/Natcom)

Pour qu'un opérateur télécom s'intègre au SNISID pour l'enregistrement obligatoire des cartes SIM :
1. L'opérateur s'inscrit sur le **Developer Portal**.
2. Il soumet ses documents légaux (Validation MANUELLE par l'AND).
3. L'AND génère des `client_id` et `client_secret` restreints au scope `kyc:verify:sim`.
4. Le système de l'opérateur appelle `POST /v1/identity/verify` avec ces credentials via mTLS.

---
*Document ID: SNISID-API-GATEWAY-001 | Approuvé par: Autorité Nationale Numérique*
