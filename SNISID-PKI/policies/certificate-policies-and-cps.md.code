# SNISID — Certificate Policies (CP) & Certification Practice Statements (CPS)
**Classification:** TOP-SECRET / RESTREINT DEFENSE  
**Version:** 5.0.0  
**Statut:** OFFICIEL — Cadre de confiance nationale SNISID

---

## Vue d'ensemble

Les **Certificate Policies (CP)** et **Certification Practice Statements (CPS)** constituent le cadre juridique, organisationnel et technique de la PKI nationale SNISID. Ils définissent :

- Les obligations de l'État en tant qu'autorité de certification souveraine
- Les procédures de délivrance, renouvellement, suspension et révocation des certificats
- Les exigences de sécurité physique, logique et cryptographique
- Les mécanismes de validation (OCSP, CRL) et de non-répudiation
- Les responsabilités des acteurs : citoyens, institutions, équipementiers, opérateurs

**Règle absolue :** Aucun certificat SNISID n'est délivré sans référence à une Certificate Policy approuvée par l'IGC et publiée dans le registre national officiel.

---

## Structure des politiques nationales

### Hiérarchie des OIDs SNISID PKI

```
1.3.6.1.4.1.SNISID          # Racine OID SNISID (à remplacer par PEN IANA attribué)
├── 0                        # PKI Framework
│   ├── 0.1                  # SNISID CP Nationale (Certificate Policy nationale)
│   └── 0.2                  # SNISID CPS Root CA (Certification Practice Statement Root)
│
├── 1                        # Intermediate CA Policies
│   ├── 1.1                  # SNISID Government Policy
│   │   └── 1.1.1            # SNISID CPS Government
│   ├── 1.2                  # SNISID Citizen Policy
│   │   └── 1.2.1            # SNISID CPS Citizen
│   ├── 1.3                  # SNISID Infrastructure Policy
│   │   └── 1.3.1            # SNISID CPS Infrastructure
│   ├── 1.4                  # SNISID Judicial Policy
│   │   └── 1.4.1            # SNISID CPS Judicial
│   │   └── 1.4.2            # SNISID Judicial Qualified Signature Policy
│   │   └── 1.4.3            # SNISID Judicial Smartcard Policy
│   └── 1.5                  # SNISID Device Trust Policy
│       └── 1.5.1            # SNISID CPS Device Trust
│       └── 1.5.2            # SNISID Edge Node Policy
│       └── 1.5.3            # SNISID Biometric Device Policy
│       └── 1.5.4            # SNISID Mobile Unit Policy
│
├── 2                        # Services complémentaires
│   ├── 2.1                  # SNISID TSA Policy (Timestamp Authority)
│   └── 2.2                  # SNISID TSA Judicial Policy
│
├── 3                        # Attributs étendus nationaux
│   ├── 3.1                  # SNISID Device Metadata Extension
│   └── 3.2                  # SNISID Audit Log Reference Extension
│
├── 4                        # Identifiants citoyens anonymisés
│   ├── 4.1                  # SNISID Citizen UUID Namespace
│   ├── 4.2                  # SNISID Enrollment Transaction ID
│   └── 4.3                  # SNISID Biometric Template Hash Reference
│
└── 5                        # Signature profiles
    ├── 5.1                  # SNISID Citizen Authentication Signature
    └── 5.2                  # SNISID Citizen Legal Signature
```

---

## Certificate Policies nationales

### CP-000 : SNISID CP Nationale (1.3.6.1.4.1.SNISID.0.1)

| Section | Contenu | Classification |
|---------|---------|----------------|
| **Introduction** | Portée, objectifs, conformité légale nationale, reconnaissance internationale | RESTREINT DEFENSE |
| **Responsabilités générales** | Engagement de l'État, obligations des souscripteurs, déclin de responsabilité (sauf faute lourde souveraine) | RESTREINT DEFENSE |
| **Identification et authentification** | Procédures d'identification physique, biométrique, documentaire pour chaque catégorie | SECRET |
| **Exigences opérationnelles** | Délibération, renouvellement, re-délivrance, modification, suspension, révocation | RESTREINT DEFENSE |
| **Exigences physiques, procédurales et de personnel** | Contrôles physiques, opérationnels, contrôles d'accès, séparation des fonctions | TOP-SECRET |
| **Exigences techniques** | Génération de clés, tailles, algorithmes, protection des clés privées, cycle de vie | TOP-SECRET |
| **Gestion du cycle de vie des certificats** | Durées de validité, extensions obligatoires, profiles X.509 | RESTREINT DEFENSE |
| **Contrôle des révoqués** | CRL, OCSP, raisons de révocation, délais de publication | RESTREINT DEFENSE |
| **Services d'audit et de non-répudiation** | Archivage, horodatage, preuve juridique | TOP-SECRET |
| **Responsabilités des parties** | CA, RA, souscripteurs, utilisateurs finaux, équipementiers | RESTREINT DEFENSE |

