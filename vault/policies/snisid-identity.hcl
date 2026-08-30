# Vault Policy - Identity Service
path "transit/encrypt/identity-ssn" {
  capabilities = ["create", "update"]
}
path "transit/decrypt/identity-ssn" {
  capabilities = ["create", "update"]
}
path "database/creds/identity-role" {
  capabilities = ["read"]
}
path "kv-v2/data/identity/*" {
  capabilities = ["read", "list"]
}
path "kv-v2/metadata/identity/*" {
  capabilities = ["read", "list"]
}
path "pki/issue/identity-service" {
  capabilities = ["create", "update"]
}
