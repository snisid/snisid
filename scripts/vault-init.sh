#!/bin/bash
# SNISID Vault Initialization Script
# Usage: ./scripts/vault-init.sh [dev|prod]
set -euo pipefail

MODE="${1:-dev}"
VAULT_ADDR="${VAULT_ADDR:-http://127.0.0.1:8200}"
VAULT_ROOT_TOKEN=""
POLICIES_DIR="vault/policies"

echo "=== SNISID Vault Initialization ==="
echo "Mode: $MODE"
echo "Vault Address: $VAULT_ADDR"

# Step 1: Initialize Vault
if vault status &>/dev/null; then
  echo "[OK] Vault is already initialized and unsealed"
else
  echo "[INIT] Initializing Vault..."
  INIT_OUTPUT=$(vault operator init -key-shares=5 -key-threshold=3 -format=json)
  echo "$INIT_OUTPUT" > vault/init-keys.json
  echo "[SAVE] Keys saved to vault/init-keys.json"

  VAULT_ROOT_TOKEN=$(echo "$INIT_OUTPUT" | jq -r '.root_token')
  UNSEAL_KEYS=$(echo "$INIT_OUTPUT" | jq -r '.unseal_keys_b64[]')

  # Unseal with 3 of 5 keys
  echo "[UNSEAL] Unsealing Vault..."
  for i in 1 2 3; do
    KEY=$(echo "$UNSEAL_KEYS" | sed -n "${i}p")
    vault operator unseal "$KEY"
  done
fi

# Step 2: Login
if [ -z "$VAULT_ROOT_TOKEN" ]; then
  VAULT_ROOT_TOKEN=$(jq -r '.root_token' vault/init-keys.json 2>/dev/null || echo "")
fi
vault login "$VAULT_ROOT_TOKEN"

# Step 3: Enable secret engines
echo "[ENGINE] Enabling secret engines..."
vault secrets enable -path=kv-v2 kv-v2 2>/dev/null || echo "  kv-v2 already enabled"
vault secrets enable -path=transit transit 2>/dev/null || echo "  transit already enabled"
vault secrets enable -path=pki pki 2>/dev/null || echo "  pki already enabled"
vault secrets enable database 2>/dev/null || echo "  database already enabled"

# Step 4: Enable Kubernetes auth (for prod)
if [ "$MODE" = "prod" ]; then
  vault auth enable kubernetes 2>/dev/null || echo "  kubernetes auth already enabled"

  K8S_HOST="${K8S_HOST:-https://kubernetes.default.svc:443}"
  K8S_CA_CERT=$(cat /var/run/secrets/kubernetes.io/serviceaccount/ca.crt 2>/dev/null || echo "")
  K8S_TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token 2>/dev/null || echo "")

  vault write auth/kubernetes/config \
    token_reviewer_jwt="$K8S_TOKEN" \
    kubernetes_host="$K8S_HOST" \
    kubernetes_ca_cert="$K8S_CA_CERT"

  echo "[AUTH] Kubernetes auth configured"
fi

# Step 5: Write policies
echo "[POLICY] Writing Vault policies..."
for policy_file in "$POLICIES_DIR"/*.hcl; do
  policy_name=$(basename "$policy_file" .hcl)
  vault policy write "$policy_name" "$policy_file"
  echo "  -> $policy_name"
done

# Step 6: Generate encryption keys
echo "[KEYS] Generating encryption keys..."
for key in identity-ssn biometric-template biometric-image fraud-model; do
  vault write -f transit/keys/$key 2>/dev/null || echo "  key $key already exists"
done

# Step 7: Write initial secrets
echo "[SECRETS] Writing initial configuration secrets..."
vault kv put kv-v2/shared/jwt \
  secret="change-me-in-production-rotate-immediately" \
  ttl="15m"

vault kv put kv-v2/shared/postgres \
  host="postgres-primary.snisid.svc.cluster.local" \
  port="5432" \
  dbname="snisid" \
  user="snisid-app" \
  password="$(openssl rand -base64 32)"

vault kv put kv-v2/shared/redis \
  host="redis.snisid.svc.cluster.local" \
  port="6379" \
  password="$(openssl rand -base64 32)"

vault kv put kv-v2/shared/kafka \
  brokers="kafka-0.kafka-headless.snisid.svc.cluster.local:9092" \
  username="snisid" \
  password="$(openssl rand -base64 32)"

# Step 8: Configure database roles (for dynamic secrets)
echo "[DB] Configuring database secret engine..."
vault write database/config/postgres-snisid \
  plugin_name=postgresql-database-plugin \
  allowed_roles="identity-role,biometrics-role,fraud-role,audit-role,gateway-role" \
  connection_url="postgresql://{{username}}:{{password}}@postgres-primary.snisid.svc.cluster.local:5432/snisid?sslmode=verify-full" \
  username="vault-db-admin" \
  password="$(openssl rand -base64 32)"

for role in identity biometrics fraud audit gateway; do
  vault write database/roles/${role}-role \
    db_name=postgres-snisid \
    creation_statements="CREATE USER \"{{name}}\" WITH PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; GRANT ${role}_role TO \"{{name}}\";" \
    default_ttl="1h" \
    max_ttl="24h"
done

echo "=== Vault initialization complete ==="
echo "Root Token: $VAULT_ROOT_TOKEN (store securely!)"
echo "Keys file: vault/init-keys.json"
