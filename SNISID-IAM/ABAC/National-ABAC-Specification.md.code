# National ABAC Specification

> **Contrôle d'Accès Basé sur les Attributs — National SNISID**  
> **Version :** 1.0.0  
> **Classification :** SOUVERAIN — RESTREINT  
> **Dernière mise à jour :** 2026-05-25

---

## 1. OBJECTIF

Créer un contrôle contextuel intelligent et adaptatif. Le système évalue en temps réel les attributs de la requête pour autoriser ou refuser l'accès.

> **Le système doit devenir adaptatif.**

---

## 2. MODÈLE ABAC

### 2.1 Attributs Évalués

#### Subject (Qui demande)
| Attribut | Source | Type | Exemple |
|----------|--------|------|---------|
| `subject.nnu` | IAM Registry | String | `NNU-HA-20260525-A7K9M2X4-73-1` |
| `subject.role` | RBAC Engine | String[] | `[enrollment_officer]` |
| `subject.mfa_verified` | MFA Engine | Boolean | `true` |
| `subject.status` | IAM Registry | String | `active` |
| `subject.clearance` | HR System | Integer | `3` |

#### Resource (Ce qui est demandé)
| Attribut | Source | Type | Exemple |
|----------|--------|------|---------|
| `resource.type` | Resource Registry | String | `citizen_identity` |
| `resource.classification` | Data Policy | String | `confidential` |
| `resource.sensitivity` | Data Policy | Integer | `4` |
| `resource.tags` | Resource Registry | String[] | `[biometric, personal]` |

#### Environment (Contexte)
| Attribut | Source | Type | Exemple |
|----------|--------|------|---------|
| `environment.time` | System Clock | Time | `14:30` |
| `environment.location` | GPS/IP | GeoPoint | `18.5944,-72.3074` |
| `environment.ip` | Network | String | `192.168.1.100` |
| `environment.device` | Device Registry | String | `trusted_terminal_42` |
| `environment.risk_score` | Risk Engine | Float | `0.23` |

---

## 3. RÈGLES ABAC

### POL-001 : Localisation invalide
```yaml
policy_id: POL-001
name: "Invalid Location Access"
effect: deny
conditions:
  - attribute: environment.location
    operator: not_in
    value: authorized_zones
```

### POL-002 : Heure inhabituelle
```yaml
policy_id: POL-002
name: "After Hours Access"
effect: conditional
conditions:
  - attribute: environment.time
    operator: not_in_range
    value: ["06:00", "20:00"]
obligations:
  - require_mfa: true
  - log_access: true
```

### POL-003 : Appareil non reconnu
```yaml
policy_id: POL-003
name: "Unrecognized Device"
effect: conditional
conditions:
  - attribute: environment.device
    operator: not_in
    value: registered_devices
obligations:
  - require_mfa: true
  - require_additional_verification: true
```

### POL-004 : Risque élevé
```yaml
policy_id: POL-004
name: "High Risk Access"
effect: deny
conditions:
  - attribute: environment.risk_score
    operator: greater_than
    value: 0.8
obligations:
  - block_access: true
  - alert_security: true
```

### POL-005 : Données biométriques
```yaml
policy_id: POL-005
name: "Biometric Data Access"
effect: conditional
conditions:
  - attribute: resource.tags
    operator: contains
    value: biometric
obligations:
  - require_role: [enrollment_officer, identity_analyst, fraud_investigator, iam_admin]
  - require_mfa: true
  - require_clearance_min: 3
```

---

## 4. MATRICE DE DÉCISION

| Condition | Action | Rôle requis | MFA |
|-----------|--------|-------------|-----|
| Localisation invalide | Refuser | — | — |
| Heure inhabituelle | MFA obligatoire | Tous | ✅ |
| Appareil inconnu | Restriction | — | ✅ |
| Risque élevé | Blocage | — | — |
| Données sensibles | Clearance minimum | — | ✅ |
| Accès hors réseau | VPN obligatoire | — | ✅ |
| Session expirée | Ré-authentification | — | ✅ |

---

## 5. MOTEUR DE DÉCISION

### 5.1 Algorithme d'Évaluation

```
[Requête d'accès] → [Collecte Attributs] → [Évaluation Politiques] → [Combinaison Décisions] → [Décision Finale]
```

### 5.2 Règles de Combinaison

```yaml
combination_rules:
  default: deny-overrides
  priority_order:
    1. judicial_restriction     # Restrictions judiciaires (toujours deny)
    2. security_threat          # Menaces sécurité (toujours deny)
    3. biometric_protection     # Protection biométrique
    4. data_sensitivity         # Sensibilité des données
    5. location_policy          # Politique de localisation
    6. time_policy              # Politique temporelle
    7. device_policy            # Politique d'appareil
    8. role_policy              # Politique de rôle
    9. default_policy           # Politique par défaut (deny)
```

---

## 6. ADAPTATIVITÉ

### 6.1 Ajustement Dynamique

```yaml
adaptive_behavior:
  risk_based_adjustment:
    low_risk:     {threshold: "0.0-0.3", action: "standard_access", mfa: "not_required"}
    medium_risk:  {threshold: "0.3-0.6", action: "enhanced_verification", mfa: "required"}
    high_risk:    {threshold: "0.6-0.8", action: "restricted_access", mfa: "required_plus", approval: "supervisor_required"}
    critical_risk:{threshold: "0.8-1.0", action: "access_denied", alert: "security_team", investigation: "automatic"}
  
  context_awareness:
    learn_patterns: true
    adjust_thresholds: true
    feedback_loop: true
    review_period: monthly
```

---

## 7. POLITIQUES SPÉCIALES

### Politique d'Urgence
```yaml
emergency_policy:
  activation:
    trigger: security_incident OR system_compromise
    authorized_by: [iam_admin, cyber_analyst]
  effects:
    - all_sessions_terminated: true
    - mfa_required_for_all: true
    - location_restriction: government_network_only
    - time_restriction: business_hours_only
  duration:
    max_emergency_hours: 72
    renewal_requires: committee_approval
```

---

## 8. CONFORMITÉ ABAC

| Exigence | Implémentation | Vérification |
|----------|---------------|-------------|
| Adaptativité | Risk-based adjustment | Monitoring continu |
| Contextualité | 4 types d'attributs | Audit trimestriel |
| Décision traçable | Logs complets | Revue mensuelle |
| Politiques signées | Signature PKI | Vérification automatique |
| Règles versionnées | Git-controlled | Audit de configuration |

---

> **Le ABAC national rend le système adaptatif, contextuel et intelligent.**
