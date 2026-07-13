---
# ============================================================
# SNISID-Data — Sovereign Data Lakehouse Architecture
# Stockage Massif, Modélisation et Requêtage
# Document ID: SNISID-DATA-LAKE-001
# Version: 1.0.0
# ============================================================

## 1. VISION: LE CERVEAU DE L'ÉTAT

L'architecture SNISID génère des millions d'événements par jour via son bus Kafka (Phase 4). Le "Sovereign Data Lakehouse" est l'endroit où tous ces événements convergent pour constituer la mémoire analytique immuable de la nation.

## 2. ARCHITECTURE TECHNIQUE (Iceberg / MinIO / Trino)

Plutôt que d'utiliser des data warehouses propriétaires coûteux (Snowflake/BigQuery), l'État Haïtien utilise une architecture 100% open-source hébergée localement sur les clusters Kubernetes (Phase 5).

- **Storage Layer (Le Lac) :** MinIO (Stockage Objet S3-compatible). Permet de stocker des Pétaoctets de données (Images biométriques, Logs SOC, Événements d'État Civil) à bas coût.
- **Table Format :** Apache Iceberg. Apporte les transactions ACID (Atomicity, Consistency, Isolation, Durability) directement sur le Data Lake, évitant la corruption des données en cas de panne de cluster.
- **Compute Layer (Le Moteur) :** Trino (anciennement Presto). Moteur SQL massivement parallèle permettant aux data scientists du gouvernement de faire des requêtes complexes (ex: "Croiser la base de la Police avec les logs de la Frontière") en quelques secondes.

## 3. MODÉLISATION EN MÉDAILLON (Medallion Architecture)

Les données ne sont jamais modifiées dans leur état brut.
- **Bronze (Brut) :** Ingestion directe depuis Kafka (JSON/Avro). Aucune transformation. L'historique immuable de l'État.
- **Silver (Nettoyé) :** Données dédoublonnées, formats de dates normalisés, numéros de téléphone validés.
- **Gold (Agrégé) :** Vues analytiques prêtes à la consommation (ex: "Tableau de Bord des Naissances par Département", "KPI des interventions de Police").

---
*Document ID: SNISID-DATA-LAKE-001 | Approuvé par: Chief Data Officer (CDO) National*
