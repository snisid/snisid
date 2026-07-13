# SNISID — National PKI Standards
**Classification:** TOP-SECRET / RESTREINT DEFENSE  
**Version:** 5.0.0  
**Statut:** OBLIGATOIRE — Toute violation bloque issuance et signale incident IGC

---

## 1. Standards cryptographiques obligatoires

### 1.1 Algorithmes asymétriques

| Usage | Algorithme | Courbe / Taille | Interdit |
|-------|-----------|-----------------|----------|
| Root CA signing | ECDSA | secp384r1 (P-384) | RSA < 4096, ECDSA P-256, P-521 sans justification |
| Intermediate CA signing | ECDSA | secp384r1 (P-384) | RSA < 4096, ECDSA P-256 |
| End-entity certificates | ECDSA | secp384r1 (P-384) | RSA < 3072, DSA, ECDSA P-256 pour Tier-0 |
| Mesh mTLS (high volume) | ECDSA | secp256r1 (P-256) **autorisé exception** | RSA, secp384r1 (performance critique uniquement) |
| Citizen wallet (mobile) | ECDSA / Ed25519 | secp384r1 / Ed25519 | RSA, secp256r1 pour signatures juridiques |
| Post-quantique (prep) | CRYSTALS-Dilithium / Falcon | NIST PQC Level 3+ | Aucun avant validation IGC |

### 1.2 Algorithmes de hachage

| Usage | Algorithme | Interdit |
|-------|-----------|----------|
| Signature certificates | SHA-384 | SHA-1, MD5, SHA-256 (Root/Intermediate uniquement) |
| OCSP responses | SHA-384 | SHA-1 |
| CRL generation | SHA-384 | SHA-1 |
| TSA timestamping | SHA-384 | SHA-1, MD5 |
| Document integrity (signatures) | SHA-384 / SHA-3-256 | SHA-1, MD5 |
| Mesh mTLS (performance) | SHA-256 **exception** | SHA-1, MD5 |

### 1.3 Chiffrement symétrique

| Usage | Algorithme | Mode / Tailles | Interdit |
|-------|-----------|---------------|----------|
| Data at rest (HSM, backups) | AES-256 | GCM (AEAD) | 3DES, RC4, ECB |
| TLS session encryption | AES-256 | GCM (TLS 1.3 native) | CBC modes, RC4 |
| Citizen wallet (PFX) | AES-256 | CBC + HMAC-SHA384 | ECB, CTR without MAC |
| Key wrapping (HSM) | AES-256 | Key Wrap (RFC 3394 / NIST SP 800-38F) | None |

### 1.4 Échange de clés

| Usage | Algorithme | Interdit |
|-------|-----------|----------|
| TLS 1.3 handshake | ECDH X25519 / secp384r1 | DH < 2048, RSA key exchange |
| WireGuard | Curve25519 | None (Noise Protocol) |
| Key agreement (HSM) | ECDH secp384r1 | DH < 2048 |

### 1.5 Génération de nombres aléatoires

| Source | Exigence |
|--------|----------|
| HSM key generation | TRNG hardware certifié (FIPS 140-2) |
| K8s / application CSPRNG | /dev/urandom avec ChaCha20 (Linux 4.8+) ou getrandom() |
| Citizen device | Secure Enclave / TEE / TPM TRNG |
| Seed material | Jamais de seed logiciel unique — pooling multiple sources |

---

## 2. Standards de certificats X.509

### 2.1 Durées de validité maximales

| Type | Durée max | Renouvellement | Alerte |
|------|-----------|----------------|--------|
| Root CA | 10 ans | Reconstruction cérémonie | 2 ans anticipation |
| Intermediate CA | 5 ans | Cérémonie HSM | 1 an anticipation |
| Government end-entity | 3 ans | Workflow BPMN annuel | 60 jours |
| Citizen end-entity | 2 ans | Re-enrollment / auto | 90 jours |
| Infrastructure end-entity | 90 jours | Auto cert-manager | 30 jours |
| Judicial end-entity | 3 ans | Workflow BPMN annuel | 60 jours |
| Device end-entity | 1 an | Auto provisioning | 30 jours |
| Mesh mTLS | 30 jours | Auto Istio SDS | 7 jours |
| OCSP responder | 1 an | Auto cert-manager | 30 jours |
| TSA certificate | 5 ans | Cérémonie HSM | 1 an |

