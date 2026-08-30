---
# ============================================================
# SNISID-Field — National Mobile Security Framework
# Hardening Matériel et Sécurité Anti-Vol
# Document ID: SNISID-MOB-SEC-001
# Version: 1.0.0
# ============================================================

## 1. SÉCURITÉ MATÉRIELLE (Device Hardening)

L'environnement physique des tablettes et camions SNISID n'est pas contrôlé (Zone Hostile). Le matériel doit se défendre lui-même.

## 2. TAMPER-RESISTANT DEVICES

- Les tablettes sont scellées à l'époxy. Toute tentative de séparation de l'écran pour accéder aux puces mémoire brise physiquement le PCB (Printed Circuit Board) contenant le module de chiffrement.
- **Bootloader Verrouillé :** Impossible de démarrer sur une clé USB ou de modifier le firmware (Android Verified Boot).

## 3. EMERGENCY LOCKOUT (Verrouillage d'Urgence)

Si le FOC détecte qu'un convoi MGU s'écarte de son itinéraire autorisé (Geofencing), ou si l'escorte policière presse le "Panic Button" :
1. Les tablettes passent instantanément en écran noir (Verrouillage).
2. Elles envoient un "Dying Gasp" (Dernier signal GPS/Batterie).
3. Le module Secure Enclave efface ses clés de chiffrement, rendant le disque illisible de manière permanente.

---
*Document ID: SNISID-MOB-SEC-001 | Approuvé par: CISO National*
