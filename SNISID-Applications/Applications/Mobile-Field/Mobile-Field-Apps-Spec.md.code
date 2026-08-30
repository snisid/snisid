# MOBILE FIELD APPS — SNISID
## Applications Terrain Nationales

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-MFA-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |
| Plateforme | Android (prioritaire), iOS |

---

## 1. PRÉSENTATION

Applications terrains conçues pour fonctionner **sans internet** dans les zones reculées d'Haïti. Ces applications permettent l'enrôlement, la vérification d'identité, la capture biométrique et la validation QR en tout lieu.

### 1.1 Contexte Haïtien

```
┌──────────────────────────────────────────────┐
│        RÉALITÉ TERRAIN HAÏTIEN               │
├──────────────────────────────────────────────┤
│  • 60% zones rurales sans internet stable    │
│  • Coupures électriques fréquentes           │
│  • Réseau mobile 2G/3G intermittent         │
│  • Agents terrain avec smartphones basique  │
│  • Forte chaleur, poussière, humidité        │
└──────────────────────────────────────────────┘
```

---

## 2. APPLICATIONS TERRAIN

### 2.1 Field Enrollment App

| Fonction | Support | Offline |
|----------|---------|---------|
| Enregistrement citoyen | ✅ | ✅ Complet |
| Capture photo | ✅ | ✅ |
| Scan documents | ✅ | ✅ |
| Collecte signatures | ✅ | ✅ |
| Saisie données personnelles | ✅ | ✅ |
| QR Generation | ✅ | ✅ |
| Biometrics (fingerprint) | ✅ | ✅ |
| Geolocation | ✅ | ✅ |
| Batch operations | ✅ | ✅ |

### 2.2 Field Identity Verification App

| Fonction | Support | Offline |
|----------|---------|---------|
| QR Code Scan | ✅ | ✅ |
| Identity Data Display | ✅ | ✅ |
| Photo Match | ✅ | ✅ |
| Fingerprint Match | ✅ | ✅ |
| Document Verify | ✅ | ✅ |
| Watchlist Check | ✅ | ✅ (cache) |
| Verification Report | ✅ | ✅ |

### 2.3 Field Biometric Capture App

| Fonction | Support | Offline |
|----------|---------|---------|
| Fingerprint Capture | ✅ (10 doigts) | ✅ |
| Face Capture | ✅ | ✅ |
| Iris Capture (optionnel) | ✅ | ✅ |
| Biometric Quality Check | ✅ | ✅ |
| Template Generation | ✅ | ✅ |
| Template Storage (encrypted) | ✅ | ✅ |
| Batch Biometric Capture | ✅ | ✅ |

### 2.4 Field QR Validation App

| Fonction | Support | Offline |
|----------|---------|---------|
| QR Scan | ✅ | ✅ |
| Identity Validation | ✅ | ✅ |
| Certificate Validation | ✅ | ✅ |
| Tamper Detection | ✅ | ✅ |
| Validation Log | ✅ | ✅ |
| Offline Audit Trail | ✅ | ✅ |

---

## 3. ARCHITECTURE OFFLINE

### 3.1 Data Flow

```
┌─────────────────────────────────────────────┐
│         FIELD DEVICE (OFFLINE)               │
├─────────────────────────────────────────────┤
│  ┌─────────────────────────────────────┐    │
│  │  Mobile Field App                    │    │
│  └────────────┬────────────────────────┘    │
│               │                              │
│  ┌────────────▼────────────────────────┐    │
│  │  Local Database (Encrypted SQLite)    │    │
│  │  • Enrollment Queue                   │    │
│  │  • Verification Logs                  │    │
│  │  • Biometric Templates                │    │
│  │  • QR Cache                           │    │
│  │  • Event Buffer                       │    │
│  └────────────┬────────────────────────┘    │
│               │                              │
│  ┌────────────▼────────────────────────┐    │
│  │  Secure Storage (AES-256)            │    │
│  │  • Biometric Data                    │    │
│  │  • Personal Information              │    │
│  │  • Audit Logs                        │    │
│  └────────────┬────────────────────────┘    │
│               │                              │
└───────────────┼──────────────────────────────┘
                │
    ┌───────────▼───────────┐
    │   SYNC WHEN ONLINE     │
    │                        │
    │  1. Compress & Encrypt │
    │  2. Push to Server     │
    │  3. Verify Checksum   │
    │  4. Clear Buffer      │
    │  5. Receive ACK       │
    └───────────────────────┘
```

### 3.2 Offline Capabilities Matrix

```
┌──────────────────────────────────────────────────────┐
│                    FIELD APP MODES                    │
├──────────────┬──────────┬────────────┬───────────────┤
│  Fonction    │  Online  │  Offline   │  Degraded     │
├──────────────┼──────────┼────────────┼───────────────┤
│ Enrollment   │  Full    │  Full      │  Full         │
│ ID Verify    │  Full    │  Full      │  Full         │
│ Biometrics   │  Full    │  Full      │  Full         │
│ QR Validate  │  Full    │  Full      │  Full         │
│ Sync Data    │  Real-t. │  Queue     │  Partial Sync │
│ Geolocation  │  Live    │  Cached    │  Cached       │
│ Updates      │  Auto    │  Manual    │  Manual       │
└──────────────┴──────────┴────────────┴───────────────┘
```