### CP-101 : SNISID Government Policy (1.3.6.1.4.1.SNISID.1.1)

| Exigence | Valeur | Justification |
|----------|--------|---------------|
| **Identification** | Pièce officielle + vérification par HR institution + IGC | Double contrôle institutionnel |
| **Authentification** | Smartcard national ou certificat précédent + OTP | Multi-facteur |
| **Délivrance** | Workflow BPMN (WF-GOV-001) | Traçabilité |
| **Renouvellement** | Annuel, re-validation HR + IGC | Rotation obligatoire |
| **Révocation** | Révocation par manager + IGC sous 24h (cessationOfOperation) | Rapidité en cas de départ |
| **Non-répudiation** | Signature codeSigning pour actes officiels | Valeur légale |
| **Audit** | Logs SIEM 10 ans + journal papier institutions 25 ans | Durée légale archives |

### CP-201 : SNISID Citizen Policy (1.3.6.1.4.1.SNISID.1.2)

| Exigence | Valeur | Justification |
|----------|--------|---------------|
| **Identification** | Pièce nationale + capture biométrique (iris/empreinte) + agent terrain validé | Identité fortement vérifiée |
| **Authentification** | Wallet mobile (Secure Enclave) + PIN/Pattern + biométrie locale | Triple facteur |
| **Délivrance** | Post-enrollment automatique (Citizen CA Manager supervisé) | Rapidité + traçabilité |
| **Renouvellement** | Biennal, notification 90j avant, re-enrollment simplifié | Usabilité + sécurité |
| **Révocation** | Citoyen (portail/terrain), agent (fraude), Justice (ordonnance), décès (état civil) | Multi-canal |
| **Non-répudiation** | digitalSignature + nonRepudiation | Signature électronique avancée |
| **Privacy** | UUID dans certificat — jamais de nom, NIF, adresse, email sans opt-in | RGPD-équivalent |
| **Audit** | Logs anonymisés 7 ans, consentement traçé 7 ans | Minimisation |

### CP-301 : SNISID Infrastructure Policy (1.3.6.1.4.1.SNISID.1.3)

| Exigence | Valeur | Justification |
|----------|--------|---------------|
| **Identification** | Service account K8s / Terraform inventory / provisioning UUID | Automatique |
| **Authentification** | Service mesh mTLS + SPIFFE ID + Istio SDS | Machine-to-machine |
| **Délivrance** | 100% automatisée — cert-manager + Vault PKI | Zero trust ops |
| **Renouvellement** | Tous les 90 jours, auto-renewal 30j avant | Short-lived, compromission limitée |
| **Révocation** | Auto (Falco/Wazuh/Tetragon) + manuel (SOC) | Immédiat si compromise |
| **Non-répudiation** | N/A (pas de valeur juridique humaine) | Technique uniquement |
| **Audit** | Logs SIEM 2 ans + Prometheus metrics | Opérationnel |

### CP-401 : SNISID Judicial Policy (1.3.6.1.4.1.SNISID.1.4)

| Exigence | Valeur | Justification |
|----------|--------|---------------|
| **Identification** | Nomination officielle + serment + vérification Conseil supérieur | Autorité judiciaire |
| **Authentification** | Smartcard judiciaire HSM-backed (Level 3) + biométrie agent | Qualifiée |
| **Délivrance** | Workflow BPMN (WF-JUD-001) avec double contrôle juridique | Cérémonial |
| **Renouvellement** | Triennal, re-validation Conseil + IGC | Magistrature |
| **Révocation** | Conseil supérieur + IGC sous 4h (privilegeWithdrawn / cessationOfOperation) | Discipline rapide |
| **Non-répudiation** | digitalSignature + nonRepudiation + CAdES-LTA | Signature qualifiée équivalente |
| **Audit** | Logs SIEM 25 ans + vidéo (si disponible) + journal tribunal | Durée légale maximale |
| **Chain of custody** | Biométrie agent + ID station + timestamp TSA + hash document | Preuve irréfutable |

