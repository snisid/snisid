---
# ============================================================
# SNISID-Interop — National Service Mesh
# Topologie Istio & mTLS Intra-Cluster
# Document ID: SNISID-SERVICE-MESH-001
# Version: 1.0.0
# ============================================================

## 1. LE SERVICE MESH GOUVERNEMENTAL (ISTIO)

Afin d'implémenter un modèle "Zero Trust" complet, la sécurité périmétrique (Firewall/Gateway) ne suffit pas. Une fois dans le réseau de l'État, aucun microservice ne fait confiance à un autre microservice par défaut.
Istio injecte un proxy "Envoy" dans chaque Pod pour gérer tout le trafic réseau.

## 2. FONCTIONNALITÉS CLÉS

### 2.1 Mutual TLS (mTLS Strict)
- Tout le trafic entre, par exemple, le `identity-service` et le `biometric-service` est chiffré.
- Les certificats TLS sont générés automatiquement par le plan de contrôle Istio (qui peut se lier au Vault SNISID) et ont une courte durée de vie (quelques heures) pour limiter l'impact d'une compromission.
- Un Pod PNH qui essaie de contacter le `biometric-service` en HTTP clair sera rejeté avec une erreur `503`.

### 2.2 Circuit Breaking & Retries
Pour éviter les pannes en cascade si la base de données de la DGI répond lentement :
```yaml
# Exemple de DestinationRule Istio pour le service DGI
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: dgi-tax-service
spec:
  host: dgi-tax-service.snisid-interop.svc.cluster.local
  trafficPolicy:
    connectionPool:
      http:
        http1MaxPendingRequests: 100
        maxRequestsPerConnection: 10
    outlierDetection:
      consecutive5xxErrors: 3
      interval: 10s
      baseEjectionTime: 30s
      maxEjectionPercent: 100
```
- Si le service DGI renvoie 3 erreurs 5xx de suite, Istio "ouvre le circuit" et rejette instantanément le trafic pendant 30 secondes pour laisser le service souffler.

### 2.3 Traffic Shifting (Canary Deployments)
Lors de la mise à jour d'un algorithme biométrique, Istio permet de rediriger 5% du trafic vers la nouvelle version (Canary) pour vérifier la stabilité avant le déploiement global, sans toucher au code applicatif.

---
*Document ID: SNISID-SERVICE-MESH-001 | Approuvé par: Platform Engineering Lead*
