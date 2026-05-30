# SNISID — MODÈLE NATIONAL DE GESTION DES ARCHIVES NUMÉRIQUES

**Classification :** CADRE OPÉRATIONNEL — ARCHIVES NATIONALES
**Référence :** SNISID-ARCH-001
**Version :** 1.0
**Date :** 25 mai 2026

---

## 1. OBJECTIF

Établir le cadre de gestion des archives numériques nationales du SNISID, garantissant la pérennité, l'intégrité et l'accessibilité à long terme des enregistrements nationaux.

---

## 2. PÉRIMÈTRE

| Type d'Archive | Description | Importance |
|---------------|-------------|-----------|
| Registre National d'Identité | Données d'identité de tous les citoyens | Vitale |
| Journaux d'audit | Historique de toutes les opérations | Critique |
| Documents administratifs | Actes, décisions, correspondances | Élevée |
| Décisions de justice numériques | Jugements, ordonnances | Vitale |
| Archives réglementaires | Lois, règlements, politiques | Élevée |
| Archives de sécurité | Rapports d'incidents, enquêtes | Critique |
| Archives techniques | Documentation technique, configurations | Moyenne |

---

## 3. FONCTIONS SUPPORTÉES

### 3.1 Archival Policies (Politiques d'Archivage)

**Politique générale :**

| Principe | Application |
|----------|------------|
| Complétude | Tout document officiel du SNISID doit être archivé |
| Authenticité | L'authenticité de chaque archive est vérifiable |
| Intégrité | L'intégrité est garantie par hash cryptographique |
| Accessibilité | Les archives sont accessibles aux personnes habilitées |
| Pérennité | Les archives sont conservées dans des formats pérennes |
| Classification | Chaque archive est classifiée selon le schéma national |

**Formats d'archivage pérennes :**

| Type de document | Format d'archivage | Standard |
|-----------------|-------------------|---------|
| Documents texte | PDF/A-3 | ISO 19005-3 |
| Données structurées | XML signé / JSON-LD signé | W3C |
| Images | TIFF non compressé / JPEG 2000 | ISO 12234 |
| Vidéo | MPEG-4 AVC | ISO 14496 |
| Audio | WAV / FLAC | Standard ouvert |
| Bases de données | Export SQL normalisé + schéma | Custom |
| Journaux | Texte structuré signé | Syslog RFC 5424 |

**Métadonnées d'archivage obligatoires :**

| Métadonnée | Description |
|-----------|-------------|
| Identifiant unique | UUID national |
| Titre | Description du document |
| Créateur | Personne ou système ayant créé le document |
| Date de création | Horodatage qualifié |
| Classification | Niveau de confidentialité |
| Type documentaire | Catégorie (acte, rapport, journal, etc.) |
| Format | Format technique du fichier |
| Hash d'intégrité | SHA-384 du contenu |
| Signature | Signature numérique du créateur |
| Durée de rétention | Période de conservation |
| Statut | Actif / archivé / en gel légal / à détruire |
| Mots-clés | Descripteurs thématiques |

### 3.2 Retention Schedules (Calendriers de Rétention)

| Catégorie | Durée Active | Durée Archive | Durée Totale | Sort Final |
|-----------|-------------|--------------|-------------|-----------|
| Registre d'identité | Vie du citoyen + 30 ans | Permanent | Permanent | Conservation permanente |
| Données biométriques | Vie du citoyen + 10 ans | 20 ans supplémentaires | Vie + 30 ans | Destruction sécurisée |
| Journaux d'audit | 2 ans | 8 ans | 10 ans | Destruction sécurisée |
| Documents administratifs courants | 5 ans | 10 ans | 15 ans | Tri puis destruction |
| Actes d'état civil | Permanent | - | Permanent | Conservation permanente |
| Décisions de justice | 30 ans | Permanent | Permanent | Conservation permanente |
| Archives réglementaires | Tant que en vigueur + 10 ans | Permanent | Permanent | Conservation permanente |
| Données d'enquête (non jugées) | 5 ans | 5 ans | 10 ans | Destruction sécurisée |
| Données d'enquête (jugées) | 30 ans | Permanent | Permanent | Conservation permanente |
| Correspondances officielles | 5 ans | 10 ans | 15 ans | Tri puis destruction |
| Documentation technique | Vie du système + 5 ans | 10 ans | Variable | Tri puis destruction |
| Rapports d'audit | 5 ans | 10 ans | 15 ans | Conservation permanente |
| Données de formation | 3 ans | 2 ans | 5 ans | Destruction |