### CP-501 : SNISID Device Trust Policy (1.3.6.1.4.1.SNISID.1.5)

| Exigence | Valeur | Justification |
|----------|--------|---------------|
| **Identification** | TPM EK / Manufacturer IDevID / Secure Boot attestation | Hardware-rooted |
| **Authentification** | Attestation quote + measured boot + firmware version | Intégrité |
| **Délivrance** | Provisioning automatisé post-attestation | Zero touch |
| **Renouvellement** | Annuel, auto si attestation valide | Maintenance |
| **Révocation** | Auto (CVE / anomalie réseau / vol) + manuel (décommission) | Temps réel |
| **Non-répudiation** | digitalSignature uniquement | Authentification machine |
| **Audit** | Logs SIEM 5 ans + device registry | Forensics |

---

## Certification Practice Statements (CPS)

### CPS-000 : SNISID CPS Root CA (1.3.6.1.4.1.SNISID.0.2)

**Classification : TOP-SECRET**

Le CPS Root CA est le document le plus sensible de la PKI nationale. Il décrit :

- **La cérémonie de génération de clé** : salle blanche, 6 custodians, biométrie, vidéo 24/7, Shamir 4-of-6
- **La topologie HSM** : Thales Luna 7 Level 4, partitions isolées, air-gap total
- **Les procédures de backup** : HSM DR replication offline, Shamir physical escrow (6 sites), LTO air-gapped
- **La révocation Root CA** : décision présidentielle, cérémonie de reconstruction, notification internationale
- **Les audits** : trimestriel IGC, annuel externe souverain, vidéo archivée 10 ans

*Document complet stocké physiquement et numériquement (chiffré) dans le coffre IGC et la salle blanche Root CA. Accès : Custodians Root CA + Directeur IGC + Présidence uniquement.*

### CPS-101 : SNISID CPS Government (1.3.6.1.4.1.SNISID.1.1.1)

**Classification : RESTREINT DEFENSE**

- **Délivrance** : BPMN workflow, 2 validateurs (manager + IGC), smartcard personnalisation
- **Renouvellement** : workflow raccourci (validation automatique HR si contrat actif + manager approval)
- **Révocation** : HR trigger → SOC vérification → Government CA Manager execution → CRL 1h
- **HSM** : Partition SNISID_GOV_CA, accès Vault agent PKI, 2-of-4 M-of-N

### CPS-201 : SNISID CPS Citizen (1.3.6.1.4.1.SNISID.1.2.1)

**Classification : SECRET**

- **Délivrance** : Enrollment station (biométrie + ID) → validation backoffice (IA + agent) → Citizen CA auto-issuance
- **Renouvellement** : Notification 90j → portail ou terrain → re-authentication biométrie simplifiée → issuance
- **Révocation** : Portail citoyen (biométrie + PIN) / agent terrain (formulaire + preuve) / Justice (ordonnance)
- **Privacy** : UUID seul dans certificat, consentement granularisé par usage, droit à l'oubli partiel (logs 7 ans)
- **Wallet** : PKCS#12 (Secure Enclave/TPM), chiffrement AES-256 PBKDF2 200k itérations, recovery codes papier

### CPS-301 : SNISID CPS Infrastructure (1.3.6.1.4.1.SNISID.1.3.1)

**Classification : RESTREINT DEFENSE**

- **Délivrance** : cert-manager + Vault PKI (100% automated), validation CI/CD + Kyverno policy
- **Renouvellement** : cert-manager auto-renew 30j avant, Istio SDS transparent reload
- **Révocation** : SOC auto (Falco/Wazuh) → Vault revoke → ArgoCD sync → CRL 1h
- **HSM** : Partition SNISID_INFRA_CA, auto-unseal Vault, jamais d'intervention humaine en production
- **Mesh mTLS** : Istio SDS + SPIFFE, certificats 30j, rotation transparente, deny-all default

### CPS-401 : SNISID CPS Judicial (1.3.6.1.4.1.SNISID.1.4.1)

**Classification : TOP-SECRET**

