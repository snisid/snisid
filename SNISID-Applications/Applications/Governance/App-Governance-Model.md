# APP GOVERNANCE MODEL — SNISID
## Modèle de Gouvernance Applicative Nationale

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-GOV-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |

---

## 1. PRÉSENTATION

Modèle de gouvernance garantissant que chaque application SNISID passe par un processus de validation officiel avant d'être déployée auprès des citoyens haïtiens.

### 1.1 Principe Fondamental

> **Aucune application n'est déployée sans validation officielle.**

---

## 2. PÉRIMÈTRE DE GOUVERNANCE

### 2.1 Domaine

| Domaine | Couvert |
|---------|---------|
| Release Governance | ✅ |
| Security Reviews | ✅ |
| UX Standards | ✅ |
| Accessibility Validation | ✅ |
| Offline Certification | ✅ |
| Performance Baseline | ✅ |
| Privacy Impact Assessment | ✅ |
| Legal Compliance | ✅ |

### 2.2 Applications Concernées

- Citizen Super App
- Government Super App
- Police & Justice Apps
- Mobile Field Apps
- Identity Wallet App
- National Admin Portal
- Notification Platform
- Mobile Sync Engine
- Toute application tierce intégrée

---

## 3. PROCESSUS DE VALIDATION

### 3.1 Cycle de Validation

```
┌──────────────────────────────────────────────┐
│        APPLICATION GOVERNANCE CYCLE          │
├──────────────────────────────────────────────┤
│                                              │
│  ┌──────────┐                               │
│  │  Design  │───▶ Architecture Review        │
│  └──────────┘       (Architecte National)    │
│                                              │
│  ┌──────────┐                               │
│  │  Review  │───▶ Security Review            │
│  └──────────┘       (CSO)                    │
│                                              │
│  ┌──────────┐                               │
│  │  Build   │───▶ UX Validation              │
│  └──────────┘       (Design Authority)       │
│                                              │
│  ┌──────────┐                               │
│  │  Pre-    │───▶ Offline Certification      │
│  │  Deploy  │       (Test Lab)               │
│  └──────────┘                               │
│                                              │
│  ┌──────────┐                               │
│  │  Deploy  │───▶ Release Governance         │
│  └──────────┘       (Change Board)           │
│                                              │
│  ┌──────────┐                               │
│  │ Monitor  │───▶ Performance & Security     │
│  └──────────┘       (Ongoing)                │
│                                              │
└──────────────────────────────────────────────┘
```

### 3.2 Release Governance

| Phase | Activité | Responsable | Durée Max |
|-------|----------|-------------|-----------|
| **Design Review** | Architecture, Data flow, Security design | Architecte National | 5 jours |
| **Security Review** | Threat model, Code audit, Pentest | CSO | 10 jours |
| **UX Validation** | Design review, Accessibility test, User test | Design Authority | 5 jours |
| **Offline Cert** | Offline test, Sync test, Battery test | Test Lab | 5 jours |
| **Performance** | Load test, Stress test, Baseline | Performance Team | 3 jours |
| **Legal** | Privacy, Compliance, Data protection | Legal Officer | 5 jours |
| **Change Board** | Release approval, Rollback plan | CAB | 2 jours |

### 3.3 Release Types

| Type | Validation | Délai | Risque |
|------|-----------|-------|--------|
| **Major** (v1.0, v2.0) | Complète (toutes étapes) | 30 jours | Élevé |
| **Minor** (v1.1, v1.2) | Sécurité + Offline + Board | 10 jours | Moyen |
| **Patch** (v1.0.1) | Sécurité + Board | 3 jours | Faible |
| **Hotfix** (urgence) | Board accéléré | 24h | Critique |

---

## 4. SECURITY REVIEWS

### 4.1 Processus

```
┌──────────────────────────────────────────────┐
│         SECURITY REVIEW PROCESS              │
├──────────────────────────────────────────────┤
│                                              │
│  Phase 1: Threat Modeling                    │
│  • STRIDE analysis                          │
│  • Attack tree construction                 │
│  • Risk assessment                          │
│                                              │
│  Phase 2: Code Review                        │
│  • Static analysis (SAST)                   │
│  • Dependency scanning                      │
│  • Secret scanning                          │
│                                              │
│  Phase 3: Dynamic Testing                    │
│  • Penetration testing                      │
│  • API security testing                     │
│  • Session management test                  │
│                                              │
│  Phase 4: Mobile Specific                    │
│  • Anti-tampering verification              │
│  • Storage encryption audit                 │
│  • Offline security validation              │
│                                              │
│  Phase 5: Report                             │
│  • Findings classification                  │
│  • Remediation plan                         │
│  • Sign-off / Rejection                     │
│                                              │
└──────────────────────────────────────────────┘
```

### 4.2 Security Classification

