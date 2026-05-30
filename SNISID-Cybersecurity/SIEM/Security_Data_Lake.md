# SNISID National Security Data Lake

## 1. Objective
To create a centralized, high-scale telemetry repository that serves as the "National Cyber Memory," enabling long-term analysis, threat hunting, and forensic reconstruction.

## 2. Data Architecture
Unlike the SIEM (which is optimized for real-time alerting), the Data Lake is optimized for mass storage and complex analytical queries.

### Storage Layers
- **Landing Zone:** Raw logs ingested in original format.
- **Processed Zone:** Normalized and enriched data (e.g., IP addresses resolved to GeoIP).
- **Curated Zone:** Aggregated security metrics and high-value datasets.

## 3. Data Domains to Ingest

| Domain | Data Type | Purpose |
| :--- | :--- | :--- |
| **Logs** | Syslog, Event Logs, App Logs | Baseline activity and audit trails. |
| **Security Events** | SIEM Alerts, EDR Detections | History of attacks and response outcomes. |
| **Threat Intelligence** | STIX/TAXII feeds, IOCs | Correlating historical data with new intel. |
| **IAM Telemetry** | Auth logs, Token usage, Permissions | Tracking identity movement over time. |
| **Network Telemetry** | Netflow, DNS logs, VPC Flow Logs | Mapping the national communication graph. |

## 4. Technology Stack
- **Storage:** Amazon S3 / Azure Blob Storage / MinIO (On-prem).
- **Query Engine:** Trino / Presto / Apache Spark.
- **Catalog:** AWS Glue / Apache Hive Metastore.
- **Format:** Apache Parquet (Columnar storage for efficiency).

## 5. Use Cases
- **Retroactive Hunting:** "We just found a new IOC from 6 months ago; was it ever in our network?"
- **Trend Analysis:** Identifying gradual shifts in adversary behavior over years.
- **Compliance:** Meeting national mandates for long-term data retention.
- **ML Training:** Providing a massive dataset to train national anomaly detection models.
