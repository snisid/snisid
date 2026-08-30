# Consent Engine Specification

> **Moteur de Gouvernance Citoyenne — SNISID**  
> **Version :** 1.0.0  
> **Classification :** SOUVERAIN — RESTREINT  
> **Dernière mise à jour :** 2026-05-25

---

## 1. OBJECTIF

Créer la gouvernance citoyenne du SNISID. Le citoyen contrôle ses données personnelles.

> **Le consentement doit être explicite et traçable.**

---

## 2. FONCTIONS CITOYEN

| Fonction | Support | Description |
|----------|---------|-------------|
| Voir accès données | ✅ | Historique complet des accès |
| Autoriser partage | ✅ | Consentement explicite par usage |
| Révoquer consentement | ✅ | Révocation immédiate |
| Voir historique | ✅ | Journal de tous les consentements |
| Recevoir alertes | ✅ | Notification en temps réel |

---

## 3. ARCHITECTURE DU CONSENT ENGINE

```
[Consent Request Handler] → [Consent Manager] → [Consent Auditor]
         ↓                       ↓                      ↓
   ┌─────────────────────────────────────────────────────────┐
   │                    CONSENT REGISTRY                     │
   │  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
   │  │ Active   │  │ Revoked  │  │ Expired  │             │
   │  │ Consents │  │ Consents │  │ Consents │             │
   │  └──────────┘  └──────────┘  └──────────┘             │
   └────────────────────────┬───────────────────────────────┘
                            ↓
   ┌─────────────────────────────────────────────────────────┐
   │                 CONSENT POLICY ENGINE                   │
   │  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
   │  │ Purpose  │  │ Scope    │  │ Duration │             │
   │  │ Checker  │  │ Checker  │  │ Checker  │             │
   │  └──────────┘  └──────────┘  └──────────┘             │
   └────────────────────────┬───────────────────────────────┘
                            ↓
   ┌─────────────────────────────────────────────────────────┐
   │                  NOTIFICATION SYSTEM                    │
   │  SMS | Email | Wallet Push                              │
   └─────────────────────────────────────────────────────────┘
```

---

## 4. MODÈLE DE CONSENTEMENT

```json
{
  "consent_id": "uuid-v4",
  "citizen_nnu": "NNU-HA-20260525-A7K9M2X4-73-1",
  "requester": {"type": "organization", "id": "org:ministere-sante", "name": "Ministère de la Santé"},
  "purpose": "Vérification d'éligibilité aux services de santé",
  "data_scope": ["identity_basic", "biometric_face", "health_eligibility"],
  "legal_basis": "article_12_loi_protection_donnees",
  "granted": true,
  "granted_at": "2026-05-25T10:00:00Z",
  "granted_method": "wallet_biometric",
  "valid_from": "2026-05-25T10:00:00Z",
  "valid_until": "2026-06-25T10:00:00Z",
  "revoked": false,
  "version": 1,
  "proof": {"type": "RsaSignature2018", "verification_method": "did:snisid:citizen#keys-1"}
}
```

---

## 5. POLITIQUES DE CONSENTEMENT

```yaml
consent_policies:
  explicit_consent:
    description: Consentement explicite requis pour toute utilisation de données
    applies_to: [biometric, health, financial, judicial]
    method: wallet_signature OR biometric_confirmation
  
  implicit_consent:
    description: Consentement implicite pour services essentiels
    applies_to: [identity_verification, basic_services]
    method: authentication_success
  
  purpose_limitation:
    description: Les données ne peuvent être utilisées que pour la finalité déclarée
    enforcement: strict
    violation_action: deny_access + alert
  
  data_minimization:
    description: Seules les données nécessaires sont partagées
    enforcement: automatic
    scope_filtering: true
  
  time_limitation:
    description: Le consentement expire après une durée définie
    default_duration: 30_days
    max_duration: 1_year
    renewal_required: true
  
  revocability:
    description: Le citoyen peut révoquer son consentement à tout moment
    effective: immediate
    notification: automatic
    audit: mandatory
```

---

## 6. WORKFLOW DE CONSENTEMENT

### Demande
```
[Demandeur (Org)] → [Consent Engine] → [Citoyen (Wallet)]
                                            ↓
                              ┌─────────────┼─────────────┐
                              ↓             ↓             ↓
                         [Consentir]   [Refuser]   [Demander infos]
                              ↓             ↓             ↓
                         [Enregistrer] [Archiver]  [Fournir infos]
                              ↓             ↓             ↓
                         [Notifier]    [Notifier]  [Re-demander]
```

### Révocation
```
[Citoyen révoque] → [Wallet signe] → [Consent Engine]
                                          ↓
                              ┌───────────┼───────────┐
                              ↓           ↓           ↓
                        [Mettre à    [Notifier   [Bloquer
                         jour Reg.]  Demandeur]   Accès]
                              ↓           ↓           ↓
                        [Archiver]  [Confirmer] [Journaliser]
                          Révocation  Citoyen     Audit
```

---

## 7. API DE CONSENTEMENT

| Endpoint | Méthode | Description | Auth |
|----------|---------|-------------|------|
| `/api/v1/consent/request` | POST | Demander un consentement | WalletAuth + BiometricAuth |
| `/api/v1/consent/{id}/grant` | POST | Accorder un consentement | WalletAuth + BiometricAuth |
| `/api/v1/consent/{id}/revoke` | POST | Révoquer un consentement | WalletAuth |
| `/api/v1/consent/history` | GET | Historique des consentements | WalletAuth |
| `/api/v1/consent/alerts` | GET | Alertes de consentement | WalletAuth |

---

## 8. CONFORMITÉ ET GOUVERNANCE

| Principe | Application |
|----------|-------------|
| Explicite | Signature ou confirmation biométrique requise |
| Traçable | Chaque action est journalisée et signée |
| Révocable | Révocation immédiate à tout moment |
| Limité | Portée et durée définies |
| Transparent | Citoyen informé de chaque usage |
| Vérifiable | Preuve cryptographique de consentement |

### Audit
```yaml
consent_audit:
  logged_events: [consent_requested, consent_granted, consent_revoked, consent_expired, consent_violation, access_granted, access_denied]
  retention_period: 10_years
  storage: immutable_ledger
  verification: cryptographic_signature
  access: [citizen, audit_officer, judicial]
  reporting:
    frequency: monthly
    recipients: [citizen, audit_officer, governance_committee]
```

---

> **Le consentement citoyen est explicite, traçable et révocable.**
