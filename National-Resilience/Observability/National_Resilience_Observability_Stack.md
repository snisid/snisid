# National Resilience Observability Stack

## 1. Objectif
Donner une visibilité nationale sur sauvegardes, DR, récupération, santé offline, énergie, crise et cyber-résilience.

## 2. Domaines monitorés
| Domaine | Monitoring | Indicateurs |
|---|---:|---|
| Backup health | Oui | succès jobs, âge backups, immutabilité, tests restore |
| DR readiness | Oui | réplication, capacité failover, runbooks prêts |
| Recovery latency | Oui | temps restauration, RTO réel, étapes bloquées |
| Offline node health | Oui | edge nodes, âge cache, batterie, file sync |
| Crisis escalation status | Oui | niveau alerte, incidents, communications |

## 3. Outils
| Domaine | Outil |
|---|---|
| Monitoring | Prometheus |
| Dashboards | Grafana |
| Logging | Loki |
| Incident management | PagerDuty ou équivalent open-source |

## 4. Dashboards
| Dashboard | Contenu |
|---|---|
| National Resilience Overview | P0/P1, sites, crise, énergie |
| Backup & Restore Health | jobs, tests, RPO, immutabilité |
| DR Readiness | réplication, failover, capacité sites |
| Offline Survival | edge nodes, cache, sync, batteries |
| Cyber Resilience | attaques, containment, backups propres |
| Power Autonomy | UPS, générateurs, carburant |

## 5. Alertes critiques
| Alerte | Sévérité | Action |
|---|---|---|
| Backup B0 failed > 1 cycle | critique | incident immédiat |
| Restore test failed | critique | correction + retest |
| Replication lag N0 > seuil | critique | analyse DR |
| Offline cache expired | high | rafraîchir ou mode restreint |
| Fuel < seuil | critique | réapprovisionnement urgent |
| IAM integrity anomaly | critique | cyber containment |

## 6. Métriques exemples
```text
snisid_backup_last_success_timestamp{class="B0"}
snisid_backup_restore_test_success{class="B0"}
snisid_dr_replication_lag_seconds{site="secondary"}
snisid_recovery_workflow_duration_seconds{scenario="primary_dc_loss"}
snisid_offline_node_sync_queue_depth{region="north"}
snisid_power_autonomy_hours{site="primary_dc"}
snisid_crisis_level{scope="national"}
```
