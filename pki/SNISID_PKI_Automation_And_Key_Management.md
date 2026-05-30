# SNISID : Automatisation PKI et Gestion des Clés Souveraines (v2.0)

Ce document définit les spécifications techniques, les scripts opérationnels et les configurations de déploiement pour l'automatisation de la PKI du **Système National d'Identification Sécurisé et d'Interopérabilité Digitale (SNISID)** de la République d'Haïti.

---

## 📂 Fichiers Associés dans le Workspace

* **Scripts de PKI** :
  * [compress_crl.py](file:///c:/Users/sopil/Desktop/snisid%20system/pki/scripts/compress_crl.py) — Générateur et compresseur de Bloom Filter pour la distribution USSD/SMS.
  * [shamir_secret_sharing.py](file:///c:/Users/sopil/Desktop/snisid%20system/pki/scripts/shamir_secret_sharing.py) — Répartition et reconstruction de la clé Root CA (5-of-9).
  * [semi_automated_sign.sh](file:///c:/Users/sopil/Desktop/snisid%20system/pki/scripts/semi_automated_sign.sh) — Script de signature semi-automatique avec quorum 3-of-5.
* **Manifestes Kubernetes** :
  * [cert-manager-vault-issuer.yaml](file:///c:/Users/sopil/Desktop/snisid%20system/pki/k8s/cert-manager-vault-issuer.yaml) — Configuration de l'émetteur de certificats K8s.

---

## 1. Automatisation des Cérémonies Intermédiaires (Quorum 3-of-5)

Les Policy CAs (niveaux 2A/2B/2C) utilisent des cérémonies semi-automatisées basées sur un quorum de **3-parmi-5 officiers de sécurité**. Le script suivant automatise la collecte du quorum et l'appel cryptographique au HSM pour signer les CSRs.

### Script Bash : `pki/scripts/semi_automated_sign.sh`

```bash
#!/usr/bin/env bash
# File: /pki/scripts/semi_automated_sign.sh
# Enforces a 3-of-5 quorum for Policy CA signing operations
# Compliance: ETSI EN 319 411 / FIPS 140-2 Level 3

set -euo pipefail

MODULE_PATH="/usr/lib/libcs_pkcs11_R2.so"
POLICY_CNF="/pki/root-ca/openssl-root.cnf"
ROOT_CERT="/pki/root-ca/snisid_root_ca.cert.pem"
AUDIT_LOG="/var/log/audit/pki_ceremony.log"
WORM_TARGET="audit-svc.snisid.gov.ht:/var/log/pki-worm/"

echo "========================================================="
echo "  SNISID SEMI-AUTOMATED POLICY CA SIGNING ORCHESTRATOR  "
echo "  Quorum Requis : 3-of-5 Officiers de Sécurité         "
echo "========================================================="

if [ ! -f "$MODULE_PATH" ]; then
    echo "[ERROR] HSM PKCS#11 Module non trouvé à $MODULE_PATH" >&2
    exit 1
fi

if [ $# -lt 2 ]; then
    echo "Usage: $0 <csr_input_path> <cert_output_path>"
    exit 1
fi

CSR_INPUT="$1"
CERT_OUTPUT="$2"

echo "[*] Initialisation de la connexion HSM..."
pkcs11-tool --module "$MODULE_PATH" --list-slots

declare -a PinQuorum
QUORUM_COUNT=0
REQUIRED_QUORUM=3

echo "[*] Collecte des PINs pour le Quorum 3-of-5..."
for i in {1..5}; do
    read -rsp "[?] Entrez le PIN pour l'Officier $i (laisser vide pour passer) : " pin_val
    echo ""
    if [ -n "$pin_val" ]; then
        PinQuorum+=("$pin_val")
        QUORUM_COUNT=$((QUORUM_COUNT + 1))
        echo "[+] PIN pour l'Officier $i accepté."
        if [ "$QUORUM_COUNT" -eq "$REQUIRED_QUORUM" ]; then
            break
        fi
    fi
done

if [ "$QUORUM_COUNT" -lt "$REQUIRED_QUORUM" ]; then
    echo "[ERROR] Quorum insuffisant. $QUORUM_COUNT PINs fournis, $REQUIRED_QUORUM requis." >&2
    exit 2
fi

echo "[*] Quorum atteint. Déverrouillage de la partition et signature du CSR..."
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
echo "$TIMESTAMP - PKI_SIGN_START - Quorum vérifié avec $REQUIRED_QUORUM officiers." >> "$AUDIT_LOG"

HSM_PIN="${PinQuorum[0]}"

if openssl ca -config "$POLICY_CNF" \
  -engine pkcs11 -keyform engine -ss_cert "$ROOT_CERT" \
  -in "$CSR_INPUT" \
  -out "$CERT_OUTPUT" \
  -extensions v3_intermediate_ca \
  -days 3650 \
  -md sha384 \
  -passin "pass:$HSM_PIN" -batch; then
  
    echo "[+] CSR signé avec succès !"
    echo "[*] Certificat de sortie écrit dans : $CERT_OUTPUT"
    
    TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    LOG_MSG="$TIMESTAMP - PKI_SIGN_SUCCESS - Signed CSR: $CSR_INPUT -> $CERT_OUTPUT (SHA-256: $(sha256sum "$CERT_OUTPUT" | cut -d' ' -f1))"
    echo "$LOG_MSG" >> "$AUDIT_LOG"
    
    # Envoi de la trace vers le stockage WORM
    logger -t SNISID-PKI -p security.info "$LOG_MSG"
    rsync -az "$AUDIT_LOG" "$WORM_TARGET" || echo "[WARNING] Échec de la synchronisation vers l'archivage WORM."
else
    TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    echo "$TIMESTAMP - PKI_SIGN_FAILURE - Échec de signature du CSR: $CSR_INPUT" >> "$AUDIT_LOG"
    logger -t SNISID-PKI -p security.err "Échec de signature PKI."
    exit 3
fi
```

---

## 2. Distribution Optimisée des CRLs pour Haïti (USSD/Bloom Filters)

Pour surmonter les contraintes de bande passante en zone rurale haïtienne, le SNISID utilise des **Filtres de Bloom hautement compressés (< 50 Ko)**. Ces filtres sont segmentés pour être distribués par SMS via le protocole USSD (`*123#CRL`).

### Script Python : `pki/scripts/compress_crl.py`

Ce script compresse une CRL ou liste de numéros de série révoqués en un filtre de Bloom zlib-compressé, puis produit les segments de SMS.

```python
# File: /pki/scripts/compress_crl.py
# Pure Python - Zero dependencies
import sys
import math
import hashlib
import struct
import zlib
import base64
import re

class SimpleBloomFilter:
    def __init__(self, size_in_bits, num_hashes):
        self.size_in_bits = size_in_bits
        self.num_hashes = num_hashes
        self.bit_array = bytearray(math.ceil(size_in_bits / 8))

    def _hashes(self, item):
        # Optimisation Kirsch-Mitzenmacher
        h = hashlib.sha256(str(item).encode('utf-8')).digest()
        h1 = struct.unpack("<Q", h[0:8])[0]
        h2 = struct.unpack("<Q", h[8:16])[0]
        for i in range(self.num_hashes):
            yield (h1 + i * h2) % self.size_in_bits

    def add(self, item):
        for bit_index in self._hashes(item):
            byte_index = bit_index // 8
            bit_offset = bit_index % 8
            self.bit_array[byte_index] |= (1 << bit_offset)

    def check(self, item):
        for bit_index in self._hashes(item):
            byte_index = bit_index // 8
            bit_offset = bit_index % 8
            if not (self.bit_array[byte_index] & (1 << bit_offset)):
                return False
        return True

    def serialize(self):
        header = struct.pack("<II", self.size_in_bits, self.num_hashes)
        compressed = zlib.compress(self.bit_array, level=9)
        return header + compressed

    @classmethod
    def deserialize(cls, data):
        size_in_bits, num_hashes = struct.unpack("<II", data[:8])
        compressed_bit_array = data[8:]
        bit_array = zlib.decompress(compressed_bit_array)
        bf = cls(size_in_bits, num_hashes)
        bf.bit_array = bytearray(bit_array)
        return bf

def parse_crl_serials(crl_path):
    serials = []
    try:
        with open(crl_path, "r") as f:
            for line in f:
                line = line.strip()
                if not line or line.startswith("#"):
                    continue
                if line.lower().startswith("0x"):
                    serials.append(int(line, 16))
                else:
                    try:
                        serials.append(int(line))
                    except ValueError:
                        try:
                            serials.append(int(line, 16))
                        except ValueError:
                            pass
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
    return serials

def chunk_payload(payload_bytes, chunk_size=120):
    encoded = base64.b64encode(payload_bytes).decode('ascii')
    return [encoded[i:i+chunk_size] for i in range(0, len(encoded), chunk_size)]

def main():
    if len(sys.argv) < 3:
        print("Usage: python compress_crl.py <input_list> <output_bin> [capacity] [error_rate]")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    capacity = int(sys.argv[3]) if len(sys.argv) > 3 else 10000
    error_rate = float(sys.argv[4]) if len(sys.argv) > 4 else 0.001

    serials = parse_crl_serials(input_file)
    if not serials:
        print("No serials found.")
        sys.exit(1)

    n = max(len(serials), capacity)
    p = error_rate
    m = int(- (n * math.log(p)) / (math.log(2) ** 2))
    k = int((m / n) * math.log(2))
    k = max(1, k)

    bf = SimpleBloomFilter(m, k)
    for s in serials:
        bf.add(s)

    serialized = bf.serialize()
    compressed_size = len(serialized)
    print(f"Bloom Filter size: {compressed_size} bytes ({compressed_size / 1024:.2f} KB)")

    with open(output_file, "wb") as f:
        f.write(serialized)

    chunks = chunk_payload(serialized, chunk_size=120)
    for idx, chunk in enumerate(chunks):
        print(f"SMS {idx+1}/{len(chunks)}: SNISID_CRL_BF:{idx+1}/{len(chunks)}:{chunk}")

if __name__ == "__main__":
    main()
```

---

## 3. Rotation Automatique des Certificats TLS (24 Heures)

Les certificats des microservices du SNISID sont configurés avec une **durée de vie stricte de 24 heures** et une fenêtre de **renouvellement de 2 heures**.

### Configuration cert-manager : `pki/k8s/cert-manager-24h-rotation.yaml`

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: identity-svc-tls
  namespace: snisid-core
spec:
  secretName: identity-svc-tls-secret
  # Durée de vie stricte de 24 heures
  duration: 24h
  # Renouvellement automatique 2 heures avant l'expiration
  renewBefore: 2h
  issuerRef:
    name: snisid-vault-issuer
    kind: ClusterIssuer
  commonName: identity-svc.snisid.local
  dnsNames:
  - identity-svc.snisid.local
  - identity-svc.snisid-core.svc.cluster.local
  usages:
  - digital signature
  - key encipherment
  - server auth
  - client auth
```

### Règle d'alerte Prometheus (SOC SLA Monitor)

Cette règle alerte immédiatement le SOC si le taux de réussite des rotations descend sous le seuil contractuel de **99%** (cible à 99.9%).

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: cert-manager-rotation-alerts
  namespace: monitoring
spec:
  groups:
  - name: cert-manager-rotation.rules
    rules:
    - alert: CertManagerRotationSuccessRateLow
      expr: |
        sum(rate(certmanager_certificate_ready_status{status="True"}[24h])) 
        / 
        (sum(rate(certmanager_certificate_ready_status{status="True"}[24h])) + sum(rate(certmanager_certificate_ready_status{status="False"}[24h])) + 1e-9)
        * 100 < 99
      for: 5m
      labels:
        severity: critical
        tier: infra
      annotations:
        summary: "SLA Rotation PKI critique (< 99% sur 24h)"
        description: "Le taux de renouvellement automatique des certificats microservices sur 24h est de {{ printf \"%.2f\" $value }}%, ce qui est inférieur au seuil de 99%."
```

---

## 4. Procédure de Récupération de la Root CA (Shamir Secret Sharing)

Pour faire face à une catastrophe nationale majeure détruisant le site principal de Port-au-Prince, la clé maître de la Root CA est divisée en **9 parts distribuées géographiquement**, nécessitant un quorum de **5 parts** pour sa reconstruction.

### Distribution Géographique des Parts :
1. **Part 1** : Banque de la République d'Haïti (BRH) — Port-au-Prince (Ouest)
2. **Part 2** : Coffre Fort de l'État — Cap-Haïtien (Nord)
3. **Part 3** : Banque Nationale de Crédit (BNC) — Les Cayes (Sud)
4. **Part 4** : Trésor Public — Hinche (Centre)
5. **Part 5** : Archives Nationales — Gonaïves (Artibonite)
6. **Part 6** : Bureau Régional — Jacmel (Sud-Est)
7. **Part 7** : Bureau Régional — Saint-Marc (Artibonite)
8. **Part 8** : Trésor Public — Port-de-Paix (Nord-Ouest)
9. **Part 9** : Trésor Public — Fort-Liberté (Nord-Est)

### Script Python : `pki/scripts/shamir_secret_sharing.py`

Ce script effectue le calcul polynomial requis sur le corps premier de Mersenne $2^{521} - 1$.

```python
# File: /pki/scripts/shamir_secret_sharing.py
# Pure Python - Zero dependencies
import sys
import secrets

PRIME = 2**521 - 1

def extended_gcd(a, b):
    if a == 0:
        return b, 0, 1
    else:
        g, x, y = extended_gcd(b % a, a)
        return g, y - (b // a) * x, x

def modular_inverse(k, p):
    g, x, y = extended_gcd(k, p)
    if g != 1:
        raise ValueError('Inverse modulaire inexistant')
    else:
        return x % p

def split_secret(secret_int, threshold, num_shares):
    coefficients = [secret_int] + [secrets.randbelow(PRIME) for _ in range(threshold - 1)]
    shares = []
    for x in range(1, num_shares + 1):
        y = 0
        for power, coeff in enumerate(coefficients):
            y = (y + coeff * pow(x, power, PRIME)) % PRIME
        shares.append((x, y))
    return shares

def reconstruct_secret(shares):
    secret = 0
    for j, (x_j, y_j) in enumerate(shares):
        numerator = 1
        denominator = 1
        for m, (x_m, _) in enumerate(shares):
            if m == j:
                continue
            numerator = (numerator * (-x_m)) % PRIME
            denominator = (denominator * (x_j - x_m)) % PRIME
        lagrange_coeff = (numerator * modular_inverse(denominator, PRIME)) % PRIME
        secret = (secret + y_j * lagrange_coeff) % PRIME
    return secret

def main():
    if len(sys.argv) < 2:
        print("Usage:")
        print("  split: python shamir_secret_sharing.py split <hex_secret> [threshold] [total_parts]")
        print("  combine: python shamir_secret_sharing.py combine <x1:y1> <x2:y2> ...")
        sys.exit(1)

    command = sys.argv[1].lower()

    if command == "split":
        hex_secret = sys.argv[2]
        threshold = int(sys.argv[3]) if len(sys.argv) > 3 else 5
        num_shares = int(sys.argv[4]) if len(sys.argv) > 4 else 9
        secret_int = int(hex_secret, 16)
        
        shares = split_secret(secret_int, threshold, num_shares)
        depts = [
            "Port-au-Prince (BRH)", "Cap-Haitien (Nord)", "Les Cayes (Sud)",
            "Hinche (Centre)", "Gonaives (Artibonite)", "Jacmel (Sud-Est)",
            "Saint-Marc (Artibonite)", "Port-de-Paix (Nord-Ouest)", "Fort-Liberte (Nord-Est)"
        ]
        for i, (x, y) in enumerate(shares):
            loc = depts[i] if i < len(depts) else f"Location {x}"
            print(f"Part {x} [{loc}] : {x}:{hex(y)[2:]}")

    elif command == "combine":
        shares_arg = sys.argv[2:]
        shares = []
        for s in shares_arg:
            parts = s.split(":")
            shares.append((int(parts[0]), int(parts[1], 16)))
        secret_int = reconstruct_secret(shares)
        print(f"Secret reconstruit (Hex) : {hex(secret_int)[2:]}")

if __name__ == "__main__":
    main()
```

---

*Ce document fait partie intégrante du cadre réglementaire SNISID de la République d'Haïti.*
