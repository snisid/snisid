# SNISID Migration Factory Architecture
## Industrialisation du Processus de Migration Nationale de Données

---

## 1. Vue d'Ensemble de la Migration Factory

La **Migration Factory** est l'infrastructure applicative et opérationnelle chargée de migrer plus de 6 millions de dossiers d'identité historiques (provenant de l'ancien système de l'ONI, de l'état civil, des archives papier et d'autres bases institutionnelles) vers la base de données unifiée et sécurisée du SNISID. 

```
                                  MIGRATION FACTORY PIPELINE
                                  
+------------------------+      +------------------------+      +------------------------+
|   Legacy Raw Ingest    | ---> | Data Cleansing Program | ---> | Identity Reconciliation |
| (ONI, Civil Registry)  |      |  (Clean, Format, Auth) |      | (Demographics/Biometrics)
+------------------------+      +------------------------+      +------------------------+
                                                                             |
                                                                             v
+------------------------+      +------------------------+      +------------------------+
|  SNISID Verified DB    | <--- |  Final Sign-Off / Audit | <--- |  Validation & Quality  |
|  (Production Cluster)  |      |  (Reversible Rollback) |      |   (Anomaly Detection)  |
+------------------------+      +------------------------+      +------------------------+
```

---

## 2. Étapes du Pipeline d'Extraction-Transformation-Chargement (ETL) Sécurisé

Le pipeline de migration s'exécute dans une enclave réseau isolée et hautement sécurisée (Data Migration Zone - DMZ-M) au sein du Datacenter Central.

### 2.1 Les Cinq Étapes Clés
1. **Ingestion Sécurisée (Secure Ingestion) :**
   *   Importation de sauvegardes chiffrées au format PostgreSQL, MySQL ou CSV.
   *   Calcul de l'empreinte cryptographique (SHA-256) du fichier source pour en garantir l'intégrité avant traitement.
2. **Nettoyage Automatique (Cleansing Program) :**
   *   Correction des structures de noms/prénoms (suppression des caractères spéciaux, normalisation de la casse, gestion des apostrophes complexes créoles et françaises).
   *   Validation et reformatage des dates de naissance (détection des dates impossibles comme le `29 Février 1900` ou les dates futures).
3. **Réconciliation d'Identité (Identity Reconciliation) :**
   *   Rapprochement démographique (Algorithme de similarité phonétique adapté au créole haïtien et au français, comme Soundex ou Levenshtein double).
   *   Rapprochement biométrique (Deduplication un-à-plusieurs via le moteur d'ABIS central).
4. **Enregistrement et Auditabilité :**
   *   Génération d'un identifiant universel d'identité (IUI) pour chaque fiche réconciliée.
   *   Génération de métadonnées de lignage de données (*Data Lineage*) permettant de remonter à la source exacte de chaque attribut migré (ex: nom issu de l'ONI, date de naissance issue de l'Acte de Naissance de l'ANH).
5. **Porte de Contrôle et Chargement (Load & Gate) :**
   *   Validation de la conformité du lot de données à 100%.
   *   Écriture transactionnelle dans la base de production du SNISID avec conservation d'un instantané de restauration immédiate (*Restore Point*) en cas d'incident.

---

## 3. Protocoles d'Auditabilité et de Restauration (Rollback Support)

Chaque exécution de migration de lot de données (*Migration Batch*) doit être réversible à 100%. 

*   **Journalisation transactionnelle (Transaction WAL) :** Toutes les opérations d'écriture s'accompagnent de la création d'un journal des transactions inversées (Undo Log).
*   **Versionnage des entités :** Le SNISID utilise le patron d'architecture d'historisation (*Event Sourcing / Temporal Tables*). Si une anomalie est détectée après la migration d'un lot, une commande de rollback de lot permet de restaurer l'état de la base à la microseconde précédant l'exécution du lot.
*   **Signature électronique :** Chaque lot de données migré est signé cryptographiquement par la clé privée de l'ingénieur responsable de la migration après validation humaine par le comité d'audit.

---

## 4. Spécifications du Pipeline Applicatif

Pour automatiser ce flux, la *Migration Factory* déploie des services conteneurisés en Go et Python, s'interfaçant avec Apache Spark pour le traitement de gros volumes de données. La logique détaillée de déduplication et de réconciliation est documentée dans les programmes et simulateurs de ce répertoire.
