# Identity Event Mesh Specification

> **Réseau d'Événements d'Identité — SNISID**  
> **Version :** 1.0.0  
> **Classification :** SOUVERAIN — RESTREINT  
> **Dernière mise à jour :** 2026-05-25

---

## 1. OBJECTIF

Créer les événements identité nationaux. Tous les événements doivent passer via **Kafka sécurisé**.

---

## 2. CATALOGUE D'ÉVÉNEMENTS

| Événement | Description | Topic | Priorité |
|-----------|-------------|-------|----------|
| `identity.created` | Création d'identité | `snisid.identity.lifecycle` | Haute |
| `identity.verified` | Vérification d'identité | `snisid.identity.lifecycle` | Moyenne |
| `identity.revoked` | Révocation d'identité | `snisid.identity.lifecycle` | Critique |
| `identity.suspended` | Suspension d'identité | `snisid.identity.lifecycle` | Haute |
| `identity.corrected` | Correction d'identité | `snisid.identity.lifecycle` | Moyenne |
| `identity.contested` | Contestation d'identité | `snisid.identity.lifecycle` | Haute |
| `biometric.captured` | Capture biométrique | `snisid.biometric.events` | Moyenne |
| `biometric.matched` | Match biométrique | `snisid.biometric.events` | Haute |
| `biometric.failed` | Échec biométrique | `snisid.biometric.events` | Moyenne |
| `fraud.detected` | Détection fraude | `snisid.fraud.events` | Critique |
| `fraud.investigated` | Investigation fraude | `snisid.fraud.events` | Haute |
| `consent.granted` | Consentement accordé | `snisid.consent.events` | Haute |
| `consent.revoked` | Consentement révoqué | `snisid.consent.events` | Haute |
| `consent.expired` | Consentement expiré | `snisid.consent.events` | Moyenne |
| `access.granted` | Accès accordé | `snisid.access.events` | Moyenne |
| `access.denied` | Accès refusé | `snisid.access.events` | Haute |
| `access.anomalous` | Accès anormal | `snisid.access.events` | Critique |
| `nnu.generated` | NNU généré | `snisid.nnu.events` | Haute |
| `nnu.revoked` | NNU révoqué | `snisid.nnu.events` | Critique |
| `certificate.issued` | Certificat émis | `snisid.pki.events` | Haute |
| `certificate.revoked` | Certificat révoqué | `snisid.pki.events` | Critique |
| `mfa.completed` | MFA terminé | `snisid.auth.events` | Moyenne |
| `mfa.failed` | MFA échoué | `snisid.auth.events` | Haute |
| `login.success` | Connexion réussie | `snisid.auth.events` | Moyenne |
| `login.failed` | Connexion échouée | `snisid.auth.events` | Haute |
| `session.created` | Session créée | `snisid.session.events` | Moyenne |
| `session.terminated` | Session terminée | `snisid.session.events` | Moyenne |
| `federation.trust.established` | Confiance établie | `snisid.federation.events` | Haute |
| `federation.trust.revoked` | Confiance révoquée | `snisid.federation.events` | Critique |

---

## 3. SCHÉMA D'ÉVÉNEMENT STANDARD

```json
{
  "event_id": "uuid-v4",
  "event_type": "identity.created",
  "event_version": "1.0.0",
  "timestamp": "2026-05-25T10:00:00Z",
  "source": "snisid-iam-service",
  "topic": "snisid.identity.lifecycle",
  "partition_key": "NNU-HA-20260525-A7K9M2X4-73-1",
  "headers": {
    "correlation_id": "uuid-corr",
    "trace_id": "otel-trace-id",
    "environment": "production",
    "priority": "high"
  },
  "data": {
    "citizen_nnu": "NNU-HA-20260525-A7K9M2X4-73-1",
    "actor_nnu": "NNU-HA-20250101-B8L3N5Y2-45-1",
    "actor_role": "enrollment_officer",
    "action": "create"
  },
  "metadata": {
    "schema_version": "1.0",
    "content_type": "application/json",
    "signature": "ECDSA-SHA384-signature"
  }
}
```

---

## 4. ARCHITECTURE KAFKA

### Topologie des Topics

