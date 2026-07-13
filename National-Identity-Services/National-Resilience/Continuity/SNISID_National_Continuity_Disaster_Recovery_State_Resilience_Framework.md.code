# SNISID National Continuity, Disaster Recovery & State Resilience Framework

## 1. Mission
Établir la doctrine nationale permettant à SNISID de soutenir la continuité gouvernementale, la continuité d'infrastructure, la continuité d'identité, les opérations d'urgence et la récupération nationale après catastrophe.

## 2. Principe souverain
SNISID doit fonctionner même si un datacenter est détruit, si Internet national s'effondre, si une cyberattaque massive survient ou si des régions sont isolées. Les services P0/P1 doivent rester disponibles, basculer en DR ou fonctionner offline.

## 3. Domaines couverts
| Domaine | Support | Objectif |
|---|---:|---|
| Government continuity | Oui | maintenir décision, coordination et services publics critiques |
| Infrastructure continuity | Oui | production, failover, restauration, mode dégradé |
| Identity continuity | Oui | vérification, émission et gouvernance de l'identité |
| Emergency operations | Oui | opérations terrain, population, registre civil d'urgence |
| National recovery | Oui | reconstruction automatisée, validation et retour normal |

## 4. Modes opérationnels
| Mode | Déclencheur | Capacités actives |
|---|---|---|
| Normal | aucune crise | production primaire, sauvegardes, monitoring |
| Dégradé | incident local | priorisation P0/P1, réduction charge |
| Crise nationale | catastrophe/cyberattaque | NRCC 24/7, DR, communications crise |
| Survie offline | Internet/DC indisponible | vérification offline, edge sync différée |
| Reconstruction | perte/corruption majeure | IaC, GitOps, backups immuables, runbooks |

## 5. Gouvernance
| Organe | Responsabilités |
|---|---|
| National Resilience Command Center (NRCC) | activation crise, orchestration DR, décisions opérationnelles 24/7 |
| Government Crisis Council | arbitrage stratégique interinstitutionnel |
| SNISID Resilience Office | doctrine, tests, KPI, audits |
| Cyber Resilience Cell | containment, forensic, clean recovery |
| Regional Continuity Cells | opérations locales, kits offline, relais terrain |

## 6. Priorités de restauration
1. réseau de crise, DNS interne, bastion ;
2. KMS/HSM ou mécanisme de clés ;
3. IAM gouvernemental ;
4. registres d'identité N0 ;
5. APIs de vérification identité ;
6. enrôlement et registre civil d'urgence ;
7. portails citoyens et reporting.

## 7. Exigences minimales
- Multi-région souverain avec sites primary, secondary, régional et vault offline.
- Sauvegardes chiffrées, immuables, hors ligne et testées.
- Fonctionnement offline pour identité et enrôlement critique.
- Recovery automation par Terraform, ArgoCD, Ansible, Velero.
- Communications de crise : messagerie sécurisée, satellite, radio, alertes multi-canaux.
- Autonomie énergétique des sites critiques.
- Exercices réguliers et mesures RTO/RPO.

## 8. Politique d'identité en crise
Les identités N0 doivent être vérifiables en ligne, en DR et offline. Les émissions d'urgence sont limitées dans le temps, signées, auditables et réconciliées après crise. Les clés sensibles exigent double contrôle et stockage souverain.

## 9. Chaîne de décision
```text
Détection → Qualification → Activation NRCC → Priorisation P0/P1 → Failover/Offline → Validation → Retour contrôlé
```

## 10. Définition de réussite
La continuité nationale est atteinte lorsque les scénarios extrêmes sont testés, que les RTO/RPO sont respectés, que les régions peuvent opérer offline, et que les équipes peuvent reconstruire les services critiques sans dépendance à un site unique.
