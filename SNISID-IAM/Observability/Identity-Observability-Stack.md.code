# Identity Observability Stack

> **Observabilité Identité Temps Réel — SNISID**  
> **Version :** 1.0.0  
> **Classification :** SOUVERAIN — RESTREINT  
> **Dernière mise à jour :** 2026-05-25

---

## 1. OBJECTIF

Voir l'identité en temps réel. Monitoring complet des opérations IAM nationales.

---

## 2. DOMAINES DE MONITORING

| Domaine | Monitoring | Seuil d'alerte |
|---------|-----------|----------------|
| Failed logins | ✅ | > 10/min |
| MFA failures | ✅ | > 5/min |
| Fraud attempts | ✅ | > 0 |
| Duplicate attempts | ✅ | > 5/h |
| Suspicious access | ✅ | > 3/h |
| Biometric failures | ✅ | > 2% |
| Consent violations | ✅ | > 0 |
| NNU operations | ✅ | > 100/min |
| PKI operations | ✅ | > 50/min |
| Federation errors | ✅ | > 1% |

---

## 3. OUTILS D'OBSERVABILITÉ

| Domaine | Outil | Usage |
|---------|-------|-------|
| Metrics | **Prometheus** | Collecte et alerting |
| Logs | **Loki** | Agrégation et recherche |
| SIEM | **OpenSearch** | Analyse sécurité |
| Tracing | **OpenTelemetry** | Traçabilité distribuée |
| Dashboards | **Grafana** | Visualisation |
| Alerting | **AlertManager** | Notifications |

### Architecture

```
[IAM Services] [Biometric] [Consent] [Federation]
         │         │           │          │
         └─────────┴───────────┴──────────┘
                           │
                   [OpenTelemetry Collector]
                           │
         ┌─────────────────┼─────────────────┐
         │                 │                 │
    [Prometheus]       [Loki]          [OpenSearch]
      Metrics           Logs              SIEM
         │                 │                 │
         └─────────────────┼─────────────────┘
                           │
                      [Grafana]
                      Dashboards
                           │
                    [AlertManager]
                   Notifications
```

---

## 4. MÉTRIQUES PROMÉTHÉUS

### Identité
- `snisid_identity_total` (counter) — Total identités créées
- `snisid_identity_active` (gauge) — Identités actives
- `snisid_identity_verification_duration_seconds` (histogram)
- `snisid_identity_duplicate_detected_total` (counter)
- `snisid_identity_fraud_detected_total` (counter)

### Authentification
- `snisid_auth_login_total` (counter) — Tentatives de connexion
- `snisid_auth_mfa_total` (counter) — Tentatives MFA
- `snisid_auth_session_total` (gauge) — Sessions actives

### Biométrique
- `snisid_biometric_capture_total` (counter)
- `snisid_biometric_match_duration_seconds` (histogram)
- `snisid_biometric_match_score` (histogram)
- `snisid_biometric_liveness_pass_rate` (gauge)

### Consentement
- `snisid_consent_granted_total` (counter)
- `snisid_consent_revoked_total` (counter)
- `snisid_consent_violation_total` (counter)
- `snisid_consent_active` (gauge)

---

## 5. RÈGLES D'ALERTE

### Alertes Critiques

| Alerte | Expression | Durée | Équipe |
|--------|-----------|-------|--------|
| IdentityFraudDetected | `increase(fraud_total[5m]) > 0` | 1m | Security |
| MFAFailureSpike | `increase(mfa_failed[5m]) > 20` | 2m | Security |
| SuspiciousAccessPattern | `increase(login_failed[5m]) > 50` | 2m | Security |
| BiometricMatchAnomaly | `liveness_pass_rate < 0.80` | 5m | Biometric |
| ConsentViolation | `increase(consent_violation[5m]) > 0` | 1m | Compliance |
| DuplicateIdentitySpike | `increase(duplicate[1h]) > 10` | 5m | Identity |

### Alertes Moyennes

| Alerte | Expression | Durée | Équipe |
|--------|-----------|-------|--------|
| HighAuthLatency | `P95 auth_duration > 2s` | 5m | Platform |
| HighBiometricLatency | `P95 biometric_duration > 5s` | 5m | Biometric |
| SessionExpiryApproaching | `expiring_sessions > 100` | 5m | Platform |

---

## 6. TRACING OPENTELEMETRY

### Services Tracés
- iam-service, nnu-service, biometric-service
- consent-service, federation-broker
- workflow-engine, fraud-engine

### Attributs de Span
```yaml
identity: [citizen.nnu, actor.role, operation.type, result.status, duration.ms]
security: [mfa.verified, abac.decision, risk.score, consent.granted]
biometric: [modality.type, quality.score, liveness.passed, match.confidence]
```

### Propagation
- Format: W3C TraceContext
- Headers: traceparent, tracestate

---

## 7. LOGGING LOKI

### Structure des Logs
```json
{
  "timestamp": "2026-05-25T10:00:00.000Z",
  "level": "info",
  "service": "iam-service",
  "trace_id": "otel-trace-id",
  "labels": {
    "environment": "production",
    "component": "identity",
    "operation": "enrollment",
    "status": "success"
  },
  "message": "Identity created successfully",
  "context": {
    "enrollment_id": "uuid-enroll",
    "nnu_generated": "NNU-HA-20260525-A7K9M2X4-73-1"
  }
}
```

### Rétention

| Type de log | Rétention | Stockage |
|-------------|-----------|----------|
| Audit logs | 10 ans | Immutable |
| Security logs | 5 ans | Chiffré |
| Operation logs | 1 an | Standard |
| Debug logs | 30 jours | Standard |

---

## 8. SIEM OPENSEARCH

### Règles de Détection

| Règle | Condition | Sévérité | Réponse |
|-------|-----------|----------|---------|
| Brute Force Authentication | `login_failed > 50 en 5min` | High | block_ip + alert |
| Identity Anomaly | `duplicate_detected > 5 en 1h` | Critical | suspend_enrollment |
| Consent Violation | `consent_violation > 0 en 5min` | Critical | block_access + alert |
| Biometric Spoofing | `liveness_pass_rate < 50% sur 10 captures` | Critical | block_device + alert |

---

> **L'observabilité identité est temps réel, complète et alertée.**
