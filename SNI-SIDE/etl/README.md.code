# SNI-SIDE: ETL Migration Toolkit
## Import des données legacy vers les 15 bases nationales

## Vue d'ensemble

```
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────────┐
│ CSV / Excel      │────>│                  │     │   PostgreSQL 16      │
│ Legacy DB        │────>│  ETL Pipeline    │────>│   CockroachDB        │
│ PDF / Scans      │────>│  (Python)        │     │   MinIO (documents)  │
│ API Legacy       │────>│                  │     │   Neo4j (graphe)     │
└──────────────────┘     └──────────────────┘     └──────────────────────┘
                               │
                               v
                          ┌────────────┐
                          │  Kafka     │
                          │  Events    │
                          └────────────┘
```

## Sources Supportées

| Source | Format | Adaptateur |
|:--|:--|:--|
| Fichiers Excel (XLSX) | `.xlsx` | `ExcelAdapter` |
| CSV (tous séparateurs) | `.csv` | `CsvAdapter` |
| PostgreSQL legacy | SQL | `PostgresLegacyAdapter` |
| MySQL/MariaDB legacy | SQL | `MysqlLegacyAdapter` |
| MongoDB legacy | JSON/BSON | `MongoAdapter` |
| PDF/Scans textuels | `.pdf` | `PdfAdapter` |
| API REST legacy | HTTP | `RestApiAdapter` |
| Fichiers JSON | `.json` | `JsonAdapter` |
| Fichiers XML | `.xml` | `XmlAdapter` |

## Types de Migration

### 1. Migration Simple (1 fichier → 1 table)
```
python etl/run.py --source persons.csv --target snisid_ncid.citizens --mapping mapping/ncid_persons.json
```

### 2. Migration Complexe (1 source → plusieurs tables liées)
```
python etl/run.py --source crimedata.xlsx --pipeline pipelines/ncid_case_pipeline.json
```

### 3. Migration Full (base legacy complète)
```
python etl/run.py --pipeline pipelines/full_migration.json --kafka-events true
```

## Mapping JSON

Chaque fichier de mapping définit la correspondance champs source → cible :

```json
{
  "target_schema": "snisid_ncid",
  "target_table": "wanted_persons",
  "batch_size": 1000,
  "field_mappings": [
    {"source": "CASE_NUMBER", "target": "case_reference", "type": "string", "required": true},
    {"source": "NOM", "target": "last_name", "type": "string", "transformer": "upper"},
    {"source": "PRENOM", "target": "first_name", "type": "string", "transformer": "capitalize"},
    {"source": "DATE_NAISS", "target": "birth_date", "type": "date", "format": "%d/%m/%Y"},
    {"source": "SEXE", "target": "gender", "type": "enum", "mapping": {"M": "MALE", "F": "FEMALE", "X": "OTHER"}},
    {"source": "NIU", "target": "niu", "type": "string", "generate_if_missing": true},
    {"source": "PHOTO_PATH", "target": "photo_url", "type": "file", "copy_to": "minio://sniside-media/photos/"},
    {"source": "CATEGORIE", "target": "risk_level", "type": "enum",
     "mapping": {"A": "CRITICAL", "B": "HIGH", "C": "MEDIUM", "D": "LOW"}},
    {"source": "COMMENTAIRES", "target": "notes", "type": "text"}
  ],
  "validations": [
    {"field": "niu", "rule": "niu_format"},
    {"field": "birth_date", "rule": "not_future"}
  ],
  "post_process": [
    {"action": "emit_kafka", "topic": "sniside.ncid.wanted.created"},
    {"action": "update_neo4j", "label": "Citizen"}
  ]
}
```

## Pipelines Prédéfinis

### NCID (National Criminal Intelligence Database)
| Fichier | Mapping | Tables cibles |
|:--|:--|:--|
| `wanted_persons.csv` | `mapping/ncid_wanted.json` | wanted_persons, warrants, aliases |
| `criminal_cases.xlsx` | `mapping/ncid_cases.json` | cases, case_subjects, case_evidence |
| `gangs.csv` | `mapping/ncid_gangs.json` | gangs, gang_members, gang_territories |

### Biometrics (HN-NGI)
| Fichier | Mapping | Tables cibles |
|:--|:--|:--|
| `fingerprints.csv` | `mapping/biometric_fingerprints.json` | fingerprint_templates |
| `face_images.csv` | `mapping/biometric_faces.json` | facial_images (MinIO + Milvus) |
| `iris_scans.csv` | `mapping/biometric_iris.json` | iris_templates |
| `voice_samples.csv` | `mapping/biometric_voice.json` | voice_templates |

