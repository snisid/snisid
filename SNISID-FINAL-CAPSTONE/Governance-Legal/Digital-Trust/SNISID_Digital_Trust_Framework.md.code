# SNISID — CADRE DE CONFIANCE NUMÉRIQUE NATIONALE

**Classification :** CADRE STRATÉGIQUE — CONFIANCE NUMÉRIQUE
**Référence :** SNISID-TRUST-001
**Version :** 1.0
**Date :** 25 mai 2026

---

## 1. OBJECTIF

Le Cadre de Confiance Numérique Nationale établit les mécanismes, standards et processus garantissant que toutes les transactions numériques au sein du SNISID sont juridiquement fiables, vérifiables et non répudiables.

---

## 2. PÉRIMÈTRE

Le cadre couvre l'ensemble des interactions numériques :
- Entre citoyens et l'État
- Entre agences gouvernementales
- Entre systèmes automatisés
- Entre l'État haïtien et des entités étrangères (dans le cadre d'accords)

---

## 3. PILIERS DE CONFIANCE

### 3.1 Signatures Numériques (Digital Signatures)

| Niveau | Technologie | Usage | Valeur Juridique |
|--------|------------|-------|------------------|
| Signature simple | HMAC / mot de passe | Validations internes | Indicative |
| Signature avancée | RSA-4096 / ECC P-384 avec certificat | Documents officiels | Probante |
| Signature qualifiée | Certificat qualifié + dispositif sécurisé (HSM/carte à puce) | Actes légaux | Équivalente manuscrite |
| Cachet serveur | Certificat machine qualifié | Transactions automatisées | Garantie d'origine |

**Standards cryptographiques obligatoires :**
- RSA : 4096 bits minimum
- ECC : P-384 minimum
- Hachage : SHA-384 minimum
- Padding : OAEP / PSS

**Processus de signature :**
```
Document/Transaction
    → Calcul du hash (SHA-384)
    → Signature avec clé privée du signataire
    → Horodatage qualifié (TSA nationale)
    → Inclusion du certificat du signataire
    → Vérification OCSP (certificat non révoqué)
    → Archivage de la preuve de signature
```

### 3.2 Confiance Certificataire (Certificate Trust)

**Architecture PKI Nationale :**
```
Root CA Nationale (Offline, HSM, Haïti)
├── CA Gouvernementale
│   ├── CA Identité (certificats citoyens CNI)
│   ├── CA Agents (certificats agents publics)
│   └── CA Serveurs (certificats systèmes)
├── CA Infrastructure
│   ├── CA Réseau (TLS interne)
│   ├── CA Code Signing (signature de code)
│   └── CA Horodatage (TSA)
└── CA Interopérabilité
    └── Cross-certification (accords bilatéraux)
```

| Composant | Standard | Localisation | Sécurité |
|-----------|---------|-------------|----------|
| Root CA | X.509 v3, RSA-4096 | HSM FIPS 140-2 L3, Haïti | Cérémonie de clé multi-parties |
| CA intermédiaires | X.509 v3, RSA-4096 | HSM en data center national | Accès multi-contrôle |
| Certificats finaux | X.509 v3, RSA-2048+ ou ECC P-256+ | Carte à puce / logiciel | PIN/biométrie |
| CRL | X.509 CRL v2 | Publication nationale | Signée, horodatée |
| OCSP | RFC 6960 | Serveurs nationaux redondants | Haute disponibilité |

**Politique de certification :**
- Vérification d'identité en personne pour certificats qualifiés
- Vérification NNI + biométrie pour certificats citoyens
- Renouvellement automatique avec ré-authentification
- Révocation en moins de 1 heure

### 3.3 Confiance Identitaire (Identity Trust)

| Niveau d'Assurance | Méthode | Usage |
|--------------------|---------|-------|
| LOA 1 — Basique | NNI + mot de passe | Consultation d'informations |
| LOA 2 — Substantiel | NNI + OTP (SMS/app) | Services en ligne standards |
| LOA 3 — Élevé | CNI électronique + PIN | Services sensibles |
| LOA 4 — Très élevé | CNI + biométrie + PIN | Transactions légales, signatures qualifiées |