### 2.2 Formats et extensions obligatoires

| Extension | Root CA | Intermediate | End-Entity | Interdit |
|-----------|---------|--------------|------------|----------|
| BasicConstraints CA=TRUE | ✅ | ✅ | ❌ (pathlen 0) | ❌ sur end-entity |
| KeyUsage keyCertSign, cRLSign | ✅ Root only | ✅ | ❌ | ❌ sur end-entity |
| KeyUsage digitalSignature | ❌ | ✅ | ✅ | — |
| KeyUsage nonRepudiation | ❌ | ❌ (sauf Judicial) | ✅ (Judicial/Citizen adv) | — |
| ExtKeyUsage serverAuth | ❌ | ✅ (Infra) | ✅ (Infra) | ❌ Gov/Judicial |
| ExtKeyUsage clientAuth | ❌ | ✅ | ✅ | — |
| ExtKeyUsage codeSigning | ❌ | ✅ (Gov/Judicial) | ✅ (Gov/Judicial) | ❌ Infra/Citizen basic |
| ExtKeyUsage emailProtection | ❌ | ✅ (Citizen opt-in) | ✅ (Citizen opt-in) | — |
| CertificatePolicies OID | ✅ | ✅ | ✅ | ❌ |
| CRLDistributionPoint | ❌ | ✅ | ✅ | ❌ Root |
| AuthorityInfoAccess (OCSP) | ❌ | ✅ | ✅ | ❌ Root |
| SubjectKeyIdentifier | ✅ | ✅ | ✅ | ❌ |
| AuthorityKeyIdentifier | ✅ | ✅ | ✅ | ❌ Root self |
| SubjectAltName | ❌ | ❌ | ✅ (selon profil) | ❌ Root/Intermediate |
| NameConstraints | ✅ (optionnel) | ✅ (optionnel) | ❌ | — |

### 2.3 Subject naming standards

| Profil | Format CN | Format OU | Format O | Interdit dans subject |
|--------|-----------|-----------|----------|----------------------|
| Root CA | `SNISID Root CA Nationale` | — | `État` | Tout PII |
| Government | `{Institution} — {Service}` | `{Service code}` | `{Ministère}` | Noms personnels |
| Citizen | `{citizen-uuid}` | `Citizen Identity` | `SNISID` | Nom, NIF, adresse, email |
| Infrastructure | `{service}.{namespace}.svc.cluster.local` | `{tier}-{region}` | `SNISID Infra` | IPs publiques |
| Judicial | `{Matricule} — {Function}` | `{Tribunal}` | `{Cour}` | Nom complet (matricule seul) |
| Device | `{device-uuid}` | `{device-class}-{tier}` | `SNISID Device` | Localisation physique |

---

## 3. Standards TLS / mTLS nationaux

### 3.1 TLS versions et suites

| Couche | Version obligatoire | Versions interdites |
|--------|---------------------|---------------------|
| Internet public (API citoyen) | TLS 1.3 | TLS 1.0, 1.1, SSLv3, SSLv2 |
| Intra-cluster (Istio mesh) | TLS 1.3 | TLS 1.2 même interne |
| HSM / Vault / SOC | TLS 1.3 | Toutes versions < 1.3 |
| OCSP / CRL / TSA | TLS 1.3 (HTTPS), HTTP acceptable (OCSP legacy) | SSLv3 |
| Citizen wallet sync | TLS 1.3 | Toutes versions < 1.3 |
| Edge node VPN | WireGuard (Noise) | IPsec IKEv1, PPTP, L2TP sans IPsec |

### 3.2 Cipher suites autorisées TLS 1.3

```
TLS_AES_256_GCM_SHA384      ✅ (Default national)
TLS_CHACHA20_POLY1305_SHA256  ✅ (Mobile/performance)
TLS_AES_128_GCM_SHA256       ❌ (Interdit national — 128 bits insuffisant pour souveraineté)
TLS_AES_128_CCM_SHA256       ❌
TLS_AES_128_CCM_8_SHA256     ❌
```

### 3.3 mTLS exigences

