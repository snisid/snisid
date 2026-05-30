# SNISID — National PKI Architecture
**Classification:** TOP-SECRET / RESTREINT DEFENSE  
**Version:** 5.0.0  
**Date:** 2026-05-25  
**Statut:** OFFICIEL — Phase 5 Trust Infrastructure Nationale

---

## 1. Vue d'ensemble stratégique

La PKI nationale SNISID constitue la **fondation de confiance cryptographique souveraine** de l'État. Elle garantit :

- L'authenticité de toute identité numérique nationale (citoyens, institutions, équipements)
- La confidentialité des communications inter-services et inter-territoires
- L'intégrité des signatures électroniques à valeur légale nationale
- La non-répudiation des transactions gouvernementales critiques

### Principes absolus
- ❌ **Aucune confiance sans cryptographie.** Toute assertion d'identité doit être vérifiable cryptographiquement.
- ❌ **Root CA jamais online.** La racine de confiance est physiquement et réseautiquement isolée.
- ❌ **Clés privées critiques jamais hors HSM.** Génération, stockage, usage dans module hardware sécurisé.
- ✅ **Multi-person control.** Cérémonies de clés nécessitant N personnes physiquement présentes (M-of-N).
- ✅ **Audit traçable.** Toute opération PKI est loggée, vidéo-enregistrée, et archivée.

---

## 2. Architecture hiérarchique nationale

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SNISID NATIONAL PKI HIERARCHY                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌──────────────────────────────────────────────────────────────┐          │
│   │           SNISID ROOT CA NATIONALE (SNISID-ROOT-CA)           │          │
│   │              Offline • HSM Air-gapped • 10 ans               │          │
│   │         État : Clé générée en salle blanche Tier-4            │          │
│   └────────┬─────────────┬─────────────┬─────────────┬───────────┘          │
│            │             │             │             │                       │
│   ┌────────▼──┐ ┌────────▼──┐ ┌────────▼──┐ ┌────────▼──┐                  │
│   │ Gov ICA    │ │ Cit ICA   │ │ Infra ICA  │ │ Jud ICA   │                  │
│   │Government │ │ Citizen   │ │Infrastructure│ │ Judicial  │                  │
│   │   CA      │ │   CA      │ │    CA      │ │   CA      │                  │
│   │(Tier-0)   │ │(Tier-1)   │ │ (Tier-0)   │ │(Tier-0)   │                  │
│   └─────┬─────┘ └─────┬─────┘ └─────┬──────┘ └─────┬─────┘                  │
│         │             │             │              │                         │
│   ┌─────▼────┐  ┌─────▼────┐ ┌─────▼────┐  ┌─────▼────┐                    │
│   │Ministries│  │Citizen   │ │ K8s Mesh │  │Courts    │                    │
│   │Agencies  │  │Wallet    │ │ API GW   │  │Tribunals │                    │
│   │Local Gov │  │Mobile ID │ │ Nodes    │  │Notaries  │                    │
│   │Embassies │  │Biometric │ │ Devices  │  │Officers  │                    │
│   └──────────┘  └──────────┘ └──────────┘  └──────────┘                    │
│                                                                              │
│   ┌──────────────────────────────────────────────────────────────┐          │
│   │           DEVICE TRUST CA (Issued under Infra ICA)          │          │
│   │    Enrollment Stations • Edge Nodes • Mobile Units • Bio      │          │
│   └──────────────────────────────────────────────────────────────┘          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.1 Hiérarchie détaillée

| Niveau | Entité | Validité | Usage | Isolation |
|--------|--------|----------|-------|-----------|
| **Root CA** | SNISID Root CA Nationale | 10 ans | Signer uniquement les Intermediate CAs | Offline, air-gapped, HSM salle blanche |
| **Intermediate CA 1** | SNISID Government CA | 5 ans | Institutions, ministères, collectivités | HSM Core DC, Vault PKI mount |
| **Intermediate CA 2** | SNISID Citizen CA | 3 ans | Identités citoyennes numériques | HSM Core DC, Vault PKI mount, consentement RGPD |
| **Intermediate CA 3** | SNISID Infrastructure CA | 5 ans | Kubernetes, APIs, réseau, service mesh | HSM Core DC + DR, auto-issuance cert-manager |
| **Intermediate CA 4** | SNISID Judicial CA | 5 ans | Juridictions, tribunaux, officiers | HSM Core DC, accès strict judiciaire |
| **Intermediate CA 5** | SNISID Device CA | 3 ans | Équipements terrain, biométrie, edge | HSM Core DC, provisioning contrôlé |

