---
# ============================================================
# SNISID-Interop — National Operations Dashboards
# Pilotage de l'Interopérabilité
# Document ID: SNISID-INTEROP-DASHBOARDS-001
# Version: 1.0.0
# ============================================================

## 1. EXECUTIVE INTEROPERABILITY DASHBOARD

Vue globale pour le Premier Ministre et le Ministre de l'Économie (KPIs économiques et d'efficacité).
- **Volume de transactions API par Ministère** (ex: 50,000 req/jour pour la DGI).
- **Taux de création d'entreprises 100% numériques** (via Orchestrateur BPMN).
- **Temps moyen de résolution des workflows inter-agences** (SLA Tracking).

## 2. TECHNICAL & SLA DASHBOARDS (NOC / SOC)

Vue pour les ingénieurs (Network/Security Operations Center).
- **API Latency (P95, P99) :** Suivi de la performance de l'API Gateway.
- **Kafka Consumer Lag :** Si la DGI a 10,000 messages de retard sur le bus, une alerte est levée.
- **Error Budgets (SRE) :** Suivi du SLA global 99.99%. Si les erreurs 5xx de la passerelle dépassent 0.01% sur 30 jours, les déploiements non-urgents sont gelés.

## 3. ANOMALY DETECTION DASHBOARD (Security)

- Alertes sur les pics inhabituels de requêtes (ex: Une banque demande 10,000 KYC en 1 heure).
- Alertes DLP (Data Loss Prevention) depuis Istio.
- Tentatives d'authentification refusées (Brute force OIDC).

---
*Document ID: SNISID-INTEROP-DASHBOARDS-001 | Approuvé par: Centre National des Opérations Numériques*
