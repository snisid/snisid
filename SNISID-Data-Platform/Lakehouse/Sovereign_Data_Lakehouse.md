# Sovereign Data Lakehouse SNISID

## Objectif
Créer une plateforme analytique nationale **souveraine, distribuée, gouvernée et auditée**.

## Stack recommandée

| Domaine | Technologie |
|---|---|
| Data Lake | MinIO compatible S3, déployé sur infrastructure nationale |
| Table Format | Apache Iceberg |
| Query Engine | Trino |
| Batch Processing | Apache Spark |
| Streaming | Kafka + Kafka Connect |
| Catalog | Iceberg REST Catalog/Hive Metastore + OpenMetadata |
| Orchestration | Airflow/Argo Workflows |
| Sécurité | TLS, KMS/HSM, IAM, Ranger/Open Policy Agent |

## Architecture logique

```text
Sources nationales -> Kafka/CDC -> Landing -> Bronze -> Silver -> Gold -> BI/API/AI
                                \-> Audit Fabric
Metadata/Lineage <----------------------------------------- chaque étape
```

## Exigences de souveraineté

- Hébergement sur cloud gouvernemental, datacenters nationaux ou infrastructure sous juridiction haïtienne.
- Chiffrement avec clés contrôlées par l'État.
- Réplication inter-sites souveraine.
- Interdiction de stockage non approuvé hors juridiction.
- Journalisation complète des accès et exports.

## Zones et politiques

| Zone | Format | Rétention par défaut | Immutabilité | Qualité requise |
|---|---|---:|---|---|
| Landing | Parquet/JSON/Avro brut | 90 jours | Oui | Signature + checksum |
| Bronze | Iceberg | 2 ans | Versionnée | Schéma valide |
| Silver | Iceberg | 7 ans | Versionnée | DQ score >= 95% |
| Gold | Iceberg | Selon domaine | Versionnée | Certifié owner |
| Archive | Iceberg/Parquet WORM | 10 ans à vie | WORM | Contrôle intégrité |

## Patterns d'ingestion

1. **CDC transactionnel** : Debezium -> Kafka -> Iceberg.
2. **Batch sécurisé** : SFTP/API -> Landing -> validation -> Bronze.
3. **Streaming events** : Kafka topics versionnés -> stream processing -> Silver/Gold.
4. **Data sharing inter-agences** : API + data contract + audit.

## Contrôles obligatoires avant publication Gold

- Dataset catalogué.
- Owner et steward identifiés.
- Classification validée.
- Contrats de données enregistrés.
- Score qualité conforme.
- Lineage complet.
- Politique d'accès approuvée.
- Rétention définie.