### 2.2 Segregation cryptographique

Chaque Intermediate CA dispose de :
- Sa propre clé privée **ECDSA P-384** (ou RSA 4096 si compatibilité legacy requise)
- Son propre path Vault PKI (`snisid-pki-gov`, `snisid-pki-citizen`, etc.)
- Ses propres Certificate Policies (OID 1.3.6.1.4.1.xxx.x)
- Son propre CRL et OCSP responder
- Son propre équipe de gestionnaires (ségrégation des rôles)

---

## 3. Root CA Nationale — SNISID-ROOT-CA

### 3.1 Caractéristiques

| Propriété | Valeur |
|-----------|--------|
| **Common Name** | `CN=SNISID Root CA Nationale, O=État, C=HT` |
| **Validité** | 10 ans (2036-05-25) |
| **Algorithme** | ECDSA avec courbe P-384 (secp384r1) |
| **Hash** | SHA-384 |
| **Key Usage** | keyCertSign, cRLSign — **UNIQUEMENT** |
| **Subject Key Identifier** | SHA-256 du clé publique |
| **Authority Key Identifier** | Self (Root) |

### 3.2 Protection physique

| Couche | Mesure |
|--------|--------|
| **Localisation** | Salle blanche Tier-4, site militaire ou institutionnel ultra-sécurisé |
| **Accès physique** | 4-eyes principle minimum (4 personnes présentes pour toute opération) |
| **HSM** | Thales Luna 7 Network HSM + Thales Luna 7 Backup HSM (DR) |
| **Network** | Aucune connectivité réseau (air-gap total). Console série locale uniquement |
| **Vidéo** | Enregistrement 24/7 de la salle, conservation 10 ans |
| **Alimentation** | UPS + générateur on-site, tamper-evident seals sur racks |
| **Médias backup** | Clés exportées via M-of-N Shamir sur supports optique sécurisés (3 sites distincts) |

### 3.3 Cérémonie de génération (Root Key Ceremony)

```
PHASE 1 — Préparation (M-3 mois)
├── Audit physique salle blanche (IGC + tiers externe)
├── Vérification HSM (firmware signé Thales, attestation)
├── Sélection 6 key-custodians (3 IGC, 1 Présidence, 1 Justice, 1 Technique)
└── Shamir config : 6 parts, threshold 4

PHASE 2 — Cérémonie (J-0)
├── Authentification biométrique + smartcard des 6 custodians
├── Vérification tamper-evident HSM (photos, sceaux)
├── Génération clé P-384 dans HSM (jamais exportée)
├── Test signature + vérification croisée
├── Création certificat auto-signé Root CA
├── Publication empreinte (hash SHA-256) dans registre national officiel
├── Export Shamir : 6 supports chiffrés, distribution 6 coffres distincts
└── Scellement HSM, vidéo archivée

PHASE 3 — Post-cérémonie (J+30j)
├── Audit externe complet (Big 4 souverain ou cabinet audité)
├── Publication trust anchor dans tous les systèmes nationaux
├── Distribution bundle Root CA aux edge nodes, citoyens, partenaires
└── Procédure de révocation Root CA (document de crise)
```

---

## 4. Intermediate CAs — Domaines nationaux

### 4.1 Government CA (`snisid-pki-gov`)

| Propriété | Valeur |
|-----------|--------|
| **Usage** | Certificats institutionnels (ministères, ambassades, collectivités) |
| **Validity** | 5 ans intermediate, 1-3 ans end-entity |
| **Subject format** | `CN={Institution} — SNISID Gov, O={Ministère}, OU={Service}, C=HT, serialNumber={SIRET-equivalent}` |
| **EKU** | serverAuth, clientAuth, codeSigning (documents officiels) |
| **Policy OID** | 1.3.6.1.4.1. SNISID.1.1 (Government Certificates) |
| **Issuance** | Manuel + workflow BPMN avec approbation hiérarchique |

