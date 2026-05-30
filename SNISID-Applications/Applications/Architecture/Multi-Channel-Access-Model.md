# MULTI-CHANNEL ACCESS MODEL — SNISID
## Modèle d'Accès National Universel

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-MCA-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |

---

## 1. PRÉSENTATION

Modèle d'accès multi-canal garantissant que chaque citoyen haïtien peut accéder aux services SNISID quel que soit son appareil, sa localisation ou sa connectivité.

### 1.1 La Fracture Numérique Haïtienne

```
┌──────────────────────────────────────────────┐
│        RÉALITÉ NUMÉRIQUE HAÏTIENNE            │
├──────────────────────────────────────────────┤
│  Smartphones: 35% de la population           │
│  Internet mobile: 38% de couverture          │
│  Ordinateurs: 10% des foyers                 │
│  Électricité: 45% des zones rurales          │
│  Alphabétisation numérique: ~30%             │
│  Créole unilingues: ~70% de la population   │
│  SMS: accessible à 95% (téléphones basiques)│
└──────────────────────────────────────────────┘
```

---

## 2. CANAUX SUPPORTÉS

### 2.1 Android

| Support | Détail |
|---------|--------|
| **Min SDK** | Android 10 (API 29) |
| **Target SDK** | Android 15 (API 35) |
| **Distribution** | Google Play Store + APK direct |
| **Offline** | ✅ Full |
| **Biometrics** | ✅ Fingerprint + Face |
| **Secure Element** | ✅ Android KeyStore |
| **Languages** | FR, HT, EN |

### 2.2 iOS

| Support | Détail |
|---------|--------|
| **Min Version** | iOS 15 |
| **Target** | iOS 19 |
| **Distribution** | Apple App Store |
| **Offline** | ✅ Full |
| **Biometrics** | ✅ Face ID + Touch ID |
| **Secure Element** | ✅ Secure Enclave |
| **Languages** | FR, HT, EN |

### 2.3 Web

| Support | Détail |
|---------|--------|
| **PWA** | ✅ Progressive Web App |
| **Browsers** | Chrome, Firefox, Safari, Edge |
| **Offline** | ✅ Service Worker |
| **Responsive** | ✅ Mobile-first |
| **Accessibility** | ✅ WCAG 2.1 AA |

### 2.4 Kiosks

| Support | Détail |
|---------|--------|
| **Type** | Bornes physiques publiques |
| **Localisation** | Bureaux d'état civil, mairies, hôpitaux |
| **Nombre** | 500 bornes (phase 1) |
| **Écran** | Tactile 21" |
| **Impression** | Reçus, QR codes |
| **Scanner** | Documents, QR |
| **Offline** | ✅ Kiosk mode |
| **Accessibilité** | Audio, Braille, Fauteuil roulant |

### 2.5 Offline Devices

| Support | Détail |
|---------|--------|
| **Type** | Appareils dédiés agents terrain |
| **Modèle** | Smartphone renforcé (Samsung XCover, etc.) |
| **Batterie** | > 8h active |
| **Stockage** | 128 GB min |
| **Résistance** | IP68, MIL-STD-810G |
| **Offline** | ✅ Full air-gap mode |
| **Biometrics** | ✅ Empreinte + Face |

---

## 3. MATRICE D'ACCÈS

### 3.1 Services par Canal

| Service | Android | iOS | Web | Kiosk | Offline Device |
|---------|---------|-----|-----|-------|---------------|
| ID Verification | ✅ | ✅ | ✅ | ✅ | ✅ |
| Identity Wallet | ✅ | ✅ | ✅ | ✅ | ✅ |
| Birth Certificate | ✅ | ✅ | ✅ | ✅ | ✅ |
| Civil Registry | ✅ | ✅ | ✅ | ❌ | ❌ |
| Notifications | ✅ | ✅ | ✅ (PWA) | ❌ | ✅ |
| QR Verify | ✅ | ✅ | ✅ (camera) | ✅ | ✅ |
| Government Ops | ✅ | ✅ | ✅ | ❌ | ❌ |
| Police Cases | ✅ | ✅ | ❌ | ❌ | ✅ |
| Field Enrollment | ✅ | ❌ | ❌ | ❌ | ✅ |
| Admin Portal | ❌ | ❌ | ✅ | ❌ | ❌ |

### 3.2 Fonctionnalités par Canal

| Fonctionnalité | Mobile App | Web | Kiosk | Offline |
|----------------|------------|-----|-------|---------|
| Push Notifications | ✅ | ✅ (PWA) | ❌ | ✅ |
| Geofencing | ✅ | ❌ | ❌ | ✅ |
| Camera Access | ✅ | ✅ (browser) | ✅ | ✅ |
| NFC | ✅ | ❌ | ❌ | ✅ |
| Biometrics | ✅ | ❌ | ❌ | ✅ |
| Voice Navigation | ⬜ | ❌ | ✅ | ❌ |
| Print | ❌ | ✅ | ✅ | ❌ |
| Offline Storage | ✅ | ✅ | ✅ | ✅ |

---

## 4. ACCESSIBILITÉ NATIONALE

### 4.1 Couverture Géographique

