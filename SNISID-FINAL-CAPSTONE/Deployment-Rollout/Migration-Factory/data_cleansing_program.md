# SNISID National Data Cleansing Program
## Programme National d'Assainissement et de Nettoyage de l'Identité

---

## 1. Introduction

Les bases de données historiques d'identification en Haïti souffrent de taux de bruit et d'incohérence extrêmement élevés (estimés à plus de 18% d'anomalies structurelles). Ces problèmes proviennent d'erreurs de saisie manuelle, de l'absence de normes orthographiques pour le créole haïtien dans les anciens systèmes, et de corruptions matérielles de bases de données obsolètes. Le **National Data Cleansing Program** est la charte d'assainissement systématique des données avant leur réconciliation biométrique.

```
                         DATA CLEANSING FLOW
                         
     +---------------------------------------------------------+
     |                  Raw Ingested Record                    |
     +---------------------------------------------------------+
                                  |
                                  v
     +---------------------------------------------------------+
     | Step 1: Structural Repair (Regex, Encoding, UTF-8)      |
     +---------------------------------------------------------+
                                  |
                                  v
     +---------------------------------------------------------+
     | Step 2: Attribute Completeness Check & Reconstruction   |
     +---------------------------------------------------------+
                                  |
                                  v
     +---------------------------------------------------------+
     | Step 3: Conflict Resolution & Heuristics                |
     | (Name Normalization, Phone/Email/Address Validation)    |
     +---------------------------------------------------------+
                                  |
                                  v
     +---------------------------------------------------------+
     | Cleaned Record Ready for Identity Reconciliation Engine |
     +---------------------------------------------------------+
```

---

## 2. Détection et Traitement des Anomalies Structurelles

Le programme d'assainissement cible 5 grandes catégories de défaillances de données :

### 2.1 Détection et Résolution des Doublons (Duplicate Detection)
Les doublons démographiques parfaits (mêmes noms, prénoms, date et lieu de naissance) ou partiels sont identifiés en appliquant un algorithme de calcul de distance textuelle.
*   **Règle :** Si la distance de Jaro-Winkler entre deux fiches est supérieure à 0.96 et que l'écart sur les dates de naissance est nul, les fiches sont marquées comme "Doublons potentiels" et transmises à l'étape de réconciliation biométrique pour validation définitive.

### 2.2 Conflits d'Identité Active (Identity Conflicts)
Un conflit d'identité survient lorsque deux fiches distinctes possèdent des attributs exclusifs identiques, comme le même numéro de CIN historique (Carte d'Identification Nationale) pour deux personnes différentes.
*   **Résolution :** L'ancien numéro de CIN est désactivé et marqué comme "Historique conflictuel". Le SNISID attribue un nouvel **Identifiant Universel d'Identité (IUI)** unique et sécurisé aux deux individus et planifie une convocation physique pour un nouvel enrôlement biométrique de contrôle.

### 2.3 Enregistrements Invalides (Invalid Records)
Les enregistrements comportant des données absurdes ou contraires aux lois de la nature sont rejetés ou mis en quarantaine.
*   **Cas Détectés :**
    *   *Âges aberrants :* Personnes nées avant 1900 (plus de 126 ans) ou enfants enrôlés comme chefs de ménage.
    *   *Séquences impossibles :* Dates d'émission de cartes antérieures aux dates de naissance, ou dates d'enrôlement futures.
*   **Résolution :** Transfert automatique de l'enregistrement vers la base de quarantaine de niveau 1 (Q1) avec une étiquette indiquant l'anomalie temporelle détectée.

### 2.4 Données Corrompues (Corrupted Data)
Chaines de caractères incompréhensibles résultant d'erreurs d'encodage (ex: `Mlissa` au lieu de `Mélissa`, ou caractères de type `null`, `NaN`, `N/A`, `---`).
*   **Résolution :** Application d'un décodeur d'encodage universel (UTF-8 / ISO-8859-1). Les valeurs factices (`null`, `N/A`, `test`, `inconnu`) sont purgées de la base de données et remplacées par la valeur système formelle `UNKNOWN` pour exiger une saisie rectificative lors de la prochaine mise à jour physique.

### 2.5 Attributs Absents (Missing Attributes)
Il s'agit d'enregistrements dont certains champs indispensables à l'identité sont vides (comme le lieu de naissance, le prénom de la mère ou le sexe).
*   **Règle de Complétude :** Un enregistrement d'identité doit contenir à minima les champs obligatoires suivants pour être migré en production :
    *   Nom et au moins un Prénom.
    *   Date de Naissance complète (ou année certifiée si l'acte de naissance est incomplet).
    *   Sexe.
    *   Lieu de naissance (Commune et Section communale).
*   **Résolution :** Si un champ obligatoire manque, l'enregistrement est classé en quarantaine de niveau 2 (Q2). Il ne peut pas être basculé en production active tant qu'une vérification manuelle aux registres papier de l'état civil n'a pas été effectuée pour combler l'attribut absent.

---

## 3. Règles de Normalisation Orthographique (Haitian Creole & French)

Les noms de famille en Haïti présentent des variations orthographiques fréquentes pour une même lignée familiale (ex : `Jean-Baptiste` écrit `Janbatis` ou `Jean Baptiste` avec espace, `Hyppolite` écrit `Hippolyte`).

### 3.1 Règles Algorithmiques Appliquées par le Pipeline de Nettoyage
1. **Suppression des espaces superflus et caractères spéciaux :** Les caractères autres que les lettres de l'alphabet latin, les tirets (`-`) et les apostrophes (`'`) sont supprimés. Les espaces multiples sont réduits à un espace unique.
2. **Normalisation de la casse :** Conversion systématique des noms de famille en majuscules (`JEAN-BAPTISTE`) et de la première lettre des prénoms en majuscule (`Melissa`).
3. **Phonétisation Spécifique (Moteur Soundex-HT) :** Un algorithme Soundex personnalisé, configuré pour prendre en compte les spécificités phonétiques du créole haïtien et du français (ex: `ch` phonétisé comme `sh`, `an` équivalent à `en`, `y` traité comme `i`, suppression du `h` muet au début des noms comme `Hyppolite` -> `Ipolite`).

---

## 4. Tableau des Seuils de Tolérance Qualité

La *Migration Factory* applique des indicateurs de qualité extrêmement stricts. Aucun lot de données ne peut être déployé en production si les taux d'anomalies après nettoyage dépassent les seuils critiques suivants :

```
+------------------------------------+--------------------------+
| Type d'Anomalie                    | Seuil de Tolérance Max   |
+------------------------------------+--------------------------+
| Doublons démographiques non résolus| 0.00% (Zéro Tolérance)   |
| Date de naissance invalide         | 0.02% (2 pour 10 000)    |
| Attributs obligatoires absents     | 0.10% (1 pour 1 000)     |
| Erreurs d'encodage (Caractères ?)  | 0.00% (Zéro Tolérance)   |
+------------------------------------+--------------------------+
```
