# Biometric Platform Specification

> **Plateforme Biométrique Nationale SNISID**  
> **Version :** 1.0.0  
> **Classification :** SOUVERAIN — RESTREINT  
> **Dernière mise à jour :** 2026-05-25

---

## 1. OBJECTIF

Créer l'identité forte nationale via une plateforme biométrique souveraine. La biométrie est :
- **Chiffrée** — Stockée et transmise de manière sécurisée
- **Auditée** — Chaque utilisation est journalisée
- **Protégée juridiquement** — Conforme aux lois nationales

---

## 2. BIOMÉTRIES SUPPORTÉES

| Type | Support | Priorité | Seuil qualité |
|------|---------|----------|--------------|
| Empreinte digitale | ✅ | Primaire | ≥ 0.85 |
| Reconnaissance faciale | ✅ | Primaire | ≥ 0.80 |
| Iris | Optionnel | Secondaire | ≥ 0.90 |
| Voix | Optionnel | Secondaire | ≥ 0.75 |

---

## 3. ARCHITECTURE

```
[Capteurs Biométriques] → [Gateway TLS 1.3 + mTLS] → [Processing Engine]
                                                        ├── Quality Check
                                                        ├── Liveness Detection
                                                        ├── Feature Extraction
                                                        └── Template Generation
                                                              ↓
                                                    [Matching Engine]
                                                        ├── 1:1 Verify
                                                        ├── 1:N Identify
                                                        └── N:N Dedup
                                                              ↓
                                                  [Storage Chiffré AES-256-GCM]
```

---

## 4. SPÉCIFICATIONS TECHNIQUES

### 4.1 Format de Gabarit Biométrique

| Section | Taille | Contenu |
|---------|--------|---------|
| Header | 128 bytes | Version, type, qualité, timestamp, device, location, officer, liveness |
| Template Data | Variable | Payload chiffré AES-256-GCM, référence clé HSM, HMAC SHA-256 |

### 4.2 Chiffrement

```yaml
encryption:
  at_rest:
    algorithm: AES-256-GCM
    key_management: HSM (FIPS 140-3 Level 3)
    key_rotation: every_90_days
  
  in_transit:
    protocol: TLS 1.3
    mutual_authentication: true (mTLS)
  
  template_protection:
    method: cancelable_biometrics
    description: Template transformé avec sel national
    reversal: impossible_without_national_salt
    salt_storage: HSM souverain
```

### 4.3 Liveness Detection

| Biométrie | Méthodes | Seuil | Timeout |
|-----------|----------|-------|---------|
| Face | Eye blink, head movement, texture, depth | 0.85 | 3000ms |
| Fingerprint | Pulse, skin conductance, ridge analysis | 0.90 | 2000ms |
| Iris | Pupil reflex, texture, spectral | 0.95 | 2000ms |
| Voice | Spectral, prosodic, challenge-response | 0.80 | 5000ms |

---

## 5. PROTECTION JURIDIQUE

### 5.1 Conformité

| Loi / Règlement | Application |
|----------------|-------------|
| Loi nationale protection des données | Conforme |
| ISO/IEC 19794 | Implémenté |
| ISO/IEC 30107 (liveness) | Implémenté |
| NIST SP 800-63 (identité numérique) | Implémenté |
| FIPS 140-3 (HSM) | Implémenté |

### 5.2 Droits du Citoyen

| Droit | Application |
|-------|-------------|
| Savoir quelles données sont collectées | Notification à l'enrôlement |
| Consentir explicitement | Signature ou confirmation biométrique |
| Accéder à ses données biométriques | Wallet biometric view |
| Corriger des données erronées | Correction workflow |
| Retirer son consentement | Consent revocation |
| Demander la suppression | Deletion workflow avec revue judiciaire |

---

## 6. PERFORMANCE

| Métrique | Cible |
|----------|-------|
| Temps d'enrôlement biométrique | < 30 secondes par modalité |
| Temps de vérification 1:1 | < 2 secondes |
| Temps d'identification 1:N | < 5 secondes (N ≤ 10M) |
| Temps de dé-duplication N:N | < 30 secondes (N ≤ 10M) |
| Taux de faux acceptation (FAR) | < 0.001% |
| Taux de faux rejet (FRR) | < 0.1% |
| Disponibilité | 99.99% |
| Capacité maximale | 15 millions citoyens |

---

> **La biométrie est chiffrée, auditée et protégée juridiquement.**
