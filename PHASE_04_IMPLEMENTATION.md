# PHASE_04_IMPLEMENTATION.md

## Nom de la phase
Phase 4 - Infrastructure-as-Code (IaC) & Automatisation Datacenter

## Objectif
Provisionner et automatiser l'architecture matérielle et cloud-native (Kubernetes/Talos sur Proxmox) pour les datacenters souverains, en utilisant Terraform, Packer et Ansible.

## Fonctionnalités ajoutées
- Modèles Packer pour générer des images de machines virtuelles immutables (Talos OS).
- Modules Terraform pour le provisionnement d'infrastructure (Proxmox K8s, Vault, Cilium).
- Playbooks Ansible de durcissement (OS Hardening) et de déploiement réseau.
- Manifests Kubernetes de base (ArgoCD, cert-manager, ingress-nginx).

## Fichiers créés / intégrés
L'ensemble de l'architecture d'automatisation `infrastructure/` a été intégré au projet sous le répertoire `SNISID-Infrastructure/` :
- `SNISID-Infrastructure/ansible/`
- `SNISID-Infrastructure/packer/`
- `SNISID-Infrastructure/terraform/`
- `SNISID-Infrastructure/kubernetes/`

## Fichiers modifiés
Aucun. L'intégration s'est faite par un ajout d'arborescence à la racine de l'écosystème.

## Dépendances ajoutées
L'exécution de cette architecture requiert les outils suivants installés sur l'environnement d'administration CI/CD :
- Terraform (>= 1.5.0)
- HashiCorp Packer
- Ansible
- Proxmox VE (environnement cible)

## Variables d’environnement
- `PM_API_URL`, `PM_USER`, `PM_PASS` : Identifiants API pour Terraform -> Proxmox.
- `VAULT_ADDR`, `VAULT_TOKEN` : Identifiants pour HashiCorp Vault.

## Migrations ou changements de base de données
- Aucun changement applicatif direct. Prépare l'infrastructure des bases.

## Commandes de test / build / déploiement
L'initialisation des modules Terraform peut être testée structurellement sans provisionnement physique (Dry Run) :
```bash
cd "SNISID-Infrastructure\terraform\providers"
terraform init
terraform validate
terraform plan
```

## Procédure de rollback
Pour retirer le code d'infrastructure IaC du projet :
```bash
Remove-Item -Path "c:\Users\sopil\Desktop\snisid system\SNISID-Infrastructure" -Recurse -Force
```

## Risques connus
- L'exécution de `terraform apply` créera physiquement des VMs sur les clusters Proxmox et engendrera des coûts de ressources. À n'exécuter qu'avec l'aval du responsable des Datacenters.

## Points à valider manuellement
- Générer l'image initiale Talos via `packer build proxmox-talos.pkr.hcl`.
- Connecter le pipeline GitOps (ArgoCD) au dépôt de configuration nouvellement créé.
