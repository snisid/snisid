# SNISID — RUNBOOKS LÉGAUX & COMPLIANCE

**Classification :** OPÉRATIONNEL — RUNBOOKS RÉGLEMENTAIRES
**Référence :** SNISID-RUN-LEG-001
**Version :** 1.0
**Date :** 25 mai 2026

---

## 1. OBJECTIF

Standardiser et industrialiser les réponses aux événements réglementaires et juridiques affectant le SNISID, garantissant des réponses rapides, cohérentes et documentées.

---

## RUNBOOK 1 — DATA BREACH DISCLOSURE (Notification de Violation de Données)

### 1.1 Déclencheur
- Violation de données personnelles détectée ou suspectée
- Accès non autorisé confirmé à des données personnelles
- Fuite de données avérée ou potentielle
- Notification par un tiers d'une exposition de données

### 1.2 Classification de Sévérité

| Niveau | Critères | Délai d'Action |
|--------|---------|---------------|
| CRITIQUE | Données biométriques / identité exposées, >10 000 personnes | Immédiat |
| ÉLEVÉ | Données personnelles sensibles, >1 000 personnes | 2 heures |
| MOYEN | Données personnelles basiques, <1 000 personnes | 4 heures |
| FAIBLE | Données non sensibles, exposition limitée | 24 heures |

### 1.3 Procédure

```
PHASE 1 — DÉTECTION ET CONFINEMENT (0-4h)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
□ 1. Alerte reçue par le SOC / CERT-HT
□ 2. Évaluation initiale de la nature et de l'étendue
□ 3. Classification de sévérité
□ 4. Confinement immédiat (isolation système, blocage accès)
□ 5. Notification interne :
    → Directeur BNC
    → RSSI
    → DPD (Délégué à la Protection des Données)
    → Directeur de l'agence concernée
□ 6. Activation de l'équipe de réponse

PHASE 2 — INVESTIGATION (4-24h)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
□ 7. Investigation forensique : cause, étendue, données affectées
□ 8. Identification des personnes concernées
□ 9. Évaluation du risque pour les personnes
□ 10. Documentation complète de l'incident
□ 11. Collecte de preuves forensiques (chain of custody)

PHASE 3 — NOTIFICATION RÉGLEMENTAIRE (24-72h)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
□ 12. Préparation du rapport de notification ANPD
    Contenu obligatoire :
    - Nature de la violation
    - Catégories et nombre de personnes concernées
    - Catégories et volume de données affectées
    - Conséquences probables
    - Mesures prises et envisagées
    - Coordonnées du DPD
□ 13. Soumission à l'ANPD (dans les 72 heures)
□ 14. Si risque élevé : notification aux personnes concernées
    Contenu :
    - Description claire de la violation
    - Conseils de protection (changement de mot de passe, etc.)
    - Mesures prises par le SNISID
    - Coordonnées pour questions
□ 15. Notification au CNGN

PHASE 4 — REMÉDIATION (72h-30 jours)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
□ 16. Correction de la vulnérabilité exploitée
□ 17. Renforcement des contrôles
□ 18. Vérification que la brèche est fermée
□ 19. Monitoring renforcé (30 jours)
□ 20. Rapport complet à l'ANPD (30 jours)
□ 21. Rapport au CNGN

PHASE 5 — POST-INCIDENT (30-90 jours)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
□ 22. Retour d'expérience (post-mortem)
□ 23. Mise à jour des procédures
□ 24. Formation complémentaire si nécessaire
□ 25. Rapport final au CNGN et à l'ANPD
□ 26. Archivage de l'ensemble du dossier
```

### 1.4 Contacts Clés

| Rôle | Contact | Disponibilité |
|------|---------|--------------|
| SOC Manager | [Coordonnées] | 24/7 |
| DPD | [Coordonnées] | Heures ouvrées + astreinte |
| Directeur BNC | [Coordonnées] | Heures ouvrées + astreinte |
| ANPD (urgences) | [Coordonnées] | 24/7 |
| Conseiller juridique | [Coordonnées] | Heures ouvrées + astreinte |

---

## RUNBOOK 2 — LEGAL INJUNCTION (Injonction Juridique)

### 2.1 Déclencheur
- Réception d'une ordonnance judiciaire
- Réquisition de données par un magistrat
- Injonction de cesser un traitement
- Ordonnance de gel de données (legal hold)

### 2.2 Procédure

