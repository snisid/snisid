---
# ============================================================
# SNISID-Edge — National Mobile Operations Platform
# Applications Terrain pour Forces de l'Ordre et État Civil
# Document ID: SNISID-MOBILE-001
# Version: 1.0.0
# ============================================================

## 1. APPLICATIONS TERRAIN (Field Apps)

Les agents de l'État utilisent des tablettes Android durcies (Rugged Tablets) pour leurs opérations quotidiennes.
Ces applications sont construites avec des technologies Offline-First (ex: React Native + WatermelonDB / SQLite local).

## 2. CAS D'USAGE

### 2.1 Enrôlement Biométrique Mobile (Campagnes Rurales)
L'Office National d'Identification (ONI) déploie des équipes dans des villages isolés (sans aucun réseau).
- L'agent capture les empreintes, la photo et les données démographiques sur la tablette.
- L'application chiffre ces données (Public Key Cryptography) et les stocke localement.
- De retour au bureau régional (Wi-Fi), la tablette transmet le "Batch" d'enrôlements vers le serveur central pour déduplication ABIS.

### 2.2 Contrôle Routier Police (PNH)
Un policier arrête un véhicule dans une zone blanche (sans 4G).
- Il scanne le QR code (vCVD) du permis de conduire.
- L'application vérifie cryptographiquement la signature du QR code avec la clé publique stockée en cache.
- L'identité est formellement validée sans aucune connexion réseau.

---
*Document ID: SNISID-MOBILE-001 | Approuvé par: Directeur des Systèmes d'Information (PNH)*
