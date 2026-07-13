# ⚖ WORKFLOWS JUDICIAIRES

> **Phase 3 / Étape 5** — Intégration de la Justice nationale dans SNISID.
> Version : 1.0.0
> Classification : OFFICIEL

---

## 1. Mission

> **Créer l'intégration justice nationale.**
> Tout acte judiciaire impactant l'identité ou l'état civil est :
> - **légalement admissible** devant les tribunaux haïtiens,
> - **signé** par les acteurs habilités (greffier + magistrat),
> - **horodaté** (RFC 3161),
> - **archivé** WORM 30 ans,
> - **traçable** via chaîne de preuve (chain of custody).

---

## 2. Cadre Légal

| Source | Référence |
|--------|-----------|
| **Code de procédure civile** | Haïti |
| **Code de procédure pénale** | Haïti |
| **Loi sur l'organisation judiciaire** | Tribunaux + Cours d'appel + Cour de Cassation |
| **Constitution 1987 amendée** | Droit à un procès équitable + recours |
| **Convention de La Haye** | Apostille |
| **eIDAS-like haïtien** | Signature qualifiée nationale (PKI SNISID) |

---

## 3. Principe d'Admissibilité Légale

Pour qu'un workflow judiciaire SNISID soit **opposable devant un tribunal**, il **doit** garantir :

| Exigence | Mise en œuvre |
|----------|---------------|
| **Authenticité** | Signature PKI du greffier + magistrat |
| **Intégrité** | Hash SHA-384 + chaîne Merkle |
| **Non-répudiation** | TSA RFC 3161 (timestamping qualifié) |
| **Chain of custody** | Trace complète des accès et transmissions |
| **Confidentialité** | Chiffrement at-rest + accès loggé |
| **Conservation** | WORM 30 ans minimum (storage immuable) |
| **Reproductibilité** | Replay possible via Kafka audit |
| **Réversibilité contrôlée** | Annulation = nouvel acte signé (jamais effacement) |

---

## 4. Catalogue des Workflows Judiciaires (5 BPMN)

