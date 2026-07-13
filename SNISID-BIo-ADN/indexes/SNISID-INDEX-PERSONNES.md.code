# SNISID-BIO-ADN — Spécifications Index Personnes
**Document ID :** SNISID-IDX-PER-001 | **Catégorie :** NCIC-HT (Personnes)

---

## INDEX PER-REC — Personnes Recherchées (Wanted Persons)

**Équivalent NCIC :** Wanted Person File  
**Alimenté par :** PNH, DCPJ, MJSP, Parquet  
**Accès :** dcpj.investigator, bio.ndis.analyst, PNH terrain (lecture)

### Types de mandats
| Code | Type | Émetteur |
|------|------|---------|
| `MAN-ARR` | Mandat d'arrêt | Juge d'instruction / Tribunal |
| `MAN-EXT` | Mandat d'extradition | MJSP + Ministère AE |
| `MAN-REC` | Mandat de recherche | DCPJ |
| `AVIS-REC` | Avis de recherche (non judiciaire) | PNH / familles |

### Champs obligatoires à la saisie
```
record_number    : auto-généré (PRE-AAAA-NNNNNN)
last_name        : Obligatoire sauf si inconnu (mettre "INCONNU")
warrant_type     : Obligatoire
warrant_number   : Obligatoire pour MAN-ARR et MAN-EXT
issuing_date     : Obligatoire
charges          : Au moins 1 chef d'accusation
entering_agency  : Code agence PNH/DCPJ
entering_officer : NIU agent saisissant
mco_contact      : Numéro de contact agence (pour vérification hit)
```

### Cycle de vie d'un enregistrement PER-REC
```
ACTIVE ──────► CLEARED  (arrestation effectuée)
  │
  ├──────────► EXPIRED   (mandat expiré, non renouvelé)
  │
  └──────────► SUSPENDED (suspendu par autorité judiciaire)
```

### Règle de vérification des hits
Inspiré NCIC : avant toute arrestation basée sur un hit PER-REC,
l'agent doit contacter l'agence entrant (champ `mco_contact`) pour confirmer
que le mandat est toujours actif. Le hit Kafka inclut ce contact.

---

## INDEX PER-FUG — Fugitifs Étrangers

