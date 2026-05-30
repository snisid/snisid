# National Resilience KPI Model

## 1. Objectif
Mesurer la résilience souveraine SNISID par des KPI suivis par le NRCC et utilisés pour améliorer architecture, doctrine et exercices.

## 2. KPI principaux
| KPI | Objectif | Définition |
|---|---|---|
| Recovery Time Objective (RTO) | Minimal | délai maximal acceptable de restauration |
| Recovery Point Objective (RPO) | Minimal | perte maximale acceptable de données |
| DR readiness score | Élevé | capacité mesurée de failover/recovery |
| Offline survivability | Maximale | durée et couverture opérations sans réseau |
| Crisis response latency | Faible | délai détection → activation → décision |

## 3. Seuils
| KPI | Vert | Orange | Rouge |
|---|---:|---:|---:|
| RTO N0 | ≤ 2 h | 2-4 h | > 4 h |
| RPO N0 | ≤ 15 min | 15-60 min | > 60 min |
| DR readiness | ≥ 90% | 75-89% | < 75% |
| Offline survivability | ≥ 90% régions prioritaires | 70-89% | < 70% |
| Crisis response latency L3 | ≤ 30 min | 30-60 min | > 60 min |

## 4. DR readiness score
| Dimension | Poids |
|---|---:|
| santé réplication | 20% |
| backups testés | 20% |
| runbooks à jour | 15% |
| automatisation recovery | 15% |
| capacité site DR | 15% |
| disponibilité équipes/contacts | 10% |
| observability/alerting | 5% |

## 5. Reporting
Daily status, rapport backup hebdomadaire, DR readiness mensuel, rapport exercice après test, revue souveraine trimestrielle.
