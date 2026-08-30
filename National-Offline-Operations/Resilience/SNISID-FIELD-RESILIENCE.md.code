---
# ============================================================
# SNISID-Edge — National Field Resilience Framework
# Continuité et Survie en Isolement Régional
# Document ID: SNISID-FIELD-RESIL-001
# Version: 1.0.0
# ============================================================

## 1. SURVIVRE À L'ISOLEMENT (Regional Isolation Survival)

L'expérience du séisme de 2010 a prouvé que les départements peuvent être physiquement et numériquement coupés de Port-au-Prince pendant des semaines. La résilience terrain est la capacité d'un Edge Node à survivre sans maintenance.

## 2. PLAYBOOKS DE RÉSILIENCE TERRAIN

### 2.1 Ouragan Catégorie 5 (Destruction VSAT)
- L'Edge Node détecte la chute barométrique (via capteurs IoT) et la perte de signal.
- **Action :** Verrouillage des disques en lecture seule (sauf pour le spooler NATS) pour protéger l'intégrité de la base. Activation des protocoles d'économie d'énergie pour maximiser la durée de vie sur batteries solaires.

### 2.2 Corruption de Données Offline (Offline Corruption Recovery)
Si le noeud local subit un crash disque et corrompt la base Edge DB :
- Les terminaux mobiles (Tablettes) basculent en "Peer-to-Peer" (Bluetooth/Wi-Fi Direct) pour continuer les opérations minimales.
- Le noeud Edge demande une "Full Sync" (Image complète) plutôt qu'un "Delta" lors de sa reconnexion au réseau national.

---
*Document ID: SNISID-FIELD-RESIL-001 | Approuvé par: Architecte Souverain*
