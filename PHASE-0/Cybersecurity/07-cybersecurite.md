# 🛡️ SNISID — National Cybersecurity Framework

**Document N° :** SNISID-SEC-007
**Étape Phase 0 :** 7/16
**Principe :** *La sécurité doit être intégrée dès le départ.*

---

## 1. Doctrine

SNISID adopte les standards :
- **NIST Cybersecurity Framework 2.0** (Identify, Protect, Detect, Respond, Recover, Govern)
- **ISO/IEC 27001 / 27002 / 27017 / 27018**
- **MITRE ATT&CK** pour la détection de menaces
- **Zero Trust Architecture** (NIST SP 800-207)
- **OWASP ASVS** pour les applications

---

## 2. Piliers (6)

### 2.1 GOVERN
- CISO national rattaché à l'AND
- Politique de Sécurité des SI de l'État (PSSIE-HT)
- Comité Cyber Crise activable 1h
- Reporting mensuel au CNN

### 2.2 IDENTIFY
- Inventaire continu des actifs (CMDB)
- Cartographie des flux de données
- Analyses de risques par actif critique (méthode EBIOS RM ou OCTAVE)
- Classification des données (Public / Interne / Confidentiel / Secret)

### 2.3 PROTECT
- **Zero Trust** : mTLS partout, identity-aware proxy, microsegmentation
- **IAM/PAM** : Keycloak + CyberArk-like (ex. Teleport OSS, HashiCorp Boundary)
  - MFA obligatoire pour tout accès admin
  - Just-In-Time access avec workflow d'approbation
  - Session recording pour comptes à privilèges
- **PKI nationale** : Autorité Racine offline (HSM), AC intermédiaires (HSM en ligne)
- **Chiffrement** : AES-256-GCM, TLS 1.3, KMS souverain
- **Hardening** : CIS Benchmarks systématiques (Linux, K8s, PostgreSQL, etc.)
- **Backups 3-2-1-1-0** : 3 copies, 2 supports, 1 hors-site, 1 immuable (offline), 0 erreur restoration

### 2.4 DETECT
- **SIEM** national (Wazuh ou Elastic Security) consolidant logs de toutes les agences
- **SOAR** pour automatiser réponses (playbooks)
- **EDR** sur tous endpoints serveurs et postes admin
- **NDR** sur points de sortie réseau
- **Threat Intelligence** : abonnement CTI + partage avec CSIRT régionaux (CARICOM, OEA)
- **UEBA** : détection comportements anormaux

### 2.5 RESPOND
- **CSIRT-HT** national 24/7
- Playbooks d'incident : ransomware, exfiltration, DDoS, fraude interne
- Procédure d'isolation < 15 min
- Communication de crise (média, citoyens, partenaires)
- Forensics : chaîne de preuve maîtrisée

### 2.6 RECOVER
- **RTO/RPO** définis par service :
  - Services critiques (identité, état civil) : RTO 1h, RPO 5 min
  - Services secondaires : RTO 4h, RPO 1h
- DRP testé semestriellement (basculement réel PaP → Cap-Haïtien)
- Tabletop exercises trimestriels

---

## 3. National SOC (Security Operations Center)

```
┌──────────────────────────────────────────────┐
│  SOC NATIONAL — Tier 1 / 2 / 3              │
├──────────────────────────────────────────────┤
│  Tier 1 — Surveillance 24/7 (3x8)            │
│  Tier 2 — Analyse approfondie + IR          │
│  Tier 3 — Threat hunting + Forensics + Dev   │
├──────────────────────────────────────────────┤
│  Outils : SIEM, SOAR, EDR, NDR, TIP, MISP   │
│  Liens : CSIRT, ANSI, Interpol, FIRST.org   │
└──────────────────────────────────────────────┘
```

**Effectif cible (2028) :** 24 analystes (8 par shift) + 5 lead + 1 CISO + 1 head SOC.

---

## 4. PKI Nationale

```
        Root CA (offline, HSM, scellée)
              │
   ┌──────────┼──────────┬──────────────┐
   │          │          │              │
Citizens   Officials   Devices       Documents
   ICA       ICA         ICA            ICA
   │          │          │              │
Cert CIN   Cert Agent  Cert IoT      Cert Sign
```

- Conformité **eIDAS** (signature qualifiée)
- HSM **FIPS 140-2 niveau 3** ou Common Criteria EAL4+
- Cérémonie de clés filmée, multi-témoins (key ceremony)
- CRL + OCSP nationaux

---

## 5. IAM (Identity & Access Management)

- **IdP unique** : Keycloak (clusters HA)
- **Standards** : OIDC, SAML 2.0, OAuth 2.1
- **MFA obligatoire** : TOTP, FIDO2/WebAuthn pour admins
- **SSO** entre toutes les agences via fédération
- **Provisioning automatisé** : SCIM 2.0
- **Revocation** sous 15 minutes pour départ/incident

---

## 6. PAM (Privileged Access Management)

- Bastion souverain (Teleport / WALLIX / HashiCorp Boundary OSS)
- Pas de SSH direct vers production
- Enregistrement vidéo + commandes de toute session admin
- Approvals multi-personnes pour actions critiques

---

## 7. Sécurité Applicative

- **SAST** intégré CI (SonarQube, Semgrep)
- **DAST** sur recette (OWASP ZAP, Burp)
- **SCA** pour dépendances (Trivy, Snyk OSS)
- **Container scan** (Trivy, Grype)
- **Secrets scanning** (Gitleaks)
- **Pentest** annuel par cabinet externe + bug bounty restreint

---

## 8. Sécurité Données

- Chiffrement at-rest : LUKS + chiffrement applicatif champs sensibles
- Tokenisation NIN dans bases secondaires
- Pseudonymisation systématique hors prod
- **DLP** sur sorties (mail, USB, cloud)
- **Watermarking** sur extractions de données

---

## 9. Sécurité Physique des Datacenters

- TIER III minimum (Uptime Institute)
- Contrôle d'accès biométrique + carte + PIN
- Vidéosurveillance H24
- Détection intrusion + incendie + inondation
- Garde armé permanent
- Sas mantrap

---

## 10. Conformité & Audit

| Audit | Fréquence | Auditeur |
|-------|-----------|----------|
| ISO 27001 | Annuel | Externe certifié |
| Pentest infrastructure | Annuel | Externe |
| Pentest applicatif | Bi-annuel | Externe + interne |
| Code review sécurité | Continu | Interne + SAST |
| Audit PKI | Annuel | Cabinet spécialisé |
| Audit conformité données (NDPA) | Annuel | NDPA |
| Audit financier sécurité | Annuel | Cour des Comptes |

---

## 11. KPI Cybersécurité

| KPI | Cible |
|-----|-------|
| MTTD (Mean Time To Detect) | < 1 h |
| MTTR (Mean Time To Respond) | < 4 h |
| Vulnérabilités critiques corrigées | < 7 j |
| Patches systèmes appliqués | < 30 j |
| Couverture MFA admin | 100 % |
| Taux de réussite phishing simulation | < 5 % |
| Incidents majeurs / an | 0 |

---
*Fin du document — Étape 7/16*
