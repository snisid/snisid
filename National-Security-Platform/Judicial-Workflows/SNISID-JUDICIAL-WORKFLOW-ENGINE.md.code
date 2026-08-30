---
# ============================================================
# SNISID-Security — National Judicial Workflow Engine
# Orchestration BPMN pour la Justice
# Document ID: SNISID-JUD-WKF-001
# Version: 1.0.0
# ============================================================

## 1. ORCHESTRATEUR JUDICIAIRE (Camunda / Temporal)

Le système judiciaire haïtien numérisé s'appuie sur un moteur de workflow robuste pour éviter les détentions préventives prolongées illégales (un problème majeur en Haïti) et garantir le respect des délais légaux (SLA).

## 2. EXEMPLE : WORKFLOW DE GARDE À VUE & INCULPATION

Ce workflow garantit que nul n'est détenu sans présentation devant un juge dans les 48 heures (Constitution haïtienne).

```mermaid
bpmn
    title Garde à vue & Comparution
    
    actor Agent_PNH
    actor Commissaire_Gouvernement
    actor Juge_Instruction
    
    Agent_PNH ->> Workflow_Engine: Enregistrer Arrestation (T0)
    
    Workflow_Engine ->> Workflow_Engine: Start 48h SLA Timer
    
    Workflow_Engine ->> Commissaire_Gouvernement: Notifier Dossier
    
    Commissaire_Gouvernement ->> Commissaire_Gouvernement: Décision (Audition)
    
    alt Libération
        Commissaire_Gouvernement ->> Workflow_Engine: Ordre de libération
        Workflow_Engine ->> Agent_PNH: Exécuter libération
    else Inculpation (Renvoi)
        Commissaire_Gouvernement ->> Workflow_Engine: Requête d'instruction
        Workflow_Engine ->> Juge_Instruction: Assigner Dossier
        Juge_Instruction ->> Workflow_Engine: Ordonnance de dépôt
        Workflow_Engine ->> Service_Penitentiaire: Ordre d'incarcération
    end
    
    Workflow_Engine ->> Workflow_Engine: Cancel 48h SLA Timer
```

## 3. GARANTIES LÉGALES ET SLA

| Étape Judiciaire | Limite Légale (SLA) | Action si Dépassement |
|------------------|---------------------|-----------------------|
| Garde à vue PNH | 48 heures | Alerte Juge de Paix + Verrouillage de prolongation |
| Ordonnance de clôture | 2 à 3 mois | Alerte Doyen du Tribunal + Ministère de la Justice |
| Détention Préventive | Max = Peine encourue | Alerte OPC (Office Protection Citoyen) |

Chaque dépassement de SLA est un événement Kafka inaltérable, permettant au Conseil Supérieur du Pouvoir Judiciaire (CSPJ) d'auditer la performance des tribunaux.

## 4. SIGNATURE ÉLECTRONIQUE (PKI)

Chaque décision judiciaire (Mandat, Ordonnance, Jugement) requiert la signature électronique (PKI) du magistrat. 
- Les magistrats utilisent une carte à puce (SmartCard) ou un token USB cryptographique pour signer les payloads JSON du workflow.
- Le moteur vérifie la chaîne de confiance (Root CA SNISID) avant d'accepter l'étape du workflow.

---
*Document ID: SNISID-JUD-WKF-001 | Approuvé par: Ministère de la Justice*
