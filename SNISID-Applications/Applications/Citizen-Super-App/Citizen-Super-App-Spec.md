# CITIZEN SUPER APP — SNISID
## Application Citoyenne Nationale

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-CSA-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |
| Plateforme | Android, iOS, Web |

---

## 1. PRÉSENTATION

### 1.1 Qu'est-ce que le Citizen Super App ?

Le **Citizen Super App** est l'application unique qui permet à chaque citoyen haïtien d'accéder à l'ensemble des services gouvernementaux numériques depuis un point d'entrée unique.

### 1.2 Principes Fondamentaux

- **Unified Access** — Tous les services en une seule app
- **Offline Partiel** — Fonctionne sans internet pour les fonctionnalités essentielles
- **Souverain** — Hébergé en Haïti, contrôle national des données
- **Accessible** — Conçu pour tous les niveaux de littératie numérique
- **Multilingue** — Français, Créole, Anglais

---

## 2. FONCTIONNALITÉS

### 2.1 Identity & Wallet

| Fonction | Support | Offline | Description |
|----------|---------|---------|-------------|
| Digital Identity | ✅ | ✅ (partiel) | Carte d'identité numérique nationale |
| Identity Wallet | ✅ | ✅ | Stockage sécurisé des certificats |
| QR Identity | ✅ | ✅ | QR code d'identité vérifiable |
| Biometric Auth | ✅ | ✅ | Empreinte / Face unlock |
| Digital Signature | ✅ | ❌ | Signature électronique des documents |

### 2.2 État Civil

| Fonction | Support | Offline | Description |
|----------|---------|---------|-------------|
| Birth Certificate | ✅ | ✅ (cache) | Consultation du certificat de naissance |
| Marriage Certificate | ✅ | ✅ (cache) | Consultation du certificat de mariage |
| Death Certificate | ✅ | ✅ (cache) | Consultation du certificat de décès |
| Civil Registry Requests | ✅ | ❌ | Demandes d'actes d'état civil |
| Certificate QR Verify | ✅ | ✅ | Vérification par QR code |

### 2.3 Services Citoyens

| Fonction | Support | Offline |
|----------|---------|---------|
| Notifications | ✅ | ✅ (buffer) |
| Inbox Messages | ✅ | ✅ (cache) |
| Service Requests | ✅ | ❌ |
| Appointment Booking | ✅ | ❌ |
| Payment History | ✅ | ✅ (cache) |
| Tax Status | ✅ | ❌ |
| Social Benefits | ✅ | ❌ |

### 2.4 Sécurité

| Fonction | Support |
|----------|---------|
| PIN Code | ✅ |
| Biometric Auth | ✅ (Face ID / Fingerprint) |
| 2FA (OTP) | ✅ |
| Device Binding | ✅ |
| Screen Lock | ✅ |
| Auto-Logout | ✅ |
| Session Management | ✅ |

---

## 3. OFFLINE CAPABILITIES

### 3.1 Ce qui fonctionne hors-ligne

```
┌─────────────────────────────────────────────┐
│              OFFLINE MODE                    │
├─────────────────────────────────────────────┤
│  ✅ Identity Wallet (lecture)                │
│  ✅ QR Identity (génération)                 │
│  ✅ QR Verification (scan)                   │
│  ✅ Birth Certificate (cache local)          │
│  ✅ Notifications (buffer)                   │
│  ✅ Profile (lecture)                        │
│  ⬜ Service Requests (buffer, sync later)    │
│  ❌ New Registration (online required)       │
│  ❌ Digital Signature (online required)      │
└─────────────────────────────────────────────┘
```

### 3.2 Sync Strategy

| Stratégie | Description |
|-----------|-------------|
| **Cache First** | Lire depuis le cache, rafraîchir depuis le serveur |
| **Write Behind** | Écrire localement, synchroniser plus tard |
| **Conflict Resolution** | "Last Write Wins" avec horodatage |
| **Selective Sync** | Prioriser les données critiques |

---

## 4. ARCHITECTURE TECHNIQUE

### 4.1 Stack

| Couche | Technologie |
|--------|-------------|
| Frontend Mobile | React Native / Flutter |
| Frontend Web | React (PWA) |
| State Management | Redux / Zustand |
| Local Storage | SQLite (encrypted) |
| Sync Engine | Custom Offline Sync |
| Push Notifications | Firebase / National |
| Biometrics | OS Native API |
| QR Library | Custom QR SDK |

### 4.2 Data Flow

