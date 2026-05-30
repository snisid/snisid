---
# ============================================================
# SNISID API Ecosystem — National API Security Model
# mTLS, WAF et Service Mesh B2B
# Document ID: SNISID-API-SEC-001
# Version: 1.0.0
# ============================================================

## 1. POSTULAT DE SÉCURITÉ B2B (Zéro Confiance Externe)

Bien que nous signions des contrats avec les banques (Sogebank, Unibank), le SNISID part du principe que **le réseau de la banque peut être compromis**.
Si un hacker pénètre la banque, il ne doit pas pouvoir utiliser la connexion B2B pour aspirer la base de données de l'État Haïtien.

## 2. MUTUAL TLS (mTLS) OBLIGATOIRE

Un simple token API (Bearer Token) n'est pas suffisant car il peut être volé.
L'API SNISID exige le **mTLS (Mutual TLS)**. 
- La banque doit présenter un certificat cryptographique X.509 signé par la PKI de l'État lors de l'établissement de la connexion TCP.
- Si le certificat est manquant, expiré, ou révoqué (CRL), Kong ferme la connexion avant même de lire la requête HTTP (Drop TCP).
- Ce certificat est physiquement lié au serveur de la banque. Il ne peut pas être volé et utilisé depuis un ordinateur portable externe.

## 3. WEB APPLICATION FIREWALL (WAF)

Chaque requête entrante B2B passe par le WAF (Coraza/ModSecurity) intégré à Kong :
- **Validation Strict du Schéma (OpenAPI Validation) :** Si la requête JSON contient un champ non documenté (tentative d'injection), elle est bloquée.
- **Payload Size Limit :** Limite stricte à 1MB par requête pour prévenir les attaques par déni de service (DDoS) asymétrique.
- **SQL/NoSQL Injection Detection :** Blocage des signatures connues.

## 4. INTÉGRATION SERVICE MESH (ISTIO)

Une fois la requête passée par Kong, elle entre dans le **Service Mesh (Istio)**.
Kong injecte un header `X-Consumer-ID`. Istio utilise les politiques d'autorisation (AuthorizationPolicy) pour vérifier que ce Consumer ID a bien le droit de parler au pod `identity-service`. S'il tente de parler au pod `justice-service`, Istio bloque le trafic latéralement.

---
*Document ID: SNISID-API-SEC-001 | Approuvé par: CISO National*
