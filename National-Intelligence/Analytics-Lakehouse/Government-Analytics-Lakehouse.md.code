# 🏞️ GOVERNMENT ANALYTICS LAKEHOUSE

> **Objectif** : Centraliser l'intelligence nationale dans un lakehouse souverain massivement scalable.

---

## 1. ARCHITECTURE MÉDAILLON

```
┌─────────────────────────────────────────────────────────┐
│  BRONZE (raw)   → Données brutes immuables              │
│  SILVER (clean) → Validées, dédupliquées, conformes     │
│  GOLD (curated) → Agrégées métier, prêtes BI/IA         │
│  PLATINUM       → Datasets décisionnels présidentiels   │
└─────────────────────────────────────────────────────────┘
```

---

## 2. CAPACITÉS SUPPORTÉES

| Domaine | Support | Mécanisme |
|---------|:-------:|-----------|
| Structured data | ✅ | Delta Lake tables ACID |
| Unstructured data | ✅ | MinIO blob + indexation Tika |
| Historical archives | ✅ | Time-travel Delta + Iceberg snapshots |
| Streaming data | ✅ | Kafka → Flink → Delta streaming sinks |
| Offline synchronization | ✅ | CDC + reprise idempotente |

---

## 3. STACK TECHNIQUE

| Domaine | Technologie | Rôle |
|---------|-------------|------|
| **Lakehouse format** | Delta Lake + Apache Iceberg | Tables ACID, schema evolution, time-travel |
| **Storage** | MinIO + Ceph | S3-compatible souverain, multi-réplica |
| **Processing batch** | Apache Spark (sur Kubernetes) | ETL massif |
| **Processing streaming** | Apache Flink | Pipelines temps réel |
| **Catalogue** | Hive Metastore + Apache Polaris | Métadonnées centralisées |
| **Orchestration** | Apache Airflow | DAGs gouvernementaux |
| **Query engine** | Trino / Presto | SQL fédéré multi-source |
| **Format colonnaire** | Parquet (compressed ZSTD) | Performance lecture |

---

## 4. DOMAINES DE DONNÉES (DATA DOMAINS)

| Domaine | Propriétaire | Tables Gold principales |
|---------|--------------|-------------------------|
| Identité | DG-SNISID | `gold.identities`, `gold.enrollments`, `gold.biometrics_meta` |
| Population | INSS | `gold.population_pyramid`, `gold.regional_density` |
| Fraude | Risk Center | `gold.fraud_alerts`, `gold.duplicate_clusters` |
| Opérations | OPS National | `gold.workflow_metrics`, `gold.service_kpis` |
| Crise | Crisis Cell | `gold.disaster_impact`, `gold.continuity_status` |
| GEOINT | Geo Cell | `gold.regional_infra`, `gold.deployment_heatmap` |
| Sécurité | CSIRT | `gold.security_events`, `gold.threat_intel` |

---

## 5. EXEMPLE DE PIPELINE (Spark / Delta)

```python
# pipeline_bronze_to_silver_enrollments.py
from pyspark.sql import SparkSession
from pyspark.sql.functions import col, to_timestamp, sha2, current_timestamp

spark = (SparkSession.builder
    .appName("snisid_enrollments_silver")
    .config("spark.sql.extensions",
            "io.delta.sql.DeltaSparkSessionExtension")
    .config("spark.sql.catalog.spark_catalog",
            "org.apache.spark.sql.delta.catalog.DeltaCatalog")
    .getOrCreate())

bronze = spark.read.format("delta").load("s3a://snisid-lake/bronze/enrollments")

silver = (bronze
    .filter(col("nin").isNotNull())
    .withColumn("event_ts", to_timestamp("event_ts"))
    .withColumn("nin_hash", sha2(col("nin"), 256))   # pseudonymisation
    .drop("nin")                                     # PII removed
    .dropDuplicates(["enrollment_id"])
    .withColumn("ingested_at", current_timestamp())
)

(silver.write
    .format("delta")
    .mode("append")
    .partitionBy("event_date")
    .option("mergeSchema", "true")
    .save("s3a://snisid-lake/silver/enrollments"))
```

---

## 6. EXEMPLE STREAMING (Flink SQL)

```sql
-- Détection en temps réel des enrôlements suspects
CREATE TABLE enrollments_kafka (
  enrollment_id STRING,
  agent_id STRING,
  region STRING,
  event_ts TIMESTAMP(3),
  WATERMARK FOR event_ts AS event_ts - INTERVAL '5' SECOND
) WITH (
  'connector' = 'kafka',
  'topic' = 'snisid.enrollments',
  'properties.bootstrap.servers' = 'kafka:9092',
  'format' = 'avro-confluent'
);

CREATE TABLE suspicious_enrollments_delta (
  agent_id STRING,
  region STRING,
  window_start TIMESTAMP(3),
  enrollment_count BIGINT
) WITH (
  'connector' = 'delta',
  'table-path' = 's3a://snisid-lake/silver/suspicious_enrollments'
);

INSERT INTO suspicious_enrollments_delta
SELECT
  agent_id,
  region,
  TUMBLE_START(event_ts, INTERVAL '1' MINUTE) AS window_start,
  COUNT(*) AS enrollment_count
FROM enrollments_kafka
GROUP BY agent_id, region, TUMBLE(event_ts, INTERVAL '1' MINUTE)
HAVING COUNT(*) > 30;   -- seuil suspect
```

---

## 7. SCALABILITÉ

| Niveau | Capacité cible |
|--------|----------------|
| Volume total | 10+ PB |
| Ingestion streaming | 500 000 événements/s |
| Requêtes concurrentes | 5 000 |
| Latence Gold query | < 2 s (P95) |
| Compaction Delta | Automatique nocturne |

**Mécanismes** :
- Auto-scaling pods Spark/Flink via KEDA
- Partitioning intelligent par date + région
- Z-Order indexing Delta sur clés analytiques
- Cache Trino par tenant

---

## 8. OFFLINE SYNC

Pour les régions à connectivité dégradée :
- Agents collectent en local (DuckDB embedded)
- Synchronisation différée via CDC compressé
- Réconciliation idempotente côté Bronze
- Mode "eventual consistency" documenté

---

## 9. SOUVERAINETÉ

- MinIO/Ceph déployés en **datacenters nationaux** uniquement
- Réplication géographique **intra-Haïti**
- Clés de chiffrement KMS souverain (HSM)
- Aucun composant SaaS externe
