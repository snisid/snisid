---
# ============================================================
# SNISID-Infra — Network Security Foundation
# Zero Trust, Micro-segmentation & Mitigation DDoS
# Document ID: SNISID-NET-SEC-001
# Version: 1.0.0
# ============================================================

## 1. STRATÉGIE DE DÉFENSE EN PROFONDEUR (Defense-in-Depth)

La sécurité du Datacenter ne repose pas sur un seul "mur" (Firewall périmétrique), mais sur de multiples couches de vérification. Si un serveur web est compromis, l'attaquant ne doit pas pouvoir "pivoter" vers la base de données.

## 2. ARCHITECTURE DE SÉCURITÉ RÉSEAU

### 2.1 DDoS Mitigation (BGP Flowspec)
Protection contre les attaques par déni de service distribué (DDoS) visant à saturer la bande passante du gouvernement.
- Netflow analyse le trafic entrant au niveau des routeurs Edge (Border Routers).
- Si un pic suspect est détecté (ex: attaque UDP d'amplification de 100 Gbps), le trafic est "blackholé" (jeté) en amont ou redirigé vers un "Scrubbing Center" (centre de nettoyage) cloud avant d'entrer en Haïti.

### 2.2 Intrusion Detection & Prevention (IDS/IPS)
Tout le trafic entrant dans les zones internes (Zones 2 et 3) est analysé par des moteurs de Deep Packet Inspection (ex: Suricata/Zeek) mis à jour quotidiennement avec les signatures des menaces avancées (APT).

### 2.3 Microsegmentation (Cilium eBPF)
Au sein du cluster Kubernetes, les règles de pare-feu (NetworkPolicies) sont appliquées au niveau de la carte réseau (eBPF) de chaque container.
- Le `pod-frontend` ne peut communiquer qu'avec le `pod-backend` sur le port 8080.
- Toute autre tentative (ex: un ping, un scan SSH) est bloquée silencieusement et loggée dans le SIEM (Security Information and Event Management) du SOC.

### 2.4 Network Access Control (NAC)
Un employé branchant son ordinateur portable sur une prise RJ45 au Ministère n'obtient pas d'IP par défaut. Le switch (802.1x) interroge le serveur NAC (Cisco ISE/PacketFence) qui vérifie le certificat machine avant d'autoriser l'accès au réseau.

---
*Document ID: SNISID-NET-SEC-001 | Approuvé par: CISO National*
