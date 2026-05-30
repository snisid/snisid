---
# ============================================================
# SNISID-Edge — National Regional Autonomy Model
# Gouvernance Déconnectée et Continuité de Décision
# Document ID: SNISID-REG-AUTO-001
# Version: 1.0.0
# ============================================================

## 1. AUTONOMIE DE DÉCISION (Delegated Governance)

La technologie Offline-First ne sert à rien si les processus bureaucratiques (BPMN) nécessitent systématiquement une signature électronique de la capitale. L'autonomie régionale permet la "Délégation de Pouvoir" numérique.

## 2. WORKFLOWS DÉGRADÉS (Degraded-Mode Operations)

### Scénario : Délivrance d'un acte de naissance tardif en mode Offline.
En temps normal, ce processus (Phase 2) requiert une validation croisée via l'API de la Justice à Port-au-Prince.

**Mode Dégradé (Réseau coupé > 24h) :**
1. L'Officier d'État Civil (OEC) initie le workflow localement.
2. Le système Edge détecte l'absence de réseau. Il autorise l'OEC à valider l'acte en activant une "Signature d'Urgence Locale".
3. Un PDF/A-3 temporaire (valide 30 jours) est généré localement avec un filigrane "VALIDÉ EN MODE OFFLINE".
4. Au retour de la connexion, le Sync Engine transmet l'acte au tribunal central pour validation finale (Post-Audit) et génération de l'acte permanent.

## 3. CONTINUITÉ OPÉRATIONNELLE RÉGIONALE

Chaque région dispose d'un "Responsable de Continuité". Si la région est coupée de l'État central suite à un séisme, ce responsable utilise la clé cryptographique du "Regional Trust Anchor" pour émettre temporairement des certificats d'identité locaux valables uniquement au sein du département, permettant de distribuer l'aide humanitaire sans fraude.

---
*Document ID: SNISID-REG-AUTO-001 | Approuvé par: Ministère de l'Intérieur*
