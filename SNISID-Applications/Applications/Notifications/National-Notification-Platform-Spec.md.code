# NATIONAL NOTIFICATION PLATFORM — SNISID
## Communication Gouvernementale Unifiée

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-NOT-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |

---

## 1. PRÉSENTATION

Plateforme centralisée de notification gouvernementale permettant à toutes les applications SNISID d'envoyer des notifications sécurisées via SMS, Email, Push et Messages In-App.

---

## 2. CANAUX DE COMMUNICATION

### 2.1 Canaux Supportés

| Canal | Support | Priorité | Volumétrie |
|-------|---------|----------|------------|
| **SMS** | ✅ | Urgences, Alertes | 1M/mois |
| **Email** | ✅ | Documents, Confirmations | 5M/mois |
| **Push Notification** | ✅ | Temps réel | 10M/mois |
| **Secure In-App Messages** | ✅ | Messages critiques | 2M/mois |

### 2.2 Matrice de Priorité

| Priorité | Canal | Délai | Signature Crypto |
|----------|-------|-------|------------------|
| **P0 — Critique** | SMS + Push + In-App | < 1 min | ✅ Obligatoire |
| **P1 — Haute** | Push + In-App | < 5 min | ✅ Obligatoire |
| **P2 — Normale** | Push + Email | < 1h | ✅ |
| **P3 — Informative** | Email + In-App | < 24h | — |
| **P4 — Promo** | In-App only | Best effort | — |

---

## 3. FONCTIONNALITÉS

### 3.1 Core Features

| Fonction | Support |
|----------|---------|
| Template Management | ✅ |
| Personalization | ✅ |
| Scheduling | ✅ |
| Batch Send | ✅ |
| Delivery Tracking | ✅ |
| Read Receipts | ✅ |
| Opt-in/Opt-out | ✅ |
| Channel Preference | ✅ |
| DLR (Delivery Reports) | ✅ |
| Analytics Dashboard | ✅ |

### 3.2 Cryptographic Signing

Toutes les notifications critiques (P0, P1) sont signées cryptographiquement :

```
┌──────────────────────────────────────────────┐
│         NOTIFICATION SIGNATURE                │
├──────────────────────────────────────────────┤
│  Header:                                     │
│  X-SNISID-Signature: sig_abc123...           │
│  X-SNISID-Timestamp: 2026-05-25T14:30:00Z   │
│  X-SNISID-Nonce: a1b2c3d4e5...              │
│                                               │
│  Body signé : SHA-256(content + timestamp    │
│               + nonce + channel)             │
│                                               │
│  Vérification côté client :                   │
│  1. Vérifier timestamp (max 5 min drift)     │
│  2. Vérifier nonce (anti-replay)             │
│  3. Vérifier signature avec clé publique     │
│  4. Afficher badge "✅ Vérifié SNISID"       │
└──────────────────────────────────────────────┘
```

---

## 4. ARCHITECTURE

```
┌─────────────────────────────────────────────┐
│         NOTIFICATION ORCHESTRATOR            │
├─────────────────────────────────────────────┤
│  ┌──────────────────────────────────────┐   │
│  │  Priority Queue                      │   │
│  │  P0 ──▶ P1 ──▶ P2 ──▶ P3 ──▶ P4    │   │
│  └──────────────────────────────────────┘   │
│                                             │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐    │
│  │ Template │ │ Consent  │ │ Pref.    │    │
│  │ Engine   │ │ Checker  │ │ Manager  │    │
│  └──────────┘ └──────────┘ └──────────┘    │
│                                             │
│  ┌──────────────────────────────────────┐   │
│  │  Channel Router                      │   │
│  │  P0→SMS+Push │ P1→Push │ P2→Push+Email│ │
│  └──────────────────────────────────────┘   │
├─────────────────────────────────────────────┤
│         CHANNEL ADAPTERS                     │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────────┐   │
│  │ SMS  │ │Email │ │ Push │ │ In-App   │   │
│  │ GW   │ │ GW   │ │ GW   │ │ Messages │   │
│  └──────┘ └──────┘ └──────┘ └──────────┘   │
├─────────────────────────────────────────────┤
│         MONITORING & ANALYTICS               │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐    │
│  │ Delivery │ │  Click   │ │  Channel │    │
│  │ Tracking │ │ Tracking │ │  Health  │    │
│  └──────────┘ └──────────┘ └──────────┘    │
└─────────────────────────────────────────────┘
```

---

## 5. SÉCURITÉ

### 5.1 Security Measures

| Measure | Implementation |
|---------|---------------|
| Channel Encryption | TLS 1.3 for all channels |
| Message Signing | ECDSA for critical messages |
| Anti-Spoofing | DKIM (Email), Sender ID (SMS) |
| Rate Limiting | Per user, per channel |
| Consent Management | GDPR-compliant |
| Opt-out | Immediate processing |
| Analytics Privacy | Anonymized aggregation |

### 5.2 Channel Security

| Canal | Chiffrement | Authentification | Anti-Spam |
|-------|-------------|-----------------|-----------|
| SMS | Transport | Sender ID verified | Rate limiting |
| Email | TLS + DKIM | DMARC + SPF | Spam score |
| Push | TLS + Token | Firebase/APNs Auth | Throttling |
| In-App | E2E Encrypted | JWT + Device bind | Server control |

---

## 6. TEMPLATES

### 6.1 Types de Templates

| Type | Variables | Canal |
|------|-----------|-------|
| Identity Verified | {name}, {date}, {id} | Push + In-App |
| Birth Certificate | {name}, {child}, {date} | Email + In-App |
| Appointment Reminder | {service}, {date}, {location} | SMS + Push |
| Sepulture Alert | {name}, {date} | SMS + Push |
| Gov Alert | {title}, {message}, {action} | SMS + Push + In-App |
| Case Update | {case_id}, {status} | Push + In-App |
| Payment Receipt | {amount}, {date}, {ref} | Email + In-App |

### 6.2 Exemple Template (Créole)

```
Sijè: {title}

Bonjour {name},

SNISID enfòme ou: {message}

Aksyon: {action_url}
Dat: {timestamp}

Mesaj sa a verifye pa SNISID
✅ Siyati elektwonik: {signature}
```

---

## 7. PERFORMANCE

| Métrique | Cible |
|----------|-------|
| P0 Delivery (SMS+Push) | < 1 min |
| P1 Delivery (Push) | < 5 min |
| P2 Delivery (Email) | < 1h |
| Platform Uptime | > 99.99% |
| Queue Capacity | > 10M messages |
| Throughput | > 10,000 msg/s |
| Delivery Rate | > 99% (SMS), > 98% (Push) |

---

## 8. INTÉGRATIONS

| Application | Type de Notification | Priorité |
|-------------|---------------------|----------|
| Citizen App | Identity, Civil Status, Payments | P0-P3 |
| Government App | Approvals, Escalations, Cases | P0-P2 |
| Police App | Alerts, Case Assignments | P0-P1 |
| Identity Wallet | Consent, Signature Requests | P0-P1 |
| Admin Portal | System Alerts, Audit Alerts | P0-P1 |

---
*Fin du document — National Notification Platform v1.0*