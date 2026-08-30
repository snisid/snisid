---
# ============================================================
# SNISID-Field — National Field Observability
# Télémétrie Mobile et Surveillance des Équipements
# Document ID: SNISID-TELEMETRY-001
# Version: 1.0.0
# ============================================================

## 1. OBSERVABILITÉ DE LA FLOTTE (Mobile Telemetry)

Le FOC (Field Operations Center) voit la flotte gouvernementale comme un système IOT géant.

## 2. DEVICE HEALTH DASHBOARDS

Les tableaux de bord Grafana du FOC affichent pour chaque MGU (Camion) et Tablette :
- **Santé Énergétique :** Niveau de charge des batteries LiFePO4, production des panneaux solaires en Watts.
- **Santé Applicative :** Nombre de "Crash" de l'application React Native d'enrôlement.
- **Santé Sync :** Taille de la file d'attente Kafka/NATS (Combien d'enregistrements attendent d'être synchronisés).
- **Santé Environnementale :** Température interne du rack serveur du camion (Alerte si > 45°C en plein soleil).

## 3. FIELD SLA MONITORING

Le système calcule l'efficacité des équipes : "L'équipe du Sud a enrôlé 450 citoyens aujourd'hui, avec un taux de rejet biométrique (NFIQ) de 5%". Si le taux de rejet monte à 30%, le superviseur est alerté que le scanner est potentiellement sale ou défectueux.

---
*Document ID: SNISID-TELEMETRY-001 | Approuvé par: Directeur des Opérations (FOC)*