### 4.2 Citizen CA (`snisid-pki-citizen`)

| Propriété | Valeur |
|-----------|--------|
| **Usage** | Identité numérique citoyenne, authentification, signature |
| **Validity** | 3 ans intermediate, 1-2 ans end-entity (citoyen) |
| **Subject format** | `CN={Nom Prénom} — {NIF}, OU=Citizen Identity, O=SNISID, C=HT, serialNumber={UUID-citoyen}` |
| **EKU** | clientAuth, emailProtection, **nonRepudiation** |
| **Policy OID** | 1.3.6.1.4.1.SNISID.1.2 (Citizen Identity Certificates) |
| **Issuance** | Automatisé post-enrollment biométrique validé, consentement électronique |
| **Privacy** | Pas de données biométriques dans le certificat. UUID anonymisé. |
| **Wallet** | PKCS#12 chiffré AES-256 + clé dérivée biométrique (template hash, pas image) |
| **Offline verification** | QR code signé contenant certificat condensé + signature CA |

### 4.3 Infrastructure CA (`snisid-pki-infra`)

| Propriété | Valeur |
|-----------|--------|
| **Usage** | Kubernetes, APIs, load balancers, service mesh, VPN nationaux |
| **Validity** | 5 ans intermediate, 90j end-entity (rotation auto) |
| **Subject format** | `CN={service}.{namespace}.snisid.gouv.local, O=SNISID Infra, C=HT` |
| **EKU** | serverAuth, clientAuth |
| **Policy OID** | 1.3.6.1.4.1.SNISID.1.3 (Infrastructure Certificates) |
| **Issuance** | 100% automatisé via cert-manager + Vault PKI (jamais manuel) |
| **Rotation** | Auto-renewal à 30j avant expiration. Emergency rotation < 2h. |

### 4.4 Judicial CA (`snisid-pki-judicial`)

| Propriété | Valeur |
|-----------|--------|
| **Usage** | Juridictions, tribunaux, officiers publics, notaires, huissiers |
| **Validity** | 5 ans intermediate, 1-3 ans end-entity |
| **Subject format** | `CN={Magistrat/Officier} — {Matricule}, O={Cour/Tribunal}, OU=Judicial, C=HT` |
| **EKU** | clientAuth, codeSigning, **nonRepudiation** |
| **Policy OID** | 1.3.6.1.4.1.SNISID.1.4 (Judicial Certificates) |
| **Issuance** | Workflow BPMN avec double validation (autorité judiciaire + IGC) |
| **Legal value** | Équivalent signature manuscrite selon cadre légal national |

### 4.5 Device CA (`snisid-pki-device`)

| Propriété | Valeur |
|-----------|--------|
| **Usage** | Équipements d'enrôlement, edge nodes, mobiles, biométriques |
| **Validity** | 3 ans intermediate, 1-2 ans end-entity |
| **Subject format** | `CN={Device-ID} — {Model}, O=SNISID Device Trust, OU={Tier}, C=HT, serialNumber={UUID-device}` |
| **EKU** | serverAuth, clientAuth |
| **Policy OID** | 1.3.6.1.4.1.SNISID.1.5 (Device Trust Certificates) |
| **Issuance** | Provisioning automatisé + attestation TPM/Secure Boot |
| **Revocation** | Révocation immédiate si vol, compromission, ou fin de vie |

---

## 5. Certificate Lifecycle Management (CLM)

### 5.1 Cycle complet

```
Enrollment ──► Validation ──► Issuance ──► Distribution ──► Usage ──► Renewal/Rotation ──► Revocation ──► Expiration ──► Archive
```

### 5.2 Policy par type

