---
# ============================================================
# SNISID-Field — National Field Workforce Management
# Suivi des Équipes et Assignation de Missions
# Document ID: SNISID-WORKFORCE-001
# Version: 1.0.0
# ============================================================

## 1. GESTION DES EFFECTIFS TERRAIN

Le FOC (Field Operations Center) supervise des milliers d'agents (Policiers, Officiers d'État Civil, Techniciens IT) déployés simultanément.

## 2. WORKFORCE SCHEDULING (Assignation Intelligente)

Le système de gestion (intégré au BPMN central) assigne les missions en fonction :
- **De la Sécurité :** L'algorithme refuse d'assigner une équipe dans une zone rouge signalée par la Threat Intel (Phase 6).
- **Des Compétences :** Un kit biométrique d'enrôlement nécessite un Agent ONI certifié (Niveau 2) et un superviseur technique.
- **De l'Autonomie :** Le système vérifie l'état de la batterie de la tablette avant d'assigner une mission longue.

## 3. MOBILE WORKFORCE TRACKING

- Les équipes en déplacement (Camions MGU) transmettent leur position GPS toutes les 5 minutes.
- **Geofencing :** Si un camion sort de son périmètre départemental assigné sans autorisation, le FOC reçoit une alerte P1 (Suspicion de vol/Détournement), et peut déclencher la procédure de *Remote Wipe* (Phase 7).

---
*Document ID: SNISID-WORKFORCE-001 | Approuvé par: Directeur des Ressources Humaines (ONI)*
