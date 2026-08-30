# Identity Workflows Specification

> **Workflows d'Identité Nationale — SNISID**  
> **Version :** 1.0.0  
> **Classification :** SOUVERAIN — RESTREINT  
> **Dernière mise à jour :** 2026-05-25

---

## 1. OBJECTIF

Industrialiser l'identité nationale. Tous les workflows doivent être :
- **Signés** — Chaque étape est signée cryptographiquement
- **Observables** — Visibilité temps réel sur l'état
- **Audité** — Historique complet et immuable
- **Event-driven** — Déclenchés par événements

---

## 2. WORKFLOWS BPMN

### 2.1 Enrollment — Création d'Identité

```
Start → [Collecter Données Citoyen] → [Vérifier Documents] → [Capture Biométrique]
  → [Liveness Check] → [Duplicate Detection]
    → Si Unique → [Créer NNU] → [Signer Identité] → [Archiver] → [Notifier Citoyen] → End
    → Si Doublon → [Investigation Manuelle]
      → Si Résolu → [Créer NNU] → [Signer] → [Archiver] → [Notifier] → End
      → Si Non Résolu → [Rejeter + Escalade] → End
```

### 2.2 Verification — Vérification d'Identité

```
Start → [Authentifier Demandeur (MFA + Rôle)] → [Vérifier Consentement]
  → Si Consent → [Rechercher Identité] → [Vérifier Signature PKI]
    → Si Valide → [Retourner Résultat + Preuve] → End
    → Si Invalide → [Alerte Fraude] → End
  → Si Non Consent → [Demander Consentement] → [Re-demander]
```

### 2.3 Recovery — Récupération d'Identité

```
Start → [Authentifier Citoyen (Questions + Biométrie alt.)]
  → Si Auth OK → [Révoquer Anciens Accès] → [Générer Nouveaux Accès] → [Notifier Citoyen] → End
  → Si Échec → [Escalade Manuelle] → [Authentifier] → [Suite]
```

### 2.4 Revocation — Révocation d'Identité

```
Start → [Vérifier Autorisation (Rôle + MFA)]
  → Si Autorisé → [Documenter Raison + Preuves] → [Révoquer Identité]
    → NNU → revoked, Certs → CRL, Tokens → invalidés, Sessions → terminées
    → [Notifier Parties] → [Archiver Décision] → End
  → Si Non Autorisé → [Rejeter Demande] → End
```

### 2.5 Correction — Correction d'Identité

```
Start → [Recevoir Demande] → [Vérifier Justificatifs]
  → Si Valide → [Créer Nouvelle Version] → [Ancienne → archive, Nouvelle → active]
    → [Signer Correction (PKI)] → [Notifier Citoyen] → End
  → Si Invalide → [Rejeter + Notifier] → End
```

### 2.6 Appeal — Contestation d'Identité

```
Start → [Recevoir Contestation] → [Enregistrer (Ticket pending)]
  → [Assigner Investigation] → [Investiguer Dossier]
    → Si Fondée → [Corriger Identité] → [Décision Finale (Comité IAM)] → End
    → Si Non Fondée → [Maintenir Décision] → [Décision Finale] → End
  → [Notifier Citoyen (Résultat + Voies recours)]
```

---

## 3. PROPRIÉTÉS DES WORKFLOWS

### Signature
```yaml
workflow_signature:
  requirement: every_step_signed
  algorithm: ECDSA P-384
  key_source: PKI Nationale
  verification: automatic
  storage: audit_trail
```

### Observabilité
```yaml
workflow_observability:
  real_time_tracking: true
  metrics: [workflow_start_time, workflow_end_time, step_duration, step_status, error_count]
  tracing: OpenTelemetry
  dashboard: Grafana
  alerts: [step_timeout, workflow_failure, anomaly_detection]
```

### Auditabilité
```yaml
workflow_audit:
  complete_history: true
  immutable_storage: true
  retention_period: 10_years
  access: [audit_officer, iam_admin, judicial]
  verification: cryptographic_signature
```

### Event-Driven
```yaml
workflow_events:
  trigger: event_bus (Kafka)
  events:
    - identity.created, identity.verified, identity.revoked
    - identity.corrected, identity.contested
    - biometric.captured, biometric.matched
    - consent.granted, consent.revoked
    - fraud.detected
  consumers: [workflow_engine, notification_system, audit_system, analytics_engine]
```

---

## 4. MOTEUR DE WORKFLOW

### Définition de Workflow (Exemple: Enrollment)

```yaml
workflow_definition:
  id: enrollment_v1
  name: Enrollment Workflow
  version: 1.0.0
  trigger:
    type: event
    event: identity.enrollment_requested
  
  steps:
    - id: collect_data
      name: Collect Citizen Data
      type: manual
      actor: enrollment_officer
      timeout: 3600
    
    - id: verify_documents
      name: Verify Documents
      type: automated
      service: document_verification
      timeout: 300
    
    - id: capture_biometrics
      name: Capture Biometrics
      type: manual
      actor: enrollment_officer
      timeout: 600
    
    - id: liveness_check
      name: Liveness Detection
      type: automated
      service: biometric_liveness
      timeout: 60
    
    - id: duplicate_detection
      name: Duplicate Detection
      type: automated
      service: duplicate_engine
      timeout: 30
    
    - id: generate_nnu
      name: Generate NNU
      type: automated
      service: nnu_generator
      timeout: 10
    
    - id: sign_identity
      name: Sign Identity
      type: automated
      service: pki_signing
      timeout: 5
    
    - id: archive_identity
      name: Archive Identity
      type: automated
      service: archive_service
      timeout: 10
    
    - id: notify_citizen
      name: Notify Citizen
      type: automated
      service: notification_service
      timeout: 30
  
  error_handling:
    retry_count: 3
    retry_delay: 5s
    fallback: manual_review
    notification: true
```

---

> **Tous les workflows sont signés, observables, audités et event-driven.**
