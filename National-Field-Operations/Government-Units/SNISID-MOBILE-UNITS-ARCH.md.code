---
# ============================================================
# SNISID-Field — Mobile Government Units Architecture
# Flotte de Camions et Autonomie Énergétique
# Document ID: SNISID-MOB-UNITS-001
# Version: 1.0.0
# ============================================================

## 1. LES "TRUCKS" D'ENRÔLEMENT (MGU - Mobile Government Units)

Pour les villes de province, le programme déploie des véhicules tout-terrain (4x4 lourds ou minibus modifiés) transformés en bureaux d'État Civil mobiles.

## 2. ARCHITECTURE D'UN CAMION SNISID

Chaque camion est un **Edge Node complet** sur roues.

### 2.1 Équipement Informatique
- Mini-serveur rackable (K3s Edge Node - Phase 7).
- Routeur multi-WAN : Starlink (Primaire) + VSAT Ku-Band (Secours) + 4G/LTE (Urbain).
- Points d'accès Wi-Fi Mesh pour connecter jusqu'à 10 tablettes d'agents autour du camion.

### 2.2 Autonomie Énergétique (Power Architecture)
Le diesel est souvent rare ou volé. L'architecture priorise le renouvelable.
- **Toit Solaire :** 4 panneaux solaires plats (Total: 1.2 kW).
- **Parc Batteries :** Batteries Lithium Fer Phosphate (LiFePO4) de 10 kWh. Assure le fonctionnement de tout l'équipement informatique et de la climatisation légère pendant 48 heures sans soleil.
- **Alternateur Secondaire :** Le moteur du camion recharge les batteries en roulant.
- **Générateur Thermique Inverter :** Réservé exclusivement à l'extrême urgence (Saison des pluies prolongée).

### 2.3 Sécurité Physique (Panic Button)
En cas d'attaque (Gang / Émeute) :
- Le chauffeur dispose d'un bouton rouge "ZEROIZE" caché.
- L'appui sur ce bouton déclenche le module TPM du serveur local, qui efface la clé de chiffrement du disque dur en moins de 1 seconde. Le camion peut être volé, mais les données biométriques des citoyens sont physiquement irrécupérables.

---
*Document ID: SNISID-MOB-UNITS-001 | Approuvé par: Directeur de la Logistique Nationale*
