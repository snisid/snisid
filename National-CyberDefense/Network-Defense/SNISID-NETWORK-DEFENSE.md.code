---
# ============================================================
# SNISID-Cyber — National Network Defense Platform
# Inspection Profonde, IDS/IPS et DNS Security
# Document ID: SNISID-NETDEF-001
# Version: 1.0.0
# ============================================================

## 1. DÉFENSE EN PROFONDEUR DU RÉSEAU

La sécurité réseau (Phase 5) est ici enrichie par l'analyse comportementale de la Phase 6.

## 2. SURVEILLANCE EST-OUEST (East-West Traffic)

La majorité des pare-feux (Nord-Sud) surveillent le trafic entre Internet et le Datacenter. Cependant, si un serveur interne est compromis (ex: un serveur de messagerie), l'attaquant essaiera de se déplacer latéralement ("Lateral Movement") vers la base de données.
Le trafic **Est-Ouest** (entre deux serveurs internes) est surveillé par des sondes Zeek/Suricata virtuelles (vIDS) déployées sur chaque hyperviseur OpenStack.

## 3. SÉCURITÉ DNS (DNS Sinkholing)

Le DNS est souvent utilisé par les malwares pour contacter leurs serveurs de contrôle (C2 - Command & Control) ou exfiltrer des données.
- Tous les serveurs et postes gouvernementaux utilisent les résolveurs DNS souverains du SNISID.
- **DNS Sinkholing :** Si un EDR ou la Threat Intelligence ajoute le domaine `bad-hacker-c2.com` à la liste noire nationale, le DNS souverain renverra une fausse adresse IP (0.0.0.0) à tout poste essayant de s'y connecter, bloquant instantanément le malware.

---
*Document ID: SNISID-NETDEF-001 | Approuvé par: Head of Network Security*