---

## 4. SYNC STRATEGIE

### 4.1 Prioritized Sync Queue

| Priorité | Data Type | Max Delay |
|----------|-----------|-----------|
| **P0** | Emergency / Critical | Immédiat (si réseau) |
| **P1** | Biometric Data | 1 heure |
| **P2** | Enrollment Records | 24 heures |
| **P3** | Verification Logs | 48 heures |
| **P4** | Audit Trails | 7 jours |

### 4.2 Compression & Encryption

```
Raw Data (10 MB)
    │
    ▼
├── Compression (gzip) → ~2 MB
│       │
│       ▼
├── Encryption (AES-256-GCM)
│       │
│       ▼
├── Chunking (1 MB chunks)
│       │
│       ▼
├── Transfer (HTTPS with Resume)
│       │
│       ▼
├── Server Verification (Checksum)
│       │
│       ▼
├── Buffer Cleanup (After ACK)
```

---

## 5. UI/UX SPÉCIFIQUE TERRAIN

### 5.1 Design Principles

| Principe | Description |
|----------|-------------|
| **Large Touch Targets** | ≥ 48px pour usage avec gants |
| **High Contrast** | Visibilité en plein soleil |
| **Minimal Steps** | ≤ 3 clics pour action principale |
| **Auto-Save** | Sauvegarde automatique toutes les 30s |
| **Battery Efficient** | Optimisé pour longue durée |
| **Offline First** | Toute action enregistrée localement |

### 5.2 Offline Status UI

```
┌────────────────────────────────┐
│  📡 Offline  │  🔋 85%        │
│  Syncing: 12 pending          │
│  Last sync: 2h ago            │
│  [Sync Now] (when connected)  │
└────────────────────────────────┘
```

### 5.3 Enrollment Screen

```
┌────────────────────────────────┐
│  ENROLLMENT — Étape 2/5       │
├────────────────────────────────┤
│  Photo du citoyen              │
│  ┌────────────────────────┐   │
│  │   [Camera Preview]     │   │
│  └────────────────────────┘   │
│  [Capture Photo]              │
│                                │
│  Nom: [________________]      │
│  Prénom: [________________]   │
│  Date Naiss: [____-__-__]    │
│                                │
│  [◀ Précédent]  [Suivant ▶]  │
└────────────────────────────────┘
```

---

## 6. SÉCURITÉ TERRAIN

### 6.1 Device Security

| Measure | Implementation |
|---------|---------------|
| Device Binding | Hardware-backed attestation |
| Remote Wipe | Si perte ou vol |
| Local Encryption | AES-256-GCM |
| Auto-Lock | After 5 min inactivity |
| Tamper Detection | Runtime integrity check |
| Session Timeout | After 30 min |

### 6.2 Data Protection

```
┌─────────────────────────────────────┐
│    LOCAL DATA PROTECTION            │
├─────────────────────────────────────┤
│  • Biometric data NEVER stored raw │
│  • Templates only (encrypted)      │
│  • PII encrypted column-level      │
│  • Auto-delete after sync (config) │
│  • Secure element for keys         │
│  • Anti-forensic on lockscreen     │
└─────────────────────────────────────┘
```

---

## 7. BATTERY & PERFORMANCE

| Métrique | Cible |
|----------|-------|
| Battery Life (active) | > 8 hours |
| Battery Life (standby) | > 24 hours |
| Enrollment Time | < 3 min per person |
| QR Scan | < 1s |
| Biometric Capture | < 10s |
| App Launch | < 2s |
| Storage per 1000 enr. | < 200 MB |
| Offline Queue capacity | > 10,000 records |

---

## 8. MATÉRIEL RECOMMANDÉ

| Device | Specs Minimales | Recommandé |
|--------|-----------------|------------|
| Smartphone | Android 10, 3GB RAM | Samsung XCover Pro |
| Biometric Scanner | USB-C / Bluetooth | MORPHO / Suprema |
| Power Bank | 10,000 mAh | 20,000 mAh |
| Micro SD | 32 GB | 128 GB |
| Protection | IP68 | IP68 + MIL-STD-810G |

---

## 9. DÉPLOIEMENT TERRAIN

| Phase | Timeline | Régions |
|-------|----------|---------|
| P12.4a | J+15 | Ouest (Port-au-Prince, Pétion-Ville) |
| P12.4b | J+30 | Artibonite, Nord |
| P12.4c | J+45 | Sud, Grand'Anse, Nippes |
| P12.4d | J+60 | Centre, Nord-Est, Sud-Est |
| P12.4e | J+75 | Nord-Ouest, toutes communes |

---
*Fin du document — Mobile Field Apps v1.0*