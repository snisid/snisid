# SNISID-BIO-ADN — Spécifications Index ADN
**Document ID :** SNISID-IDX-ADN-001 | **Catégorie :** CODIS-HT

---

## INDEX BIO-CON — ADN Condamnés

**Équivalent CODIS :** Convicted Offender Index  
**Niveau minimum :** SDIS-HT  
**Accès :** lab.supervisor, ndis.analyst, dcpj.director

### Critères d'éligibilité (Haïti)
- Condamnation définitive pour : meurtre, viol, enlèvement, trafic de personnes,
  vol à main armée, terrorisme, crimes organisés.
- Le prélèvement est ordonné par le juge compétent dans le jugement.
- Consentement **non requis** après condamnation définitive (comme DNA Act USA).

### Workflow de collecte
```
1. Jugement définitif rendu
       │
       ▼
2. Ordre de prélèvement transmis à l'établissement pénitentiaire
       │
       ▼
3. Prélèvement buccal par professionnel médical accrédité
       │
       ▼
4. Envoi scellé au laboratoire LDIS-HT accrédité
       │
       ▼
5. Génération profil STR (20 loci CODIS Core)
       │
       ▼
6. Validation qualité (quality_score ≥ 0.90)
       │
       ▼
7. Upload SDIS-HT → NDIS-HT (schedule hebdomadaire)
```

### Champs obligatoires
| Champ | Type | Description |
|-------|------|-------------|
| specimen_number | VARCHAR(100) | Numéro de scellé unique |
| loci_encrypted | BYTEA | 20 loci STR chiffrés HSM |
| amelogenin | CHAR(2) | XX ou XY |
| quality_score | DECIMAL | ≥ 0.90 obligatoire |
| lab_id | UUID | Laboratoire accrédité |
| case_number | VARCHAR | Numéro dossier judiciaire |

---

## INDEX BIO-ARR — ADN Arrestés

**Équivalent CODIS :** Arrestee Index  
**Niveau :** LDIS-HT uniquement (pas de partage NDIS sans condamnation)  
**Accès :** lab.technician (propre lab uniquement)

### Critères d'éligibilité
- Arrestation pour crime (meurtre, viol, enlèvement, trafic d'armes).
- Ordonnance du juge d'instruction **requise**.
- **Expungement automatique** si non-lieu ou acquittement dans les 90 jours.

### Règles d'expungement BIO-ARR
```
Arrestation
    │
    ├── Condamnation → Profil migré vers BIO-CON
    ├── Acquittement → Expungement automatique J+30
    └── Classement sans suite → Expungement J+90
```

---

## INDEX BIO-FSC — ADN Scènes de Crime

**Équivalent CODIS :** Forensic/Case Index  
**Niveau :** LDIS → SDIS → NDIS (tous niveaux)  
**Accès :** ndis.analyst, dcpj.investigator

### Types de traces forensiques acceptées
- Sang, salive, sperme, cheveux avec follicule, peau
- Prélèvement réalisé par technicien de scène de crime certifié PNH/DCPJ

### Critères de qualité minimaux
| Critère | Seuil |
|---------|-------|
| Nombre de loci valides | ≥ 10 sur 20 |
| Quality score | ≥ 0.60 |
| Absence de mélange (mixture) | Ratio ≥ 60% contributeur principal |

### Algorithme de matching BIO-FSC → BIO-CON
```
BIO-FSC (scène de crime, loci_count ≥ 10)
        │
        ▼
Recherche LDIS → SDIS → NDIS
        │
        ├── FULL_MATCH (confidence ≥ 0.999) → Alerte CRITIQUE + rapport DCPJ
        ├── PARTIAL   (confidence 0.85–0.999) → Alerte HAUTE + investigation
        └── FAMILIAL  (confidence 0.40–0.84)  → Signalement discret
```

---

## INDEX BIO-DIS — ADN Personnes Disparues

**Équivalent CODIS :** Missing Persons + Relatives of Missing Persons  
**Niveau :** LDIS → NDIS  
**Accès :** ndis.analyst, dcpj.investigator

### Deux sous-types
1. **Direct** : ADN de la personne disparue (si disponible — ex. brosse à dents)
2. **Familial** : ADN d'un parent biologique (pour recherche par parenté)

### Matching familial
Le matching familial utilise un calcul de **Likelihood Ratio (LR)** tenant compte :
- Probabilité allélique par population (fréquences alléliques haïtiennes)
- Relation attendue (parent, enfant, fratrie)
- Seuil LR ≥ 1000 pour signalement, LR ≥ 10 000 pour alerte formelle

> **Note :** Les fréquences alléliques haïtiennes doivent être établies via
> une étude de population MSPP/ANH avant la mise en production de cet index.

---

## INDEX BIO-RNI — Restes Humains Non Identifiés

**Équivalent CODIS :** Unidentified Human Remains  
**Niveau :** LDIS → NDIS  
**Accès :** ndis.analyst, médecins légistes accrédités

### Sources
- Morgues MSPP (HUP, HUEH, Cap-Haïtien, Gonaïves)
- Scènes de catastrophes naturelles (séismes, ouragans)
- Scènes de crimes avec victimes non identifiées

### Flux de croisement automatique
```
BIO-RNI (restes non identifiés)
    │
    ├── Croisé avec BIO-DIS (personnes disparues)
    ├── Croisé avec BIO-CON (historique condamnés)
    └── Croisé avec INTERPOL DNA Gateway (victimes internationales)
```

---

## STANDARD TECHNIQUE : 20 LOCI CODIS CORE (NIST 2017)

```json
{
  "codis_core_loci": [
    "CSF1PO", "D3S1358", "D5S818",  "D7S820",  "D8S1179",
    "D13S317","D16S539","D18S51",   "D21S11",  "FGA",
    "TH01",   "TPOX",   "vWA",      "D1S1656", "D2S441",
    "D2S1338","D10S1248","D12S391", "D19S433", "D22S1045"
  ],
  "sex_marker": "Amelogenin",
  "total_loci": 21,
  "minimum_loci_for_upload": {
    "BIO-CON": 20,
    "BIO-ARR": 18,
    "BIO-FSC": 10,
    "BIO-DIS": 15,
    "BIO-RNI": 8
  }
}
```
