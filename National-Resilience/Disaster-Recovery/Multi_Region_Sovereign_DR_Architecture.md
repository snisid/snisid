# Multi-Region Sovereign Disaster Recovery Architecture

## 1. Objectif
Créer une reprise après désastre souveraine, multi-sites, capable de restaurer ou maintenir SNISID après destruction, corruption ou indisponibilité d'un site majeur.

## 2. Topologie
| Site | Fonction | Mode | Données |
|---|---|---|---|
| Primary National DC | production | hot | données actives complètes |
| Secondary National DC | failover national | hot/warm | réplication quasi temps réel |
| Regional DR Sites | continuité régionale | warm | caches critiques, services P0/P1 |
| Offline Vault Site | sauvegardes critiques | cold/air-gapped | backups immuables, clés recovery |

## 3. Capacités
| Fonction | Support | Exigence |
|---|---:|---|
| Multi-datacenter replication | Oui | réplication chiffrée, surveillée, suspendable |
| Cross-region failover | Oui | bascule orchestrée par NRCC |
| Immutable backups | Oui | WORM/Object Lock/offline vault |
| Cold/warm/hot sites | Oui | selon criticité N0-N3 |
| Autonomous recovery | Oui | IaC, GitOps, Ansible, Velero |

## 4. Design logique
```text
Agences/Citoyens → Traffic Manager souverain/DNS crise
                  ├→ Primary National DC
                  ├→ Secondary National DC
                  ├→ Regional DR Sites
                  └→ Offline Vault Site
```

## 5. Stratégie de réplication
| Donnée | Méthode | Fréquence | Protection |
|---|---|---:|---|
| Registre identité N0 | sync/near-sync | secondes-minutes | chiffrement, intégrité |
| Événements identité | log shipping | continu | append-only |
| Config plateformes | GitOps mirror | continu | signatures |
| Biométrie critique | réplication contrôlée | politique souveraine | chiffrement fort |
| Sauvegardes complètes | backup chiffré | quotidien/hebdo | immuable + offline |

## 6. Failover
Déclencheurs : perte primary DC, corruption majeure, cyberattaque, blackout prolongé ou décision NRCC.

Séquence : geler changements non essentiels → qualifier données → isoler site compromis → promouvoir secondary/DR → vérifier IAM/registre → ouvrir P0 puis P1 → surveiller.

## 7. Règles d'intégrité
- Ne jamais propager une corruption connue.
- Suspendre réplication lors d'attaque ou corruption.
- Conserver snapshots avant promotion.
- Séparer clés de récupération et production.

## 8. Tests
| Test | Fréquence |
|---|---:|
| failover applicatif P0 | mensuel |
| failover DC complet | semestriel |
| restauration backup immuable | mensuel |
| exercice offline vault | trimestriel |
| reconstruction cluster | trimestriel |
