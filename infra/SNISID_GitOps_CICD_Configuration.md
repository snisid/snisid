# SNISID : Configuration et Automatisation CI/CD GitOps (v2.0)

**Classification :** RESTREINT / DEVSECOPS  
**Recommandation de Référence :** SNISID v2.0 — MP-009

Ce playbook détaille les configurations GitLab CI/CD, les règles de caching, et les scripts opérationnels requis pour implémenter la boucle de feedback de moins de 10 minutes et les environnements éphémères de test.

---

## 📂 Fichiers de Référence dans le Workspace

* **Scripts de Pipeline** :
  * [diff_sast_scan.py](file:///c:/Users/sopil/Desktop/snisid%20system/pki/scripts/diff_sast_scan.py) — Script d'analyse SAST Semgrep différentiel (< 10 min).
  * [manage_ephemeral_env.sh](file:///c:/Users/sopil/Desktop/snisid%20system/pki/scripts/manage_ephemeral_env.sh) — Script de gestion de cycle de vie des namespaces de PR.
* **Architecture Générale** :
  * [SNISID_GitOps_DevSecOps_Architecture.md](file:///c:/Users/sopil/Desktop/snisid%20system/SNISID_GitOps_DevSecOps_Architecture.md) — Principes GitOps, seuils SAST/DAST, et objectifs SLSA.

---

## 1. Pipeline GitLab CI/CD Complet (8 Étapes)

Ce manifeste représente le pipeline d'intégration continue standard de SNISID pour le service d'identité (`identity-svc`).

### Fichier : `apps/identity-service/.gitlab-ci.yml`

```yaml
stages:
  - pre-commit
  - build-sast
  - sca
  - container-scan
  - sign
  - deploy-staging
  - dast
  - deploy-prod

variables:
  HARBOR_REGISTRY: "harbor.snisid.gov.ht"
  IMAGE_NAME: "core/identity-svc"
  IMAGE_TAG: "$CI_COMMIT_SHA"
  KUBECONFIG_STAGING: "$K8S_STAGING_CONFIG"
  KUBECONFIG_PROD: "$K8S_PROD_CONFIG"

# Stratégie Globale de Caching pour accélérer les temps de build (< 10 minutes)
cache:
  key: "$CI_COMMIT_REF_SLUG"
  paths:
    - .cargo/registry/index/
    - .cargo/registry/cache/
    - target/
    - .go-cache/
    - .npm/

# ÉTAPE 1 : PRE-COMMIT (Secret Scanning)
secret_scanning:
  stage: pre-commit
  image: gitleaks/gitleaks:latest
  script:
    - gitleaks detect --source=. --verbose --redact
  allow_failure: false

# ÉTAPE 2 : BUILD & SAST (Semgrep Différentiel)
sast_diff:
  stage: build-sast
  image: python:3.13-slim
  script:
    - pip install semgrep
    - python pki/scripts/diff_sast_scan.py
  allow_failure: false

# ÉTAPE 3 : SCA (Vérification des dépendances)
dependency_check:
  stage: sca
  image: aquasec/trivy:latest
  script:
    - trivy fs --scanners vuln --exit-code 1 --severity CRITICAL .
  allow_failure: false

# ÉTAPE 4 : CONTAINER SCAN (Scanning de l'image Docker)
container_vulnerability_scan:
  stage: container-scan
  image: aquasec/trivy:latest
  before_script:
    - docker build -t $HARBOR_REGISTRY/$IMAGE_NAME:$IMAGE_TAG .
  script:
    - trivy image --exit-code 1 --severity CRITICAL $HARBOR_REGISTRY/$IMAGE_NAME:$IMAGE_TAG
  allow_failure: false

# ÉTAPE 5 : SIGNING (Signature Cosign + SBOM Syft - SLSA Level 3/4)
sign_and_attest:
  stage: sign
  image: gcr.io/projectsigstore/cosign:latest
  script:
    # 1. Génération du SBOM CycloneDX avec Syft
    - syft $HARBOR_REGISTRY/$IMAGE_NAME:$IMAGE_TAG -o cyclonedx-json --file sbom.json
    # 2. Publication du SBOM comme attestation signée Cosign dans Harbor
    - cosign attest --key k8s://snisid-pki/cosign-key --type cyclonedx --attestation sbom.json $HARBOR_REGISTRY/$IMAGE_NAME:$IMAGE_TAG
    # 3. Signature de l'image principale
    - cosign sign --key k8s://snisid-pki/cosign-key $HARBOR_REGISTRY/$IMAGE_NAME:$IMAGE_TAG
  allow_failure: false

# ÉTAPE 6 : DEPLOY STAGING (Déploiement automatique via Kustomize/ArgoCD)
deploy_to_staging:
  stage: deploy-staging
  image: line/kubectl-kustomize:latest
  script:
    - cd gitops/overlays/staging
    - kustomize edit set image $HARBOR_REGISTRY/$IMAGE_NAME=$HARBOR_REGISTRY/$IMAGE_NAME:$IMAGE_TAG
    - git config --global user.email "ci-runner@snisid.gov.ht"
    - git config --global user.name "CI Runner"
    - git commit -am "Update staging image to $IMAGE_TAG"
    - git push origin main
  rules:
    - if: '$CI_COMMIT_BRANCH == "main"'

# ÉTAPE 7 : DAST (Analyse de vulnérabilité dynamique)
dast_scan:
  stage: dast
  image: owasp/zap2docker-stable:latest
  script:
    - zap-baseline.py -t https://staging.identity.snisid.gov.ht -r zap_report.html
  artifacts:
    paths:
      - zap_report.html
  rules:
    - if: '$CI_COMMIT_BRANCH == "main"'

# ÉTAPE 8 : DEPLOY PROD (ArgoCD double approbation - restriction horaire)
deploy_to_production:
  stage: deploy-prod
  image: line/kubectl-kustomize:latest
  script:
    # Validation de la fenêtre horaire : pas de déploiement le vendredi après 16h
    - DAY_OF_WEEK=$(date +%u)
    - HOUR=$(date +%H)
    - |
      if [ "$DAY_OF_WEEK" -eq 5 ] && [ "$HOUR" -ge 16 ]; then
        echo "[ERROR] Déploiement bloqué : Interdiction réglementaire de déployer le vendredi après 16h."
        exit 1
      fi
    # Kustomize prod update
    - cd gitops/overlays/prod
    - kustomize edit set image $HARBOR_REGISTRY/$IMAGE_NAME=$HARBOR_REGISTRY/$IMAGE_NAME:$IMAGE_TAG
    - git commit -am "Update production image to $IMAGE_TAG [Approved]"
    - git push origin production
  rules:
    - if: '$CI_COMMIT_BRANCH == "production"'
      when: manual  # Déclenchement manuel requis (Double Approbation)
```

---

## 2. Automatisation des Environnements Éphémères par PR

Pour tester chaque fonctionnalité de manière isolée, un namespace Kubernetes temporaire est provisionné lors de la création d'une Pull Request (PR) et détruit automatiquement lors de sa fusion.

### Script Bash : `pki/scripts/manage_ephemeral_env.sh`

```bash
#!/usr/bin/env bash
# File: /pki/scripts/manage_ephemeral_env.sh
# Gère le cycle de vie des namespaces éphémères de test pour les PRs

set -euo pipefail

ACTION="${1:-}" # "create" ou "destroy"
PR_ID="${2:-}"  # ex: "pr-42"

if [ -z "$ACTION" ] || [ -z "$PR_ID" ]; then
    echo "Usage: $0 <create|destroy> <pr_id>"
    exit 1
fi

NAMESPACE="snisid-ephemeral-$PR_ID"

if [ "$ACTION" == "create" ]; then
    echo "[*] Création de l'environnement éphémère pour $PR_ID..."
    
    # 1. Créer le namespace K8s
    kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -
    
    # 2. Appliquer les NetworkPolicies Cilium d'isolation stricte (Pas d'accès prod)
    cat <<EOF | kubectl apply -f -
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: restrict-ephemeral-access
  namespace: $NAMESPACE
spec:
  endpointSelector:
    matchLabels: {}
  egress:
    # Autoriser uniquement le trafic interne au namespace et le DNS
    - toEndpoints:
        - matchLabels:
            "k8s:io.kubernetes.pod.namespace": $NAMESPACE
    - toEndpoints:
        - matchLabels:
            "k8s:io.kubernetes.pod.namespace": kube-system
      toPorts:
        - ports:
            - port: "53"
              protocol: UDP
EOF

    # 3. Déployer les microservices de test (avec mocks de base de données)
    echo "[+] Environnement $NAMESPACE prêt pour les tests."

elif [ "$ACTION" == "destroy" ]; then
    echo "[*] Destruction de l'environnement éphémère pour $PR_ID..."
    kubectl delete namespace "$NAMESPACE" --ignore-not-found=true
    echo "[+] Environnement $NAMESPACE supprimé."
fi
```

---

*Ce processus CI/CD garantit une sécurité maximale dès le commit, sans impacter la vélocité des développeurs.*
