---
# ============================================================
# SNISID Capstone — Government Workflow Factory (Phase 11)
# Usine Low-Code/No-Code pour l'État
# Document ID: SNISID-CAP-WORKFLOW-001
# Version: 1.0.0
# ============================================================

## 1. DIGITALISATION ACCÉLÉRÉE

Il est impossible de recruter assez de développeurs (Phase 16) pour coder à la main les milliers de procédures administratives existantes (Demande de passeport, Permis de construire, Certificat de mariage). 
La "Workflow Factory" est une plateforme Low-Code/No-Code hébergée sur Kubernetes (Phase 1).

## 2. MOTEUR D'ORCHESTRATION (Camunda Zeebe)

- Les chefs de service (non-développeurs) utilisent une interface graphique BPMN (Business Process Model and Notation) pour dessiner leurs processus administratifs ("Glisser-Déposer").
- Exemple : Une demande de permis de construire nécessite l'approbation de l'urbanisme, puis le paiement d'une taxe. Le workflow dessiné est automatiquement transformé en code exécutable.

## 3. INTÉGRATION SNISID (Single Sign-On)

La Workflow Factory est nativement connectée à :
1. **L'Identité (Phase 2) :** Les formulaires générés sont pré-remplis avec les données du citoyen connecté.
2. **L'API Gateway (Phase 4) :** Si un workflow demande la vérification du casier judiciaire, le composant graphique se connecte automatiquement au système de la Police (Phase 3).

---
*Document ID: SNISID-CAP-WORKFLOW-001 | Approuvé par: Directeur de la Modernisation de l'État*
