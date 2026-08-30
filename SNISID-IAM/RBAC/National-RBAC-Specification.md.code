# National RBAC Specification

> **Contrôle d'Accès Basé sur les Rôles — National SNISID**  
> **Version :** 1.0.0  
> **Classification :** SOUVERAIN — RESTREINT  
> **Dernière mise à jour :** 2026-05-25

---

## 1. OBJECTIF

Contrôler les accès gouvernementaux selon le **principe du moindre privilège**. Chaque accès est :
- **Minimal** — Juste ce qui est nécessaire
- **Temporaire** — Expirable si nécessaire
- **Auditable** — Traçable et journalisé

---

## 2. HIÉRARCHIE DES RÔLES

| Rôle | Niveau | Description | MFA | ABAC | Session max |
|------|--------|-------------|-----|------|-------------|
| `citizen` | Faible | Accès à ses propres données | ✅ | ✅ | 30 min |
| `enrollment_officer` | Moyen | Création et modification identité | ✅ | ✅ | 8h |
| `identity_analyst` | Moyen | Investigation et corrélation | ✅ | ✅ | 8h |
| `judge` | Élevé | Décisions judiciaires identité | ✅ | ✅ | 4h |
| `police_officer` | Élevé | Vérification et investigation | ✅ | ✅ | 4h |
| `cyber_analyst` | Critique | Surveillance et réponse sécurité | ✅ | ✅ | 2h |
| `pki_admin` | Critique | Certificats et clés | ✅ | ✅ | 1h |
| `iam_admin` | Critique | Gestion complète IAM | ✅ | ✅ | 1h |
| `audit_officer` | Élevé | Consultation logs et rapports | ✅ | ✅ | 8h |
| `fraud_investigator` | Critique | Enquête et action fraude | ✅ | ✅ | 4h |

### Niveaux d'accès

```
CRITIQUE ██████████████████████████████████████
  ├── cyber_analyst, pki_admin, iam_admin, sys_admin, fraud_investigator

ÉLEVÉ  ████████████████████████
  ├── judge, police_officer, audit_officer

MOYEN  ████████████
  ├── enrollment_officer, identity_analyst

FAIBLE ███
  └── citizen
```

---

## 3. PERMISSIONS PAR RÔLE

### citizen
| Ressource | read | write | delete | admin |
|-----------|------|-------|--------|-------|
| Ses propres données | ✅ | ❌ | ❌ | ❌ |
| Son wallet | ✅ | ✅ (profil) | ❌ | ❌ |
| Ses consentements | ✅ | ✅ | ✅ (révoquer) | ❌ |
| Son historique d'accès | ✅ | ❌ | ❌ | ❌ |
| Données d'autrui | ❌ | ❌ | ❌ | ❌ |

### enrollment_officer
| Ressource | read | write | delete | admin |
|-----------|------|-------|--------|-------|
| Données enrôlement | ✅ | ✅ | ❌ | ❌ |
| Registre NNU | ✅ (création) | ✅ | ❌ | ❌ |
| Biométrie | ✅ (capture) | ✅ | ❌ | ❌ |
| Identités existantes | ✅ | ✅ (correction) | ❌ | ❌ |

### judge
| Ressource | read | write | delete | admin |
|-----------|------|-------|--------|-------|
| Identités | ✅ | ❌ | ❌ | ❌ |
| Statut judiciaire | ✅ | ✅ | ❌ | ❌ |
| Restrictions | ✅ | ✅ | ✅ | ❌ |
| NNU | ✅ | ✅ (suspension/révocation) | ❌ | ❌ |

### iam_admin
| Ressource | read | write | delete | admin |
|-----------|------|-------|--------|-------|
| Rôles | ✅ | ✅ | ❌ | ✅ |
| Permissions | ✅ | ✅ | ✅ | ✅ |
| Utilisateurs | ✅ | ✅ | ✅ | ✅ |
| Politiques | ✅ | ✅ | ✅ | ✅ |
| Configuration MFA | ✅ | ✅ | ❌ | ✅ |
| Federation | ✅ | ✅ | ✅ | ✅ |

---

## 4. SÉPARATION DES TÂCHES

| Règle | Application |
|-------|-------------|
| Dual control | Les rôles critiques nécessitent 2 personnes |
| Rotation | Les administrateurs PKI/IAM tournent tous les 90 jours |
| Vacances obligatoires | 5 jours consécutifs minimum par an |
| Audit indépendant | L'audit officer est séparé des opérations |

---

## 5. CONTRAINTES TEMPorelles

```yaml
temporal_constraints:
  citizen:           {access_hours: "00:00-23:59", max_session: 30}
  enrollment_officer:{access_hours: "06:00-20:00", max_session: 480}
  judge:             {access_hours: "07:00-19:00", max_session: 240, weekdays_only: true}
  police_officer:    {access_hours: "00:00-23:59", max_session: 240}
  cyber_analyst:     {access_hours: "00:00-23:59", max_session: 120}
  pki_admin:         {access_hours: "08:00-17:00", max_session: 60, weekdays_only: true}
  iam_admin:         {access_hours: "08:00-17:00", max_session: 60, weekdays_only: true}
  audit_officer:     {access_hours: "07:00-19:00", max_session: 480, weekdays_only: true}
  fraud_investigator:{access_hours: "00:00-23:59", max_session: 240}
```

---

## 6. GESTION DES RÔLES

### Cycle de vie
```
[Création (Admin)] → [Attribution (HR+Admin)] → [Activation (MFA)] → [Utilisation (Session)]
                                                                                           ↓
[Archive (Audit)] ← [Révocation (Admin)] ← [Suspension (Auto)]
```

### Attribution
```yaml
role_assignment:
  process:
    1. Demande HR → validation fonction
    2. Vérification identité → MFA
    3. Attribution rôle → par IAM admin
    4. Activation → après MFA
    5. Notification → à l'utilisateur et audit
  requirements:
    - identity_verified: true
    - mfa_enabled: true
    - background_check: true (rôles élevés et critiques)
    - approval: HR + IAM Admin
    - audit_logged: true
```

### Révocation
```yaml
role_revocation:
  triggers: [departure, rotation, security_incident, policy_violation, manual]
  process:
    1. Déclenchement → automatique ou manuel
    2. Suspension immédiate → session en cours terminée
    3. Révocation → rôle retiré
    4. Nettoyage → certificats, tokens, sessions
    5. Audit → journalisation complète
  sla:
    immediate: rôles critiques
    within_1h: rôles élevés
    within_4h: rôles moyens
```

---

## 7. CONFORMITÉ RBAC

| Exigence | Implémentation | Vérification |
|----------|---------------|-------------|
| Moindre privilège | Rôles minimalistes | Audit trimestriel |
| Séparation tâches | Dual control | Revue mensuelle |
| Rotation | 90 jours pour critiques | Automatique |
| Auditabilité | Logs complets | Monitoring continu |
| Révocabilité | Suspension/révocation | Test semestriel |
| Temporalité | Contraintes horaires | Enforcement système |

---

> **Le RBAC national garantit que chaque accès est minimal, temporaire et auditable.**
