# NATIONAL UX/UI DESIGN SYSTEM — SNISID
## Système de Design National

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-UX-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |

---

## 1. PRÉSENTATION

Système de design unifié pour toutes les applications SNISID garantissant une expérience cohérente, accessible et sécurisée à tous les citoyens haïtiens.

### 1.1 Principes Directeurs

| Principe | Description |
|----------|-------------|
| **Mobile-First** | Conçu d'abord pour mobile, adapté au desktop |
| **Accessible** | WCAG 2.1 AA minimum |
| **Sécurisé** | UX intégrant la sécurité sans friction |
| **Offline-First** | Expérience cohérente en ligne et hors-ligne |
| **Multilingue** | Design adaptable au Français, Créole, Anglais |
| **Souverain** | Identité visuelle nationale haïtienne |

---

## 2. COMPOSANTS

### 2.1 Core Components

| Component | Status | Accessible | Mobile-First |
|-----------|--------|------------|--------------|
| Button | ✅ | ✅ | ✅ |
| Input Field | ✅ | ✅ | ✅ |
| Dropdown | ✅ | ✅ | ✅ |
| Checkbox | ✅ | ✅ | ✅ |
| Radio Button | ✅ | ✅ | ✅ |
| Toggle | ✅ | ✅ | ✅ |
| Date Picker | ✅ | ✅ | ✅ |
| Search Bar | ✅ | ✅ | ✅ |
| Badge | ✅ | ✅ | ✅ |
| Avatar | ✅ | ✅ | ✅ |

### 2.2 Complex Components

| Component | Status | Description |
|-----------|--------|-------------|
| Identity Card | ✅ | Affichage identité numérique |
| QR Display | ✅ | Affichage QR code avec sécurité |
| QR Scanner | ✅ | Scanner avec retour haptique |
| Offline Banner | ✅ | Indicateur d'état de connexion |
| Sync Status | ✅ | Progression synchro |
| Biometric Prompt | ✅ | Interface déverrouillage biométrique |
| Signature Pad | ✅ | Zone signature électronique |
| Notification Card | ✅ | Notification avec badge vérifié |
| Timeline | ✅ | Suivi d'événements |
| Secure File View | ✅ | Visualisation documents sécurisés |

### 2.3 Navigation Components

| Component | Status |
|-----------|--------|
| Bottom Navigation | ✅ |
| Top App Bar | ✅ |
| Side Menu | ✅ |
| Tab Bar | ✅ |
| Stepper | ✅ |
| Breadcrumbs | ✅ |
| Pagination | ✅ |

---

## 3. ACCESSIBILITÉ

### 3.1 Standards

| Standard | Niveau | Statut |
|----------|--------|--------|
| WCAG 2.1 | AA | ✅ Implémenté |
| WCAG 2.1 | AAA | 🟡 En cours |
| RNIB Guidelines | — | ✅ |
| Mobile Accessibility | iOS/Android | ✅ |

### 3.2 Accessibility Features

| Feature | Support | Description |
|---------|---------|-------------|
| Screen Reader | ✅ | TalkBack, VoiceOver |
| High Contrast | ✅ | Mode contraste élevé |
| Large Text | ✅ | Jusqu'à 200% |
| Reduced Motion | ✅ | Pour épilepsie |
| Focus Indicator | ✅ | Navigation clavier |
| Color Blind Mode | ✅ | Daltonien-friendly |
| Font Scaling | ✅ | Dynamique |
| Touch Target | ✅ | ≥ 48px |

### 3.3 Color Palette Accessible

```
┌──────────────────────────────────────────────┐
│           PALETTE COULEURS PRINCIPALE         │
├──────────────────────────────────────────────┤
│  Primary:    #1B4F72  (Bleu National)        │
│  Secondary:  #E74C3C  (Rouge)                │
│  Accent:     #F39C12  (Jaune)                │
│  Success:    #27AE60  (Vert)                 │
│  Warning:    #F1C40F  (Jaune clair)          │
│  Error:      #C0392B  (Rouge foncé)          │
│  Info:       #2980B9  (Bleu clair)           │
│                                              │
│  Text:       #2C3E50  (Presque noir)         │
│  Background: #FFFFFF  (Blanc)                │
│  Surface:    #F8F9FA  (Gris clair)           │
│                                              │
│  ✅ Tous les contrastes ≥ 4.5:1 (AA)        │
│  ✅ ≥ 7:1 pour textes < 18px                │
└──────────────────────────────────────────────┘
```

---

## 4. OFFLINE UX

### 4.1 Offline States

| State | Visual | Action |
|-------|--------|--------|
| **Online** | 🟢 Connected | Normal operations |
| **Degraded** | 🟡 Limited | Cache mode, some features disabled |
| **Offline** | 🔴 Offline | Local-only operations |
| **Reconnecting** | 🔄 Syncing | Background sync in progress |
| **Sync Error** | ⚠️ Sync Error | Manual retry needed |

