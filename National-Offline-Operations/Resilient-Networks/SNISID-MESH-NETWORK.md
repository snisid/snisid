---
# ============================================================
# SNISID-Edge — National Resilient Field Network
# Réseaux Maillés et Communications Radio Tactiques
# Document ID: SNISID-MESH-001
# Version: 1.0.0
# ============================================================

## 1. CONNECTIVITÉ INTERMITTENTE

Dans les zones reculées d'Haïti, le réseau 4G est inexistant ou instable. L'architecture réseau terrain (Field Network) du SNISID est conçue pour fonctionner avec des réseaux non fiables.

## 2. MESH NETWORKING (Réseaux Maillés Locaux)

Si un poste frontière est coupé d'Internet, les terminaux mobiles des agents (tablettes) peuvent tout de même communiquer entre eux.
- Utilisation de protocoles Mesh (ex: Bluetooth LE / Wi-Fi Aware / B.A.T.M.A.N. Advanced).
- Si l'Agent A n'a pas de réseau, mais que l'Agent B réussit à attraper un signal 3G en haut d'une colline, la tablette de l'Agent A routéra ses paquets chiffrés à travers la tablette de l'Agent B pour atteindre le serveur central.

## 3. RADIO TACTIQUE (Fallback Ultime)

En cas de black-out total (ouragan majeur), les réseaux Mesh locaux ne peuvent plus joindre la capitale.
- Les Noeuds Edge majeurs sont équipés de modems radio tactiques (ondes courtes / LoRaWAN).
- Bien que le débit soit extrêmement faible (quelques kilooctets par seconde), il permet de transmettre des "Pings de Survie" (Health Checks) et de petits messages textes chiffrés d'urgence au Datacenter Central.

---
*Document ID: SNISID-MESH-001 | Approuvé par: Architecte Réseau (AND)*
