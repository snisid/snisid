#!/bin/bash
# SNISID PKI Configuration Script
# Automates the configuration of the Vault Intermediate CA Roles and Revocation (OCSP/CRL)

set -e
export VAULT_ADDR="https://vault.snisid-security.svc.cluster.local:8200"

echo "[*] Tuning PKI Secrets Engine..."
vault secrets enable -path=pki_int pki
vault secrets tune -max-lease-ttl=87600h pki_int # 10 years

echo "[*] Configuring Revocation Endpoints (CRL / OCSP)..."
vault write pki_int/config/urls \
    issuing_certificates="$VAULT_ADDR/v1/pki_int/ca" \
    crl_distribution_points="$VAULT_ADDR/v1/pki_int/crl" \
    ocsp_servers="$VAULT_ADDR/v1/pki_int/ocsp"

echo "[*] Creating Role: Server TLS (snisid.local)..."
vault write pki_int/roles/snisid-dot-local \
    allowed_domains="snisid-internal.svc.cluster.local,snisid.local" \
    allow_subdomains=true \
    max_ttl="720h" \
    require_cn=false \
    generate_lease=true

echo "[*] Creating Role: Client mTLS (SPIFFE / Istio)..."
vault write pki_int/roles/istio-mtls \
    allowed_uri_sans="spiffe://cluster.local/ns/*" \
    client_flag=true \
    server_flag=true \
    max_ttl="72h"

echo "[*] Creating Role: Code Signing..."
vault write pki_int/roles/code-signing \
    allow_any_name=true \
    client_flag=false \
    server_flag=false \
    code_signing_flag=true \
    max_ttl="8760h"

echo "[*] Configuring Kubernetes Authentication for cert-manager..."
vault auth enable kubernetes
vault write auth/kubernetes/config \
    kubernetes_host="https://kubernetes.default.svc.cluster.local:443"

vault write auth/kubernetes/role/cert-manager \
    bound_service_account_names=cert-manager \
    bound_service_account_namespaces=cert-manager \
    policies=pki-issuer-policy \
    ttl=1h

echo "[*] PKI Configuration Complete."