### 4.2 Offline UI Pattern

```
┌──────────────────────────────────────────────┐
│              OFFLINE MODE                     │
├──────────────────────────────────────────────┤
│  🔴 Vous êtes hors-ligne                     │
│  • Consultation des documents: ✅             │
│  • QR Verification: ✅                        │
│  • Nouvelles demandes: ⬜                    │
│  • Synchronisation: 12 éléments en attente   │
│                                              │
│  Dernière synchro: il y a 2h                 │
│  [Synchroniser maintenant]                   │
└──────────────────────────────────────────────┘
```

### 4.3 Offline Data Indicator

```
┌────────────────────────────────┐
│  📦 Offline data available     │
│  ID Card: ✅ (cached)         │
│  Birth Cert: ✅ (cached)      │
│  Notifications: 3 new         │
└────────────────────────────────┘
```

---

## 5. SECURITY UX

### 5.1 Authentication UX

```
┌────────────────────────────────┐
│  🇭🇹 SNISID                    │
│  Déverrouillez votre compte    │
│                                │
│  [👤 Face ID / Touch ID]      │
│  ou                            │
│  [____] [____] [____] [____]  │
│    ●    ●    ●    ●           │
│                                │
│  [Mot de passe oublié]        │
└────────────────────────────────┘
```

### 5.2 Security Indicators

| Indicator | Meaning |
|-----------|---------|
| 🔒 | Connection sécurisée |
| ✅ Vérifié | Notification signée SNISID |
| 🛡️ | Protection active |
| ⚠️ | Attention requise |
| 🚫 | Accès refusé |

---

## 6. RESPONSIVE DESIGN

### 6.1 Breakpoints

| Device | Width | Layout |
|--------|-------|--------|
| Mobile S | 320px | Single column |
| Mobile M | 375px | Single column |
| Mobile L | 425px | Single column |
| Tablet | 768px | Two columns |
| Desktop | 1024px | Multi columns |
| Wide | 1440px | Full layout |

### 6.2 Grid System

```
Mobile (4 columns):
[ 1 ][ 2 ][ 3 ][ 4 ]

Tablet (8 columns):
[ 1 ][ 2 ][ 3 ][ 4 ][ 5 ][ 6 ][ 7 ][ 8 ]

Desktop (12 columns):
[ 1 ][ 2 ][ 3 ][ 4 ][ 5 ][ 6 ][ 7 ][ 8 ][ 9 ][10][11][12]
```

---

## 7. TYPOGRAPHIE

| Style | Police (FR/HT/EN) | Size | Weight |
|-------|-------------------|------|--------|
| H1 | Poppins | 32px | Bold |
| H2 | Poppins | 24px | Bold |
| H3 | Poppins | 20px | Semi-Bold |
| Body | Inter | 16px | Regular |
| Body Small | Inter | 14px | Regular |
| Caption | Inter | 12px | Regular |
| Button | Poppins | 16px | Medium |
| Monospace | JetBrains Mono | 14px | Regular |

---

## 8. ICONOGRAPHIE

| Type | Style | Format |
|------|-------|--------|
| System Icons | Outline, 24dp | SVG |
| Action Icons | Filled, 24dp | SVG |
| Status Icons | Colored | SVG |
| Identity Icons | Custom shield | SVG |
| National Symbols | Coat of Arms | SVG |

---

## 9. DESIGN TOKENS

```json
{
  "colors": {
    "primary": "#1B4F72",
    "secondary": "#E74C3C",
    "success": "#27AE60",
    "warning": "#F1C40F",
    "error": "#C0392B",
    "info": "#2980B9"
  },
  "spacing": {
    "xs": 4,
    "sm": 8,
    "md": 16,
    "lg": 24,
    "xl": 32,
    "xxl": 48
  },
  "borderRadius": {
    "sm": 4,
    "md": 8,
    "lg": 12,
    "xl": 16,
    "full": 9999
  },
  "shadows": {
    "sm": "0 1px 2px rgba(0,0,0,0.1)",
    "md": "0 4px 6px rgba(0,0,0,0.1)",
    "lg": "0 10px 15px rgba(0,0,0,0.1)"
  },
  "animation": {
    "fast": "150ms ease",
    "normal": "300ms ease",
    "slow": "500ms ease"
  }
}
```

---

## 10. COMPLIANCE

| Standard | Status |
|----------|--------|
| WCAG 2.1 AA | ✅ |
| Mobile Accessibility | ✅ |
| Platform Guidelines (Material/Apple) | ✅ |
| SNISID Security Standards | ✅ |
| Offline UX Certification | ✅ |

---
*Fin du document — National UX/UI Design System v1.0*