| Communication | Client cert | Server cert | Vérification |
|-------------|-------------|-------------|--------------|
| Citoyen → API publique | Citizen CA (sensible) ou OAuth2 (basique) | Infra CA | OCSP + cert chain |
| Institution → API interne | Gov CA | Infra CA | OCSP + SPIFFE + AuthZ |
| Service → Service (mesh) | Istio SDS (Infra CA mesh-mtls) | Istio SDS | SPIFFE ID + AuthZ |
| Kubelet → API Server | Infra CA (node-auth) | Infra CA (API server) | CN + SAN IP |
| etcd peer → etcd peer | Infra CA (etcd-peer) | Infra CA (etcd-peer) | SAN hostname |
| Edge → Core Kafka | Infra CA (device/node) | Infra CA (Kafka) | SASL/SCRAM + mTLS |
| Bastion → Node SSH | Gov CA (admin) ou Infra CA (host) | Infra CA (host) | Certificate + CA pinning |
| Vault consumer → Vault | Infra CA (service) | Infra CA (Vault) | namespace + AuthZ policy |
| HSM client → HSM | HSM partition cert | HSM server cert | TLS 1.3 + mutual |

---

## 4. Standards HSM et sécurité physique

### 4.1 HSM requirements

| Propriété | Root CA | Intermediate CAs | End-entity signing |
|-----------|---------|-------------------|-------------------|
| FIPS 140-2 level | Level 4 (tamper-responsive + environmental) | Level 3 minimum | Level 3 ou TEE/Secure Enclave |
| Tamper evidence | Physical + electrical + temperature | Physical + electrical | Minimum physical |
| Key extraction | IMPOSSIBLE (design) | IMPOSSIBLE | Hardware-bound |
| M-of-N | 4-of-6 | 2-of-4 | N/A (service accounts) |
| Network | Air-gap total | Isolated VLAN Tier-0 | PKCS#11 proxy / SDS |
| Backup | Shamir SSS 4-of-6 | HSM DR replication + Shamir | Auto-issued (replaceable) |
| Zeroization | Auto on tamper | Auto on tamper | Manual decommission |

### 4.2 Salle blanche PKI (Root CA ceremony)

| Exigence | Valeur |
|----------|--------|
| Access | Biometric 2-factor + smartcard + armed guard + registry |
| Video | 4K, 24/7, 10 years retention, dual power + UPS |
| RF shielding | TEMPEST level (Faraday cage) |
| Climate | 18-24°C, 40-60% RH, positive pressure, HEPA filtration |
| Fire | FM-200 + water mist backup, pre-action sprinklers |
| Power | Dual utility + UPS 4h + diesel 72h + solar emergency |
| Network | Aucun — console série locale uniquement |
| Maintenance | Présence IGC obligatoire, pré-announcement 72h |

---

## 5. Standards de gestion du cycle de vie (CLM)

### 5.1 Issuance

| Type | Validation | Approval | Audit |
|------|-----------|----------|-------|
| Root CA | 6 custodians biométrie | Présidence + IGC | Vidéo + logs HSM + registre |
| Intermediate | 4 custodians HSM | IGC | Logs HSM + SIEM |
| Government | HR + manager + IGC check | Gov CA Manager | SIEM + BPMN workflow |
| Citizen | Biometric + ID doc + agent | Automated (enrollment) | SIEM + consent log |
| Infrastructure | Terraform / Helm / CI | Automated (cert-manager) | CI logs + SIEM |
| Judicial | Nomination + serment + collège | Judicial Council + IGC | SIEM + BPMN + vidéo (option) |
| Device | TPM attestation + registry | Automated (provisioning) | SIEM + attestation log |
| Mesh mTLS | SPIFFE + service account | Automated (Istio SDS) | Envoy access logs |

### 5.2 Renewal / Rotation

| Type | Auto | Window | Méthode |
|------|------|--------|---------|
| Root CA | ❌ | N/A | Reconstruction cérémonie (10 ans) |
| Intermediate | ❌ | 12 mois | HSM ceremony + new key pair |
| Government | ❌ | 60 jours | BPMN workflow + re-authentication |
| Citizen | ❌ (notification) | 90 jours | Re-enrollment or online consent + biometric |
| Infrastructure | ✅ | 30 jours | cert-manager + Vault PKI auto |
| Judicial | ❌ | 6 mois | BPMN workflow + re-authentication |
| Device | ✅ | 30 jours | Provisioning service auto + attestation |
| Mesh mTLS | ✅ | 7 jours | Istio SDS auto-reload |

### 5.3 Revocation

