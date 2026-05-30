# National Recovery Automation Platform

## 1. Objectif
Automatiser la reconstruction nationale de SNISID : infrastructure, Kubernetes, bases de données et services d'identité.

## 2. Technologies
| Domaine | Technologie | Usage |
|---|---|---|
| IaC | Terraform | réseaux, compute, stockage, DNS interne |
| GitOps | ArgoCD | restauration déclarative Kubernetes |
| Automation | Ansible | bootstrap, procédures OS/applicatives |
| Cluster restore | Velero | restauration Kubernetes et volumes |

## 3. Capacités
| Fonction | Support | Mécanisme |
|---|---:|---|
| Infrastructure rebuild | Oui | modules Terraform + états sauvegardés |
| Kubernetes restoration | Oui | bootstrap cluster + ArgoCD + Velero |
| Database restoration | Oui | PITR, snapshots, log replay, checksum |
| Identity restoration | Oui | IAM, annuaires, clés, politiques accès |

## 4. Pipeline
1. sélectionner scénario ;
2. isoler sources compromises ;
3. choisir dernier point sain ;
4. provisionner par Terraform ;
5. restaurer réseau, bastion, clés ;
6. reconstruire Kubernetes ;
7. restaurer bases et volumes ;
8. déployer via ArgoCD ;
9. valider IAM, registre, APIs ;
10. ouvrir progressivement au trafic.

## 5. Standards
Scripts versionnés et signés, playbooks idempotents, secrets hors Git, dépendances mirorées offline, logs d'exécution archivés.

## 6. Workflow exemple
```yaml
recovery_workflow:
  scenario: primary_dc_loss
  target_site: secondary_national_dc
  restore_point: last_verified_clean
  priority_order: [network, keys, iam, identity_registry, verification_api, enrollment]
  approvals_required: 2
  validation_required: true
```
