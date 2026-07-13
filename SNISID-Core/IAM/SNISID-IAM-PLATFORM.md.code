# ============================================================
# SNISID — IAM Platform
# Keycloak + HashiCorp Vault + Teleport + OPA
# Document ID: SNISID-IAM-001
# Version: 1.0.0
# ============================================================

---
## 1. KEYCLOAK — Realm SNISID

Keycloak est le **fournisseur d'identité centralisé** (IdP) de la plateforme SNISID. Il implémente OAuth 2.1 + OpenID Connect 1.0.

### 1.1 Realm Configuration

```json
{
  "realm": "snisid",
  "displayName": "SNISID — Système National d'Identification",
  "displayNameHtml": "<img src='/snisid-logo.png'/> SNISID",
  "enabled": true,
  "sslRequired": "all",
  "registrationAllowed": false,
  "loginWithEmailAllowed": false,
  "duplicateEmailsAllowed": false,
  "resetPasswordAllowed": false,
  "editUsernameAllowed": false,
  "bruteForceProtected": true,
  "permanentLockout": false,
  "maxFailureWaitSeconds": 900,
  "minimumQuickLoginWaitSeconds": 60,
  "waitIncrementSeconds": 60,
  "quickLoginCheckMilliSeconds": 1000,
  "maxDeltaTimeSeconds": 43200,
  "failureFactor": 5,
  "defaultRoles": ["offline_access"],
  "requiredCredentials": ["password"],
  "passwordPolicy": "length(16) and upperCase(2) and lowerCase(2) and digits(2) and specialChars(2) and notUsername and notEmail and passwordHistory(12) and forceExpiredPasswordChange(90)",
  "otpPolicyType": "totp",
  "otpPolicyAlgorithm": "HmacSHA256",
  "otpPolicyDigits": 6,
  "otpPolicyPeriod": 30,
  "webAuthnPolicySignatureAlgorithms": ["RS256", "ES256"],
  "webAuthnPolicyAttestationConveyancePreference": "direct",
  "webAuthnPolicyAuthenticatorAttachment": "cross-platform",
  "webAuthnPolicyUserVerificationRequirement": "required",
  "accessTokenLifespan": 300,
  "accessTokenLifespanForImplicitFlow": 900,
  "ssoSessionIdleTimeout": 1800,
  "ssoSessionMaxLifespan": 36000,
  "offlineSessionIdleTimeout": 2592000,
  "accessCodeLifespan": 60,
  "smtpServer": {
    "host": "smtp.snisid.gov.ht",
    "port": "587",
    "starttls": true
  },
  "eventsEnabled": true,
  "eventsListeners": ["jboss-logging", "kafka-event-listener"],
  "enabledEventTypes": ["LOGIN", "LOGOUT", "REGISTER", "CLIENT_LOGIN", "ADMIN"],
  "adminEventsEnabled": true,
  "adminEventsDetailsEnabled": true
}
```

### 1.2 Clients (Agences)

```json
{
  "clients": [
    {
      "clientId": "identity-service",
      "name": "SNISID Identity Service",
      "enabled": true,
      "clientAuthenticatorType": "client-jwt",
      "redirectUris": ["https://api.snisid.gov.ht/v1/callback"],
      "standardFlowEnabled": false,
      "directAccessGrantsEnabled": false,
      "serviceAccountsEnabled": true,
      "authorizationServicesEnabled": true,
      "protocol": "openid-connect",
      "defaultScopes": ["identity:read", "identity:write"],
      "attributes": {
        "use.jwks.url": "true",
        "jwks.url": "https://identity-service.snisid-identity.svc.cluster.local:8443/jwks"
      }
    },
    {
      "clientId": "oni-portal",
      "name": "ONI — Office National d'Identification",
      "enabled": true,
      "clientAuthenticatorType": "client-jwt",
      "serviceAccountsEnabled": true,
      "authorizationServicesEnabled": true,
      "defaultScopes": ["identity:read", "identity:write", "biometric:enroll"]
    },
    {
      "clientId": "pnh-enforcement",
      "name": "PNH — Police Nationale d'Haïti",
      "enabled": true,
      "serviceAccountsEnabled": true,
      "defaultScopes": ["identity:verify", "biometric:verify"]
    },
    {
      "clientId": "dgi-fiscal",
      "name": "DGI — Direction Générale des Impôts",
      "enabled": true,
      "serviceAccountsEnabled": true,
      "defaultScopes": ["identity:read"]
    },
    {
      "clientId": "cep-electoral",
      "name": "CEP — Conseil Électoral Permanent",
      "enabled": true,
      "serviceAccountsEnabled": true,
      "defaultScopes": ["identity:read", "civil:read"]
    }
  ]
}
```

