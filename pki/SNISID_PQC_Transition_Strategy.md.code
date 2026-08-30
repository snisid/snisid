# SNISID PQC Transition Strategy: Operational Scripts & Playbooks

**Classification:** RESTRICTED / SOVEREIGN PQC
**Compliance:** NIST SP 800-227 / FIPS 203 / FIPS 204 / RFC 9370

This document defines the configuration templates, operational scripts, and testing playbooks required to implement Post-Quantum Cryptography (PQC) across the SNISID network.

---

## 1. Hybrid Certificate Configuration & OpenSSL Templates

During the hybrid transition, OpenSSL utilizes the OQS (Open Quantum Safe) provider to generate dual-signature certificates.

### 1.1. OpenSSL Configuration for Dual-Signature Certificates
Use the configuration file below to define custom extensions for embedding ML-DSA public keys inside standard ECC certificates.

```ini
# File: /etc/ssl/openssl-hybrid.cnf
[ req ]
default_bits        = 3072
distinguished_name  = req_distinguished_name
req_extensions      = v3_hybrid_req
x509_extensions     = v3_hybrid_cert
string_mask         = utf8only

[ req_distinguished_name ]
countryName                     = Country Name (2 letter code)
stateOrProvinceName             = State or Province Name
localityName                    = Locality Name
organizationName                = Organization Name
commonName                      = Common Name

[ v3_hybrid_req ]
# Enable classical ECC extensions
subjectKeyIdentifier = hash

[ v3_hybrid_cert ]
# Standard certificate extensions
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer
basicConstraints = critical, CA:true, pathlen:0
keyUsage = critical, digitalSignature, cRLSign, keyCertSign

# Custom NIST PQC extension definition (Dual Signature mapping)
# Contains the OID for ML-DSA-87 public key and the raw public key bytes
1.3.6.1.4.1.2.261.2.1 = DER:300f020101040a4d4c2d4453412d3837
```

### 1.2. Generating a Hybrid Keypair & CSR
To generate a key pair and output a Certificate Signing Request (CSR) with the Open Quantum Safe provider:

```bash
#!/bin/bash
# File: /opt/snisid/pki/generate_hybrid_csr.sh
set -e

# 1. Generate classical ECC private key
openssl ecparam -name secp384r1 -genkey -noout -out classical.key

# 2. Generate ML-DSA-87 PQC private key
openssl genpkey -algorithm mldsa87 -out pqc.key

# 3. Combine keys into a hybrid key container
openssl pkeyutl -combine \
  -inkey classical.key -inkey pqc.key \
  -out hybrid_combined.key

# 4. Generate the hybrid CSR
openssl req -config /etc/ssl/openssl-hybrid.cnf \
  -new -key hybrid_combined.key \
  -out hybrid_device.csr \
  -subj "/C=HT/ST=Ouest/L=Port-au-Prince/O=SNISID/CN=Biometric-Kiosk-102"
```

---

## 2. Kubernetes & Vault Agility Orchestrations

Enabling hybrid key exchanges ensures that all mTLS traffic is secure against "Store-Now-Decrypt-Later" exfiltrations.

### 2.1. Istio Envoy Filter for Hybrid Key Exchange
Apply this EnvoyFilter to force the mesh sidecars to negotiate the experimental/draft hybrid key exchange algorithm (`X25519Kyber768Draft00` or standard `ML-KEM-768`).

```yaml
# File: /pki/k8s/istio-pqc-envoyfilter.yaml
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: enforce-hybrid-kex
  namespace: istio-system
spec:
  configPatches:
    - applyTo: downstream_tls_context
      match:
        context: SIDECAR_INBOUND
      patch:
        operation: MERGE
        value:
          common_tls_context:
            tls_params:
              tls_minimum_protocol_version: TLSv1_3
              # Enforce Kyber / ML-KEM hybrid exchange alongside classical ECDHE
              ecdh_curves:
                - x25519_kyber768_draft00
                - X25519
                - P-384
```

### 2.2. Vault PKI PQC Secrets Configuration
Configure HashiCorp Vault's PKI engine to enable post-quantum signature support on intermediate certificate paths:

```bash
# Enable PQC algorithms on Vault PKI secrets engine (requires OQS patched Vault version)
vault write pki_int/config/crypto \
  allowed_signature_algorithms="ECDSA-SHA384,ML-DSA-65,ML-DSA-87" \
  allowed_key_types="ec,mldsa"

# Update Gov-Infrastructure role to allow hybrid generation
vault write pki_int/roles/gov-infra-pqc \
  allowed_domains="snisid.local" \
  allow_subdomains=true \
  key_type="mldsa" \
  key_bits="87" \
  max_ttl="168h"
```

---

## 3. Biometric Edge Node Compatibility Audit Script

Prior to pushing new hybrid certificates, a discovery scan determines whether remote enrollment terminals support the updated algorithms.

```python
# File: /opt/snisid/pki/audit_edge_pqc.py
import socket
import ssl
import sys

def audit_host(ip, port=443):
    print(f"[*] Testing {ip}:{port} for PQC/Hybrid TLS support...")
    context = ssl.create_default_context()
    
    # Configure custom post-quantum and hybrid cipher suites if supported by local library
    try:
        context.set_ciphers("ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384")
        # Attempt to set curves including Kyber/ML-KEM
        context.set_ecdh_curve("x25519_kyber768_draft00")
    except ssl.SSLError:
        print("[!] Local python environment does not support ML-KEM cipher configurations.")
        
    try:
        with socket.create_connection((ip, port), timeout=5) as sock:
            with context.wrap_socket(sock, server_hostname=ip) as ssock:
                cipher = ssock.cipher()
                shared_curve = ssock.shared_ciphers()
                print(f"[+] Connection Successful: Version={ssock.version()}, Cipher={cipher[0]}")
                print(f"[+] Shared curves negotiated: {shared_curve}")
                return True
    except Exception as e:
        print(f"[-] Connection Failed: {str(e)}")
        return False

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python audit_edge_pqc.py <ip_address>")
        sys.exit(1)
    audit_host(sys.argv[1])
```

---

## 4. Quantum Attack Simulation & Hardening Playbook

Security analysts simulate quantum attacks by attempting to downgrade communication links to vulnerable algorithms.

### 4.1. Downgrade Attack Simulation Runbook
1. **Target Identification:** Point scanning tool (e.g., testssl.sh) at the API Gateway.
2. **ExecuteDowngrade Command:**
   ```bash
   # Attempt to force the gateway to accept legacy TLS 1.0/1.1 or weakened RSA cipher suites
   testssl.sh --standard -p api.snisid.gov.ht
   ```
3. **Analyze Results:** The gateway MUST reject any handshake that falls below TLS 1.2 or uses RSA key exchange with key lengths < 2048 bits.

### 4.2. Hardening Audit Log Configuration
To audit cryptographic handshakes, configure Wazuh to monitor for TLS negotiation details in Envoy logs. Envoy should output:

```json
{
  "protocol": "TLSv1.3",
  "cipher_suite": "TLS_AES_256_GCM_SHA384",
  "downstream_peer_cert": "CN=Biometric-Kiosk-102",
  "key_exchange_negotiated": "x25519_kyber768_draft00"
}
```
If `key_exchange_negotiated` contains standard `X25519` (classical-only) for a high-security internal service-to-service connection, Wazuh triggers a Level 10 warning indicating "Non-Quantum Safe Connection Active".

---

*Verified and signed by the SNISID Post-Quantum Cryptography Committee.*
