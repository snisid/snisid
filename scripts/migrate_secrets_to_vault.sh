#!/bin/bash
set -euo pipefail
VAULT_ADDR="${VAULT_ADDR:-https://vault.snisid.gouv.ht}"
echo "=== Migration Secrets SNISID vers Vault ==="
 
# Auth Vault via Kubernetes ServiceAccount
SA_TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
VAULT_TOKEN=$(curl -sf \
  --data '{"role":"snisid-migration","jwt":"'"$SA_TOKEN"'"}' \
  "$VAULT_ADDR/v1/auth/kubernetes/login" | jq -r '.auth.client_token')
 
vault_write() {
  curl -sf -H "X-Vault-Token: $VAULT_TOKEN" \
    -X POST "$VAULT_ADDR/v1/$1" --data "$2" > /dev/null
  echo "  OK Migre: $1"
}
 
# Phase 13 - DB Signing Service
vault_write "secret/data/snisid/phase13/database" \
  '{"data":{"dsn":"postgres://signing_svc:CHANGE_ME@postgres-signing:5432/snisid_signing"}}'
 
# Phase 15 - Grafana
vault_write "secret/data/snisid/phase15/grafana" \
  '{"data":{"admin_password":"GENERE_VIA_VAULT_GENERATE"}}'
 
# Phase 18 - Intelligence Nationale
vault_write "secret/data/snisid/phase18/mlflow" \
  '{"data":{"tracking_uri":"postgresql://mlflow:CHANGE_ME@postgres-ml:5432/mlflow"}}'
vault_write "secret/data/snisid/phase18/superset" \
  '{"data":{"admin_password":"GENERE_VIA_VAULT","secret_key":"GENERE_VIA_VAULT"}}'
 
echo "OK Migration terminee. Supprimer les anciens docker-compose secrets!"