### 1.3 Roles Hiérarchiques SNISID

```
SNISID Realm Roles:
├── snisid_admin                    → Accès complet (AND Direction)
│   ├── Hérite: identity:*, biometric:*, civil:*, fraud:*, audit:*
│
├── national_ciso                   → Sécurité (logs, alerts, revocation)
│   ├── Hérite: audit:read, security:*, identity:suspend
│
├── enrollment_supervisor           → Chef d'équipe ONI
│   ├── Hérite: enrollment:supervise, identity:read, agent:manage
│
├── enrollment_officer              → Agent terrain ONI
│   ├── Hérite: identity:write, biometric:enroll, civil:write
│   └── Restriction: commune_assignee only (ABAC OPA)
│
├── oec_officer                     → Officier d'État Civil
│   ├── Hérite: civil:write, identity:read
│   └── Restriction: commune_assignee only
│
├── agency_verifier                 → Agent vérificateur (DGI, MSPP, etc.)
│   └── Hérite: identity:verify (niveau 1 seul)
│
├── pnh_officer                     → Officier PNH vérification
│   └── Hérite: identity:verify (niveau 2), biometric:verify (niveau 3)
│
└── dcpj_investigator               → DCPJ (Investigation fraude)
    └── Hérite: fraud:investigate, audit:read, identity:suspend
```

---

## 2. HASHICORP VAULT — Secrets Management

### 2.1 Architecture Vault

```yaml
# vault-values.yaml (Helm)
vault:
  global:
    enabled: true
    tlsDisable: false

  server:
    ha:
      enabled: true
      replicas: 3
      raft:
        enabled: true
        setNodeId: true
        config: |
          ui = true
          listener "tcp" {
            tls_disable = 0
            address = "[::]:8200"
            tls_cert_file = "/vault/userconfig/vault-tls/vault.crt"
            tls_key_file = "/vault/userconfig/vault-tls/vault.key"
            tls_client_ca_file = "/vault/userconfig/vault-tls/ca.crt"
          }
          storage "raft" {
            path = "/vault/data"
            retry_join {
              leader_tls_servername = "vault-0.vault-internal"
              leader_api_addr = "https://vault-0.vault-internal:8200"
            }
            retry_join {
              leader_tls_servername = "vault-1.vault-internal"
              leader_api_addr = "https://vault-1.vault-internal:8200"
            }
            retry_join {
              leader_tls_servername = "vault-2.vault-internal"
              leader_api_addr = "https://vault-2.vault-internal:8200"
            }
          }
          seal "pkcs11" {
            lib = "/usr/lib/softhsm/libsofthsm2.so"
            slot = "0"
            pin = env:VAULT_HSM_PIN
            key_label = "snisid-vault-seal-key"
            hmac_key_label = "snisid-vault-hmac-key"
          }
          telemetry {
            prometheus_retention_time = "30s"
            disable_hostname = true
          }
          log_level = "info"
          log_format = "json"
          api_addr = "https://vault.snisid-security.svc.cluster.local:8200"
          cluster_addr = "https://$(POD_IP):8201"
          ui = true
          disable_mlock = true

    resources:
      requests:
        memory: 256Mi
        cpu: 250m
      limits:
        memory: 512Mi
        cpu: 500m

    auditStorage:
      enabled: true
      size: 50Gi
      storageClass: ceph-rbd-xfs
```

