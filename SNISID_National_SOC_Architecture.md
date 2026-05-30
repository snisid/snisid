# SNISID : COMMANDEMENT NATIONAL DE CYBERDÉFENSE (SOC)
## Index Principal et Architecture des Opérations de Sécurité
**République d'Haïti — Infrastructure Souveraine de Cybersécurité**

---

## PRÉAMBULE
Ce document constitue l'index central de l'architecture du **Commandement National de Cyberdéfense** du SNISID. Ce SOC (Security Operations Center) de niveau étatique (Government-Grade) assure la protection 24/7 de l'identité numérique, de l'état civil, et des frontières d'Haïti contre les menaces avancées (APT), les cyberattaques d'États-nations et le cyberterrorisme.

L'architecture est fondée sur les principes du **Zero Trust Architecture (ZTA)**, intégrant des capacités souveraines de SIEM, SOAR, UEBA, et de réponse aux incidents (DFIR).

---

## 📑 SOMMAIRE DE LA CYBERDÉFENSE NATIONALE

La doctrine défensive et opérationnelle est découpée en 4 volumes d'ingénierie tactique :

### 🛡️ [VOLUME 1 : Structure du Commandement et du SOC](cyber-defense/SNISID_SOC_01_Command_Structure.md)
Organisation hiérarchique et unités tactiques :
- SOC Tiers 1 (Triage), Tier 2 (Analyse), Tier 3 (Threat Hunting).
- Unités DFIR (Digital Forensics & Incident Response) et opérations CERT.
- La "War Room" Nationale et la gouvernance de crise cyber.

### 👁️ [VOLUME 2 : Opérations de Sécurité et Télémétrie](cyber-defense/SNISID_SOC_02_Security_Operations.md)
Architecture de détection des menaces :
- Architecture SIEM (Ingestion massive Kubernetes, X-Road, Kafka).
- UEBA (Détection des menaces internes et anomalies comportementales).
- Intégration de la Threat Intelligence et mapping MITRE ATT&CK.
- Sécurité des APIs et surveillance des flux est-ouest (Service Mesh).

### ⚡ [VOLUME 3 : SOAR et Réponse aux Incidents (IR)](cyber-defense/SNISID_SOC_03_SOAR_and_Response.md)
Orchestration et escalade défensive :
- Playbooks d'automatisation (SOAR).
- Workflows d'incidents et matrice d'escalade cyber.
- Pipelines forensiques, préservation des preuves légales (Chain of Custody).

### 🚨 [VOLUME 4 : Playbooks de Crise et Cyber-Guerre](cyber-defense/SNISID_SOC_04_Crisis_Playbooks.md)
Scénarios d'urgence extrême :
- Réponse aux Ransomwares d'États (Confinement et Isolation).
- Procédures de "Disaster Cyber Response" (Zero Trust Lockdown).
- Protocoles de rétablissement post-violation majeure.

---

## 🔐 DOCTRINE "DEFEND FORWARD & ISOLATE"
La survie de l'État Civil haïtien dépend de la capacité du SOC à confiner instantanément toute menace avant qu'elle ne compromette l'intégrité de la base de données (WORM) ou l'autorité de certification (AN-PKI). Le Commandement de Cyberdéfense dispose des pleins pouvoirs techniques pour isoler n'importe quel ministère, agence ou département géographique du reste du réseau national en moins de 5 secondes.

*Spécifié, modélisé et ratifié par l'Agence Nationale de Cybersécurité (ANCD).*
*Classification : TRÈS SECRET / SOUVERAIN / DÉFENSE NATIONALE*
