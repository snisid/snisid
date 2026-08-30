# APPLICATION SECURITY FRAMEWORK — SNISID
## Cadre de Sécurité Applicative National

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-SEC-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |
| Classification | CONFIDENTIEL |

---

## 1. PRÉSENTATION

Framework de sécurité unifié pour toutes les applications SNISID, garantissant une protection résistante aux attaques quel que soit le canal ou le contexte d'utilisation.

### 1.1 Menaces Couvertes

| Menace | Couverture |
|--------|-----------|
| Reverse Engineering | ✅ Anti-tampering, Obfuscation |
| Man-in-the-Middle | ✅ Certificate Pinning, TLS 1.3 |
| Device Compromise | ✅ Root/Jailbreak Detection |
| Data Theft | ✅ Encryption, Secure Storage |
| Replay Attacks | ✅ Nonce, Timestamp |
| Cloning | ✅ Device Attestation |
| Phishing | ✅ MFA, Biometric Verification |
| Malware | ✅ App Integrity Check |

---

## 2. SÉCURITÉ PAR COUCHE

### 2.1 Application Layer

| Measure | Implementation | Niveau |
|---------|---------------|--------|
| **Obfuscation** | ProGuard, R8, DexGuard | Toutes apps |
| **Anti-Tampering** | Integrity verification at startup | Toutes apps |
| **Root Detection** | Runtime check (multiple methods) | Toutes apps |
| **Debug Detection** | Anti-debugging (TracerPid, ptrace) | Toutes apps |
| **Emulator Detection** | Runtime environment check | Sensibles |
| **Screen Protection** | FLAG_SECURE for sensitive screens | Sensibles |
| **Memory Protection** | Sensitive data in native memory | Ultra sécurisées |

### 2.2 Transport Layer

| Measure | Implementation |
|---------|---------------|
| **TLS 1.3** | All network communications |
| **Mutual TLS** | Server-to-server and critical endpoints |
| **Certificate Pinning** | SHA-256 pins in app binary |
| **HSTS** | Strict transport security |
| **ALPN** | Protocol negotiation |
| **OCSP Stapling** | Certificate status check |

### 2.3 Data Layer

```
┌──────────────────────────────────────────────┐
│            DATA PROTECTION                    │
├──────────────────────────────────────────────┤
│  At Rest:                                    │
│  • Database: SQLCipher (AES-256)             │
│  • Files: AES-256-GCM per file               │
│  • Keys: Android KeyStore / iOS Keychain     │
│  • SharedPrefs: EncryptedSharedPreferences   │
│                                               │
│  In Transit:                                 │
│  • TLS 1.3 with Pinning                      │
│  • JWT with short expiry                     │
│  • Payload encryption for sensitive data     │
│                                               │
│  In Use:                                     │
│  • Memory encryption (native)                │
│  • Secure enclave for biometric data         │
│  • Zero-copy for sensitive processing        │
└──────────────────────────────────────────────┘
```

### 2.4 Authentication Layer

| Facteur | Support | Usage |
|---------|---------|-------|
| **Something you know** | PIN, Password | Base |
| **Something you have** | OTP Token, Device | Medium |
| **Something you are** | Biometric (Face/Finger) | Medium |
| **Location** | Geo-fencing | Contextual |
| **Hardware Token** | HSM, Smart Card | Ultra sécurisé |

---

## 3. DEVICE ATTESTATION

### 3.1 Android

```
┌──────────────────────────────────────────────┐
│        GOOGLE PLAY INTEGRITY API              │
├──────────────────────────────────────────────┤
│  Checks performed:                           │
│  ✅ Device integrity (unmodified)            │
│  ✅ App integrity (signed by SNISID)         │
│  ✅ Play Protect enabled                     │
│  ✅ CTS profile match                        │
│  ✅ No known vulnerabilities                 │
└──────────────────────────────────────────────┘
```

### 3.2 iOS

```
┌──────────────────────────────────────────────┐
│          APPLE DEVICECHECK                    │
├──────────────────────────────────────────────┤
│  Checks performed:                           │
│  ✅ Device integrity                         │
│  ✅ App integrity                            │
│  ✅ No jailbreak                             │
│  ✅ Secure Enclave available                 │
└──────────────────────────────────────────────┘
```

