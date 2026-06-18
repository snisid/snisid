# Vault Policy - Fraud Engine
path "transit/encrypt/fraud-model" {
  capabilities = ["create", "update"]
}
path "transit/decrypt/fraud-model" {
  capabilities = ["create", "update"]
}
path "database/creds/fraud-role" {
  capabilities = ["read"]
}
path "kv-v2/data/fraud/*" {
  capabilities = ["read", "list"]
}
path "pki/issue/fraud-service" {
  capabilities = ["create", "update"]
}
