---
# ============================================================
# SNISID-Edge — National Offline BPMN Engine
# Exécution Locale des Workflows (Camunda Edge)
# Document ID: SNISID-BPMN-EDGE-001
# Version: 1.0.0
# ============================================================

## 1. ORCHESTRATION DÉCONNECTÉE

Le moteur central BPMN (Phase 4) gère les processus inter-ministères. Cependant, les processus locaux (ex: "Procédure d'arrestation en commissariat") doivent pouvoir s'exécuter sans interroger la capitale.

## 2. ARCHITECTURE EDGE BPMN (Camunda Zeebe Edge)

Un micro-moteur d'orchestration (ex: Zeebe) tourne dans le cluster K3s du commissariat.

### 2.1 Sync-Aware BPMN
Les workflows BPMN sont modélisés pour être "conscients" de l'état du réseau (Sync-Aware).

```mermaid
graph TD
    Start[Début: Arrestation Suspect] --> CheckNet{Réseau WAN actif ?}
    
    CheckNet -- OUI --> ApiCentral[Vérification Casier Judiciaire (Central)]
    CheckNet -- NON --> ApiLocal[Vérification Cache Local (Edge DB)]
    
    ApiCentral --> Decision{Recherché ?}
    ApiLocal --> Decision
    
    Decision -- OUI --> Cellule[Placement en Cellule]
    Decision -- NON --> Relache[Relâchement ou Garde à Vue Simple]
    
    Cellule --> Queue[Mise en file d'attente (NATS) pour rapport central]
```

### 2.2 Delayed Escalation
Si un processus requiert l'approbation d'un juge basé à la capitale, mais que le réseau est coupé, le workflow BPMN passe en état `PAUSED_AWAITING_SYNC`. L'agent local voit que la demande est "En attente de transmission". Dès que le réseau remonte, l'événement est émis vers Kafka.

---
*Document ID: SNISID-BPMN-EDGE-001 | Approuvé par: Architecte Processus Métiers*