---

## 4. SECURE STORAGE

### 4.1 Key Management

| Key Type | Storage | Algorithm | Rotation |
|----------|---------|-----------|----------|
| **Device Key** | Secure Enclave | ECDSA P-256 | Device reset |
| **Auth Key** | Keychain/Keystore | HMAC-SHA256 | 90 days |
| **Session Key** | Memory only | AES-256 | Per session |
| **Data Key** | Keychain/Keystore | AES-256-GCM | Per data item |
| **Sync Key** | Secure Enclave | ECDH | Per sync batch |

### 4.2 Certificate Pinning

```json
{
  "pins": [
    {
      "hostname": "api.snisid.gouv.ht",
      "sha256_pins": [
        "sha256/abc123def456...",
        "sha256/789ghi012jkl..."
      ],
      "expiry": "2027-12-31"
    },
    {
      "hostname": "identity.snisid.gouv.ht",
      "sha256_pins": [
        "sha256/def456abc789..."
      ],
      "expiry": "2027-12-31"
    }
  ]
}
```

---

## 5. ANTI-TAMPERING

### 5.1 Protection Layers

```
┌──────────────────────────────────────────────┐
│         ANTI-TAMPERING STACK                  │
├──────────────────────────────────────────────┤
│  Layer 1: App Integrity Check                 │
│  • SHA-256 hash of APK/IPA at startup        │
│  • Compare with expected hash                │
│                                              │
│  Layer 2: Runtime Integrity                   │
│  • Check code signatures                     │
│  • Detect debugger attachment                │
│  • Detect hooking frameworks (Frida, Xposed) │
│                                              │
│  Layer 3: Environment Check                   │
│  • Root/Jailbreak detection                  │
│  • Emulator detection                        │
│  • VPN/Proxy detection (for sensitive ops)   │
│                                              │
│  Layer 4: Behavioral Detection                │
│  • Anomalous usage patterns                  │
│  • Rapid successive operations               │
│  • Unusual geolocation                       │
└──────────────────────────────────────────────┘
```

### 5.2 Response to Tampering

| Détection | Action | Niveau |
|-----------|--------|--------|
| App Hash mismatch | 🔴 App crash + Remote wipe | Critique |
| Root detected | 🟡 Limited mode | Élevé |
| Debugger detected | 🟡 Warning + Log | Moyen |
| Hooking detected | 🔴 App crash | Critique |
| Emulator detected | 🟡 Limited mode | Moyen |
| Anomalous behavior | 🟡 Audit + Lock | Contextuel |

---

## 6. SÉCURITÉ PAR TYPE D'APPLICATION

| Application Type | Auth | Storage | Transport | Attestation |
|-----------------|------|---------|-----------|-------------|
| Citizen App | PIN + Biometric | Encrypted Local | TLS 1.3 | Basic |
| Government App | MFA (2FA) | Encrypted + HSM | TLS 1.3 + mTLS | Full |
| Police/Justice | MFA (3FA) | HSM + Airgap | mTLS + E2E | Full + HW |
| Wallet | Biometric + PIN | Secure Enclave | TLS 1.3 | Full |
| Field Apps | Biometric + Token | Encrypted Local | TLS 1.3 + Sync | Medium |

---

## 7. AUDIT & LOGGING

| Type | Détail | Rétention |
|------|--------|-----------|
| Authentication | User, Time, Method, Success/Fail | 1 an |
| Data Access | Resource, Action, User, Timestamp | 1 an |
| Security Events | Type, Severity, Details, User | 5 ans |
| Tamper Attempts | Full forensic data | 5 ans |
| Sync Operations | Data type, Size, Timestamp, Success | 90 jours |

---

## 8. COMPLIANCE

| Standard | Status |
|----------|--------|
| OWASP MASVS (Mobile) | L1 + L2 |
| ISO 27001 | ✅ |
| NIST SP 800-163 | ✅ |
| PCI-DSS (mobile) | ✅ |
| GDPR/Loi Haïtienne | ✅ |

---
*Fin du document — Application Security Framework v1.0*