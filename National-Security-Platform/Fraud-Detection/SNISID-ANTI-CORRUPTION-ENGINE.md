---
# ============================================================
# SNISID-Security — National Fraud & Corruption Detection Engine
# Moteur Anti-Corruption & Insider Threat
# Document ID: SNISID-ANTI-CORRUPTION-001
# Version: 1.0.0
# ============================================================

## 1. LUTTE CONTRE LA CORRUPTION (INSIDER THREAT)

Un des objectifs clés du SNISID est de combattre la corruption endémique. Le moteur d'Anti-Corruption analyse le comportement des utilisateurs autorisés (juges, policiers, greffiers) pour détecter des anomalies.

## 2. SCÉNARIOS DE DÉTECTION (ANOMALY ANALYTICS)

### 2.1 "The Missing File" (Le dossier disparu)
- **Scénario :** Un greffier tente d'effacer une preuve d'un dossier.
- **Détection :** Impossible techniquement (Event Sourcing). Si tentative d'injection SQL directe -> `MerkleTreeVerificationFailed` -> Alerte CRIT au CSPJ.

### 2.2 "The VIP Treatment" (Traitement de faveur)
- **Scénario :** Un juge libère un suspect classé "GANG_MEMBER" de manière anormalement rapide (< 24h) sans ordonnance détaillée.
- **Détection :** Le workflow judiciaire croise le label DCPJ avec le délai de décision. Une alerte "Suspicious Judicial Action" est générée pour audit.

### 2.3 "The Snooper" (Le fouineur)
- **Scénario :** Un agent de police recherche le même NIU (ex: une célébrité ou un politicien) 50 fois par jour sans rapport d'incident lié.
- **Détection :** Kafka Streams analyse les logs de recherche. `SearchVelocity > Threshold && NoLinkedIncident == true` -> Alerte DG PNH.

## 3. TABLEAU DE BORD "INSIDER THREAT"

Seul le Conseil Supérieur (CSPJ/CSPN) a accès à ce tableau de bord qui liste le "Risk Score" de chaque agent de l'État basé sur son activité numérique.

---
*Document ID: SNISID-ANTI-CORRUPTION-001 | Approuvé par: ULCC (Unité de Lutte Contre la Corruption)*
