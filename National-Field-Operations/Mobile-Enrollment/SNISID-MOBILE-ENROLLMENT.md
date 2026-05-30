---
# ============================================================
# SNISID-Field — National Mobile Enrollment Platform
# Capture Biométrique Hors-Ligne & Anti-Fraude
# Document ID: SNISID-MOB-ENROL-001
# Version: 1.0.0
# ============================================================

## 1. STRATÉGIE D'ENRÔLEMENT MOBILE

Afin d'atteindre 100% de la population haïtienne, l'Office National d'Identification (ONI) utilise des "Kits Valises" (Suitcase Kits) transportables à dos de mule ou en moto pour les sections communales inaccessibles en voiture.

## 2. ARCHITECTURE DU KIT D'ENRÔLEMENT

Un kit mobile comprend :
- 1 Tablette Android durcie (Rugged, MIL-STD-810G).
- 1 Scanner d'empreintes (FAP 50 / 10 doigts).
- 1 Appareil photo / Scanner d'iris (ISO 19794).
- 1 Imprimante thermique de reçus.
- 1 Batterie externe (Power Bank) + Panneau solaire pliable 100W.

## 3. CONTRÔLES ANTI-FRAUDE TERRAIN

L'agent enrôleur est potentiellement isolé. L'application intègre des contrôles stricts pour empêcher la création de faux citoyens :
1. **Vérification Liveness (Preuve de vie) :** Le logiciel bloque la capture si l'on présente une photo imprimée au lieu d'un visage réel (Anti-Spoofing de Phase 2).
2. **Double Signature :** Chaque enregistrement doit être signé par le certificat (Smartcard) de l'agent ONI **et** validé par les empreintes d'un "Témoin Notable" (ex: Maire, Pasteur, Juge de paix) préalablement certifié dans le système.
3. **Géolocalisation Immuable :** Les coordonnées GPS sont hachées avec les données biométriques. Si un agent est assigné à Jérémie, mais que le GPS indique Port-au-Prince, l'enrôlement sera automatiquement flaggé par le SOC lors de la synchronisation.

---
*Document ID: SNISID-MOB-ENROL-001 | Approuvé par: Directeur Général de l'ONI*