### 2.2 Secrets Engine Configuration

```hcl
# vault/setup/secrets-engines.hcl

# PKI — Certificates dynamiques
vault secrets enable -path=pki pki
vault secrets tune -max-lease-ttl=87600h pki

# Database — Credentials dynamiques (PostgreSQL)
vault secrets enable -path=database database
vault write database/config/snisid-identity \
    plugin_name=postgresql-database-plugin \
    allowed_roles="identity-service-role,registry-service-role,audit-service-role" \
    connection_url="postgresql://{{username}}:{{password}}@cockroachdb.snisid-databases.svc.cluster.local:5432/snisid" \
    max_open_connections=5 \
    max_connection_lifetime="5m" \
    username="vault-rotation-user" \
    password="{{.InitialPassword}}"

# Rôle identity-service (minimal privileges)
vault write database/roles/identity-service-role \
    db_name=snisid-identity \
    creation_statements="CREATE ROLE '{{name}}' WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}';
      GRANT USAGE ON SCHEMA snisid_identity TO '{{name}}';
      GRANT SELECT, INSERT, UPDATE ON snisid_identity.citizens TO '{{name}}';
      GRANT INSERT ON snisid_identity.identity_events TO '{{name}}';
      SET ROLE '{{name}}';" \
    revocation_statements="REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA snisid_identity FROM '{{name}}'; DROP ROLE '{{name}}';" \
    default_ttl="1h" \
    max_ttl="8h"

# KV v2 — Secrets d'application
vault secrets enable -path=secret -version=2 kv

# Transit — Chiffrement des données sensibles
vault secrets enable transit
vault write transit/keys/citizen-data type=aes256-gcm96 deletion_allowed=false
vault write transit/keys/biometric-template type=aes256-gcm96 deletion_allowed=false

# AWS — Credentials cloud (si applicable)
# vault secrets enable aws
```

### 2.3 Policies Vault

```hcl
# vault/policies/identity-service.hcl
path "database/creds/identity-service-role" {
  capabilities = ["read"]
}
path "secret/data/identity-service/*" {
  capabilities = ["read"]
}
path "transit/encrypt/citizen-data" {
  capabilities = ["update"]
}
path "transit/decrypt/citizen-data" {
  capabilities = ["update"]
}
path "pki/issue/identity-service-cert" {
  capabilities = ["create", "update"]
}

# vault/policies/biometric-service.hcl
path "database/creds/biometric-service-role" {
  capabilities = ["read"]
}
path "secret/data/biometric-service/*" {
  capabilities = ["read"]
}
path "transit/encrypt/biometric-template" {
  capabilities = ["update"]
}
path "transit/decrypt/biometric-template" {
  capabilities = ["update"]
}

# vault/policies/emergency-admin.hcl
# Accès bris-de-glace en cas de crise
path "*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
# Utilisé uniquement si identités des 3 administrateurs de sécurité sont disponibles (Shamir 3/5)
```

### 2.4 Kubernetes Auth Method

```hcl
# Activer l'authentification Kubernetes
vault auth enable kubernetes

vault write auth/kubernetes/config \
    kubernetes_host="https://kubernetes.default.svc.cluster.local:443" \
    token_reviewer_jwt=@/var/run/secrets/kubernetes.io/serviceaccount/token \
    kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt \
    issuer="https://kubernetes.default.svc.cluster.local"

# Binding: identity-service ServiceAccount → policy
vault write auth/kubernetes/role/identity-service \
    bound_service_account_names=identity-service \
    bound_service_account_namespaces=snisid-identity \
    policies=identity-service \
    ttl=1h

vault write auth/kubernetes/role/biometric-service \
    bound_service_account_names=biometric-service \
    bound_service_account_namespaces=snisid-biometrics \
    policies=biometric-service \
    ttl=1h

vault write auth/kubernetes/role/civil-registry \
    bound_service_account_names=civil-registry-service \
    bound_service_account_namespaces=snisid-civil-registry \
    policies=civil-registry-service \
    ttl=1h
```

