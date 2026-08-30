# SNISID-BIO-ADN — Spécifications Index Biens
**Document ID :** SNISID-IDX-BIE-001 | **Catégorie :** NCIC-HT (Biens)

---

## INDEX BIE-VEH — Véhicules Volés

**Équivalent NCIC :** Stolen Vehicle File  
**Intégration :** FOVeS/SIV (MP-15) — synchronisation bidirectionnelle  
**Accès :** PNH terrain (lecture), dcpj.investigator, LAPI (lecture automatique)

### Champs d'identification
| Champ | Format | Obligatoire | Source |
|-------|--------|-------------|--------|
| VIN | 17 caractères alphanumériques | Oui si connu | FOVeS |
| Plaque | Dép-Type-NNNN | Oui | FOVeS |
| Marque/Modèle/Année | VARCHAR | Oui | FOVeS |
| Couleur | VARCHAR | Oui | Déclaration |
| NIU propriétaire | VARCHAR(20) | Oui si connu | SNISID Core |

### Intégration LAPI (MP-16)
```
Agent LAPI scanne plaque ───► Requête temps réel BIE-VEH (< 200ms)
                                        │
                         ┌──────────────┴──────────────┐
                         ▼                             ▼
                    HIT TROUVÉ                   Pas de hit
                         │                             │
              Alerte PNH + photo du             Rien (passage
              véhicule + localisation GPS       normal)
```

### Procédure de récupération
```
1. PNH immobilise le véhicule
2. Interrogation BIE-VEH pour confirmation
3. Contact obligatoire agence entrante (même règle que PER-REC)
4. Mise à jour statut → RECOVERED
5. Notification automatique propriétaire (DIDComm → wallet SNISID)
6. Synchronisation FOVeS/SIV
```

---

## INDEX BIE-ARM — Armes à Feu Volées

**Équivalent NCIC :** Gun File  
**Alimenté par :** PNH (déclarations), DCPJ, Douanes haïtiennes  
**Accès :** dcpj.director, dcpj.investigator (restreint)

### Numéro de série (champ critique)
Le numéro de série est le champ primaire d'identification.
En cas de numéro effacé ou altéré, le dossier est marqué
`serial_obliterated = TRUE` et transmis au labo balistique.

### Types d'armes cataloguées
```
PISTOL          — Pistolet semi-automatique
REVOLVER        — Revolver
RIFLE           — Fusil semi-auto ou à pompe
SHOTGUN         — Fusil de chasse
MACHINEGUN      — Arme automatique (trafic grave)
EXPLOSIVE       — Grenades, explosifs (DCPJ Anti-Gang)
OTHER           — Autre arme à feu
```

### Alerte sur scènes de crime
Si une arme inscrite BIE-ARM est retrouvée sur une scène de crime,
un hit croise automatiquement avec les dossiers pénaux DCPJ
(via Kafka topic `snisid.bio.arm.hit`).

---

## INDEX BIE-DOC — Documents Volés

