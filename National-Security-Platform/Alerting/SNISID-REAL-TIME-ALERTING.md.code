---
# ============================================================
# SNISID-Security — National Real-Time Alerting Platform
# Routage d'Incidents et Escalade
# Document ID: SNISID-ALERTING-001
# Version: 1.0.0
# ============================================================

## 1. CONCEPT DE L'ALERTING

L'architecture d'alerte SNISID gère la diffusion de notifications push critiques vers les agents (via PNH Mobile) et les centres de commandement (via Web).

## 2. NIVEAUX DE SÉVÉRITÉ

| Niveau | Définition | Action Requise | Délai Escalade |
|--------|------------|----------------|----------------|
| **INFO** | Information standard (Rapport soumis) | Aucune | N/A |
| **WARN** | Attention (Ex: 40h de Garde à vue sur 48h) | Proactive | 8 heures |
| **HIGH** | Événement grave (Ex: Arrestation Fugitif) | Réactive (Intervention) | 15 minutes |
| **CRIT** | Sécurité Nationale (Ex: Évasion, Altération Preuve) | Immédiate | 5 minutes |

## 3. INCIDENT ROUTING MODEL (Kafka -> PagerDuty souverain)

Les événements transitent par Kafka. Le service `alert-manager` évalue les règles.

Exemple : ÉVASION DE PRISON.
1. `EscapeEvent` généré. (Severity: CRIT).
2. `alert-manager` notifie le Directeur de la Prison (SMS + Call automatique).
3. Si non acquitté en 5 mins, escalade à la DG PNH et Ministre de la Justice.
4. Broadcast simultané à toutes les unités PNH mobiles dans un rayon de 50km de la prison (Push notification sur kit MEK).

---
*Document ID: SNISID-ALERTING-001 | Approuvé par: Centre Renseignement Opérationnel (CRO)*
