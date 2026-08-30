---
# ============================================================
# SNISID-Data — Real-Time Analytics Platform
# Moteur de Streaming et Alerting Immédiat
# Document ID: SNISID-DATA-STREAM-001
# Version: 1.0.0
# ============================================================

## 1. STREAM PROCESSING (Le temps réel)

Le Lakehouse (Iceberg/Trino) est conçu pour l'analyse historique profonde (Big Data). Mais certaines situations exigent une réaction à la milliseconde (ex: Détection de fraude bancaire croisée avec l'identité, Alertes de sécurité).

## 2. ARCHITECTURE DE STREAMING (Flink / ClickHouse)

- **Le Cœur (Apache Flink) :** Se branche directement sur les topics Kafka (Phase 4). Flink analyse les flux d'événements à la volée sans les écrire sur le disque.
- **La Base Ultra-Rapide (ClickHouse) :** Flink pousse les résultats des calculs dans ClickHouse, une base de données analytique en colonnes optimisée pour le temps réel (OLAP).

## 3. CAS D'USAGE GOUVERNEMENTAUX (Use Cases)

1. **Détection de "Voyage Impossible" (Impossible Travel) :** 
   - *Événement A :* Un passeport est scanné à la frontière de Ouanaminthe.
   - *Événement B :* La carte d'identité du même citoyen est utilisée pour un retrait bancaire à Jérémie (10 heures de route), 30 minutes plus tard.
   - *Action Flink :* Détecte l'anomalie en temps réel et envoie une alerte critique au SOC (Phase 6) et à la DCPJ (Phase 3).
2. **Surveillance des Opérations Mobiles (Phase 8) :**
   - Suivi en temps réel de la température des batteries au lithium des camions MGU via la télémétrie ingérée à très haute fréquence.

---
*Document ID: SNISID-DATA-STREAM-001 | Approuvé par: Chief Data Officer (CDO) National*
