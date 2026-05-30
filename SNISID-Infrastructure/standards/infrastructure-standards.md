# SNISID — National Infrastructure Standards
**Classification:** RESTREINT DEFENSE  
**Version:** 4.0.0  
**Statut:** OBLIGATOIRE — Toute violation bloque CI/CD

---

## 1. Kubernetes Standards

### 1.1 Versions & Lifecycle
| Règle | Valeur |
|-------|--------|
| Version K8s | 1.28.x LTS (aligné RKE2/K3s) |
| Support window | N-2 versions maximum |
| Upgrade window | Dimanche 02h00-06h00 national |
| Feature gates | Uniquement GA (pas de beta/alpha en prod) |
| CRDs | Pas de CRD alpha. Review IGC obligatoire pour CRD custom. |

### 1.2 Resource Standards
```yaml
# Pod minimal national
apiVersion: v1
kind: Pod
metadata:
  labels:
    snisid.gov.tier: "tier-1"          # Obligatoire
    snisid.gov.region: "core"          # Obligatoire
    snisid.gov.service: "mon-service"   # Obligatoire
    snisid.gov.owner: "equipe-x"        # Obligatoire
    snisid.gov.data-classification: "restricted"  # Obligatoire
    app.kubernetes.io/name: "mon-service"         # Obligatoire
    app.kubernetes.io/version: "v4.2.1"             # Obligatoire
    app.kubernetes.io/component: "api"              # Recommandé
    app.kubernetes.io/part-of: "snisid-core"      # Recommandé
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 65534
    runAsGroup: 65534
    fsGroup: 65534
    seccompProfile:
      type: RuntimeDefault
    sysctls: []  # Aucun sysctl custom sans whitelist IGC
  containers:
    - name: app
      image: "registry.interne.snisid.gouv.local/mon-service:v4.2.1@sha256:abc123..."  # Digest obligatoire
      imagePullPolicy: Always
      securityContext:
        readOnlyRootFilesystem: true
        allowPrivilegeEscalation: false
        capabilities:
          drop:
            - ALL
        runAsNonRoot: true
        runAsUser: 65534
      resources:
        requests:
          memory: "128Mi"   # Jamais unspecified
          cpu: "100m"       # Jamais unspecified
        limits:
          memory: "512Mi"
          cpu: "1000m"
      livenessProbe:
        httpGet:
          path: /health/live
          port: 8080
        initialDelaySeconds: 10
        periodSeconds: 15
      readinessProbe:
        httpGet:
          path: /health/ready
          port: 8080
        initialDelaySeconds: 5
        periodSeconds: 5
      startupProbe:
        httpGet:
          path: /health/startup
          port: 8080
        failureThreshold: 30
        periodSeconds: 10
      volumeMounts:
        - name: tmp
          mountPath: /tmp
        - name: cache
          mountPath: /var/cache
  volumes:
    - name: tmp
      emptyDir:
        sizeLimit: 100Mi
    - name: cache
      emptyDir:
        sizeLimit: 500Mi
```

### 1.3 Namespace Standards
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: snisid-core
  labels:
    snisid.gov.tier: tier-1
    snisid.gov.region: core
    istio-injection: enabled
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

## 2. GitOps Standards

### 2.1 Repository Structure
```
applications/national.git
├── helm/
│   ├── snisid-core/
│   │   ├── Chart.yaml
│   │   ├── values.yaml (DEFAULTS UNIQUEMENT — jamais de secrets)
│   │   ├── values-core-prod.yaml
│   │   ├── values-dr-prod.yaml
│   │   ├── values-edge-regional.yaml
│   │   └── templates/
├── ci/
│   └── .gitlab-ci.yml / github-actions.yml
└── README.md

infrastructure/national.git
├── kubernetes/
├── terraform/
├── gitops/
├── security/
├── observability/
└── standards/
```