**Protocoles d'authentification :**
- SAML 2.0 pour fédération d'identité inter-agences
- OpenID Connect pour services en ligne citoyens
- FIDO2/WebAuthn pour authentification forte
- Certificats X.509 pour authentification machine-to-machine

**Fédération d'identité nationale :**
- Identity Provider (IdP) national centralisé
- Single Sign-On (SSO) inter-agences
- Politique d'attributs standardisée
- Consentement citoyen pour partage d'attributs

### 3.4 Intégrité Transactionnelle (Transaction Integrity)

| Mécanisme | Description | Application |
|-----------|-------------|------------|
| Hash d'intégrité | SHA-384 sur chaque transaction | Toutes transactions |
| Signature transactionnelle | Signature numérique de chaque opération | Transactions sensibles |
| Horodatage qualifié | TSA nationale RFC 3161 | Toutes transactions |
| Journalisation immutable | Append-only, chiffrée, signée | Toutes transactions |
| Merkle Tree | Arbre de hash pour groupes de transactions | Lots de transactions |
| Numéro de séquence | Détection de transactions manquantes | Toutes transactions |

**Garanties :**
- Toute transaction est datée avec précision (NTP synchronisé, horodatage qualifié)
- Toute transaction est attribuable à un acteur identifié
- Toute modification est détectable
- Toute suppression est impossible (journaux immutables)
- Tout rejeu est détecté (nonce + séquence)

### 3.5 Non-Répudiation

| Type | Mécanisme | Preuve |
|------|-----------|--------|
| Non-répudiation d'origine | Signature numérique de l'émetteur | Certificat + timestamp |
| Non-répudiation de réception | Accusé de réception signé | Certificat + timestamp |
| Non-répudiation de soumission | Horodatage qualifié à la soumission | TSA + hash |
| Non-répudiation de transport | Journaux réseau signés | Logs + signatures |
| Non-répudiation de stockage | Hash de stockage vérifié | Hash + audit trail |

**Protocole de non-répudiation :**
```
1. Émetteur signe le message/document
2. Système horodatage qualifie la soumission
3. Transport chiffré et journalisé
4. Destinataire vérifie signature + certificat + horodatage
5. Destinataire émet accusé de réception signé
6. Tout est archivé dans le journal immutable
→ Aucune partie ne peut nier l'échange
```

---

## 4. SERVICES DE CONFIANCE

| Service | Description | Disponibilité |
|---------|-------------|--------------|
| Service de signature | API de signature numérique | 24/7, 99.99% |
| Service de vérification | Vérification de signatures et certificats | 24/7, 99.99% |
| Service d'horodatage (TSA) | Horodatage qualifié RFC 3161 | 24/7, 99.99% |
| Service OCSP | Vérification statut certificat | 24/7, 99.99% |
| Service d'archivage | Conservation à long terme des preuves | 24/7, 99.9% |
| Service de validation | Validation complète (signature + certificat + horodatage) | 24/7, 99.99% |

---

## 5. GOUVERNANCE DE LA CONFIANCE

### 5.1 Comité de Confiance Numérique
- Revue trimestrielle des politiques de confiance
- Approbation des changements cryptographiques
- Gestion des incidents de confiance
- Planification de la migration cryptographique (post-quantique)

### 5.2 Audits de Confiance
- Audit annuel de la PKI par auditeur qualifié
- Tests de pénétration semestriels sur les services de confiance
- Revue annuelle des algorithmes et tailles de clé
- Exercice de révocation d'urgence annuel

---

## 6. PLAN DE MIGRATION POST-QUANTIQUE

| Phase | Horizon | Action |
|-------|---------|--------|
| Veille | 2026-2027 | Surveillance des standards NIST post-quantiques |
| Préparation | 2027-2028 | Évaluation des algorithmes candidats |
| Pilote | 2028-2029 | Implémentation hybride (classique + PQ) |
| Migration | 2029-2031 | Migration progressive |
| Complétion | 2031+ | Infrastructure post-quantique complète |

---

*Document cadre préparé dans le cadre de la Phase 14 — SNISID National Legal Framework*
