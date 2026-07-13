---
# ============================================================
# SNISID-Security — Forensic & Audit Platform
# Traçabilité immuable & Investigation Replay
# Document ID: SNISID-FORENSIC-001
# Version: 1.0.0
# ============================================================

## 1. INVESTIGATION REPLAY

Puisque le SNISID repose sur l'Event Sourcing, le système d'Audit (Forensic) permet un "Investigation Replay".
Il est possible de "rejouer" la base de données criminelle pour voir exactement quel était l'état d'un dossier à une date précise passée.

**Cas d'usage :**
Un juge est accusé d'avoir ignoré une preuve. L'inspecteur (CSPJ) utilise le Replay pour voir l'état exact de l'écran du juge le `12 Mai 2026 à 14h00`. La preuve y était-elle attachée à ce moment-là ?

## 2. IMMUTABLE LOGS (Kafka Audit Topic)

Tous les événements du système sont dupliqués vers un topic Kafka spécial : `snisid.audit.system`.
- **Rétention :** Infinie.
- **Accès :** Ségrégation stricte (Seulement le CSPJ, l'AND, et la Cour des Comptes).

## 3. TAMPER DETECTION (Détection d'Altération)

Un job récurrent (CronJob Kubernetes) lit l'Event Store et recalcule la chaîne des signatures cryptographiques (Merkle Tree). Si un administrateur base de données malveillant (Insider Threat) altère une ligne directement dans CockroachDB, la vérification du hash échouera, déclenchant une alerte nationale (Severity: CRITICAL).

---
*Document ID: SNISID-FORENSIC-001 | Approuvé par: Cour Supérieure des Comptes (CSCCA)*