**Équivalent NCIC :** Article File (documents d'identité)  
**Alimenté par :** ONI (passeports), ONACA (titres fonciers), PNH (déclarations)  
**Accès :** Agents frontières, PNH terrain, ONI

### Types de documents prioritaires pour Haïti

| Type | Émetteur | Critique |
|------|---------|---------|
| Passeport haïtien | ONI | OUI — trafic humain |
| CIN (Carte d'Identité Nationale) | ONI | OUI — usurpation identité |
| Acte de naissance | ONEC | OUI — fraude état civil |
| Permis de conduire | MTPTC | MOYEN |
| Titre foncier | ONACA | OUI — fraude immobilière |

### Intégration frontières
Les postes frontières (aéroport PAP, Malpasse, Ouanaminthe) interrogent
BIE-DOC en temps réel à chaque présentation d'un document d'identité.

```
Document présenté au poste frontière
            │
            ▼
Interrogation BIE-DOC (numéro + type + date émission)
            │
      ┌─────┴──────┐
      ▼            ▼
  HIT TROUVÉ   Pas de hit
      │
  Alerte agent + protocole interpellation
```

### Synchronisation ONI ↔ BIE-DOC
Lorsque l'ONI annule un document (signalement perte/vol ou révocation),
l'événement Kafka `snisid.oni.document.revoked` déclenche
automatiquement la création d'un enregistrement BIE-DOC.

---

## INDEX BIE-PLQ — Plaques Minéralogiques

**Équivalent NCIC :** License Plate File  
**Intégration :** LAPI (MP-16) — usage LAPI exclusivement  
**Accès :** LAPI automatique, PNH terrain (lecture)

### Format des plaques haïtiennes
```
Format standard : [Dép]-[Type]-[4 chiffres]
Exemples :
  OE-A-1234   (Ouest, particulier)
  NO-G-5678   (Nord, gouvernement)
  NA-D-0001   (Nord, diplomatique)
  SU-T-2345   (Sud, transport public)
```

### Cas de plaques suspectes
- Plaque clonée (même numéro sur deux véhicules) : alerte haute
- Plaque d'un véhicule inscrit BIE-VEH : alerte critique
- Plaque périmée avec véhicule encore en circulation : notification PNH

---

## INDEX BIE-EMB — Embarcations (Navires et Bateaux)

**Équivalent NCIC :** Boat File  
**Spécificité Haïti :** Index stratégique — 1771 km de côtes, trafic humain, drogue  
**Accès :** CGFADH (Garde-Côtes), DCPJ Maritime, BDNAH

### Importance pour Haïti

Haïti dispose d'une côte de **1771 km** et fait face à :
- Trafic de migrants (boat people) — embarcations de fortune volées
- Trafic de drogue (Amérique du Sud → USA via Haïti)
- Piraterie maritime dans le Golfe de la Gonâve et Canal du Vent
- Contrebande d'armes

BIE-EMB est donc **plus critique pour Haïti** que pour la plupart des pays.

### Types d'embarcations cataloguées
```
FISHING_CANOE    — Canot de pêche traditionnel (yole)
MOTORBOAT        — Bateau à moteur < 12m
SAILBOAT         — Voilier de plaisance
FERRY            — Traversier inter-îles (ex: PAP–Jeremie)
CARGO_SMALL      — Caboteur < 500 tonnes
PATROL_BOAT      — Bateau de patrouille (CGFADH)
OTHER            — Autre embarcation
```

### Champs d'identification spéciaux
| Champ | Description |
|-------|-------------|
| hull_id_number (HIN) | Identifiant coque international (si gravé) |
| vessel_name | Nom peint sur la coque |
| home_port | Port d'attache (PAP, Cap-Haïtien, Jacmel, etc.) |
| engine_serial | Numéro série moteur hors-bord |
| distinctive_marks | Peinture, marquages, réparations visibles |

### Coordination Garde-Côtes (CGFADH)
```
CGFADH intercepte embarcation
        │
        ▼
Interrogation BIE-EMB (nom + HIN + plaque si applicable)
        │
        ├── HIT → Procédure d'arraisonnement + transmission DCPJ
        └── Pas de hit → Contrôle standard (documents, cargaison)
```

---

## INDEX BIE-OBJ — Objets et Biens Précieux Volés

**Équivalent NCIC :** Article File (valeurs/bijoux)  
**Alimenté par :** PNH (déclarations citoyens), compagnies d'assurance  
**Accès :** PNH terrain, dcpj.investigator, douanes

### Catégories d'objets
```
JEWELRY          — Bijoux (or, argent, diamants)
ART              — Œuvres d'art, sculptures, peintures
ELECTRONICS      — Téléphones, ordinateurs, équipements
CURRENCY         — Monnaie (billets volés, numéros de série si connus)
CATTLE           — Bétail (critique en zone rurale haïtienne)
MACHINERY        — Équipements agricoles et industriels
OTHER            — Autre bien déclaré
```

> **Note :** Le bétail (bovins, équins) est inclus car les vols de bétail
> constituent une problématique majeure dans les zones rurales haïtiennes.

---

## INDEX BIE-TIT — Titres et Valeurs

**Équivalent NCIC :** Securities File  
**Alimenté par :** BRH (Banque de la République d'Haïti), MJSP  
**Accès :** BRH, institutions financières agréées, dcpj.investigator

### Types de titres
- Chèques en blanc volés (numéros de série)
- Obligations d'État
- Titres de propriété (doublon avec BIE-DOC pour les titres fonciers)
- Lettres de crédit frauduleuses
