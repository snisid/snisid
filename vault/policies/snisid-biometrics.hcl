# Vault Policy - Biometrics Service
path "transit/encrypt/biometric-template" {
  capabilities = ["create", "update"]
}
path "transit/decrypt/biometric-template" {
  capabilities = ["create", "update"]
}
path "transit/encrypt/biometric-image" {
  capabilities = ["create", "update"]
}
path "transit/decrypt/biometric-image" {
  capabilities = ["create", "update"]
}
path "database/creds/biometrics-role" {
  capabilities = ["read"]
}
path "kv-v2/data/biometrics/*" {
  capabilities = ["read", "list"]
}
path "pki/issue/biometrics-service" {
  capabilities = ["create", "update"]
}
