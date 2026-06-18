# Vault Policy - API Gateway
path "transit/verify/jwt-*" {
  capabilities = ["create", "update"]
}
path "transit/decrypt/session" {
  capabilities = ["create", "update"]
}
path "database/creds/gateway-role" {
  capabilities = ["read"]
}
path "kv-v2/data/gateway/*" {
  capabilities = ["read", "list"]
}
path "pki/issue/gateway-service" {
  capabilities = ["create", "update"]
}
path "auth/token/lookup-self" {
  capabilities = ["read"]
}
