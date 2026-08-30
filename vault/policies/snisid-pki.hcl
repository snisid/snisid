# Vault Policy - PKI Management
path "pki/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
path "pki/cert/ca" {
  capabilities = ["read"]
}
path "pki/issue/*" {
  capabilities = ["create", "update"]
}
path "pki/revoke" {
  capabilities = ["create", "update"]
}
path "pki/crl" {
  capabilities = ["read"]
}
path "pki/ocsp" {
  capabilities = ["read"]
}
path "pki/config/*" {
  capabilities = ["create", "read", "update"]
}
path "pki/roles/*" {
  capabilities = ["create", "read", "update", "delete"]
}