| Type | Auto | Timeline | Notification |
|------|------|----------|--------------|
| Infrastructure compromise | ✅ (SOC) | < 5 min | SIEM + IGC + ArgoCD sync |
| Device compromise | ✅ (SOC/Wazuh) | < 15 min | SIEM + Device CA Manager |
| Citizen fraud/theft | ❌ (ticket) | < 1h | SMS + email + app push |
| Government dismissal | ❌ (HR ticket) | < 2h | Manager + HR + SIEM |
| Judicial discipline | ❌ (Council ticket) | < 4h | Judicial Council + SIEM |
| Mass revocation (CA compromise) | ✅ + manuel | < 1h | Présidence + IGC + national broadcast |

---

## 6. Standards OCSP et CRL

### 6.1 OCSP

| Propriété | Valeur |
|-----------|--------|
| Protocole | OCSP over HTTP (port 80) + HTTPS (port 443 interne) |
| Réponse | DER-encoded OCSPResponse, signed ECDSA P-384 SHA-384 |
| Nonce | Requis pour Citizen, Gov, Judicial (anti-replay) |
| Cache | 5 minutes max (good responses), pas de cache revoked |
| HA | 3 Core + 2 DR, load-balancer internal |
| SLA | 99.99%, p99 < 200ms |

### 6.2 CRL

| Propriété | Valeur |
|-----------|--------|
| Format | X.509 v2 CRL, DER, signed ECDSA P-384 SHA-384 |
| Fréquence | 1h (Gov/Infra), 4h (Citizen), immédiat si emergency |
| Distribution | HTTP/HTTPS + LDAP (legacy institutions), CDN interne Ceph |
| Delta-CRL | Activé pour Citizen/Device (volume élevé) |
| Taille max | 10 MB → fragmentation ou delta-CRL |
| Staleness edge | 72h max (offline grace period) |

---

## 7. Standards de signatures numériques

### 7.1 Formats et profils

| Usage | Format | Niveau | Timestamp | LTV | Archive |
|-------|--------|--------|-----------|-----|---------|
| Citizen basic | CAdES-BES / XAdES-B-B / JAdES-B-B | Basique | Optionnel | ❌ | 7 ans |
| Citizen advanced | CAdES-T / XAdES-T / JAdES-T | Avancée | Requis | ❌ | 7 ans |
| Government official | PAdES-LTV / CAdES-LT / XAdES-LT | Officielle | Requis | ✅ | 10 ans |
| Judicial qualified | CAdES-LTA / XAdES-LTA / PAdES-LTA | Qualifiée | Requis | ✅ | 25 ans |
| API integrity | JWS (detached) | Technique | iat claim | ❌ | 2 ans |

### 7.2 TSA (Timestamp Authority)

| Propriété | Valeur |
|-----------|--------|
| Précision | 1 ms |
| Source horloge | GNSS discipliné + atomic clock national |
| Protocole | RFC 3161 |
| Signature | ECDSA P-384 SHA-384 |
| Redondance | Core + DR + edge buffer |
| Rétention | 25 ans |

---

## 8. Conformité et audit

### 8.1 Audits obligatoires

| Fréquence | Type | Réalisé par | Portée |
|-----------|------|-------------|--------|
| Trimestriel | HSM ops + CLM metrics | IGC interne | Tous les CAs, expirations, révocations |
| Semestriel | PKI penetration test | Équipe Red Team nationale (souversaine) | Vault PKI, cert-manager, Istio mTLS, OCSP |
| Annuel | External audit | Cabinet audité souverain / alliance technique | Root CA ceremony review, HSM physical, CPS compliance |
| Mensuel | DR drill | Équipe Infra Nationale | Basculement CRL/OCSP, HSM DR sync, cert rotation |

### 8.2 Documentation obligatoire

| Document | Mise à jour | Approbation | Access |
|----------|-------------|-------------|--------|
| CPS Root CA | Après chaque cérémonie | Présidence + IGC | Custodians + IGC |
| CPS Intermediates | Annuel | IGC | CA Managers + IGC + auditeurs |
| CP Nationale | Biennal | Présidence + Parlement équivalent | Public restreint (institutions) |
| Standards PKI (ce doc) | Annuel | IGC | Tous les ops PKI |
| Runbooks PKI | Post-incident | IGC | SOC + CA Managers + Infra |

---

*Standards validés par IGC. Toute exception nécessite dérogation signée Présidence + IGC + justification cryptographique.*