---

## 3. TELEPORT — Zero Trust Access Platform

Teleport fournit l'accès **bastion Zero Trust** pour les administrateurs SNISID.

```yaml
# teleport-values.yaml
teleport:
  enabled: true
  version: 15
  
  auth:
    enabled: true
    cluster_name: "snisid.gov.ht"
    public_addr: "teleport.snisid.gov.ht:443"
    authentication:
      second_factor: otp    # TOTP obligatoire
      type: local
    audit_log:
      storage:
        type: dynamodb    # ou PostgreSQL
      regions: [snisid-prod]
    session_recording: node-sync    # Sessions enregistrées
    
  proxy:
    enabled: true
    public_addr: "teleport.snisid.gov.ht:443"
    ssh_public_addr: "ssh.teleport.snisid.gov.ht:3023"
    kube_public_addr: "k8s.teleport.snisid.gov.ht:3026"
    
  teleportConfig: |
    teleport:
      log:
        severity: INFO
        format:
          output: json
    auth_service:
      authentication:
        second_factor: otp
        webauthn:
          rp_id: teleport.snisid.gov.ht
      session_recording: node-sync
      audit_events_uri: ["kafka://kafka.snisid-event-bus:9092?topic=snisid.audit.teleport"]
    kubernetes_service:
      enabled: true
      kubeconfig_file: ""
      labels:
        env: production
        cluster: snisid-prod
```

---

## 4. OPA GATEKEEPER — Admission Control

```yaml
# gatekeeper-config.yaml
apiVersion: config.gatekeeper.sh/v1alpha1
kind: Config
metadata:
  name: config
  namespace: gatekeeper-system
spec:
  sync:
    syncOnly:
    - group: ""
      version: v1
      kind: Namespace
    - group: ""
      version: v1
      kind: Pod
  validation:
    traces:
    - user: "k8s-admin"
      kind:
        group: "*"
        version: "*"
        kind: "*"

---
# Constraint Template: Require non-root containers
apiVersion: templates.gatekeeper.sh/v1
kind: ConstraintTemplate
metadata:
  name: snisidnorootcontainers
spec:
  crd:
    spec:
      names:
        kind: SNISIDNoRootContainers
  targets:
    - target: admission.k8s.gatekeeper.sh
      rego: |
        package snisid.norootcontainers
        
        violation[{"msg": msg}] {
          c := input.review.object.spec.containers[_]
          not c.securityContext.runAsNonRoot
          msg := sprintf("Container '%v' must set runAsNonRoot=true", [c.name])
        }

        violation[{"msg": msg}] {
          c := input.review.object.spec.containers[_]
          c.securityContext.runAsUser == 0
          msg := sprintf("Container '%v' must not run as root (uid=0)", [c.name])
        }

---
# Constraint: Apply to SNISID namespaces
apiVersion: constraints.gatekeeper.sh/v1beta1
kind: SNISIDNoRootContainers
metadata:
  name: no-root-containers-snisid
spec:
  enforcementAction: deny
  match:
    namespaces:
    - snisid-identity
    - snisid-biometrics
    - snisid-civil-registry
    - snisid-databases
    - snisid-security
    kinds:
    - apiGroups: [""]
      kinds: ["Pod"]

---
# Constraint Template: Require resource limits
apiVersion: templates.gatekeeper.sh/v1
kind: ConstraintTemplate
metadata:
  name: snisidrequirerlimits
spec:
  crd:
    spec:
      names:
        kind: SNISIDRequireResourceLimits
  targets:
    - target: admission.k8s.gatekeeper.sh
      rego: |
        package snisid.requirelimits
        
        violation[{"msg": msg}] {
          c := input.review.object.spec.containers[_]
          not c.resources.limits.memory
          msg := sprintf("Container '%v' must have memory limit", [c.name])
        }
        
        violation[{"msg": msg}] {
          c := input.review.object.spec.containers[_]
          not c.resources.limits.cpu
          msg := sprintf("Container '%v' must have CPU limit", [c.name])
        }

---
# Constraint: Required image from approved registry
apiVersion: templates.gatekeeper.sh/v1
kind: ConstraintTemplate
metadata:
  name: snisidapprovedregistry
spec:
  crd:
    spec:
      names:
        kind: SNISIDApprovedRegistry
  targets:
    - target: admission.k8s.gatekeeper.sh
      rego: |
        package snisid.approvedregistry
        
        approved_registries = {
          "harbor.snisid.gov.ht",
          "registry.k8s.io",
          "gcr.io/google_containers"
        }
        
        violation[{"msg": msg}] {
          c := input.review.object.spec.containers[_]
          not any_approved(c.image)
          msg := sprintf("Container '%v' uses unapproved image: %v", [c.name, c.image])
        }
        
        any_approved(image) {
          registry := approved_registries[_]
          startswith(image, registry)
        }

---
apiVersion: constraints.gatekeeper.sh/v1beta1
kind: SNISIDApprovedRegistry
metadata:
  name: approved-registry-enforcement
spec:
  enforcementAction: deny
  match:
    namespaces:
    - snisid-identity
    - snisid-biometrics
    - snisid-civil-registry
    - snisid-databases
    - snisid-security
    - snisid-event-bus
    kinds:
    - apiGroups: [""]
      kinds: ["Pod"]
```

