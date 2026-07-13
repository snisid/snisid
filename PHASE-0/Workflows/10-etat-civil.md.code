# 📜 SNISID — Domaine État Civil (Workflows complets)

**Document N° :** SNISID-EC-010
**Étape Phase 0 :** 10/16
**Principe :** *L'état civil doit devenir totalement numérique et traçable.*

---

## 1. Contexte Haïtien

L'état civil haïtien souffre de :
- Estimations ONU : ~30-40 % des naissances non enregistrées à temps
- Registres papier vulnérables (incendies, séismes, inondations)
- Délais d'obtention d'actes parfois > 6 mois
- Falsifications fréquentes
- Multiples copies sans source unique de vérité

**SNISID répond** par un **registre national numérique unique** + **5 types de procédures de naissance** + mariage, divorce, décès, adoption.

---

## 2. Cadre Légal de Référence

- **Code Civil haïtien** (à moderniser dans le cadre légal SNISID)
- **Loi sur l'État Civil** (réforme en cours)
- **Décret du 17 mai 1990** sur l'organisation des Offices d'État Civil
- **Convention ONU droits de l'enfant** (déclaration immédiate)
- **ODD 16.9** (identité légale pour tous d'ici 2030)

---

## 3. Workflows de NAISSANCE (5 procédures)

### 3.1 EC-N01 — Naissance Simple
**Conditions :** Déclaration dans les **2 ans** suivant la naissance, par parents ou témoins légaux.

```
1. Identification déclarant (CIN scannée)
2. Saisie données enfant (nom, prénoms, date, lieu, sexe)
3. Identification parents (NIN + acte mariage si applicable)
4. Présence témoins (≥ 2, identifiés)
5. Production éventuelle attestation médicale (MSPP via FHIR)
6. Vérifications automatiques (DMN) :
   - Dates cohérentes
   - Parents existants
   - Pas de doublon enfant
7. Validation Officier État Civil (OEC)
8. Génération NIN enfant
9. Signature électronique OEC (PKI)
10. Génération acte de naissance numérique (PDF/A-3 + QR)
11. Émission événement BirthRegistered sur bus
12. Notifications : parents, CEP (future inscription électorale), MSPP, MENFP
```

**SLA cible :** < 24h.

### 3.2 EC-N02 — Naissance par Reconnaissance
**Conditions :** Père non marié reconnaît un enfant a posteriori.

```
1. Identification père reconnaissant (CIN)
2. Référencement acte de naissance existant
3. Consentement mère (signature ou présence)
4. Si enfant majeur (≥ 16 ans) : consentement enfant requis
5. Validation OEC
6. Génération acte de reconnaissance (annexe à acte naissance)
7. Mise à jour fiche Person de l'enfant (champ father)
8. Signature électronique OEC
9. Émission événement RecognitionRegistered
10. Notification parties
```

### 3.3 EC-N03 — Naissance par Déclaration Tardive
**Conditions :** Déclaration **> 2 ans** après la naissance. Procédure renforcée.

```
1. Demande déposée à l'OEC compétent
2. Production preuves : témoins (≥ 2 majeurs), attestation médicale, baptême, etc.
3. Enquête administrative OEC (15 j max)
4. Avis du Parquet (Procureur)
5. Si accord Parquet → enregistrement direct
6. Si désaccord → renvoi vers jugement (passe en EC-N05)
7. Génération acte avec mention "déclaration tardive"
8. Signature OEC + visa Procureur
9. Émission événement + notifications
```

**SLA cible :** < 60 j.

### 3.4 EC-N04 — Naissance par Décret
**Conditions :** Naturalisation, octroi nationalité par décret présidentiel.

```
1. Référence décret présidentiel (numéro, date, Moniteur)
2. Validation conformité par MJSP
3. Création fiche Person + NIN
4. Génération acte de naissance "par décret"
5. Signature OEC + cachet officiel
6. Notification : intéressé, DIE (passeport éligible)
```

### 3.5 EC-N05 — Naissance par Jugement au rang des Minutes
**Conditions :** Aucune autre voie possible — jugement du Tribunal de Première Instance.

```
1. Requête au TPI compétent (introduit par avocat ou en personne)
2. Audience, témoins, preuves
3. Jugement rendu
4. Transmission expédition jugement à l'OEC du lieu d'origine
5. Transcription au rang des minutes du registre d'état civil
6. Génération acte avec mention "par jugement transcrit"
7. Signature OEC + référence jugement
8. Émission événement + notifications
```

**SLA cible :** dépend tribunal (cible 6 mois max).

---

## 4. Workflow MARIAGE (EC-M01)

```
1. Demande conjointe des futurs époux (en ligne ou guichet)
2. Identification (CIN + NIN)
3. Vérifications :
   - Âge légal (≥ 18 ans, ou ≥ 16 avec autorisation parentale + procureur)
   - Pas de mariage en cours (interrogation registre national)
   - Pas de lien parenté prohibé (DMN)
4. Choix régime matrimonial (séparation/communauté)
5. Publication des bans (10 jours, en ligne + affichage commune)
6. Cérémonie devant OEC + témoins (≥ 2)
7. Signatures électroniques époux + témoins + OEC
8. Génération acte de mariage
9. Mise à jour fiches Person (champ civil_status, spouse)
10. Émission événement MarriageRegistered
```

**Variante EC-M02 :** Mariage religieux transcrit (Église catholique, protestants, vodou enregistré) — vérification autorisation officiant + transcription dans 8 jours.

---

## 5. Workflow DIVORCE

### EC-D01 — Divorce par Consentement Mutuel
```
1. Convention déposée par avocat commun ou avocats séparés
2. Inventaire patrimoine + accord garde enfants
3. Présentation TPI
4. Délai de réflexion (15 j)
5. Confirmation époux
6. Jugement de divorce
7. Transcription en marge de l'acte de mariage
8. Mise à jour fiches Person (civil_status: Divorced)
9. Émission événement DivorceRegistered
```

### EC-D02 — Divorce Contentieux
- Procédure judiciaire complète, jugement TPI, voies de recours
- Transcription identique après définitivité du jugement.

---

## 6. Workflow DÉCÈS (EC-X01)

```
1. Constatation médicale (médecin ou agent santé MSPP)
2. Certificat médical décès numérique (FHIR)
3. Déclaration par proche ou hôpital, dans les 24h
4. Identification défunt (NIN si dispo, sinon enquête)
5. Vérifications automatiques :
   - Personne existante et vivante dans registre
   - Cohérence date/lieu
6. Validation OEC
7. Génération acte de décès (PDF/A-3 + signature PKI)
8. Mise à jour Person.status = Deceased
9. Émission événement DeathRegistered
10. Notifications cascadées :
    - CEP (radiation liste électorale)
    - DGI (clôture obligations fiscales)
    - OFATMA/ONA (pensions, capital décès)
    - Banques (via API consentie)
    - Notaire (si succession)
```

**SLA cible :** < 48h.

---

## 7. Workflow ADOPTION

### EC-A01 — Adoption Simple
- Lien biologique préservé, adoptant ajouté.

### EC-A02 — Adoption Plénière
- Rupture avec famille biologique, nouvelle filiation complète.

**Procédure (commune) :**
```
1. Requête au tribunal compétent
2. Enquête sociale (IBESR — Institut Bien-Être Social)
3. Période probatoire si applicable
4. Jugement d'adoption
5. Transcription par OEC (création nouvel acte naissance pour plénière)
6. Mise à jour Person : parents adoptifs
7. Conservation lien biologique sécurisé (accès restreint)
8. Notifications + génération nouvel acte
```

---

## 8. Sécurité & Intégrité

- Chaque acte généré est **signé PKI** (XAdES-LTA)
- Chaque modification = **événement immuable** (event-sourcing)
- **Watermark + QR code** pour vérification publique
- **Archivage légal** sur stockage WORM (Write Once Read Many)
- **Reconstruction historique** complète possible (audit)

---

## 9. Modèle Offline pour État Civil

Tous les workflows ci-dessus disposent d'une **version offline** sur kit mobile :
- Saisie complète
- Signature locale (HSM portable opérateur)
- Génération provisoire acte avec mention "en attente validation centrale"
- Acte définitif émis après sync + validation finale
- Le citoyen reçoit un récépissé + SMS lors de la validation

---

## 10. KPI État Civil

| KPI | Baseline (2026) | Cible 2030 |
|-----|------------------|------------|
| % naissances enregistrées | ~60-70 % | ≥ 95 % |
| Délai moyen acte naissance | 30+ j | < 24 h |
| Taux d'erreur sur actes | ~10 % | < 1 % |
| Falsifications détectées | ND | tracking actif |
| Digitalisation actes existants | 0 % | ≥ 80 % du stock |

---
*Fin du document — Étape 10/16*