```
RÉCEPTION (0-2h)
━━━━━━━━━━━━━━━
□ 1. Réception de l'acte judiciaire
□ 2. Vérification de l'authenticité :
    → Signature du magistrat vérifiée
    → Juridiction compétente vérifiée
    → Sceau officiel vérifié
□ 3. Enregistrement dans le registre des actes judiciaires
□ 4. Notification immédiate :
    → Directeur BNC
    → Conseiller juridique
    → Directeur de l'agence concernée
    → DPD

ANALYSE (2-24h)
━━━━━━━━━━━━━━
□ 5. Analyse juridique de l'injonction :
    → Portée exacte (quelles données, quel périmètre)
    → Délai de conformité
    → Obligations et restrictions
    → Impact opérationnel
□ 6. Évaluation de la faisabilité technique
□ 7. Identification des données/systèmes concernés
□ 8. Avis juridique (conformité constitutionnelle, proportionnalité)

EXÉCUTION (selon délai)
━━━━━━━━━━━━━━━━━━━━━━
□ 9. Si injonction de production de données :
    → Extraction par agent habilité
    → Procès-verbal de collecte
    → Hash d'intégrité
    → Transmission sécurisée au tribunal
    → Décharge de réception
□ 10. Si injonction de gel (legal hold) :
    → Application immédiate du flag LEGAL_HOLD
    → Copie de sauvegarde supplémentaire
    → Notification aux équipes techniques
    → Suspension de toute destruction programmée
□ 11. Si injonction de cessation de traitement :
    → Arrêt du traitement concerné
    → Documentation de l'arrêt
    → Évaluation de l'impact sur les services
    → Plan de continuité si nécessaire
□ 12. Si mandat de perquisition numérique :
    → Coopération avec les forces de l'ordre
    → Présence du conseiller juridique
    → Documentation de la perquisition
    → Procès-verbal contradictoire

DOCUMENTATION (post-exécution)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
□ 13. Rapport d'exécution au tribunal
□ 14. Mise à jour du registre des actes judiciaires
□ 15. Notification au CNGN si impact significatif
□ 16. Archivage complet du dossier
```

### 2.3 Cas Particulier : Contestation

Si l'injonction est jugée disproportionnée ou inconstitutionnelle :
1. Avis juridique documenté
2. Consultation du CNGN (urgence)
3. Recours devant la juridiction compétente
4. Exécution provisoire sauf en cas de dommage irréparable
5. Documentation de la contestation

---

## RUNBOOK 3 — COMPLIANCE FAILURE (Défaut de Conformité)

### 3.1 Déclencheur
- Non-conformité détectée par audit interne ou externe
- Constatation par le BNC
- Signalement par un lanceur d'alerte
- Rapport de l'ANPD ou de l'ANC

### 3.2 Procédure

```
DÉTECTION ET ÉVALUATION (0-48h)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
□ 1. Identification de la non-conformité
□ 2. Classification de sévérité :
    → CRITIQUE : violation légale active, risque immédiat
    → MAJEURE : écart significatif, risque élevé
    → SIGNIFICATIVE : écart notable, risque moyen
    → MINEURE : amélioration souhaitable
□ 3. Documentation de la non-conformité
□ 4. Notification à l'agence concernée

ESCALADE (selon sévérité)
━━━━━━━━━━━━━━━━━━━━━━━━
□ 5. CRITIQUE → Directeur BNC + CNGN immédiat + ANPD si données
    → Mesures d'urgence (arrêt du traitement si nécessaire)
□ 6. MAJEURE → Directeur BNC + Directeur agence
    → Plan d'action sous 7 jours
□ 7. SIGNIFICATIVE → Chef division BNC + Responsable SNISID agence
    → Plan d'action sous 30 jours
□ 8. MINEURE → Analyste BNC + Point de contact agence
    → Plan d'action sous 90 jours

REMÉDIATION
━━━━━━━━━━━
□ 9. L'agence soumet un plan de remédiation
□ 10. Validation du plan par le BNC
□ 11. Exécution du plan
□ 12. Vérification par le BNC
    → Conforme → Clôture
    → Non conforme → Escalade au niveau supérieur
□ 13. Documentation et leçons apprises

SUIVI
━━━━━
□ 14. Inscription au registre des non-conformités
□ 15. Suivi des KPIs de conformité
□ 16. Inclusion dans le rapport trimestriel
```

---

## RUNBOOK 4 — PRIVACY COMPLAINT (Plainte Vie Privée)

### 4.1 Déclencheur
- Plainte d'un citoyen auprès de l'ANPD ou directement auprès d'une agence
- Signalement d'une violation de vie privée
- Demande d'exercice de droits non satisfaite

### 4.2 Procédure