| Classe | Définition | Review | Fréquence |
|--------|-----------|--------|-----------|
| **Standard** | Données publiques | SAST + DAST | Par release |
| **Sensible** | Données personnelles | SAST + DAST + Pentest | Par release |
| **Critique** | Identité, biométrie | SAST + DAST + Pentest + Audit | Par release + trimestriel |

---

## 5. UX STANDARDS

### 5.1 Validation Checklist

| Critère | Standard | Vérification |
|---------|----------|--------------|
| Mobile-first | UI priorise mobile | Design review |
| Accessibility | WCAG 2.1 AA | Automated + Manual |
| Multilingue | FR + HT + EN | Translation audit |
| Offline UX | État offline clair | Offline test |
| Security UX | MFA fluide | User testing |
| Loading states | Skeleton screens | Visual audit |
| Error handling | Messages clairs | Test scenarios |
| Consistency | Design system | Visual audit |

### 5.2 UX Sign-off

```
┌──────────────────────────────────────────────┐
│          UX VALIDATION REPORT                │
├──────────────────────────────────────────────┤
│  Application: Citizen Super App v1.0         │
│  Validator: National Design Authority        │
│  Date: 2026-05-25                           │
├──────────────────────────────────────────────┤
│  ✅ National design system compliance        │
│  ✅ Accessibility WCAG 2.1 AA verified       │
│  ✅ Mobile-first experience validated        │
│  ✅ Offline UX pattern implemented           │
│  ✅ Multilingual (FR/HT/EN) verified         │
│  ✅ Security UX integration reviewed         │
│  ⚠  High-contrast mode optimization          │
│  ⚠  Screen reader testing partial           │
├──────────────────────────────────────────────┤
│  Status: APPROUVÉ AVEC RÉSERVES             │
│  Date limite corrections: 2026-06-01        │
└──────────────────────────────────────────────┘
```

---

## 6. ACCESSIBILITY VALIDATION

### 6.1 Test Coverage

| Test | Outil | Seuil |
|------|-------|-------|
| Color Contrast | axe DevTools | > 4.5:1 |
| Screen Reader | TalkBack / VoiceOver | 100% labels |
| Touch Targets | Manual | ≥ 48px |
| Keyboard Nav | Manual | 100% reachable |
| Focus Order | Manual | Logical order |
| Alt Text | axe DevTools | 100% images |
| Form Labels | axe DevTools | 100% inputs |

### 6.2 Accessibility Pass Criteria

```
┌──────────────────────────────────────────────┐
│    ACCESSIBILITY PASS CRITERIA               │
├──────────────────────────────────────────────┤
│  Critical: 0 blockers                        │
│  High: < 5 violations                        │
│  Medium: < 10 violations                     │
│  Low: < 20 violations                        │
│                                              │
│  Total score: > 90/100                       │
└──────────────────────────────────────────────┘
```

---

## 7. OFFLINE CERTIFICATION

### 7.1 Test Scenarios

| Test | Condition | Pass |
|------|-----------|------|
| No network | App launch and core functions | ✅ |
| Airplane mode | Full workflow offline | ✅ |
| Intermittent network | Sync stability | ✅ |
| No network 7 days | Data survival | ✅ |
| Storage full | Graceful handling | ✅ |
| Low battery | Priority operations | ✅ |
| Concurrent sync | Conflict resolution | ✅ |
| Data corruption | Integrity check | ✅ |
| Power loss mid-sync | Atomic rollback | ✅ |
| Network resume | Sync continuation | ✅ |

### 7.2 Offline Certification Levels

| Niveau | Description | Requis |
|--------|-------------|--------|
| **Bronze** | Basic offline read | Consultation seule |
| **Silver** | Offline read + write | Pending sync queue |
| **Gold** | Full offline operations | Complete workflows offline |
| **Platinum** | Air-gapped operations | No network required ever |

---

## 8. RELEASE BOARD

### 8.1 Membres

| Rôle | Nom (exemple) | Décision |
|------|---------------|----------|
| Architecte National | DNI | Approbation architecture |
| CSO | DNI | Approbation sécurité |
| UX Director | DNI | Approbation UX |
| Test Manager | DNI | Approbation qualité |
| Product Owner | DNI | Approbation fonctionnelle |
| CIO | DNI | Décision finale |

### 8.2 Décisions

| Décision | Condition |
|----------|-----------|
| **Approuvé** | Tous les critères remplis |
| **Approuvé avec réserves** | Critères mineurs non remplis, date limite |
| **Rejeté** | Critères majeurs non remplis |
| **Reporté** | Dépendances non résolues |

---

## 9. COMPLIANCE

| Standard | Validation |
|----------|-----------|
| OWASP MASVS | Security review |
| WCAG 2.1 | Accessibility audit |
| RGPD / Loi Haïtienne | Privacy assessment |
| ISO 27001 | ISMS compliance |
| National Identity Law | Legal review |

---
*Fin du document — App Governance Model v1.0*