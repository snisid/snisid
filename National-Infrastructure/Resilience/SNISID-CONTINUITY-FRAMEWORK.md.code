---
# ============================================================
# SNISID-Infra — National Resilience & Continuity Framework
# Plan de Continuité (PCA) et de Reprise (PRA)
# Document ID: SNISID-CONTINUITY-001
# Version: 1.0.0
# ============================================================

## 1. GESTION DES CATASTROPHES (Disaster Recovery Playbooks)

Le système SNISID doit survivre aux scénarios catastrophes majeurs de la région des Caraïbes.
Ce framework définit les processus (Playbooks) exécutés par le NOC/SOC.

## 2. SCÉNARIOS DE CRISE (Crisis Scenarios)

### 2.1 Ouragan Catégorie 5 (Hurricane)
- **Préparation (H-48) :** Ravitaillement prioritaire en diesel des Datacenters. Basculement (Scale-up) automatique des bases de données sur le Datacenter DR si le Primary est sur la trajectoire directe.
- **Impact :** Perte probable des antennes VSAT et fibres optiques aériennes.
- **Continuité :** Les Noeuds Edge passent en mode "Offline-First" (autonomie complète sur énergie solaire + batteries locales).

### 2.2 Séisme Majeur (Earthquake)
- **Impact :** Destruction physique possible du bâtiment Primary Datacenter et coupures massives de fibre optique sous-marine.
- **Continuité :** Basculement automatique du BGP (Failover) vers le Datacenter DR en moins d'une heure. Le RPO (Recovery Point Objective) garantit moins de 5 minutes de perte de données.

### 2.3 Crise Socio-Politique (Instabilité / Émeutes)
- **Impact :** Impossibilité pour le personnel de se rendre au Datacenter ou de ravitailler en carburant. Sabotage des câbles.
- **Continuité :** Le NOC (National Operations Center) active le pilotage 100% à distance via VPN satellitaire d'urgence. Le Datacenter est verrouillé (Lockdown physique). 

## 3. EXERCICES DE CHAOS (Chaos Engineering)

L'équipe "Chaos" de l'infrastructure coupe volontairement un câble d'alimentation (PDU) ou un lien fibre optique tous les mois en pleine journée pour valider que le système de redondance fonctionne parfaitement (cf. *Chaos Monkey*).

---
*Document ID: SNISID-CONTINUITY-001 | Approuvé par: Premier Ministre*
