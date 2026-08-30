# Identity Federation Specification

> **Fédération d'Identité Nationale — SNISID**  
> **Version :** 1.0.0  
> **Classification :** SOUVERAIN — RESTREINT  
> **Dernière mise à jour :** 2026-05-25

---

## 1. OBJECTIF

Permettre la confiance inter-agences. Toutes les agences gouvernementales partagent la même confiance identité.

---

## 2. STANDARDS SUPPORTÉS

| Standard | Support | Usage |
|----------|---------|-------|
| OIDC (OpenID Connect) | ✅ | Authentification |
| OAuth 2.0 | ✅ | Autorisation |
| SAML 2.0 | ✅ | Fédération legacy |
| mTLS | ✅ | Communication sécurisée |
| DID/VC | ✅ | Identité décentralisée |

---

## 3. ARCHITECTURE DE FÉDÉRATION

```
                    ┌──────────────┐
                    │  SNISID IAM  │
                    │  (IdP Nat.)  │
                    └──────┬───────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
     ┌────▼────┐      ┌────▼────┐      ┌────▼────┐
     │Ministère│      │Ministère│      │ Police  │
     │Intérieur│      │  Santé  │      │ Nat.    │
     │  (SP)   │      │  (SP)   │      │  (SP)   │
     └─────────┘      └─────────┘      └─────────┘
          │                │                │
     ┌────▼────┐      ┌────▼────┐      ┌────▼────┐
     │ Justice │      │Finances │      │ Douane  │
     │  (SP)   │      │  (SP)   │      │  (SP)   │
     └─────────┘      └─────────┘      └─────────┘

Protocoles : OIDC / OAuth2 / SAML / mTLS
Trust Anchor : PKI Nationale SNISID
```

---

## 4. IDENTITY BROKER

### Architecture du Broker

```
[OIDC Adapter] [SAML Adapter] [mTLS Adapter] [DID/VC Adapter]
         │            │              │              │
         └────────────┼──────────────┼──────────────┘
                      ↓
         ┌─────────────────────────────┐
         │    PROTOCOL TRANSLATOR      │
         │  OIDC ←→ SAML ←→ mTLS      │
         └──────────────┬──────────────┘
                        ↓
         ┌─────────────────────────────┐
         │       CLAIMS MAPPING        │
         │  NNU → subject              │
         │  Nom → name                 │
         │  Rôle → roles               │
         │  Clearance → clearance      │
         │  Consent → consent_status   │
         └──────────────┬──────────────┘
                        ↓
         ┌─────────────────────────────┐
         │     TRUST VALIDATION        │
         │  - Vérification signature   │
         │  - Vérification certificat  │
         │  - Vérification scope       │
         │  - Vérification expiration  │
         └─────────────────────────────┘
```

---

## 5. REGISTRE DES PARTENAIRES

### Service Providers

```yaml
service_providers:
  - id: "sp:ministere-interieur"
    name: "Ministère de l'Intérieur"
    protocols: [oidc, mtls]
    trust_level: 4
    scopes: [identity_basic, identity_full]
    claims: [nnu, name, role]
  
  - id: "sp:ministere-sante"
    name: "Ministère de la Santé Publique"
    protocols: [oidc, saml]
    trust_level: 3
    scopes: [identity_basic, health_eligibility]
    claims: [nnu, name, health_status]
  
  - id: "sp:police-nationale"
    name: "Police Nationale d'Haïti"
    protocols: [oidc, mtls, saml]
    trust_level: 5
    scopes: [identity_full, biometric_verify, criminal_record]
    claims: [nnu, name, biometric_refs, judicial_status]
  
  - id: "sp:justice"
    name: "Ministère de la Justice"
    protocols: [oidc, mtls]
    trust_level: 5
    scopes: [identity_full, judicial_status, criminal_record]
    claims: [nnu, name, judicial_status, court_records]
```

### Enregistrement d'un Partenaire

```
[Demande d'accès] → [Vérif Sécurité] → [Config Broker] → [Tests Intégration]
                                                                                    ↓
[Production Access] ← [Audit Sécurité] ← [Approbation Comité]
```

---

## 6. GESTION DE CONFIANCE

### PKI Fédérée

```yaml
federated_pki:
  trust_anchor:
    name: SNISID Root CA
    algorithm: ECDSA P-384
    key_size: 384
    validity_years: 10
  
  intermediate_cas:
    - name: SNISID Federation CA (5 ans)
    - name: Agency CA - Ministère Intérieur (3 ans)
    - name: Agency CA - Police Nationale (3 ans)
    - name: Agency CA - Ministère Justice (3 ans)
  
  certificate_lifecycle:
    issuance: automated_via_api
    renewal: 30_days_before_expiry
    revocation: immediate_via_crl
    validation: online_ocsp + crl
```

### Mutual TLS

```yaml
mtls_configuration:
  enabled: true
  client_cert_required: true
  server_cert_required: true
  verify_depth: 3
  crl_check: true
  ocsp_stapling: true
  tls_versions: [TLSv1.3, TLSv1.2_temporary]
```

---

## 7. SCOPES ET CLAIMS

### Scopes Standardisés

| Scope | Description | Accès |
|-------|-------------|-------|
| `identity_basic` | NNU + nom + date naissance | Tous les partenaires |
| `identity_full` | Identité complète + biométrie | Partenaires niveau 4+ |
| `biometric_verify` | Vérification biométrique | Partenaires niveau 5 |
| `judicial_status` | Statut judiciaire | Justice, Police |
| `health_eligibility` | Éligibilité santé | Ministère Santé |
| `criminal_record` | Casier judiciaire | Justice, Police |

### Claims Standardisés

```yaml
standard_claims:
  identity: [nnu, given_name, family_name, date_of_birth, place_of_birth, nationality, sex]
  security: [mfa_verified, role, clearance, session_id]
  biometric: [biometric_refs, biometric_status]
  judicial: [judicial_status, restrictions, court_reference]
  consent: [granted_scopes, revoked_scopes, last_consent_date]
```

---

## 8. SÉCURITÉ FÉDÉRÉE

| Menace | Mitigation |
|--------|-----------|
| Token volé | mTLS + courte validité (5 min) |
| SP compromis | Révocation immédiate du certificat |
| Replay attack | Nonce + timestamp |
| Man-in-the-middle | mTLS + signature |
| Claim manipulation | Signature JWT + validation |
| Usurpation identité | MFA + biométrie |

### Surveillance Fédérée

```yaml
federation_monitoring:
  metrics: [auth_requests_per_second, failed_authentications, token_issuance_rate, cert_expiry_warnings, trust_level_changes, scope_violations]
  alerts:
    - failed_auth_threshold: 10_per_minute
    - token_reuse_detected
    - certificate_revoked
    - trust_level_downgrade
    - unauthorized_scope_request
  audit: [all_authentication_events, all_token_issuance, all_certificate_operations, all_trust_changes]
```

---

> **Toutes les agences partagent la même confiance identité via la fédération SNISID.**