**Processus de révision du calendrier :**
- Revue annuelle par le Comité des Archives Numériques
- Consultation de l'ANPD (protection données) et des Archives Nationales
- Approbation par le CNGN

### 3.3 Immutable Records (Enregistrements Immutables)

| Mécanisme | Description |
|-----------|-------------|
| Stockage WORM | Write Once Read Many pour les archives permanentes |
| Signature numérique | Chaque archive signée par le système d'archivage |
| Chaîne de hash | Lien cryptographique entre archives successives |
| Horodatage qualifié | TSA nationale sur chaque archive |
| Réplication | 3 copies minimum sur sites géographiquement distincts |
| Vérification d'intégrité | Vérification automatique mensuelle de tous les hash |
| Contrôle d'accès strict | Aucun droit de modification/suppression sans procédure |

**Architecture d'archivage immutable :**
```
Document source
    → Conversion en format pérenne (PDF/A, XML signé)
    → Calcul hash SHA-384
    → Ajout métadonnées normalisées
    → Signature numérique du système d'archivage
    → Horodatage qualifié TSA
    → Stockage WORM (site primaire)
    → Réplication automatique (site secondaire + tertiaire)
    → Inscription au registre des archives
    → Vérification d'intégrité immédiate (hash check)
```

### 3.4 Legal Holds (Gels Légaux)

| Aspect | Description |
|--------|-------------|
| Déclenchement | Sur demande de la Justice (ordonnance, mandat, réquisition) |
| Effet | Suspension de toute destruction, modification ou transfert |
| Périmètre | Ciblé (documents spécifiques) ou large (catégorie) |
| Durée | Jusqu'à levée par l'autorité judiciaire |
| Notification | L'administration est notifiée du gel (pas de son contenu) |
| Non-contournement | Techniquement impossible de supprimer un document gelé |

**Processus de gel légal :**
```
Réception de l'ordonnance judiciaire
    → Vérification de l'authenticité (signature du magistrat)
    → Identification des documents concernés
    → Application du flag LEGAL_HOLD
    → Copie de sécurité supplémentaire
    → Notification au responsable des archives
    → Journalisation de l'action
    → Confirmation au tribunal
    → Monitoring : aucune action sur les documents gelés
    → Levée du gel uniquement sur nouvelle ordonnance
```

### 3.5 Historical Preservation (Préservation Historique)

| Fonction | Description |
|----------|-------------|
| Migration de format | Conversion proactive vers les nouveaux formats pérennes |
| Migration de support | Transfert vers les nouveaux supports de stockage |
| Vérification de lisibilité | Test de lecture de chaque archive tous les 5 ans |
| Émulation | Capacité d'émuler les anciens systèmes si nécessaire |
| Documentation contextuelle | Contexte historique accompagnant chaque archive |
| Accès public | Archives historiques accessibles au public (délai de communicabilité respecté) |

**Calendrier de migration technologique :**

| Action | Fréquence | Responsable |
|--------|-----------|-------------|
| Audit des formats | Annuel | Archives numériques |
| Migration proactive | Tous les 5 ans | Archives numériques |
| Test de lisibilité | Tous les 5 ans | Archives numériques |
| Renouvellement du stockage | Tous les 7 ans | Infrastructure |
| Re-signature cryptographique | Avant expiration des algorithmes | PKI nationale |

---

## 4. STRUCTURE DE GOUVERNANCE DES ARCHIVES

```
Comité National des Archives Numériques
├── Président : Directeur des Archives Nationales
├── Membre : Responsable SNISID
├── Membre : Représentant ANPD
├── Membre : Représentant Justice
├── Membre : Archiviste en chef
└── Membre : Expert en préservation numérique

Missions :
- Politique d'archivage
- Calendriers de rétention
- Standards de format
- Budget de préservation
- Accès aux archives historiques
```

---

## 5. INDICATEURS

| KPI | Cible | Mesure |
|-----|-------|--------|
| Taux d'intégrité des archives | 100% | Mensuel (vérification hash) |
| Taux de conformité au calendrier de rétention | 100% | Trimestriel |
| Taux de migration proactive | 100% dans les délais | Annuel |
| Délai de réponse aux demandes d'accès | ≤ 5 jours | Continu |
| Disponibilité du système d'archivage | 99.99% | Continu |
| Nombre de gels légaux actifs | Tracking | Continu |

---

*Document cadre préparé dans le cadre de la Phase 14 — SNISID National Legal Framework*