### 2.2 Commit & Merge Rules
| Règle | Valeur |
|-------|--------|
| Branche production | `main` (protégée) |
| Branche staging | `staging` |
| PR vers main | 2 reviewers minimum, dont 1 IGC pour Tier-0 |
| CI obligatoire | Lint, Test, Scan Trivy, Validate Kyverno, Terraform plan |
| Secrets dans repo | **INTERDIT** — bloqué par pre-commit + gitleaks + CI |
| Sign commits | GPG obligatoire pour `main` |
| ArgoCD sync | Auto-sync ON, selfHeal ON, prune ON (prod) |

## 3. Terraform Standards

### 3.1 Backend & State
```hcl
terraform {
  backend "s3" {
    bucket                      = "snisid-terraform-state"
    key                         = "${var.environment}/${var.component}.tfstate"
    region                      = "national"
    endpoint                    = "https://s3.interne.snisid.gouv.local"
    encrypt                     = true
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_region_validation      = true
    force_path_style            = true
  }
}
```

### 3.2 Module Standards
- Tout module dans `infrastructure/terraform/modules/`
- Versioning par Git tag (`git tag modules/proxmox-k8s/v1.2.3`)
- `terraform plan` obligatoire en CI, `terraform apply` uniquement via CI/CD (pas local)
- State locking via DynamoDB-compatible (ceph-rgw) ou Consul

## 4. Helm Standards

| Règle | Valeur |
|-------|--------|
| Chart lint | `ct lint` obligatoire en CI |
| Values chiffrées | SOPS + age (clés IGC) pour environments prod |
| Image tags | Jamais `latest`. Tag sémantique + digest SHA256 |
| Réplicas Tier-0 | Minimum 3, PodDisruptionBudget minAvailable=2 |
| Réplicas Tier-1 | Minimum 2, PDB minAvailable=1 |
| Hooks | Éviter helm hooks pré/post — utiliser ArgoCD sync waves |

## 5. Security Policies

### 5.1 Kyverno (obligatoire, enforcement)
- `require-labels` : tous les labels nationaux présents
- `require-resources` : limits/requests CPU+memory
- `require-ro-rootfs` : readOnlyRootFilesystem: true
- `require-run-as-non-root` : runAsNonRoot + runAsUser > 0
- `require-drop-all` : capabilities.drop: [ALL]
- `restrict-image-registries` : uniquement `registry.interne.snisid.gouv.local/*`
- `require-image-digest` : tag digest obligatoire
- `block-latest` : interdiction tag `latest`

### 5.2 Network Policies
- **Default deny** sur tous les namespaces
- Chaque service doit exposer explicitement ses ALLOW via CiliumNetworkPolicy
- Pas de `0.0.0.0/0` en egress sans justification IGC
- Istio Ingress Gateway est le seul point d'entrée North-South autorisé

## 6. Logging & Observability Standards

### 6.1 Format JSON (ECS — Elastic Common Schema)
```json
{
  "@timestamp": "2026-05-25T14:30:00.000Z",
  "log.level": "info",
  "message": "Citoyen enrollé avec succès",
  "service.name": "snisid-core-api",
  "service.version": "v4.2.1",
  "snisid.gov.tier": "tier-1",
  "snisid.gov.region": "core",
  "snisid.gov.trace_id": "abc123-def456",
  "snisid.gov.tenant": "national",
  "http.request.method": "POST",
  "http.response.status_code": 201,
  "source.ip": "10.1.30.45",
  "user.id": "agent-1234",
  "event.action": "citizen_enrollment",
  "event.outcome": "success",
  "citizen.id_hash": "sha256:abc..."  # Jamais PII en clair dans les logs
}
```

### 6.2 Rétention
| Signal | Hot | Warm | Cold (air-gapped) |
|--------|-----|------|-------------------|
| Metrics | 2 ans Prometheus/Thanos | — | — |
| Logs | 30 jours Loki | 1 an S3 Ceph | 7 ans bande LTO |
| Traces | 30 jours Jaeger | 90 jours S3 | — |
| Audit K8s | 1 an | — | 7 ans bande LTO |
| Audit Vault | 2 ans | — | 10 ans bande LTO |

---

*Standards validés par IGC. Toute exception nécessite dérogation signée DG + IGC.*