```
RÉCEPTION (0-5 jours)
━━━━━━━━━━━━━━━━━━━━
□ 1. Réception de la plainte (portail / courrier / en personne)
□ 2. Enregistrement dans le système de gestion des plaintes
□ 3. Accusé de réception au citoyen (5 jours max)
□ 4. Assignation à un analyste
□ 5. Notification au DPD de l'agence concernée

INVESTIGATION (5-30 jours)
━━━━━━━━━━━━━━━━━━━━━━━━━
□ 6. Analyse de la plainte :
    → Recevabilité (le citoyen est-il la personne concernée ?)
    → Fondement (y a-t-il un manquement apparent ?)
    → Urgence (y a-t-il un risque en cours ?)
□ 7. Si irrecevable → Notification motivée + voies de recours
□ 8. Si recevable :
    → Contact avec l'agence mise en cause
    → Demande d'explications (15 jours)
    → Examen des journaux d'accès pertinents
    → Entretien avec les parties si nécessaire

MÉDIATION (si applicable, 30-60 jours)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
□ 9. Tentative de médiation entre le citoyen et l'agence
□ 10. Si accord → Mise en œuvre et vérification → Clôture

DÉCISION (si pas de médiation, 60-90 jours)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
□ 11. Rapport d'investigation
□ 12. Décision :
    → Plainte non fondée → Notification motivée au citoyen
    → Plainte fondée → Mesures correctives
        → Mise en demeure de l'agence
        → Demande de rectification / effacement
        → Sanction si nécessaire
□ 13. Notification de la décision au citoyen
□ 14. Voies de recours indiquées

SUIVI
━━━━━
□ 15. Vérification de l'exécution des mesures
□ 16. Retour au citoyen
□ 17. Clôture du dossier
□ 18. Statistiques et leçons apprises
```

---

## RUNBOOK 5 — REGULATORY AUDIT (Audit Réglementaire)

### 5.1 Déclencheur
- Audit planifié (ANPD, ANC, Cour des Comptes, externe)
- Audit inopiné de l'ANPD
- Demande d'audit du CNGN

### 5.2 Procédure

```
PRÉPARATION (J-60 à J-0 pour audits planifiés)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
□ 1. Réception de la notification d'audit
□ 2. Identification du périmètre d'audit
□ 3. Désignation du coordinateur d'audit (BNC)
□ 4. Constitution de l'équipe interne
□ 5. Préparation de la documentation :
    → Registre des traitements
    → Politiques et procédures
    → Journaux d'audit
    → Rapports précédents
    → Preuves de conformité
□ 6. Pré-audit interne (gap analysis)
□ 7. Briefing de l'équipe
□ 8. Logistique (salle, accès, équipements)

EXÉCUTION (Pendant l'audit)
━━━━━━━━━━━━━━━━━━━━━━━━━━
□ 9. Réunion d'ouverture
    → Présentation du périmètre et de la méthodologie
    → Questions / clarifications
□ 10. Mise à disposition des documents demandés
    → Délai : 24h maximum pour les demandes courantes
    → Documents classifiés : avec approbation du RSSI
□ 11. Facilitation des interviews
□ 12. Accès supervisé aux systèmes
□ 13. Suivi quotidien avec le coordinateur
□ 14. Documentation des échanges

RÉCEPTION DU RAPPORT (Post-audit)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
□ 15. Réception du rapport provisoire
□ 16. Analyse des constatations
□ 17. Rédaction des observations (30 jours)
□ 18. Réception du rapport définitif

PLAN D'ACTION (J+30 à J+90)
━━━━━━━━━━━━━━━━━━━━━━━━━━━
□ 19. Élaboration du plan d'action pour chaque constatation :
    → Action corrective
    → Responsable
    → Délai
    → Indicateur de résolution
□ 20. Validation du plan par le BNC
□ 21. Soumission à l'auditeur
□ 22. Approbation

SUIVI (J+90 à J+365)
━━━━━━━━━━━━━━━━━━━━
□ 23. Exécution des actions correctives
□ 24. Reporting mensuel de progression
□ 25. Vérification par le BNC
□ 26. Rapport de suivi à l'auditeur (6 mois)
□ 27. Clôture des constatations
□ 28. Archivage du dossier d'audit complet
```

---

## 6. MODÈLE DE REGISTRE DES ÉVÉNEMENTS RÉGLEMENTAIRES

| Champ | Description |
|-------|-------------|
| ID | Identifiant unique (REGEV-YYYY-XXXX) |
| Type | Data Breach / Injunction / Compliance Failure / Privacy Complaint / Audit |
| Date de détection | Date et heure |
| Sévérité | Critique / Élevée / Moyenne / Faible |
| Agence concernée | Agence SNISID impliquée |
| Description | Résumé de l'événement |
| Statut | Ouvert / En cours / Résolu / Clôturé |
| Responsable | Personne en charge |
| Actions prises | Liste des actions |
| Notifications | Entités notifiées et dates |
| Documents | Pièces jointes |
| Date de clôture | Date de résolution |
| Leçons apprises | Améliorations identifiées |

---

*Runbooks opérationnels préparés dans le cadre de la Phase 14 — SNISID National Legal Framework*
