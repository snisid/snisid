---
# ============================================================
# SNISID-Cyber — National Threat Intelligence Platform
# Anticipation, Profilage et IOCs
# Document ID: SNISID-CTI-001
# Version: 1.0.0
# ============================================================

## 1. THREAT INTELLIGENCE (Le Renseignement Cyber)

La défense passive ne suffit pas. Le SOC doit connaître les tactiques, techniques et procédures (TTP) des groupes d'attaquants (ex: APT29, Lazarus) avant même qu'ils ne frappent.

## 2. SOURCES DE RENSEIGNEMENT (Feeds)

La plateforme centralise les flux (via MISP - Malware Information Sharing Platform) :
- **OSINT :** Flux publics (Abuse.ch, AlienVault OTX).
- **Gouvernemental :** Flux partagés par des CERTs alliés (ex: CISA, ANSSI française).
- **Dark Web :** Monitoring automatique des forums clandestins pour détecter des fuites de mots de passe d'agents de l'État (Credential Stuffing).

## 3. CORRÉLATION AUTOMATIQUE (Matching)

Si un flux Threat Intel signale qu'un nouveau ransomware utilise l'IP `198.51.100.45` pour contacter son serveur de contrôle :
1. L'IP est automatiquement ingérée dans MISP.
2. Le SIEM (OpenSearch) scanne instantanément les **3 derniers mois** de logs (Retro-Hunting) pour vérifier si une machine de l'État n'a pas déjà contacté cette IP.
3. Les Firewalls (Phase 5) bloquent automatiquement toute nouvelle connexion sortante vers cette IP.

---
*Document ID: SNISID-CTI-001 | Approuvé par: Head of Cyber Intelligence*