| Phase | Citizen | Government | Infrastructure | Judicial | Device |
|-------|---------|------------|----------------|----------|--------|
| **Enrollment** | Biométrie + pièce ID | Décret + autorité | Auto (Terraform/Helm) | Nomination + serment | Provisioning + TPM |
| **Validation** | Agent terrain + backoffice | IGC + RH | CI/CD policy | Collège + IGC | Secure Boot attestation |
| **Issuance** | Auto (Vault API) | Manuel (BPMN) | Auto (cert-manager) | Manuel (BPMN) | Auto (provisioning) |
| **Distribution** | Wallet mobile/cloud | Smartcard physique | Secret K8s | Smartcard + HSM judiciaire | Injected via provisioning |
| **Validity** | 1-2 ans | 1-3 ans | 90j | 1-3 ans | 1-2 ans |
| **Renewal** | Push notification + re-enrollment | Workflow annuel | Auto 30j avant | Workflow annuel | Auto provisioning refresh |
| **Revocation** | Portail citoyen + backoffice | IGC + hiérarchie | Auto (compromise detect) | Collège + IGC | Auto (MDM/Wazuh detect) |
| **Archive** | 7 ans (log uniquement) | 10 ans | 2 ans | 25 ans (légal) | 5 ans |

---

## 6. HSM Infrastructure Nationale

### 6.1 Topology

```
Core DC (Capitale)
├── HSM Thales Luna 7 — Slot 0 : Root CA (partition isolée, jamais online)
├── HSM Thales Luna 7 — Slot 1 : Intermediate CAs (Gov, Infra, Judicial)
├── HSM Thales Luna 7 — Slot 2 : Citizen CA (partition dédiée RGPD)
└── HSM Thales Luna 7 — Slot 3 : Device CA + Key Escrow temporaire

DR DC (Region Securisee)
├── HSM Thales Luna 7 — Slot 0 : Root CA Backup (clé répliquée, offline sauf cérémonie)
├── HSM Thales Luna 7 — Slot 1 : Intermediate CAs DR (mirror)
├── HSM Thales Luna 7 — Slot 2 : Citizen CA DR (mirror)
└── HSM Thales Luna 7 — Slot 3 : Device CA DR (mirror)

Key Escrow (3 sites indépendants)
├── Site A : Shamir parts 1+2 (Coffre présidence)
├── Site B : Shamir parts 3+4 (Coffre IGC)
└── Site C : Shamir parts 5+6 (Coffre Banque Centrale / Institution monétaire)
```

### 6.2 Governance HSM

| Aspect | Règle |
|--------|-------|
| **Key ceremonies** | 4-eyes minimum, 6-eyes pour Root CA. Vidéo, logs, témoins. |
| **Key escrow** | Shamir 4-of-6. Aucune personne ne détient >1 part. |
| **Backup HSM** | Synchronisation offline (transport physique scellé). Jamais en ligne simultané. |
| **Physical protection** | Salle blanche Tier-4, accès biométrie + smartcard + garde armée. |
| **Audit** | Toute opération HSM loggée SIEM temps réel. Vidéo 24/7. |
| **Compromise** | Procédure d'urgence : destruction clés HSM + révocation massive + reconstruction Root. |

---

## 7. National mTLS Model

Toutes les communications nationales critiques utilisent mTLS avec certificats SNISID Infrastructure CA.

### 7.1 Coverage

| Domaine | Certificat | Issuer | Rotation |
|---------|-----------|--------|----------|
| Kubernetes API servers | Server + Client | Infra CA | 90j auto |
| etcd peers | Peer + Client | Infra CA | 90j auto |
| Kubelet | Client | Infra CA | 90j auto |
| Istio Ingress Gateway | Server (wildcard) | Infra CA | 90j auto |
| Istio Sidecars (mTLS) | Client/Server | Infra CA | 90j auto |
| Kafka brokers + clients | Server + Client | Infra CA | 90j auto |
| Vault nodes | Server + Client | Infra CA | 90j auto |
| PostgreSQL (Patroni) | Server + Client | Infra CA | 90j auto |
| Edge nodes VPN/WireGuard | Client | Infra CA | 30j auto |
| CoreDNS DoT/DoH | Server | Infra CA | 90j auto |
| Management bastions SSH | Host + User | Infra CA + Gov CA | 30j manuel |

### 7.2 Mutual authentication matrix

