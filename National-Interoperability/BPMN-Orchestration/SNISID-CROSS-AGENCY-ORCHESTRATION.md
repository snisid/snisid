---
# ============================================================
# SNISID-Interop — National Cross-Agency Orchestration Engine
# BPMN Multi-Ministères (Camunda/Temporal)
# Document ID: SNISID-BPMN-ORCH-001
# Version: 1.0.0
# ============================================================

## 1. ORCHESTRATION INTER-AGENCES

Certains processus gouvernementaux requièrent l'action de multiples ministères (Silo-breaking). Un orchestrateur central (Camunda) coordonne ces flux complexes.

## 2. EXEMPLE : WORKFLOW "CRÉATION D'ENTREPRISE" (GUICHET UNIQUE)

Au lieu de faire physiquement le tour de 4 ministères, le citoyen lance un seul processus.

```mermaid
bpmn
    title Guichet Unique - Création d'Entreprise
    
    actor Citoyen
    actor Orchestrateur (Camunda)
    actor Ministère_Commerce (MCI)
    actor DGI (Impôts)
    actor OFATMA (Sécurité Sociale)
    
    Citoyen ->> Orchestrateur: Formulaire Soumis
    Orchestrateur ->> MCI: Validation Nom/Statuts (API)
    
    alt Nom Invalide
        MCI ->> Orchestrateur: Rejet
        Orchestrateur ->> Citoyen: Notification Rejet
    else Nom Valide
        MCI ->> Orchestrateur: Approuvé (Patente Provisoire)
        
        par Action Parallèle
            Orchestrateur ->> DGI: Demande NIF (Numéro Fiscal)
            Orchestrateur ->> OFATMA: Inscription Employeur
        end
        
        DGI ->> Orchestrateur: NIF Généré
        OFATMA ->> Orchestrateur: Matricule Généré
        
        Orchestrateur ->> Citoyen: Dossier Complet (Patente + NIF)
    end
```

## 3. SLA ET ESCALATION CHAINS

L'Orchestrateur surveille les délais d'exécution de chaque Ministère.
- Si le MCI prend plus de 7 jours pour valider un nom, un événement d'escalade est envoyé au Directeur du MCI.
- Toutes les étapes et latences sont loggées pour générer le "Bulletin de Performance" du gouvernement (Dashboard exécutif).

---
*Document ID: SNISID-BPMN-ORCH-001 | Approuvé par: Ministère de l'Économie*