### CODIS (DNA)
| Fichier | Mapping | Tables cibles |
|:--|:--|:--|
| `dna_profiles.csv` | `mapping/codis_profiles.json` | dna_profiles |
| `crime_scene_dna.xlsx` | `mapping/codis_scene.json` | crime_scene_dna, dna_matches |

### Others (similar structure for all 15 bases)

## Transformateurs Disponibles

| Transformer | Description |
|:--|:--|
| `upper` | Met en majuscules |
| `lower` | Met en minuscules |
| `capitalize` | Première lettre en majuscule |
| `trim` | Supprime espaces |
| `normalize_name` | Normalise les accents/casse |
| `niu_generator` | Génère un NIU valide (10 caractères) |
| `date_parse` | Parse une date selon format |
| `phone_normalize` | Normalise numéro en international |
| `plate_normalize` | Normalise plaque (AZ-123-AB) |
| `hash_pii` | Hash PII pour anonymisation |
| `file_copy` | Copie fichier vers MinIO |
| `base64_decode` | Décode base64 vers fichier |
| `enum_map` | Mapping valeur → valeur cible |
| `concat` | Concatène plusieurs champs |
| `split` | Split un champ en plusieurs |
| `regex_extract` | Extraction par regex |
| `lookup` | Recherche dans table de référence |
| `geo_encode` | Coordonnées → PostGIS point |

## Validateurs Disponibles

| Validateur | Description |
|:--|:--|
| `niu_format` | Valide format NIU (10 chars alphanum) |
| `niu_unique` | Vérifie unicité du NIU |
| `email` | Valide format email |
| `phone` | Valide format téléphone |
| `plate` | Valide format plaque |
| `date` | Valide date |
| `not_future` | Vérifie que la date n'est pas future |
| `required` | Champ obligatoire |
| `min_length` | Longueur minimale |
| `max_length` | Longueur maximale |
| `regex` | Validation par expression régulière |
| `enum` | Valeur dans liste autorisée |
| `range` | Valeur dans intervalle |
| `unique_check` | Vérifie unicité dans table cible |
| `foreign_key` | Vérifie existence FK dans table source |

## Utilisation

### Migration rapide d'un CSV
```bash
python run.py --source legacy/wanted.csv \
    --target snisid_ncid.wanted_persons \
    --mapping mappings/ncid_wanted.json \
    --batch-size 500 \
    --kafka-events true
```

### Migration complète d'une base legacy
```bash
python run.py --pipeline pipelines/full_migration.json \
    --db-source postgresql://legacy:password@192.168.1.100:5432/crimedb \
    --db-target postgresql://sniside:password@postgres:5432/sniside \
    --kafka-events true \
    --neo4j-update true \
    --dry-run false
```

### Validation uniquement (sans écriture)
```bash
python run.py --pipeline pipelines/ncid_migration.json --dry-run true
```

### Résumé de migration
```bash
python run.py --pipeline pipelines/full_migration.json --report-only true
```

## Architecture ETL

```
run.py                          → Point d'entrée CLI
├── adapters/                   → Adaptateurs source
│   ├── csv_adapter.py          → Lecture CSV
│   ├── excel_adapter.py        → Lecture Excel
│   ├── db_adapter.py           → Lecture base SQL
│   ├── pdf_adapter.py          → Extraction PDF
│   ├── json_adapter.py         → Lecture JSON
│   ├── xml_adapter.py          → Lecture XML
│   └── rest_adapter.py         → API REST legacy
│
├── transformers/               → Transformateurs
│   ├── string_transformers.py  → upper, lower, trim, etc.
│   ├── date_transformers.py    → parse, format, etc.
│   ├── id_generators.py        → NIU generator, UUID
│   ├── file_handlers.py        → MinIO copy, base64
│   └── lookup_transformers.py  → FK lookup, geo encoding
│
├── validators/                 → Validateurs
│   ├── field_validators.py     → par champ
│   ├── row_validators.py       → validations inter-champs
│   └── batch_validators.py     → validations inter-lignes
│
├── writers/                    → Cibles
│   ├── pg_writer.py            → PostgreSQL 16
│   ├── cockroach_writer.py     → CockroachDB
│   ├── neo4j_writer.py         → Neo4j graph
│   ├── minio_writer.py         → MinIO objects
│   ├── milvus_writer.py        → Milvus vectors
│   └── kafka_writer.py         → Kafka events
│
├── pipeline.py                 → Orchestrateur de pipeline
├── mapping.py                  → Mapping engine
├── models.py                   → Data models
├── progress.py                 → Progress tracking
└── report.py                   → Migration report
```