| Source → Destination | Auth method | Certificate type |
|---------------------|-------------|------------------|
| Citoyen → API SNISID | mTLS + OIDC | Citizen cert + Bearer token |
| Institution → API SNISID | mTLS + OAuth2 | Government cert + service account |
| Service A → Service B (mesh) | mTLS STRICT | Istio sidecar certs (Infra CA) |
| Edge → Core Kafka | mTLS + SASL/SCRAM | Edge node cert (Infra CA) + Kafka user cert |
| Admin → Bastion → K8s | Cert SSH + mTLS | Gov admin cert + Infra node cert |
| Device → Enrollment station | mTLS + TPM attestation | Device cert + station cert |

---

## 8. OCSP & CRL Infrastructure

### 8.1 OCSP Responder HA

| Paramètre | Valeur |
|-----------|--------|
| **Déploiement** | 3 replicas Core DC + 2 replicas DR DC (active-active) |
| **Protocole** | OCSP over HTTP (port 80) + HTTPS (port 443) interne |
| **Cache** | Varnish/Squid 5 minutes (soumis à politique nationale) |
| **Signature** | OCSP responses signées par OCSP Signer dédié (Infra CA subordinate) |
| **Nonce** | Obligatoire pour Citizen/Judicial (anti-replay) |
| **Monitoring** | Prometheus `ocsp_responder_latency`, `ocsp_cache_hit_ratio` |

### 8.2 CRL Distribution

| Paramètre | Valeur |
|-----------|--------|
| **Fréquence** | 1h pour Gov/Infra, 4h pour Citizen, immédiat si révocation d'urgence |
| **Distribution** | CDNs internes (Ceph RGW), HTTP/HTTPS, LDAP (institutions legacy) |
| **Taille max** | Fragmentation si > 10MB. Delta-CRL supporté. |
| **URL format** | `http://crl.snisid.gouv.local/{ca-name}/{date}/{serial}.crl` |

### 8.3 Révocation d'urgence

```
Détection compromission ──► SOC National alert ──► IGC validation ──►
Root CA admin (air-gap laptop) signe CRL emergency ──►
Publication immédiate OCSP + CRL ──►
ArgoCD sync cert-manager force-renew ──►
Notification broadcast (Kafka topic `pki.revocation.emergency`) ──►
Edge nodes sync via burst ──► Fin incident
```

---

## 9. Digital Signature Services

### 9.1 Types de signatures supportées

| Type | Format | Usage | Validation |
|------|--------|-------|------------|
| **Citizen Basic** | CAdES-BES / XAdES-B-B | Consentement, formulaires | Online OCSP |
| **Citizen Advanced** | CAdES-T / XAdES-T | Contrats, démarches | Online OCSP + TSA |
| **Government Official** | PAdES-LTV / CAdES-LT | Actes administratifs | Online OCSP + TSA + archive |
| **Judicial Qualified** | CAdES-LTA / XAdES-LTA | Jugements, assignations | OCSP + TSA + long-term archive |
| **API Integrity** | JWS (RFC 7515) | Webhooks, callbacks | SNISID Infra CA chain |

### 9.2 Timestamp Authority (TSA) nationale

| Paramètre | Valeur |
|-----------|--------|
| **Hardware** | HSM dédié (slot séparé, jamais mixé avec signing) |
| **Précision** | 1 ms maximum (NTP stratum 0 / PTP grandmaster) |
| **Format** | RFC 3161 (Time-Stamp Protocol) |
| **Policy OID** | 1.3.6.1.4.1.SNISID.2.1 (TSA Policy) |
| **Archivage** | Tous les timestamps loggés 25 ans (valeur légale) |
| **Redondance** | TSA Core + TSA DR + TSA satellite (edge, 1h buffer) |

---

## 10. PKI Observability & Audit

### 10.1 Métriques critiques

| Métrique | Seuil | Action |
|----------|-------|--------|
| `cert_expiring_7d` | > 0 | Alert warning — rotation auto |
| `cert_expiring_1d` | > 0 | Alert critical — rotation forcée |
| `cert_expired_active` | > 0 | **EMERGENCY** — service impacté |
| `ocsp_responder_latency_p99` | > 200ms | Scale up + cache tuning |
| `ocsp_failure_rate` | > 0.1% | Investigation SOC + possible attaque |
| `hsm_operation_error_rate` | > 0.01% | Inspection physique HSM |
| `revocation_queue_depth` | > 100 | Escalade IGC + procédure urgence |
| `ca_certificate_issuance_rate` | > 1000/h | Détection anomalie / attaque |