```
┌──────────────────────────────────────────────┐
│           COUVERTURE NATIONALE                │
├──────────────────────────────────────────────┤
│  Départements: 10/10                         │
│  Arrondissements: 42/42                      │
│  Communes: 145/145                           │
│  Sections communales: ~570                   │
│                                              │
│  Kiosks:                                     │
│  • Chef-lieu département: 10                 │
│  • Chef-lieu arrondissement: 42              │
│  • Communes principales: 100                 │
│  • Hôpitaux publics: 50                      │
│  • Bureaux état civil: 200                   │
│  • Frontières: 10                            │
│  • Aéroports: 5                              │
│                                              │
│  Bornes connectées: 🛰️ Satellite + 4G       │
│  Bornes isolées: Mode kiosk offline         │
└──────────────────────────────────────────────┘
```

### 4.2 Accessibilité pour Tous

| Groupe | Android | iOS | Web | Kiosk |
|--------|---------|-----|-----|-------|
| **Voyants** | ✅ | ✅ | ✅ | ✅ |
| **Malvoyants** | ✅ TalkBack | ✅ VoiceOver | ✅ Screen reader | ✅ Audio guidance |
| **Non-voyants** | ✅ TalkBack | ✅ VoiceOver | ✅ Screen reader | ✅ Braille pad |
| **Sourds** | ✅ Texte | ✅ Texte | ✅ Texte | ✅ Texte |
| **Mobilité réduite** | ✅ Voice control | ✅ Voice control | ✅ Keyboard nav | ✅ Hauteur adaptée |
| **Analphabètes** | ✅ Icônes + Audio | ✅ Icônes + Audio | ⬜ | ✅ Audio + Icônes |
| **Créole unilingues** | ✅ HT | ✅ HT | ✅ HT | ✅ HT |

---

## 5. ARCHITECTURE MULTI-CANAL

```
┌──────────────────────────────────────────────┐
│          MULTI-CHANNEL ORCHESTRATOR           │
├──────────────────────────────────────────────┤
│                                              │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐     │
│  │  Mobile  │ │   Web    │ │  Kiosk   │     │
│  │  App     │ │   (PWA)  │ │   App    │     │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘     │
│       │            │            │            │
│       └────────────┼────────────┘            │
│                    │                         │
│  ┌─────────────────▼──────────────────────┐ │
│  │         API GATEWAY                     │ │
│  │    Rate limiting • Auth • Routing      │ │
│  └─────────────────┬──────────────────────┘ │
│                    │                         │
│  ┌─────────────────▼──────────────────────┐ │
│  │         CORE SERVICES                   │ │
│  └────────────────────────────────────────┘ │
│                                              │
│  ┌─────────────────┬──────────────────────┐ │
│  │  Mobile Sync    │   Kiosk Sync         │ │
│  │  Engine         │   Engine             │ │
│  └─────────────────┴──────────────────────┘ │
│                                              │
└──────────────────────────────────────────────┘
```

---

## 6. STRATÉGIE DE DÉPLOIEMENT

### 6.1 Phases de Déploiement

| Phase | Timeline | Canaux | Régions |
|-------|----------|--------|---------|
| **P1** | J0-J30 | Android + iOS | Ouest |
| **P2** | J15-J45 | Web (PWA) | Ouest + Artibonite |
| **P3** | J30-J60 | Kiosks (50) | Capitales départements |
| **P4** | J45-J75 | Kiosks (200) | Arrondissements |
| **P5** | J60-J90 | Offline devices | Zones rurales |
| **P6** | J75-J100 | Kiosks (500) | Communes + frontières |
| **P7** | J90-J120 | Full deployment | National |

### 6.2 Priorisation Géographique

```
┌──────────────────────────────────────────────┐
│        PRIORITÉ GÉOGRAPHIQUE                 │
├──────────────────────────────────────────────┤
│  P1: Ouest (Port-au-Prince, Pétion-Ville,   │
│      Delmas, Carrefour)                      │
│  P2: Artibonite (Gonaïves, Saint-Marc)      │
│  P3: Nord (Cap-Haïtien)                     │
│  P4: Sud (Les Cayes)                        │
│  P5: Grand'Anse, Nippes                     │
│  P6: Centre, Nord-Est, Nord-Ouest           │
│  P7: Sud-Est, toutes communes               │
└──────────────────────────────────────────────┘
```

---

## 7. PERFORMANCE PAR CANAL

| Métrique | Mobile | Web | Kiosk | Offline |
|----------|--------|-----|-------|---------|
| Page Load | < 2s | < 3s | < 2s | < 1s |
| Transaction | < 5s | < 5s | < 10s | < 3s |
| QR Scan | < 1s | < 2s | < 1s | < 0.5s |
| Sync | < 30s | < 1min | < 1min | < 2min |
| Uptime | > 99.9% | > 99.95% | > 99% | > 99.9% |
| Battery Life | > 8h | N/A | N/A | > 8h |

---

## 8. SÉCURITÉ PAR CANAL

| Canal | Auth | Storage | Transport |
|-------|------|---------|-----------|
| Android | PIN + Biometric + 2FA | KeyStore + Encrypted | TLS 1.3 |
| iOS | PIN + Biometric + 2FA | Keychain + Encrypted | TLS 1.3 |
| Web | Password + 2FA | Session + HSTS | TLS 1.3 |
| Kiosk | Card + PIN + Agent | HSM + Encrypted | TLS 1.3 + mTLS |
| Offline Device | Biometric + Token | Encrypted + Secure Element | TLS 1.3 + Sync |

---
*Fin du document — Multi-Channel Access Model v1.0*