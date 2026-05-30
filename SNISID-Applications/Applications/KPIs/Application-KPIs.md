# APPLICATION KPIs — SNISID
## Indicateurs de Performance Applicative Nationale

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-KPI-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |

---

## 1. PRÉSENTATION

KPIs pour mesurer, piloter et améliorer l'expérience numérique nationale. Toutes les applications SNISID sont pilotées par les données.

---

## 2. KPI PRINCIPAUX

### 2.1 App Uptime

| KPI | Cible | Mesure | Source |
|-----|-------|--------|--------|
| **App Uptime** | > 99.9% | Crash-free sessions / Total sessions | Crashlytics |
| **API Uptime** | > 99.95% | HTTP 200 / Total requests | Prometheus |
| **Service Availability** | > 99.99% | Service up / Total time | Prometheus |
| **Kiosk Uptime** | > 99% | Kiosk operational hours | Kiosk monitoring |

**Formule :**
```
App Uptime = (1 - (Downtime in minutes / Total minutes in period)) × 100
```

### 2.2 Offline Success Rate

| KPI | Cible | Mesure |
|-----|-------|--------|
| **Offline Success Rate** | > 98% | Successful offline ops / Total offline ops |
| **Sync Success Rate** | > 99% | Successful syncs / Total sync attempts |
| **Cache Hit Rate** | > 95% | Cache hits / Total reads |
| **Conflict Rate** | < 1% | Conflicts / Total syncs |

**Formule :**
```
Offline Success Rate = (Successful offline operations / Total offline operations) × 100
```

### 2.3 Crash Rate

| KPI | Cible | Seuil |
|-----|-------|-------|
| **Android Crash Rate** | < 0.1% | > 0.3% |
| **iOS Crash Rate** | < 0.1% | > 0.3% |
| **ANR Rate** | < 0.05% | > 0.1% |
| **OOM Rate** | < 0.01% | > 0.05% |
| **Crash-free Users** | > 99.5% | < 99% |

**Formule :**
```
Crash Rate = (Number of crash sessions / Total sessions) × 100
```

### 2.4 User Satisfaction

| KPI | Cible | Mesure | Fréquence |
|-----|-------|--------|-----------|
| **App Store Rating** | > 4.0/5 | Average rating | Mensuel |
| **NPS** | > 50 | Survey score | Trimestriel |
| **CSAT** | > 85% | Satisfaction survey | Mensuel |
| **Task Success Rate** | > 90% | Completed tasks / Started tasks | Continu |
| **User Retention (D1/D7/D30)** | > 80%/>60%/>40% | Users returning | Continu |

**Formule NPS :**
```
NPS = % Promoters (9-10) - % Detractors (0-6)
```

### 2.5 Sync Reliability

| KPI | Cible | Mesure |
|-----|-------|--------|
| **Sync Reliability** | > 99% | Successful syncs / Total syncs |
| **P0 Sync Latency** | < 1 min | Time from creation to server ACK |
| **P1 Sync Latency** | < 1h | Time from creation to server ACK |
| **Max Offline Duration** | > 7 days | Data survival without network |
| **Data Integrity** | 100% | Checksum verification |

---

## 3. KPI PAR APPLICATION

### 3.1 Citizen Super App

| KPI | Cible | Période |
|-----|-------|---------|
| DAU (Daily Active Users) | > 500,000 | Journalier |
| MAU (Monthly Active Users) | > 5,000,000 | Mensuel |
| ID Verification Success | > 99% | Continu |
| QR Verification Count | > 100,000/jour | Journalier |
| Birth Certificate Access | > 50,000/jour | Journalier |
| Civil Registry Requests | > 10,000/jour | Journalier |
| App Open Rate | > 3x/semaine | Hebdomadaire |
| Session Duration | > 5 min | Continu |

### 3.2 Government Super App