| Workflow | Fichier | SLA | Description |
|----------|---------|-----|-------------|
| `judicial.validation.act` | `BPMN/Judicial/judicial-validation.v1.0.0.bpmn` | **48h** | Validation judiciaire d'un acte administratif |
| `judicial.suspension.identity` | `BPMN/Judicial/identity-suspension.v1.0.0.bpmn` | **4h** | Blocage identité sur ordre judiciaire |
| `judicial.investigation.fraud` | `BPMN/Judicial/fraud-investigation.v1.0.0.bpmn` | **90j** | Enquête fraude (analyste + parquet + tribunal) |
| `judicial.appeal.management` | `BPMN/Judicial/appeal-management.v1.0.0.bpmn` | **90j** | Gestion appels (Cour d'appel) |
| `judicial.court.integration` | `BPMN/Judicial/court-integration.v1.0.0.bpmn` | **24h** | Synchronisation tribunaux ↔ SNISID |

---

## 5. Détail des Workflows

### 5.1 Judicial Validation (`judicial.validation.act`)

**Objectif :** valider légalement un acte (ex : naissance par jugement, mariage prononcé par tribunal).

**Étapes :**
1. Audit start
2. **Hash + chaîne de preuve** (`evidence.chain.create`)
3. **Greffier — réception** (`court-clerks`)
4. **Magistrat — validation** (`civil-judges`)
5. **Signature qualifiée** (PKI eIDAS-like)
6. **Horodatage RFC 3161** (TSA)
7. **Archivage WORM** (preuve légale 30 ans)
8. Kafka `judicial.validated.v1`

**Garanties :**
- Aucun acte sans greffier ET magistrat
- Aucun acte sans signature + timestamp + WORM
- Chaîne de preuve auditée à tout moment

### 5.2 Identity Suspension (`judicial.suspension.identity`)

**Objectif :** suspension immédiate d'une identité sur ordre judiciaire (mesure conservatoire, mandat, etc.).

**Étapes :**
1. Audit
2. Hash + chain of custody de l'ordre judiciaire
3. **Vérification signature PKI** du greffier + magistrat
4. Gateway : ordre valide ?
   - Non → fin avec rejet (notification au tribunal)
5. **Validation juridique (LVB)** (4h max)
6. Fraud detection (cohérence)
7. **Suspension immédiate** (NIN + carte + credentials)
8. **Publication CRL/OCSP**
9. **Archivage WORM**
10. Signature PKI de l'acte de suspension
11. Kafka `judicial.order.suspension.v1`
12. Notification citoyen + administrations dépendantes

**SLA critique :** 4h (sécurité publique).

### 5.3 Fraud Investigation (`judicial.investigation.fraud`)

**Objectif :** enquête judiciaire en cas de fraude identité confirmée.

**Étapes :**
1. Audit
2. Ouverture dossier (case ID)
3. **Suspension préventive** identité
4. Enquête par analyste fraude + police judiciaire (`fraud-investigators`)
5. Saisine parquet si confirmé (`prosecutors`)
6. Jugement (`civil-judges`)
7. Signature PKI
8. Kafka `judicial.fraud.case.v1`

**Durée :** jusqu'à 90 jours, avec étapes intermédiaires.

### 5.4 Appeal Management (`judicial.appeal.management`)

**Objectif :** gestion d'un appel devant la Cour d'appel.

**Étapes :**
1. Audit
2. Dépôt requête appel par greffier
3. Notification parties (avocats, citoyens, ONI)
4. Audience cour d'appel (`appeal-judges`)
5. Arrêt cour d'appel (signé)
6. Signature PKI
7. Kafka `judicial.appeal.ruled.v1`

**Effet :** peut annuler / modifier / confirmer la décision initiale.

### 5.5 Court Integration (`judicial.court.integration`)

**Objectif :** ingérer les décisions des tribunaux haïtiens dans SNISID en temps réel.

**Étapes :**
1. Audit
2. **Ingestion gRPC** depuis les SI des tribunaux (formats HL7 / JSON / XML)
3. **Vérification signatures greffier** (PKI)
4. Normalisation (schéma SNISID)
5. Routing vers le domaine concerné (civil / pénal / admin)
6. Kafka émission événement métier (`civil.*` ou `judicial.*`)

**Note :** ce workflow est la **passerelle officielle** entre SNISID et les tribunaux.

---

## 6. Acteurs Judiciaires

| Acteur | Rôle | Groupes Camunda |
|--------|------|-----------------|
| **Greffier** | Réception, enregistrement, transmission | `court-clerks` |
| **Magistrat civil** | Décision état civil, naissance, mariage | `civil-judges` |
| **Magistrat appel** | Cour d'appel | `appeal-judges` |
| **Procureur** | Parquet, poursuites | `prosecutors` |
| **Police judiciaire** | Enquête | `fraud-investigators` |
| **Avocat** | Représentation | `lawyers` |
| **Officier juridique ONI** | Vérification juridique interne | `legal-officers` |
| **Ombudsman** | Médiation citoyen | `ombudsman-office` |
| **MJSP** | Coordination ministérielle | `moj-officers` |

---

## 7. Chaîne de Preuve (Chain of Custody)

Chaque acte judiciaire suit une **chaîne immuable** :

```
1. RÉCEPTION  (greffier) → hash + signature greffier
2. VALIDATION (magistrat) → signature magistrat
3. SIGNATURE  (PKI eIDAS) → signature qualifiée nationale
4. HORODATAGE (TSA RFC 3161) → timestamp opposable
5. ARCHIVAGE  (WORM) → impossible à modifier
6. ÉMISSION   (Kafka audit.*) → tracé multi-DC
7. CONSULTATION (audit) → tout accès tracé
```

À tout moment, on peut **prouver mathématiquement** :
- Qui a signé
- Quand (au millième de seconde)
- Que le contenu n'a pas été altéré
- Qui y a accédé depuis

---

## 8. Garde-fous (chaque BPMN judiciaire)

| # | Garde-fou | Obligatoire |
|---|-----------|:-----------:|
| 1 | Versioning SemVer | ✅ |
| 2 | Audit trail Merkle | ✅ |
| 3 | Signature qualifiée greffier + magistrat | ✅ |
| 4 | Horodatage RFC 3161 | ✅ |
| 5 | Archivage WORM 30 ans | ✅ |
| 6 | Event sourcing Kafka | ✅ |
| 7 | Validation juridique LVB | ✅ |
| 8 | Chaîne de preuve | ✅ |
| 9 | Anti-tampering (chaîne Merkle vérifiée) | ✅ |
| 10 | Notification parties | ✅ |

---

## 9. Topics Kafka Judiciaires

| Topic | Description | Rétention |
|-------|-------------|-----------|
| `judicial.flagged.v1` | Drapeau judiciaire posé | 10 ans |
| `judicial.validated.v1` | Acte validé légalement | 10 ans |
| `judicial.case.opened.v1` | Dossier ouvert | 10 ans |
| `judicial.case.closed.v1` | Dossier clôturé | 10 ans |
| `judicial.order.suspension.v1` | Ordre suspension identité | 10 ans |
| `judicial.appeal.filed.v1` | Appel déposé | 10 ans |
| `judicial.appeal.ruled.v1` | Arrêt rendu | 10 ans |
| `judicial.fraud.case.v1` | Dossier fraude | 10 ans |

> Les **audit** correspondants sont conservés **30 ans** sur WORM séparé.

---

## 10. Métriques (SLO Judiciaires)

| Indicateur | Cible |
|------------|-------|
| Disponibilité workflows judiciaires | 99,95 % |
| p95 validation acte | < 24 h |
| p95 suspension identité | < 1 h |
| Délai escalade Sev1 | < 1 min |
| Signatures vérifiées | 100 % |
| Chaîne preuve intacte | 100 % (alerte Sev0 sinon) |
| Audit complet | 100 % |

---

## 11. Interopérabilité Justice ↔ SNISID

```
┌──────────────────┐       gRPC       ┌──────────────────┐
│ Tribunal X       │ ◄──── mTLS ────► │ Court-Integration│
│ (SI tribunal)    │   signatures     │ Service          │
└──────────────────┘                  └────────┬─────────┘
                                               │
                                               ▼
                                    ┌─────────────────────┐
                                    │  Kafka Event Mesh    │
                                    │ (civil.* judicial.*) │
                                    └─────────────────────┘
```

Chaque tribunal dispose :
- d'un certificat PKI dédié (`spiffe://snisid.ht/court/<tribunal>`)
- d'un greffier numérique habilité (signature personnelle)
- d'un point d'API gRPC sécurisé
- de logs audit cross-référencés

---

## 12. Cas spécifiques

### 12.1 Tribunal hors-ligne (région isolée)
- Activation du mode **offline judicial** :
  - Greffier signe localement avec son token HSM
  - Stockage en outbox
  - Sync différée à reconnexion
- Audit local Merkle préservé

### 12.2 Décision urgente (référé liberté)
- Workflow accéléré (SLA 1h max)
- Notification immédiate ONI + suspension provisoire
- Validation rétroactive sous 24h

### 12.3 Crise nationale
- Possibilité de "batch judicial" (ex : libération massive prisonniers après séisme)
- Validation collective signée + audit individuel préservé

### 12.4 Erreur judiciaire (réhabilitation)
- Workflow `identity.un-suspend` invoqué
- Acte de réhabilitation signé
- Notification + indemnité administrative
- **Aucun acte original n'est effacé** (event sourcing)

---

## 13. Conformité Internationale

| Standard | Conformité |
|----------|------------|
| **ETSI EN 319 122** (CAdES) | Signatures CAdES-B-T |
| **RFC 3161** | TSA timestamping |
| **eIDAS** (modèle européen) | Signature qualifiée nationale |
| **ICAO 9303** | Documents de voyage (Phase 4) |
| **Convention La Haye 1961** | Apostille (Phase 4) |
| **Mutual Legal Assistance** | Via Kafka topic `judicial.mla.v1` (Phase 4) |

---

## 14. Gouvernance & Contrôle

- **LVB** (Legal Validation Board) approuve chaque BPMN judiciaire avant prod
- **Cour Supérieure des Comptes** audit annuel
- **Ministère Justice et Sécurité Publique** : coordination institutionnelle
- **Avocats Barreau** : consultation triennale sur procédures
- **Ombudsman national** : recours citoyen

---

## 15. Roadmap

| Version | Contenu |
|---------|---------|
| v1.0 | 5 workflows fondamentaux ✅ |
| v1.1 | Workflow `judicial.unsuspend` (réhabilitation explicite) |
| v1.2 | Intégration Cour de Cassation |
| v1.3 | Workflows pénaux (mandat, libération conditionnelle) |
| v2.0 | Tribunal numérique complet (audiences à distance signées) |

---

**Maintenu par :** Workflow Governance Office + LVB + MJSP + Direction Justice Numérique