**Équivalent NCIC :** Foreign Fugitive File  
**Alimenté par :** DCPJ INTERPOL NCB (Bureau Central National d'Haïti)  
**Accès :** dcpj.director, ndis.analyst

### Intégration INTERPOL
```
INTERPOL I-24/7 (notice rouge)
        │
        ▼
BCN Haïti (DCPJ)
        │
        ▼
Saisie PER-FUG dans SNISID-BIO-ADN
        │
        ▼
Diffusion automatique vers PER-REC (copie avec flag interpol_notice)
```

> Les agents de terrain haïtiens ne peuvent interroger INTERPOL directement.
> Toute requête INTERPOL passe par le BCN DCPJ.

---

## INDEX PER-DIS — Personnes Disparues

**Équivalent NCIC :** Missing Person File (5 catégories)  
**Alimenté par :** PNH, DCPJ, familles via portail citoyen SNISID  
**Accès :** dcpj.investigator, ndis.analyst, portail citoyen (saisie seulement)

### Catégories et priorités

| Catégorie | Code | Priorité | Description |
|-----------|------|----------|-------------|
| Enfant | `CHILD` | P1 CRITIQUE | Moins de 18 ans, toute circonstance |
| En danger | `ENDANGERED` | P1 CRITIQUE | Tout âge, risque physique confirmé |
| Involontaire | `INVOLUNTARY` | P2 HAUTE | Disparition contre la volonté |
| Catastrophe | `CATASTROPHE` | P2 HAUTE | Suite séisme, ouragan, inondation |
| Autre | `OTHER` | P3 NORMALE | Adulte, aucun facteur aggravant |

### Procédure spéciale enfants (CHILD)
1. Alerte immédiate tous les postes PNH du département
2. Diffusion portrait sur portail SNISID public (avec accord famille)
3. Croisement automatique BIO-DIS si ADN familial disponible
4. Notification automatique Brigades de Protection des Mineurs (BPM)
5. Si non-localisé J+30 : transmission INTERPOL Missing Persons

### Portail citoyen (saisie PER-DIS)
Les citoyens peuvent signaler une disparition via le portail SNISID.
La saisie est validée par un agent PNH dans les 24h avant activation.

---

## INDEX PER-NID — Personnes Non Identifiées

**Équivalent NCIC :** Unidentified Person File  
**Alimenté par :** Médecins légistes, morgues accréditées, PNH  
**Accès :** ndis.analyst, médecins légistes accrédités

### Sources principales en Haïti
- HUP (Hôpital de l'Université d'État d'Haïti)
- HUEH — Hôpital Universitaire Estrabon Haïti
- Morgues départementales (10 départements)
- Scènes de catastrophes naturelles

### Croisement automatique PER-NID ↔ PER-DIS
À chaque nouvelle entrée PER-NID, un job Kafka déclenche une recherche
automatique dans PER-DIS sur :
1. Correspondance physique (âge, sexe, signalement)
2. Correspondance ADN si échantillon BIO-RNI disponible
3. Correspondance empreintes digitales si disponibles (lien NGI-HT)

---

## INDEX PER-SEX — Registre Délinquants Sexuels

**Équivalent USA :** NSSOR (National Sex Offender Registry)  
**Alimenté par :** Tribunaux (après condamnation définitive)  
**Accès :** dcpj.director, PNH/DCPJ vérifications d'antécédents

### Niveaux de risque
| Niveau | Conditions | Obligations |
|--------|-----------|-------------|
| LOW | Délit unique, victime adulte consentante | Enregistrement annuel |
| MEDIUM | Récidive ou victime mineure | Enregistrement semestriel + notification voisinage |
| HIGH | Récidive avec violence ou multiple victimes mineures | Enregistrement trimestriel + restriction géographique |

### Restrictions HIGH
- Interdiction de résider à moins de 500m d'une école, église, terrain de jeu
- Obligation de déclarer tout changement d'adresse dans les 72h
- Vérification physique par agent PNH tous les 90 jours

---

## INDEX PER-GNG — Membres de Gangs

**Équivalent NCIC :** Gang File  
**Alimenté par :** DCPJ Unité Anti-Gang, PNH Renseignement  
**Accès :** dcpj.director, dcpj.investigator (restreint, audit systématique)

### Critères d'inscription (inspiré CalGang — USA)
Un individu peut être inscrit PER-GNG si **au moins 3 critères** parmi :
1. Déclaration d'appartenance à un gang
2. Appréhendé en compagnie de membres de gang connus
3. Tatouages ou insignes gang identifiés
4. Documents ou photos le montrant avec membres de gang
5. Renseignement d'un informateur fiable (niveau A/B)
6. Condamnation pour crime lié à une activité de gang

### Règles de révision obligatoire
- Révision annuelle par superviseur DCPJ
- Radiation automatique après 5 ans sans nouvelle activité
- La personne inscrite peut contester son inscription devant le MJSP

---

## INDEX PER-TER — Terrorisme et Financement

**Équivalent NCIC :** Known or Appropriately Suspected Terrorist (KST)  
**Alimenté par :** DCPJ Unité CT, MJSP, transmission BCN INTERPOL  
**Accès :** dcpj.director uniquement (plus restreint de tous les index)

### Règles strictes
- Inscription requiert approbation écrite du Directeur DCPJ + AG (Procureur Général)
- Révision obligatoire tous les 6 mois par comité tri-partite (DCPJ, MJSP, PGR)
- Partage INTERPOL via BCN uniquement, pas de diffusion directe

---

## INDEX PER-OPR — Ordres de Protection

**Équivalent NCIC :** Protection Order File  
**Alimenté par :** Tribunaux civils et pénaux, MCFDF (femmes/enfants)  
**Accès :** PNH terrain (lecture urgente), dcpj.investigator

### Priorité d'accès terrain
Les agents PNH disposent d'un accès **lecture urgente** (< 5 secondes)
à cet index lors d'interventions domestiques. Le système retourne :
- Existence d'une ordonnance
- Type d'ordonnance et restrictions
- Agence et juge émetteur
- Contact d'urgence bénéficiaire

---

## INDEX PER-LIB — Libération Conditionnelle

**Équivalent NCIC :** Supervised Release File  
**Alimenté par :** DGAP (Direction Générale Administration Pénitentiaire)  
**Accès :** dcpj.investigator, agents surveillance judiciaire

### Conditions de supervision trackées
```json
{
  "niu": "HTI-...",
  "supervision_type": "CONDITIONAL_RELEASE",
  "start_date": "2026-01-15",
  "end_date": "2027-01-15",
  "conditions": [
    "Interdiction de quitter le département Ouest",
    "Pointage hebdomadaire PNH Delmas",
    "Interdiction de contact avec la victime"
  ],
  "supervising_officer": "Agent NIU",
  "agency": "DGAP"
}
```
