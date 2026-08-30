# SNISID-BIO-ADN — Normes Qualité ADN Forensique
**Document ID :** SNISID-BIO-QUA-001 | **Version :** 1.0.0

---

## 1. NORMES DE RÉFÉRENCE

| Norme | Organisme | Application |
|-------|-----------|-------------|
| ISO/IEC 17025:2017 | ISO | Accréditation labs d'essais |
| SWGDAM Guidelines 2020 | FBI/USA | Validation méthodes STR |
| ISO 18385:2016 | ISO | Minimisation contamination ADN |
| ENFSI Guidelines | Europe | Interprétation profils mixtes |
| CODIS QAS (FBI) | FBI | Standards contrôle qualité NDIS |

---

## 2. SEUILS DE QUALITÉ PAR INDEX

| Index | Quality Score Min | Loci Min | Commentaire |
|-------|-------------------|----------|-------------|
| BIO-CON | 0.95 | 20/20 | Profil complet exigé |
| BIO-ARR | 0.90 | 18/20 | Quasi-complet |
| BIO-FSC | 0.60 | 10/20 | Tolérance pour traces dégradées |
| BIO-DIS | 0.85 | 15/20 | Matching familial nécessite haute qualité |
| BIO-RNI | 0.50 | 8/20 | Restes dégradés — seuil bas accepté |

---

## 3. CONTRÔLES QUALITÉ OBLIGATOIRES

### 3.1 Par analyse (interne labo)
- [ ] Témoin positif inclus (ADN de référence certifié NIST)
- [ ] Témoin négatif inclus (eau stérile)
- [ ] Blanc réactif inclus
- [ ] Ladder allélique validé
- [ ] Électrophorégramme interprété par 2 techniciens indépendants

### 3.2 Externe (inter-laboratoires)
- [ ] Participation au programme PT (Proficiency Testing) annuel
- [ ] Échange d'échantillons témoins entre LDIS-PAP-001 et LDIS-CAP-001
- [ ] Audit SNISID semestriel documenté

### 3.3 Calcul du Quality Score

```python
# bio-adn-service/internal/quality/scorer.py

def calculate_quality_score(electropherogram: dict) -> float:
    """
    Calcule le quality score d'un profil STR.
    Score = (loci_valides / 20) × facteur_intensité × facteur_bruit
    """
    codis_loci = [
        "CSF1PO","D3S1358","D5S818","D7S820","D8S1179",
        "D13S317","D16S539","D18S51","D21S11","FGA",
        "TH01","TPOX","vWA","D1S1656","D2S441",
        "D2S1338","D10S1248","D12S391","D19S433","D22S1045"
    ]

    valid_loci = 0
    intensity_scores = []
    noise_penalties = 0

    for locus in codis_loci:
        if locus not in electropherogram:
            continue

        peak = electropherogram[locus]
        if peak["height"] < 150:       # RFU minimum
            noise_penalties += 1
            continue
        if peak["height"] > 30000:     # Pull-up / saturation
            noise_penalties += 0.5
            continue

        valid_loci += 1
        intensity_scores.append(min(peak["height"] / 5000, 1.0))

    base_score = valid_loci / len(codis_loci)
    intensity_factor = sum(intensity_scores) / max(len(intensity_scores), 1)
    noise_factor = max(0, 1 - (noise_penalties * 0.05))

    return round(base_score * 0.6 + intensity_factor * 0.3 + noise_factor * 0.1, 3)
```

---

## 4. FRÉQUENCES ALLÉLIQUES HAÏTIENNES

> ⚠️ **ACTION REQUISE** : Les fréquences alléliques doivent être établies
> sur une population haïtienne représentative avant d'utiliser les
> calculs de Likelihood Ratio (LR) et de matching familial.

### Plan d'étude de population
1. **Partenariat :** MSPP + Université d'État d'Haïti + partenaire international (FBI, ENFSI)
2. **Taille d'échantillon :** Minimum 200 individus non apparentés par groupe
3. **Groupes :** Population générale haïtienne (3 régions : Ouest, Nord, Sud)
4. **Durée :** 18 mois
5. **Base de données temporaire :** Utiliser fréquences alléliques afro-caribéennes
   publiées (base NIST CODIS) en attendant l'étude nationale

---

## 5. ARCHIVAGE DES ÉCHANTILLONS BIOLOGIQUES

| Index | Durée conservation | Température | Support |
|-------|-------------------|-------------|---------|
| BIO-CON | Durée condamnation + 10 ans | -80°C | Congélateur forensique |
| BIO-ARR | 3 ans max | -20°C | Congélateur standard |
| BIO-FSC | Durée prescription + 5 ans | -80°C | Congélateur forensique |
| BIO-DIS | Jusqu'à identification | -80°C | Congélateur forensique |
| BIO-RNI | Jusqu'à identification | -80°C | Congélateur forensique |