```
User ──▶ App UI ──▶ State Manager ──▶ API Layer
                                        │
                                  ┌─────┴──────┐
                                  │  Online?    │
                                  └─────┬──────┘
                                        │
                          ┌─────────────┴─────────────┐
                          │                           │
                    ┌─────▼─────┐              ┌─────▼─────┐
                    │   Online  │              │  Offline  │
                    └─────┬─────┘              └─────┬─────┘
                          │                          │
                    ┌─────▼─────┐              ┌─────▼─────┐
                    │  API Call │              │  Local DB │
                    └─────┬─────┘              └─────┬─────┘
                          │                          │
                    ┌─────▼─────┐              ┌─────▼─────┐
                    │  Return   │              │  Queue    │
                    │  & Cache  │              │  For Sync │
                    └───────────┘              └───────────┘
```

---

## 5. UI/UX SPECIFICATIONS

### 5.1 Navigation

```
┌────────────────────────────────┐
│          Home Screen           │
├────────────────────────────────┤
│  ┌──────┐ ┌──────┐ ┌──────┐  │
│  │  ID   │ │Wallet│ │Serv. │  │
│  │Card   │ │      │ │      │  │
│  └──────┘ └──────┘ └──────┘  │
│                                │
│  ┌────────────────────────┐   │
│  │  Quick Actions         │   │
│  │  • Voir mon identité   │   │
│  │  • Scanner un QR       │   │
│  │  • Mes notifications   │   │
│  └────────────────────────┘   │
│                                │
│  ┌────────────────────────┐   │
│  │  Recent Activity       │   │
│  └────────────────────────┘   │
├────────────────────────────────┤
│  Home  │  ID  │ Wallet│ More  │
└────────────────────────────────┘
```

### 5.2 Offline Indicator

```
┌────────────────────────────────┐
│  🔴 Offline Mode              │
│  ⬜ Services limited to cache │
│  ⏳ Pending sync: 3 items     │
└────────────────────────────────┘
```

---

## 6. SÉCURITÉ

### 6.1 Authentication Flow

```
┌──────────┐    ┌──────────┐    ┌──────────┐
│  Launch   │───▶│  Auth    │───▶│  Home    │
│   App     │    │  Screen  │    │  Screen  │
└──────────┘    └────┬─────┘    └──────────┘
                     │
               ┌─────▼─────┐
               │  PIN /    │
               │ Biometric │
               └─────┬─────┘
                     │
            ┌────────▼────────┐
            │  Validate Local │
            │  (No network)   │
            └────────┬────────┘
                     │
            ┌────────▼────────┐
            │    Session OK   │
            └─────────────────┘
```

### 6.2 Security Measures

| Measure | Implementation |
|---------|---------------|
| App Hardening | ProGuard, R8, Obfuscation |
| Root/Jailbreak Detection | Runtime check |
| SSL Pinning | Certificate pinning |
| Secure Storage | Android Keystore / iOS Keychain |
| Data Encryption | AES-256-GCM |
| Anti-Tampering | Integrity verification |
| Screen Capture Protection | FLAG_SECURE |

---

## 7. PERFORMANCE

| Métrique | Cible |
|----------|-------|
| App Launch Time | < 2s |
| Screen Load | < 500ms |
| QR Scan | < 200ms |
| Cache Read | < 50ms |
| Sync Operation | < 5s (10MB data) |
| Battery Impact | < 5% / heure |
| Memory Usage | < 150 MB |
| APK Size | < 30 MB |

---

## 8. MULTILINGUE

| Screen | Français | Créole | Anglais |
|--------|----------|--------|---------|
| Welcome Screen | Bienvenue | Byenveni | Welcome |
| My Identity | Mon Identité | Idantite Mwen | My Identity |
| Wallet | Portefeuille | Pòtfèy | Wallet |
| QR Scan | Scanner QR | Eskane QR | QR Scan |
| Notifications | Notifications | Notifikasyon | Notifications |

---

## 9. COMPLIANCE

| Requirement | Status |
|-------------|--------|
| RGPD (Loi Haïtienne) | ✅ |
| WCAG 2.1 AA | ✅ |
| Offline Certification | ✅ |
| Security Audit | ✅ |
| UX National Review | ✅ |

---

## 10. DÉPLOIEMENT

| Phase | Date | Canal |
|-------|------|-------|
| Beta Privée | J+15 | TestFlight / Play Store Internal |
| Beta Publique | J+30 | Play Store / App Store |
| Production | J+45 | Tous canaux |
| Kiosks | J+60 | Kiosks nationaux |

---
*Fin du document — Citizen Super App v1.0*
