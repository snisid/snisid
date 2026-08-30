# SNISID — PKI Nationale & Trust Infrastructure
**Système National d'Identification et d'Inscription Digitale**

**Classification:** TOP-SECRET / RESTREINT DEFENSE  
**Statut:** Phase 5 — Fondation de Confiance Cryptographique Souveraine  
**Date:** 2026-05-25

---

## Objectif

Construire la fondation de confiance numérique souveraine du SNISID. Cette phase est l'une des plus critiques du projet entier.

**Règle absolue:**
> **Dans SNISID : aucune confiance sans cryptographie.**

---

## Structure du Repository

```
PKI/
├── docs/
│   └── SNISID-National-PKI-Architecture.md        # Architecture officielle nationale
├── root-ca/
│   └── root-ca-ceremony-procedure.md              # Procédure cérémonie Root CA (6 custodians, air-gap, HSM)
├── intermediate-cas/
│   ├── government/
│   │   └── government-ca-policy.yaml              # Government CA — institutions, ministries, agencies
│   ├── citizen/
│   │   └── citizen-ca-policy.yaml                 # Citizen CA — identité numérique, wallet, privacy
│   ├── infrastructure/
│   │   └── infrastructure-ca-policy.yaml          # Infra CA — K8s, APIs, mesh, auto-issuance
│   ├── judicial/
│   │   └── judicial-ca-policy.yaml                # Judicial CA — tribunaux, magistrats, notaires
│   └── device/
│       └── device-ca-policy.yaml                  # Device CA — enrollment stations, edge, mobile, biometric
├── hsm/
│   └── hsm-topology.yaml                          # Thales Luna 7 topology Core/DR, partitions, Shamir escrow
├── certificates/
│   └── certificate-lifecycle-management.yaml        # CLM — issuance, renewal, rotation, revocation, recovery
├── policies/
│   └── (CPS/CP references — see docs/)             # Certificate Policies & Practice Statements
├── workflows/
│   └── bpm-pki-workflows.yaml                     # BPMN workflows — Gov issuance, Judicial issuance, Emergency revocation, CA rotation
├── mtls/
│   └── national-mtls-model.yaml                   # mTLS STRICT mesh-wide — Istio, Cilium, Kafka, etcd, K8s
├── citizen-certs/
│   └── citizen-identity-certs.yaml                # Citizen X.509 profile, wallet (PKCS#12), offline QR verification
├── kubernetes/
│   └── kubernetes-pki.yaml                        # K8s PKI — etcd, API server, kubelet, ingress, mesh, cert-manager
├── device-trust/
│   └── device-enrollment.yaml                     # Device attestation (TPM/Secure Boot), provisioning, revocation
├── ocsp-crl/
│   └── ocsp-crl-infrastructure.yaml               # OCSP responders HA, CRL generation, distribution, delta-CRL
├── digital-signatures/
│   └── signature-services.yaml                    # TSA, CAdES/XAdES/PAdES/JAdES, legal value levels
├── audit/
│   └── pki-observability.yaml                     # Prometheus rules — cert expiry, HSM health, OCSP, revocation, TSA
├── runbooks/
│   ├── root-ca-compromise.md                      # Cellule Crise PKI, reconstruction nationale, cérémonie accélérée
│   ├── hsm-failure.md                             # HSM hardware/tamper/network failure, DR failover
│   └── emergency-revocation.md                    # Mass revocation, intermediate compromise, government/judicial revocation
├── standards/
│   └── pki-standards.md                           # Algorithmes, durées, formats TLS, HSM, CLM, OCSP/CRL, signatures
└── README.md                                      # Ce document
```

---

## Hiérarchie PKI Nationale

| Niveau | Entité | Validité | Usage | Isolation |
|--------|--------|----------|-------|-----------|
| **Root CA** | SNISID Root CA Nationale | 10 ans | Signer Intermediates uniquement | Offline, air-gapped, HSM salle blanche |
| **Gov ICA** | Government CA | 5 ans | Institutions, ministères, ambassades | HSM Core DC, Vault PKI mount |
| **Citizen ICA** | Citizen CA | 3 ans | Identités citoyennes numériques | HSM Core DC, privacy partition |
| **Infra ICA** | Infrastructure CA | 5 ans | Kubernetes, APIs, service mesh | HSM Core DC, 100% automated |
| **Judicial ICA** | Judicial CA | 5 ans | Tribunaux, magistrats, notaires | HSM Core DC, 2-person control |
| **Device ICA** | Device CA | 3 ans | Équipements terrain, edge, biométrie | HSM Core DC, attestation-based |

---

## Principes absolus

1. **Root CA jamais online.** Physiquement et réseautiquement isolée. Clé jamais hors HSM.
2. **Clés privées critiques jamais hors HSM.** Génération, stockage, usage dans module hardware sécurisé.
3. **Multi-person control.** Cérémonies de clés nécessitant N personnes physiquement présentes (M-of-N).
4. **Audit traçable.** Toute opération PKI est loggée, vidéo-enregistrée, et archivée.
5. **Aucun certificat n'expire silencieusement.** Monitoring Prometheus 7j/1j/0j + auto-renewal + runbooks.
6. **Révocation immédiate possible.** CRL/OCSP 1h-24h, emergency < 5 min pour infrastructure.
7. **Citizen privacy by design.** Pas de PII dans les certificats — UUID anonymisé uniquement.
8. **Standards cryptographiques stricts.** P-384, SHA-384, TLS 1.3, AES-256-GCM — jamais de fallback faible.

---

## Métriques critiques monitorées

| Métrique | Seuil | Action |
|----------|-------|--------|
| `cert_expiring_7d` | > 0 | Warning — auto-renewal |
| `cert_expiring_1d` | > 0 | Critical — force renewal NOW |
| `cert_expired_active` | > 0 | **Emergency** — service impact |
| `hsm_tamper_detected` | == 1 | **Catastrophe** — Cellule Crise PKI |
| `ocsp_failure_rate` | > 0.1% | Investigation SOC |
| `revocation_queue_depth` | > 100 | Escalade IGC |

---

## Contacts & Gouvernance

| Rôle | Entité | Décision |
|------|--------|----------|
| Root CA Owner | Présidence / Haute autorité | Cérémonies Root, politique générale |
| Intermediate CA Owner | IGC | Création/revocation Intermediates, audits |
| Government CA Manager | Direction numérique | Délivrance institutions |
| Citizen CA Manager | Agence identité numérique | Délivrance citoyens |
| Infrastructure CA Manager | Équipe Infra Nationale | Délivrance auto services |
| Judicial CA Manager | Conseil supérieur judiciaire | Délivrance magistrats |
| Device CA Manager | Équipe Terrain & IoT | Provisioning équipements |
| Revocation Authority | IGC + SOC National | Révocations d'urgence |
| Audit Authority | Cour des comptes / Audit externe souverain | Audits annuels |

---

*SNISID PKI Nationale — Fondation de confiance cryptographique souveraine.*
*Classification: TOP-SECRET — SNISID Phase 5 Trust Infrastructure*
