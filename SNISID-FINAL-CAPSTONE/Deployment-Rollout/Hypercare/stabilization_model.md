# SNISID Post-GoLive Stabilization & Hypercare Model
## Modèle de Stabilisation et de Support Post-Production

---

## 1. Le Cadre Général du Support Hypercare

Le déploiement initial du SNISID dans un nouveau département ou une nouvelle agence étatique s'accompagne d'un afflux important de requêtes d'assistance, de signalements de dysfonctionnements mineurs par les opérateurs, et de demandes d'ajustement de données par les citoyens. 

Pour éviter l'engorgement du support technique standard et assurer une stabilisation rapide du système, la phase de **Hypercare** est activée pour chaque site de déploiement pendant une **durée fixe de 30 jours** à compter du Go-Live.

```
                           HYPERCARE SUPPORT PIPELINE
                           
       [Incident Signalé par un Opérateur / Citoyen]
                            |
                            v (Automatic routing - Tier 1 Support)
       [Triage Hypercare Central - Qualification en < 10 mins]
                            |
            +---------------+---------------+
            | (Incident Connu)              | (Incident Complexe / Bloquant)
            v                               v
[Résolution Immédiate par Base de Connaissances]  [Escalade d'Urgence Tier 2 & 3]
- Temps de réponse moyen < 15 mins          - Correction de bug en patch chaud (Hotfix)
- Validation et fermeture immédiate         - Prise en charge par l'équipe d'ingénierie
```

---

## 2. Structure et Organisation du Support Hypercare

Le dispositif s'articule autour de trois niveaux de support dédiés et hautement réactifs.

### 2.1 Les Niveaux de Support (Tiers 1, 2, 3)

1. **Niveau 1 : Triage et Assistance de Premier Niveau (Tier 1 Support)**
   *   *Rôle :* Prise en charge des appels et des tickets émis par les bureaux de liaison. Aide à l'utilisation du matériel biométrique (ex : "Le capteur d'empreinte digitale ne s'allume pas").
   *   *Objectif de SLA :* Prise en charge en moins de 10 minutes, résolution en moins de 1 heure.
2. **Niveau 2 : Expertise Métier et Technique de Proximité (Tier 2 Support)**
   *   *Rôle :* Gestion des conflits d'identité démographique complexes non résolus automatiquement par le NIRE. Réinstallation rapide à distance des agents d'Edge Nodes.
   *   *Objectif de SLA :* Résolution en moins de 4 heures.
3. **Niveau 3 : Équipe d'Ingénierie Noyau et Éditeurs (Tier 3 Support)**
   *   *Rôle :* Correction de failles de sécurité, résolution de crashs d'infrastructure ou de bugs critiques dans le code de l'ABIS central ou de l'API Gateway.
   *   *Objectif de SLA :* Déploiement d'un patch d'urgence (Hotfix) en moins de 12 heures.

---

## 3. Gestion de l'Afflux d'Incidents (Incident Surge Support)

Afin d'éviter la paralysie des équipes de support face à un volume d'appels exceptionnel durant la première semaine de déploiement national :
*   **La Base de Connaissances Interactive (Self-Help KB) :** Les LEN intègrent un chatbot d'assistance hors-ligne (SNISID Assistant) capable d'aider les opérateurs à résoudre eux-mêmes 80% des problèmes matériels connus sans ouvrir de ticket.
*   **Les Équipes Volantes (Field Strike Teams) :** Deux ingénieurs réseaux et un formateur SNISID chevronné sont pré-positionnés dans le chef-lieu de chaque département prêt à intervenir physiquement en moins de 2 heures dans n'importe quel bureau de liaison en difficulté.

---

## 4. Métriques de Stabilisation (Critères de Sortie d'Hypercare)

Un département ne peut quitter la phase d'assistance renforcée Hypercare pour basculer dans le régime de maintenance opérationnelle standard (BAU - Business As Usual) que si l'intégralité des indicateurs de stabilisation suivants sont validés sur une période glissante de **7 jours consécutifs** :

```
                  HYPERCARE EXIT CRITERIA SCORECARD
                  
+------------------------------------+--------------------------+
| Indicateur de Stabilisation        | Seuil de Validation      |
+------------------------------------+--------------------------+
| Disponibilité globale des LEN      | >= 99.95%                |
| Tickets d'urgence de niveau S1     | 0 ticket actif           |
| Tickets de niveau S2 en souffrance | < 2 tickets par dépt     |
| Temps de réponse moyen de l'ABIS   | < 1.5 seconde            |
| Taux de satisfaction des opérateurs| >= 92%                   |
+------------------------------------+--------------------------+
```
---

## 5. Registre et Outils d'Enregistrement

Pour industrialiser cette phase, les incidents sont suivis et classés au travers d'un référentiel de ticketing automatisé, dont la structure et le moteur de tri sont détaillés par le script d'automatisation de ce répertoire.
