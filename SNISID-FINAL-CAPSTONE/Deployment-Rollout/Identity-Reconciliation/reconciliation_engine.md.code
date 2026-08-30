# SNISID National Identity Reconciliation Engine (NIRE)
## Spécification Technique du Moteur National de Réconciliation d'Identité

---

## 1. Introduction & Objectif du Moteur

Le **National Identity Reconciliation Engine** (NIRE) est le noyau technologique du SNISID responsable de garantir le principe d'**Unicité de l'Identité** sur l'ensemble du territoire national. Aucun citoyen ne doit pouvoir s'enrôler sous de multiples identités (fraude à l'identité multiple), et aucun dossier civil ne doit être dupliqué.

```
                              NIRE CROSS-MATCHING GATE
                              
+------------------------------------+       +------------------------------------+
|        DEMOGRAPHIC VECTOR         |       |         BIOMETRIC VECTOR           |
| (Name, DoB, Birthplace, Parentage) |       |     (Fingerprints, Face, Iris)     |
+------------------------------------+       +------------------------------------+
                   \                                          /
                    \                                        /
                     v                                      v
              +----------------------------------------------------+
              |          HYBRID RECONCILIATION LAYER (NIRE)        |
              |       (Weight-based Score & Thresholding)          |
              +----------------------------------------------------+
                                       |
                   +-------------------+-------------------+
                   |                                       |
                   v (Score >= 0.85)                       v (Score < 0.85)
       [DEDUPLICATED / MERGED]                     [NEW UNIQUE IDENTITY]
       - Update existing profile                   - Generate New IUI
       - Flag cross-system reference               - Create master security token
```

Le NIRE effectue un recoupement à deux facteurs (Démographique et Biométrique) pour calculer un indice de confiance de correspondance (*Identity Match Score - IMS*).

---

## 2. Rapprochement Démographique (Demographic Matching Model)

Le recoupement purement textuel utilise un modèle pondéré comparant plusieurs attributs clés d'identité avec tolérance aux fautes d'orthographe (distances de Levenshtein et Jaro-Winkler adaptées).

### 2.1 Attribution des Poids d'Évaluation (Poids Démographiques)

| Attribut Comparé | Algorithme Appliqué | Poids Relatif | Description / Tolérance |
| :--- | :--- | :--- | :--- |
| **Nom de Famille** | Jaro-Winkler | 25% | Prise en compte de l'inversion d'orthographe et phonétique créole. |
| **Prénoms** | Jaro-Winkler | 20% | Tolérance aux variations de traits d'union et abréviations (ex : Jn-Baptiste / Jean-Baptiste). |
| **Date de Naissance**| Écart temporel exact | 25% | Pénalisation graduelle si écart de jours/mois/années (tolérance d'inversion jour/mois). |
| **Lieu de Naissance**| Similarité textuelle | 15% | Standardisation par dictionnaire de synonymes géographiques d'Haïti. |
| **Prénom de la Mère**| Levenshtein | 15% | Variable d'ancrage d'état civil extrêmement efficace en Haïti. |

*   **Seuil de suspicion démographique :** Tout score démographique combiné $\ge 0.75$ déclenche une suspicion d'homonymie ou de doublon, forçant un contrôle biométrique strict.

---

## 3. Rapprochement Biométrique (Biometric Reconciliation)

La biométrie constitue le garde-fou ultime contre la fraude. Le NIRE s'interface directement avec l'ABIS (Automated Biometric Identification System) souverain d'Haïti.

```
                                  ABIS SCORE ENGINE
                                  
   +----------------------+   +----------------------+   +----------------------+
   | Ten-Print Finger     |   | Facetime Match (ICAO)|   | Iris Scan Match      |
   | (FAP-45 Standard)    |   | (Neural Vector Match)|   | (Dual-Iris Vector)   |
   +----------------------+   +----------------------+   +----------------------+
              |                          |                          |
              +--------------------------+--------------------------+
                                         |
                                         v
                      [Biometric Fusion Score (0.00 - 1.00)]
```

### 3.1 Les Trois Piliers Biométriques
1. **Empreintes Digitales (Ten-Print Fingerprints) :**
   *   *Standard :* Format NIST Record Type-9, minutie ANSI/INCITS 378.
   *   *Algorithme :* Comparaison d'empreintes 1-to-N (1-à-plusieurs). Le score de minutie du doigt ayant la meilleure correspondance détermine l'indice de confiance.
2. **Reconnaissance Faciale (Facial Recognition) :**
   *   *Standard :* Conformité ISO/IEC 19794-5 (Norme OACI pour passeports).
   *   *Algorithme :* Modèle de Deep Learning (Réseau de neurones convolutif - CNN) projetant le visage sur un vecteur de caractéristiques de 512 dimensions. Mesure de distance cosinus entre vecteurs.
3. **Reconnaissance de l'Iris (Iris Scan) :**
   *   *Standard :* ISO/IEC 19794-6.
   *   *Algorithme :* Moteur Daugman d'analyse de texture d'iris. Utilisé en cas de dégradation sévère des empreintes digitales (ex: travailleurs agricoles).

### 3.2 Seuil de Fusion Biométrique
*   **Match Biométrique Confirmé :** Si le score consolidé de fusion biométrique est $\ge 0.88$, l'identité est considérée comme biométriquement identique de manière absolue (Probabilité d'erreur de correspondance $< 1 \times 10^{-7}$).

---

## 4. Stratégie de Résolution des Conflits de Fraude (Fraud Detection & Duplicate Resolution)

Lorsqu'un doublon ou une tentative de fraude est identifié par le NIRE, les profils sont isolés selon les règles suivantes :

```
                                 RESOLVING DETECTED MATCHES
                                 
           [Biometric Match AND Demographic Match]
                      |
                      +---> IDEMPOTENT MERGE (Identité légitime unique, fusion automatique)
                      
           [Biometric Match BUT Demographic Mismatch]
                      |
                      +---> IDENTITY FRAUD DETECTED (Usurpation / double enregistrement)
                            - Lock both accounts immediately
                            - Generate PNH (Police) War Room alert
                            - Quarantine records
```

1. **Fusion Idempotente (Idempotent Merge) :**
   *   *Condition :* Correspondance biométrique validée AND Correspondance démographique $\ge 0.85$.
   *   *Action :* Les dossiers sont fusionnés automatiquement sous l'IUI le plus ancien. Les données plus récentes enrichissent le profil de manière asynchrone.
2. **Suspicion de Fraude / Usurpation (Identity Fraud Warning) :**
   *   *Condition :* Correspondance biométrique validée BUT Correspondance démographique $< 0.50$ (Même personne physique enregistrée sous deux noms totalement différents).
   *   *Action :* Les deux dossiers sont immédiatement verrouillés (Locked Status). Une alerte de fraude de niveau 1 est envoyée au bureau de liaison concerné et à la DCPJ pour enquête.
3. **Quarantaine de Rapprochement (Manual Audit) :**
   *   *Condition :* Correspondance biométrique ambiguë (Score de fusion entre 0.70 et 0.87).
   *   *Action :* Les dossiers sont envoyés vers la file d'attente d'audit humain où des experts certifiés du SNISID analysent visuellement les images d'empreintes et de visages pour statuer.
