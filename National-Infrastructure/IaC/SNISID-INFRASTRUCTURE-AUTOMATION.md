---
# ============================================================
# SNISID-Infra — Infrastructure Automation Platform (IaC)
# Terraform, Ansible, et GitOps
# Document ID: SNISID-IAC-001
# Version: 1.0.0
# ============================================================

## 1. INFRASTRUCTURE AS CODE (La fin du Clic)

Pour garantir la reproductibilité et éviter le "Configuration Drift" (un admin modifie une règle de firewall à la main et l'oublie), TOUTE l'infrastructure est gérée par du code.
Personne ne se connecte en SSH sur un serveur de production pour faire une modification manuelle (Immutable Infrastructure).

## 2. GITOPS PIPELINE

1. **Le Code :** Un ingénieur réseau veut ouvrir le port 443 pour un nouveau service. Il modifie le fichier Terraform (`firewall.tf`) sur son poste.
2. **Review :** Il pousse le code sur le Gitlab du gouvernement (Git). Une Merge Request (MR) est créée. Deux autres ingénieurs doivent approuver (Principe des 4 yeux).
3. **Automated Testing :** La CI/CD vérifie que cette règle ne viole pas le *Security Framework* (ex: OPA/Checkov scan).
4. **Déploiement (ArgoCD/Terraform Cloud) :** Une fois mergé sur la branche `main`, le robot d'automatisation applique la modification sur le pare-feu physique.

## 3. BARE METAL AUTOMATION (Metal-as-a-Service)

Lorsqu'un nouveau rack de serveurs physiques est branché :
- Ils bootent en PXE.
- Le service d'automatisation (MaaS / Ironic) détecte le CPU/RAM/Disques, installe l'OS standardisé (ex: Ubuntu Server durci / Talos Linux).
- Ansible configure les clés SSH (Uniquement pour le robot), durcit le kernel (Sysctl), et joint le noeud au cluster Kubernetes/OpenStack.
- Tout se fait sans intervention humaine au-delà du branchement initial.

---
*Document ID: SNISID-IAC-001 | Approuvé par: Head of Platform Engineering*