| Topic | Partitions | Réplication | Rétention |
|-------|-----------|-------------|-----------|
| `snisid.identity.lifecycle` | 12 | 3 | 10 ans |
| `snisid.biometric.events` | 8 | 3 | 5 ans |
| `snisid.fraud.events` | 6 | 3 | 10 ans |
| `snisid.consent.events` | 8 | 3 | 10 ans |
| `snisid.access.events` | 12 | 3 | 10 ans |
| `snisid.nnu.events` | 6 | 3 | 10 ans |
| `snisid.pki.events` | 4 | 3 | 10 ans |
| `snisid.auth.events` | 12 | 3 | 1 an |
| `snisid.session.events` | 12 | 3 | 30 jours |
| `snisid.federation.events` | 4 | 3 | 5 ans |

### Configuration Sécurité

```yaml
kafka_cluster:
  brokers: 3
  min_insync_replicas: 2
  default_replication_factor: 3
  security:
    protocol: SASL_SSL
    sasl_mechanism: SCRAM-SHA-512
    ssl_endpoint_identification: https
  topics:
    compaction_enabled: true
    cleanup_policy: compact,delete
```

---

## 5. PRODUCTEURS ET CONSOMMATEURS

### Producteurs
- IAM Service → `snisid.identity.lifecycle`
- Biometric Service → `snisid.biometric.events`
- Fraud Engine → `snisid.fraud.events`
- Consent Engine → `snisid.consent.events`
- Access Gateway → `snisid.access.events`
- NNU Service → `snisid.nnu.events`
- PKI Service → `snisid.pki.events`
- Auth Service → `snisid.auth.events`
- Session Manager → `snisid.session.events`
- Federation Service → `snisid.federation.events`

### Consommateurs
- Workflow Engine ← tous les topics
- Notification Service ← lifecycle, consent, fraud
- Audit System ← tous les topics
- Analytics Engine ← tous les topics
- SIEM ← fraud, access, auth
- Dashboard ← tous les topics (agrégé)
- Archive Service ← tous les topics (long terme)

---

## 6. GARANTIES DE LIVRAISON

| Niveau | Topics | Garantie |
|--------|--------|----------|
| Exactly-once | lifecycle, fraud, consent, pki | Transactionnel |
| At-least-once | biometric, access, auth | Ack all |
| At-most-once | session, federation | Fire-and-forget |

### Dead Letter Queue
```yaml
dead_letter_queue:
  topic: snisid.dlq
  retention_ms: 7776000000  # 90 jours
  alerting: true
  manual_review: true
```

---

## 7. INTÉGRATION WORKFLOWS

```yaml
workflow_integration:
  triggers:
    identity.created:
      - activate_wallet
      - send_welcome_notification
      - log_audit_entry
    identity.revoked:
      - invalidate_certificates
      - terminate_sessions
      - notify_federation_partners
      - send_revocation_notification
    fraud.detected:
      - suspend_identity
      - alert_security_team
      - create_investigation_ticket
      - notify_judicial_if_needed
    consent.granted:
      - update_consent_registry
      - notify_requester
      - log_audit_entry
    consent.revoked:
      - update_consent_registry
      - block_pending_access
      - notify_requester
      - log_audit_entry
```

---

## 8. SÉCURITÉ DU MESH

| Couche | Mécanisme |
|--------|-----------|
| Transport | TLS 1.3 + mTLS |
| Authentification | SCRAM-SHA-512 + certificats |
| Données sensibles | Chiffrement applicatif (AES-256) |
| Signature événements | ECDSA P-384 |

### Contrôle d'Accès ACL

```yaml
acl_configuration:
  producers:
    iam_service: {topics: [identity.lifecycle, nnu.events], operations: [write]}
    biometric_service: {topics: [biometric.events], operations: [write]}
    fraud_engine: {topics: [fraud.events], operations: [write]}
    consent_engine: {topics: [consent.events], operations: [write]}
  consumers:
    workflow_engine: {topics: [all], operations: [read]}
    audit_system: {topics: [all], operations: [read]}
    siem: {topics: [fraud.events, access.events, auth.events], operations: [read]}
```

---

## 9. SURVEILLANCE

| Métrique | Seuil | Alerte |
|----------|-------|--------|
| Latence production | < 100ms | > 500ms |
| Latence consommation | < 200ms | > 1s |
| Taille DLQ | < 100 | > 100 |
| Erreurs production | < 0.01% | > 0.1% |
| Erreurs consommation | < 0.01% | > 0.1% |

---

> **Tous les événements passent via Kafka sécurisé avec garanties de livraison.**
