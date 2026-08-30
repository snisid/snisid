# SNISID National Deployment War Room Governance
## Cadre de Pilotage en Temps Réel et de Gestion des Incidents Majeurs

---

## 1. Rôle et Objectif de la War Room Nationale

Pendant toute la durée de la phase de déploiement national et de migration (Phase 15), une cellule de crise permanente appelée **National Deployment War Room** est activée. Sa mission est de superviser en temps réel le déploiement sur l'ensemble des 10 départements d'Haïti, de piloter la migration de données historique, de coordonner les résolutions de pannes d'infrastructure, et de débloquer les obstacles opérationnels sur le terrain.

```
                           WAR ROOM ESCALATION PATH
                           
  [Local Edge Node Operator] (Level 1 - Operator)
              |
              v (Field Incident - unresolved in 15 mins)
  [Departmental Coordinator] (Level 2 - Regional Supervisor)
              |
              v (Major Blocker / Connectivity / Security)
  [National War Room Command Center] (Level 3 - Technical & Operational Directors)
              |
              v (National Sovereignty Impact / Disaster)
  [Executive Crisis Board] (Level 4 - Ministers & Directors General)
```

---

## 2. Structure Organisationnelle de la War Room

La War Room réunit des experts multidisciplinaires organisés en 4 cellules spécialisées travaillant sous la supervision d'un **Incident Commander**.

```
                           [INCIDENT COMMANDER]
                                    |
       +--------------------+-------+-------+--------------------+
       |                    |               |                    |
       v                    v               v                    v
 [Infra-Ops Desk]     [Field-Ops Desk] [Security Desk]    [Comm Desk]
 - Cloud & Datacenters- Operator Training- Cryptography   - Citizen updates
 - Starlink & Fiber   - Device Delivery - Fraud Audits    - Radio & Press
```

### 2.1 Les Quatre Desks de Crise
1. **Infra-Ops Desk (Pupitre Réseau & Datacenters) :**
   *   *Rôle :* Assurer la disponibilité du DC principal et secondaire, de l'ABIS, et surveiller l'état de la connectivité Starlink et fibre de chaque commune.
2. **Field-Ops Desk (Pupitre Déploiement Terrain) :**
   *   *Rôle :* Interface avec les bureaux de liaison communaux (BLC). Gérer la logistique de livraison de matériel d'enrôlement et l'assistance technique directe aux opérateurs.
3. **Security & Fraud Desk (Pupitre Sécurité & Lutte contre la Fraude) :**
   *   *Rôle :* Analyser les suspicions d'usurpation biométrique transmises par le moteur de réconciliation, investiguer les attaques réseau ou les intrusions physiques de nœuds d'Edge.
4. **Communication & Public Relations Desk (Pupitre Média & Citoyen) :**
   *   *Rôle :* Rédiger les alertes à destination de la population, surveiller les rumeurs sur les réseaux sociaux (Sogemedia, Twitter, Facebook) et rédiger les communiqués de presse officiels.

---

## 3. Matrice d'Escalade et Niveaux de Gravité (SLA & Severity Levels)

Chaque anomalie détectée sur le terrain ou sur l'infrastructure centrale se voit attribuer un niveau de sévérité dictant le canal de traitement et le délai maximum de résolution (SLA).

| Niveau | Désignation | Description Technique | SLA de Résolution | Protocole de Notification |
| :--- | :--- | :--- | :--- | :--- |
| **S1** | **Critique / Bloquant** | Panne complète du DC Central, corruption de base de données, ou interruption de connectivité de plus de 5 départements. | **$< 30$ minutes** | SMS automatique + Appel vocal de l'Incident Commander au Directeur National. |
| **S2** | **Majeur** | Panne de connectivité d'un département complet, ou dysfonctionnement de l'ABIS biométrique 1-to-N. | **$< 2$ heures** | Canal Slack `#war-room-s2-alerts` + Notification par e-mail. |
| **S3** | **Moyen** | Panne d'un Edge Node d'une commune isolée (bascule automatique sur mode offline). | **$< 6$ heures** | Ticket consigné dans le système d'incident avec assignation à l'équipe locale de maintenance. |
| **S4** | **Mineur** | Problème d'impression d'une carte d'identité individuelle, ou question d'usage d'un opérateur. | **$< 24$ heures** | Traité par la file d'attente standard du Support Hypercare. |

---

## 4. Rituels Opérationnels de la War Room
*   **Daily Standup (07:30 - Matin) :** Revue rapide de l'état des systèmes avant l'ouverture des bureaux de liaison. Validation du score de readiness de la journée.
*   **Hourly Dashboard Review (Toutes les heures) :** Analyse du tableau de bord de déploiement (progression des enrôlements, taux de rejet biométrique, taux de synchro différée).
*   **Daily Retrospective (18:30 - Soir) :** Bilan des incidents résolus de la journée, clôture comptable des enrôlements nationaux, et ajustement de la logistique du lendemain.
