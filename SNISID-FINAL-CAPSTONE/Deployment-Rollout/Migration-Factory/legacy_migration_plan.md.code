# SNISID Legacy Systems Migration Plan
## Plan Stratégique de Migration des Anciens Systèmes Gouvernementaux d'Haïti

---

## 1. Cartographie des Systèmes Legacy à Migrer

Pour constituer la base de données unifiée du SNISID, cinq grands systèmes historiques indépendants et hétérogènes doivent être consolidés et migrés. Chacun de ces systèmes présente des défis d'intégrité de données, de formats de stockage obsolètes et de doublons massifs.

```
+---------------------------------------------------------------------------------+
|                               MIGRATION TARGETS                                 |
+---------------------------------------------------------------------------------+
  | (ONI Legacy)      | (Civil Registry)  | (Police Records) | (Judiciary Sys)  | (Immigration)
  | - Dec-Alpha SQL   | - Actes Papier    | - Fiches PNH     | - Registre MJSP  | - Sygma DIE
  v                   v                   v                  v                  v
+---------------------------------------------------------------------------------+
|                       SNISID SECURE ADAPTER LAYER (ETL)                         |
+---------------------------------------------------------------------------------+
                                      |
                                      v
+---------------------------------------------------------------------------------+
|                       SNISID RECONCILED NATIONAL SYSTEM                         |
+---------------------------------------------------------------------------------+
```

---

## 2. Stratégies de Migration par Domaine

### 2.1 Office National d'Identification (ONI Legacy)
*   **Description :** Base de données des anciennes Cartes d'Identification Nationale (CIN). Environ 5.8 millions d'enregistrements.
*   **Format Source :** Base SQL Server héritée et fichiers plats indexés. Biométrie au format propriétaire WSQ (Wavelet Scalar Quantization) pour les empreintes.
*   **Stratégie de Migration :**
    1. Extraction de la base démographique et conversion en UTF-8.
    2. Extraction des fichiers d'empreintes WSQ, conversion par lot aux normes ouvertes ISO/IEC 19794-4 (Finger Minutiae Data) et NIST Record Type-9.
    3. Injection des données démographiques nettoyées dans l'enclave temporaire de réconciliation.
    4. Envoi des empreintes à l'ABIS central pour création de la galerie biométrique maîtresse et déduplication automatique.

### 2.2 État Civil et Archives Nationales d'Haïti (Civil Registry)
*   **Description :** Actes de naissance, mariage, divorce et décès. Des millions d'actes papier, certains numérisés sous forme d'images JPEG non structurées.
*   **Format Source :** Registres physiques (papier), base de données de métadonnées partielles (MS Access, SQL Server obsolètes).
*   **Stratégie de Migration :**
    1. Déploiement d'un pipeline d'OCR (Reconnaissance Optique de Caractères) intelligent, pré-entraîné sur l'écriture cursive historique haïtienne, couplé à une double saisie manuelle de contrôle pour les actes illisibles.
    2. Rapprochement par filiation : structuration des relations parentales (père, mère) pour construire le premier graphe relationnel de parenté nationale (Arbre de Filiation Civil).
    3. Résolution des dates aberrantes et des homonymies géographiques (ex: communes nées d'un découpage territorial récent).

### 2.3 Fiches de la Police Nationale d'Haïti (Police Records)
*   **Description :** Registres d'antécédents judiciaires, mandats de recherche, casiers judiciaires de la DCPJ.
*   **Format Source :** Fiches papier indexées, base de données locale isolée.
*   **Stratégie de Migration :**
    *   *Avertissement Critique :* Pour des raisons de protection des libertés individuelles et de sécurité nationale, ces données ne sont **jamais fusionnées** de manière visible dans l'identité civile.
    *   *Mécanisme :* Hachage cryptographique fort (HMAC-SHA-256 avec sel dynamique géré par le Ministère de la Justice) des attributs de base (Nom + Prénom + Date de Naissance + Empreintes digitales). Seul ce "token d'antécédent" est enregistré dans un index de sécurité crypté. Si un citoyen est enrôlé, la gateway interroge de manière aveugle cet index. En cas de correspondance positive ("Hit"), une alerte sécurisée cryptée est envoyée à la DCPJ pour traitement manuel.

### 2.4 Registres du Système Judiciaire (Judicial Systems / MJSP)
*   **Description :** Décisions de déchéance des droits civiques, peines de prison fermes, décrets de naturalisation.
*   **Format Source :** Registres papier des greffes des tribunaux de première instance (TPI).
*   **Stratégie de Migration :**
    1. Saisie numérique des ordonnances de restriction des droits civiques en cours de validité.
    2. Liaison stricte avec l'IUI (Identifiant Universel d'Identité) du SNISID pour marquer l'état civil de l'individu (ex: interdiction de vote ou d'émission de passeport).

### 2.5 Système d'Immigration (Immigration Systems / DIE)
*   **Description :** Base de données des passeports émis par la Direction de l'Immigration et de l'Émigration (DIE) et fiches d'entrées/sorties transfrontalières (Système SECURI-PORT ou SYGMA).
*   **Format Source :** Base Oracle et flux de données XML.
*   **Stratégie de Migration :**
    1. Rapprochement direct entre le numéro de passeport historique et le numéro d'identité nationale CIN.
    2. Établissement d'une passerelle d'API temps réel (gRPC) pour alimenter le dossier de voyageur du citoyen dans le SNISID, sécurisant la validation d'identité biométrique aux douanes.

---

## 3. Matrice de Criticité et Risques de Perte de Données

| Source | Niveau de Risque de Perte | Mesure de Protection d'Intégrité | Critère de Validation Finale |
| :--- | :--- | :--- | :--- |
| **ONI Legacy** | Très Élevé (Corruption de tables) | Copie physique miroir en lecture seule sur stockage WORM (Write Once Read Many). | Taux d'écart de bilan comptable entre fichiers lus et écrits inférieur à 0.0001%. |
| **Civil Registry** | Critique (Perte papier, incendies) | Numérisation haute définition (TIFF non compressé) avec stockage répliqué sur 2 sites physiques géographiquement distants. | Validation d'intégrité cryptographique des images stockées. |
| **DIE Passports** | Moyen | Synchronisation différentielle journalière et double écriture transactionnelle. | Rapport de conformité de l'audit de clé étrangère unifiée. |
| **Police Records** | Élevé | Chiffrement matériel (HSM - Hardware Security Module) empêchant toute extraction illégale des clés d'anonymisation. | Vérification par test d'intrusion externe de l'impossibilité de fuite démographique. |
