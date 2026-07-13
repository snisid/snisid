---
# ============================================================
# SNISID-Cyber — National SIEM Platform
# Détection des Menaces et Journalisation Centralisée
# Document ID: SNISID-SIEM-001
# Version: 1.0.0
# ============================================================

## 1. CENTRALIZED LOGGING (La Mémoire du Système)

Le SIEM (Security Information and Event Management) ingère des dizaines de milliers de logs par seconde provenant de :
- Firewalls et Routeurs (Phase 5).
- API Gateway (Kong) et Service Mesh (Istio) (Phase 4).
- Base de données Identité (CockroachDB) (Phase 2).
- EDR (Antivirus Next-Gen) des postes de travail.

## 2. ARCHITECTURE TECHNOLOGIQUE SOUVERAINE

Pour éviter la dépendance aux licences logicielles coûteuses (Vendor Lock-in), le gouvernement utilise une stack Open Source "Enterprise Grade".

- **Agent de Collecte :** Wazuh Agent / FluentBit / Filebeat.
- **Queueing Buffer :** Apache Kafka (Évite la perte de logs si la base de données sature).
- **Moteur d'Indexation & Détection :** OpenSearch (Fork d'Elasticsearch) + Wazuh Server (Rules Engine).

### 2.1 Rétention des Logs (Compliance)
- **Hot Storage (OpenSearch) :** 90 Jours. Accès immédiat (millisecondes) pour les analystes SOC.
- **Cold Storage (S3 MinIO WORM) :** 10 Ans. Stockage immuable sur disque dur à bas coût. Impossible à effacer, même par un Super Admin. Nécessaire pour les enquêtes forensiques au long cours.

## 3. RÈGLES DE CORRÉLATION (Use Cases)

Exemple d'une règle (SIGMA/Wazuh) : "Détection de fuite de données d'État Civil".
```yaml
rule:
  name: "Massive Export of Identity Records"
  condition: "User == ANY_OFFICIER_ETAT_CIVIL AND API_Endpoint == '/v1/identities/*' AND Query_Count > 500 AND Timeframe < 10m"
  action: "Trigger_Alert_CRITICAL"
```
Si un officier d'état civil interroge plus de 500 dossiers en moins de 10 minutes (comportement non-humain/scraping), une alerte critique est levée et son compte Keycloak est temporairement suspendu (via SOAR).

---
*Document ID: SNISID-SIEM-001 | Approuvé par: Head of Threat Detection*
