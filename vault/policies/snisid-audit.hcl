# Vault Policy - Audit Service
path "transit/verify/audit-trail" {
  capabilities = ["create", "update"]
}
path "transit/hash/audit-event" {
  capabilities = ["create", "update"]
}
path "database/creds/audit-role" {
  capabilities = ["read"]
}
path "kv-v2/data/audit/*" {
  capabilities = ["read", "list"]
}
path "pki/issue/audit-service" {
  capabilities = ["create", "update"]
}