- **Délivrance** : Conseil supérieur validation → IGC background check → double contrôle (CA Manager + Conseil membre) → HSM 2-person → smartcard personalization 2-person
- **Renouvellement** : Re-validation Conseil triennale, même cérémonial
- **Révocation** : Conseil supérieur + IGC sous 4h, CRL 1h, notification tribunaux
- **Non-répudiation** : CAdES-LTA, TSA national, chain of custody biométrique, 25 ans archive
- **HSM** : Partition SNISID_JUDICIAL_CA, 2-of-4 M-of-N, salle blanche judicial IT center

### CPS-501 : SNISID CPS Device Trust (1.3.6.1.4.1.SNISID.1.5.1)

**Classification : RESTREINT DEFENSE**

- **Provisioning** : Secure Boot → measured boot → TPM attestation → registry check → Vault auto-issuance → 802.1X admission
- **Renouvellement** : Auto-refresh annuel si attestation validée, CVE scan clean, registry active
- **Révocation** : Auto (CVE 9+ unpatched, anomalie Cilium, Falco privesc) → immediate 802.1X disable + WireGuard peer remove
- **HSM** : Partition SNISID_DEVICE_CA, auto-issuance via provisioning service, revocation auto via Wazuh integration

---

## Publication et distribution

| Document | Version | Dernière mise à jour | Prochaine révision | Diffusion |
|----------|---------|---------------------|---------------------|-----------|
| CP-000 Nationale | 5.0.0 | 2026-05-25 | 2028-05-25 | Institutions + IGC + auditeurs |
| CPS-000 Root CA | 5.0.0 | 2026-05-25 | Root CA reconstruction | Custodians + Présidence + IGC |
| CP-101 Government | 5.0.0 | 2026-05-25 | 2027-05-25 | Ministères + CA Manager Gov + IGC |
| CPS-101 Government | 5.0.0 | 2026-05-25 | 2027-05-25 | CA Manager Gov + IGC + auditeurs |
| CP-201 Citizen | 5.0.0 | 2026-05-25 | 2027-05-25 | Agence identité + CA Manager Citizen + IGC |
| CPS-201 Citizen | 5.0.0 | 2026-05-25 | 2027-05-25 | CA Manager Citizen + IGC + DPO équivalent |
| CP-301 Infrastructure | 5.0.0 | 2026-05-25 | 2027-05-25 | Équipe Infra + CA Manager Infra + IGC |
| CPS-301 Infrastructure | 5.0.0 | 2026-05-25 | 2027-05-25 | Équipe Infra + IGC + auditeurs |
| CP-401 Judicial | 5.0.0 | 2026-05-25 | 2027-05-25 | Conseil supérieur + CA Manager Judicial + IGC |
| CPS-401 Judicial | 5.0.0 | 2026-05-25 | 2027-05-25 | Conseil supérieur + CA Manager Judicial + IGC |
| CP-501 Device | 5.0.0 | 2026-05-25 | 2027-05-25 | Équipe Terrain + CA Manager Device + IGC |
| CPS-501 Device | 5.0.0 | 2026-05-25 | 2027-05-25 | Équipe Terrain + CA Manager Device + IGC |

---

## Révisions et gouvernance

### Processus de révision

1. **Proposition** : IGC ou CA Manager identifie besoin de révision (CVE, évolution légale, audit)
2. **Rédaction** : Groupe de travail IGC + CA Manager concerné + juriste national
3. **Revue interne** : Peer review au sein IGC (4 reviews minimum)
4. **Revue externe** : Audit externe souverain + cabinet juridique national
5. **Approbation** :
   - CP Nationale : Présidence + Haute autorité
   - CP/CPS Intermédiaires : IGC + CA Manager + juriste
6. **Publication** : Registre national officiel + distribution aux parties prenantes
7. **Entrée en vigueur** : 90 jours après publication (sauf urgence sécurité — 7 jours)

### Traçabilité des versions

Chaque CPS/CP est :
- **Signé numériquement** : Government Official signature (CAdES-LTV) par IGC Director
- **Horodaté** : TSA national (RFC 3161)
- **Archivé** : Ceph RGW WORM bucket, 25 ans, immutable
- **Diffusé** : Git national (repository sécurisé) + papier signé (IGC archive)

---

*Documents CP/CPS complets stockés dans le vault physique IGC et le registre national numérique (accès RBAC RESTREINT DEFENSE).*  
*Ce référentiel contient les spécifications techniques, OIDs, et mappings. Les documents juridiques complets sont gouvernés par la procédure de révision ci-dessus.*
