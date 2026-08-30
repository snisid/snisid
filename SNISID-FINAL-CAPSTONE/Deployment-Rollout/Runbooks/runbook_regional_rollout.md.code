# SNISID Runbook — Regional Rollout Procedure
**Code de Procédure :** SNISID-RB-01  
**Statut :** Approuvé & Obligatoire  
**Audience :** Coordonnateurs Départementaux et Ingénieurs Logistiques  

---

## 1. Objectif

Ce runbook décrit la procédure pas-à-pas pour préparer, exécuter et valider l'activation physique et technique d'un nouveau département géographique de la République d'Haïti dans l'écosystème SNISID.

---

## 2. Déroulement Chronologique (Séquence de 6 Semaines)

```
S-6 to S-4: Audit & Logistics ----> S-3 to S-2: Installation ----> S-1: Training ----> S-0: Go-Live
```

### PHASE 1 : Préparation & Audit de Readiness (Semaine -6 à Semaine -4)
1. **Évaluation Logistique :** Remplir et valider la grille d'évaluation de maturité régionale (*Regional Readiness Scorecard*). Le score doit atteindre **au moins 85%**.
2. **Énergie de Secours :** Livrer sur chaque BLC cible les panneaux solaires et le boitier d'autonomie batterie.
3. **Sécurisation :** Prendre contact avec le Directeur Départemental de la PNH pour convenir du protocole d'escorte du matériel sensible (Edge Nodes, terminaux biométriques) et du gardiennage des locaux.

### PHASE 2 : Installation Technique & Connectivité (Semaine -3 à Semaine -2)
1. **Assemblage Matériel :** Poser le kit Starlink Business, câbler la liaison Ethernet locale vers la station de travail de l'opérateur et alimenter le Local Edge Node (LEN).
2. **Initialisation Logique :** Démarrer le LEN. Introduire la clé de licence d'activation départementale signée par le Central.
3. **Test Réseau :** Exécuter le script de validation de latence et de bande passante :
   `ping -c 100 central-dc.snisid.gouv.ht`
   *La latence moyenne doit être inférieure à 120ms, et aucune perte de paquets ne doit être tolérée (0.00% packet loss).*

### PHASE 3 : Formation & Habilitation des Agents (Semaine -1)
1. **Délivrance de la Formation :** Dispenser le programme de formation standardisé de 5 jours aux opérateurs locaux.
2. **Examen Pratique :** Évaluer chaque candidat opérateur. L'opérateur doit réussir un enrôlement complet de test (Saisie démographique + Capture 10 empreintes + Capture visage et iris) en moins de 3 minutes.
3. **Génération d'accès YubiKey :** Émettre pour chaque opérateur certifié sa clé cryptographique matérielle personnelle d'accès et enregistrer le certificat client dans le LEN.

### PHASE 4 : Go-Live & Bascule Officielle (Semaine 0)
1. **Bascule d'Autorité :** Exécuter le technical runbook de cutover (DNS, redirection API).
2. **Ouverture Civile :** Ouvrir les portes du BLC à la population. Lancer la campagne de communication SMS et radio locale.
3. **Surveillance Hypercare :** Placer le département sous surveillance renforcée pendant 30 jours (activation du support Tier 1 de proximité).
