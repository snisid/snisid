---
# ============================================================
# SNISID-Interop — Agency Integration Map
# Matrice d'Intégration Inter-Gouvernementale
# Document ID: SNISID-INT-MAP-001
# Version: 1.0.0
# ============================================================

## 1. CARTOGRAPHIE DES FLUX NATIONAUX

Ce diagramme représente les flux vitaux d'intégration entre les systèmes SNISID (Core/Identité) et les autres instances gouvernementales.

```mermaid
graph TD
    subgraph "SNISID Sovereign Core"
        ID[Identity Registry\n(Master Data)]
        BUS[Kafka Interop Bus]
        GATE[National API Gateway]
    end

    subgraph "Agences Consommatrices & Productrices"
        DGI[DGI - Impôts\n(Tax & Revenue)]
        MSPP[MSPP - Santé\n(Hospitals & Births)]
        JUSTICE[Ministère Justice\n(Courts & Police)]
        CEP[CEP - Élections\n(Voter Rolls)]
        MENFP[MENFP - Éducation\n(Student Tracking)]
        BANQUES[Secteur Bancaire\n(KYC Verification)]
    end

    MSPP -- "Publie Naissance/Décès" --> GATE
    JUSTICE -- "Publie Condamnation" --> GATE
    
    GATE --> BUS
    BUS --> ID
    
    ID -- "Publie Identity Update" --> BUS
    
    BUS -- "Consomme Changements" --> DGI
    BUS -- "Consomme Changements" --> CEP
    
    BANQUES -- "Vérification KYC (API Monétisée)" --> GATE
    MENFP -- "Vérification NIU Étudiant" --> GATE
```

## 2. MODÈLE DE GOUVERNANCE DES DONNÉES (MASTER DATA)

Pour éviter la duplication et les incohérences, le concept de **"Single Source of Truth" (SSOT)** est imposé légalement.

| Domaine de Donnée | Propriétaire (Owner) | Règle d'Interopérabilité |
|-------------------|----------------------|--------------------------|
| **Identité Civile** | ONI / SNISID | Seul le SNISID peut modifier un NIU, nom, ou date de naissance. La DGI et le CEP doivent consommer l'état civil du SNISID. |
| **Statut Médical/Décès** | MSPP | Seul un médecin certifié MSPP peut émettre un certificat de décès qui modifiera le statut de l'identité dans le SNISID. |
| **Statut Criminel** | Justice / PNH | Seule la Justice peut altérer la liberté d'un individu. |
| **Passeport/Visas** | DOI (Immigration) | L'émission d'un passeport doit requérir un appel temps réel à l'API SNISID pour vérifier l'identité et les mandats (Warrants). |

---
*Document ID: SNISID-INT-MAP-001 | Classification: OFFICIELLE*
