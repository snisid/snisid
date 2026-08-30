# SNISID HashiCorp Vault Server Configuration (Intermediate CA)
# Integrating with the online Thales/Utimaco Network HSM via PKCS#11

storage "raft" {
  path    = "/vault/data"
  node_id = "snisid-vault-node-1"
}

listener "tcp" {
  address       = "0.0.0.0:8200"
  tls_cert_file = "/vault/tls/vault-cert.pem"
  tls_key_file  = "/vault/tls/vault-key.pem"
  
  # TLS 1.3 enforced natively
  tls_min_version = "tls13"
}

# Requirement: HSM Integration (PKCS#11)
# Auto-unseal Vault using the HSM and enable Seal Wrap to cryptographically
# protect all CA private keys stored in the Raft backend.
seal "pkcs11" {
  lib            = "/usr/lib/libcs_pkcs11_R2.so" # Vendor specific PKCS#11 library
  slot           = "0"
  pin            = "SNISID_HSM_PIN_INJECTED_AT_RUNTIME"
  key_label      = "vault-master-key"
  hmac_key_label = "vault-hmac-key"
}

api_addr = "https://vault.snisid-security.svc.cluster.local:8200"
cluster_addr = "https://snisid-vault-node-1:8201"
ui = true

# Audit logging to local file, which Filebeat/Logstash picks up for the SIEM
audit {
  type = "file"
  path = "/vault/logs/vault-audit.log"
}
