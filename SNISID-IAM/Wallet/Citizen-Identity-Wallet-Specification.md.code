# Citizen Identity Wallet Specification

> **Portefeuille d'Identité Numérique Citoyen — SNISID**  
> **Version :** 1.0.0  
> **Classification :** SOUVERAIN — RESTREINT  
> **Dernière mise à jour :** 2026-05-25

---

## 1. OBJECTIF

Créer un portefeuille d'identité numérique portable pour le citoyen. Le wallet doit fonctionner **offline**.

---

## 2. FONCTIONS SUPPORTÉES

| Fonction | Support | Description |
|----------|---------|-------------|
| Digital Identity | ✅ | Identité numérique vérifiée (VC/DID) |
| Certificates | ✅ | Certificats officiels dématérialisés |
| Signatures | ✅ | Signature électronique qualifiée (QES) |
| Consent Management | ✅ | Gestion des consentements |
| Offline Verification | ✅ | Vérification sans connexion |
| QR Validation | ✅ | Validation par code QR |

---

## 3. ARCHITECTURE DU WALLET

```
┌────────────────────────────────────────────────────┐
│              CITIZEN IDENTITY WALLET               │
├────────────────────────────────────────────────────┤
│  PRESENTATION: Mobile App | Web Portal | Kiosk    │
│  SECURITY: Biometric Auth | PIN | Secure Enclave  │
│  IDENTITY: NNU Profile | VCs | Presentation Defs  │
│  CONSENT: Consent Manager | Access Control | Hist. │
└────────────────────────────────────────────────────┘
```

---

## 4. CONTENU DU WALLET

### 4.1 Identité Numérique (Verifiable Credential)

```json
{
  "@context": "https://www.w3.org/2018/credentials/v1",
  "type": ["VerifiableCredential", "NationalIdentityCredential"],
  "id": "vc:snisid:identity:uuid",
  "issuer": "https://identity.snisid.ht",
  "issuanceDate": "2026-05-25T10:00:00Z",
  "credentialSubject": {
    "id": "did:snisid:nnu:NNU-HA-20260525-A7K9M2X4-73-1",
    "given_name": "Jean",
    "family_name": "Baptiste",
    "date_of_birth": "1990-01-15",
    "place_of_birth": {"country": "HT", "department": "Ouest"},
    "sex": "M",
    "nationality": "HT"
  },
  "proof": {
    "type": "RsaSignature2018",
    "created": "2026-05-25T10:00:00Z",
    "verificationMethod": "did:snisid:authority#keys-1"
  }
}
```

### 4.2 Certificats

| Certificat | Émetteur | Validité | Usage |
|------------|----------|----------|-------|
| Naissance | État Civil | À vie | Preuve d'identité |
| Identité nationale | SNISID IAM | 5 ans | Accès services |
| Résidence | Commune | 2 ans | Preuve domicile |
| Biométrie | Plateforme Biométrique | À vie | Vérification forte |

### 4.3 Signature Électronique

```yaml
electronic_signature:
  type: Qualified Electronic Signature (QES)
  standard: eIDAS / ISO 14888
  algorithm: ECDSA P-384
  key_storage: Secure Enclave (mobile) / HSM (serveur)
  usage: [sign_consent, sign_documents, sign_access_requests, sign_legal_declarations]
```

---

## 5. MODE OFFLINE

### 5.1 Fonctionnement

| Scénario | Processus |
|----------|-----------|
| Vérification sans connexion | Citoyen présente QR → Agent vérifie signature locale (certificat cache) |
| Accès sans connexion | Auth biométrique locale → Affichage identité stockée → Sync à reconnexion |
| Signature sans connexion | Clé privée dans Secure Enclave → Document signé → Transmission à reconnexion |

### 5.2 Cache Local

```yaml
local_cache:
  certificates:
    validity: 30_days
    storage: encrypted_keychain
  identity_data:
    fields: [nnu, name, dob, photo_hash]
    encryption: AES-256-GCM
    key_derivation: PBKDF2 (100k iterations)
  consent_records:
    max_stored: 100
  access_history:
    max_stored: 500
```

### 5.3 QR Code Validation

```
Format QR: SNISID:ID:YYYYMMDDHHMMSS:NNN:SIG:EXP

SNISID = Préfixe système
ID     = Identifiant session
YYYYMMDDHHMMSS = Timestamp génération
NNN    = NNU masqué (****M2X4**)
SIG    = Signature courte (Base64)
EXP    = Expiration (minutes)

Vérification: Scanner → Extraire → Vérifier signature (cert cache) → Vérifier expiration → Résultat
Offline: Vérification locale uniquement
Online:  Vérification + synchronisation registre
```

---

## 6. SÉCURITÉ DU WALLET

| Méthode | Usage | Sécurité |
|---------|-------|----------|
| Biométrie (empreinte/face) | Ouverture quotidienne | Haute |
| PIN code | Authentification principale | Moyenne |
| MFA | Actions sensibles | Très haute |

### Protection des Données

| Donnée | Protection | Stockage |
|--------|-----------|----------|
| Identité | Chiffrée (AES-256) | Secure Enclave |
| Certificats | Signés (ECDSA) | Keychain chiffré |
| Signatures | Clé privée isolée | HSM / Secure Element |
| Consentements | Hashés + Signés | Local + Cloud |
| Historique | Chiffré | Local + Cloud |

### Gestion des Risques

| Risque | Mitigation |
|--------|-----------|
| Perte appareil | Révocation à distance + Recovery |
| Vol appareil | Verrouillage biométrique + PIN |
| Clonage | Secure Enclave + Anti-tampering |
| Usurpation | MFA + Liveness detection |
| Fuite données | Chiffrement E2E + Zero-knowledge |

---

## 7. WORKFLOWS WALLET

| Workflow | Étapes |
|----------|--------|
| Activation | Télécharg App → Vérif Identité → Enrôlement → Activ Wallet |
| Récupération | Perte Appareil → Auth Secours → Vérif Identité → Restore Wallet |
| Révocation | Demande → Vérif Identité → Révocat Certificats → Confirm Archive |

---

## 8. CONFORMITÉ

| Standard | Support |
|----------|---------|
| W3C Verifiable Credentials | ✅ |
| Decentralized Identifiers (DID) | ✅ |
| OpenID Connect for Verifiable Presentations | ✅ |
| ISO/IEC 18013-5 (mDL) | ✅ |
| NIST SP 800-63-3 (Digital Identity) | ✅ |

---

> **Le wallet citoyen fonctionne offline, est sécurisé et portable.**
