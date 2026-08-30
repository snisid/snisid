# National Workflow Registry (Catalogue Officiel des Processus)
## Registre de Gouvernance et de Traçabilité des Flux Étatiques — SNISID v4.0

---

## 1. OBJECTIF & CADRE DE GOUVERNANCE
Le **National Workflow Registry** est le registre unique d'autorité répertoriant tous les workflows BPMN 2.0 certifiés et approuvés pour exécution au sein de l'État Digital SNISID. Aucun workflow ne peut être exécuté s'il n'est pas préalablement enregistré, validé pour conformité de schéma, et signé cryptographiquement par le bureau de gouvernance de la Direction Nationale du Numérique.

---

## 2. SPÉCIFICATIONS DES FLUX ENREGISTRÉS

### 2.1 Domaine : État Civil (Namespace: `civil.registry`)

#### 2.1.1 Naissance Simple
- **Process ID** : `snisid-civil-birth-simple`
- **Version Active** : `1.2.0` (Date de certification: 2026-01-15)
- **Propriétaire (Ownership)** : Direction de l'État Civil (Ministère de l'Intérieur)
- **SLA Légal** : 4 Heures
- **Niveau de Classification** : CONFIDENTIEL (PII - Données d'identité nominatives)
- **Événements Déclencheurs** : `birth.created`
- **Description** : Processus d'enregistrement standard d'une naissance biologique déclarée en maternité ou centre de santé agréé sous 3 jours.

#### 2.1.2 Naissance par Reconnaissance
- **Process ID** : `snisid-civil-birth-recognition`
- **Version Active** : `2.0.1` (Date de certification: 2026-03-10)
- **Propriétaire (Ownership)** : Direction de l'État Civil / Direction des Affaires Juridiques
- **SLA Légal** : 8 Heures
- **Niveau de Classification** : CONFIDENTIEL
- **Événements Déclencheurs** : `birth.recognition.initiated`
- **Description** : Permet la reconnaissance filiale volontaire hors mariage. Intègre la validation des consentements parentaux signés cryptographiquement.

#### 2.1.3 Naissance par Déclaration Tardive
- **Process ID** : `snisid-civil-birth-late-declaration`
- **Version Active** : `1.4.0` (Date de certification: 2026-02-28)
- **Propriétaire (Ownership)** : Direction de l'État Civil & Officiers d'État Civil Municipaux
- **SLA Légal** : 24 Heures
- **Niveau de Classification** : CONFIDENTIEL
- **Événements Déclencheurs** : `birth.late_declaration.initiated`
- **Description** : Processus applicable après le délai légal d'enregistrement de naissance. Exige une audition des témoins et une validation administrative renforcée.

#### 2.1.4 Naissance par Décret
- **Process ID** : `snisid-civil-birth-decree`
- **Version Active** : `1.0.0` (Date de certification: 2025-11-12)
- **Propriétaire (Ownership)** : Secrétariat Général de la Présidence & Ministère de la Justice
- **SLA Légal** : 48 Heures
- **Niveau de Classification** : CONFIDENTIEL
- **Événements Déclencheurs** : `birth.decree.published`
- **Description** : Intégration à l'état civil suite à l'attribution de nationalité par décret officiel ou naturalisation.

#### 2.1.5 Naissance par Jugement au Rend des Minutes
- **Process ID** : `snisid-civil-birth-court-judgment`
- **Version Active** : `1.1.2` (Date de certification: 2026-04-05)
- **Propriétaire (Ownership)** : Greffe Civil des Tribunaux de Première Instance
- **SLA Légal** : 96 Heures
- **Niveau de Classification** : CONFIDENTIEL
- **Événements Déclencheurs** : `birth.judgment.registered`
- **Description** : Régularisation d'état civil ordonnée par jugement d'un tribunal civil compétent.

---

### 2.2 Domaine : Identité Biométrique (Namespace: `identity.registry`)

#### 2.2.1 Enrôlement National (Enrollment)
- **Process ID** : `snisid-identity-enrollment`
- **Version Active** : `3.1.0` (Date de certification: 2026-02-10)
- **Propriétaire (Ownership)** : Office National de l'Identité (ONI)
- **SLA Légal** : 8 Heures
- **Niveau de Classification** : CONFIDENTIEL (Données biométriques sensibles)
- **Événements Déclencheurs** : `identity.enrollment.started`
- **Description** : Capture d'empreintes, iris, photographie faciale et raccordement à l'acte de naissance validé.

#### 2.2.2 Vérification (Verification)
- **Process ID** : `snisid-identity-verification`
- **Version Active** : `2.1.5` (Date de certification: 2026-05-01)
- **Propriétaire (Ownership)** : ONI - Division Déduplication Biométrique
- **SLA Légal** : 2 Heures
- **Niveau de Classification** : CONFIDENTIEL
- **Événements Déclencheurs** : `identity.verified`
- **Description** : Comparaison un-à-plusieurs (1:N) via l'ABIS national pour éliminer les doublons et usurpations d'identité.

#### 2.2.3 Révocation & Suspension Urgente
- **Process ID** : `snisid-identity-revocation`
- **Version Active** : `1.0.1` (Date de certification: 2025-10-20)
- **Propriétaire (Ownership)** : Direction de la Sécurité ONI / Police Nationale
- **SLA Légal** : 30 Minutes (Priorité P1 - Critique)
- **Niveau de Classification** : CONFIDENTIEL
- **Événements Déclencheurs** : `identity.revocation.initiated`
- **Description** : Révocation immédiate d'un titre d'identité ou de l'état d'identité d'un individu en cas de fraude avérée, perte signalée ou décès constaté.

---

### 2.3 Domaine : Justice & Police (Namespace: `justice.registry` & `police.registry`)

#### 2.3.1 Émission de Mandat d'Arrêt (Warrant)
- **Process ID** : `snisid-justice-warrant-arrest`
- **Version Active** : `2.0.0` (Date de certification: 2026-03-22)
- **Propriétaire (Ownership)** : Ministère de la Justice / Cabinet d'Instruction
- **SLA Légal** : 30 Minutes (Priorité P1 - Critique)
- **Niveau de Classification** : SECRET-DÉFENSE
- **Événements Déclencheurs** : `judicial.warrant.issued`
- **Description** : Émission et diffusion instantanée d'un mandat d'arrêt vers toutes les forces de police aux frontières (DCPJ, Polifront).

#### 2.3.2 Arrestation & Garde à Vue (GAV)
- **Process ID** : `snisid-police-arrest-gav`
- **Version Active** : `1.8.2` (Date de certification: 2026-04-18)
- **Propriétaire (Ownership)** : Direction Centrale de la Police Judiciaire (DCPJ)
- **SLA Légal** : 24 Heures strictes (Priorité P1 - Limite Constitutionnelle)
- **Niveau de Classification** : SECRET-DÉFENSE
- **Événements Déclencheurs** : `police.arrestation.created`
- **Description** : Orchestration et surveillance en temps réel de la durée légale de rétention en cellule de police. Déclenche une alarme d'urgence rouge vers le parquet si aucune prorogation n'est validée à T-1 heure.

---

## 3. PROCESSUS D'HOMOLOGATION D'UN NOUVEAU WORKFLOW
Pour qu'un nouveau workflow soit déployé et enregistré dans ce registre, il doit valider :
1. **La validation syntaxique BPMN 2.0** sans boucle infinie.
2. **Le schéma CloudEvents** pour tous les messages publiés et consommés.
3. **La définition explicite des SLAs** (T1 et T2) et de la stratégie de repli d'urgence.
4. **La clé cryptographique** d'authentification de l'administration propriétaire.
5. **Le certificat de signature** du Code Review de la Direction du Numérique de l'État.