| KPI | Cible | Période |
|-----|-------|---------|
| DAU | > 20,000 | Journalier |
| Approvals Processed | > 5,000/jour | Journalier |
| Case Resolution Time | < 48h | Continu |
| Secure Messages Sent | > 10,000/jour | Journalier |
| Escalations Handled | < 1h P0 | Continu |
| MFA Success Rate | > 99.5% | Continu |

### 3.3 Police & Justice Apps

| KPI | Cible | Période |
|-----|-------|---------|
| Cases Created | > 500/jour | Journalier |
| Biometric Matches | > 1,000/jour | Journalier |
| Watchlist Hits | > 100/jour | Journalier |
| Case Closure Rate | > 80% within 30 days | Mensuel |
| Border Verifications | > 10,000/jour | Journalier |

### 3.4 Mobile Field Apps

| KPI | Cible | Période |
|-----|-------|---------|
| Enrollment/Day/Agent | > 50 | Journalier |
| Offline Operations | > 90% of total | Continu |
| Sync Success Rate | > 99% | Continu |
| Max Offline Days | > 7 jours | Continu |
| Biometric Quality | > 95% pass rate | Continu |

### 3.5 Identity Wallet

| KPI | Cible | Période |
|-----|-------|---------|
| Wallet Activations | > 10,000/jour | Journalier |
| QR Verifications | > 200,000/jour | Journalier |
| Digital Signatures | > 5,000/jour | Journalier |
| Consent Grants | > 1,000/jour | Journalier |
| Wallet Recovery Success | > 95% | Continu |

---

## 4. TABLEAU DE BORD KPI

```
┌──────────────────────────────────────────────────────────┐
│  🇭🇹 SNISID Performance Dashboard │ 🟢 98.5% Overall     │
├──────────────────────────────────────────────────────────┤
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐   │
│  │ App      │ │  Offline  │ │  Crash   │ │  User    │   │
│  │ Uptime   │ │  Success  │ │   Rate   │ │  Sat.   │   │
│  │ 99.92%   │ │  98.7%   │ │  0.08%   │ │  4.2/5  │   │
│  │ 🟢       │ │  🟢      │ │  🟢      │ │  🟢     │   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘   │
│                                                          │
│  ⚠ KPIs Under Threshold:                                │
│  • Citizen App DAU: 423,000 (target: 500,000)           │
│  • Sync Reliability: 98.2% (target: 99%)                │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │  DAU Trend — Last 30 Days                        │   │
│  │  ████████████████████████████████████░░ 423K    │   │
│  └──────────────────────────────────────────────────┘   │
│                                                          │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                │
│  │  Today   │ │  7 Days  │ │  30 Days │                │
│  │  +2.3%   │ │  +5.1%   │ │ +12.4%   │                │
│  └──────────┘ └──────────┘ └──────────┘                │
├──────────────────────────────────────────────────────────┤
│  Last updated: 2026-05-25T14:30:00Z │ Auto-refresh: 5m │
└──────────────────────────────────────────────────────────┘
```

---

## 5. ALERTES KPI

| KPI | Seuil Avertissement | Seuil Critique | Action |
|-----|-------------------|----------------|--------|
| App Uptime | < 99.9% | < 99.5% | Investigation immédiate |
| Offline Success | < 98% | < 95% | Sync engine review |
| Crash Rate | > 0.1% | > 0.5% | Rollback si nécessaire |
| User Satisfaction | < 4.0 | < 3.5 | UX improvement plan |
| Sync Reliability | < 99% | < 95% | Infrastructure review |
| Response Time p95 | > 1s | > 3s | Performance optimization |

---

## 6. REPORTING

| Rapport | Fréquence | Destinataires | Format |
|---------|-----------|---------------|--------|
| **Daily Health** | Quotidien | Ops, Dev | Dashboard |
| **Weekly Report** | Hebdomadaire | Management | PDF + Dashboard |
| **Monthly Review** | Mensuel | DNI, CIO | Presentation |
| **Quarterly Deep** | Trimestriel | Gouvernement | Executive report |
| **Annual Report** | Annuel | Parlement | Public report |

---
*Fin du document — Application KPIs v1.0*