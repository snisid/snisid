---
# ============================================================
# SNISID-Cyber — Security Observability Platform
# Visualisation et Tableaux de Bord SOC
# Document ID: SNISID-OBS-SEC-001
# Version: 1.0.0
# ============================================================

## 1. SECURITY OBSERVABILITY (Eyes on Glass)

L'observabilité de sécurité consolide les métriques de la Phase 5 avec une optique "Attaque/Défense". Elle est visualisée via Grafana et OpenSearch Dashboards (Kibana-like) sur les écrans géants du SOC.

## 2. DASHBOARDS CRITIQUES (Visualisations)

### 2.1 Threat Map (Carte Nationale des Menaces)
Carte géographique d'Haïti affichant le statut de chaque noeud régional (Edge Node).
- Vert : RAS.
- Jaune : Anomalie de trafic détectée (ex: Pic de bande passante inhabituel).
- Rouge : Incident de sécurité critique en cours (ex: Détection Ransomware).

### 2.2 IAM Anomalies Dashboard
Surveille l'utilisation des identités (Keycloak) en temps réel.
- **Impossible Travel Alert :** Si un utilisateur se connecte depuis Port-au-Prince, puis 10 minutes plus tard depuis une adresse IP en Russie, le compte est instantanément bloqué (Suspicious Login).

### 2.3 Attack Surface Dashboard
Montre le nombre de vulnérabilités critiques non patchées (CVEs) sur les serveurs gouvernementaux, triées par score CVSS.

---
*Document ID: SNISID-OBS-SEC-001 | Approuvé par: SOC Manager*
