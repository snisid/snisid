---
# ============================================================
# SNISID-Cyber — National Forensics Platform
# Collecte de Preuves et Chaîne de Traçabilité
# Document ID: SNISID-FORENSIC-001
# Version: 1.0.0
# ============================================================

## 1. FORENSIC READINESS (Prêt pour la Justice)

Lorsqu'un piratage majeur a lieu (ex: tentative de falsification des résultats électoraux via l'état civil), le SOC ne se contente pas de réparer : il collecte des preuves numériques admissibles devant un tribunal.

## 2. CHAIN OF CUSTODY (Chaîne de Possession)

Toute preuve numérique (ex: un dump de la RAM d'un serveur infecté, un fichier log) doit être traitée selon un protocole strict :
1. **Acquisition :** L'image disque est capturée bit-à-bit.
2. **Hashing (Intégrité) :** L'empreinte cryptographique (SHA-256) du fichier est calculée immédiatement. Cette empreinte est insérée dans la base de données immuable (CockroachDB/Event Store) pour garantir qu'elle n'a jamais été altérée.
3. **Stockage WORM :** La preuve est stockée sur les serveurs MinIO en mode Write-Once-Read-Many (cf. Phase 5). Elle ne peut être ni modifiée ni supprimée, même par l'administrateur système.

## 3. WORKFLOW D'INVESTIGATION (Incident Replay)

Les experts de la DCPJ (Police Judiciaire) utilisent une "Clean Room" (réseau totalement isolé) pour analyser les malwares.
Le SIEM permet de faire un "Incident Replay", c'est-à-dire de rejouer les logs seconde par seconde pour comprendre exactement comment l'attaquant a contourné l'API Gateway, puis pivoté vers la base de données.

---
*Document ID: SNISID-FORENSIC-001 | Approuvé par: Chef de l'Unité Cybercrime (DCPJ)*
