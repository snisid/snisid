# Vault Policy - SNISID Administrator
path "sys/mounts" {
  capabilities = ["read", "list"]
}
path "sys/policies/*" {
  capabilities = ["read", "list"]
}
path "identity/*" {
  capabilities = ["read", "list"]
}
path "sys/health" {
  capabilities = ["read"]
}
path "sys/seal-status" {
  capabilities = ["read"]
}
path "kv-v2/data/shared/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
path "kv-v2/metadata/shared/*" {
  capabilities = ["read", "delete", "list"]
}
path "transit/keys/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
path "pki/*" {
  capabilities = ["read", "list"]
}