### 10.2 Audit trail

Toute opération PKI est enregistrée :
- **Vault audit log** → Loki → SIEM national (7 ans hot, 25 ans bande)
- **HSM audit log** → Forward temps réel SOC
- **Vidéo cérémonies** → Archive 10 ans air-gapped
- **CRL/OCSP publications** → Immutable log (append-only Ceph RGW)

---

## 11. Gouvernance PKI

### 11.1 Organigramme responsabilités

| Rôle | Entité | Décision |
|------|--------|----------|
| **Root CA Owner** | Présidence / Haute autorité | Cérémonies Root, politique générale, révocation Root |
| **Intermediate CA Owner** | IGC (Inspection Générale Cyber) | Création/revocation Intermediates, audits |
| **Government CA Manager** | Direction numérique | Délivrance institutions, suspensions |
| **Citizen CA Manager** | Agence identité numérique | Délivrance citoyens, portails, consentements |
| **Infrastructure CA Manager** | Équipe Infra Nationale | Délivrance auto services, rotation, révocation auto |
| **Judicial CA Manager** | Conseil supérieur judiciaire | Délivrance magistrats, révocation disciplinaire |
| **Device CA Manager** | Équipe Terrain & IoT | Provisioning équipements, révocation matérielle |
| **Revocation Authority** | IGC + SOC National | Révocations d'urgence, CRL emergency |
| **Audit Authority** | Cour des comptes / Audit externe souverain | Audits annuels, rapports publics restreints |

### 11.2 Certificate Policies (CP) & Certification Practice Statement (CPS)

| Document | OID | Contenu |
|----------|-----|---------|
| SNISID CP Nationale | 1.3.6.1.4.1.SNISID.0.1 | Cadre général, hiérarchie, algorithmes, rôles |
| SNISID CPS Root CA | 1.3.6.1.4.1.SNISID.0.2 | Procédures Root CA, cérémonies, HSM |
| SNISID CPS Government | 1.3.6.1.4.1.SNISID.1.1.1 | Délivrance institutionnelle, validation |
| SNISID CPS Citizen | 1.3.6.1.4.1.SNISID.1.2.1 | Enrollment citoyen, wallet, privacy |
| SNISID CPS Infrastructure | 1.3.6.1.4.1.SNISID.1.3.1 | Auto-issuance, rotation, révocation |
| SNISID CPS Judicial | 1.3.6.1.4.1.SNISID.1.4.1 | Procédures judiciaires, sécurité renforcée |
| SNISID CPS Device | 1.3.6.1.4.1.SNISID.1.5.1 | Provisioning, attestation, cycle de vie |

---

## 12. Standards cryptographiques nationaux

| Domaine | Standard obligatoire | Interdit |
|---------|-------------------|----------|
| **Asymmetric** | ECDSA P-384 (secp384r1), Ed25519 (citizen lightweight) | RSA < 3072, ECDSA P-256 pour Tier-0 |
| **Hash** | SHA-384, SHA-3-256, BLAKE3 (non-cryptographique) | SHA-1, MD5 |
| **Symmetric** | AES-256-GCM (perf), AES-256-CTR + HMAC (FIPS) | 3DES, RC4 |
| **Key exchange** | ECDH P-384, X25519 | DH < 2048 |
| **TLS** | 1.3 uniquement (pas de downgrade) | TLS 1.0, 1.1, SSLv3 |
| **Certificate validity** | Root 10y, Intermediate 5y, End-entity 90j-3y | > 10y Root, > 5y Intermediate, > 3y end-entity |
| **Key rotation** | Infrastructure 90j, Citizen/Gov 1-2y, Root 10y | Aucun certificat sans rotation planifiée |
| **HSM** | FIPS 140-2 Level 3 minimum (Level 4 Root) | Software keys pour signing CA |
| **Random** | TRNG hardware HSM + /dev/random vérifié | PRNG non audité |

---

**Document approuvé pour la fondation de confiance nationale.**  
*Classification: TOP-SECRET — SNISID PKI Nationale*
