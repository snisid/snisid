# DIGITAL IDENTITY WALLET — SNISID
## Portefeuille Identité National

| Attribut | Valeur |
|----------|--------|
| Document ID | SNISID-PH12-WAL-001 |
| Version | 1.0 |
| Statut | APPROUVÉ |
| Plateforme | Android, iOS |

---

## 1. PRÉSENTATION

Portefeuille d'identité numérique nationale hautement sécurisé permettant aux citoyens haïtiens de stocker, gérer et partager leurs certificats d'identité de manière sécurisée, même hors-ligne.

### 1.1 Principes Clés

```
┌──────────────────────────────────────────────┐
│        HAUTEMENT SÉCURISÉ                     │
├──────────────────────────────────────────────┤
│  ✅ Stockage local chiffré AES-256           │
│  ✅ Clés dans Secure Enclave / TEE          │
│  ✅ Zero-knowledge proofs disponibles        │
│  ✅ Consentement explicite requis            │
│  ✅ Offline QR verification                  │
│  ✅ Anti-tampering + Anti-cloning            │
│  ✅ Révocation à distance                    │
└──────────────────────────────────────────────┘
```

---

## 2. FONCTIONNALITÉS

### 2.1 Identity Certificates

| Certificat | Support | Offline |
|------------|---------|---------|
| National ID Card | ✅ | ✅ |
| Birth Certificate | ✅ | ✅ |
| Digital Driver License | ✅ | ✅ |
| Voter ID | ✅ | ✅ |
| Tax ID (NIF) | ✅ | ✅ |
| Social Security | ✅ | ✅ |
| Professional License | ✅ | ✅ |
| Passport | ✅ | ✅ |

### 2.2 Digital Signatures

| Fonction | Support | Offline |
|----------|---------|---------|
| Document Signing | ✅ | ✅ |
| Qualified Signatures | ✅ (with HSM) | ❌ |
| Signature Verification | ✅ | ✅ |
| Batch Signing | ✅ | ❌ |
| Signature Expiry | ✅ | ✅ |
| Revocation Check | ✅ | ❌ |

### 2.3 Consent Management

| Fonction | Support |
|----------|---------|
| Selective Disclosure | ✅ |
| Data Minimization | ✅ |
| Consent Grant | ✅ |
| Consent Revoke | ✅ |
| Consent History | ✅ |
| Expiring Consent | ✅ |
| Granular Permissions | ✅ |
| Audit Log | ✅ |

### 2.4 Offline QR Verification

| Fonction | Support |
|----------|---------|
| QR Generation | ✅ |
| QR Scan | ✅ |
| Signed QR | ✅ |
| Tamper Detection | ✅ |
| Expiry Check | ✅ |
| Selective QR Data | ✅ |
| Batch QR Verify | ✅ |

### 2.5 Secure Storage

| Fonction | Support |
|----------|---------|
| Hardware-backed Encryption | ✅ |
| Biometric Lock | ✅ |
| PIN + Biometric | ✅ |
| Auto-Lock | ✅ |
| Remote Wipe | ✅ |
| Backup & Recovery | ✅ (encrypted) |
| Multiple Profiles | ✅ |

---

## 3. ARCHITECTURE

### 3.1 Wallet Structure

```
┌─────────────────────────────────────────────┐
│          DIGITAL IDENTITY WALLET             │
├─────────────────────────────────────────────┤
│  ┌──────────────────────────────────────┐   │
│  │  LOCK SCREEN                          │   │
│  │  [PIN] [Biometric] [PIN+Biometric]   │   │
│  └──────────────────────────────────────┘   │
│                                             │
│  ┌──────────────────────────────────────┐   │
│  │  CERTIFICATE STORAGE                  │   │
│  │  ┌──────────┐ ┌──────────┐           │   │
│  │  │ National │ │  Birth   │           │   │
│  │  │    ID    │ │  Cert.   │           │   │
│  │  └──────────┘ └──────────┘           │   │
│  │  ┌──────────┐ ┌──────────┐           │   │
│  │  │ License  │ │  Voter   │           │   │
│  │  └──────────┘ └──────────┘           │   │
│  └──────────────────────────────────────┘   │
│                                             │
│  ┌──────────────────────────────────────┐   │
│  │  KEY MANAGEMENT                       │   │
│  │  ┌──────────┐ ┌──────────┐           │   │
│  │  │ Signing  │ │Encryption│           │   │
│  │  │  Keys    │ │  Keys    │           │   │
│  │  └──────────┘ └──────────┘           │   │
│  │  ┌──────────┐ ┌──────────┐           │   │
│  │  │ Consent  │ │Recovery  │           │   │
│  │  │  Keys    │ │  Keys    │           │   │
│  │  └──────────┘ └──────────┘           │   │
│  └──────────────────────────────────────┘   │
│                                             │
│  ┌──────────────────────────────────────┐   │
│  │  AUDIT LOG                            │   │
│  │  • All access logged                  │   │
│  │  • All sharing logged                 │   │
│  │  • All signatures logged              │   │
│  │  • Immutable local ledger             │   │
│  └──────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

### 3.2 Key Hierarchy

```
┌─────────────────────────────────────┐
│        MASTER KEY (Secure Enclave)   │
└────────────────┬────────────────────┘
                 │
    ┌────────────┴────────────┐
    │                         │
