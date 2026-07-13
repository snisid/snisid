---
# ============================================================
# SNISID-Interop — Audit & Traceability Engine
# Logs immuables et Traces Forensiques
# Document ID: SNISID-AUDIT-001
# Version: 1.0.0
# ============================================================

## 1. TRACEABILITY ENGINE

Pour qu'un écosystème inter-agences fonctionne, la confiance est requise. La confiance repose sur la vérifiabilité. 
Le moteur de traçabilité enregistre "Qui a fait quoi, quand, au nom de qui, et depuis quelle adresse IP".

## 2. API AUDIT (Kong + OpenTelemetry)

Toutes les requêtes passant par Kong sont tracées de bout en bout (End-to-End Tracing avec Tempo).
- **Correlation ID :** Un en-tête `X-Correlation-ID` est généré à l'entrée de l'API Gateway.
- Il est propagé par le Service Mesh (Istio) jusqu'à la base de données.
- Si une transaction échoue au niveau de l'Identity Registry, le log dans Loki contiendra ce Correlation ID, permettant de lier l'erreur à la requête originale de la banque ou du ministère.

## 3. IMMUTABLE FORENSIC LOGGING (Kafka)

Toute requête en modification (POST/PUT/DELETE) génère un événement d'audit asynchrone envoyé sur le topic Kafka protégé `snisid.platform.audit`.
Ce topic a une rétention **infinie** et son accès en lecture est physiquement restreint au CISO National et à la Cour Supérieure des Comptes (CSCCA).

```json
{
  "timestamp": "2026-05-25T14:30:00Z",
  "correlation_id": "req-9f8a-4c21",
  "agency_id": "DGI",
  "agent_niu": "1234567890",
  "action": "READ_IDENTITY",
  "target_niu": "0987654321",
  "authorized_by_policy": "dgi-tax-read-policy-v2",
  "response_status": 200
}
```

---
*Document ID: SNISID-AUDIT-001 | Approuvé par: CSCCA*
