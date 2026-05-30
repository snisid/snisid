# Glossaire Technique SNISID — v1.0 (MP-013)

> **Document officiel** — Validé par le Conseil Technique SNISID  
> **Règle** : Toute modification requiert validation du Conseil Technique avant intégration  
> **Hash SHA-256** : `[Calculé lors de la cérémonie Root CA]`

---

## Métadonnées

| Champ | Valeur |
|---|---|
| Version | 1.0 |
| Date création | 2026-05-24 |
| Termes | 20 |
| Standards | 11 |
| Référence | Master Prompt MP-013 |

---

## Termes Techniques

| Terme | Définition | v | MP |
|---|---|---|---|
| **ABIS** | Automated Biometric Identification System — Moteur déduplication 1:N GPU cluster. Seuils: <85% Nouveau, 85-95% Tier-2, ≥95% Hard Match→DCPJ | v1.0 | MP-005, MP-006 |
| **ABAC** | Attribute-Based Access Control — Contrôle accès par attributs contextuels via OPA/Rego en sidecar | v1.0 | MP-001, MP-004 |
| **CockroachDB** | Base SQL distribuée ACID, multi-région, compatible PostgreSQL — Réplication DC1↔DC2 | v1.0 | MP-007, MP-010 |
| **CRDT** | Conflict-free Replicated Data Type — Résolution automatique des conflits offline MEK | v1.0 | MP-008, MP-010 |
| **EJBCA** | Enterprise JavaBeans Certificate Authority — PKI open-source, Issuing CA eID (~500K certs/an) | v1.0 | MP-009 |
| **FIDO2** | Fast Identity Online 2 (W3C WebAuthn) — Auth sans mot de passe via YubiKey (AAL3) | v1.0 | MP-001 |
| **HSM** | Hardware Security Module — Protection matérielle clés crypto. FIPS 140-2 L3 minimum requis | v1.0 | MP-009 |
| **MEK** | Mobile Enrollment Kit — Kit biométrique ruggedisé terrain. Solaire 200W, 72h autonomie, IP67 | v1.0 | MP-008, MP-010 |
| **mTLS** | Mutual TLS — Auth mutuelle entre services (Zero Trust pilier 3) via Istio + SPIFFE/SPIRE | v1.0 | MP-001, MP-003 |
| **NATS JetStream** | Message broker léger orienté edge — Queuing offline MEK, resync automatique | v1.0 | MP-008, MP-010 |
| **NNI** | Numéro National d'Identification — ID unique pérenne citoyen haïtien. Format: HTI-AAAA-DPT-NNNNN | v1.0 | MP-002, MP-005, MP-006 |
| **ONI** | Office National d'Identification — Agence nationale identité, opérateur SNISID | v1.0 | MP-002, MP-011 |
| **OPA** | Open Policy Agent (CNCF Graduated) — Moteur politique ABAC, langage Rego, sidecar K8s | v1.0 | MP-001, MP-004 |
| **PAD** | Presentation Attack Detection — Anti-spoofing biométrique (deepfakes, faux doigts, photos) | v1.0 | MP-005, MP-006 |
| **RKE2** | Rancher Kubernetes Engine 2 — Distribution K8s sécurisée FIPS-compatible, CIS L2 | v1.0 | MP-003, MP-007 |
| **RPO** | Recovery Point Objective — Perte données maximale acceptable. **SNISID : < 1 minute** | v1.0 | MP-010 |
| **RTO** | Recovery Time Objective — Temps rétablissement maximum. **SNISID : < 15 minutes** | v1.0 | MP-010 |
| **SPIFFE** | Secure Production Identity Framework — Identités cryptographiques microservices via SPIRE | v1.0 | MP-001, MP-003 |
| **WORM** | Write-Once-Read-Many — Stockage immuable audit logs légaux. Rétention 10 ans Cold Tier | v1.0 | MP-004 |
| **X-Road** | Standard interopérabilité estonien — Échange données inter-agences sécurisé (30+ pays) | v1.0 | MP-003, MP-011 |

---

## Standards Internationaux

| Standard | Domaine | Applicabilité SNISID | Phase |
|---|---|---|---|
| NIST SP 800-63-3 | Digital Identity Guidelines | IAM complet AAL1/2/3 | Phase 1 |
| NIST SP 800-207 | Zero Trust Architecture | Architecture réseau 7 piliers | Phase 1 |
| ISO/IEC 27001:2022 | Information Security Management | Certification complète | Phase 5 |
| ISO/IEC 27701:2019 | Privacy Information Management | Données biométriques citoyens | Phase 2 |
| ISO/IEC 19794-2 | Biometric Data Format Fingerprint | ABIS + MEK templates | Phase 2 |
| FIPS 140-2 | Cryptographic Modules Security | HSM L3 min + PKI + MEK | Phase 1 |
| X-Road Protocol | Estonian Interoperability | Échanges inter-agences | Phase 3 |
| OWASP API Top 10 | API Security Guidelines | API Gateway Kong WAF | Phase 1 |
| MITRE ATT&CK | Adversary Tactics & Techniques | SOC + SIEM playbooks | Phase 2 |
| CNCF Landscape | Cloud Native Technologies | Stack technique complète | Phase 1 |
| SLSA Framework | Software Supply Chain Security | DevSecOps SLSA L3→4 | Phase 2 |

---

## Scores de Maturité — Projections M36

| Domaine | Actuel | M36 | Gain |
|---|:---:|:---:|:---:|
| Architecture Globale & Microservices | 87/100 | 97/100 | **+10** |
| Cybersécurité & Zero Trust | 85/100 | 96/100 | **+11** |
| IAM National | 83/100 | 95/100 | **+12** |
| Biométrie ABIS | 81/100 | 94/100 | **+13** |
| SOC/SIEM/SOAR | 79/100 | 93/100 | **+14** |
| PKI & HSM | 82/100 | 96/100 | **+14** |
| Offline-First & Résilience | 88/100 | 98/100 | **+10** |
| DevSecOps/GitOps | 76/100 | 93/100 | **+17** |
| Gouvernance des Données | 74/100 | 90/100 | **+16** |
| Disaster Recovery PRA/PCA | 80/100 | 96/100 | **+16** |
| **SCORE GLOBAL MOYEN** | **81.5/100** | **94.8/100** | **+13.3** |

---

## Recommandations MP-013

### 📌 Versioning Terminologique Formel
- Numéro de version `vX.Y` sur chaque terme
- Date de mise à jour obligatoire
- Validation Conseil Technique avant toute modification
- Hash SHA-256 signé par le Secrétaire du Conseil

### 🔗 Liens Croisés vers Modules
- Matrice Terme × Master Prompt maintenue à jour
- Références bidirectionnelles (glossaire → MP, MP → glossaire)
- Navigation rapide dans les 55+ fichiers de référence

---

*SNISID_Glossaire_Technique_v1.0 — © République d'Haïti — Document CONFIDENTIEL*
