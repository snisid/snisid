# Vault Policy - Interoperability
path "kv-v2/data/interop/*" {
  capabilities = ["read", "list"]
}
path "kv-v2/metadata/interop/*" {
  capabilities = ["read", "list"]
}
path "pki/issue/interop-*" {
  capabilities = ["create", "update"]
}
path "transit/encrypt/interop-*" {
  capabilities = ["create", "update"]
}
path "transit/decrypt/interop-*" {
  capabilities = ["create", "update"]
}