┌───▼───┐               ┌───▼───┐
│Device │               │Recovery│
│ Key   │               │  Key   │
└───┬───┘               └───┬───┘
    │                         │
┌───▼───────────────────────▼───┐
│       DERIVED KEYS             │
├───────────────────────────────┤
│  • Signing Key (ECDSA)        │
│  • Encryption Key (AES)       │
│  • Consent Key (HMAC)         │
│  • QR Signing Key (Ed25519)   │
│  • Audit Key (SHA-256)        │
└───────────────────────────────┘
```

---

## 4. SÉCURITÉ

### 4.1 Protection Measures

| Measure | Implementation |
|---------|---------------|
| **Root/Jailbreak Detection** | Runtime + Startup check |
| **Debug Detection** | Anti-debugger |
| **Screen Capture** | FLAG_SECURE / Obscured |
| **Key Storage** | Android KeyStore / iOS Keychain |
| **Data at Rest** | AES-256-GCM with IV |
| **Data in Use** | Memory encryption |
| **App Attestation** | Google Play Integrity / iOS DeviceCheck |
| **Certificate Pinning** | SHA-256 pins |
| **Anti-Replay** | Timestamp + Nonce |

### 4.2 Offline Security

```
┌──────────────────────────────────────────────┐
│          OFFLINE SECURITY MODEL               │
├──────────────────────────────────────────────┤
│  Toutes les vérifications sont locales :     │
│                                              │
│  1. QR signé avec clé privée du wallet       │
│  2. Clé publique connue du verificateur      │
│  3. Signature vérifiée localement            │
│  4. Horodatage via timestamp local           │
│  5. Anti-replay via nonce + counter          │
│                                              │
│  Aucun appel réseau nécessaire              │
│  pour vérification de base                   │
└──────────────────────────────────────────────┘
```

---

## 5. UI/UX

### 5.1 Wallet Home

```
┌────────────────────────────────┐
│  🇭🇹 SNISID Wallet  │ 🔒      │
├────────────────────────────────┤
│                                │
│  ┌────────────────────────┐   │
│  │    [Photo ID]          │   │
│  │    Jean-Marie DUPONT   │   │
│  │    ID: 1234-5678-90    │   │
│  └────────────────────────┘   │
│                                │
│  ┌──────┐ ┌──────┐ ┌──────┐  │
│  │ID Card│ │Birth │ │License│  │
│  │      │ │ Cert │ │       │  │
│  └──────┘ └──────┘ └──────┘  │
│                                │
│  ┌──────┐ ┌──────┐ ┌──────┐  │
│  │  QR   │ │Sign  │ │Consent│  │
│  │ Show  │ │Doc   │ │ Mgmt  │  │
│  └──────┘ └──────┘ └──────┘  │
│                                │
│  ┌────────────────────────┐   │
│  │  Recent Activity        │   │
│  │  • Shared ID (2m ago)  │   │
│  │  • Signed contract     │   │
│  │    (1h ago)            │   │
│  └────────────────────────┘   │
├────────────────────────────────┤
│  Home │ Certs │ Sign │ More   │
└────────────────────────────────┘
```

---

## 6. PERFORMANCE

| Métrique | Cible |
|----------|-------|
| Wallet Unlock | < 500ms |
| QR Generation | < 200ms |
| QR Verification | < 100ms |
| Document Sign | < 2s |
| Certificate Load | < 300ms |
| Storage per 10 certs | < 50 MB |
| Battery Impact | < 3% / heure |

---

## 7. BACKUP & RECOVERY

### 7.1 Recovery Options

| Méthode | Sécurité | Offline |
|---------|----------|---------|
| Recovery Phrase | Haute (BIP39) | ✅ |
| Government Recovery | Très haute | ❌ (nécessite bureau) |
| Secure Cloud Backup | Haute | ❌ |
| Local Export | Moyenne | ✅ |

### 7.2 Recovery Flow

```
┌──────────┐    ┌──────────┐    ┌──────────┐
│ Lost     │───▶│  Enter   │───▶│  Verify  │
│ Device   │    │ Recovery │    │ Identity │
│          │    │  Phrase  │    │          │
└──────────┘    └──────────┘    └────┬─────┘
                                     │
                            ┌────────▼────────┐
                            │  Re-generate     │
                            │  Key Hierarchy   │
                            └────────┬────────┘
                                     │
                            ┌────────▼────────┐
                            │  Restore         │
                            │  Certificates    │
                            │  (from server)   │
                            └────────┬────────┘
                                     │
                            ┌────────▼────────┐
                            │  Wallet Ready    │
                            └─────────────────┘
```

---

## 8. INTÉGRATIONS

| Système | Type |
|---------|------|
| Identity Hub | Certificate issuance & verification |
| QR Verification API | Validation service |
| Digital Signature Service | Qualified signatures |
| Consent Registry | Consent management |
| National Notification | Alerts and updates |

---

## 9. COMPLIANCE

| Standard | Conformité |
|----------|------------|
| eIDAS (aligné) | ✅ |
| W3C VC Data Model | ✅ |
| ISO 27001 | ✅ |
| PCI-DSS (wallet) | ✅ |
| GDPR/Loi Haïtienne | ✅ |

---
*Fin du document — Digital Identity Wallet v1.0*