---

## 5. KYVERNO — Policy Engine

```yaml
# kyverno-policies/require-signed-images.yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-signed-images
  annotations:
    policies.kyverno.io/title: Require Signed Images
    policies.kyverno.io/description: >
      Toutes les images déployées dans les namespaces SNISID
      doivent être signées avec Cosign via Harbor.
spec:
  validationFailureAction: Enforce
  background: false
  rules:
  - name: check-image-signature
    match:
      any:
      - resources:
          kinds: [Pod]
          namespaces:
          - snisid-identity
          - snisid-biometrics
          - snisid-civil-registry
          - snisid-security
    verifyImages:
    - imageReferences:
      - "harbor.snisid.gov.ht/*"
      attestors:
      - count: 1
        entries:
        - keys:
            publicKeys: |-
              -----BEGIN PUBLIC KEY-----
              MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE... (SNISID Cosign Public Key)
              -----END PUBLIC KEY-----

---
# kyverno-policies/disallow-privileged.yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: disallow-privileged-containers
spec:
  validationFailureAction: Enforce
  rules:
  - name: no-privileged
    match:
      any:
      - resources:
          kinds: [Pod]
          namespaces:
          - snisid-*
    validate:
      message: "Les conteneurs privilégiés sont interdits dans SNISID"
      pattern:
        spec:
          containers:
          - securityContext:
              privileged: "false | null"
          =(initContainers):
          - securityContext:
              privileged: "false | null"

---
# kyverno-policies/require-network-policy.yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-namespace-network-policy
spec:
  validationFailureAction: Audit    # Audit d'abord, puis Enforce après validation
  rules:
  - name: check-network-policy-exists
    match:
      any:
      - resources:
          kinds: [Namespace]
          selector:
            matchLabels:
              snisid.gov.ht/tier: critical
    validate:
      message: "Les namespaces SNISID critiques doivent avoir une NetworkPolicy default-deny"
      deny:
        conditions:
        - key: "{{ request.object.metadata.name }}"
          operator: NotIn
          value: "{{ networkpolicies.items[?name=='default-deny-all'].metadata.namespace }}"
```

---

*Document ID : SNISID-IAM-001 v1.0.0 — Mai 2026*  
*Keycloak + HashiCorp Vault + Teleport + OPA Gatekeeper + Kyverno*  
*Approuvé par : CISO National | DG-AND | Platform Security Lead*
